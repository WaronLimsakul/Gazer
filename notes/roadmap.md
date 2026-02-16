## List of what I want to do

### To-do Lists
- [ ] Support table element
- [ ] `info.cern.ch`, `berkshirehathaway.com` <- pure html and CSS. Can use for testing
- [ ] Have a global error logging system
    - [ ] Bubbling all errors up to log at the highest level.
- [ ] Check bold/italic inside bullet point rendenring
- [x] Support SVG file
- [ ] Button component <- wrap gioui button, so we can inactivate
- [ ] Tab navigation like window navigation
- [x] Image size setting (`width`, `height` attribute)
    - Supports, "500", "500px", "500rem" kind of input
- [ ] Have user input dispatcher, so I can do `range input` and `switch-case` the input type




### HTML tags support
- [x] A 
- [x] Img 
- [x] Ul, LI 
- [x] Ol
- [x] B (or Strong) 
- [x] I (or Em) 
- [x] Hr 
- [x] Div 
- [x] Span
- [x] Section
- [x] Button
- [x] Input
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
- [x] CSS front-end processor
  - [x] Inline: using `style` attr
  - [x] Internal: using `<style>` elements in `<head>`
  - [x] External: using `<link>` element
- [x] CSS rendering process:
  - [x] Text Size: `font-size`
  - [x] Text Color `color`
  - [x] Element margin size `margin top right bot left`: `margin-left` `margin-right` `margin-bottom` `margin-top`
  - [x] Element border size : `border-width`, `border-color`, `border-radius`
    - [x] `border-style` 
    - [x] border shorthand
  - [x] Element padding size `padding`
  - [ ] Flex model 
- [x] Comments
- [ ] At-rule

### JavaScript support
- [ ] Light JavaScript support
- [ ] DOM Manipulation
    - document.getElementById, querySelector, querySelectorAll
    - element.innerHTML, element.textContent
    - element.setAttribute, getAttribute, removeAttribute
    - element.style (basic property setting)
    - element.appendChild, removeChild, createElement
    - element.classList.add/remove/toggle/contains
- [ ] Events
    - element.addEventListener, removeEventListener
    - Basic event object (event.target, event.preventDefault, event.stopPropagation)
    - Mouse events: click, mouseover, mouseout
    - Keyboard events: keydown, keyup
    - window.onload / DOMContentLoaded
- [ ] Timers
    - setTimeout, clearTimeout
    - setInterval, clearInterval


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
