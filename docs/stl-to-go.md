# C++ STL to Go: choose by operations

| C++ | Go approach | Ordering / complexity / caveat |
| --- | --- | --- |
| `std::array<T,N>` | `[N]T` | Length is part of type; assignment copies the array |
| `std::vector<T>` | `[]T`, make/append/copy | Index O(1), append amortized O(1); descriptor copies share storage |
| `std::deque<T>` | Ring buffer or slice with head | No standard general-purpose deque; reslicing alone can retain memory |
| `std::list<T>` | `container/list` | Doubly linked node handles; O(1) known-node move/remove; Value is any |
| `std::stack<T>` | Slice Push/Pop at end | Amortized O(1) Push, O(1) Pop; clear removed pointer slots |
| `std::queue<T>` | Slice+head / ring buffer | FIFO; avoid shifting O(n) per Pop |
| `std::priority_queue<T>` | `container/heap` adapter | Push/Pop O(log n), root O(1); define Less for min/max priority |
| `std::unordered_map<K,V>` | `map[K]V` | Expected O(1) lookup; comparable keys; no stable order |
| `std::unordered_set<T>` | `map[T]struct{}` | Expected O(1) membership; no stable order |
| `std::map<K,V>` | Sorted key slice + map, or custom balanced tree | Standard Go has no ordered tree map; sorted insertion O(n) |
| `std::set<T>` | Sorted deduplicated slice, or custom balanced tree | Binary search O(log n), insertion O(n) for slice approach |
| `std::multiset<T>` | `map[T]int` counts plus sorted keys as needed | Duplicates explicit; ordering is a separate operation |
| `std::multimap<K,V>` | `map[K][]V` | Multiple values per key; does not provide key ordering |
| `std::pair<A,B>` | Named struct | Express meaning with fields; multiple returns for transient results |
| `std::tuple` | Named struct / multiple returns | No general tuple literal type |
| `std::bitset<N>` | Array/slice of uint64 | Implement bit operations; `math/bits` helps count/search bits |
| `std::lower_bound` | `sort.Search` / `slices.BinarySearch` | Sorted precondition required; define exact lower/upper-bound behavior |
| `std::sort` | `slices.Sort` / `SortFunc` | SortFunc comparator returns negative/zero/positive |
| `std::stable_sort` | `slices.SortStableFunc` | Preserves order among equal keys |
| `std::unique` | Sort if needed, then `slices.Compact` | Compact removes adjacent equal values |
| `std::find` | `slices.Index` / `Contains` | Linear search, O(n) |
| `std::accumulate` | Explicit loop | State overflow and numeric type policy |
| `std::transform` / `copy_if` | Explicit loop | Preallocate capacity when useful; handle per-item errors explicitly |

```go
set := make(map[string]struct{})
set["go"] = struct{}{}
_, present := set["go"]
delete(set, "go")

counts := map[string]int{"go": 2}
value, ok := counts["missing"] // 0, false
_ = value; _ = ok; _ = present
```

A sorted-slice set inserts by finding an index, appending room, copying the suffix right,
and writing the new value. Lookup is logarithmic, but insertion/removal is linear. It is
useful for small/read-heavy collections; it does not have red-black-tree guarantees.

For a heap adapter, its Pop method removes the final slice element; `heap.Pop` first moves
the root there and restores heap order. Calling the adapter Pop directly is incorrect.
An LRU cache combines map lookup with linked-node recency; neither structure alone provides
both efficient key lookup and efficient arbitrary recency updates.

Architecture: `required operations -> ordering + ownership -> structure -> API -> tests`.
For concurrent programs keep the same structure but assign one owner or add synchronization.
Do not reach for sync.Map as the universal replacement for map plus mutex: its operations
and intended workloads differ, and a multi-key invariant still needs coordination.

Practice: medium 01/02/03 and hard 01/07. Reference APIs: [container](https://pkg.go.dev/container),
[slices](https://pkg.go.dev/slices), [sort](https://pkg.go.dev/sort), [maps](https://pkg.go.dev/maps).
