package test

import (
	"github.com/k-samuel/go-faceted-search/pkg/filter"
	"github.com/k-samuel/go-faceted-search/pkg/index"
	"github.com/k-samuel/go-faceted-search/pkg/search"
	"github.com/k-samuel/go-faceted-search/pkg/sorter"
	"reflect"
	"testing"
)

func getSortTestStringData() []map[string]interface{} {
	return []map[string]interface{}{
		{"id": 1, "size": "AA", "tag": 1},
		{"id": 2, "size": "C", "tag": 1},
		{"id": 3, "size": "AAAA", "tag": 1},
		{"id": 4, "size": "B", "tag": 1},
		{"id": 5, "size": "O", "tag": 2},
	}
}
func getSortTestIntData() []map[string]interface{} {
	return []map[string]interface{}{
		{"id": 1, "size": 12, "tag": 1},
		{"id": 2, "size": 6, "tag": 1},
		{"id": 3, "size": 100, "tag": 1},
		{"id": 4, "size": 8, "tag": 1},
		{"id": 5, "size": 8, "tag": 2},
	}
}

func creteStringSortTestFacet() *search.Search {

	idx := index.NewIndex()
	facet := search.NewSearch(idx)
	records := getSortTestStringData()

	for _, v := range records {
		id := v["id"]
		delete(v, "id")
		if dat, ok := id.(int); ok {
			idx.Add(int64(dat), v)
		}
	}
	return facet
}
func creteIntSortTestFacet() *search.Search {

	idx := index.NewIndex()
	facet := search.NewSearch(idx)
	records := getSortTestIntData()

	for _, v := range records {
		id := v["id"]
		delete(v, "id")
		if dat, ok := id.(int); ok {
			idx.Add(int64(dat), v)
		}
	}
	return facet
}

func TestSortStringDesc(t *testing.T) {

	facet := creteStringSortTestFacet()
	filters := []filter.FilterInterface{
		&filter.ValueFilter{FieldName: "tag", Values: []string{"1"}},
	}
	res, _ := facet.Find(filters, []int64{})

	srt := sorter.NewStringSorter(facet.GetIndex())
	res, _ = srt.Sort(res, "size", sorter.SORT_DESC)
	exp := []int64{2, 4, 3, 1}
	if !reflect.DeepEqual(exp, res) {
		t.Errorf("results not match\nGot:\n%v\nExpected:\n%v", res, exp)
	}
}

func TestSortStringAsc(t *testing.T) {

	facet := creteStringSortTestFacet()
	filters := []filter.FilterInterface{
		&filter.ValueFilter{FieldName: "tag", Values: []string{"1"}},
	}
	res, _ := facet.Find(filters, []int64{})

	srt := sorter.NewStringSorter(facet.GetIndex())
	res, _ = srt.Sort(res, "size", sorter.SORT_ASC)
	exp := []int64{1, 3, 4, 2}
	if !reflect.DeepEqual(exp, res) {
		t.Errorf("results not match\nGot:\n%v\nExpected:\n%v", res, exp)
	}
}

func TestSortIntDesc(t *testing.T) {

	facet := creteIntSortTestFacet()
	filters := []filter.FilterInterface{
		&filter.ValueFilter{FieldName: "tag", Values: []string{"1"}},
	}
	res, _ := facet.Find(filters, []int64{})

	srt := sorter.NewIntSorter(facet.GetIndex())
	res, _ = srt.Sort(res, "size", sorter.SORT_DESC)
	exp := []int64{3, 1, 4, 2}
	if !reflect.DeepEqual(exp, res) {
		t.Errorf("results not match\nGot:\n%v\nExpected:\n%v", res, exp)
	}
}

func TestSortIntAsc(t *testing.T) {

	facet := creteIntSortTestFacet()
	filters := []filter.FilterInterface{
		&filter.ValueFilter{FieldName: "tag", Values: []string{"1"}},
	}
	res, _ := facet.Find(filters, []int64{})

	srt := sorter.NewIntSorter(facet.GetIndex())
	res, _ = srt.Sort(res, "size", sorter.SORT_ASC)
	exp := []int64{2, 4, 1, 3}
	if !reflect.DeepEqual(exp, res) {
		t.Errorf("results not match\nGot:\n%v\nExpected:\n%v", res, exp)
	}
}
