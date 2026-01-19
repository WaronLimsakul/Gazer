## List of what I want to do

### To-do Lists
- [x] Support none-same-addr
- [x] Test all the css rendering from different places 

- [ ] Separate tab goroutine?
  - not right now, we render one by one anyway
- [ ] Optimize some more

- [x] Close tab button
- [x] Support `<header>`, `<footer>` 
- [x] Support `<main>`, `<article>`
- [ ] Support table element
- [x] CSS comment
- [x] Support container style support
- [ ] Support local files traversal
- [ ] Support pdf if possible
- [ ] Update serve.sh to observe change and reload
- [ ] Fix loading bar, I saw it firing like a Gatling gun when load a large page.
- [ ] Update serve.sh to observe the change and reload


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
