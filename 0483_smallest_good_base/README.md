# 0483 — Smallest Good Base

> LeetCode #483 · Difficulty: Hard
> **Categories:** Math, Binary Search

---

## Problem Statement

Given an integer `n` represented as a string, return *the smallest **good base** of* `n`.

We call `k >= 2` a **good base** of `n`, if all digits of `n` base `k` are `1`'s.

**Example 1:**

```
Input: n = "13"
Output: "3"
Explanation: 13 base 3 is 111.
```

**Example 2:**

```
Input: n = "4681"
Output: "8"
Explanation: 4681 base 8 is 11111.
```

**Example 3:**

```
Input: n = "1000000000000000000"
Output: "999999999999999999"
Explanation: 1000000000000000000 base 999999999999999999 is 11.
```

**Constraints:**

- `n` is an integer in the range `[3, 10^18]`.
- `n` does not contain any leading zeros.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★☆☆ Medium     | 2023          |
| Amazon     | ★★☆☆☆ Low        | 2022          |
| Uber       | ★★☆☆☆ Low        | 2022          |
| Microsoft  | ★☆☆☆☆ Rare       | 2021          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Math / number theory (geometric series & repunits)** — a good base means `n = 1 + k + k² + … + k^m`, a base-`k` *repunit*. The whole problem is reasoning about this geometric sum: bounding the digit count by `log₂ n`, and using the dominant term `k^m ≈ n` to pin the base → see [`/dsa/math_number_theory.md`](/dsa/math_number_theory.md)
- **Binary search** — for a fixed digit count `m+1`, the repunit `f(k) = 1 + k + … + k^m` is strictly increasing in `k`, so the base that makes `f(k) = n` is found by binary search over `k` → see [`/dsa/binary_search.md`](/dsa/binary_search.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force Over Bases | O(n · log n) | O(1) | Tiny `n` only; scans every base 2..n-1 (TLE near 10^18) |
| 2 | Binary Search on the Base per Digit-Count | O(log³ n) | O(1) | Robust, no reliance on floating-point roots |
| 3 | Direct m-th Root Estimate (Optimal) | O(log² n) | O(1) | Fewest ops: one candidate base per digit count |

> Key structural fact used by 2 & 3: **smaller base ⇔ more digits**, so scanning digit count from large to small and returning the first match yields the smallest base. `k = n-1` (digits `"11"`) is a guaranteed fallback.

---

## Approach 1 — Brute Force Over Bases

### Intuition

"Smallest good base" literally asks for the least `k ≥ 2` such that `n` written in base `k` is all 1's. So test `k = 2, 3, 4, …` and return the first that works. Verifying a base is digit-peeling: while `cur > 0`, every remainder `cur % k` must equal 1, then `cur /= k`. A single non-1 remainder disqualifies the base. `k = n-1` always works (`n = 11` in base `n-1`), so we never run off the end.

### Algorithm

1. For `k = 2 … n-2`:
   a. `cur = n`. While `cur > 1`: if `cur % k != 1`, break (fail); else `cur /= k`.
   b. If the loop drained `cur` to exactly `1`, base `k` is good → return it.
2. Fallback: return `n-1`.

### Complexity

- **Time:** O(n · log n) — up to ~`n` candidate bases, each checked in `O(log_k n)` divisions. Fine for small `n`; catastrophic near `10^18` (hence shown only on small inputs).
- **Space:** O(1) — a couple of scalars.

### Code

```go
func bruteForce(nStr string) string {
	n, _ := strconv.ParseUint(nStr, 10, 64)
	for k := uint64(2); k < n-1; k++ {
		cur := n
		good := true
		for cur > 1 { // peel digits of n in base k, low to high
			if cur%k != 1 { // every digit must be exactly 1
				good = false
				break
			}
			cur /= k // drop the digit we just verified
		}
		if good && cur == 1 { // consumed everything and ended on the leading 1
			return strconv.FormatUint(k, 10)
		}
	}
	return strconv.FormatUint(n-1, 10) // n = "11" in base n-1 always works
}
```

### Dry Run

Example 1: `n = 13`. Test bases upward.

| k | digit-peel of 13 in base k | all 1's? | result |
|---|-----------------------------|----------|--------|
| 2 | 13%2=1 → 6; 6%2=0 ✗ | no | reject |
| 3 | 13%3=1 → 4; 4%3=1 → 1; stop (cur=1) | yes | **return 3** |

`13 = 111₃` ✔

---

## Approach 2 — Binary Search on the Base per Digit-Count

### Intuition

Fix the number of digits `m+1`. Then `f(k) = 1 + k + … + k^m` is **strictly increasing** in `k`, so at most one base satisfies `f(k) = n`, and we can binary-search it. Because the smallest base corresponds to the **most** digits, scan `m` from its maximum down and return the first base found — it necessarily uses the maximum digit count and is therefore minimal. Bounds: since `n ≥ 2^m`, the top exponent `m ≤ ⌊log₂ n⌋`; and since `n > k^m`, the base satisfies `k < n^(1/m)`, giving a tight search interval. All repunit comparisons are done **overflow-safe** (bail out the moment a partial sum exceeds `n` or 64 bits).

### Algorithm

1. `maxM = ⌊log₂ n⌋` (largest possible top exponent).
2. For `m = maxM` down to `1`:
   - `lo = 2`, `hi = ⌊n^(1/m)⌋ + 1`.
   - Binary-search `k`: compare `repunit(k, m)` to `n` with the overflow-safe comparator; return `k` on an exact match, move `lo`/`hi` otherwise.
3. If nothing matched, return `n-1` (the two-digit `m = 1` case).

### Complexity

- **Time:** O(log³ n) — ~`log n` values of `m`, each an `O(log n)` binary search whose comparisons cost `O(m) = O(log n)`. In practice ≤ a few thousand operations.
- **Space:** O(1).

### Code

```go
func binarySearch(nStr string) string {
	n, _ := strconv.ParseUint(nStr, 10, 64)

	maxM := int(math.Log2(float64(n))) // most possible digits minus one
	for m := maxM; m >= 1; m-- {       // more digits first ⇒ smaller base first
		lo := uint64(2)
		hi := uint64(math.Pow(float64(n), 1.0/float64(m))) + 1 // k < n^(1/m), +1 for float slack
		for lo <= hi {
			mid := lo + (hi-lo)/2
			switch repunitCmp(mid, m, n) {
			case 0:
				return strconv.FormatUint(mid, 10) // exact repunit → smallest base found
			case -1:
				lo = mid + 1 // repunit too small → need a bigger base
			default:
				hi = mid - 1 // repunit too big → need a smaller base
			}
		}
	}
	return strconv.FormatUint(n-1, 10) // guaranteed two-digit fallback
}
```

Supporting overflow-safe comparator:

```go
func repunitCmp(k uint64, m int, target uint64) int {
	var sum uint64 = 1  // the k^0 = 1 term
	var term uint64 = 1 // running k^i, starts at k^0
	for i := 1; i <= m; i++ {
		hi, lo := bits.Mul64(term, k) // 128-bit product to catch overflow
		if hi != 0 {                  // needs > 64 bits → too big
			return 1
		}
		term = lo
		if sum > math.MaxUint64-term { // sum would overflow → too big
			return 1
		}
		sum += term
		if sum > target { // already past goal
			return 1
		}
	}
	if sum < target {
		return -1
	}
	return 0
}
```

### Dry Run

Example 1: `n = 13`, `maxM = ⌊log₂ 13⌋ = 3`. Scan `m` high→low.

| m | hi = ⌊13^(1/m)⌋+1 | binary search over k | repunit found? |
|---|--------------------|-----------------------|----------------|
| 3 | ⌊13^(1/3)⌋+1 = 2+1 = 3 | k=2: 1+2+4+8=15 > 13 (hi→1); k range empties | no |
| 2 | ⌊13^(1/2)⌋+1 = 3+1 = 4 | k=3: 1+3+9=13 == 13 ✓ | **return 3** |

`13 = 111₃` ✔ (found at `m = 2`, i.e. 3 digits)

---

## Approach 3 — Direct m-th Root Estimate (Optimal)

### Intuition

Skip the inner binary search entirely. For `m+1` digits, `n = k^m + k^(m-1) + … + 1` is dominated by `k^m`, so the base is essentially forced: `k = ⌊n^(1/m)⌋`. That is the *only* base that could produce `m+1` digits, so we compute it directly with a real `m`-th root and verify that single candidate's exact repunit (guarding against the floating-point estimate being off by one). Scanning `m` from large to small and returning the first candidate that verifies gives the smallest base with just **one** check per digit count.

### Algorithm

1. For `m = ⌊log₂ n⌋` down to `2`:
   - `k = ⌊n^(1/m)⌋` (the mandatory base for `m+1` digits).
   - If `k < 2`, skip; else verify `repunit(k, m) == n` (overflow-safe). If it matches, return `k`.
2. Fallback: return `n-1` (the `m = 1` / `"11"` case).

### Complexity

- **Time:** O(log² n) — ~`log n` values of `m`, each an `O(m) = O(log n)` verification. No inner search.
- **Space:** O(1).

### Code

```go
func mthRootEstimate(nStr string) string {
	n, _ := strconv.ParseUint(nStr, 10, 64)

	// m is the top exponent; digit count is m+1. Largest m first ⇒ smallest base.
	for m := int(math.Log2(float64(n))); m >= 2; m-- {
		// Dominant-term estimate: n ≈ k^m ⇒ k ≈ n^(1/m).
		k := uint64(math.Pow(float64(n), 1.0/float64(m)))
		if k < 2 {
			continue // a valid base must be ≥ 2
		}
		// The float root can be off by one; verify the true repunit exactly.
		if repunitCmp(k, m, n) == 0 {
			return strconv.FormatUint(k, 10) // first verified ⇒ maximum digits ⇒ min base
		}
	}
	return strconv.FormatUint(n-1, 10) // m = 1: n written as "11" in base n-1
}
```

### Dry Run

Example 1: `n = 13`, start `m = ⌊log₂ 13⌋ = 3`.

| m | k = ⌊13^(1/m)⌋ | repunit 1+k+…+k^m | == 13? | action |
|---|-----------------|--------------------|--------|--------|
| 3 | ⌊13^(1/3)⌋ = 2 | 1+2+4+8 = 15 | no | continue |
| 2 | ⌊13^(1/2)⌋ = 3 | 1+3+9 = 13 | yes | **return 3** |

`13 = 111₃` ✔ — one root estimate + one verify per `m`, no binary search.

---

## Key Takeaways

- **Repunit reframing**: "all digits are 1 in base `k`" ⇔ `n = 1 + k + … + k^m`. Recognising the geometric series is the whole unlock.
- **Smallest base ⇔ most digits.** Iterate digit count from large to small; the first feasible base is automatically the minimum, so no global comparison is needed.
- **Bound the exponent by `log₂ n`.** With `k ≥ 2`, `n ≥ 2^m`, so there are only ~60 digit counts to try for `n ≤ 10^18` — turning a hopeless search into a handful of checks.
- **Dominant-term estimate `k ≈ n^(1/m)`** collapses the inner binary search to a single candidate; always *verify* it exactly because `math.Pow` can be off by one.
- **Overflow discipline**: with `n` up to `10^18`, `k^m` overflows `uint64` easily. Use `bits.Mul64` (128-bit product) and early bail-out instead of computing then comparing.

---

## Related Problems

- LeetCode #50 — Pow(x, n) (fast exponentiation, same power arithmetic)
- LeetCode #69 — Sqrt(x) (integer root via binary search)
- LeetCode #372 — Super Pow (modular exponentiation on big values)
- LeetCode #1281 — Subtract the Product and Sum of Digits (base/digit reasoning)
