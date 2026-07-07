# 0466 — Count The Repetitions

> LeetCode #466 · Difficulty: Hard
> **Categories:** String, Dynamic Programming (cycle detection)

---

## Problem Statement

We define `str = [s, n]` as the string `str` which consists of the string `s` concatenated `n` times.

- For example, `str == ["abc", 3] == "abcabcabc"`.

We define that string `s1` can be obtained from string `s2` if we can remove some characters from `s2` such that it becomes `s1`.

- For example, `s1 = "abc"` can be obtained from `s2 = "ab**dbe**c"` based on our definition by removing the bolded underlined characters.

You are given two strings `s1` and `s2` and two integers `n1` and `n2`. You have the two strings `str1 = [s1, n1]` and `str2 = [s2, n2]`.

Return *the maximum integer* `m` *such that* `str = [str2, m]` *can be obtained from* `str1`.

**Example 1:**

```
Input: s1 = "acb", n1 = 4, s2 = "ab", n2 = 2
Output: 2
```

**Example 2:**

```
Input: s1 = "acb", n1 = 1, s2 = "acb", n2 = 1
Output: 1
```

**Constraints:**

- `1 <= s1.length, s2.length <= 100`
- `s1` and `s2` consist of lowercase English letters.
- `1 <= n1, n2 <= 10^6`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★☆☆ Medium     | 2023          |
| Amazon     | ★★☆☆☆ Low        | 2022          |
| Uber       | ★★☆☆☆ Low        | 2022          |
| ByteDance  | ★★☆☆☆ Low        | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Subsequence matching (greedy two-pointer scan)** — "`p` can be obtained from `t`" means `p` is a *subsequence* of `t`; the greedy left-to-right pointer advance is the canonical string-scan primitive → see [`/dsa/string_algorithms.md`](/dsa/string_algorithms.md)
- **Cycle detection on a finite-state machine** — the matcher's only state between `s1` copies is the current index inside `s2` (≤ 100 states), so the per-copy behaviour must repeat; detecting the cycle lets us fast-forward over the (up to 10⁶) copies of `s1` arithmetically instead of simulating each → see [`/dsa/string_algorithms.md`](/dsa/string_algorithms.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute-Force Simulation | O(n1 · \|s1\|) | O(1) | Correctness baseline; up to 10⁸ steps → can TLE on max inputs |
| 2 | Cycle Detection / Pattern Fast-Forward (Optimal) | O(\|s2\| · \|s1\|) | O(\|s2\|) | Always; independent of n1's magnitude, jumps whole cycles in O(1) |

---

## Approach 1 — Brute-Force Simulation

### Intuition

"`str2 = [s2, n2]` can be obtained from `str1`" is precisely "`s2` repeated `n2·m` times is a subsequence of `str1 = [s1, n1]`". Subsequence matching is greedy: to match a pattern `p` inside a text `t`, scan `t` left to right and advance a pointer into `p` whenever `t`'s current character equals `p[pointer]`; each time the pointer wraps past the end of `p`, one full copy of `p` has been matched. Take `p = s2`, `t = s1` repeated `n1` times, count how many complete `s2` copies fall out, then divide by `n2`.

### Algorithm

1. `j = 0` (index into `s2`), `cntS2 = 0` (completed `s2` copies).
2. Repeat `n1` times, once per copy of `s1`:
   - For each character `c` of `s1`: if `c == s2[j]`, do `j++`; when `j == len(s2)`, reset `j = 0` and `cntS2++`.
3. Return `cntS2 / n2` — how many `[s2, n2]` blocks fit into the matched copies.

### Complexity

- **Time:** O(n1 · |s1|) — visits every character of the fully expanded `str1` exactly once; with `n1 ≤ 10⁶` and `|s1| ≤ 100` that is up to 10⁸ steps.
- **Space:** O(1) — only the two integer counters `j` and `cntS2`.

### Code

