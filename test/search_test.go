package test

import (
	"github.com/k-samuel/go-faceted-search/pkg/filter"
	"github.com/k-samuel/go-faceted-search/pkg/index"
	"github.com/k-samuel/go-faceted-search/pkg/search"
	"github.com/k-samuel/go-faceted-search/pkg/utils"
	"reflect"
	"testing"
)

func TestFind(t *testing.T) {
	facet := getSearch()
	filters := []filter.FilterInterface{
		&filter.ValueFilter{FieldName: "vendor", Values: []string{"Samsung", "Apple"}},
		&filter.ValueFilter{FieldName: "sale", Values: []string{"1"}},
		&filter.RangeFilter{FieldName: "cam_mp", Values: filter.Range{Min: 16, Type: filter.RANGE_MIN}},
		&filter.RangeFilter{FieldName: "price", Values: filter.Range{Max: 80000, Type: filter.RANGE_MAX}},
	}

	res, _ := facet.Find(filters, []int64{})
	exp := utils.FlipInt64ToMap([]int64{3, 4})

	if !reflect.DeepEqual(exp, utils.FlipInt64ToMap(res)) {
		t.Errorf("results not match\nGot:\n%v\nExpected:\n%v", res, exp)
	}

	res, _ = facet.Find(
		[]filter.FilterInterface{&filter.ValueFilter{FieldName: "vendor", Values: []string{"Google"}}},
		[]int64{},
	)

	if len(res) != 0 {
		t.Errorf("results not match\nGot:\n%v\nExpected:[]", res)
	}

}

func TestFindWithLimit(t *testing.T) {

	facet := getSearch()

	filters := make([]filter.FilterInterface, 0, 1)
	filters = append(filters, &filter.ValueFilter{FieldName: "vendor", Values: []string{"Samsung", "Apple"}})

	res, _ := facet.Find(filters, []int64{1, 3})
	exp := utils.FlipInt64ToMap([]int64{1, 3})

	if !reflect.DeepEqual(exp, utils.FlipInt64ToMap(res)) {
		t.Errorf("results not match\nGot:\n%v\nExpected:\n%v", exp, res)
	}
}

func TestGetAggregate(t *testing.T) {
	search := getSearch()
	filters := make([]filter.FilterInterface, 0, 1)
	filters = append(filters, &filter.ValueFilter{FieldName: "color", Values: []string{"black"}})

	res, _ := search.AggregateFilters(filters, []int64{})
	exp := map[string]map[string]int{
		"vendor":     {"Apple": 1, "Samsung": 2, "Xiaomi": 1},
		"model":      {"Iphone X Pro Max": 1, "Galaxy S20": 1, "Galaxy A5": 1, "MI 9": 1},
		"price":      {"80999": 1, "70599": 1, "15000": 1, "26000": 1},
		"color":      {"black": 4, "white": 1, "yellow": 1},
		"has_phones": {"1": 4},
		"cam_mp":     {"40": 1, "105": 1, "12": 1, "48": 1},
		"sale":       {"1": 3, "0": 1},
	}

	for k, v := range exp {
		if _, ok := res[k]; !ok {
			t.Errorf("Result has no expected field %v", k)
			return
		}
		for val, count := range v {
			if _, ok := res[k][val]; !ok {
				t.Errorf("Result has no expected field value %v -> %v ", k, val)
				return
			}
			if count != res[k][val] {
				t.Errorf("Unexpected count field value %v -> %v  \nGot:\n%v\nExpected:\n%v", k, val, res[k][val], count)
				return
			}
		}
	}

	if !reflect.DeepEqual(exp, res) {
		t.Errorf("results not match\nGot:\n%v\nExpected:\n%v", exp, res)
	}
}

func TestAggregateMultiFilter(t *testing.T) {
	idx := index.NewIndex()
	facet := search.NewSearch(idx)
	data := []map[string]interface{}{
		{"color": "black", "size": 7, "group": "A"},
		{"color": "black", "size": 8, "group": "A"},
		{"color": "white", "size": 7, "group": "B"},
		{"color": "yellow", "size": 7, "group": "C"},
		{"color": "black", "size": 7, "group": "C"},
	}
	for i, v := range data {
		idx.Add(int64(i), v)
	}

	filters := []filter.FilterInterface{
		&filter.ValueFilter{FieldName: "color", Values: []string{"black"}},
		&filter.ValueFilter{FieldName: "size", Values: []string{"7"}},
	}

	res, _ := facet.Find(filters, []int64{})
	info, _ := facet.AggregateFilters(filters, []int64{})
	exp := map[string]map[string]int{
		"color": {"black": 2, "white": 1, "yellow": 1},
		"size":  {"7": 2, "8": 1},
		"group": {"A": 1, "C": 1},
	}

	if !reflect.DeepEqual(exp, info) {
		t.Errorf("results not match\nGot:\n%v\nExpected:\n%v", exp, res)
	}
}

func TestAggregateNoFilter(t *testing.T) {
	idx := index.NewIndex()
	data := []map[string]interface{}{
		{"color": "black", "size": 7, "group": "A"},
		{"color": "black", "size": 8, "group": "A"},
		{"color": "white", "size": 7, "group": "B"},
		{"color": "yellow", "size": 7, "group": "C"},
		{"color": "black", "size": 7, "group": "C"},
	}
	for i, v := range data {
		idx.Add(int64(i), v)
	}
	facet := search.NewSearch(idx)

	res, _ := facet.AggregateFilters([]filter.FilterInterface{}, []int64{})
	exp := map[string]map[string]int{
		"color": {"black": 3, "white": 1, "yellow": 1},
		"size":  {"7": 4, "8": 1},
		"group": {"A": 2, "B": 1, "C": 2},
	}

	if !reflect.DeepEqual(exp, res) {
		t.Errorf("results not match\nGot:\n%v\nExpected:\n%v", res, exp)
	}
}

func getTestData() []map[string]interface{} {
	data := []map[string]interface{}{
		{
			"id":         1,
			"vendor":     "Apple",
			"model":      "Iphone X",
			"price":      80999,
			"color":      "white",
			"has_phones": false,
			"cam_mp":     20,
			"sale":       true,
		}, {
			"id":         2,
			"vendor":     "Apple",
			"model":      "Iphone X Pro Max",
			"price":      80999,
			"color":      "black",
			"has_phones": true,
			"cam_mp":     40,
			"sale":       true,
		}, {
			"id":         3,
			"vendor":     "Samsung",
			"model":      "Galaxy S20",
			"price":      70599,
			"color":      "yellow",
			"has_phones": true,
			"cam_mp":     105,
			"sale":       true,
		}, {
			"id":         4,
			"vendor":     "Samsung",
			"model":      "Galaxy S20",
			"price":      70599,
			"color":      "black",
			"has_phones": true,
			"cam_mp":     105,
			"sale":       true,
		}, {
			"id":         5,
			"vendor":     "Samsung",
			"model":      "Galaxy A5",
			"price":      15000,
			"color":      "black",
			"has_phones": true,
			"cam_mp":     12,
			"sale":       true,
		}, {
			"id":         6,
			"vendor":     "Xiaomi",
			"model":      "MI 9",
			"price":      26000,
			"color":      "black",
			"has_phones": true,
			"cam_mp":     48,
			"sale":       false,
		},
	}
	return data
}

func getSearch() *search.Search {
	idx := index.NewIndex()
	facet := search.NewSearch(idx)
	records := getTestData()

	for _, v := range records {
		id := v["id"]
		delete(v, "id")
		if dat, ok := id.(int); ok {
			idx.Add(int64(dat), v)
		}
	}
	return facet
}
