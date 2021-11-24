# Experimental port of PHP k-samuel/faceted-search

PHP Library https://github.com/k-samuel/faceted-search


Bench v1.3.0 PHP 7.4.25 (no xdebug extension)

| Items count     | Memory   | Find             | Get Filters (aggregates) | Sort by field| Results Found    |
|----------------:|---------:|-----------------:|-------------------------:|-------------:|-----------------:|
| 10,000          | ~7Mb     | ~0.0004 s.       | ~0.003 s.                | ~0.0002 s.   | 907              |
| 50,000          | ~49Mb    | ~0.003 s.        | ~0.019 s.                | ~0.0008 s.   | 4550             |
| 100,000         | ~98Mb    | ~0.007 s.        | ~0.042 s.                | ~0.002 s.    | 8817             |
| 300,000         | ~242Mb   | ~0.021 s.        | ~0.167 s.                | ~0.009 s.    | 26891            |
| 1,000,000       | ~812Mb   | ~0.107 s.        | ~0.687 s.                | ~0.036 s.    | 90520            |

Bench v0.1.0 golang 1.17.3

| Items count     | Memory   | Find             | Get Filters (aggregates) | Sort by field| Results Found    |
|----------------:|---------:|-----------------:|-------------------------:|-------------:|-----------------:|
| 10,000          | ~5Mb     | ~0.001 s.        |                          |              | 907              |
| 50,000          | ~15Mb    | ~0.006 s.        |                          |              | 4550             |
| 100,000         | ~30Mb    | ~0.014 s.        |                          |              | 8817             |
| 300,000         | ~128Mb   | ~0.064 s.        |                          |              | 26891            |
| 1,000,000       | ~284Mb   | ~0.217 s.        | ~0.740                   |              | 90520            |
