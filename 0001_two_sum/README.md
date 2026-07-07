# 0001 — Two Sum

> LeetCode #1 · Difficulty: Easy
> **Categories:** Array, Hash Map, Two Pointers, Sorting

---

## Problem Statement

Given an array of integers `nums` and an integer `target`, return **indices** of the two numbers such that they add up to `target`.

You may assume that each input would have **exactly one solution**, and you may not use the same element twice.

You can return the answer in any order.

**Example 1**
```
Input:  nums = [2,7,11,15], target = 9
Output: [0,1]
Explanation: nums[0] + nums[1] = 2 + 7 = 9, so return [0, 1].
```

**Example 2**
```
Input:  nums = [3,2,4], target = 6
Output: [1,2]
```

**Example 3**
```
Input:  nums = [3,3], target = 6
Output: [0,1]
```

**Constraints**
- `2 <= nums.length <= 10⁴`
- `-10⁹ <= nums[i] <= 10⁹`
- `-10⁹ <= target <= 10⁹`
- Only one valid answer exists.

**Follow-up:** Can you come up with an algorithm that is less than O(n²) time complexity?

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Google    | ★★★★★ Very High | 2024          |
| Amazon    | ★★★★★ Very High | 2024          |
| Meta      | ★★★★★ Very High | 2024          |
| Microsoft | ★★★★★ Very High | 2024          |
| Apple     | ★★★★☆ High      | 2024          |
| Adobe     | ★★★★☆ High      | 2024          |
| Uber      | ★★★★☆ High      | 2023          |
| Bloomberg | ★★★★☆ High      | 2024          |
| Netflix   | ★★★☆☆ Medium    | 2023          |
| LinkedIn  | ★★★☆☆ Medium    | 2023          |
| Flipkart  | ★★★☆☆ Medium    | 2023          |
| Salesforce| ★★☆☆☆ Low       | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community
> interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Array** — the input is an index-addressable sequence; returning indices is the core output requirement. → see [`/dsa/arrays.md`](/dsa/arrays.md)
- **Hash Map** — stores value→index pairs for O(1) complement lookup, turning the O(n²) brute force into O(n). → see [`/dsa/hash_map.md`](/dsa/hash_map.md)
- **Two Pointers** — after sorting, left and right pointers converge from both ends; if the sum is too small advance left, if too large retreat right. → see [`/dsa/two_pointers.md`](/dsa/two_pointers.md)
- **Sorting** — enables the two-pointer technique; costs O(n log n) but eliminates the need for extra hash map space conceptually. → see [`/dsa/sorting.md`](/dsa/sorting.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force | O(n²) | O(1) | n is tiny (< 100); no extra memory allowed |
| 2 | Two-Pass Hash Map | O(n) | O(n) | Clearest code; useful as a teaching example |
| 3 | One-Pass Hash Map ✅ | O(n) | O(n) | General case — fastest with minimal code |
| 4 | Sort + Two Pointers | O(n log n) | O(n) | When input is already sorted or space is a hard constraint |

---

## Approach 1 — Brute Force

### Intuition
Try every possible pair `(i, j)` with `i < j`. If `nums[i] + nums[j] == target`, that is the answer. No data structure needed — just raw enumeration.

### Algorithm
1. Loop `i` from `0` to `n-2`.
2. Loop `j` from `i+1` to `n-1`.
3. If `nums[i] + nums[j] == target` → return `[i, j]`.

### Complexity
- **Time:** O(n²) — for each of the n elements we scan up to n-1 elements after it.
- **Space:** O(1) — only loop counters; no extra allocation.

### Code
```go
func bruteForce(nums []int, target int) []int {
    n := len(nums)
    for i := 0; i < n-1; i++ {
        for j := i + 1; j < n; j++ {
            if nums[i]+nums[j] == target {
                return []int{i, j}
            }
        }
    }
    return nil
}
```

### Dry Run — Example 1: `nums = [2,7,11,15], target = 9`

| i | j | nums[i] | nums[j] | sum | == 9? |
|---|---|---------|---------|-----|-------|
| 0 | 1 | 2 | 7 | 9 | ✅ → return [0,1] |

Found on the very first pair.

---

## Approach 2 — Two-Pass Hash Map

### Intuition
Instead of scanning the entire remaining array for the complement of `nums[i]`, pre-load all values into a hash map. Then a second pass finds the complement in O(1) per element. Two passes, but each is O(n).

### Algorithm
1. **Pass 1:** for each `i`, insert `nums[i] → i` into `indexMap`.
2. **Pass 2:** for each `i`, compute `complement = target - nums[i]`. Look up `complement` in `indexMap`. If found at index `j` and `j != i` → return `[i, j]`.

The guard `j != i` prevents using the same element twice (e.g. `nums = [3,3], target = 6` — we must not return `[0,0]`).

### Complexity
- **Time:** O(n) — two separate linear scans.
- **Space:** O(n) — the hash map stores up to n entries.

### Code
```go
func twoPassHashMap(nums []int, target int) []int {
    indexMap := make(map[int]int)
    for i, v := range nums {
        indexMap[v] = i          // build: value → index
    }
    for i, v := range nums {
        complement := target - v
        if j, ok := indexMap[complement]; ok && j != i {
            return []int{i, j}   // found the pair
        }
    }
    return nil
}
```

### Dry Run — Example 1: `nums = [2,7,11,15], target = 9`

**Pass 1 — build map:**
| i | nums[i] | indexMap after |
|---|---------|----------------|
| 0 | 2 | {2:0} |
| 1 | 7 | {2:0, 7:1} |
| 2 | 11 | {2:0, 7:1, 11:2} |
| 3 | 15 | {2:0, 7:1, 11:2, 15:3} |

**Pass 2 — find complement:**
| i | nums[i] | complement | in map? | j | j!=i? | result |
|---|---------|------------|---------|---|-------|--------|
| 0 | 2 | 7 | ✅ | 1 | ✅ | return [0,1] |

---

## Approach 3 — One-Pass Hash Map (Optimal)

### Intuition
While inserting elements into the map, simultaneously check if the current element's complement was already inserted. If yes — the pair is complete right now, and we return immediately. Because we check **before** inserting the current element, we never pair an element with itself.

This collapses two passes into one without changing complexity — it just finds the answer earlier on average.

### Algorithm
1. For each index `i`:
   a. Compute `complement = target - nums[i]`.
   b. If `complement` is in `seen` → return `[seen[complement], i]`.
   c. Otherwise, insert `nums[i] → i` into `seen`.

### Complexity
- **Time:** O(n) — single pass; each element is processed once.
- **Space:** O(n) — the map holds at most n entries.

### Code
```go
func onePassHashMap(nums []int, target int) []int {
    seen := make(map[int]int) // value → index of previously seen elements
    for i, v := range nums {
        complement := target - v
        if j, ok := seen[complement]; ok {
            return []int{j, i}
        }
        seen[v] = i
    }
    return nil
}
```

### Dry Run — Example 1: `nums = [2,7,11,15], target = 9`

| i | v | complement | seen before check | found? | seen after insert |
|---|---|------------|-------------------|--------|-------------------|
| 0 | 2 | 7 | {} | ❌ | {2:0} |
| 1 | 7 | 2 | {2:0} | ✅ → return [0,1] | — |

### Dry Run — Example 3: `nums = [3,3], target = 6` (duplicate values)

| i | v | complement | seen | found? | seen after |
|---|---|------------|------|--------|------------|
| 0 | 3 | 3 | {} | ❌ | {3:0} |
| 1 | 3 | 3 | {3:0} | ✅ → return [0,1] | — |

The second `3` finds the first `3` in the map. No same-index issue because the first `3` was inserted **before** the second `3` is checked.

---

## Approach 4 — Sort + Two Pointers

### Intuition
In a sorted array, the two-pointer technique lets us discard half the remaining candidates at each step. Start with pointers at both ends. If the sum is too small, advance the left pointer (we need a bigger left value). If too large, retreat the right pointer (we need a smaller right value). No hash map needed.

**Catch:** the problem asks for original indices, not values. We must preserve original indices through the sort by sorting `(value, originalIndex)` pairs.

### Algorithm
1. Create `pairs = [(nums[i], i) for all i]`.
2. Sort `pairs` by value ascending.
3. Set `l = 0`, `r = n-1`.
4. While `l < r`:
   - `sum = pairs[l].val + pairs[r].val`
   - If `sum == target` → return `[pairs[l].idx, pairs[r].idx]` (sorted ascending).
   - If `sum < target` → `l++`.
   - If `sum > target` → `r--`.

### Complexity
- **Time:** O(n log n) — dominated by the sort; the two-pointer scan is O(n).
- **Space:** O(n) — the auxiliary `pairs` slice.

### Code
```go
func sortAndTwoPointers(nums []int, target int) []int {
    type pair struct{ val, idx int }
    pairs := make([]pair, len(nums))
    for i, v := range nums {
        pairs[i] = pair{v, i}
    }
    sort.Slice(pairs, func(i, j int) bool {
        return pairs[i].val < pairs[j].val
    })
    l, r := 0, len(pairs)-1
    for l < r {
        sum := pairs[l].val + pairs[r].val
        switch {
        case sum == target:
            a, b := pairs[l].idx, pairs[r].idx
            if a > b { a, b = b, a }
            return []int{a, b}
        case sum < target:
            l++
        default:
            r--
        }
    }
    return nil
}
```

### Dry Run — Example 2: `nums = [3,2,4], target = 6`

**After sort:** `pairs = [{2,1}, {3,0}, {4,2}]`

| l | r | pairs[l] | pairs[r] | sum | action |
|---|---|----------|----------|-----|--------|
| 0 | 2 | {2,1} | {4,2} | 6 | ✅ → return [1,2] |

Original indices: `pairs[l].idx=1`, `pairs[r].idx=2` → `[1,2]` ✓

---

## Key Takeaways

- **The complement trick** — transforming "find A+B=T" into "find T-A" is the core insight for nearly all two-sum-style problems. It appears in 3Sum, 4Sum, and countless variants.
- **Hash map = O(1) lookup** — any time you need to check "have I seen X before?", a hash map replaces an O(n) scan with O(1). This is one of the most frequently used tricks in array problems.
- **One-pass beats two-pass** — inserting and checking in the same loop is not just more elegant; it terminates earlier on average (the moment the second element of the pair is encountered).
- **Sort + two pointers trades time for space** — O(n log n) vs O(n), but avoids the hash map entirely. Useful when space is constrained or the array is pre-sorted.
- **Duplicate handling** — when `nums` contains duplicates (e.g. `[3,3]`), the one-pass map handles it correctly because we check before inserting. The two-pass map needs the `j != i` guard.

---

## Related Problems

- LeetCode #167 — Two Sum II (sorted array → direct two pointers, no hash map needed)
- LeetCode #15 — 3Sum (fix one element + two-sum on the rest)
- LeetCode #18 — 4Sum (fix two elements + two-sum on the rest)
- LeetCode #560 — Subarray Sum Equals K (complement trick with prefix sums)
- LeetCode #1 variants — Two Sum III (data structure design), Two Sum IV (BST)
