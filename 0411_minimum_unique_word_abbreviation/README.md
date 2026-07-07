# 0411 — Minimum Unique Word Abbreviation

> LeetCode #411 · Difficulty: Hard 🔒
> **Categories:** Bit Manipulation, Backtracking, String, Bitmask

---

## Problem Statement

A string can be **abbreviated** by replacing any number of **non-adjacent, non-empty** substrings with their lengths. The lengths **should not** have leading zeros.

For example, a string such as `"substitution"` could be abbreviated as (but not limited to):

- `"s10n"` (`"s ubstitutio n"`)
- `"sub4u4"` (`"sub stit u tion"`)
- `"12"` (`"substitution"`)
- `"su3i1u2on"` (`"su bst i t u ti on"`)
- `"substitution"` (no substrings replaced)

Note that `"s55n"` (`"s ubsti tutio n"`) is not a valid abbreviation of `"substitution"` because the replaced substrings are adjacent.

The **length** of an abbreviation is the number of letters that were **not** replaced plus the number of substrings that **were** replaced. For example, the abbreviation `"s10n"` has a length of `3` (`2` letters + `1` substring) and `"su3i1u2on"` has a length of `9` (`6` letters + `3` substrings).

Given a target string `target` and an array of strings `dictionary`, return *an abbreviation of* `target` *with the **shortest possible length** such that it is **not an abbreviation of** any string in* `dictionary`. If there are multiple shortest abbreviations, return any of them.

**Example 1:**

```
Input: target = "apple", dictionary = ["blade"]
Output: "a4"
Explanation: The shortest abbreviation of "apple" is "5", but this is also an abbreviation of "blade".
The next shortest abbreviations are "a4" and "4e". "4e" is an abbreviation of blade while "a4" is not.
Hence, return "a4".
```

**Example 2:**

```
Input: target = "apple", dictionary = ["plain","amber","blade"]
Output: "1p3"
Explanation: "1p3" has length 3 and "apple" has length 5. We can see that "1p3" is not an abbreviation of any word in the dictionary.
Note that "5" is an abbreviation of both "apple" but this is not the answer since it conflicts with the dictionary.
There may be multiple valid answers such as "ap3", "a3e", "2p2", "3le", "3l1".
```

**Constraints:**

