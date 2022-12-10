[![Go](https://github.com/k-samuel/go-faceted-search/actions/workflows/go.yml/badge.svg)](https://github.com/k-samuel/go-faceted-search/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/k-samuel/go-faceted-search?style=flat-square)](https://goreportcard.com/report/github.com/k-samuel/go-faceted-search)
[![Release](https://img.shields.io/github/release/golang-standards/project-layout.svg?style=flat-square)](https://github.com/k-samuel/go-faceted-search/releases/latest)

# Experimental port of PHP k-samuel/faceted-search v0.3.3

PHP Library https://github.com/k-samuel/faceted-search

### Golang version benchmark

Bench v0.3.3 golang 1.19.4 with parallel aggregates

| Items count     | Memory   | Find             | Get Filters (aggregates) | Sort by field| Results Found    |
|----------------:|---------:|-----------------:|-------------------------:|-------------:|-----------------:|
| 10,000          | ~7Mb     | ~0.0003 s.       | ~0.002 s.                | ~0.0002 s.   | 907              |
| 50,000          | ~14Mb    | ~0.001 s.        | ~0.012 s.                | ~0.001 s.    | 4550             |
| 100,000         | ~21Mb    | ~0.003 s.        | ~0.025 s.                | ~0.002 s.    | 8817             |
| 300,000         | ~47Mb    | ~0.010 s.        | ~0.082 s.                | ~0.006 s.    | 26891            |
| 1,000,000       | ~140Mb   | ~0.037 s.        | ~0.285 s.                | ~0.026 s.    | 90520            |

Bench v0.3.3 golang 1.17.3 with parallel aggregates

| Items count     | Memory   | Find             | Get Filters (aggregates) | Sort by field| Results Found    |
|----------------:|---------:|-----------------:|-------------------------:|-------------:|-----------------:|
| 10,000          | ~7Mb     | ~0.0003 s.       | ~0.002 s.                | ~0.0002 s.   | 907              |
| 50,000          | ~14Mb    | ~0.002 s.        | ~0.015 s.                | ~0.001 s.    | 4550             |
| 100,000         | ~20Mb    | ~0.004 s.        | ~0.030 s.                | ~0.002 s.    | 8817             |
| 300,000         | ~44Mb    | ~0.012 s.        | ~0.086 s.                | ~0.007 s.    | 26891            |
| 1,000,000       | ~142Mb   | ~0.046 s.        | ~0.297 s.                | ~0.027 s.    | 90520            |


### PHP version benchmark

PHP 8.2  v2.1.5 ArrayIndex  JIT + opcache (no xdebug extension)

| Items count     | Memory   | Find             | Get Filters & Count (aggregate)| Sort by field| Results Found    |
|----------------:|---------:|-----------------:|-------------------------------:|-------------:|-----------------:|
| 10,000          | ~6Mb     | ~0.0004 s.       | ~0.002 s.                      | ~0.0001 s.   | 907              |
| 50,000          | ~40Mb    | ~0.001 s.        | ~0.011 s.                      | ~0.0005 s.   | 4550             |
| 100,000         | ~80Mb    | ~0.003 s.        | ~0.024 s.                      | ~0.001 s.    | 8817             |
| 300,000         | ~189Mb   | ~0.010 s.        | ~0.082 s                       | ~0.003 s.    | 26891            |
| 1,000,000       | ~657Mb   | ~0.046 s.        | ~0.306 s.                      | ~0.015 s.    | 90520            |

PHP 8.2  v2.1.5 FixedArrayIndex JIT + opcache (no xdebug extension)

| Items count     | Memory   | Find             | Get Filters & Count (aggregate)| Sort by field| Results Found    |
|----------------:|---------:|-----------------:|-------------------------------:|-------------:|-----------------:|
| 10,000          | ~2Mb     | ~0.0006 s.       | ~0.003 s.                      | ~0.0002 s.   | 907              |
| 50,000          | ~12Mb    | ~0.003 s.        | ~0.017 s.                      | ~0.0009 s.   | 4550             |
| 100,000         | ~23Mb    | ~0.006 s.        | ~0.040 s.                      | ~0.001 s.    | 8817             |
| 300,000         | ~70Mb    | ~0.019 s.        | ~0.120 s.                      | ~0.006 s.    | 26891            |
| 1,000,000       | ~233Mb   | ~0.077 s.        | ~0.455 s.                      | ~0.023 s.    | 90520            |


# Note

Search index should be created in one thread before using. Currently, Index hash map access not using mutex. 
It can cause problems in concurrent writes and reads.

## Example
```go
    package main

    import (
    "github.com/k-samuel/go-faceted-search/pkg/filter"
    "github.com/k-samuel/go-faceted-search/pkg/index"
    "github.com/k-samuel/go-faceted-search/pkg/search"
    )

    idx := index.NewIndex()
    facet := search.NewSearch(idx)
    // example data
    data := []map[string]interface{}{
        {"color": "black", "size": 7, "group": "A"},
        {"color": "black", "size": 8, "group": "A"},
        {"color": "white", "size": 7, "group": "B"},
        {"color": "yellow", "size": 7, "group": "C"},
        {"color": "black", "size": 7, "group": "C"},
    }
    // Add values using i as recordId
    for i, v := range data {
        idx.Add(int64(i), v)
    }
	
    // create some filters
    filters := []filter.FilterInterface{
        & filter.ValueFilter{FieldName: "color", Values: []string{"black"}},
        & filter.ValueFilter{FieldName: "size", Values: []string{"7"}},
    }
    // find records
    res, _ := facet.Find(filters, []int64{})
    // aggregate filters
    info, _ := facet.AggregateFilters(filters, []int64{})
```

### More examples

[Web Server](./example/)

### Test
` go test  ./test  -coverpkg  ./pkg/... -v -coverprofile=cover.out && go tool cover -html=cover.out -o cover.html `