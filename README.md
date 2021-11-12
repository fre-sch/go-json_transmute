# JSON Transmute

Define transformations of JSON with JSON.

This document does not describe a specific implementation,
it describes the  Transmute *"protocol"* that implementations
should adhere to.

```typescript
type JSON = null|boolean|number|string|array|object

Transmute(Expression: JSON, Context: JSON) -> JSON
```

**`Expression`**
: Describes a new JSON object created from `Context`.  `Expression` is always fully recursive.

**`Context`**
: Application specific data.


## Expressions

`expression` describes the resulting JSON. It can contain operator expressions
referencing `context`. Any valid JSON is a valid transmute `expression`. Without
any operator, `expression` simply describes the resulting JSON value.

```javascript
let expression = {
    "person": {
        "firstname": "Alice"
    }
}
let context = {
    "description": "this context isn't being referenced in expression"
}
Transmute(expression, context) === {
    "person": {
        "firstname": "Alice"
    }
}
```

### JSON-Path

Property values must be checked if they're JSON-Path strings, that is, they're
starting with a `$` (dollar) symbol. In that case, the JSON-Path is evaluated
against `context` and matching value is retrieved.

```javascript
let expression = "$.person.firstname"
let context = {
    "person": {
        "firstname": "Albert"
    }
}
Transmute(expression, context) === ["Albert"]
```

```javascript
let expression = "$.products.*.price"
let context = {
    "products": [
        {
            "price": "$4.00"
        },
        {
            "price": "$12.99"
        }
    ]
}
Transmute(expression, context) === ["$4.00", "$12.99"]
```

```javascript
let expression = "$.products.*.title"
let context = {
    "products": [
        {
            "price": "$4.00"
        },
        {
            "price": "$12.99"
        }
    ]
}
Transmute(expression, context) === []
```

## Operators

### `#format`

`#format` provides template strings and uses `Context` to create a new string.
Template system is implementation specific, however it must provide access to
`context`.

```javascript
let expression = {
    "#format": "Hello {person.firstname}!"
}
let context = {
    "person": {
        "firstname": "Berta"
    }
}

Transmute(expression, context) === "Hello Berta!"
```


### `#first`

`#first` retrieves the first item of an array.

```javascript
let expression = {
    "#first": "$.*.firstname"
}
let context = {
    "ID1": {
        "firstname": "Alice"
    },
    "ID2": {
        "firstname": "Bob"
    }
}
Transmute(expression, context) === "Albert"
```


### `#map`

`#map` creates arrays from `context`. `#map`s value can be a JSON-Path
referencing `context` or an array. `#map` iterates the items in array, and
provides access to the current items as `$.it` while `$.parent` refers to the
original `context`.

```javascript
let expression = {
    "#map": "$.tags",
    "title": "$.it",
    "price": "$.parent.defaultPrice"
}
let context = {
    "defaultPrice": 3.00,
    "tags": ["beer", "wine", "peanuts"]
}
Transmute(expression, context) === [
    {
        "title": "beer",
        "price": 3.00
    },
    {
        "title": "wine",
        "price": 3.00
    },
    {
        "title": "peanuts",
        "price": 3.00
    }
]
```

```javascript
let expression = {
    "#map": "$.products.*",
    "title": "$.it.title.en",
    "price": "$.it.price.USD"
}
let context = {
    "products": [
        {
            "title": {
                "en": "Beer",
                "de": "Bier",
            },
            "price": {
                "USD": 12.00,
                "EUR": 9.99
            }
        },
        {
            "title": {
                "ja": "豚カツ",
                "de": "Schweineschnitzel",
                "en": "Pork cutlet, breaded",
            },
            "price": {
                "USD": 3.00,
                "EUR": 2.49,
                "JPY": 300,
            }
        }
    ]
}

Transmute(expression, context) === [
    {
        "title": "Beer",
        "price": 9.99
    },
    {
        "title": "Pork cutlet, breaded",
        "price": 3.00
    }
]
```

### `#sum`

Sum of values.

```javascript
let expression = {
    "#sum": "$.products.*.price"
}
let context = {
    "products": [
        {
            "price": 1
        },
        {
            "price": 2
        },
        {
            "price": "invalid value",
        },
        {
            "price": 3,
        }
    ]
}
Transmute(expression, context) === 6
```

### `#join`

```javascript
let expression = {
    "#join": "$.products.*.title",
    "#separator": " | "
}
let context = {
    "products": [
        {
            "title": "Beer",
            "price": 1
        },
        {
            "title": "Pizza",
            "price": 2
        },
        {
            "title": "Brownies",
            "price": 3
        }
    ]
}
Transmute(expression, context) === "Beer | Pizza | Brownies"
```


### `#coalesce`

```javascript
let expression = {
    "#coalesce": "$.products"
}
let context = {
    "products": [
        {
            "title": "one"
        },
        {
            "title": "two"
        },
        null,
        {
            "title": "three",
        }
    ]
}
Transmute(expression, context) === [
    {
        "title": "one"
    },
    {
        "title": "two"
    },
    {
        "title": "three"
    }
]
```

### `#transmute`

Defer value evaluation

```javascript
let expression = {
    "#map": {
        "#transmute": "$.propname"
    },
    "label": "$.it.title"
}
let context = {
    "propname": "products",
    "products": [
        {
            "title": "one"
        },
        {
            "title": "two"
        },
        {
            "title": "three"
        }
    ]
}
Transmute(expression, context) === [
    {
        "label": "one"
    },
    {
        "label": "two"
    },
    {
        "label": "three"
    }
]
```
