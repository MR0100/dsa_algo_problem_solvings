# 0471 — Encode String with Shortest Length

> LeetCode #471 · Difficulty: Hard
> **Categories:** String, Dynamic Programming, Interval DP

---

## Problem Statement

Given a string `s`, encode the string such that its encoded length is the shortest.

The encoding rule is: `k[encoded_string]`, where the `encoded_string` inside the square brackets is being repeated exactly `k` times. `k` should be a positive integer.

If an encoding process does not make the string shorter, then do not encode it. If there are several solutions, return **any of them**.

**Example 1:**

```
Input: s = "aaa"
Output: "aaa"
Explanation: There is no way to encode it such that it is shorter than the input string, so we do not encode it.
```

**Example 2:**

```
Input: s = "aaaaa"
Output: "5[a]"
Explanation: "5[a]" is shorter than "aaaaa" by 1 character.
```

**Example 3:**

```
Input: s = "aaaaaaaaaa"
Output: "10[a]"
Explanation: "a9[a]" or "9[a]a" are also valid solutions, both of them have the same length = 5, which is the same as "10[a]".
```

**Constraints:**

- `1 <= s.length <= 150`
- `s` consists of only lowercase English letters.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★☆☆ Medium     | 2023          |
| Amazon     | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2022          |
| Apple      | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Interval DP** — `dp[i][j]` = the shortest encoding of the substring `s[i..j]`; every state combines strictly shorter substrings (splits and repeated blocks), the defining shape of interval dynamic programming → see [`/dsa/interval_dp.md`](/dsa/interval_dp.md)
- **String period / repetition detection** — the `k[…]` form only applies when a substring is a repeated block; the `(s+s).indexOf(s, 1)` trick finds the smallest period in O(n) → see [`/dsa/string_algorithms.md`](/dsa/string_algorithms.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force interval recursion (no memo) | Exponential | O(n) | Shows the recurrence; recomputes substrings endlessly |
| 2 | Interval DP (top-down memoised) | O(n⁴) | O(n²) | Natural to write; each substring solved once |
| 3 | Interval DP (bottom-up, Optimal) | O(n⁴) | O(n²) | Iterative by length; no recursion, cache-friendly |

*(n ≤ 150. The `O(n⁴)` bound: `O(n²)` substrings × `O(n)` split points, plus `O(n)` period detection per state.)*

---

## Approach 1 — Brute Force Interval Recursion (No Memo)

### Intuition

The shortest encoding of a substring is the **best of three** options:

1. **Raw** — leave the substring as-is.
2. **Split** — cut it into a left and right part at some boundary and concatenate their own shortest encodings.
3. **Repeat** — if the substring is a block `P` repeated `k` times, write `k[bestEncoding(P)]`, recursively encoding the *block* (a block like `abbbabbb` may itself compress).

Take whichever yields the shortest string. Crucially, "do not encode if it isn't shorter" falls out for free: the raw form is always a candidate, so a non-shrinking `k[…]` simply loses on length (e.g. `2[abc]` and `abc abc` are both 6 characters → raw wins).

### Algorithm

1. `encode(sub)`: start `best = sub` (raw).
2. For each split `k` in `1..len-1`: `cand = encode(sub[:k]) + encode(sub[k:])`; keep the shorter.
3. If `sub` is periodic with smallest block `P` (`count = len/|P|` copies): `cand = count[encode(P)]`; keep the shorter.
4. Return `best`.

### Complexity

- **Time:** Exponential — without memoisation the same substrings are re-encoded across overlapping splits (fan-out on every boundary).
- **Space:** O(n) recursion depth (plus transient substrings).

### Code

```go
func bruteForce(s string) string {
	var encode func(sub string) string
	encode = func(sub string) string {
		if len(sub) == 0 {
			return ""
		}
		best := sub // option A: keep it raw
		// Option B: try every split into two non-empty halves.
		for k := 1; k < len(sub); k++ {
			left := encode(sub[:k])
			right := encode(sub[k:])
			if len(left)+len(right) < len(best) {
				best = left + right
			}
		}
		// Option C: if periodic, wrap the block's best encoding as count[...].
		enc := collapse(sub, encode)
		if len(enc) < len(best) {
			best = enc
		}
		return best
	}
	return encode(s)
}
```

Where `collapse` detects a repeated block and wraps it (the block is encoded via the passed-in `encode`):

```go
func collapse(s string, encodedBlock func(block string) string) string {
	n := len(s)
	doubled := (s + s)[1 : 2*n] // drop first char so index 0 can't match
	idx := strings.Index(doubled, s)
	p := idx + 1 // period length (offset back for the dropped char)
	if p < n {
		count := n / p
		block := encodedBlock(s[:p])
		return strconv.Itoa(count) + "[" + block + "]"
	}
	return s // not periodic: no run-length form
}
```

### Dry Run

Example 2: `s = "aaaaa"`. `encode("aaaaa")`:

| Candidate | value | length |
|-----------|-------|--------|
| raw | `aaaaa` | 5 |
| any split, e.g. `encode("a")+encode("aaaa")` | `a` + `aaaa` = `aaaaa` | 5 |
| repeat: period `p=1`, block `a`, count `5` → `5[encode("a")]` | `5[a]` | **4** |

`4 < 5`, so `best = "5[a]"`.

Result: `"5[a]"` ✔

---

## Approach 2 — Interval DP (Top-Down Memoised)

### Intuition

There are only `O(n²)` distinct substrings `s[i..j]`, and the best encoding of each is fixed regardless of how the recursion reaches it. Cache it. The recurrence is unchanged — raw / best split / periodic wrap — but now the block's encoding inside `collapse` is looked up by index in the memo, so the exponential fan-out becomes polynomial.

### Algorithm

1. `memo[i][j]` = best encoding of `s[i..j]`, filled lazily (empty = not yet solved).
2. `solve(i, j)`: if cached, return it. Start `best = s[i..j]` raw. For each split `k` in `i..j-1`: combine `solve(i,k) + solve(k+1,j)`. If `s[i..j]` is periodic, form `count[solve(block bounds)]`. Cache and return `best`.
3. Answer: `solve(0, n-1)`.

### Complexity

- **Time:** O(n⁴) — `O(n²)` states, each scanning `O(n)` split points and doing `O(n)` period/substring work.
- **Space:** O(n²) — the memo table (plus recursion).

### Code

```go
func memoDP(s string) string {
	n := len(s)
	if n == 0 {
		return ""
	}
	memo := make([][]string, n)
	for i := range memo {
		memo[i] = make([]string, n)
	}

	var solve func(i, j int) string
	solve = func(i, j int) string {
		if memo[i][j] != "" {
			return memo[i][j] // already computed this substring
		}
		sub := s[i : j+1]
		best := sub // raw substring
		// Split into s[i..k] + s[k+1..j] for every internal boundary k.
		for k := i; k < j; k++ {
			cand := solve(i, k) + solve(k+1, j)
			if len(cand) < len(best) {
				best = cand
			}
		}
		// Periodic wrap: encode block via the memoised solver (over block bounds).
		enc := collapseIndexed(s, i, j, solve)
		if len(enc) < len(best) {
			best = enc
		}
		memo[i][j] = best
		return best
	}
	return solve(0, n-1)
}
```

Where `collapseIndexed` runs the same period test on `s[i..j]` but resolves the block through the indexed memo:

```go
func collapseIndexed(s string, i, j int, solve func(a, b int) string) string {
	sub := s[i : j+1]
	n := len(sub)
	doubled := (sub + sub)[1 : 2*n]
	idx := strings.Index(doubled, sub)
	p := idx + 1 // smallest period length
	if p < n {
		count := n / p
		block := solve(i, i+p-1) // encode the block via the memo
		return strconv.Itoa(count) + "[" + block + "]"
	}
	return sub
}
```

### Dry Run

Example 2: `s = "aaaaa"`, `n = 5`. `solve(0,4)`:

| Sub-state consulted | result | note |
|---------------------|--------|------|
| period test on `aaaaa` | `p=1` (repeats every char) | count = 5, block bounds `(0,0)` |
| `solve(0,0)` | `a` (single char, raw) | memoised |
| wrap | `5[a]` (len 4) | beats raw `aaaaa` (len 5) |
| splits (e.g. `solve(0,0)+solve(1,4)`) | ≥ 5 chars | never shorter |

`memo[0][4] = "5[a]"`.

Result: `"5[a]"` ✔

---

## Approach 3 — Interval DP (Bottom-Up, Optimal)

### Intuition

`dp[i][j]` depends only on **strictly shorter** substrings: any split yields two shorter pieces, and the repeated block is shorter than the whole. So fill the table in order of increasing substring length — by the time we compute a length-`L` substring, all length `< L` substrings are ready, and no recursion is needed. Same three candidates per cell.

### Algorithm

1. `dp[i][i] = s[i]` (single characters).
2. For length `L = 2..n`, for each start `i` (with `j = i+L-1`):
   - `dp[i][j] = s[i..j]` raw.
   - For each split `k` in `i..j-1`: keep the shorter of `dp[i][j]` and `dp[i][k] + dp[k+1][j]`.
   - If `s[i..j]` is periodic with smallest period `p`: candidate `count[dp[i][i+p-1]]`; keep if shorter.
3. Answer: `dp[0][n-1]`.

### Complexity

- **Time:** O(n⁴) — `O(n²)` cells × `O(n)` splits (+ `O(n)` period detection per cell).
- **Space:** O(n²) — the `dp` table.

### Code

```go
func bottomUpDP(s string) string {
	n := len(s)
	if n == 0 {
		return ""
	}
	dp := make([][]string, n)
	for i := range dp {
		dp[i] = make([]string, n)
	}
	// Base case: length-1 substrings encode to themselves.
	for i := 0; i < n; i++ {
		dp[i][i] = string(s[i])
	}
	// Grow by substring length so dependencies (shorter substrings) are ready.
	for length := 2; length <= n; length++ {
		for i := 0; i+length-1 < n; i++ {
			j := i + length - 1
			sub := s[i : j+1]
			best := sub // raw
			// Best two-way split.
			for k := i; k < j; k++ {
				if len(dp[i][k])+len(dp[k+1][j]) < len(best) {
					best = dp[i][k] + dp[k+1][j]
				}
			}
			// Periodic wrap using the already-computed block cell.
			doubled := (sub + sub)[1 : 2*length]
			idx := strings.Index(doubled, sub)
			p := idx + 1 // smallest period
			if p < length {
				count := length / p
				cand := strconv.Itoa(count) + "[" + dp[i][i+p-1] + "]"
				if len(cand) < len(best) {
					best = cand
				}
			}
			dp[i][j] = best
		}
	}
	return dp[0][n-1]
}
```

### Dry Run

Example 2: `s = "aaaaa"`, `n = 5`. Base: `dp[i][i] = "a"` for all `i`. Filling by length, tracking each row's cell that stays optimal:

| length L | cell (i,j) | best split | period? | dp[i][j] |
|----------|-----------|-----------|---------|----------|
| 1 | (i,i) | — | — | `a` |
| 2 | (0,1) `aa` | `a`+`a` = `aa` (2) | p=1, `2[a]` (4) | `aa` |
| 3 | (0,2) `aaa` | `a`+`aa` = `aaa` (3) | p=1, `3[a]` (4) | `aaa` |
| 4 | (0,3) `aaaa` | `aa`+`aa` = `aaaa` (4) | p=1, `4[a]` (4) | `aaaa` |
| 5 | (0,4) `aaaaa` | `a`+`aaaa` = `aaaaa` (5) | p=1, `5[a]` (**4**) | **`5[a]`** |

`dp[0][4] = "5[a]"`.

Result: `"5[a]"` ✔

For the nested case `abbbabbbcabbbabbbc`, the block `abbbabbbc` first collapses to `2[abbb]c` (its own `dp` cell), and then the whole string, being that block twice, becomes `2[2[abbb]c]` — the interval DP composes the inner and outer encodings automatically.

---

## Key Takeaways

- **`dp[i][j]` over a substring + "compose from shorter substrings" = interval DP.** Fill by increasing length (bottom-up) or memoise by `(i,j)` (top-down); the split loop is the workhorse.
- **The raw string is always a candidate**, which is exactly how "don't encode unless shorter" is enforced — no special-casing needed. Note `2[abc]` (6) does **not** beat `abcabc` (6); the `k[]` overhead is `len(str(k)) + 2`.
- **Smallest-period trick:** `s` is a repeated block iff `(s+s).indexOf(s, 1) < len(s)`; the found index is the smallest period. This turns "is this a repetition?" into one O(n) search.
- **Recursively encode the block.** The optimal answer can be nested (`2[2[abbb]c]`), so the repeat candidate must wrap the block's *own* shortest encoding, not the raw block — which the DP gives for free by reading the block's cell.

---

## Related Problems

- LeetCode #394 — Decode String (the inverse: expand `k[...]`)
- LeetCode #664 — Strange Printer (interval DP over a string)
- LeetCode #546 — Remove Boxes (interval DP with an extra dimension)
- LeetCode #459 — Repeated Substring Pattern (the same period-detection trick)
