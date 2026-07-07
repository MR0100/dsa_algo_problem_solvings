# Hash Map (Hash Table)

## What it is

A hash map stores **key → value** pairs and provides O(1) average-time insert, lookup, and delete. Internally it uses a hash function to map keys to array slots (buckets). Collisions are handled via chaining or open addressing.

In Go: `map[KeyType]ValueType`

```go
m := make(map[int]int)
m[key] = value          // insert / update  O(1) avg
val, ok := m[key]       // lookup           O(1) avg  (ok = false if absent)
delete(m, key)          // delete           O(1) avg
```

---

## When to recognise it

Use a hash map when the problem requires any of:

| Signal in the problem | Map role |
|-----------------------|----------|
| "find if X exists in a collection" | set (map to bool or struct{}) |
| "find the index / position of X" | value → index map |
| "count occurrences of each element" | value → count map |
| "group elements by some property" | property → list of elements |
| "two-sum / complement lookup" | value → index, check complement |
| "detect duplicates in one pass" | value → bool/count |

---

## General templates

### Frequency counter
```go
freq := make(map[int]int)
for _, v := range nums {
    freq[v]++
}
```

### Existence / set
```go
seen := make(map[int]bool)
for _, v := range nums {
    if seen[v] {
        // duplicate found
    }
    seen[v] = true
}
```

### Complement lookup (Two Sum pattern)
```go
seen := make(map[int]int) // value → index
for i, v := range nums {
    complement := target - v
    if j, ok := seen[complement]; ok {
        return []int{j, i}
    }
    seen[v] = i
}
```

---

## Complexity

| Operation | Average | Worst (hash collision) |
|-----------|---------|------------------------|
| Insert | O(1) | O(n) |
| Lookup | O(1) | O(n) |
| Delete | O(1) | O(n) |
| Space | O(n) | O(n) |

Worst case is rare in practice. Go's built-in map uses a well-distributed hash and rehashes automatically.

---

## Common pitfalls

1. **Zero value confusion** — `m[key]` returns `0` / `false` / `""` if the key is absent, not an error. Always use the two-value form `val, ok := m[key]` when absence matters.
2. **Iterating while modifying** — safe in Go (unlike some languages), but the iteration order is deliberately randomised.
3. **Non-comparable key types** — slices and maps cannot be map keys. Use strings, ints, structs of comparables, or arrays.
4. **Nil map panic** — reading from a nil map is safe (returns zero value), but writing panics. Always `make` before writing.

---

## Problems in this repo that use Hash Map

- [0001 — Two Sum](/0001_two_sum/README.md) — complement lookup (value → index)
