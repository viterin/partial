# partial

Fast, generic partial sorting algorithms in Go.

Finding the top k items in a list is often done by sorting it first. This approach is slow for large lists as time is wasted sorting irrelevant items. The partial package provides a number of generic partial sorting algorithms to find them more efficiently. The interface mirrors that of Go's slices package. See below for benchmarks comparing performance to Go's standard sort.

## Installation

```shell
go get -u github.com/viterin/partial
```

## Example usage

### Partial sort

The `partial.Sort` function takes a slice of any ordered type and the number of elements to sort:

```go
s := []int{9, 2, 5, -1, 4}
partial.Sort(s, 2)
fmt.Println(s[:2]) // [-1 2]
```

The `partial.SortFunc` function accepts a slice of any type and sorts using a custom ordering function:

```go
type Person struct {
    Name string
    Age  int
}
s := []Person{
    {"Karl", 39},
    {"Jane", 31},
    {"Bob", 45},
    {"Ann", 19},
}
partial.SortFunc(s, 2, func(a, b Person) bool { return a.Age > b.Age })
fmt.Println(s[:2]) // [{Bob 45}, {Karl 39}]
```

### Top k elements

The `partial.TopK` function places the smallest k elements of a slice in front, but only guarantees that the kth element is in sorted order. This can be significantly faster than a partial sort if the order among the elements is unimportant, especially for large k:

```go
s := []int{9, 2, 5, -1, 4}
partial.TopK(s, 3)
fmt.Println(s[:3]) // [-1 2 4] or [2 -1 4]
```

The `partial.TopKFunc` accepts a slice of any type and a custom ordering function.

## Benchmarks

### Partial sort

Sorting the first 100 elements of various sized slices (times in microseconds):

|                      | **1,000** | **10,000** | **100,000** |
|----------------------|-----------|------------|-------------|
| **sort.Slice**       | 22        | 597        | 8,237       |
| **slices.Sort**      | 9         | 344        | 5,080       |
| **slices.SortFunc**  | 21        | 495        | 7,121       |
| **partial.Sort**     | **3**     | **7**      | **107**     |
| **partial.SortFunc** | 8         | 28         | 347         |

### Top k elements

Selecting the top 50% of various sized slices (times in microseconds):

|                      | **1,000** | **10,000** | **100,000** |
|----------------------|-----------|------------|-------------|
| **slices.Sort**      | 9         | 344        | 5,080       |
| **slices.SortFunc**  | 21        | 495        | 7,121       |
| **partial.Sort**     | 6         | 194        | 3,010       |
| **partial.SortFunc** | 14        | 297        | 4,295       |
| **partial.TopK**     | **2**     | **15**     | **585**     |
| **partial.TopKFunc** | 6         | 61         | 913         |
