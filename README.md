jsonget
-------

A simple command-line utility for extracting plain-text values from JSON.

Installation
============

    go install github.com/tysontate/jsonget

Binaries forthcoming.

Examples
========

    % cat test.json 
    {
      "foo": true,
      "bar": {
        "baz": 5
      },
      neat: "You bet."
    }

    % jsonget test.json bar.baz
    5
    % jsonget test.json bar    
    {"baz":5}
    % jsonget test.json foo
    true
    % jsonget test.json neat
    You bet.
    % jsonget test_json/test.json bar.baz bar foo
    5
    {"baz":5,"biz":5.5}
    true

TODO
====

* Code documentation
* Handle improper access (i.e. bar.baz.wooo)
* Support accessing inside arrays
* Support piping in / out

