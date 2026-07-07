# Two Pointers

## What it is

Two Pointers is a technique where two index variables traverse a data structure
(usually an array or string) — either towards each other, away from each other,
or in the same direction at different speeds. It converts many O(n²) brute-force
pair-enumeration problems into O(n) single-pass solutions.

---

## Variants

### 1. Converging pointers (opposite ends → middle)
Start one pointer at the left end, one at the right end. Move them inward based
on a comparison condition. Stop when they meet.

```
[  ←  ...  →  ]
 l             r
```

**Use when:** array is sorted; you need a pair/triplet that satisfies a sum or
difference condition.

**Template:**
```go
l, r := 0, len(arr)-1
for l < r {
    if condition(arr[l], arr[r]) {
        // found answer
    } else if needLarger {
        l++
    } else {
        r--
    }
}
```

### 2. Same-direction pointers (slow + fast)
Both pointers start at the left. The fast pointer explores ahead; the slow
pointer marks the "boundary" of a processed region.

```
[  s → → → f → ]
```

**Use when:** removing duplicates in-place, finding a cycle, sliding window
variants, partitioning.

**Template:**
```go
slow := 0
for fast := 0; fast < len(arr); fast++ {
    if shouldKeep(arr[fast]) {
        arr[slow] = arr[fast]
        slow++
    }
}
// arr[:slow] is the result
```

### 3. Two arrays / strings
One pointer per sequence, advance whichever is lagging.

**Use when:** merging two sorted arrays, comparing sequences character by character.

---

## When to recognise it

| Signal | Pointer style |
|--------|--------------|
| Sorted array + find pair with target sum | Converging |
| Find triplet / k-sum | Converging (fix one, converge on rest) |
| Palindrome check | Converging |
| Remove duplicates in-place | Slow + fast |
| Linked list cycle detection | Slow + fast (Floyd's) |
| Merge two sorted arrays | Two-array |

---

## Complexity

| | Time | Space |
|-|------|-------|
| Converging | O(n) | O(1) |
| Slow + fast | O(n) | O(1) |
| Two arrays | O(m + n) | O(1) |

The key advantage over brute force: **each pointer moves at most n steps** total,
so the whole loop is O(n) regardless of what happens inside.

---

## Common pitfalls

1. **Forgetting the sorted precondition** — converging pointers only work correctly
   when the array is sorted. Always sort first (O(n log n)) if not given sorted.
2. **Off-by-one on loop condition** — use `l < r` (strict) not `l <= r`;
   when `l == r` there is only one element left, not a valid pair.
3. **Duplicate skipping** — in 3Sum / 3Sum Closest, after recording an answer
   you must skip duplicate values of each pointer to avoid recording the same
   triplet multiple times.

---

## Problems in this repo that use Two Pointers

- [0001 — Two Sum](/0001_two_sum/README.md) — Sort + two pointers (Approach 4)
