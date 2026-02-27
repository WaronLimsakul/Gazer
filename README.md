# Gazer: The "I Can't Believe It's Not V8" Browser Engine

Welcome to **Gazer**, a low-level browser engine built in Go using the Gio UI framework. While modern browsers consume RAM like a competitive eater at a buffet, Gazer is handcrafted to render the web one state-machine transition at a time. It’s small, it’s fast-ish, and it’s mostly pointer-safe.

---

## Technical Overview: The State Machine

At the core of our lexer is a Finite State Machine (FSM). It treats HTML strings like a suspicious obstacle course, moving from a "Void" state into specific paths for tags, comments, and doctypes.

```json
{
  "lexer_name": "HTML_Tag_Lexer",
  "initial_state": "Void",
  "final_state": "return",
  "description": "A state machine for identifying HTML-style tags, comments, and doctypes."
}

```

### High-Level Architecture

Gazer follows an MVC-adjacent pattern split into two main packages:

* **The Engine:** The "brains" that handle state, tab logic, and network fetching.
* **The Renderer:** The "brawn" that takes a DOM tree and forces Gio to draw it on your screen.

---

## Engineering Philosophies (The "Hard Way" Lessons)

### 1. Style Inheritance via Accumulative Recursion

Originally, we tried a decorator pattern where styles flowed up. That was a mistake. We now use **accumulative recursion** where the parent’s style is passed down the tree.

* **The H1 Test:** If you put an `<h1>` inside an `<h2>`, Gazer now correctly realizes that the universe shouldn't implode and renders the inner tag's style.
* **Rendering Context:** We pass a `RenderingContext` stack rather than just a style struct. This allows a `<li>` to look up the stack, see its `<ul>` ancestor, and realize it needs a bullet point.

### 2. The DOM and the Text Node Revolution

We realized that `Inner string` fields are where dreams go to die. If you have `<p>Hello <br> world</p>`, a single string field can't handle the line break. We transitioned to a pure node-based system where everything—even "Hello "—is its own `Text` node.

### 3. Gio Power User Moves

* **op.Record:** This is our secret weapon. We record operations to a macro, measure the resulting dimensions, and then retrospectively draw backgrounds. It’s basically time travel for UI drawing.
* **Immediate Mode Gifs:** Rendering Gifs in an immediate-mode UI is awkward. We precompute frame time bounds and use `op.InvalidateCmd` to tell the GPU exactly when to wake up for the next frame.

### 4. Reflection and "Pointer-o-phobia"

* **Reflect Package:** We use Go's `reflect` to loop over CSS `Style` structs. It’s significantly cleaner than writing 500 `if` statements for every possible CSS property.
* **Pass-by-Value:** We have a strict "no pointer" policy unless the struct is "huge," needs mutation, or is optional. If you see a pointer in our rendering logic, it’s there for a very good reason.

---

## Project Status

### Supported HTML Elements

| Category | Supported |
| --- | --- |
| **Layout** | `div`, `span`, `section`, `hr` |
| **Typography** | `b`, `strong`, `i`, `em`, `h1-h6` |
| **Lists** | `ul`, `ol`, `li` (with auto-incrementing counts) |
| **Media** | `img` (including Gif and SVG), `a` (links) |
| **Forms** | `button`, `input` (text, password, email, number) |

### CSS Implementation

* **Sources:** Supports Inline `style` attributes, Internal `<style>` tags, and External `<link>` stylesheets.
* **Box Model:** Basic support for `margin`, `padding`, and `border` (including shorthand and radius).
* **Typography:** `font-size` and `color` inheritance is fully functional.

---

## The Roadmap (Or: The "Why Did I Start This" List)

* [ ] **Flexbox:** Currently, we use a mix of `List` for rows and `Flex` for inline elements. Implementing a full Flexbox model is the final boss.
* [ ] **Table Support:** Because 1990s web design deserves to be rendered correctly.
* [ ] **JavaScript (The Light Version):** Support for basic DOM manipulation (`getElementById`) and timers (`setTimeout`).
* [ ] **Custom Window Decorations:** Because the default Gio window is, let's be honest, a bit plain.

---

## Getting Started

1. **Clone the repo:** `git clone https://github.com/WaronLimsakul/Gazer`
2. **Build:** `go build`
3. **Run:** `./Gazer`
4. **Test:** Point it at `info.cern.ch` and witness the glorious dawn of the internet.
