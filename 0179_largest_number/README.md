# 0179 — Largest Number

> LeetCode #179 · Difficulty: Medium
> **Categories:** Array, String, Greedy, Sorting

---

## Problem Statement

Given a list of non-negative integers `nums`, arrange them such that they form the **largest number** and return it.

Since the result may be very large, so you need to return a string instead of an integer.

**Example 1:**

```
Input: nums = [10,2]
Output: "210"
```

**Example 2:**

```
Input: nums = [3,30,34,5,9]
Output: "9534330"
```

**Constraints:**

- `1 <= nums.length <= 100`
- `0 <= nums[i] <= 10^9`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★★ Very High  | 2024          |
| Microsoft  | ★★★★☆ High       | 2024          |
| Google     | ★★★★☆ High       | 2024          |
| Apple      | ★★★☆☆ Medium     | 2023          |
| Bloomberg  | ★★★☆☆ Medium     | 2024          |
| Salesforce | ★★☆☆☆ Low        | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Greedy (Exchange Argument)** — if neighbours satisfy `a+b < b+a`, swapping them strictly improves the result; an arrangement with no improving swap is globally optimal → see [`/dsa/greedy.md`](/dsa/greedy.md)
- **Sorting with a Custom Comparator** — the whole problem reduces to sorting under the relation "`a` before `b` iff `a+b > b+a`" → see [`/dsa/sorting.md`](/dsa/sorting.md)
- **String Algorithms** — numbers are treated as digit strings; concatenation-comparison sidesteps both overflow and the `"3" vs "30"` prefix trap → see [`/dsa/string_algorithms.md`](/dsa/string_algorithms.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (All Permutations) | O(n! · n·k) | O(n·k) | Ground truth for n ≤ ~8; validates the comparator |
| 2 | Greedy Selection | O(n² · k) | O(n·k) | Shows the greedy choice explicitly; fine for n ≤ 100 |
| 3 | Custom Comparator Sort (Optimal) | O(n log n · k) | O(n·k) | Always — the interview answer |

*(k = maximum digit count of one number, ≤ 10 here.)*

---

## Approach 1 — Brute Force (All Permutations)

### Intuition

The answer *is* some ordering of the input concatenated together — so try every ordering. One subtlety makes comparison cheap: every candidate uses exactly the same digits, so all candidates have the **same length**, and for equal-length digit strings lexicographic comparison equals numeric comparison. No big-integer arithmetic needed, just `>` on strings. Factorial growth caps this at tiny n, but it is the oracle the clever approaches must agree with.

### Algorithm

1. Convert every number to its decimal string.
2. Generate all n! permutations with swap-based backtracking: for slot `k`, try each remaining string, recurse, swap back.
3. For each complete permutation, join the strings and keep the lexicographic maximum.
4. Normalize (`"00…0"` → `"0"`) and return.

### Complexity

- **Time:** O(n! · n·k) — n! permutations, each joined and compared in O(n·k).
- **Space:** O(n·k) — the string slice and O(n) recursion depth.

### Code

```go
func bruteForce(nums []int) string {
	strs := toStrings(nums)
	best := "" // lexicographic max so far (all candidates share one length)
	var permute func(k int)
	permute = func(k int) {
		// A full permutation is fixed — evaluate its concatenation.
		if k == len(strs) {
			cand := strings.Join(strs, "")
			// Same length ⇒ lexicographic comparison == numeric comparison.
			if cand > best {
				best = cand
			}
			return
		}
		for i := k; i < len(strs); i++ {
			strs[k], strs[i] = strs[i], strs[k] // choose strs[i] for slot k
			permute(k + 1)                      // permute the remaining slots
			strs[k], strs[i] = strs[i], strs[k] // undo the choice (backtrack)
		}
	}
	permute(0)
	return normalize(best)
}
```

### Dry Run

Example 1: `nums = [10, 2]` → strs = `["10", "2"]`, 2! = 2 permutations

| Step | Permutation | Candidate | best before | cand > best? | best after |
|------|-------------|-----------|-------------|--------------|------------|
| 1 | ["10", "2"] | "102" | "" | yes | "102" |
| 2 | ["2", "10"] | "210" | "102" | yes ("2" > "1") | "210" |
| 3 | normalize("210") | — | — | no leading '0' | **"210"** |

Result: `"210"` ✔ (Example 2 enumerates 5! = 120 candidates; the max is `"9534330"` ✔)

---

## Approach 2 — Greedy Selection

### Intuition

Build the answer left to right, always asking: *which remaining number should lead?* Raw string comparison picks wrong — `"30" > "3"` lexicographically, yet 3 must come first because `"330" > "303"`. The fix: compare candidates by **what they produce**, i.e. compare the two concatenations `s+t` vs `t+s` directly. The string winning that duel against every other remaining string takes the next slot. This is selection sort powered by the concatenation duel, making the greedy choice visible one slot at a time.

### Algorithm

1. Convert numbers to strings.
2. For each output slot `pos = 0 … n−1`:
   1. Scan positions `pos+1 … n−1`; whenever `strs[i]+strs[best] > strs[best]+strs[i]`, update `best = i`.
   2. Swap the winner into `pos`.
3. Join all strings, normalize, return.

### Complexity

- **Time:** O(n² · k) — ~n²/2 duels, each comparing two ~2k-digit strings in O(k).
- **Space:** O(n·k) — the string copies; selection happens in place.

### Code

```go
func greedySelection(nums []int) string {
	strs := toStrings(nums)
	for pos := 0; pos < len(strs); pos++ {
		bestIdx := pos // assume the current occupant leads best
		for i := pos + 1; i < len(strs); i++ {
			// Does strs[i] lead better than the current best? Compare the two
			// possible concatenations instead of the raw strings.
			if strs[i]+strs[bestIdx] > strs[bestIdx]+strs[i] {
				bestIdx = i
			}
		}
		// Place this slot's winner; the rest stay in the pool.
		strs[pos], strs[bestIdx] = strs[bestIdx], strs[pos]
	}
	return normalize(strings.Join(strs, ""))
}
```

### Dry Run

Example 1: `nums = [10, 2]` → strs = `["10", "2"]`

| Step | pos | duel | winner | strs after |
|------|-----|------|--------|------------|
| 1 | 0 | "2"+"10" = "210" vs "10"+"2" = "102" → "210" wins | index 1 ("2") | ["2", "10"] |
| 2 | 1 | no rivals left | index 1 ("10") | ["2", "10"] |
| 3 | join + normalize | — | — | **"210"** |

Example 2 slot-by-slot (`[3,30,34,5,9]`): slot 0 → "9" (beats all in duels), slot 1 → "5", slot 2 → "34" (`"343" > "334"` vs "3"; `"3430" > "3034"` vs "30"), slot 3 → "3" (`"330" > "303"`), slot 4 → "30" → `"9534330"` ✔

---

## Approach 3 — Custom Comparator Sort (Optimal)

### Intuition

Approach 2's duel is a **binary relation** — so let a real sort use it. Define `a ≺ b` ("a goes first") iff `a+b > b+a`. Two facts make this correct:

1. **Exchange argument (optimality):** in any arrangement where some adjacent pair violates the relation (`a+b < b+a`), swapping that pair strictly increases the result and changes nothing else. The optimal arrangement therefore has every adjacent pair in relation order — exactly what sorting produces.
2. **Total order (safety):** the relation is transitive — treating a digit string `s` of length L as the rational value `v(s) = int(s)/(10^L − 1)` (its infinite repetition), `a+b ≥ b+a` ⇔ `v(a) ≥ v(b)`; real-number order is transitive, so the comparator is consistent and `sort.Slice` cannot misbehave.

### Algorithm

1. Convert numbers to strings.
2. `sort.Slice` descending under `strs[i]+strs[j] > strs[j]+strs[i]`.
3. Concatenate in sorted order.
4. If the result starts with `'0'`, every number was 0 → return `"0"`; else return the concatenation.

### Complexity

- **Time:** O(n log n · k) — O(n log n) comparisons, each doing O(k) string concatenation/comparison.
- **Space:** O(n·k) — the string slice (plus the sort's O(log n) stack).

### Code

```go
func customSort(nums []int) string {
	strs := toStrings(nums)
	sort.Slice(strs, func(i, j int) bool {
		// "i before j" exactly when i leading yields the bigger digit string.
		return strs[i]+strs[j] > strs[j]+strs[i]
	})
	return normalize(strings.Join(strs, ""))
}

// normalize collapses results like "00" to "0". If the largest arrangement
// starts with '0', every number must be 0 (a non-zero leader would have been
// placed first), so the whole answer is just "0".
func normalize(s string) string {
	if len(s) > 0 && s[0] == '0' {
		return "0" // all zeros — don't return "000...0"
	}
	return s
}
```

### Dry Run

Example 1: `nums = [10, 2]` → strs = `["10", "2"]`

| Step | Comparison | Result | strs |
|------|-----------|--------|------|
| 1 | "10"+"2" = "102" vs "2"+"10" = "210" | "210" bigger → "2" first | ["2", "10"] |
| 2 | join | "210" | — |
| 3 | normalize: first char '2' ≠ '0' | keep as is | **"210"** |

Example 2 key comparator calls (`[3,30,34,5,9]`): `9 ≺ 5` ("95" > "59"), `5 ≺ 34` ("534" > "345"), `34 ≺ 3` ("343" > "334"), `3 ≺ 30` ("330" > "303") → sorted `["9","5","34","3","30"]` → `"9534330"` ✔

---

## Key Takeaways

- **Compare by outcome, not by value:** for ordering puzzles where elements combine, define the comparator on the *combined result* (`a+b` vs `b+a`), never on the raw elements. The `"3" vs "30"` prefix trap is why plain descending string sort fails.
- **Exchange argument** is the standard proof template for greedy orderings: "any violating adjacent pair can be swapped for strict improvement ⇒ sorted order is optimal."
- **Custom comparators must be total orders** — Go's `sort.Slice` (like C++ `std::sort`) has undefined behaviour on inconsistent comparators. Here transitivity holds because each string behaves like the fixed repeating decimal `int(s)/(10^len − 1)`.
- **Same-length digit strings compare lexicographically = numerically** — the trick that lets the brute force avoid big integers.
- **Normalize the all-zero case** (`[0,0]` → `"0"`): if the sorted concatenation starts with '0', the maximum itself is 0.

---

## Related Problems

- LeetCode #321 — Create Maximum Number (largest number from two arrays with order constraints)
- LeetCode #402 — Remove K Digits (greedy digit arrangement, minimising instead)
- LeetCode #1754 — Largest Merge of Two Strings (the same suffix-aware greedy comparison)
- LeetCode #2165 — Smallest Value of the Rearranged Number (digit arrangement with sign/zero rules)
- LeetCode #905 — Sort Array By Parity (warm-up on sorting with a custom predicate)
