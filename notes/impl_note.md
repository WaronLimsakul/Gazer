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


### High-level Architecture

2 main components: `engine` and `ui`
- `ui` (`ui.Draw`) uses `engine.State.Notifier` channel to notify the engine to do something
- `engine` (`engine.start`) uses `app.Window.Invalidate` to trigger redraw 

### DOM rendering
- Somehow, `p` inside `h1` cannot override `h1` while `h2` can. So I don't want this confusing heirachy. I will just separate
