# Theming Klarity

Klarity supports theming your documentation site through the `[visual.vars]` section in your `klarity.toml` file.

---

## 🎨 Customizing Colors

You can override the default colors used by Klarity by specifying CSS variables in the `[visual.vars]` section.  
These variables control backgrounds, borders, accents, and text colors.

**Available variables:**
- `bg_main`
- `bg_panel`
- `bg_hover`
- `bg_active`
- `border_soft`
- `border_hard`
- `accent_primary`
- `accent_secondary`
- `text_main`
- `text_dim`

You can set any or all of these. Unset variables will use Klarity's defaults.

Catppuccin/Dracula like example:

```toml
[visual.vars]
bg_main = "#181825"
bg_panel = "#1e1e2e"
bg_hover = "#313244"
bg_active = "#45475a"
border_soft = "#585b70"
border_hard = "#6c7086"
accent_primary = "#f5c2e7"
accent_secondary = "#b4befe"
text_main = "#cdd6f4"
text_dim = "#a6adc8"
```

> [!TIP]
> You don't have to use hex colors here, any valid css color syntax will work.

---

## 🖌️ Code Highlighting Theme

The `[visual] theme` option controls code block highlighting.  
See the [[Config.md|config]] page for details and available themes.

---

## 🧩 Custom CSS (coming soon)

A future version of Klarity will allow you to supply your own CSS file for full customization.  
Stay tuned for updates!

