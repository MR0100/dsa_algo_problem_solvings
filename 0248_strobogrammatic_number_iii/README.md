# 0248 — Strobogrammatic Number III

> LeetCode #248 · Difficulty: Hard
> **Categories:** Array, String, Recursion, Math

---

## Problem Statement

Given two strings `low` and `high` that represent two integers `low` and `high` where `low <= high`, return *the number of **strobogrammatic numbers** in the range* `[low, high]`.

A **strobogrammatic number** is a number that looks the same when rotated `180` degrees (looked at upside down).

**Example 1:**

```
Input: low = "50", high = "100"
Output: 3
```

**Example 2:**

```
Input: low = "0", high = "0"
Output: 1
```

**Constraints:**

- `1 <= low.length, high.length <= 15`
- `low <= high`
- `low` and `high` consist of only digits.
- `low` and `high` do not contain any leading zeros except for zero itself.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★★☆ High       | 2024          |
| Facebook   | ★★★☆☆ Medium     | 2023          |
| Uber       | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Backtracking / Recursion** — fill a length-`L` buffer from both ends inward, choosing mirror pairs, pruning positions that break the leading-zero or center rules → see [`/dsa/backtracking.md`](/dsa/backtracking.md)
- **String Algorithms** — equal-length numeric strings compare lexicographically, so range membership is a length-then-lexicographic test → see [`/dsa/string_algorithms.md`](/dsa/string_algorithms.md)
- **Math / Number Theory** — counting objects by length class and comparing digit strings as numbers → see [`/dsa/math_number_theory.md`](/dsa/math_number_theory.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Generate All, Then Count in Range | O(Σ 5^(L/2)) | O(5^(L/2)·L) | Simple; reuses the #247 builder + a range filter |
| 2 | Recursive Count with Bounds Pruning (Optimal) | O(Σ 5^(L/2)) | O(L) depth | O(L) space; boundary-aware construction, prunes early |

---

## Approach 1 — Generate All, Then Count in Range

### Intuition

Any number in `[low, high]` has a length between `len(low)` and `len(high)`. For each such length `L`, enumerate ALL strobogrammatic numbers of that length (exactly the outside-in build from problem 247) and keep those inside `[low, high]`. Because none of these strings has a leading zero (except the single `"0"`), equal-length strings compare numerically by simple lexicographic order.

### Algorithm

1. For each length `L` from `len(low)` to `len(high)`:
   1. Generate all strobogrammatic strings of length `L` (skip leading zeros for `L>1`).
   2. Count each candidate that satisfies `inRange(cand, low, high)`.
2. Return the total.

### Complexity

- **Time:** O(Σ 5^(L/2)) over the lengths in range.
- **Space:** O(5^(L/2)·L) to hold one length's candidates.

### Code

```go
func generateAndCount(low, high string) int {
	count := 0
	for L := len(low); L <= len(high); L++ {
		for _, cand := range buildLength(L) {
			if inRange(cand, low, high) {
				count++
			}
		}
	}
	return count
}

func inRange(s, low, high string) bool {
	if len(s) < len(low) || (len(s) == len(low) && s < low) {
		return false
	}
	if len(s) > len(high) || (len(s) == len(high) && s > high) {
		return false
	}
	return true
}
```

### Dry Run

Trace `generateAndCount("50", "100")` (lengths 2 and 3):

| Step | L | candidates generated | in [50,100]? | running count |
|------|---|----------------------|--------------|---------------|
| 1 | 2 | `11,69,88,96` | 11<50 ✗; 69 ✓; 88 ✓; 96 ✓ | 3 |
| 2 | 3 | `101,111,181,609,...,986` | all ≥101 > 100 ✗ | 3 |

Result = `3`. ✓

---

## Approach 2 — Recursive Count with Bounds Pruning (Optimal)

### Intuition

Instead of materialising every candidate, fill a length-`L` byte buffer from both ends toward the middle, choosing a mirror pair at each layer. When the buffer is complete, count it if it lies in `[low, high]`. Filling both ends at once fixes both the most- and least-significant digits, and the leading-zero / odd-center rules become simple guards that prune invalid branches immediately. Uses only O(L) recursion depth.

### Algorithm

1. For each length `L` in `[len(low), len(high)]`, call `dfs(buf, 0, L-1, L)` on a fresh buffer.
2. `dfs`: if `left > right`, the buffer is complete — count it if `inRange`.
3. Otherwise for each mirror pair `(a,b)`:
   - Skip `a=='0'` when `left==0 && L>1` (leading zero).
   - Skip when `left==right && a!=b` (odd center must be self-symmetric).
   - Place `buf[left]=a`, `buf[right]=b`, recurse into `dfs(buf, left+1, right-1, L)`.

### Complexity

- **Time:** O(Σ 5^(L/2)) — same asymptotics; pruning helps in practice.
- **Space:** O(L) recursion depth plus the running count.

### Code

```go
func countStrobogrammaticInRange(low, high string) int {
	count := 0

	var dfs func(buf []byte, left, right, L int)
	dfs = func(buf []byte, left, right, L int) {
		if left > right {
			if inRange(string(buf), low, high) {
				count++
			}
			return
		}
		for _, p := range strobPairs {
			if left == 0 && L > 1 && p[0] == '0' {
				continue
			}
			if left == right && p[0] != p[1] {
				continue
			}
			buf[left] = p[0]
			buf[right] = p[1]
			dfs(buf, left+1, right-1, L)
		}
	}

	for L := len(low); L <= len(high); L++ {
		dfs(make([]byte, L), 0, L-1, L)
	}
	return count
}
```

### Dry Run

Trace `countStrobogrammaticInRange("50", "100")`, length `L=2`, `dfs(buf, 0, 1, 2)`:

| Step | pair (a,b) | guard | buffer | recurse → complete? | in [50,100]? | count |
|------|-----------|-------|--------|---------------------|--------------|-------|
| 1 | (0,0) | left==0,L>1 → **skip** | — | — | — | 0 |
| 2 | (1,1) | ok | `11` | left=1>right=0 done | 11<50 ✗ | 0 |
| 3 | (6,9) | ok | `69` | done | ✓ | 1 |
| 4 | (8,8) | ok | `88` | done | ✓ | 2 |
| 5 | (9,6) | ok | `96` | done | ✓ | 3 |

Length 3 adds nothing (all ≥101). Result = `3`. ✓

---

## Key Takeaways

- Reduce a range-count into a **per-length** enumeration; length bounds the search and equal-length strings compare numerically by lexicographic order.
- Building from both ends inward turns the leading-zero and odd-center constraints into cheap `if` guards that prune whole subtrees.
- The optimal version trades exponential storage for O(L) recursion depth — the count is accumulated, never the list.

---

## Related Problems

- LeetCode #246 — Strobogrammatic Number (validate one number)
- LeetCode #247 — Strobogrammatic Number II (generate all of a length)
- LeetCode #233 — Number of Digit One (digit-position counting in a range)
- LeetCode #357 — Count Numbers with Unique Digits (constructive counting)
