# JSON Transmute

Define transformations of JSON with JSON.

This document does not describe a specific implementation,
it describes the  Transmute *"protocol"* that implementations
should adhere to.

```typescript
type JSON =  null|boolean|number|string|array|object

Transmute(Expression: JSON, Context: JSON) -> JSON
```

**`Expression`**
: Describes a new JSON object created from `Context`.  `Expression` is always fully recursive.

**`Context`**
: Application specific data.


## Expressions

### `#format`

`#format` provides template strings and uses `Context` to create a new string.
Template system is implementation specific.

```javascript
let expression = {
    "#format": "Hello {{path \"$.name\"}}!"
}
let context = {
    "name": "Albert"
}
Transmute(expression, context)
// => "Hello Albert!"
```
