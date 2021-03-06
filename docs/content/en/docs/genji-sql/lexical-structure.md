---
title: "Lexical Structure"
date: 2020-03-29T17:48:51+04:00
weight: 30
description: >
    Describes of how SQL components are parsed by Genji.
---

Whenever Genji receives a query, it will be parsed according to the following rules and transformed into components Genji can understand.

## Literals

### Strings

A string is a sequence of characters surrounded by double or single quotes. They may contain any unicode character or escaped single or double quotes \(i.e `\'` or `\"`\)

```sql
foo
"l'école des fans"
'(╯ಠ_ಠ）╯︵ ┳━┳'
'foo \''
```

### Integers

An integer is a sequence of characters that only contain digits. They may start with a `+` or `-` sign.

```sql
123456789
+100
-455
```

### Floats

A float is a sequence of characters that contains three parts:

* a sequence of digits
* a decimal point \(i.e. `.`\)
* a sequence of digits

They may start with a `+` or a `-` sign.

```sql
123.456
+3.14
-1.0
```

### Booleans

A boolean is any sequence of character that is written as `true` or `false`, regardless of the case.

```sql
true
false
TRUE
FALSE
tRUe
FALse
```

### Durations

A duration is any sequence of character that starts with an [integer](#integers) and ends with one of the following unit:

* `ns` for nanoseconds (0.000000001 seconds)
* `us` or `µs` for microseconds (0.000001 seconds)
* `ms` for microseconds (0.001 seconds)
* `s` for seconds
* `m` for minutes
* `h` for hours
* `d` for days
* `w` for weeks

```sql
1ms
40ns
6w
```

### Arrays

An array is any sequence of character that starts and ends with either:

* `(` and `)`
* `[` and `]`

and that contains a coma-separated list of [expressions]({{< relref "/docs/genji-sql/expressions" >}}).

```python
[1.5, "hello", 1 > 10, [true, -10], {foo: "bar"}]
```

### Documents

A document is any sequence of character that starts and ends with `{` and `}` and that contains a list of pairs.
Each pair associates an identifier with an [expression]({{< relref "/docs/genji-sql/expressions" >}}), both separated by a colon. Each pair must be separated by a coma.

```js
{
    foo: 1,
    bar: "hello",
    baz: true AND false,
    "long field": {
        a: 10
    }
}
```

In a document, the identifiers are referred to as **fields**.
In the example above, the document has four top-level fields (`foo`, `bar`, `baz` and `long field`) and one nested field `a`.

Note that any JSON object is a valid document.

## Identifiers

Identifiers are a sequence of characters that refer to table names, field names and index names.

Identifiers may be unquoted or surrounded by backquotes. Depending on that, different rules may apply.

| Unquoted identifiers | Identifiers surrounded by backquotes |
| ---|--- |
| Must begin with an uppercase or lowercase ASCII character or an underscore | May contain any unicode character, other than the new line character (i.e. `\n`) |
| May contain only ASCII letters, digits and underscore | May contain escaped `"` character (i.e. `\"`) |

```text
foo
_foo_123_
`頂きます (*｀▽´)_旦~~`
`foo \` bar`
```

## Dot notation

The [dot notation]({{< relref "/docs/genji-sql/documents" >}}#dot-notation) is any sequence of characters that one or more [identifiers](#identifiers) separated by dots.

```text
foo
foo.bar
foo."bar baz".0.bat
```

Depending on the context, a single identifier with no dot will be parsed an identifier or as dot notation.
