package partial

import (
	"fmt"
	"golang.org/x/exp/slices"
	"math/rand"
	"sort"
	"testing"
)

func TestSort(t *testing.T) {
	cases := []testCase[int]{
		{[]int{}, 0},
		{[]int{2}, 1},
		{[]int{2, 1}, 1},
		{[]int{2, 1}, 2},
		{[]int{1, 1, 1}, 2},
		{[]int{5, 0, 0, 0, 1}, 2},
		{[]int{5, 0, 0, 0, 1}, 5},
	}
	rand.Seed(2)
	big := make([]int, 100_000)
	for i := 0; i < 100_000; i++ {
		big[i] = rand.Intn(10_000)
	}
	cases = append(cases, testCase[int]{big, 10_000})
	for _, c := range cases {
		Sort(c.x, c.k)
		if !slices.IsSorted(c.x[:c.k]) {
			t.Errorf("Not sorted, out=%v, k=%v", c.x, c.k)
		}
	}
}

func TestSortFunc(t *testing.T) {
	cases := []testCase[person]{
		{[]person{{"bob", 45}, {"jane", 31}}, 1},
		{[]person{{"bob", 45}, {"jane", 31}}, 2},
		{[]person{{"bob", 45}, {"jane", 31}, {"karl", 39}}, 2},
		{[]person{{"bob", 45}, {"jane", 31}, {"karl", 39}}, 3},
	}
	less := func(x, y person) bool { return x.age < y.age }
	for _, c := range cases {
		SortFunc(c.x, c.k, less)
		if !slices.IsSortedFunc(c.x[:c.k], less) {
			t.Errorf("Not sorted, out=%v, k=%v", c.x, c.k)
		}
	}
}

func TestSortOutOfBounds(t *testing.T) {
	less := func(x, y int) bool { return x < y }

	x := []int{9, 2, 5}
	Sort(x, -1)
	if !slices.Equal(x, []int{9, 2, 5}) {
		t.Errorf("Negative k should be treated as zero and sort nothing")
	}

	y := []int{9, 2, 5}
	SortFunc(y, 5, less)
	if !slices.Equal(y, []int{2, 5, 9}) {
		t.Errorf("Entire slice should be sorted when k is greater than len")
	}
}

func BenchmarkSort(b *testing.B) {
	sizes := []int{1_000, 10_000, 100_000}
	k := 100
	for _, size := range sizes {
		var x []int
		for i := 0; i < size; i++ {
			x = append(x, rand.Intn(size/10))
		}
		b.Run(fmt.Sprintf("sort.Slice_%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				y := slices.Clone(x)
				b.StartTimer()
				sort.Slice(y, func(i, j int) bool { return y[i] < y[j] })
			}
		})
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
		b.Run(fmt.Sprintf("partial.Sort_%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				y := slices.Clone(x)
				b.StartTimer()
				Sort(y, k)
			}
		})
		b.Run(fmt.Sprintf("partial.SortFunc_%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				y := slices.Clone(x)
				b.StartTimer()
				SortFunc(y, k, func(i, j int) bool { return i < j })
			}
		})
	}
}
