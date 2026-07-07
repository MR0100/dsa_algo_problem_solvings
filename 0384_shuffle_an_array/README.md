# 0384 — Shuffle an Array

> LeetCode #384 · Difficulty: Medium
> **Categories:** Array, Math, Randomized, Design, Fisher–Yates

---

## Problem Statement

Given an integer array `nums`, design an algorithm to randomly shuffle the array. All permutations of the array should be **equally likely** as a result of the shuffling.

Implement the `Solution` class:

- `Solution(int[] nums)` Initializes the object with the integer array `nums`.
- `int[] reset()` Resets the array to its original configuration and returns it.
- `int[] shuffle()` Returns a random shuffling of the array.

**Example 1:**
```
Input
["Solution", "shuffle", "reset", "shuffle"]
[[[1, 2, 3]], [], [], []]
Output
[null, [3, 1, 2], [1, 2, 3], [1, 3, 2]]

Explanation
Solution solution = new Solution([1, 2, 3]);
solution.shuffle();    // Shuffle the array [1,2,3] and return its result.
                       // Any permutation of [1,2,3] must be equally likely to be returned.
                       // Example: return [3, 1, 2]
solution.reset();      // Resets the array back to its original configuration [1,2,3]. Return [1, 2, 3]
solution.shuffle();    // Returns the random shuffling of array [1,2,3]. Example: return [1, 3, 2]
```

**Constraints:**
- `1 <= nums.length <= 50`
- `-10⁶ <= nums[i] <= 10⁶`
- All the elements of `nums` are **unique**.
- At most `10⁴` calls **in total** will be made to `reset` and `shuffle`.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★★☆ High       | 2024          |
| Amazon     | ★★★☆☆ Medium     | 2024          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Facebook   | ★★☆☆☆ Low        | 2023          |
| Apple      | ★★☆☆☆ Low        | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Fisher–Yates (Knuth) Shuffle** — the in-place O(n) algorithm that produces each of n! permutations with equal probability → see [`/dsa/shuffle.md`](/dsa/shuffle.md)
- **Randomized / Sampling** — same family of uniform-random techniques as reservoir sampling → see [`/dsa/reservoir_sampling.md`](/dsa/reservoir_sampling.md)
- **Design / Data-Structure API** — stateful class keeping a pristine original for `reset` alongside `shuffle` → see [`/dsa/design_data_structures.md`](/dsa/design_data_structures.md)

---

## Approaches Overview

Let n = len(nums).

| # | Approach | Time (shuffle) | Space | When to use |
|---|----------|----------------|-------|-------------|
| 1 | Brute Force (draw from bag) | O(n) | O(n) | Intuitive baseline; conceptually the definition of a shuffle |
| 2 | Fisher–Yates (Optimal) | O(n) | O(n) result | The standard fair shuffle; least code, in-place over the copy |

---

## Approach 1 — Brute Force (Draw-from-a-bag)

### Intuition
A fair shuffle is "draw items one at a time, uniformly, without replacement". Model the bag as a mutable copy. Each round pick a random remaining index, take that element, and remove it (swap-with-last + shrink) so it can't be drawn again. Provably uniform; the removal from an arbitrary index is the part the optimal version streamlines.

### Algorithm
1. Copy the original into a `bag` slice.
2. While the bag is non-empty: pick `k = rand.Intn(len(bag))`; append `bag[k]` to the result; overwrite `bag[k]` with the last element and shrink the bag.
3. Return the result. `Reset()` returns a copy of the original.

### Complexity
- **Time:** reset O(n); shuffle O(n) — each of n draws is O(1) via swap-remove.
- **Space:** O(n) — original copy plus bag/result.

### Code
```go
func (s *BruteForceSolution) Shuffle() []int {
	bag := make([]int, len(s.original))
	copy(bag, s.original)          // mutable working copy
	res := make([]int, 0, len(bag))
	for len(bag) > 0 {
		k := s.rng.Intn(len(bag)) // pick a random remaining element
		res = append(res, bag[k]) // draw it
		bag[k] = bag[len(bag)-1]  // swap the hole with the last element
		bag = bag[:len(bag)-1]    // shrink: that element is now used
	}
	return res
}
```

### Dry Run
`original = [1,2,3]`, `bag = [1,2,3]`.

| Round | len(bag) | k | drawn | bag after (swap last into k, shrink) | res |
|-------|----------|---|-------|--------------------------------------|-----|
| 1 | 3 | 2 | 3 | `[1,2]` | `[3]` |
| 2 | 2 | 0 | 1 | `[2]` | `[3,1]` |
| 3 | 1 | 0 | 2 | `[]` | `[3,1,2]` |

Result `[3,1,2]` — one of the 6 equally likely permutations.

---

## Approach 2 — Fisher–Yates (Optimal)

### Intuition
Walk `i` from the last index down to 1. At each step pick `j` uniformly in `[0, i]` and swap `arr[i], arr[j]`. This locks in position `i` with a uniformly chosen element from the still-unfixed prefix. Because index `i` is filled from `i+1` equally likely candidates and later steps never touch it, all n! orderings occur with probability 1/n!.

### Algorithm
1. Copy the current array (so shuffle doesn't destroy the original).
2. For `i` from `n-1` down to 1: `j = rand.Intn(i+1)`; swap `arr[i], arr[j]`.
3. Return `arr`. `Reset()` returns a copy of the original.

### Complexity
- **Time:** reset O(n); shuffle O(n) — one pass, O(1) per step.
- **Space:** O(n) — the returned shuffled copy (shuffle is in-place over it).

### Code
```go
func (s *FisherYatesSolution) Shuffle() []int {
	arr := make([]int, len(s.original))
	copy(arr, s.original) // shuffle a copy, keep original intact
	for i := len(arr) - 1; i >= 1; i-- {
		j := s.rng.Intn(i + 1)          // uniform in [0, i]
		arr[i], arr[j] = arr[j], arr[i] // fix position i with a random unfixed element
	}
	return arr
}
```

### Dry Run
`arr = [1,2,3]`.

| i | range [0,i] | j (example) | swap arr[i]↔arr[j] | arr |
|---|-------------|-------------|--------------------|-----|
| 2 | [0,2] | 0 | swap idx2,0 | `[3,2,1]` |
| 1 | [0,1] | 1 | swap idx1,1 (no-op) | `[3,2,1]` |
| end | — | — | — | `[3,2,1]` |

Each of the 6 permutations of `[1,2,3]` is produced with probability 1/6.

---

## Key Takeaways
- **Fisher–Yates** is the canonical unbiased shuffle: iterate from the end, swap each element with a random one at index ≤ current. O(n) time, O(1) extra beyond the array.
- A common **bug** is picking `j` from the full range `[0, n)` every step (the "naïve shuffle") — that produces n^n outcomes, which is not divisible by n!, so it is biased. Always restrict `j` to `[0, i]`.
- Keep a pristine copy of the input so `reset` is trivial and `shuffle` never mutates the source of truth.
- Randomized correctness is validated by a **distribution test** over many trials, not a single output.

---

## Related Problems
- LeetCode #382 — Linked List Random Node (reservoir sampling sibling)
- LeetCode #398 — Random Pick Index
- LeetCode #528 — Random Pick with Weight
- LeetCode #470 — Implement Rand10() Using Rand7()
