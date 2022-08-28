package partial

import (
	"golang.org/x/exp/constraints"
	"math"
)

// TopK reorders a slice such that x[:k] contains the first k elements of the
// slice when sorted in ascending order. Only the kth element x[k-1] is
// guaranteed to be in sorted order. All elements in x[:k-1] are less than or
// equal to the kth element, all elements in x[k:] are greater than or equal.
// This is faster than using slices.Sort.
func TopK[E constraints.Ordered](x []E, k int) {
	k = min(k, len(x))
	if k > 0 {
		floydRivest(x, 0, len(x)-1, k-1) // 0-indexed
	}
}

// TopKFunc reorders a slice such that x[:k] contains the first k elements of
// the slice when sorted in ascending order as determined by the less function.
// Only the kth element x[k-1] is guaranteed to be in sorted order. All elements
// in x[:k-1] are less than or equal to the kth element, all elements in x[k:]
// are greater than or equal. This is faster than using slices.SortFunc.
func TopKFunc[E any](x []E, k int, less func(E, E) bool) {
	k = min(k, len(x))
	if k > 0 {
		floydRivestFunc(x, 0, len(x)-1, k-1, less)
	}
}

// https://en.wikipedia.org/wiki/Floyd%E2%80%93Rivest_algorithm
func floydRivest[E constraints.Ordered](x []E, left, right, k int) {
	// left is the left index for the interval
	// right is the right index for the interval
	// k is the desired index value, where x[k] is the (k+1)th smallest element when left = 0
	length := len(x)
	for right > left {
		// Use select recursively to sample a smaller set of size s
		// the arbitrary constants 600 and 0.5 are used in the original
		// version to minimize execution time.
		if right-left > 600 {
			var n = float64(right - left + 1)
			var i = float64(k - left + 1)
			var z = math.Log(n)
			var s = 0.5 * math.Exp(2*z/3)
			var sd = 0.5 * math.Sqrt(z*s*(n-s)/n) * float64(sign(i-n/2))
			var kf = float64(k)
			var newLeft = max(left, int(math.Floor(kf-i*s/n+sd)))
			var newRight = min(right, int(math.Floor(kf+(n-i)*s/n+sd)))
			floydRivest(x, newLeft, newRight, k)
		}
		// partition the elements between left and right around t
		var t = x[k]
		var i = left
		var j = right
		x[left], x[k] = x[k], x[left]
		if t < x[right] {
			x[left], x[right] = x[right], x[left]
		}
		for i < j {
			x[i], x[j] = x[j], x[i]
			i++
			j--
			for i < length && x[i] < t {
				i++
			}
			for j >= 0 && t < x[j] {
				j--
			}
		}
		if x[left] == t {
			x[left], x[j] = x[j], x[left]
		} else {
			j++
			x[j], x[right] = x[right], x[j]
		}
		// Adjust left and right towards the boundaries of the subset
		// containing the (k − left + 1)th smallest element.
		if j <= k {
			left = j + 1
		}
		if k <= j {
			right = j - 1
		}
	}
}

func floydRivestFunc[E any](x []E, left, right, k int, less func(E, E) bool) {
	// left is the left index for the interval
	// right is the right index for the interval
	// k is the desired index value, where x[k] is the (k+1)th smallest element when left = 0
	length := len(x)
	for right > left {
		// Use select recursively to sample a smaller set of size s
		// the arbitrary constants 600 and 0.5 are used in the original
		// version to minimize execution time.
		if right-left > 600 {
			var n = float64(right - left + 1)
			var i = float64(k - left + 1)
			var z = math.Log(n)
			var s = 0.5 * math.Exp(2*z/3)
			var sd = 0.5 * math.Sqrt(z*s*(n-s)/n) * float64(sign(i-n/2))
			var kf = float64(k)
			var newLeft = max(left, int(math.Floor(kf-i*s/n+sd)))
			var newRight = min(right, int(math.Floor(kf+(n-i)*s/n+sd)))
			floydRivestFunc(x, newLeft, newRight, k, less)
		}
		// partition the elements between left and right around t
		var t = x[k]
		var i = left
		var j = right
		x[left], x[k] = x[k], x[left]
		if less(t, x[right]) {
			x[left], x[right] = x[right], x[left]
		}
		for i < j {
			x[i], x[j] = x[j], x[i]
			i++
			j--
			for i < length && less(x[i], t) {
				i++
			}
			for j >= 0 && less(t, x[j]) {
				j--
			}
		}
		if !(less(x[left], t) || less(t, x[left])) { // x[left] == t
			x[left], x[j] = x[j], x[left]
		} else {
			j++
			x[j], x[right] = x[right], x[j]
		}
		// Adjust left and right towards the boundaries of the subset
		// containing the (k − left + 1)th smallest element.
		if j <= k {
			left = j + 1
		}
		if k <= j {
			right = j - 1
		}
	}
}

func min[E constraints.Ordered](x, y E) E {
	if x < y {
		return x
	}
	return y
}

func max[E constraints.Ordered](x, y E) E {
	if x > y {
		return x
	}
	return y
}

func sign(x float64) int {
	if x < 0 {
		return -1
	}
	return 1
}
