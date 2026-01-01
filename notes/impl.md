## Implementation Note

### Gio
I kinda want low-level control over what I am drawing, so I chose Gio, which give
me low-level control but also some widgets to play with.

### Testing 
I just realize don't there is `t.Run()` which makes subtests very clean and easy

### Redesign `parser.Node`
Initially, design a DOM node to have a field `Inner string` as inner content but then I realize that something like
```html
<p>Hello <br> world<p>
```
is also valid and would break my implementation (I mean... how can you determine where the `<br>` is with just one `Inner` and a `Children` field). 
Therefore, I'd add new node type called `Text` and use it instead of `Inner`. However, I'd reuse `Inner` field as a field only for `Text` to holds its content. Now I can just do `children{Text{"Hello "}, Br{}, Text{" world"}}`

### Reimplement `lexer.GetNextToken`
As you know, lexer/parser is just a state machine. But no one tells me how to implement it.
So in the `for` loop, I `switch` the current character first and then `switch` the current state inside,
It works but so hard to read. So I flipped the order to reflect the state machine, easier to read now.

### High-level Architecture

2 main components (package): `engine` and `renderer`
- `renderer` (`renderer.Draw`) uses `engine.State.Notifier` channel to notify the engine to do something
- `engine` (`engine.start`) uses `app.Window.Invalidate` to trigger redraw 

### DOM rendering: tag heirachy
Somehow, `p` inside `h1` cannot override `h1` while `h2` can. So I don't want this confusing heirachy. I will just separate

### DOM rendering: `DomRenderer`
Gio need me to hold the state of "selectable" text separate from the material (nothing unexpected), so I implementd
`DomRenderer` struct (1 of this per 1 website) so that it can hold any widget state related to that website

### DOM label style inheritance
A node like this `H1 -> i -> Text` means we want the render text to be big and italic.
Therefore, we need style inheritance system. Decorator pattern seems like a good idea.
But I don't want to make it super OOP, so I modified it a bit just be `LabelFunc` and `LabelDecorator` 
type in `ui` package. 

- `type LabelFunc = func(*Theme, *Selectable, string) LabelStyle` is a base text
- `type LabelDecorator = func(*Theme, LabelStyle) LabelStyle` is a decorated text

Now we can just do `H1(thm, I(thm, Text(thm, "hello")))` for the example node (easy for recursion).

### Use `List` to layout element instead of `Flex`
At first, I only use how to use `Flex` container but I see `List` is a better choice because
1. It lays out only what is visible (don't have to waste time computing what we can't see)
2. It supports scrolling

However, I will try to use `Flex` when rendering text inline e.g. `<p>hello <i>world</i></p>` should be same line.

Therefore, now `DomRenderer.render` return `[][]Element` (`Element` is `interface {Layout(gtx C) D}`)
outer layer is rendering row, inner layer is each element in a row from left to right.

We use `List` to render the row and `Flex` to render each one from left to right.

