package sliceutil

import (
	"reflect"
	"sort"
	"testing"
)

func TestSubsetUint64(t *testing.T) {
	testCases := []struct {
		setA []uint64
		setB []uint64
		out  bool
	}{
		{[]uint64{1}, []uint64{1, 2, 3, 4}, true},
		{[]uint64{1, 2, 3, 4}, []uint64{1, 2, 3, 4}, true},
		{[]uint64{1, 1}, []uint64{1, 2, 3, 4}, false},
		{[]uint64{}, []uint64{1}, true},
		{[]uint64{1}, []uint64{}, false},
		{[]uint64{1, 2, 3, 4, 5}, []uint64{1, 2, 3, 4}, false},
	}
	for _, tt := range testCases {
		result := SubsetUint64(tt.setA, tt.setB)
		if result != tt.out {
			t.Errorf("%v, got %v, want %v", tt.setA, result, tt.out)
		}
	}
}

func TestIntersectionUint64(t *testing.T) {
	testCases := []struct {
		setA []uint64
		setB []uint64
		setC []uint64
		out  []uint64
	}{
		{[]uint64{2, 3, 5}, []uint64{3}, []uint64{3}, []uint64{3}},
		{[]uint64{2, 3, 5}, []uint64{3, 5}, []uint64{5}, []uint64{5}},
		{[]uint64{2, 3, 5}, []uint64{3, 5}, []uint64{3, 5}, []uint64{3, 5}},
		{[]uint64{2, 3, 5}, []uint64{5, 3, 2}, []uint64{3, 2, 5}, []uint64{2, 3, 5}},
		{[]uint64{3, 2, 5}, []uint64{5, 3, 2}, []uint64{3, 2, 5}, []uint64{2, 3, 5}},
		{[]uint64{3, 3, 5}, []uint64{5, 3, 2}, []uint64{3, 2, 5}, []uint64{3, 5}},
		{[]uint64{2, 3, 5}, []uint64{2, 3, 5}, []uint64{2, 3, 5}, []uint64{2, 3, 5}},
		{[]uint64{2, 3, 5}, []uint64{}, []uint64{}, []uint64{}},
		{[]uint64{2, 3, 5}, []uint64{2, 3, 5}, []uint64{}, []uint64{}},
		{[]uint64{2, 3}, []uint64{2, 3, 5}, []uint64{5}, []uint64{}},
		{[]uint64{2, 2, 2}, []uint64{2, 2, 2}, []uint64{}, []uint64{}},
		{[]uint64{}, []uint64{2, 3, 5}, []uint64{}, []uint64{}},
		{[]uint64{}, []uint64{}, []uint64{}, []uint64{}},
		{[]uint64{1}, []uint64{1}, []uint64{}, []uint64{}},
		{[]uint64{1, 1, 1}, []uint64{1, 1}, []uint64{1, 2, 3}, []uint64{1}},
	}
	for _, tt := range testCases {
		setA := append([]uint64{}, tt.setA...)
		setB := append([]uint64{}, tt.setB...)
		setC := append([]uint64{}, tt.setC...)
		result := IntersectionUint64(setA, setB, setC)
		sort.Slice(result, func(i, j int) bool {
			return result[i] < result[j]
		})
		if !reflect.DeepEqual(result, tt.out) {
			t.Errorf("got %d, want %d", result, tt.out)
		}
		if !reflect.DeepEqual(setA, tt.setA) {
			t.Errorf("slice modified, got %v, want %v", setA, tt.setA)
		}
		if !reflect.DeepEqual(setB, tt.setB) {
			t.Errorf("slice modified, got %v, want %v", setB, tt.setB)
		}
		if !reflect.DeepEqual(setC, tt.setC) {
			t.Errorf("slice modified, got %v, want %v", setC, tt.setC)
		}
	}
}

func TestIsSortedUint64(t *testing.T) {
	testCases := []struct {
		setA []uint64
		out  bool
	}{
		{[]uint64{1, 2, 3}, true},
		{[]uint64{3, 1, 3}, false},
		{[]uint64{1}, true},
		{[]uint64{}, true},
	}
	for _, tt := range testCases {
		result := IsUint64Sorted(tt.setA)
		if result != tt.out {
			t.Errorf("got %v, want %v", result, tt.out)
		}
	}
}

