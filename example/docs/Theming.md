# Theming Klarity

Klarity supports theming your documentation site through the `[visual.vars]` section in your `klarity.toml` file.

---

## 🎨 Customizing Colors

You can override the default colors used by Klarity by specifying CSS variables in the `[visual.vars]` section.  
These variables control backgrounds, borders, accents, and text colors.

**Available variables:**
- `bg_main`: the background behind the main content of the page
- `bg_panel`: background of things like the sidebar
- `bg_hover`: background of things being hovered
- `bg_active`: mostly unused, *may be removed in the future*
- `border_soft`: borders like the codeblock border
- `border_hard`: mostly used on scrollbars and small underlines
- `accent_primary`: the main contrast color
- `accent_secondary`: only really appears in blockquotes
- `text_main`: main text color, despite this most common text uses `text_dim`
- `text_dim`: the color for most of the content text

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

## 🧩 Custom CSS

Klarity also supports providing your own custom CSS file, which gives you complete control over the styling of the site.

To use this feature simply place your css file anywhere in the project and set the following in your config:

```toml
visual.custom_css = "custom.css"
```

### What are the drawbacks ?

Well Klarity ships with a default stylesheet containing ~700 lines of CSS, you view it [here](https://github.com/kociumba/klarity/blob/main/assets/style.css), due to how Klarity works, there's no elaborate theming API, so creating your own theme means working with raw CSS and overriding existing styles directly.

### Advice for theming with custom CSS

If you are an experienced CSS veteran, go ahead and create your own themes without listening to my advice. But for anyone less experienced with this kind of customization:

- Use dev tools extensively to see how the default styles apply and what which rules are taking precedence.
- Some selectors Klarity uses are pretty specific, e. g. `.custom-block.danger[data-callout-type="github-style"]`, to override them effectively, you may need to match or exceed their specificity.
- Make focused changes, large-scale overrides across many elements are more likely to introduce subtle breakage.
- Some global rules like:
    ```css
    *, *::before, *::after {
        box-sizing: border-box;
    }
    ```
    are core to the layout and should not be reset or overridden in themes.

In general, if you're aiming for a full "overhaul" theme, I recommend starting from Klarity's default [style.css](https://github.com/kociumba/klarity/blob/main/assets/style.css) and modifying it directly. Just be aware: this stylesheet may change without much notice, which could break compatibility with your custom theme.
