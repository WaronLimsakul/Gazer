# Gazer

A lightweight browser engine written in Go using the Gio UI framework.
It parses HTML and CSS, builds a DOM tree, and renders it using an immediate mode UI approach.

---

## What It Does

- Fetches and parses HTML
- Applies CSS from inline, internal, and external stylesheets
- Builds a DOM tree with proper text nodes
- Renders layout, typography, lists, images, links, and basic form elements
- Supports GIF rendering with timed frame updates

---

## How to Run

1. Install Go
2. Clone the repository
  ```
    git clone https://github.com/WaronLimsakul/Gazer
  ```
3. Build
  ```
   go build
  ```

You can then enter a URL such as info.cern.ch to test rendering.

---

## How to Contribute

1. Fork the repository
2. Create a new branch
3. Make your changes
4. Submit a pull request

Areas that need work:
- Flexbox layout system
- Table rendering
- Basic JavaScript support
- Rendering improvements and performance optimization
