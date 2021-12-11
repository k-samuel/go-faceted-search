[![Go](https://github.com/k-samuel/go-faceted-search/actions/workflows/go.yml/badge.svg)](https://github.com/k-samuel/go-faceted-search/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/k-samuel/go-faceted-search?style=flat-square)](https://goreportcard.com/report/github.com/k-samuel/go-faceted-search)
[![Release](https://img.shields.io/github/release/golang-standards/project-layout.svg?style=flat-square)](https://github.com/k-samuel/go-faceted-search/releases/latest)
# Experimental port of PHP k-samuel/faceted-search v0.3.0

PHP Library https://github.com/k-samuel/faceted-search

Bench v1.3.3 PHP 8.1.0 + JIT + opcache (no xdebug extension)

| Items count     | Memory   | Find             | Get Filters (aggregates) | Sort by field| Results Found    |
|----------------:|---------:|-----------------:|-------------------------:|-------------:|-----------------:|
| 10,000          | ~7Mb     | ~0.0007 s.       | ~0.004 s.                | ~0.0003 s.   | 907              |
| 50,000          | ~49Mb    | ~0.002 s.        | ~0.014 s.                | ~0.0009 s.   | 4550             |
| 100,000         | ~98Mb    | ~0.004 s.        | ~0.028 s.                | ~0.002 s.    | 8817             |
| 300,000         | ~242Mb   | ~0.012 s.        | ~0.112 s.                | ~0.007 s.    | 26891            |
| 1,000,000       | ~812Mb   | ~0.057 s.        | ~0.443 s.                | ~0.034 s.    | 90520            |

Bench v0.3.0 golang 1.17.3 with parallel aggregates

| Items count     | Memory   | Find             | Get Filters (aggregates) | Sort by field| Results Found    |
|----------------:|---------:|-----------------:|-------------------------:|-------------:|-----------------:|
| 10,000          | ~5Mb     | ~0.0004 s.       | ~0.001 s.                | ~0.0002 s.   | 907              |
| 50,000          | ~15Mb    | ~0.003 s.        | ~0.012 s.                | ~0.001 s.    | 4550             |
| 100,000         | ~21Mb    | ~0.006 s.        | ~0.028 s.                | ~0.003 s.    | 8817             |
| 300,000         | ~47Mb    | ~0.018 s.        | ~0.91 s.                 | ~0.010 s.    | 26891            |
| 1,000,000       | ~150Mb   | ~0.094 s.        | ~0.408 s.                | ~0.040 s.    | 90520            |

Bench v0.2.4 golang 1.17.3 with parallel aggregates

| Items count     | Memory   | Find             | Get Filters (aggregates) | Sort by field| Results Found    |
|----------------:|---------:|-----------------:|-------------------------:|-------------:|-----------------:|
| 10,000          | ~5Mb     | ~0.0009 s.       | ~0.002 s.                | ~0.0004 s.   | 907              |
| 50,000          | ~15Mb    | ~0.005 s.        | ~0.019 s.                | ~0.002 s.    | 4550             |
| 100,000         | ~30Mb    | ~0.011 s.        | ~0.043 s.                | ~0.005 s.    | 8817             |
| 300,000         | ~128Mb   | ~0.053 s.        | ~0.144 s.                | ~0.014 s.    | 26891            |
| 1,000,000       | ~284Mb   | ~0.130 s.        | ~0.577 s.                | ~0.057 s.    | 90520            |

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
` go test facet -coverprofile=cover.out && go tool cover -html=cover.out -o cover.html `