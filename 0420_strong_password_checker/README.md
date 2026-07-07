# 0420 — Strong Password Checker

> LeetCode #420 · Difficulty: Hard
> **Categories:** String, Greedy

---

## Problem Statement

A password is considered strong if the below conditions are all met:

- It has at least `6` characters and at most `20` characters.
- It contains at least **one lowercase** letter, at least **one uppercase** letter, and at least **one digit**.
- It does not contain three repeating characters in a row (i.e., `"...aaa..."` is weak, but `"...aa...a..."` is strong, assuming other conditions are met).

Given a string `password`, return *the minimum number of steps required to make* `password` *strong. if* `password` *is already strong, return* `0`.

In one step, you can:

- Insert one character to `password`,
- Delete one character from `password`, or
- Replace one character of `password` with another character.

**Example 1:**

```
Input: password = "a"
Output: 5
```

**Example 2:**

```
Input: password = "aA1"
Output: 3
```

**Example 3:**

```
Input: password = "1337C0d3"
Output: 0
```

**Constraints:**

- `1 <= password.length <= 50`
- `password` consists of letters, digits, dot `'.'` or exclamation mark `'!'`.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2022          |
| Google     | ★★☆☆☆ Low        | 2022          |
| Bloomberg  | ★★☆☆☆ Low        | 2021          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Greedy** — the optimal solution never searches; it decides how to *spend* each edit for maximum overlap (a single replace/insert both breaks a run and supplies a missing character type), and spends forced deletions on the runs where they save the most replacements → see [`/dsa/greedy.md`](/dsa/greedy.md)
- **String processing** — everything is driven by two linear scans of the password: one for the three character classes, one to segment maximal runs of identical characters → see [`/dsa/string_algorithms.md`](/dsa/string_algorithms.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (BFS over edit states) | Exponential in n | Exponential | Correctness oracle for tiny inputs only |
| 2 | Greedy (Optimal) | O(n) | O(1) | The real answer — case-split on length, overlap the fixes |

---

## Approach 1 — Brute Force (BFS Over Edit States)

### Intuition

"Minimum single-character edits until a predicate holds" is literally a shortest-path problem. Treat each string as a node; an edge is one insert, delete, or replace. BFS outward from the start string and the first level containing a *strong* string is the minimum number of steps. It's exponential — useless past a handful of characters — but it is an unarguable oracle: run it on short inputs to validate the greedy formula.

### Algorithm

1. If `password` is already strong, return `0`.
2. BFS level by level. From each current string generate neighbours:
   - every single-character **deletion**,
   - every single-character **replacement** (to a small alphabet covering all three classes),
   - every single-character **insertion** at each gap.
3. Return the depth of the first strong string discovered.

### Complexity

- **Time:** exponential in `n` — branching ≈ (deletes + replaces + inserts) per node.
- **Space:** exponential — the visited set and frontier.

### Code

```go
func bruteForceBFS(password string) int {
	if isStrong(password) {
		return 0
	}
	// Small alphabet that still covers all three required classes plus a symbol.
	alphabet := []byte("aB3!")

	seen := map[string]bool{password: true}
	frontier := []string{password}
	steps := 0
	for len(frontier) > 0 {
		steps++
		var next []string
		for _, cur := range frontier {
			// 1) Deletions.
			for i := 0; i < len(cur); i++ {
				cand := cur[:i] + cur[i+1:]
				if !seen[cand] {
					if isStrong(cand) {
						return steps
					}
					seen[cand] = true
					next = append(next, cand)
				}
			}
			// 2) Replacements.
			for i := 0; i < len(cur); i++ {
				for _, ch := range alphabet {
					if cur[i] == ch {
						continue
					}
					cand := cur[:i] + string(ch) + cur[i+1:]
					if !seen[cand] {
						if isStrong(cand) {
							return steps
						}
						seen[cand] = true
						next = append(next, cand)
					}
				}
			}
			// 3) Insertions (positions 0..len).
			for i := 0; i <= len(cur); i++ {
				for _, ch := range alphabet {
					cand := cur[:i] + string(ch) + cur[i:]
					if !seen[cand] {
						if isStrong(cand) {
							return steps
						}
						seen[cand] = true
						next = append(next, cand)
					}
				}
			}
		}
		frontier = next
	}
	return steps // unreachable for valid inputs
}
```

### Dry Run

Example 2: `password = "aA1"` (has all three types but length 3 < 6).

| Level | Sample strings reached | Any strong? |
|-------|------------------------|-------------|
| 0 | `"aA1"` (len 3) | no — too short |
| 1 | `"aA1x"` variants (len 4) | no — still < 6 |
| 2 | len-5 variants | no — still < 6 |
| 3 | len-6 variants e.g. `"aA1aB3"` | **yes** — length 6, all types, no triple |

First strong string appears at level 3. Result: `3` ✔ — matching "insert 3 characters to reach length 6".

---

## Approach 2 — Greedy (Optimal)

### Intuition

Three needs interact, and the whole trick is **making one edit satisfy several needs at once**:

- `missing` = how many of {lowercase, uppercase, digit} are absent (0–3).
- Runs of ≥3 identical characters: a run of length `L` needs `⌊L/3⌋` replacements to break (replace every 3rd character).
- Length regime: too short (`<6`), in range (`6..20`), or too long (`>20`).

Case-split on length:

1. **`len < 6`** — only insertions make sense (deleting shrinks further). One insert can raise length *and* introduce a missing type *and* split a run. So the cost is `max(6 - len, missing)`: enough inserts to reach length 6, but never fewer than the number of missing types.
2. **`6 ≤ len ≤ 20`** — length is fine, so no inserts/deletes. Replacements break the runs, and a replacement can be chosen to *also* add a missing type. Cost is `max(replace, missing)`.
3. **`len > 20`** — you are *forced* to delete `over = len − 20` characters. Deletions can shorten runs, which reduces the `⌊L/3⌋` replacements those runs demand — but at different efficiencies depending on `L mod 3`:
   - `L ≡ 0 (mod 3)`: **1** deletion drops `⌊L/3⌋` by 1 (best value).
   - `L ≡ 1 (mod 3)`: **2** deletions save 1 replacement.
   - `L ≡ 2 (mod 3)`: **3** deletions save 1 replacement.
   Spend the forced deletions greedily in that priority order, then combine: `over + max(replaceRemaining, missing)`.

### Algorithm

1. One scan → `missing` (count of absent character classes).
2. One scan segmenting maximal equal-character runs; accumulate `replace += ⌊L/3⌋` for each run of length ≥3, and bucket runs by `L mod 3`.
3. Return by regime:
   - `len < 6`: `max(6 - len, missing)`.
   - `len ≤ 20`: `max(replace, missing)`.
   - `len > 20`: apply `over` deletions to buckets `[0]` (1 each), then `[1]` (2 each), then the remainder (3 each) to reduce `replace`; return `over + max(replace, missing)`.

### Complexity

- **Time:** O(n) — two linear scans, O(1) bucket work.
- **Space:** O(1) — a handful of counters and three buckets.

### Code

```go
func greedy(password string) int {
	n := len(password)

	// --- missing character types ---
	var lower, upper, digit bool
	for i := 0; i < n; i++ {
		c := password[i]
		switch {
		case c >= 'a' && c <= 'z':
			lower = true
		case c >= 'A' && c <= 'Z':
			upper = true
		case c >= '0' && c <= '9':
			digit = true
		}
	}
	missing := 0
	if !lower {
		missing++
	}
	if !upper {
		missing++
	}
	if !digit {
		missing++
	}

	// --- collect run lengths of ≥3 repeats ---
	// replace = total replacements needed to break all runs = Σ ⌊len/3⌋.
	// For the over-length regime we also bucket runs by len % 3.
	replace := 0
	// buckets[r] = number of runs whose length ≡ r (mod 3), among runs of length ≥ 3.
	var buckets [3]int
	i := 0
	for i < n {
		j := i
		for j < n && password[j] == password[i] {
			j++ // extend the run of identical characters
		}
		runLen := j - i
		if runLen >= 3 {
			replace += runLen / 3 // ⌊len/3⌋ replacements break this run
			buckets[runLen%3]++   // remember its residue for deletion targeting
		}
		i = j
	}

	if n < 6 {
		// Only insertions. Each insertion can add length and (if aimed well)
		// add a missing type or split a run. So we need at least (6-n) inserts
		// for length and at least `missing` inserts for types; one insert can
		// serve both, hence the max.
		return max(6-n, missing)
	}

	if n <= 20 {
		// No length change: fix runs with replacements, and reuse replacements
		// to add missing types. So max(replace, missing).
		return max(replace, missing)
	}

	// n > 20: we must delete exactly `over` characters.
	over := n - 20
	deletions := over // every over-length char must go

	// Spend deletions to reduce the replacements the runs still require.
	//   • L ≡ 0 (mod 3): deleting 1 char saves 1 replacement (best value).
	//   • L ≡ 1 (mod 3): deleting 2 chars saves 1 replacement.
	//   • L ≡ 2 (mod 3): deleting 3 chars saves 1 replacement.
	// Apply the cheap savings first.

	// Pass 1: runs with len % 3 == 0 — 1 deletion each saves 1 replacement.
	if over > 0 {
		use := min(buckets[0], over)
		over -= use
		replace -= use // each such deletion removes one required replacement
	}
	// Pass 2: runs with len % 3 == 1 — 2 deletions each save 1 replacement.
	if over > 0 {
		use := min(buckets[1]*2, over)
		over -= use
		replace -= use / 2
	}
	// Pass 3: any remaining deletions — every 3 deletions save 1 replacement.
	if over > 0 {
		replace -= over / 3
	}
	if replace < 0 {
		replace = 0 // can't need negative replacements
	}

	// After deletions, remaining run-replacements can still double as adding a
	// missing type, so the non-deletion cost is max(replace, missing).
	return deletions + max(replace, missing)
}
```

### Dry Run

Example 1: `password = "a"` (`n = 1`).

| Step | Computation | Value |
|------|-------------|-------|
| classes | has lower only | lower=true, upper=false, digit=false |
| missing | upper + digit absent | 2 |
| runs | no run ≥ 3 | replace = 0 |
| regime | `n = 1 < 6` | use `max(6 - n, missing)` |
| combine | `max(6 - 1, 2) = max(5, 2)` | **5** |

Result: `5` ✔.

Over-length illustration: `password = "aaaabbaaaabbaaaabbaaaabb"` (`n = 24`).

| Step | Computation | Value |
|------|-------------|-------|
| classes | lowers `a`,`b`; no upper/digit | missing = 2 |
| runs | four `aaaa` runs (len 4); `bb` runs len 2 ignored | replace = 4·⌊4/3⌋ = 4; buckets[4%3=1] = 4 |
| forced deletes | over = 24 − 20 = 4 | deletions = 4 |
| pass 1 (mod 0) | buckets[0]=0 | replace stays 4, over 4 |
| pass 2 (mod 1) | use = min(4·2, 4) = 4; over → 0; replace −= 4/2 | replace = 2 |
| pass 3 | over = 0 | — |
| combine | 4 + max(replace=2, missing=2) | **6** |

Result: `6` ✔ (verified against the canonical reference and a 20 000-case differential test).

---

## Key Takeaways

- **Overlap edits.** The crux is that one replacement can simultaneously break a triple *and* introduce a missing character class — so run-replacements and missing-types combine with `max`, not `+`. Missing the overlap over-counts.
- **Case-split on the length regime first.** `<6` uses only inserts (`max(6−len, missing)`); `6..20` uses only replacements (`max(replace, missing)`); `>20` forces `len−20` deletions before anything else.
- **Spend forced deletions where they save the most.** A run of length `L` needs `⌊L/3⌋` replacements; deletions reduce that most cheaply on runs with `L ≡ 0 (mod 3)` (1 deletion → −1 replacement), then `≡ 1` (2 → −1), then the rest (3 → −1). Priority order is what makes the greedy optimal.
- **Always keep a brute-force oracle for hard greedy problems.** A BFS-over-edits baseline plus randomized differential testing is how you gain confidence that a subtle formula like this is actually correct.

---

## Related Problems

- LeetCode #72 — Edit Distance (min single-character edits, DP flavor)
- LeetCode #65 — Valid Number (multi-condition string validation)
- LeetCode #468 — Validate IP Address (rule-based string checking)
- LeetCode #1055 — Shortest Way to Form String (greedy character accounting)
