# 0041 — First Missing Positive

> LeetCode #41 · Difficulty: Hard
> **Categories:** Array, Hash Table

---

## Problem Statement

Given an unsorted integer array `nums`, return the smallest missing positive integer.

You must implement an algorithm that runs in `O(n)` time and uses `O(1)` auxiliary space.

**Example 1**
```
Input:  nums = [1,2,0]
Output: 3
```

**Example 2**
```
Input:  nums = [3,4,-1,1]
Output: 2
```

**Example 3**
```
Input:  nums = [7,8,9,11,12]
Output: 1
```

**Constraints**
- `1 <= nums.length <= 10⁵`
- `-2³¹ <= nums[i] <= 2³¹ - 1`

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Amazon    | ★★★★★ Very High | 2024          |
| Google    | ★★★★★ Very High | 2024          |
| Microsoft | ★★★★☆ High      | 2024          |
| Meta      | ★★★★☆ High      | 2024          |
| Bloomberg | ★★★★☆ High      | 2023          |
| Apple     | ★★★☆☆ Medium    | 2023          |
| Adobe     | ★★★☆☆ Medium    | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Index as Hash / Cyclic Sort** — the key insight: place each value `v ∈ [1,n]` at index `v-1`. This turns the input array into an in-place hash map, enabling O(n)/O(1) solution.
- **Pigeonhole Principle** — the answer is always in `[1, n+1]`. If all integers 1..n are present, the answer is n+1. Otherwise some value in 1..n is missing.

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Sorting | O(n log n) | O(1) | Fails O(n) requirement |
| 2 | Hash Set | O(n) | O(n) | Fails O(1) space requirement |
| 3 | Index as Hash ✅ | O(n) | O(1) | The only approach satisfying both constraints |

---

## Approach 1 — Sorting

### Intuition
Sort the array; walk and track the expected next integer starting at 1. Simple but O(n log n).

### Complexity
- **Time:** O(n log n).
- **Space:** O(1) (in-place sort).

---

## Approach 2 — Hash Set

### Intuition
Put all values in a set; scan 1, 2, 3, … and return the first not in the set.

### Complexity
- **Time:** O(n).
- **Space:** O(n).

---

## Approach 3 — Index as Hash (Recommended ✅)

### Intuition
The answer is always in `[1, n+1]`. We only care about values in `[1, n]`. Use the array itself as a hash:
- Value `v ∈ [1, n]` belongs at index `v-1`.
- After placement, scan: the first index `i` where `nums[i] != i+1` gives the answer `i+1`.

**Placement (cyclic sort variant):** at each position `i`, while `nums[i]` is in `[1,n]` and is not already at its correct position (`nums[nums[i]-1] != nums[i]`), swap `nums[i]` with the element at its correct destination.

### Algorithm
```
for i = 0 to n-1:
  while nums[i] in [1,n] and nums[nums[i]-1] != nums[i]:
    swap(nums[i], nums[nums[i]-1])

for i = 0 to n-1:
  if nums[i] != i+1: return i+1
return n+1
```

### Complexity
- **Time:** O(n) — each element is swapped at most once (it either reaches its correct slot or gets displaced by its rightful owner). Total swaps ≤ n.
- **Space:** O(1) — in-place.

### Code
```go
func indexAsHash(nums []int) int {
    n := len(nums)
    for i := 0; i < n; i++ {
        for nums[i] > 0 && nums[i] <= n && nums[nums[i]-1] != nums[i] {
            dest := nums[i] - 1
            nums[i], nums[dest] = nums[dest], nums[i]
        }
    }
    for i := 0; i < n; i++ {
        if nums[i] != i+1 { return i+1 }
    }
    return n+1
}
```

### Dry Run — `nums = [3,4,-1,1]`
```
i=0: nums[0]=3 ∈ [1,4], nums[2]=-1≠3 → swap(nums[0],nums[2]) → [-1,4,3,1]
     nums[0]=-1 ∉ [1,4] → stop
i=1: nums[1]=4 ∈ [1,4], nums[3]=1≠4 → swap(nums[1],nums[3]) → [-1,1,3,4]
     nums[1]=1 ∈ [1,4], nums[0]=-1≠1 → swap(nums[1],nums[0]) → [1,-1,3,4]
     nums[1]=-1 ∉ [1,4] → stop
i=2: nums[2]=3 ∈ [1,4], nums[2]=3=3 → already in place → stop
i=3: nums[3]=4 ∈ [1,4], nums[3]=4=4 → already in place → stop

Scan: nums=[1,-1,3,4]
i=0: 1==1 ✓
i=1: -1≠2 → return 2 ✓
```

---

## Key Takeaways

- **Answer is always in [1, n+1]** — this bounds the search and is why using indices 0..n-1 works perfectly.
- **Guard `nums[nums[i]-1] != nums[i]`** — without this, duplicate values in [1,n] cause infinite swapping (e.g., `[1,1]` would endlessly swap index 0 and 1).
- **This is cyclic sort** — the same "place elements at their natural index" technique solves several "find missing/duplicate" problems (#268, #287, #442, #448).
- **Negation variant** — an alternative O(n)/O(1) approach marks values by negating `nums[v-1]` when `v ∈ [1,n]`, then scans for the first positive. Cyclic sort is arguably cleaner.

---

## Related Problems

- LeetCode #268 — Missing Number (find missing in [0,n])
- LeetCode #287 — Find the Duplicate Number (find one duplicate in [1,n])
- LeetCode #442 — Find All Duplicates in an Array (all duplicates in [1,n])
- LeetCode #448 — Find All Numbers Disappeared in an Array
