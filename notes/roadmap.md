## List of what I want to do

### To-do Lists
- [x] Fix `<!...>` comment style. what if `<!-- <h1>Hello, world</h1> -->`
- [x] Update default font
- [x] Has text rendering custom component
- [x] `<h1><i>hello</i></h1>` has to be big and italic
- [x] Lay the texts out vertically using `List` because
    1. It only lays what is visible
    2. Scrollable default
- [x] Use flex for inline text e.g. `<p>hello <i>world</i></p>` should be same line
- [x] Tabs inside state is a UI component. It shouldn't be.
- [x] Learn the MVC arch
- [ ] Entire theme set
- [x] Some metadata support
  - [x] title
  - [x] favicon
- [x] Set up supporting CSS field in components
- [x] Have CSS parser `StyleParser`?: parse CSS string into `Style` struct
- [ ] `DomRenderer` should have `StyleParser`
While `RenderDOM`: NOTE, source priority if conflict: inline `style` attr > `style` tag > external css file
- [ ] If saw`stylesheet`, fetch it, parse and store in the function scope (might consider cache it)[ ] [ ] 
- [ ] If saw `style` tag, parse it and merge with the old one if exists .
Call it "global style" for now (in this function term)
- [ ] Start rendering like normal. But modify the component according to the global style
- [ ] If saw `style` attr during rendering, apply the global style first, then parse this one and apply after to override (higher priority)

### HTML tags support
- [x] A 
- [x] Img 
- [x] Ul, LI 
- [x] Ol
- [ ] B (or Strong) 
- [x] I (or Em) 
- [x] Hr 
- [x] Div 
- [x] Span
- [x] Section
- [x] Button
- [ ] Input
  - [x] type text
  - [x] type password
  - [x] type number
  - [x] type email
  - [ ] type checkbox
  - [ ] type radio
  - [ ] type date
  - [ ] type submit
- [ ] Table, Tr, Td, Th

### CSS Support
[src](https://www.w3schools.com/html/html_css.asp)
- [ ] CSS front-end processor
  - [ ] Inline: using `style` attr
  - [ ] Internal: using `<style>` elements in `<head>`
  - [ ] External: using `<link>` element
- [ ] CSS rendering process:
  - [ ] Text Size: `font-size`
  - [ ] Text Color `color`
  - [ ] Element margin size `margin top right bot left`: `margin-left` `margin-right` `margin-bottom` `margin-top`
  - [ ] Element border size : `border-width`, `border-color`, `border-radius`
    - [ ] `border-style` 
    - [ ] border shorthand
  - [ ] Element padding size `padding`
  - [ ] Flex model 

### JavaScript support
- [ ] Light JavaScript support

### Other UI features
- [x] Search bar's search button 
- [ ] Search bar's suggestion
- [x] Tab system
- [ ] GioUI normal window is super ugly. Turn-off window decoration and handroll the window.
  - [ ] Wait, I think we can just `Decorate` it. Oh, it's the same way.
- [ ] Custom theme system for Gazer
- [ ] Tab tooltip
- [ ] Close tab button
- [ ] Keybinding for manipulating tab
- [ ] CSS support structure
