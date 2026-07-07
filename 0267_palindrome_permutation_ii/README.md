# 0267 — Palindrome Permutation II

> LeetCode #267 · Difficulty: Medium
> **Categories:** Hash Table, String, Backtracking

---

## Problem Statement

Given a string `s`, return all the palindromic permutations (without duplicates)
of it. You may return the answer in **any order**. If `s` has no palindromic
permutation, return an empty list.

**Example 1:**
```
Input: s = "aabb"
Output: ["abba","baab"]
```

**Example 2:**
```
Input: s = "abc"
Output: []
```

**Constraints:**
- `1 <= s.length <= 16`
- `s` consists of only lowercase English letters.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★☆☆ Medium     | 2023          |
| Google     | ★★★☆☆ Medium     | 2023          |
| Uber       | ★★☆☆☆ Low        | 2022          |
| Bloomberg  | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used
- **Backtracking** — enumerate distinct permutations of one half of the string → see [`/dsa/backtracking.md`](/dsa/backtracking.md)
- **Hash Map / Frequency Count** — parity check for feasibility, and half-count construction → see [`/dsa/hash_map.md`](/dsa/hash_map.md)
- **String Algorithms** — palindromes are determined by their left half + a middle char → see [`/dsa/string_algorithms.md`](/dsa/string_algorithms.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Backtracking over Half Permutations (Optimal) | O((n/2)!·n) | O(n) | Standard approach; the only practical one |

---

## Approach 1 — Backtracking over Half Permutations (Optimal)

### Intuition
A palindrome is completely determined by its **left half** plus at most one
**middle character**. So we never permute the whole string. First check
feasibility (at most one character may have an odd count, exactly as in #266).
Then take `count/2` copies of each character to form the multiset of the left
half, permute that multiset without producing duplicates, and mirror each
permutation: `half + middle + reverse(half)`.

### Algorithm
1. Count characters. If more than one has an odd count, return `[]`.
2. Build `half` = each character repeated `count/2` times; if a character has an
   odd count, remember it as the single `mid` (middle) character.
3. Backtrack to permute `half` uniquely: sort the multiset, and when choosing a
   value equal to its previous sibling, skip it if the previous sibling was not
   used at this level (standard "permutations II" duplicate pruning).
4. For each complete half, emit `half + mid + reversed(half)`.

### Complexity
- **Time:** O((n/2)!·n) — permutations of the `n/2`-length half, each costing
  O(n) to build the full palindrome. Duplicate pruning makes the practical count
  far smaller.
- **Space:** O(n) — the counts, the `half` buffer, and recursion depth `n/2`.

### Code
```go
func backtrack(s string) []string {
	counts := make(map[byte]int)
	for i := 0; i < len(s); i++ {
		counts[s[i]]++ // frequency of each character
	}

	// Feasibility + build the half multiset and the middle character.
	var mid string
	oddCount := 0
	half := make([]byte, 0, len(s)/2)
	// Iterate characters in sorted order for deterministic, sorted output.
	keys := make([]byte, 0, len(counts))
	for c := range counts {
		keys = append(keys, c)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	for _, c := range keys {
		if counts[c]%2 == 1 { // this char has an odd count
			oddCount++
			mid = string(c) // candidate middle character
		}
		for k := 0; k < counts[c]/2; k++ {
			half = append(half, c) // each char contributes count/2 to the half
		}
	}
	if oddCount > 1 {
		return []string{} // more than one odd char -> impossible
	}

	var result []string
	used := make([]bool, len(half))
	cur := make([]byte, 0, len(half))

	var dfs func()
	dfs = func() {
		if len(cur) == len(half) { // a full left half is assembled
			// Build palindrome: half + mid + reverse(half).
			left := string(cur)
			rev := make([]byte, len(cur))
			for i := range cur {
				rev[len(cur)-1-i] = cur[i] // reversed copy
			}
			result = append(result, left+mid+string(rev))
			return
		}
		for i := 0; i < len(half); i++ {
			if used[i] {
				continue // already placed at an earlier position
			}
			// Skip duplicates: if this char equals the previous one and the
			// previous is unused at this level, choosing it now repeats work.
			if i > 0 && half[i] == half[i-1] && !used[i-1] {
				continue
			}
			used[i] = true
			cur = append(cur, half[i])
			dfs()
			cur = cur[:len(cur)-1] // undo choice
			used[i] = false
		}
	}
	dfs()
	return result
}
```

### Dry Run
Input `s = "aabb"`:

| Step | Action                                             | State |
|------|----------------------------------------------------|-------|
| 1    | Count chars                                        | {a:2, b:2} |
| 2    | Odd counts? none → `oddCount=0`, `mid=""`          | feasible |
| 3    | Build half: a→1 copy, b→1 copy                     | `half=[a,b]` |
| 4    | DFS pick index 0 = 'a'                             | `cur=[a]` |
| 5    | DFS pick index 1 = 'b'; half complete              | `cur=[a,b]` |
| 6    | Emit `"ab" + "" + "ba"`                            | result=["abba"] |
| 7    | Backtrack; DFS pick index 1 = 'b'                  | `cur=[b]` |
| 8    | DFS pick index 0 = 'a'; half complete              | `cur=[b,a]` |
| 9    | Emit `"ba" + "" + "ab"`                            | result=["abba","baab"] |

Return **["abba", "baab"]**. ✅

---

## Key Takeaways
- Reduce the problem to permuting **half** the string — the palindrome mirrors automatically.
- Reuse #266's parity test as the feasibility gate before doing any work.
- Duplicate-free permutation is the classic **"Permutations II"** pattern: sort, then skip `nums[i]==nums[i-1] && !used[i-1]`.
- The odd-count character (if any) is fixed as the middle — never part of the permuted half.

---

## Related Problems
- LeetCode #266 — Palindrome Permutation (the feasibility check alone)
- LeetCode #47 — Permutations II (duplicate-free permutation pruning)
- LeetCode #46 — Permutations (base permutation backtracking)
- LeetCode #31 — Next Permutation (permutation ordering)