```go
func bruteForce(s1 string, n1 int, s2 string, n2 int) int {
	j := 0     // current position we are trying to match inside s2
	cntS2 := 0 // how many complete copies of s2 we have matched so far
	// Expand str1 = [s1, n1] copy by copy without materialising the big string.
	for i := 0; i < n1; i++ {
		// Consume one copy of s1, advancing the s2 pointer on every match.
		for k := 0; k < len(s1); k++ {
			if s1[k] == s2[j] { // this s1 char matches the char s2 needs next
				j++ // move on to the next character of s2
				if j == len(s2) {
					j = 0     // finished a whole s2 — wrap the pointer
					cntS2++   // and record the completed copy
				}
			}
		}
	}
	// cntS2 copies of s2 = cntS2/n2 copies of [s2, n2]; integer division floors it.
	return cntS2 / n2
}
```

### Dry Run

Example 1: `s1 = "acb", n1 = 4, s2 = "ab", n2 = 2`. Pointer `j` targets `s2[j]` (`s2[0]='a'`, `s2[1]='b'`).

| s1 copy | char | s2[j] wanted | match? | j after | cntS2 after |
|---------|------|--------------|--------|---------|-------------|
| 1 | a | a | yes | 1 | 0 |
| 1 | c | b | no  | 1 | 0 |
| 1 | b | b | yes | 0 (wrap) | 1 |
| 2 | a | a | yes | 1 | 1 |
| 2 | c | b | no  | 1 | 1 |
| 2 | b | b | yes | 0 (wrap) | 2 |
| 3 | a | a | yes | 1 | 2 |
| 3 | c | b | no  | 1 | 2 |
| 3 | b | b | yes | 0 (wrap) | 3 |
| 4 | a | a | yes | 1 | 3 |
| 4 | c | b | no  | 1 | 3 |
| 4 | b | b | yes | 0 (wrap) | 4 |

`cntS2 = 4` complete copies of `s2`. Answer `= 4 / n2 = 4 / 2 = 2` ✔

---

## Approach 2 — Cycle Detection / Pattern Fast-Forward (Optimal)

### Intuition

The subsequence matcher is a finite-state machine whose **only** state between `s1` copies is *which index of `s2` we are partway through* — at most `len(s2) ≤ 100` distinct values. By the pigeonhole principle, within the first `len(s2)+1` copies of `s1` some `s2`-index must recur. From the copy where a state first appears (`start`) to the copy where it recurs (`i`), the machine forms a **cycle**: that block of `i − start` copies always yields the same number of `s2` copies. So simulate copies one at a time only until a state repeats, then jump over all the whole cycles that fit in the remaining `n1` copies with pure arithmetic, and finish by replaying the short leftover tail.

### Algorithm

1. Keep `countRec[i]` = total `s2` copies after `i` copies of `s1`, and `indexRec[i]` = the `s2` index `j` after `i` copies. Initialise index 0 to `(0, 0)`.
2. For `i = 1 … n1`: consume one copy of `s1` (advancing `j`/`cnt` exactly as in brute force), then store `countRec[i]`, `indexRec[i]`.
3. After storing, scan earlier copies for a `start < i` with `indexRec[start] == j`. On a hit:
   - `cycleLen = i − start`, `cycleCnt = countRec[i] − countRec[start]`.
   - `remaining = n1 − start`; `cyclesLeft = remaining / cycleLen`; `tail = remaining % cycleLen`.
   - `total = countRec[start] + cyclesLeft·cycleCnt + (countRec[start+tail] − countRec[start])`.
   - Return `total / n2`.
4. If no cycle is found within `n1` copies (small `n1`), return `countRec[n1] / n2`.

### Complexity

- **Time:** O(|s2| · |s1|) — at most `|s2|+1` copies are simulated before a state repeats; everything past that is O(1) arithmetic, so runtime is independent of `n1`.
- **Space:** O(|s2|) — the `countRec` / `indexRec` snapshots, bounded by the number of distinct `s2` indices before a repeat (plus the leftover tail).

### Code

