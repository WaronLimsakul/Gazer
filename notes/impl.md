## Implementation Note

### Gio
I kinda want low-level control over what I am drawing, so I chose Gio, which give
me low-level control but also some widgets to play with.

### Testing 
I just realize don't there is `t.Run()` which makes subtests very clean and easy

### `layout.Context`
I just realize that I have `gtx.Constraints.Constrain()`, easy way to sneak custom drawing in.

### `op.Record(*op.Ops)`
This one is super OP. I can do `macro := op.Record(gtx.Ops)` then do whatever I want with the `gtx` and all
the operations will be recorded to `macro`. Then I can do `savedOp := macro.Stop()` to stop recording and have all
the operations saved into `savedOp`. I use it to check the size of the content of component then draw a background with that size.

### Render component interaction
I used to let the component expose the `.Update(gtx)` to let the caller use when got the `gtx`, then handle it separately.
but I just realize we can just put that in `.Layout(gtx)`, basically let the component update itself from the previous frame.

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

#### Update 1: MVC
I do further research on what other people have done for architecture. I found that there is something call MVC (Model-View-Controller)
which I like the idea of letting the Model (business logic processor unit) handles the app main state. Other than that
my `renderer.Draw` should be View and Controller since I notify the engine directly.

### DOM rendering: tag heirachy
Somehow, `p` inside `h1` cannot override `h1` while `h2` can. So I don't want this confusing heirachy. I will just separate

### DOM rendering: `DomRenderer`
Gio need me to hold the state of "selectable" text separate from the material (nothing unexpected), so I implementd
`DomRenderer` struct (1 of this per 1 website) so that it can hold any widget state related to that website

#### Update 1: `clickables`
Now I also have `linkClickables` for storing the state of the clickable anchor inside `DomRenderer`

#### Update 2: `domRenderers`
After I introduced a tab system, I am afraid that the ui state that `DomRenderer` holds will fighting each other
if I use one dom renderer to render all the tabs I have. So I have `domRenderers := map[*ui.Tab]*DomRenderer` to store it.
I know making a struct that hold both of them looks better, but I feel like that's overkill.

#### Update 3: memory leaks concern
The `selectables map[*Node]*Selectable` is easy to use but what if the node got cleaned? will it segfault?
I came to realize that the node will not be cleaned since it map itself refer to it. But this will still be a (lowkey) memory leak until
the dom renderer itself be freed.

### DOM label style inheritance
A node like this `H1 -> i -> Text` means we want the render text to be big and italic.
Therefore, we need style inheritance system. Decorator pattern seems like a good idea.
But I don't want to make it super OOP, so I modified it a bit just be `LabelFunc` and `LabelDecorator` 
type in `ui` package. 

- `type LabelFunc = func(*Theme, *Selectable, string) LabelStyle` is a base text
- `type LabelDecorator = func(*Theme, LabelStyle) LabelStyle` is a decorated text

Now we can just do `H1(thm, I(thm, Text(thm, "hello")))` for the example node (easy for recursion).

#### Update 1: type check
I just realize that in case that there is other type of node that is not text node as a child 
of the text node, I have to call the `renderNode` inside `renderTextNode` anyway. So I might as well
just let `renderTextNode` return `[][]Element`. Now it's like a lowkey mutual recursion.

#### Update 2: prefix
I used to just sneak "• " in the `Label` content. But sine I introduced button, I can't do that, otherwise the button 
border will takes the bullet point out of it as well.

### Use `List` to layout element instead of `Flex`
At first, I only use how to use `Flex` container but I see `List` is a better choice because
1. It lays out only what is visible (don't have to waste time computing what we can't see)
2. It supports scrolling

However, I will try to use `Flex` when rendering text inline e.g. `<p>hello <i>world</i></p>` should be same line.

Therefore, now `DomRenderer.render` return `[][]Element` (`Element` is `interface {Layout(gtx C) D}`)
outer layer is rendering row, inner layer is each element in a row from left to right.

We use `List` to render the row and `Flex` to render each one from left to right.

### For rendering inheritance
Sometimes, renderer has to know the relationship between node to render correct (e.g. `li` will have
bullet point if and only if it is under `ul`). I solve it by let the custom `Label` widget has `tags map[Tag]bool`
to record the inhertance of the renderee (because rendering flow down, but the widgets flow up).

### Tabs system
This is the part that I think my architecture design is at lowest part. I don't know how to separate
between the business logic and rendering anymore. When I create `ui.Tabs`, I unintentionally put things like
`IsSelected`, `IsLoading` and DOM `root` in the `Tab`. That makes me put `tabs` inside the `engine.State`
and now it becomes that `State` also holds the ui component. I have to solve this

#### Update 1: separate `ui.Tab` and `engine.Tab`
I want the state to only hold the business data. So I pull `url` and `root` out of `ui.Tab` and create `engine.Tab` to hold them.
Now I hold `Tabs []*engine.Tab` in the state and let the `engine.Start` handle it.

One thing I realize is that `ui.Tab` rendering system doesn't depend on those fields I pulled out at all. Good call.
However, one problem is that I still have to keep the `state.Tabs` and `tabs` inside `renderer.Draw` in sync (e.g. same amount, and same order).

