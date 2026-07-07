# 0491 — Non-decreasing Subsequences

> LeetCode #491 · Difficulty: Medium
> **Categories:** Backtracking, Array, Hash Table, Bit Manipulation

---

## Problem Statement

Given an integer array `nums`, return *all the different possible non-decreasing subsequences of the given array with at least two elements*. You may return the answer in **any order**.

**Example 1:**

```
Input: nums = [4,6,7,7]
Output: [[4,6],[4,6,7],[4,6,7,7],[4,7],[4,7,7],[6,7],[6,7,7],[7,7]]
```

**Example 2:**

```
Input: nums = [4,4,3,2,1]
Output: [[4,4]]
```

**Constraints:**

- `1 <= nums.length <= 15`
- `-100 <= nums[i] <= 100`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★☆☆ Medium     | 2023          |
| Amazon     | ★★★☆☆ Medium     | 2023          |
| Facebook   | ★★☆☆☆ Low        | 2022          |
| Microsoft  | ★★☆☆☆ Low        | 2022          |
| Bloomberg  | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Backtracking (subset enumeration)** — we build subsequences with a take/skip DFS and undo each choice after recursing; the twist is producing *distinct* results despite duplicate values → see [`/dsa/backtracking.md`](/dsa/backtracking.md)
- **Hash Table (deduplication)** — either a global set keyed on the serialized subsequence (Approach 1) or a per-recursion-level "value already used here" set (Approach 2) removes duplicate branches → see [`/dsa/hash_map.md`](/dsa/hash_map.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Backtracking + Hash-Set Dedup | O(2ⁿ · n) | O(2ⁿ · n) | Simplest to reason about; dedup after the fact |
| 2 | Backtracking + Per-Level Used Set (Optimal) | O(2ⁿ · n) | O(n) aux | Never generates a duplicate; the standard interview answer |

Both are exponential because the *output itself* can be exponential (a strictly increasing array of length `n` has `2ⁿ − n − 1` valid subsequences). Approach 2 wins by never doing wasted work on duplicates.

---

## Approach 1 — Backtracking + Hash-Set Dedup

### Intuition

Every subsequence is the result of a sequence of "take or skip" decisions made left→right. To keep results **non-decreasing**, we only ever take `nums[i]` when it is `>=` the last value already in the path. Because `nums` can repeat (e.g. the two `7`s in `[4,6,7,7]`), two different take/skip paths can yield the *same* subsequence, so we serialize each finished subsequence into a string key and store it in a set — duplicates collapse automatically.

### Algorithm

1. Run a DFS carrying `start` (next index to consider) and `path` (subsequence so far).
2. Whenever `len(path) >= 2`, serialize `path` to a comma-joined string and insert a copy into a `seen` map keyed by that string.
3. For `i` from `start` to `n-1`: if `path` is empty **or** `nums[i] >= path[last]`, append `nums[i]`, recurse from `i+1`, then pop to backtrack.
4. Collect the map's values into the result and return.

### Complexity

- **Time:** O(2ⁿ · n) — up to `2ⁿ` subsequences reach the base case, and each costs O(n) to copy and hash.
- **Space:** O(2ⁿ · n) — the set can hold every distinct subsequence, plus O(n) recursion depth.

### Code

```go
func backtrackSetDedup(nums []int) [][]int {
	seen := map[string][]int{} // serialized path -> the path itself (dedup store)
	path := []int{}            // current subsequence under construction

	var dfs func(start int)
	dfs = func(start int) {
		if len(path) >= 2 { // valid answer: at least two elements
			key := serialize(path) // stable string key for this exact sequence
			if _, ok := seen[key]; !ok {
				cp := make([]int, len(path)) // snapshot: path is mutated later
				copy(cp, path)
				seen[key] = cp // remember this distinct subsequence
			}
		}
		for i := start; i < len(nums); i++ {
			// Only extend when the non-decreasing property is preserved.
			if len(path) == 0 || nums[i] >= path[len(path)-1] {
				path = append(path, nums[i]) // take nums[i]
				dfs(i + 1)                   // recurse on the remaining suffix
				path = path[:len(path)-1]    // undo the take (backtrack)
			}
		}
	}
	dfs(0)

	res := make([][]int, 0, len(seen))
	for _, v := range seen {
		res = append(res, v) // collect all distinct subsequences
	}
	return canonical(res)
}
```

### Dry Run

Example 1: `nums = [4,6,7,7]`. We trace the DFS; `✔` marks a subsequence recorded into `seen` (length ≥ 2). Indices in `nums` are 0:4, 1:6, 2:7, 3:7.

| Call (start, path) | Action | Recorded? |
|--------------------|--------|-----------|
| (0, []) | take 4 → | |
| (1, [4]) | take 6 → | |
| (2, [4,6]) | record `[4,6]` ✔; take 7 → | ✔ |
| (3, [4,6,7]) | record `[4,6,7]` ✔; take 7 (7≥7) → | ✔ |
| (4, [4,6,7,7]) | record `[4,6,7,7]` ✔; loop empty | ✔ |
| back to (2, [4,6]) | skip index 2 (7); take index 3 (7) → | |
| (4, [4,6,7]) | key `4,6,7` already in `seen` → no new record | — |
| back to (1, [4]) | skip 6; take 7 (idx 2) → record `[4,7]` ✔ … | ✔ |
| … continues … | yields `[4,7,7]`, `[6,7]`, `[6,7,7]`, `[7,7]` | ✔ |

After canonicalizing (sort by length, then values): `[[4,6],[4,7],[6,7],[7,7],[4,6,7],[4,7,7],[6,7,7],[4,6,7,7]]` — the 8 expected subsequences (any order is accepted).

---

## Approach 2 — Backtracking + Per-Level Used Set (Optimal)

### Intuition

Duplicates in Approach 1 come from one specific place: at a single recursion depth, choosing the *same value* from two different indices launches two branches that produce identical subsequences (picking the first `7` vs the second `7` as "the next element"). Fix it at the source — keep a tiny `used` set **local to each recursion call**, and the second time a value appears as a candidate at that level, skip it. Now every path is unique by construction: no global set, no post-dedup pass. The non-decreasing rule (`nums[i] >= path[last]`) still gates every take.

### Algorithm

1. DFS carrying `start` and `path`.
2. If `len(path) >= 2`, append a copy to `res` immediately — it is guaranteed distinct.
3. Create a fresh `used := map[int]bool{}` for this call (values already branched at this level).
4. For `i` from `start` to `n-1`: skip if `used[nums[i]]`; skip if `nums[i] < path[last]`; else mark `used[nums[i]] = true`, take, recurse from `i+1`, backtrack.

### Complexity

- **Time:** O(2ⁿ · n) worst case (all-distinct increasing input still has `~2ⁿ` real subsequences), but zero duplicate branches — strictly less work than Approach 1 on repeated values.
- **Space:** O(n) recursion depth plus O(n) per-level `used` sets — no global dedup store (output not counted).

### Code

```go
func backtrackLevelDedup(nums []int) [][]int {
	res := [][]int{}
	path := []int{}

	var dfs func(start int)
	dfs = func(start int) {
		if len(path) >= 2 {
			cp := make([]int, len(path)) // snapshot the current subsequence
			copy(cp, path)
			res = append(res, cp) // guaranteed distinct — no set needed
		}
		used := map[int]bool{} // values already used as a candidate at THIS level
		for i := start; i < len(nums); i++ {
			if used[nums[i]] {
				continue // same value already branched here → would duplicate
			}
			if len(path) > 0 && nums[i] < path[len(path)-1] {
				continue // taking it would break non-decreasing order
			}
			used[nums[i]] = true         // block this value for the rest of the level
			path = append(path, nums[i]) // take nums[i]
			dfs(i + 1)                   // explore subsequences continuing after i
			path = path[:len(path)-1]    // backtrack
		}
	}
	dfs(0)
	return canonical(res)
}
```

### Dry Run

Example 1: `nums = [4,6,7,7]`. Focus on the level where the two `7`s (indices 2 and 3) are both candidates — this is where dedup happens.

| Call (start, path) | Candidate i | `used` before | Decision |
|--------------------|-------------|---------------|----------|
| (2, [4,6]) | 2 (val 7) | {} | 7 not used, 7≥6 → take, mark used{7} |
| (3, [4,6,7]) | 3 (val 7) | {} | 7 not used, 7≥7 → take → records `[4,6,7,7]` |
| back to (2, [4,6]) | 3 (val 7) | {7} | `used[7]` true → **skip** (prevents duplicate `[4,6,7]`) |

The same per-level skip at the top level (`start=0`) means the two `7`s never both launch a branch as the first element, so `[7,7]` is produced exactly once. Result after canonicalizing: `[[4,6],[4,7],[6,7],[7,7],[4,6,7],[4,7,7],[6,7,7],[4,6,7,7]]` ✔ — identical set to Approach 1, with no duplicates ever generated.

---

## Key Takeaways

- **"Distinct subsequences/combinations with duplicates" → dedup at the recursion level, not with a global set.** The rule "at one level, do not reuse the same value twice" (`used := map[T]bool{}` per call) is the canonical trick, shared with Subsets II (#90) and Combination Sum II (#40).
- Do **not** sort the input here. Sorting would destroy subsequence order (it must respect original positions). Because we cannot sort, the usual `nums[i] == nums[i-1]` adjacency check is replaced by a per-level `used` *set*.
- The non-decreasing constraint is enforced by a single guard: only take `nums[i]` when `nums[i] >= path[last]`. Combined with the take/skip DFS, that guarantees validity for free.
- Record answers at the top of the DFS call (whenever length ≥ 2), not only at the leaves — every internal node is itself a valid subsequence.

---

## Related Problems

- LeetCode #90 — Subsets II (distinct subsets with duplicate elements; same per-level dedup)
- LeetCode #40 — Combination Sum II (skip duplicate candidates at each level)
- LeetCode #78 — Subsets (take/skip enumeration, no duplicates)
- LeetCode #46 / #47 — Permutations / Permutations II (backtracking with used-tracking)
- LeetCode #673 — Number of Longest Increasing Subsequence (non-decreasing/increasing structure)
