# Experimental port of PHP k-samuel/faceted-search

PHP Library https://github.com/k-samuel/faceted-search


Bench v1.3.0 PHP 7.4.25 + opcache (no xdebug extension)

| Items count     | Memory   | Find             | Get Filters (aggregates) | Sort by field| Results Found    |
|----------------:|---------:|-----------------:|-------------------------:|-------------:|-----------------:|
| 10,000          | ~7Mb     | ~0.0004 s.       | ~0.003 s.                | ~0.0002 s.   | 907              |
| 50,000          | ~49Mb    | ~0.003 s.        | ~0.019 s.                | ~0.0008 s.   | 4550             |
| 100,000         | ~98Mb    | ~0.007 s.        | ~0.042 s.                | ~0.002 s.    | 8817             |
| 300,000         | ~242Mb   | ~0.021 s.        | ~0.167 s.                | ~0.009 s.    | 26891            |
| 1,000,000       | ~812Mb   | ~0.107 s.        | ~0.687 s.                | ~0.036 s.    | 90520            |

Bench v0.2.0 golang 1.17.3 with parallel aggregates

| Items count     | Memory   | Find             | Get Filters (aggregates) | Sort by field| Results Found    |
|----------------:|---------:|-----------------:|-------------------------:|-------------:|-----------------:|
| 10,000          | ~5Mb     | ~0.002 s.        | ~0.007 s.                | ~0.0009 s.   | 907              |
| 50,000          | ~15Mb    | ~0.011 s.        | ~0.038 s.                | ~0.002 s.    | 4550             |
| 100,000         | ~30Mb    | ~0.018 s.        | ~0.061 s.                | ~0.004 s.    | 8817             |
| 300,000         | ~128Mb   | ~0.083 s.        | ~0.268 s.                | ~0.015 s.    | 26891            |
| 1,000,000       | ~284Mb   | ~0.344 s.        | ~1.020 s.                | ~0.055 s.    | 90520            |

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
	for i, v := range data {
		index.Add(int64(i), v)
	}
	
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
` go test -coverprofile=cover.out && go tool cover -html=cover.out -o cover.html `