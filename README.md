jsonget
-------

A simple command-line utility for extracting plain-text values from JSON.

Installation
============

    go install github.com/tysontate/jsonget

Usage
=====

Given `my.json`:

```json
{
  "foo": true,
  "bar": {
    "baz": 5
  },
  "neat": "You bet."
}
```

`jsonget` can read from files:

```bash
% jsonget --file my.json foo
true
```

Or from stdin:

```bash
% cat my.json | jsonget foo
true
```

JSON strings are returned without surrounding quotes:

```bash
% cat my.json | jsonget neat
You bet.
```

And JSON objects and arrays are returned as JSON:

```bash
% cat my.json | jsonget bar
{"baz":5}
```

TODO
====

* Code documentation
* Support accessing inside arrays

