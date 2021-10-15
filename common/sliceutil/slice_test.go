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
