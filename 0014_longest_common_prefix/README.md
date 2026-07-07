# 0014 — Longest Common Prefix

> LeetCode #14 · Difficulty: Easy
> **Categories:** String, Trie, Binary Search

---

## Problem Statement

Write a function to find the longest common prefix string amongst an array of strings.

If there is no common prefix, return an empty string `""`.

**Example 1**
```
Input:  strs = ["flower","flow","flight"]
Output: "fl"
```

**Example 2**
```
Input:  strs = ["dog","racecar","car"]
Output: ""
Explanation: There is no common prefix among the input strings.
```

**Constraints**
- `1 <= strs.length <= 200`
- `0 <= strs[i].length <= 200`
- `strs[i]` consists of only lowercase English letters.

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Amazon    | ★★★★★ Very High | 2024          |
| Google    | ★★★★☆ High      | 2024          |
| Microsoft | ★★★★☆ High      | 2024          |
| Bloomberg | ★★★★☆ High      | 2024          |
| Meta      | ★★★☆☆ Medium    | 2023          |
| Apple     | ★★★☆☆ Medium    | 2023          |
| Adobe     | ★★☆☆☆ Low       | 2022          |
| Uber      | ★★☆☆☆ Low       | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **String** — all approaches iterate over string characters by position.
- **Sorting** — Approach 3 leverages lexicographic sort so that only the extreme elements need comparison.
- **Binary Search** — Approach 4 binary-searches the prefix length, exploiting the monotonic property: "if a prefix of length L is common, all shorter prefixes are also common."

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Horizontal Scan | O(S) | O(1) | Simple; good when many strings have long common prefixes |
| 2 | Vertical Scan ✅ | O(S) | O(1) | Best for short common prefixes; early exit column by column |
| 3 | Sort + Compare Endpoints | O(n log n + m) | O(1) | Elegant; only compares two strings after sort |
| 4 | Binary Search on Length | O(S log m) | O(1) | Shows binary search pattern; slight overhead from repeated checks |

S = total characters across all strings. m = length of shortest string. All approaches are O(S) worst case; sorting adds O(n log n).

---

## Approach 1 — Horizontal Scanning

### Intuition
Start with `strs[0]` as the candidate prefix. For each subsequent string, shrink the prefix from the right until it matches the start of that string.

### Algorithm
1. `prefix = strs[0]`.
2. For each `s` in `strs[1:]`:
   - While `!strings.HasPrefix(s, prefix)`: `prefix = prefix[:len(prefix)-1]`.
   - If `prefix == ""` → return `""`.
3. Return `prefix`.

### Complexity
- **Time:** O(S) — in the worst case (all strings equal) we scan every character.
- **Space:** O(1).

### Code
```go
// horizontalScan starts with strs[0] as the prefix and progressively trims it
// until it is a prefix of every subsequent string.
//
// Time:  O(S) where S = sum of all character lengths in strs.
// Space: O(1) extra beyond the output.
func horizontalScan(strs []string) string {
	if len(strs) == 0 {
		return ""
	}
	prefix := strs[0]
	for _, s := range strs[1:] {
		// Trim prefix from the right until it matches s's beginning.
		for !strings.HasPrefix(s, prefix) {
			prefix = prefix[:len(prefix)-1]
			if prefix == "" {
				return ""
			}
		}
	}
	return prefix
}
```

### Dry Run — `strs = ["flower","flow","flight"]`
`prefix` starts as `strs[0]` and is trimmed from the right until it prefixes each next string.

| next string | prefix before | trimming steps (HasPrefix?) | prefix after |
|-------------|---------------|-----------------------------|--------------|
| `"flow"`    | `"flower"`    | "flow" HasPrefix "flower"? no → "flowe" no → "flow" yes | `"flow"` |
| `"flight"`  | `"flow"`      | "flight" HasPrefix "flow"? no → "flo" no → "fl" yes | `"fl"` |

All strings consumed → return `"fl"` ✓

---

## Approach 2 — Vertical Scanning (Recommended ✅)

### Intuition
Compare character by character across all strings at the same column index. The moment any string is too short or has a mismatching character, stop. This short-circuits immediately at the first differing column, making it better than horizontal scan when the common prefix is short.

### Algorithm
1. For each column `i` from 0 to `len(strs[0])-1`:
   - For each string `s`: if `i >= len(s)` or `s[i] != strs[0][i]` → return `strs[0][:i]`.
2. Return `strs[0]`.

### Complexity
- **Time:** O(S) worst case; faster in practice for short common prefixes.
- **Space:** O(1).

### Code
```go
func verticalScan(strs []string) string {
    if len(strs) == 0 { return "" }
    for i := 0; i < len(strs[0]); i++ {
        ch := strs[0][i]
        for _, s := range strs[1:] {
            if i >= len(s) || s[i] != ch { return strs[0][:i] }
        }
    }
    return strs[0]
}
```

### Dry Run — `strs = ["flower","flow","flight"]`
```
i=0: ch='f'. "flow"[0]='f' ✓, "flight"[0]='f' ✓
i=1: ch='l'. "flow"[1]='l' ✓, "flight"[1]='l' ✓
i=2: ch='o'. "flow"[2]='o' ✓, "flight"[2]='i' ✗ → return strs[0][:2] = "fl" ✓
```

---

## Approach 3 — Sort + Compare First and Last

