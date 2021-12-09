package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/k-samuel/go-faceted-search/pkg/filter"
	idx "github.com/k-samuel/go-faceted-search/pkg/index"
	facet "github.com/k-samuel/go-faceted-search/pkg/search"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"
)

var oilsSearch *facet.Search
var shoeSearch *facet.Search

// inmemory databases (to simplify example)
type inmemoryDb map[int64]map[string]interface{}

var oilsDb = make(inmemoryDb, 7000)
var shoeDb = make(inmemoryDb, 12000)

type recordExtractor func(result map[string]interface{}) map[string]interface{}

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	fmt.Println("Loading indexes")
	loadIndexes()

	http.HandleFunc("/catalog/oils/", func(w http.ResponseWriter, r *http.Request) {
		facetHandler(w, r, oilsSearch, oilsDb)
	})

	http.HandleFunc("/catalog/clothing/", func(w http.ResponseWriter, r *http.Request) {
		facetHandler(w, r, shoeSearch, shoeDb)
	})

	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)

	fmt.Println("Starting server at http://localhost:8081")
	err := http.ListenAndServe(":8081", nil)

	if err != nil {
		log.Fatal(err)
	}
}

func loadIndexes() {

	oilsSearch = loadFacet(
		"data/oils.db.txt",
		oilsDb,
		func(result map[string]interface{}) map[string]interface{} {
			exclude := []string{"model"}
			if s, ok := result["fields"].(map[string]interface{}); ok {
				if len(exclude) > 0 {
					for _, name := range exclude {
						delete(s, name)
					}
				}
				return s
			}
			return nil
		})

	shoeSearch = loadFacet(
		"data/shoe.db.txt",
		shoeDb,
		func(result map[string]interface{}) map[string]interface{} {
			if s, ok := result["features"].(map[string]interface{}); ok {
				return s
			}
			return nil
		})

}

func loadFacet(filePath string, db map[int64]map[string]interface{}, extract recordExtractor) *facet.Search {
	start := time.Now()
	fmt.Print("Loading ", filePath)
	var result map[string]interface{}
	index := idx.NewIndex()
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	counter := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		result = make(map[string]interface{})
		json.Unmarshal([]byte(scanner.Text()), &result)
		id := int64(result["id"].(float64))
		db[id] = copyMap(result)
		if extract != nil {
			result = extract(result)
		}
		index.Add(id, result)
		counter++
	}
	fmt.Println(" records:", counter, " time:", time.Since(start))
	return facet.NewSearch(index)
}

func facetHandler(w http.ResponseWriter, r *http.Request, search *facet.Search, db inmemoryDb) {

	filters, err := extractFilters(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	filterResult, _ := search.AggregateFilters(filters, []int64{})
	result := map[string]interface{}{"filters": filterResult, "results": findResults(search, filters, db)}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func extractFilters(r *http.Request) (filters []filter.FilterInterface, err error) {
	filters = make([]filter.FilterInterface, 0, 0)
	err = r.ParseForm()
	s := r.FormValue("filters")
	if s == "" {
		return
	}
	var fMap map[string]interface{}
	err = json.Unmarshal([]byte(s), &fMap)
	if err != nil {
		return
	}

	for n, v := range fMap {
		if values, ok := v.([]interface{}); ok {
			list := make([]string, 0, len(values))
			for _, str := range values {
				if str, ok := str.(string); ok {
					list = append(list, str)
				}
			}
			filters = append(filters, &filter.ValueFilter{FieldName: n, Values: list})
		} else {
			fmt.Println(n, v)
		}
	}
	return
}

func findResults(search *facet.Search, filters []filter.FilterInterface, db inmemoryDb) (result map[string]interface{}) {
	pageLimit := 25
	resultIds, _ := search.Find(filters, []int64{})

	records := make([]map[string]interface{}, 0, pageLimit)
	result = map[string]interface{}{"count": len(resultIds), "limit": pageLimit, "data": &records}
	if len(resultIds) > 0 {
		for _, v := range resultIds {
			if len(records) == pageLimit {
				return
			}
			if dat, ok := db[v]; ok {
				records = append(records, dat)
			}
		}
	}
	return
}

func copyMap(m map[string]interface{}) map[string]interface{} {
	cp := make(map[string]interface{})
	for k, v := range m {
		vm, ok := v.(map[string]interface{})
		if ok {
			cp[k] = copyMap(vm)
		} else {
			cp[k] = v
		}
	}

	return cp
}
