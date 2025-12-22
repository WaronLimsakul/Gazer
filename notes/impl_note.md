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


