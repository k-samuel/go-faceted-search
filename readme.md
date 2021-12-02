[![Go](https://github.com/k-samuel/go-faceted-search/actions/workflows/go.yml/badge.svg)](https://github.com/k-samuel/go-faceted-search/actions/workflows/go.yml)
# Experimental port of PHP k-samuel/faceted-search v0.2.1

PHP Library https://github.com/k-samuel/faceted-search


Bench v1.3.1 PHP 8.1.0 + JIT + opcache (no xdebug extension)

| Items count     | Memory   | Find             | Get Filters (aggregates) | Sort by field| Results Found    |
|----------------:|---------:|-----------------:|-------------------------:|-------------:|-----------------:|
| 10,000          | ~7Mb     | ~0.0007 s.       | ~0.003 s.                | ~0.0004 s.   | 907              |
| 50,000          | ~49Mb    | ~0.004 s.        | ~0.016 s.                | ~0.0009 s.   | 4550             |
| 100,000         | ~98Mb    | ~0.007 s.        | ~0.036 s.                | ~0.002 s.    | 8817             |
| 300,000         | ~242Mb   | ~0.022 s.        | ~0.135 s.                | ~0.009 s.    | 26891            |
| 1,000,000       | ~812Mb   | ~0.095 s.        | ~0.572 s.                | ~0.035 s.    | 90520            |

Bench v0.2.1 golang 1.17.3 with parallel aggregates

| Items count     | Memory   | Find             | Get Filters (aggregates) | Sort by field| Results Found    |
|----------------:|---------:|-----------------:|-------------------------:|-------------:|-----------------:|
| 10,000          | ~5Mb     | ~0.001 s.        | ~0.005 s.                | ~0.0004 s.   | 907              |
| 50,000          | ~15Mb    | ~0.011 s.        | ~0.036 s.                | ~0.002 s.    | 4550             |
| 100,000         | ~30Mb    | ~0.024 s.        | ~0.069 s.                | ~0.005 s.    | 8817             |
| 300,000         | ~128Mb   | ~0.094 s.        | ~0.220 s.                | ~0.015 s.    | 26891            |
| 1,000,000       | ~284Mb   | ~0.306 s.        | ~0.912 s.                | ~0.058 s.    | 90520            |

# Note

Search index should be created in one thread before using. Currently Index hash map access not using mutex. 
It can cause problems in concurrent writes

## Example
```go
    index := facet.NewIndex()
    search := facet.NewSearch(index)
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
        index.Add(int64(i), v)
    }
	
    // create some filters
    filters := []facet.FilterInterface{
        & facet.ValueFilter{FieldName: "color", Values: []string{"black"}},
        & facet.ValueFilter{FieldName: "size", Values: []string{"7"}},
    }
    // find records
    res, _ := search.Find(filters, []int64{})
    // aggregate filters
    info, _ := search.AggregateFilters(filters, []int64{})
```

### Test
` go test facet -coverprofile=cover.out && go tool cover -html=cover.out -o cover.html `