### Theme UI system
2 problems
1. The palette that `material.Theme` gives me are not enough, I'm thinking of having `GazerTheme`
2. If I want my component to be more flexible, I might change the universal `Layout(gtx)` into `Layout(thm, gtx)`
    Or hybrid (keep the original but talso)

### `strings.Builder` struct
When implement `String() string` method, I did lots of `res += something` (e.g. in a loop), lsp said it's inefficient, I agree,
because string in Go is just immutable []char reallocate repeatedly is not a good idea. Thus, I search up some soln and found `string.Builder` struct.

`Builder` have it's own resizable buffer, we can do `.WriteString()` or `.Fprintf()` into its address to append whatever we want.
Then at the end, I do `.String()` to get the result. More efficient.

### CSS parsing
I don't wanna really touch the `parser` and `lexer` package so I just add `css` package that do both lexing and parsing css strings.

My idea right now is to represent entire css document in a struct called `StyleSet` which is something like this.

```go
// StyleSet is a (almost) ready-to-use style set of one CSS file (or more?)
type StyleSet struct {
	universal   *Style
	idStyles    map[string]*Style
	classStyles map[string]*Style
	tagStyles   map[parser.Tag]*Style
}
```

I merge all possible styling you can do with an element into `Style` struct. I know it's not efficient, but it's easy to deal with.

```go
// Style is a property to style the rendering of any argument.
// The responsibility to intepret the struct is on caller.
type Style struct {
	color    *color.NRGBA
	bgColor  *color.NRGBA
	margin   *layout.Inset
    ...
}
```

There will be a `StyleParser` that can `Parse` string into `StyleSet`. I plan to let my `DomRenderer` (DR) hold a `StyleParser`.
I will let DR look for `<link rel="styleshee">` and/or `<style></style>` in the `<head>` to get a style set first before render `<body>`.

Then while traversing and rendering a dom-tree, my DR should always test the node against the style set to see if it spit out
any `Style` the  DR has to add when render the element. 

Oh, it also needs some inline parsing for `style` attribute.

#### Update 1: engine
I feel like the task of parsing something should be with engine, and also we have to fetch something, that convinces
me to let engine handle the fetch+parse CSS. For now, I'll let it stores the style in the `Styles` field of `Tab`.

#### Update 2: Dom Renderer
Since I store a `StyleSet` in the state.tab, dom renderer won't have access to it unless it has been told
in the function parameter (DR has `ui.Tab`, not `engine.Tab`). So I'm thinking of passing `*StyleSet` into `render` 
and let all those helper recursive function pass them around as well.

One good thing about this approach is if I found inline `style` that will affect the children of the node, I can just modify 
my current `style` and pass it to the recursion.

#### Update 3: Wait... Then my decoration and style goes different way?
Take a look at this:
```html
<h2>
    <h1>101</h1>
 </h2>
 <h1>
     <h2>102</h2>
 </h1>
```
V8 displays the first one as h1 and second one as h2. So I guess my way of decoration flowing up doesn't work.
Then, instead of using label decorator as `Label -> Label`, I'll do `Style -> Style`, thus my rendering process
will become accumulative recursion instead of a normal one. At the basecase, `Text()` should takes the style and 
produce a Label correctly. NOTE: I might rename `Text()` to be something like `Label()`.

#### Update 4: `<Ol>`...
A little special case is this guy. I have to find a way to pass the order count of the `<li>` inside an ol.

#### Update 5: I might be a genius
So my Label component doesn't only need `css.Style` to create a `Label`, it needs some more fields (e.g. some `clickable`, `seletable`).
So I defined `ui.LabelStyle`. Initially, I just embed `css.Style` into a struct, but then I realize I can pass `ui.LabelStyle` around 
in the `renderNode` function since we don't know what node are we rendering. Therefore, we need to pass around something that are uniform
and friendly for all types of node. So I separate `css.Style` with the extra fields I need in label



```go
// in ui/label.go
type LabelStyle struct {
	base  css.Style
	extra LabelExtraStyle
}

type LabelExtraStyle struct {
	tags      map[parser.Tag]bool
	clickable *widget.Clickable
	prefix    string
	count     *int // for <ol>
}

// in renderer/dom.go
// This is what we will use to pass around in RenderNode
type RenderingStyle struct {
	base  *css.Style
	label *ui.LabelExtraStyle
	// NOTE: can add more extra style for other type of node
}
```

Now when we traverse the tree and render, we can just check the node type, then build the style from `RenderingStyle` according to the node type.

### Pass by value? Wait, or should it be pointer?
I kinda let my intuition decide these kinda of questions, but to make things rigid. Here are some rules I will try to obey

#### Pass by value when
1. Type is small
2. I don't wanna mutate anything
3. Type is simple and immutable
4. Clean and simple stupid API

#### Pass by pointer when
1. Type is big
2. Type has something mutable (e.g. map, and slice), in that case, there is no point passing by value
3. I want to mutate something client pass in
4. I have to represent optionality

