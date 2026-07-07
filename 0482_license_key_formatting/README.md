# 0482 — License Key Formatting

> LeetCode #482 · Difficulty: Easy
> **Categories:** String

---

## Problem Statement

You are given a license key represented as a string `s` that consists of only alphanumeric characters and dashes. The string is separated into `n + 1` groups by `n` dashes. You are also given an integer `k`.

We want to reformat the string `s` such that each group contains exactly `k` characters, except for the first group, which could be shorter than `k` but still must contain at least one character. Furthermore, there must be a dash inserted between two groups, and you should convert all lowercase letters to uppercase.

Return *the reformatted license key*.

**Example 1:**

```
Input: s = "5F3Z-2e-9-w", k = 4
Output: "5F3Z-2E9W"
Explanation: The string s has been split into two parts, each part has 4 characters.
Note that the two extra dashes are not needed and can be removed.
```

**Example 2:**

```
Input: s = "2-5g-3-J", k = 2
Output: "2-5G-3J"
Explanation: The string s has been split into three parts, each part has 2 characters except the first part as it could be shorter as mentioned above.
```

**Constraints:**

- `1 <= s.length <= 10^5`
- `s` consists of English letters, digits, and dashes `'-'`.
- `1 <= k <= 10^4`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★★☆ High       | 2024          |
| Amazon     | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Capital One| ★★☆☆☆ Low        | 2022          |
| Bloomberg  | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **String building & normalisation** — the task is pure string manipulation: filter out dashes, upper-case letters, and re-chunk. Building the result with a byte buffer (and a single reversal in the optimal version) avoids repeated allocations → see [`/dsa/string_algorithms.md`](/dsa/string_algorithms.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Clean, Then Slice From the Left | O(n) | O(n) | Most readable; explicit first-group arithmetic (`L % k`) |
| 2 | Build From the Right, One Char at a Time (Optimal) | O(n) | O(n) | No length math; grouping falls out of a backward walk + one reversal |

> Both are linear. Approach 2 avoids computing the first-group size explicitly — the short leading group emerges naturally because we count groups from the right.

---

## Approach 1 — Clean, Then Slice From the Left (Brute Force)

### Intuition

Strip the string down to its raw alphanumeric characters (upper-cased). Grouping is defined **from the right**: every group is exactly `k` except the first, which holds the remainder. So if there are `L` clean characters, the first group's size is `firstLen = L % k` — unless that is `0` (an exact multiple), in which case the first group is a full `k`. Once `firstLen` is known, emit the leading chunk and then successive `k`-sized chunks left to right, joining with dashes.

### Algorithm

1. Build `clean` = every non-dash character of `s`, converting lowercase to uppercase.
2. If `clean` is empty (string was only dashes), return `""`.
3. `firstLen = len(clean) % k`; if `firstLen == 0`, set `firstLen = k`.
4. Emit `clean[0:firstLen]`, then `clean[firstLen:firstLen+k]`, `clean[firstLen+k : firstLen+2k]`, … as groups.
5. Join all groups with `"-"`.

### Complexity

- **Time:** O(n) — one pass to clean, one pass to slice into groups (`n = len(s)`).
- **Space:** O(n) — the cleaned string and the output buffer.

### Code

```go
func bruteForceLeftSlice(s string, k int) string {
	// Step 1: strip dashes and upper-case in a single pass.
	var clean strings.Builder
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '-' {
			continue // dashes carry no information; drop them
		}
		if c >= 'a' && c <= 'z' {
			c -= 32 // ASCII lower→upper ('a'-'A' == 32)
		}
		clean.WriteByte(c)
	}
	cs := clean.String()
	if len(cs) == 0 {
		return "" // nothing but dashes → empty result
	}

	// Step 3: size of the (possibly short) first group.
	firstLen := len(cs) % k
	if firstLen == 0 {
		firstLen = k // exact multiple → first group is a full k, not empty
	}

	// Step 4-5: emit the first group, then successive full k-groups.
	groups := []string{cs[:firstLen]}       // leading, possibly-short chunk
	for i := firstLen; i < len(cs); i += k { // every remaining chunk is exactly k
		groups = append(groups, cs[i:i+k])
	}
	return strings.Join(groups, "-")
}
```

### Dry Run

Example 1: `s = "5F3Z-2e-9-w", k = 4`.

| Step | Detail | Value |
|------|--------|-------|
| 1 | clean (dashes removed, upper-cased) | `"5F3Z2E9W"` (length 8) |
| 3 | `firstLen = 8 % 4` | `0` → reset to `k = 4` |
| 4 | first group `cs[0:4]` | `"5F3Z"` |
| 4 | next group `cs[4:8]` | `"2E9W"` |
| 5 | join with `-` | `"5F3Z-2E9W"` |

Result: `"5F3Z-2E9W"` ✔

---

## Approach 2 — Build From the Right, One Char at a Time (Optimal)

### Intuition

Because groups are anchored on the **right**, the cleanest construction walks `s` backwards. Keep a counter of how many real characters have been placed since the last dash; whenever it becomes a positive multiple of `k`, drop a dash *before* the next character. This automatically leaves the leftmost group short (whatever remains), with no length arithmetic or slicing. We build the answer reversed and flip it once at the end.

### Algorithm

1. Walk `i` from `len(s)-1` down to `0`.
2. Skip dashes. For each real character: if `count > 0` and `count % k == 0`, append `'-'`; then append the upper-cased character and increment `count`.
3. Reverse the accumulated bytes to restore left-to-right order and return.

### Complexity

- **Time:** O(n) — one backward pass plus one in-place reversal.
- **Space:** O(n) — the output buffer.

### Code

```go
func rightToLeftBuild(s string, k int) string {
	buf := make([]byte, 0, len(s)) // reversed output accumulator
	count := 0                     // real chars placed since the last dash

	for i := len(s) - 1; i >= 0; i-- {
		c := s[i]
		if c == '-' {
			continue // ignore existing separators
		}
		if count > 0 && count%k == 0 {
			buf = append(buf, '-') // completed a group of k → separator goes here
		}
		if c >= 'a' && c <= 'z' {
			c -= 32 // normalise to uppercase
		}
		buf = append(buf, c)
		count++ // one more real character in the current group
	}

	// buf currently holds the answer reversed; flip it in place.
	for l, r := 0, len(buf)-1; l < r; l, r = l+1, r-1 {
		buf[l], buf[r] = buf[r], buf[l]
	}
	return string(buf)
}
```

### Dry Run

Example 1: `s = "5F3Z-2e-9-w", k = 4`. Walk right→left; `buf` grows reversed.

| i | s[i] | dash? | count%k==0 & count>0? → dash | char appended (upper) | buf (reversed) | count after |
|---|------|-------|------------------------------|-----------------------|----------------|-------------|
| 10 | `w` | no | count=0, no | `W` | `W` | 1 |
| 9 | `-` | yes (skip) | — | — | `W` | 1 |
| 8 | `9` | no | 1%4≠0 | `9` | `W9` | 2 |
| 7 | `-` | yes (skip) | — | — | `W9` | 2 |
| 6 | `e` | no | 2%4≠0 | `E` | `W9E` | 3 |
| 5 | `2` | no | 3%4≠0 | `2` | `W9E2` | 4 |
| 4 | `-` | yes (skip) | — | — | `W9E2` | 4 |
| 3 | `Z` | no | 4%4==0 & count>0 → append `-` | `Z` | `W9E2-Z` | 5 |
| 2 | `3` | no | 5%4≠0 | `3` | `W9E2-Z3` | 6 |
| 1 | `F` | no | 6%4≠0 | `F` | `W9E2-Z3F` | 7 |
| 0 | `5` | no | 7%4≠0 | `5` | `W9E2-Z3F5` | 8 |

Reverse `W9E2-Z3F5` → `"5F3Z-2E9W"` ✔

---

## Key Takeaways

- **Anchor grouping on the correct side.** The "first group may be short" rule means groups are counted from the *right*; either compute `L % k` up front (Approach 1) or walk backwards so the short group falls out for free (Approach 2).
- **`L % k == 0` is the trap.** An exact multiple must yield a full first group, not an empty one — always special-case it.
- **Build reversed, flip once** is a clean idiom when the natural construction order is opposite to the desired output order — cheaper than prepending (which is O(n) per insert).
- Case-fold with `c -= 32` for ASCII letters, or use `strings.ToUpper` for clarity; both are fine at this scale.

---

## Related Problems

- LeetCode #6 — Zigzag Conversion (re-chunking characters by position)
- LeetCode #68 — Text Justification (grouping tokens into fixed-width lines)
- LeetCode #38 — Count and Say (string transformation / rebuild)
- LeetCode #1108 — Defanging an IP Address (simple string normalisation)
