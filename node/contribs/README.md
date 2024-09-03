# Standard libraries usages

- `MemberExpression`

  - `Ident[Member]`
  - `Ident.Member`
    - `Ident.Object`
    - `Ident.Function`
    - `Ident.Identifier`

- Class
  - `new Ident()` (e. g. `new Session()`)
  - `new Ident.Member` (e. g.: `new console.Console(out, err)`)
  - `function (Ident)`
  - `Key: Ident`
- Functions
  - `Ident.Member()` (e. g.: `console.log('hello world')`)
  - `Ident.Member` (e. g.: `func(console.log)`)
- Primitives
  - `Ident` (e. g.: `http.maxHeaderSize`)
- Object
  - `Ident.Key` (e. g.: `inspector.Network.requestWillBeSent()`)