- `m == target.length`
- `n == dictionary.length`
- `1 <= m <= 21`
- `0 <= n <= 1000`
- `1 <= dictionary[i].length <= 100`
- If `n > 0`, then `log2(n) + m <= 21`.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★☆☆ Medium     | 2023          |
| Facebook   | ★★☆☆☆ Low        | 2022          |
| Amazon     | ★★☆☆☆ Low        | 2022          |
| Microsoft  | ★☆☆☆☆ Rare       | 2021          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Bit Manipulation (diff-masks & the `keep & diff` test)** — encode each same-length dictionary word as a bitmask of positions where it differs from `target`; a set of kept letters (also a bitmask) makes the abbreviation unique against that word iff their AND is non-zero → see [`/dsa/bit_manipulation.md`](/dsa/bit_manipulation.md)
- **Backtracking** — Approach 2 builds the abbreviation position-by-position (skip a run as a number, or keep the next letter), pruning branches that cannot beat the best length → see [`/dsa/backtracking.md`](/dsa/backtracking.md)
- **String manipulation** — rendering a kept-positions mask back into an abbreviation string (letters + run-length numbers) → see [`/dsa/string_algorithms.md`](/dsa/string_algorithms.md)
- **Bitmask subset enumeration** — Approach 1 iterates all `2^m` subsets of kept positions; the constraint `log2(n)+m ≤ 21` is precisely what keeps this exponential search feasible → see [`/dsa/bitmask.md`](/dsa/bitmask.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force Over All Kept-Masks | O(2^m · D) | O(D) | Direct and provably correct; feasible under `log2(n)+m ≤ 21` |
| 2 | Backtracking Abbreviation Builder | O(2^m · D) worst, pruned | O(m) | Same space of choices but prunes early; the usual interview answer |

> `m = len(target)`, `D` = number of dictionary words whose length equals `m` (the only ones that can collide).

---

## Key Reductions (shared by both approaches)

1. **Only same-length words matter.** An abbreviation encodes `target`'s exact length, so a dictionary word of a different length can never be one of its expansions. Drop them.
2. **An abbreviation ≙ a set of kept letter positions.** Whatever you don't keep collapses into run-length numbers. Encode "which positions are kept" as a bitmask `cand` over `0..m-1`.
3. **Diff-mask per word.** For a same-length word `w`, let `diff` have bit `i` set exactly where `w[i] != target[i]`. Keeping any one of those differing positions is enough to tell `target` and `w` apart, so `cand` distinguishes them **iff `cand & diff != 0`**.
4. **Valid ⇔ distinguishes every word.** `cand` is a valid abbreviation iff `cand & diff != 0` for *all* diff-masks. Among valid `cand`, minimise the **display length** (kept letters + number of runs), not the popcount.

---

## Approach 1 — Brute Force Over All Kept-Masks

### Intuition

Enumerate every possible "kept positions" bitmask `cand` from `0` to `2^m − 1`. A mask is a valid abbreviation iff it ANDs non-zero with every same-length word's diff-mask (it keeps at least one distinguishing letter for each). Compute the display length of each valid mask and remember the shortest. The constraint `log2(n) + m ≤ 21` bounds `2^m` so this exponential sweep stays in budget.

### Algorithm

1. Build diff-masks for all same-length dictionary words. If there are none, the shortest abbreviation is just the number `m` (abbreviate everything) — return `"m"`.
2. For `cand` in `0 .. 2^m − 1`:
   - `valid` iff `cand & diff != 0` for every diff-mask.
   - if valid and its `abbrevLen` beats the best so far, record `cand`.
3. Render the best mask into an abbreviation string.

### Complexity

- **Time:** O(2^m · D) — `2^m` masks, each tested against `D` same-length words.
- **Space:** O(D) — the diff-mask list (plus O(m) to render the answer).

### Code

```go
func bruteForce(target string, dictionary []string) string {
	n := len(target)
	diffs := buildDiffMasks(target, dictionary)

	// No same-length word to avoid → abbreviate everything to just the number n.
	if len(diffs) == 0 {
		return strconv.Itoa(n)
	}

	bestMask := (1 << n) - 1 // fallback: keep all letters (always valid, length n)
	bestLen := n             // its display length is exactly n
	for cand := 0; cand < (1 << n); cand++ {
		valid := true
		for _, d := range diffs {
			if cand&d == 0 { // fails to distinguish target from this word
				valid = false
				break
			}
		}
		if !valid {
			continue
		}
		if l := abbrevLen(cand, n); l < bestLen { // strictly shorter → adopt
			bestLen = l
			bestMask = cand
		}
	}
	return buildAbbrev(target, bestMask)
}
```

Supporting helpers (`abbrevLen`, `buildAbbrev`, `buildDiffMasks`) are in `main.go`.

### Dry Run

Example 1: `target = "apple"`, `dictionary = ["blade"]`. Positions `a0 p1 p2 l3 e4`.

Diff-mask of `blade` vs `apple`: `b≠a`(0), `l≠p`(1), `a≠p`(2), `d≠l`(3), `e=e`(—) → `diff = 0b01111 = 15` (bits 0,1,2,3).

| cand (kept bits) | abbrev | `cand & 15` | valid? | display length |
|------------------|--------|-------------|--------|----------------|
| `00000` (keep none) | `"5"` | `0` | **no** | — |
| `10000` (keep e, bit 4) | `"4e"` | `0` | **no** (bit 4 ∉ diff) | — |
| `00001` (keep a, bit 0) | `"a4"` | `1` | **yes** | 2 |

The first *valid* mask of length 2 is `cand = 1` → `"a4"`. No length-1 mask can be valid (only `"5"` has length 1 and it fails), so `"a4"` is optimal. Result: `"a4"` ✔

---

## Approach 2 — Backtracking Abbreviation Builder

### Intuition

Instead of enumerating raw bitmasks, **construct** the abbreviation left to right. At each position choose either to *keep* the current letter (1 token, set its bit) or to *skip* a run of `k ≥ 1` characters as a single number (1 token). Track the kept-positions mask so we can validate at the end with the same `cand & diff` test. Prune aggressively: if the running token count already reaches the best complete answer, abandon the branch — a longer or equal abbreviation can never win.

### Algorithm

1. Build diff-masks; if none, return `"m"`.
2. Recurse on `(pos, tokens, keptMask)`:
   - If `tokens >= bestLen`, prune.
   - If `pos == m`: if `keptMask` distinguishes all words and `tokens < bestLen`, record it.
   - Else, branch: for each run length `1 .. m−pos`, recurse `(pos+run, tokens+1, keptMask)`; and recurse keeping the letter `(pos+1, tokens+1, keptMask | (1<<pos))`.
3. Render the best mask.

### Complexity

- **Time:** O(2^m · D) worst case (same choice space), but the `tokens >= bestLen` prune cuts most branches; each accepted leaf costs O(D) to validate.
- **Space:** O(m) — recursion depth is at most the word length.

### Code

```go
func backtracking(target string, dictionary []string) string {
	n := len(target)
	diffs := buildDiffMasks(target, dictionary)
	if len(diffs) == 0 {
		return strconv.Itoa(n)
	}

	bestLen := n + 1 // no valid abbreviation found yet (n letters is the ceiling)
	bestMask := (1 << n) - 1

	// distinguishes reports whether keeping `mask` separates target from every word.
	distinguishes := func(mask int) bool {
		for _, d := range diffs {
			if mask&d == 0 {
				return false
			}
		}
		return true
	}

	var dfs func(pos, tokens, keptMask int)
	dfs = func(pos, tokens, keptMask int) {
		if tokens >= bestLen { // cannot improve on the best complete answer → prune
			return
		}
		if pos == n {
			// Completed an abbreviation of display length `tokens`.
			if distinguishes(keptMask) && tokens < bestLen {
				bestLen = tokens
				bestMask = keptMask
			}
			return
		}
		// Option A: abbreviate a run of length runLen as a single number token.
		for runLen := 1; pos+runLen <= n; runLen++ {
			dfs(pos+runLen, tokens+1, keptMask) // one token for the whole run
		}
		// Option B: keep the letter at pos as a literal (also one token).
		dfs(pos+1, tokens+1, keptMask|(1<<pos))
	}
	dfs(0, 0, 0)

	return buildAbbrev(target, bestMask)
}
```

### Dry Run

Example 1: `target = "apple"` (`m = 5`), `diff("blade") = 15`. We trace the branches that discover the answer (`bestLen` starts at 6).

| Call `(pos, tokens, keptMask)` | Action taken | Outcome |
|--------------------------------|--------------|---------|
| `(0, 0, 00000)` | Option A, run = 5 → `(5, 1, 00000)` | leaf: `keptMask=0`, `0 & 15 == 0` → **not** distinguishing, reject (this is `"5"`) |
| `(0, 0, 00000)` | Option B, keep `a` → `(1, 1, 00001)` | descend |
| `(1, 1, 00001)` | Option A, run = 4 → `(5, 2, 00001)` | leaf: `1 & 15 = 1 ≠ 0` → distinguishes, `tokens=2 < 6` → **record** `bestLen=2`, mask=`00001` (this is `"a4"`) |
| remaining branches | any branch reaching `tokens = 2` | pruned by `tokens >= bestLen` (bestLen now 2) |

Best mask `00001` → `buildAbbrev("apple", 1) = "a4"`. Result: `"a4"` ✔

---

## Key Takeaways

- **Turn "does abbreviation X match word Y?" into a bitmask AND.** Diff-mask per word + kept-mask for the abbreviation, and uniqueness is a single `keep & diff != 0` test per word — the crux that makes the whole search cheap per candidate.
- **Filter by length first.** Only equal-length words can collide; this alone shrinks `D` dramatically and is trivial to forget.
- **Minimise *display length*, not kept letters.** Two masks with the same popcount can have different token counts because adjacent skipped runs merge into one number — always score with the run-aware `abbrevLen`.
- **The odd constraint `log2(n) + m ≤ 21` is a hint.** It bounds `2^m · n`, signalling that an exponential-in-`m` search over subsets (or a pruned backttrack) is the intended solution.
- **Impossible inputs:** if a dictionary word equals `target` (diff-mask 0) no unique abbreviation exists; LeetCode's tests never include this, and both approaches degrade gracefully to the full word.

---

## Related Problems

- LeetCode #320 — Generalized Abbreviation (enumerate all abbreviations via bitmask)
- LeetCode #408 — Valid Word Abbreviation (check one abbreviation against a word)
- LeetCode #527 — Word Abbreviation (shortest unique abbreviations for a whole list)
- LeetCode #288 — Unique Word Abbreviation (dictionary abbreviation collisions)