func TestIntersectionInt64(t *testing.T) {
	testCases := []struct {
		setA []int64
		setB []int64
		setC []int64
		out  []int64
	}{
		{[]int64{2, 3, 5}, []int64{3}, []int64{3}, []int64{3}},
		{[]int64{2, 3, 5}, []int64{3, 5}, []int64{5}, []int64{5}},
		{[]int64{2, 3, 5}, []int64{3, 5}, []int64{3, 5}, []int64{3, 5}},
		{[]int64{2, 3, 5}, []int64{5, 3, 2}, []int64{3, 2, 5}, []int64{2, 3, 5}},
		{[]int64{3, 2, 5}, []int64{5, 3, 2}, []int64{3, 2, 5}, []int64{2, 3, 5}},
		{[]int64{3, 3, 5}, []int64{5, 3, 2}, []int64{3, 2, 5}, []int64{3, 5}},
		{[]int64{2, 3, 5}, []int64{2, 3, 5}, []int64{2, 3, 5}, []int64{2, 3, 5}},
		{[]int64{2, 3, 5}, []int64{}, []int64{}, []int64{}},
		{[]int64{2, 3, 5}, []int64{2, 3, 5}, []int64{}, []int64{}},
		{[]int64{2, 3}, []int64{2, 3, 5}, []int64{5}, []int64{}},
		{[]int64{2, 2, 2}, []int64{2, 2, 2}, []int64{}, []int64{}},
		{[]int64{}, []int64{2, 3, 5}, []int64{}, []int64{}},
		{[]int64{}, []int64{}, []int64{}, []int64{}},
		{[]int64{1}, []int64{1}, []int64{}, []int64{}},
		{[]int64{1, 1, 1}, []int64{1, 1}, []int64{1, 2, 3}, []int64{1}},
	}
	for _, tt := range testCases {
		setA := append([]int64{}, tt.setA...)
		setB := append([]int64{}, tt.setB...)
		setC := append([]int64{}, tt.setC...)
		result := IntersectionInt64(setA, setB, setC)
		sort.Slice(result, func(i, j int) bool {
			return result[i] < result[j]
		})
		if !reflect.DeepEqual(result, tt.out) {
			t.Errorf("got %d, want %d", result, tt.out)
		}
		if !reflect.DeepEqual(setA, tt.setA) {
			t.Errorf("slice modified, got %v, want %v", setA, tt.setA)
		}
		if !reflect.DeepEqual(setB, tt.setB) {
			t.Errorf("slice modified, got %v, want %v", setB, tt.setB)
		}
		if !reflect.DeepEqual(setC, tt.setC) {
			t.Errorf("slice modified, got %v, want %v", setC, tt.setC)
		}
	}
}

func TestUnionUint64(t *testing.T) {
	testCases := []struct {
		setA []uint64
		setB []uint64
		out  []uint64
	}{
		{[]uint64{2, 3, 5}, []uint64{4, 6}, []uint64{2, 3, 5, 4, 6}},
		{[]uint64{2, 3, 5}, []uint64{3, 5}, []uint64{2, 3, 5}},
		{[]uint64{2, 3, 5}, []uint64{2, 3, 5}, []uint64{2, 3, 5}},
		{[]uint64{2, 3, 5}, []uint64{}, []uint64{2, 3, 5}},
		{[]uint64{}, []uint64{2, 3, 5}, []uint64{2, 3, 5}},
		{[]uint64{}, []uint64{}, []uint64{}},
		{[]uint64{1}, []uint64{1}, []uint64{1}},
	}
	for _, tt := range testCases {
		result := UnionUint64(tt.setA, tt.setB)
		if !reflect.DeepEqual(result, tt.out) {
			t.Errorf("got %d, want %d", result, tt.out)
		}

	}
	items := [][]uint64{
		{3, 4, 5},
		{6, 7, 8},
		{9, 10, 11},
	}
	variadicResult := UnionUint64(items...)
	want := []uint64{3, 4, 5, 6, 7, 8, 9, 10, 11}
	if !reflect.DeepEqual(want, variadicResult) {
		t.Errorf("Received %v, wanted %v", variadicResult, want)
	}
}

func TestUnionInt64(t *testing.T) {
	testCases := []struct {
		setA []int64
		setB []int64
		out  []int64
	}{
		{[]int64{2, 3, 5}, []int64{4, 6}, []int64{2, 3, 5, 4, 6}},
		{[]int64{2, 3, 5}, []int64{3, 5}, []int64{2, 3, 5}},
		{[]int64{2, 3, 5}, []int64{2, 3, 5}, []int64{2, 3, 5}},
		{[]int64{2, 3, 5}, []int64{}, []int64{2, 3, 5}},
		{[]int64{}, []int64{2, 3, 5}, []int64{2, 3, 5}},
		{[]int64{}, []int64{}, []int64{}},
		{[]int64{1}, []int64{1}, []int64{1}},
	}
	for _, tt := range testCases {
		result := UnionInt64(tt.setA, tt.setB)
		if !reflect.DeepEqual(result, tt.out) {
			t.Errorf("got %d, want %d", result, tt.out)
		}
	}
	items := [][]int64{
		{3, 4, 5},
		{6, 7, 8},
		{9, 10, 11},
	}
	variadicResult := UnionInt64(items...)
	want := []int64{3, 4, 5, 6, 7, 8, 9, 10, 11}
	if !reflect.DeepEqual(want, variadicResult) {
		t.Errorf("Received %v, wanted %v", variadicResult, want)
	}
}

func TestCleanUint64(t *testing.T) {
	testCases := []struct {
		in  []uint64
		out []uint64
	}{
		{[]uint64{2, 4, 4, 6, 6}, []uint64{2, 4, 6}},
		{[]uint64{3, 5, 5}, []uint64{3, 5}},
		{[]uint64{2, 2, 2}, []uint64{2}},
		{[]uint64{1, 4, 5, 9, 9}, []uint64{1, 4, 5, 9}},
		{[]uint64{}, []uint64{}},
		{[]uint64{1}, []uint64{1}},
	}
	for _, tt := range testCases {
		result := SetUint64(tt.in)
		if !reflect.DeepEqual(result, tt.out) {
			t.Errorf("got %d, want %d", result, tt.out)
		}
	}
}