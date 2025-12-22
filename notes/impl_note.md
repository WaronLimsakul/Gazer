## Implementation Note

### Gio
I kinda want low-level control over what I am drawing, so I chose Gio, which give
me low-level control but also some widget to play with.

### Testing 
I just realize don't there is `t.Run()` which makes subtests very clean and easy

### Redesign `parser.Node`
Initially, design a DOM node to have a field `Inner string` as inner content but then I realize that something like
```html
<p>Hello <br> world<p>
```
is also valid and would break my implementation (I mean... how can you determine where the `<br>` is with just one `Inner` and a `Children` field).
Therefore, I'd redesign it to not have `Inner` anymore, and move them to a new `Node` type called `Text`.
Now I can just do `children{Text{"Hello "}, Br{}, Text{" world"}}`


