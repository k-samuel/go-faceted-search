package test

import (
	"github.com/k-samuel/go-faceted-search/pkg/utils"
	"reflect"
	"testing"
)

func TestIntersectSortedInt(t *testing.T) {
	src := []int64{1, 2}
	cmp := []int64{1, 2}

	res := utils.IntersectSortedInt(src, cmp)
	exp := []int64{1, 2}
	if !reflect.DeepEqual(exp, res) {
		t.Errorf("results not match\nGot:\n%v\nExpected:\n%v", res, exp)
	}

	src = []int64{1, 2, 7}
	cmp = []int64{1, 8}

	res = utils.IntersectSortedInt(src, cmp)
	exp = []int64{1}
	if !reflect.DeepEqual(exp, res) {
		t.Errorf("results not match\nGot:\n%v\nExpected:\n%v", res, exp)
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
