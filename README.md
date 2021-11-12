# Go implementation of JSON Transmute

See [JSON-Transmute](https://github.com/fre-sch/json-transmute) for the full
documentation of JSON-Transmute.

## Notes

* For JSON-Path functionality, [oliveagle/jsonpath](github.com/oliveagle/jsonpath)
  is used, which deviates from the [reference document](https://goessner.net/articles/JsonPath/)
  where it returns arrays for multiple matches, but only first match of single
  matches.

* `#format` is implemented using go `text/template`. An additional template
  function `path` is available, it uses JSON-Path to reference `context`.
