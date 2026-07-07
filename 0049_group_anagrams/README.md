# 0049 — Group Anagrams

> LeetCode #49 · Difficulty: Medium
> **Categories:** Array, Hash Table, String, Sorting

---

## Problem Statement

Given an array of strings `strs`, group **the anagrams** together. You can return the answer in **any order**.

**Example 1**
```
Input:  strs = ["eat","tea","tan","ate","nat","bat"]
Output: [["bat"],["nat","tan"],["ate","eat","tea"]]
```

**Example 2**
```
Input:  strs = [""]
Output: [[""]]
```

**Example 3**
```
Input:  strs = ["a"]
Output: [["a"]]
```

**Constraints**
- `1 <= strs.length <= 10⁴`
- `0 <= strs[i].length <= 100`
- `strs[i]` consists of lowercase English letters.

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Amazon    | ★★★★★ Very High | 2024          |
| Google    | ★★★★★ Very High | 2024          |
| Meta      | ★★★★★ Very High | 2024          |
| Microsoft | ★★★★☆ High      | 2024          |
| Bloomberg | ★★★★☆ High      | 2023          |
| Apple     | ★★★☆☆ Medium    | 2023          |
| Adobe     | ★★★☆☆ Medium    | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Hash Map with Canonical Key** — the core idea: represent each anagram class with a unique canonical form (sorted string or frequency array) and bucket strings by that key.
- **Sorting as Canonical Form** — sorted string is unique for each anagram class.
- **Frequency Array as Key** — `[26]int` array avoids O(k log k) sort per word.

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Sort Each String as Key | O(n × k log k) | O(n × k) | Simpler code; dominant in practice |
| 2 | Frequency Array as Key ✅ | O(n × k) | O(n × k) | Optimal; avoids sort; better for long strings |

n = number of strings, k = max string length.

---

## Approach 1 — Sort Each String as Key

### Intuition
Two strings are anagrams iff they have the same characters in any order. Sorting each string gives a canonical form: `"eat"` → `"aet"`, `"tea"` → `"aet"`. Group by this key.

### Algorithm
```
groups = {}
for s in strs:
  key = sort(s)
  groups[key].append(s)
return groups.values()
```

### Complexity
- **Time:** O(n × k log k) — n strings, each sorted in O(k log k).
- **Space:** O(n × k) — map storage.

### Code
```go
func groupBySorted(strs []string) [][]string {
    groups := make(map[string][]string)
    for _, s := range strs {
        key := sortedKey(s)
        groups[key] = append(groups[key], s)
    }
    result := make([][]string, 0, len(groups))
    for _, g := range groups {
        result = append(result, g)
    }
    return result
}

func sortedKey(s string) string {
    b := []byte(s)
    sort.Slice(b, func(i, j int) bool { return b[i] < b[j] })
    return string(b)
}
```

### Dry Run — `strs = ["eat","tea","tan","ate","nat","bat"]`
| String | `sortedKey` | `groups` after |
|--------|-------------|----------------|
| "eat" | "aet" | {aet:["eat"]} |
| "tea" | "aet" | {aet:["eat","tea"]} |
| "tan" | "ant" | {aet:[…], ant:["tan"]} |
| "ate" | "aet" | {aet:["eat","tea","ate"], ant:["tan"]} |
| "nat" | "ant" | {aet:[…], ant:["tan","nat"]} |
| "bat" | "abt" | {aet:[…], ant:["tan","nat"], abt:["bat"]} |

Collect map values → [["eat","tea","ate"],["tan","nat"],["bat"]] — 3 groups ✓

---

## Approach 2 — Frequency Array as Key (Recommended ✅)

### Intuition
Character frequency uniquely identifies an anagram class. Represent each string as a `[26]int` array (count of each letter a–z). In Go, `[26]int` arrays are comparable and can serve directly as map keys.

### Algorithm
```
groups = {}
for s in strs:
  freq = [26]int{}
  for ch in s: freq[ch-'a']++
  groups[freq].append(s)
return groups.values()
```

### Complexity
- **Time:** O(n × k) — O(k) to build frequency array; O(26) lookup.
- **Space:** O(n × k).

### Code
```go
func groupByFreq(strs []string) [][]string {
    groups := make(map[[26]int][]string)
    for _, s := range strs {
        var freq [26]int
        for _, ch := range s { freq[ch-'a']++ }
        groups[freq] = append(groups[freq], s)
    }
    result := make([][]string, 0, len(groups))
    for _, g := range groups { result = append(result, g) }
    return result
}
```

### Dry Run — `strs = ["eat","tea","tan","ate","nat","bat"]`
```
"eat": freq[4]++,freq[0]++,freq[19]++ → key=[1,0,0,0,1,0,...,1,...] (a=1,e=1,t=1)
"tea": same key → same group as "eat"
"tan": freq[19]++,freq[0]++,freq[13]++ → key=[1,0,0,0,0,...,1,...,1,...] (a=1,n=1,t=1)
"ate": same key as "eat" → same group
"nat": same key as "tan" → same group
"bat": key=[1,1,0,...,1,...] (a=1,b=1,t=1) → new group

Result: 3 groups: ["eat","tea","ate"], ["tan","nat"], ["bat"] ✓
```

---

## Key Takeaways

- **`[26]int` as map key** — in Go, arrays (unlike slices) are comparable and hashable, making them ideal as map keys. This is a useful Go-specific technique.
- **Sorted string key is simpler to implement** — O(k log k) per string is usually fast enough (k ≤ 100). Use Approach 2 only if k is large or the interviewer asks for optimal.
- **This is the canonical "canonical form + hash map" pattern** — the same technique applies to valid anagram checking (#242), find all anagrams (#438), and similar string grouping problems.

---

## Related Problems

- LeetCode #242 — Valid Anagram (check if two strings are anagrams)
- LeetCode #438 — Find All Anagrams in a String (sliding window + frequency)
- LeetCode #266 — Palindrome Permutation (check if any permutation is a palindrome)