```go
func cycleDetection(s1 string, n1 int, s2 string, n2 int) int {
	if n1 == 0 { // no copies of s1 at all → nothing can be matched
		return 0
	}
	// countRec[i] / indexRec[i]: state AFTER i copies of s1 have been consumed.
	// Index 0 means "before any copy": 0 completed, pointer at s2 index 0.
	countRec := make([]int, n1+1)
	indexRec := make([]int, n1+1)

	j := 0     // current index inside s2 we are matching
	cnt := 0   // total complete s2 copies matched so far
	countRec[0] = 0
	indexRec[0] = 0

	for i := 1; i <= n1; i++ { // i = number of s1 copies consumed so far
		// Consume the i-th copy of s1.
		for k := 0; k < len(s1); k++ {
			if s1[k] == s2[j] { // matched the char s2 currently needs
				j++
				if j == len(s2) {
					j = 0   // completed one full s2
					cnt++   // count it
				}
			}
		}
		countRec[i] = cnt // snapshot totals after this copy
		indexRec[i] = j

		// Look for an earlier copy that ended in the SAME s2 index j.
		for start := 0; start < i; start++ {
			if indexRec[start] == j { // state repeats → cycle from start..i
				cycleLen := i - start                        // s1 copies per cycle
				cycleCnt := countRec[i] - countRec[start]    // s2 copies per cycle
				remaining := n1 - start                      // copies left after prefix
				cyclesLeft := remaining / cycleLen           // whole cycles that fit
				tail := remaining % cycleLen                 // leftover copies to simulate

				// Prefix copies (0..start) already contribute countRec[start].
				total := countRec[start] + cyclesLeft*cycleCnt
				// The leftover `tail` copies produce the same delta the pattern
				// produced over its first `tail` copies past `start`.
				total += countRec[start+tail] - countRec[start]
				return total / n2 // how many [s2, n2] blocks fit
			}
		}
	}
	// No cycle detected within n1 copies (n1 small): plain division of the total.
	return countRec[n1] / n2
}
```

### Dry Run

Example 1: `s1 = "acb", n1 = 4, s2 = "ab", n2 = 2`.

Each copy of `s1 = "acb"` matches `a` then `b` (skipping `c`) = one full `s2`, leaving `j = 0` afterwards.

| i (copies) | j after copy i | cnt after copy i | state `j` seen before? |
|------------|----------------|------------------|-------------------------|
| 0 | 0 | 0 | — (base) |
| 1 | 0 | 1 | **yes**, at `start = 0` |

Cycle found immediately at `i = 1`, `start = 0`:

- `cycleLen = 1 − 0 = 1` copy per cycle.
- `cycleCnt = countRec[1] − countRec[0] = 1 − 0 = 1` s2 copy per cycle.
- `remaining = n1 − start = 4 − 0 = 4`; `cyclesLeft = 4 / 1 = 4`; `tail = 4 % 1 = 0`.
- `total = countRec[0] + 4·1 + (countRec[0] − countRec[0]) = 0 + 4 + 0 = 4`.

Answer `= total / n2 = 4 / 2 = 2` ✔ — same result as brute force, but the 4 copies were collapsed into one arithmetic jump.

---

## Key Takeaways

- **"Can be obtained by removing characters" = subsequence.** Reframe the wordy statement as "is `p` a subsequence of `t`?" and the greedy one-pointer scan appears immediately.
- **A bounded state that is re-entered forces a cycle.** Whenever an iterative process has only `k` possible states (here `≤ len(s2)`), its behaviour must repeat within `k+1` steps — record `(step → state)` and look for a repeat to fast-forward over enormous repetition counts (`n1` up to 10⁶) in O(1).
- **Snapshot cumulative totals, not just the state.** Storing `countRec[i]` alongside the recurring state is what lets you compute the contribution of skipped cycles *and* the leftover tail by subtraction.
- This "detect the cycle, jump the middle, replay the tail" template also solves problems like *Super Pow*, *Sum of the first n terms of a repeating decimal*, and any "repeat this transformation 10⁹ times" question.

---

## Related Problems

- LeetCode #392 — Is Subsequence (the core greedy subsequence check)
- LeetCode #686 — Repeated String Match (concatenate copies until a pattern fits)
- LeetCode #726 — Number of Atoms (parsing repeated structure)
- LeetCode #957 — Prison Cells After N Days (cycle detection to fast-forward huge N)
- LeetCode #1780 — Check if Number is a Sum of Powers of Three (state-repetition reasoning)
