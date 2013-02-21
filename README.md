jsonget
-------

`jsonget` is a command-line tool for extracting values out from JSON. This is
useful, for example, when you're `curling` a JSON api and just want to get a
single value from it:

```bash
% curl http://openweathermap.org/data/2.0/weather/city/524901 | jsonget main.temp
259.92
```

Or multiple newline-separated values:

```bash
% curl http://openweathermap.org/data/2.0/weather/city/524901 | jsonget name main.temp wind.speed
Moscow
259.92
2
```

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

