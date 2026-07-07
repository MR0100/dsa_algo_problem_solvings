# 0293 — Flip Game

> LeetCode #293 · Difficulty: Easy
> **Categories:** String, Simulation

---

## Problem Statement

You are playing a Flip Game with your friend.

You are given a string `currentState` that contains only `'+'` and `'-'`. You and your friend take turns to flip **two consecutive** `"++"` into `"--"`. The game ends when a person can no longer make a move, and therefore the other person will be the winner.

Return all possible states of the string `currentState` after **one valid move**. You may return the answer in **any order**. If there is no valid move, return an empty list `[]`.

**Example 1:**
```
Input: currentState = "++++"
Output: ["--++","+--+","++--"]
```

**Example 2:**
```
Input: currentState = "+"
Output: []
```

**Constraints:**
- `1 <= currentState.length <= 500`
- `currentState[i]` is either `'+'` or `'-'`.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★☆☆☆ Low        | 2022          |
| Amazon     | ★★☆☆☆ Low        | 2022          |
| Microsoft  | ★☆☆☆☆ Rare       | 2021          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used
- **String Algorithms** — scanning for a pattern and building new strings by slicing → see [`/dsa/string_algorithms.md`](/dsa/string_algorithms.md)
- **Two Pointers (adjacent pair scan)** — inspect `s[i]` and `s[i+1]` while sweeping → see [`/dsa/two_pointers.md`](/dsa/two_pointers.md)

---

## Approaches Overview
| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Byte Slice Scan | O(n²) | O(n)/result | Explicit mutable copy |
| 2 | Slice Concatenation (Optimal) | O(n²) | O(n)/result | Tightest, clearest |

---

## Approach 1 — Byte Slice Scan

### Intuition
A move flips a consecutive `"++"` into `"--"`. Every position where `s[i]=='+'` and `s[i+1]=='+'` is one legal move. For each such spot, produce the resulting string. Copying `s` into a fresh byte slice per move keeps the mutation obvious, at the cost of one allocation per result.

### Algorithm
1. For `i = 0..len-2`:
2. If `s[i]=='+'` and `s[i+1]=='+'`:
3. Copy `s` into a byte slice, set positions `i` and `i+1` to `'-'`, append the string.
4. Return the collected results.

### Complexity
- **Time:** O(n²) — up to `n` moves, each rebuilds an O(n) string.
- **Space:** O(n) per result string (the total output is inherently O(n²)).

### Code
```go
func byteSliceScan(s string) []string {
	res := []string{}
	for i := 0; i+1 < len(s); i++ {
		if s[i] == '+' && s[i+1] == '+' { // a flippable "++" pair
			b := []byte(s) // fresh mutable copy of the whole string
			b[i] = '-'     // flip left plus
			b[i+1] = '-'   // flip right plus
			res = append(res, string(b))
		}
	}
	return res
}
```

### Dry Run
`s = "++++"`

| i | s[i] s[i+1] | flippable? | copy → flip i,i+1 | result added |
|---|-------------|-----------|-------------------|--------------|
| 0 | `+ +` | yes | `--++` | `--++` |
| 1 | `+ +` | yes | `+--+` | `+--+` |
| 2 | `+ +` | yes | `++--` | `++--` |

Result: `["--++","+--+","++--"]` ✓

---

## Approach 2 — Slice Concatenation (Optimal)

### Intuition
Only the two characters at `i` and `i+1` change. So the result is exactly `s[:i] + "--" + s[i+2:]` — no per-character copy, just three slices glued together. Same asymptotics, tighter code.

### Algorithm
1. For `i = 0..len-2` where `s[i]==s[i+1]=='+'`:
2. Append `s[:i] + "--" + s[i+2:]`.
3. Return results.

### Complexity
- **Time:** O(n²) — `n` candidate positions, each builds an O(n) string.
- **Space:** O(n) per output string.

### Code
```go
func sliceConcat(s string) []string {
	res := []string{}
	for i := 0; i+1 < len(s); i++ {
		if s[i] == '+' && s[i+1] == '+' { // flippable pair at i, i+1
			res = append(res, s[:i]+"--"+s[i+2:])
		}
	}
	return res
}
```

### Dry Run
`s = "++++"`

| i | flippable? | `s[:i]` | `"--"` | `s[i+2:]` | result |
|---|-----------|---------|--------|-----------|--------|
| 0 | yes | `""` | `--` | `++` | `--++` |
| 1 | yes | `+` | `--` | `+` | `+--+` |
| 2 | yes | `++` | `--` | `""` | `++--` |

Result: `["--++","+--+","++--"]` ✓

---

## Key Takeaways
- To enumerate one-move states, scan **adjacent pairs** and produce a spliced copy per valid position.
- Prefer `s[:i] + mid + s[i+2:]` over a full byte-copy when only a fixed window changes — clearer and no manual mutation.
- Output size alone forces O(n²): up to `n/2` results, each length `n`.

---

## Related Problems
- LeetCode #294 — Flip Game II (can the first player win?)
- LeetCode #17 — Letter Combinations of a Phone Number (enumerate states)
- LeetCode #46 — Permutations (enumerate all configurations)
