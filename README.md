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
        "firstname": "Alice",
    }
}
let context = {
    "description": "this context isn't being referenced in expression"
}
Transmute(expression, context) === {
    "person": {
        "firstname": "Alice".
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


### `#each`

`#each` creates arrays from `context`. `#each`s value can be a JSON-Path
referencing `context` or an array. `#each` iterates the items in array, and
provides access to the current items as `$.it` while `$.parent` refers to the
original `context`.

```javascript
let expression = {
    "#each": "$.tags",
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
    "#each": "$.products.*",
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
