package partial

import (
	"fmt"
	"golang.org/x/exp/slices"
	"math/rand"
	"testing"
)

func checkTopKInvariants[E any](x []E, k int, less func(E, E) bool) bool {
	sorted := slices.Clone(x)
	slices.SortFunc(sorted, less)

	if len(x) < 2 {
		return true
	}

	// Kth element should be in sorted position
	if less(x[k-1], sorted[k-1]) || less(sorted[k-1], x[k-1]) { // x[k-1] != sorted[k-1]
		return false
	}

	// All elements before the kth should be less or equal
	for _, v := range x[:k-1] {
		if less(x[k-1], v) {
			return false
		}
	}

	// All elements following the kth should be greater or equal
	for _, v := range x[k:] {
		if less(v, x[k-1]) {
			return false
		}
	}

	return true
}

type testCase[E any] struct {
	x []E
	k int
}

func TestTopK(t *testing.T) {
	rand.Seed(2)
	cases := []testCase[int]{
		{[]int{}, 1},
		{[]int{2}, 1},
		{[]int{2, 1}, 1},
		{[]int{2, 1}, 2},
		{[]int{1, 1, 1}, 2},
		{[]int{5, 0, 0, 0, 1}, 2},
		{[]int{5, 0, 0, 0, 1}, 5},
	}
	big := make([]int, 100_000)
	for i := 0; i < 100_000; i++ {
		big[i] = rand.Intn(10_000)
	}
	cases = append(cases, testCase[int]{big, 10_000})
	less := func(x, y int) bool { return x < y }
	for _, c := range cases {
		x := slices.Clone(c.x)
		TopK(x, c.k)
		if !checkTopKInvariants(x, c.k, less) {
			t.Errorf("Invariants failed, in=%v, k=%v, out=%v.", c.x, c.k, x)
		}
	}
}

type person struct {
	name string
	age  int
}

func TestTopKFunc(t *testing.T) {
	cases := []testCase[person]{
		{[]person{{"bob", 45}, {"jane", 31}}, 1},
		{[]person{{"bob", 45}, {"jane", 31}}, 2},
		{[]person{{"bob", 45}, {"jane", 31}, {"karl", 31}}, 2},
		{[]person{{"bob", 45}, {"jane", 31}, {"karl", 31}}, 3},
	}
	less := func(x, y person) bool { return x.age < y.age }
	for _, c := range cases {
		x := slices.Clone(c.x)
		TopKFunc(x, c.k, less)
		if !checkTopKInvariants(x, c.k, less) {
			t.Errorf("Invariants failed, in=%v, k=%v, out=%v.", c.x, c.k, x)
		}
	}
}

func TestTopKOutOfBounds(t *testing.T) {
	less := func(x, y int) bool { return x < y }

	x := []int{9, 2, 5}
	TopK(x, -1)
	if !slices.Equal(x, []int{9, 2, 5}) {
		t.Errorf("Negative k should be treated as zero and sort nothing")
	}

	y := []int{9, 2, 5}
	TopK(y, 5)
	if !checkTopKInvariants(y, 3, less) {
		t.Errorf("Should take TopK of entire slice when k is greater than len")
	}
}

func BenchmarkTopK(b *testing.B) {
	sizes := []int{1_000, 10_000, 100_000}
	for _, size := range sizes {
		var x []int
		for i := 0; i < size; i++ {
			x = append(x, rand.Intn(size/10))
		}
		k := size / 2
		b.Run(fmt.Sprintf("slices.Sort_%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				y := slices.Clone(x)
				b.StartTimer()
				slices.Sort(y)
			}
		})
		b.Run(fmt.Sprintf("slices.SortFunc_%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				y := slices.Clone(x)
				b.StartTimer()
				slices.SortFunc(y, func(i, j int) bool { return i < j })
			}
		})
		b.Run(fmt.Sprintf("partial.Sort%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				y := slices.Clone(x)
				b.StartTimer()
				Sort(y, k)
			}
		})
		b.Run(fmt.Sprintf("partial.SortFunc%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				y := slices.Clone(x)
				b.StartTimer()
				SortFunc(y, k, func(i, j int) bool { return i < j })
			}
		})
		b.Run(fmt.Sprintf("partial.TopK_%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				y := slices.Clone(x)
				b.StartTimer()
				TopK(y, k)
			}
		})
		b.Run(fmt.Sprintf("partial.TopKFunc_%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				y := slices.Clone(x)
				b.StartTimer()
				TopKFunc(y, k, func(i, j int) bool { return i < j })
			}
		})
	}
}
