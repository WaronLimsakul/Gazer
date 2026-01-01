### Lexer `GetNextToken` state diagram

```mermaid
stateDiagram-v2
    [*] --> Void

    Void --> Comment : &lt;!--
    Comment --> Return : --&gt;
    Comment --> Comment : else

    Void --> DocType : &lt;!DOCTYPE
    DocType --> Return : &gt;
    DocType --> DocType : else

    Void --> Open : &lt;
    Open --> SelfClose : /&gt;
    Open --> NoTag_L : else
    SelfClose --> Return

    Void --> Close : &lt;/
    Close --> Return : &gt;
    Close --> Close : else

    Void --> NoTag_R : else
    NoTag_R --> Return : &lt;
    NoTag_R --> NoTag_R : else

    NoTag_L --> Return

    state Return <<final>>
```
