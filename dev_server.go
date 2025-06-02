package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/alecthomas/kong"
	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/websocket"
)

type wsHub struct {
	mu      sync.Mutex
	clients map[*websocket.Conn]struct{}
}

func newWsHub() *wsHub {
	return &wsHub{
		clients: make(map[*websocket.Conn]struct{}),
	}
}

func (h *wsHub) add(conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[conn] = struct{}{}
}

func (h *wsHub) remove(conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.clients, conn)
}

func (h *wsHub) broadcast(msg string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for c := range h.clients {
		c.SetWriteDeadline(time.Now().Add(2 * time.Second))
		if err := c.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
			c.Close()
			delete(h.clients, c)
		}
	}
}

var dev_server = false

func (d *DevServer) Run(ctx *kong.Context) error {
	dev_server = true
	pwd = d.Path
	InitMarkdown(pwd)
	projectPath, err := filepath.Abs(d.Path)
	if err != nil {
		return fmt.Errorf("failed to resolve path: %w", err)
	}

	if err := buildKlarity(projectPath); err != nil {
		return fmt.Errorf("initial build failed: %w", err)
	}

	cfg := ReadConfig(projectPath)
	outputDir, err := filepath.Abs(filepath.Join(projectPath, cfg.Output_dir))
	if err != nil {
		return fmt.Errorf("failed to get output dir: %w", err)
	}

	hub := newWsHub()

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("watcher error: %w", err)
	}
	defer watcher.Close()

	watchDirs := []string{projectPath}
	for _, dir := range cfg.Doc_dirs {
		watchDirs = append(watchDirs, filepath.Join(projectPath, dir))
	}
	for _, dir := range watchDirs {
		filepath.Walk(dir, func(path string, info os.FileInfo, _ error) error {
			if info == nil {
				return nil
			}
			if info.IsDir() {
				watcher.Add(path)
			}
			return nil
		})
	}

	var debounceTimer *time.Timer
	var mu sync.Mutex
	triggerRebuild := func() {
		mu.Lock()
		defer mu.Unlock()
		if debounceTimer != nil {
			debounceTimer.Stop()
		}
		debounceTimer = time.AfterFunc(400*time.Millisecond, func() {
			fmt.Println("[Klarity] Change detected, rebuilding...")
			if err := buildKlarity(projectPath); err != nil {
				fmt.Printf("[Klarity] Rebuild error: %v\n", err)
			} else {
				hub.broadcast("reload")
			}
		})
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		p := r.URL.Path
		if p == "/" {
			p = "/index.html"
		}
		file := filepath.Join(outputDir, filepath.Clean(p))
		if _, err := os.Stat(file); os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}
		if strings.HasSuffix(file, ".html") {
			raw, err := os.ReadFile(file)
			if err != nil {
				http.Error(w, "Internal server error", 500)
				return
			}
			mod := injectLiveReload(string(raw))
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(mod))
			return
		}
		http.ServeFile(w, r, file)
	})

	http.HandleFunc("/klarity-livereload", func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("ws upgrade:", err)
			return
		}
		hub.add(conn)
		defer hub.remove(conn)
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				break
			}
		}
	})

	srv := &http.Server{Addr: "localhost:5173"}
	if cfg.Dev.Port != 0 && cfg.Dev.Port > 1024 && cfg.Dev.Port < 49151 {
		srv.Addr = fmt.Sprintf("localhost:%d", cfg.Dev.Port)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	done := make(chan struct{})

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				isMd := strings.HasSuffix(event.Name, ".md")
				isToml := strings.HasSuffix(event.Name, ".toml")

				if isMd || isToml {
					if event.Op&fsnotify.Create == fsnotify.Create {
						fmt.Println("[Watcher] Detected new file:", event.Name)
						triggerRebuild()
					}
					if event.Op&fsnotify.Write == fsnotify.Write {
						fmt.Println("[Watcher] Detected modification:", event.Name)
						triggerRebuild()
					}
					if event.Op&fsnotify.Remove == fsnotify.Remove {
						fmt.Println("[Watcher] Detected deletion:", event.Name)
						triggerRebuild()
					}
				}

				// TODO: add a check for cfg.Doc_dirs
				if event.Op&fsnotify.Create == fsnotify.Create {
					info, err := os.Stat(event.Name)
					if err == nil && info.IsDir() && slices.Contains(cfg.Doc_dirs, info.Name()) {
						watcher.Add(event.Name)
						filepath.Walk(event.Name, func(path string, info os.FileInfo, _ error) error {
							if info != nil && info.IsDir() {
								watcher.Add(path)
							}
							return nil
						})
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("[Klarity] Watcher error:", err)
			case <-done:
				return
			}
		}
	}()

	go func() {
		fmt.Printf("[Klarity] Dev server running on http://%s\n", srv.Addr)
		fmt.Println("[Klarity] Watching for changes...")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[Klarity] HTTP server error: %v", err)
		}
	}()

	<-quit
	fmt.Println("[Klarity] Shutting down...")

	// Stop file watcher goroutine
	close(done)

	// Graceful HTTP shutdown
	ctxTimeout, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctxTimeout); err != nil {
		log.Printf("[Klarity] HTTP server shutdown error: %v", err)
	}

	fmt.Println("[Klarity] Goodbye!")
	return nil
}

func injectLiveReload(html string) string {
	script := `<script>
(function() {
	var ws = new WebSocket((location.protocol === 'https:' ? 'wss://' : 'ws://') + location.host + '/klarity-livereload');
	ws.onmessage = function(event) {
		if (event.data === 'reload') location.reload();
	};
})();
</script>`
	if strings.Contains(html, "</body>") {
		return strings.Replace(html, "</body>", script+"</body>", 1)
	}
	return html + script
}
