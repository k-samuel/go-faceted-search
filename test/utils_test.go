package test

import (
	"github.com/k-samuel/go-faceted-search/pkg/utils"
	"reflect"
	"testing"
)

type testCaseData struct {
	src []int64
	cmp []int64
	exp []int64
}

func intersectData() []testCaseData {
	return []testCaseData{
		{src: []int64{1, 2}, cmp: []int64{1, 2}, exp: []int64{1, 2}},
		{src: []int64{10, 21, 123, 124}, cmp: []int64{1, 2, 22, 123, 127}, exp: []int64{123}},
		{src: []int64{1}, cmp: []int64{2, 3, 6}, exp: []int64{}},
		{src: []int64{1, 7, 8, 9}, cmp: []int64{2, 7}, exp: []int64{7}},
		{src: []int64{1, 2, 3, 4, 6}, cmp: []int64{3, 4, 5, 6}, exp: []int64{3, 4, 6}},
	}
}

func TestIntersectSortedInt(t *testing.T) {

	data := intersectData()
	var res []int64

	for _, v := range data {
		res = utils.IntersectSortedInt(v.src, v.cmp)
		if !reflect.DeepEqual(v.exp, res) {
			t.Errorf("results not match\nGot:\n%v\nExpected:\n%v", res, v.exp)
		}
	}
}

func BenchmarkIntersectSortedInt(b *testing.B) {
	data := intersectData()
	for i := 0; i < b.N; i++ {
		for _, v := range data {
			utils.IntersectSortedInt(v.src, v.cmp)
		}
	}
}

func TestIntersectCountSortedInt(t *testing.T) {
	src := []int64{1, 2}
	cmp := []int64{1, 2}

	res := utils.IntersectCountSortedInt(src, cmp)
	exp := 2
	if !reflect.DeepEqual(exp, res) {
		t.Errorf("results not match\nGot:\n%v\nExpected:\n%v", res, exp)
	}

	src = []int64{1, 2, 7}
	cmp = []int64{1, 8}

	res = utils.IntersectCountSortedInt(src, cmp)
	exp = 1
	if !reflect.DeepEqual(exp, res) {
		t.Errorf("results not match\nGot:\n%v\nExpected:\n%v", res, exp)
	}
}