### Intuition
After lexicographic sort, the first and last strings are maximally different. The LCP of the entire set equals the LCP of just these two: any character position where `first[i] != last[i]` would also differ in at least one intermediate pair.

### Algorithm
1. `sort.Strings(strs)`.
2. Compare `strs[0]` and `strs[n-1]` character by character.

### Complexity
- **Time:** O(n log n) for sort + O(m) for comparison.
- **Space:** O(1) extra (sort is in-place in Go).

### Code
```go
// sortEndpoints sorts the slice lexicographically and compares only the
// first and last strings — the LCP of all strings equals the LCP of these two.
//
// Time:  O(n log n) for the sort + O(m) for the comparison, where m = min length.
// Space: O(1) extra (sort may be in-place depending on implementation).
func sortEndpoints(strs []string) string {
	if len(strs) == 0 {
		return ""
	}
	sort.Strings(strs)
	first, last := strs[0], strs[len(strs)-1]

	i := 0
	for i < len(first) && i < len(last) && first[i] == last[i] {
		i++
	}
	return first[:i]
}
```

### Dry Run — `strs = ["flower","flow","flight"]`
```
sort.Strings → ["flight","flow","flower"]
first = "flight", last = "flower"
```

| i | first[i] | last[i] | equal? |
|---|----------|---------|--------|
| 0 | 'f'      | 'f'     | yes → i=1 |
| 1 | 'l'      | 'l'     | yes → i=2 |
| 2 | 'i'      | 'o'     | no → stop |

Return `first[:2] = "fl"` ✓

---

## Approach 4 — Binary Search on Prefix Length

### Intuition
The LCP length lies in `[0, minLen]`. This range has a monotonic property: if length `L` works (all strings share a common prefix of length `L`), then all `L' < L` also work. Binary search exploits this to find the maximum valid `L` in O(log m) iterations, each costing O(n·L) to verify.

### Algorithm
1. Find `minLen = min(len(s) for s in strs)`.
2. Binary search `lo=0, hi=minLen`:
   - `mid = (lo+hi+1)/2`.
   - If `isCommonPrefix(strs, mid)`: `lo = mid`.
   - Else: `hi = mid - 1`.
3. Return `strs[0][:lo]`.

### Complexity
- **Time:** O(S log m) — O(log m) iterations × O(S/n per iter) ≈ O(S log m).
- **Space:** O(1).

### Code
```go
// binarySearchLen binary-searches the length of the LCP in range [0, minLen].
// For a given length L, check if strs[0][:L] is a prefix of all strings.
//
// Time:  O(S log m) — O(log m) iterations × O(S/n) per check, where m = min length.
// Space: O(1).
func binarySearchLen(strs []string) string {
	if len(strs) == 0 {
		return ""
	}

	// Find the minimum string length to bound the search.
	minLen := len(strs[0])
	for _, s := range strs[1:] {
		if len(s) < minLen {
			minLen = len(s)
		}
	}

	lo, hi := 0, minLen
	for lo < hi {
		mid := (lo + hi + 1) / 2 // bias up so we converge on the maximum valid length
		if isCommonPrefix(strs, mid) {
			lo = mid
		} else {
			hi = mid - 1
		}
	}
	return strs[0][:lo]
}

// isCommonPrefix checks whether strs[0][:length] is a prefix of every string.
func isCommonPrefix(strs []string, length int) bool {
	prefix := strs[0][:length]
	for _, s := range strs[1:] {
		if !strings.HasPrefix(s, prefix) {
			return false
		}
	}
	return true
}
```

### Dry Run — `strs = ["flower","flow","flight"]`
```
minLen = min(6,4,6) = 4  → search length range [0,4]
```

| lo | hi | mid=(lo+hi+1)/2 | isCommonPrefix(mid) → prefix strs[0][:mid] | move |
|----|----|-----------------|--------------------------------------------|------|
| 0  | 4  | 2 | "fl" prefix of "flow"✓, "flight"✓ → true | lo=2 |
| 2  | 4  | 3 | "flo" prefix of "flow"✓, "flight"✗ → false | hi=2 |
| 2  | 2  | — | lo==hi, loop ends | — |

Return `strs[0][:2] = "fl"` ✓

---

## Key Takeaways

- **Vertical scan short-circuits the earliest** — it stops at the first mismatching column across all strings. For inputs with short common prefixes, it's faster in practice than horizontal scan.
- **Sort trick reduces the problem to one comparison** — LCP(all strings) = LCP(min, max) after sorting. A clean interview trick when you want to impress.
- **Binary search on the answer length** — when a "does length L work?" check is cheap, binary searching L is a valid strategy. The key is spotting the monotonic property: "if L works, all shorter lengths work too."
- **Empty strings in the input** — if any string has length 0, the answer is immediately `""`. Both vertical scan and sort handle this naturally.

---

## Implementation (Go)

See [main.go](main.go).

### Verification
```
["flower","flow","flight"]        → "fl"    ✓
["dog","racecar","car"]           → ""      ✓
["interview","interact","interior"]→ "inter" ✓
["a"]                             → "a"     ✓
```

---

## Related Problems

- LeetCode #208 — Implement Trie (Prefix Tree) — the data structure that directly models common prefixes
- LeetCode #720 — Longest Word in Dictionary (longest word where all prefixes exist)
- LeetCode #1268 — Search Suggestions System (prefix matching with sorted candidates)
