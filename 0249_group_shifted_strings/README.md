# 0249 — Group Shifted Strings

> LeetCode #249 · Difficulty: Medium
> **Categories:** Array, Hash Table, String

---

## Problem Statement

We can shift a string by shifting each of its letters to its successive letter.

- For example, `"abc"` can be shifted to be `"bcd"`.

We can keep shifting the string to form a sequence.

- For example, we can keep shifting `"abc"` to form the sequence: `"abc" -> "bcd" -> ... -> "xyz"`.

Given an array of strings `strings`, group all `strings[i]` that belong to the same shifting sequence. You may return the answer in **any order**.

**Example 1:**

```
Input: strings = ["abc","bcd","acef","xyz","az","ba","a","z"]
Output: [["acef"],["a","z"],["abc","bcd","xyz"],["az","ba"]]
```

**Example 2:**

```
Input: strings = ["a"]
Output: [["a"]]
```

**Constraints:**

- `1 <= strings.length <= 200`
- `1 <= strings[i].length <= 50`
- `strings[i]` consists of lowercase English letters.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★★☆ High       | 2024          |
| Amazon     | ★★★☆☆ Medium     | 2023          |
| Uber       | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2022          |
| Bloomberg  | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Hash Map** — bucket every string under a canonical key so same-sequence strings collide into one group → see [`/dsa/hash_map.md`](/dsa/hash_map.md)
- **String Algorithms** — the invariant of a shift is the sequence of consecutive letter gaps (mod 26); computing that signature is a linear string scan → see [`/dsa/string_algorithms.md`](/dsa/string_algorithms.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Hash Map by Shift Signature (Optimal) | O(N·L) | O(N·L) | Canonical: key is the mod-26 gap tuple, robust and clear |
| 2 | Normalize to Canonical Form | O(N·L) | O(N·L) | Equivalent; key is the string shifted so it starts at 'a' |

---

## Approach 1 — Hash Map by Shift Signature (Optimal)

### Intuition

Shifting moves every letter forward by the same amount with `z→a` wraparound, so it preserves the **gap** between each pair of consecutive letters. Two strings are in the same shifting sequence iff their gap tuples (mod 26) are identical. So use that gap tuple as a hash key. Single-character strings have no gaps — they all share the empty key and group together.

### Algorithm

1. For each string `s`, compute a signature: for `i` in `1..len-1`, `diff = (s[i]-s[i-1]+26) % 26`; join the diffs (with a delimiter) into a key.
2. Append `s` to the map bucket for that key.
3. Return all buckets.

### Complexity

- **Time:** O(N·L) — N strings, each length L, one linear pass to build a key.
- **Space:** O(N·L) — the map holds every string plus its key.

### Code

```go
func hashMapBySignature(strings []string) [][]string {
	groups := map[string][]string{}

	for _, s := range strings {
		key := shiftSignature(s)
		groups[key] = append(groups[key], s)
	}

	result := make([][]string, 0, len(groups))
	for _, g := range groups {
		result = append(result, g)
	}
	return result
}

func shiftSignature(s string) string {
	if len(s) <= 1 {
		return ""
	}
	key := ""
	for i := 1; i < len(s); i++ {
		diff := (int(s[i]) - int(s[i-1]) + 26) % 26
		key += fmt.Sprintf("%d,", diff)
	}
	return key
}
```

### Dry Run

Trace `["abc","bcd","acef","xyz","az","ba","a","z"]`:

| String | gaps (mod 26) | key | bucket after |
|--------|---------------|-----|--------------|
| `abc`  | (1,1)         | `1,1,` | `{1,1,: [abc]}` |
| `bcd`  | (1,1)         | `1,1,` | `{1,1,: [abc,bcd]}` |
| `acef` | (2,2,1)       | `2,2,1,` | new bucket `[acef]` |
| `xyz`  | (1,1)         | `1,1,` | `{1,1,: [abc,bcd,xyz]}` |
| `az`   | (25)          | `25,` | new bucket `[az]` |
| `ba`   | (25)          | `25,` | `{25,: [az,ba]}` |
| `a`    | ()            | `` (empty) | new bucket `[a]` |
| `z`    | ()            | `` (empty) | `{: [a,z]}` |

Groups = `[[abc,bcd,xyz],[acef],[az,ba],[a,z]]` → 4 groups. ✓

---

## Approach 2 — Normalize to Canonical Form

### Intuition

If we shift each string back by the exact amount that turns its first letter into `'a'`, every string in the same sequence collapses to the identical normalized string. That normalized form is a natural, human-readable key that yields the same grouping as the gap signature.

### Algorithm

1. For each string `s`, `shift = s[0]-'a'`.
2. Build `normalized`: each char `c → 'a' + (c-'a'-shift+26) % 26`.
3. Bucket `s` under `normalized`; return the buckets.

### Complexity

- **Time:** O(N·L).
- **Space:** O(N·L).

### Code

```go
func normalizeToBase(strings []string) [][]string {
	groups := map[string][]string{}

	for _, s := range strings {
		key := normalize(s)
		groups[key] = append(groups[key], s)
	}

	result := make([][]string, 0, len(groups))
	for _, g := range groups {
		result = append(result, g)
	}
	return result
}

func normalize(s string) string {
	if len(s) == 0 {
		return ""
	}
	shift := int(s[0] - 'a')
	buf := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		buf[i] = byte('a' + (int(s[i]-'a')-shift+26)%26)
	}
	return string(buf)
}
```

### Dry Run

Trace a few from `["abc","bcd","acef","xyz","az","ba"]`:

| String | shift = s[0]-'a' | normalized (start at 'a') | bucket |
|--------|------------------|---------------------------|--------|
| `abc`  | 0 | `abc` | `[abc]` |
| `bcd`  | 1 | `abc` | `[abc,bcd]` |
| `acef` | 0 | `acef` | new `[acef]` |
| `xyz`  | 23 | `abc` | `[abc,bcd,xyz]` |
| `az`   | 0 | `az` | new `[az]` |
| `ba`   | 1 | `az` | `[az,ba]` |

Same grouping as Approach 1. ✓

---

## Key Takeaways

- The invariant of "same shifting sequence" is the tuple of **consecutive letter gaps mod 26** — turn any equivalence into a canonical key, then group with a hash map.
- Always delimit multi-value keys (`"1,2,"` not `"12"`) so distinct gap sequences never collide.
- Normalizing to a base form (shift so it starts at `'a'`) is an equivalent, often more intuitive key.

---

## Related Problems

- LeetCode #49 — Group Anagrams (same "canonical key + hash map" pattern)
- LeetCode #205 — Isomorphic Strings (structural equivalence of strings)
- LeetCode #290 — Word Pattern (mapping-based grouping)
