# 0017 — Letter Combinations of a Phone Number

> LeetCode #17 · Difficulty: Medium
> **Categories:** Hash Table, String, Backtracking

---

## Problem Statement

Given a string containing digits from `2-9` inclusive, return all possible letter combinations that the number could represent. Return the answer in **any order**.

A mapping of digits to letters (just like on the telephone buttons) is given below. Note that `1` does not map to any letters.

```
2 → abc    3 → def    4 → ghi    5 → jkl
6 → mno    7 → pqrs   8 → tuv    9 → wxyz
```

**Example 1**
```
Input:  digits = "23"
Output: ["ad","ae","af","bd","be","bf","cd","ce","cf"]
```

**Example 2**
```
Input:  digits = ""
Output: []
```

**Example 3**
```
Input:  digits = "2"
Output: ["a","b","c"]
```

**Constraints**
- `0 <= digits.length <= 4`
- `digits[i]` is a digit in the range `['2', '9']`.

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Amazon    | ★★★★★ Very High | 2024          |
| Google    | ★★★★★ Very High | 2024          |
| Meta      | ★★★★☆ High      | 2024          |
| Microsoft | ★★★★☆ High      | 2024          |
| Bloomberg | ★★★☆☆ Medium    | 2023          |
| Apple     | ★★★☆☆ Medium    | 2023          |
| Uber      | ★★☆☆☆ Low       | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Backtracking** — Approach 1 builds combinations by choosing one letter per digit and undoing that choice (backtrack) to try the next letter. This is the standard backtracking template.
- **BFS / Queue expansion** — Approach 2 models each partial combination as a BFS state and expands it by all letters of the next digit.

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Backtracking (recursive) | O(4^n · n) | O(n) | The canonical pattern; interview standard |
| 2 | Iterative BFS ✅ | O(4^n · n) | O(4^n · n) | No recursion; easy to reason about iteratively |

Both have the same time complexity. Backtracking uses O(n) stack space; BFS uses O(4^n · n) queue space. For n ≤ 4, both are tiny.

---

## Approach 1 — Backtracking (Recursive)

### Intuition
Think of building the output string position by position. At position `idx`, we have `len(phoneMap[digits[idx]])` choices (2, 3, or 4 letters). Try each, recurse to fill the next position, then undo the choice (backtrack). When all positions are filled (`idx == len(digits)`), record the current path.

### Algorithm
```
btHelper(digits, idx, path, result):
  if idx == len(digits):
    result.append(string(path))
    return
  for each letter in phoneMap[digits[idx]]:
    path.append(letter)
    btHelper(digits, idx+1, path, result)
    path.pop()   // backtrack
```

### Complexity
- **Time:** O(4^n · n) — 4^n combinations, each O(n) to copy into result.
- **Space:** O(n) — recursion stack depth + path buffer of length n.

### Code
```go
func backtracking(digits string) []string {
    if len(digits) == 0 { return []string{} }
    var result []string
    btHelper(digits, 0, []byte{}, &result)
    return result
}
func btHelper(digits string, idx int, path []byte, result *[]string) {
    if idx == len(digits) {
        *result = append(*result, string(path))
        return
    }
    for _, ch := range phoneMap[digits[idx]] {
        path = append(path, byte(ch))
        btHelper(digits, idx+1, path, result)
        path = path[:len(path)-1]
    }
}
```

### Dry Run — `digits = "23"`
```
phoneMap['2']="abc", phoneMap['3']="def"

Call tree (pruned for brevity):
btHelper(idx=0, path=""):
  ch='a' → btHelper(idx=1, path="a"):
    ch='d' → btHelper(idx=2, path="ad") → record "ad"
    ch='e' → record "ae"
    ch='f' → record "af"
  ch='b' → "bd","be","bf"
  ch='c' → "cd","ce","cf"

Result: [ad ae af bd be bf cd ce cf] ✓
```

---

## Approach 2 — Iterative BFS (Recommended ✅)

### Intuition
Maintain a queue of all partial combinations built so far. For each digit, expand every existing partial by all letters that digit maps to. After processing all digits, the queue contains all complete combinations.

### Algorithm
1. Start: `queue = [""]`.
2. For each digit in `digits`:
   - For each partial in queue, for each letter in `phoneMap[digit]`:
     - Append `partial + letter` to the next queue.
   - Replace queue with next queue.
3. Return queue.

### Complexity
- **Time:** O(4^n · n) — same as backtracking.
- **Space:** O(4^n · n) — queue holds all partials at the widest level.

### Code
```go
func iterativeBFS(digits string) []string {
    if len(digits) == 0 { return []string{} }
    queue := []string{""}
    for i := 0; i < len(digits); i++ {
        letters := phoneMap[digits[i]]
        next := make([]string, 0, len(queue)*len(letters))
        for _, partial := range queue {
            for _, ch := range letters {
                next = append(next, partial+string(ch))
            }
        }
        queue = next
    }
    return queue
}
```

### Dry Run — `digits = "23"`
```
Initial: queue = [""]

Digit '2' (letters="abc"):
  "" + 'a' = "a", "" + 'b' = "b", "" + 'c' = "c"
  queue = ["a","b","c"]

Digit '3' (letters="def"):
  "a"+'d'="ad", "a"+'e'="ae", "a"+'f'="af"
  "b"+'d'="bd", "b"+'e'="be", "b"+'f'="bf"
  "c"+'d'="cd", "c"+'e'="ce", "c"+'f'="cf"
  queue = [ad,ae,af,bd,be,bf,cd,ce,cf]

Return [ad ae af bd be bf cd ce cf] ✓
```

---

## Key Takeaways

- **Backtracking template** — `choose → recurse → unchoose`. This is the universal template for combinatorial generation. Memorise it.
- **BFS for combinations** — the BFS queue expansion is an elegant iterative alternative. Each "level" of BFS corresponds to one digit position.
- **Both produce the same total count** — for `n` digits with at most 4 letters each: `3^k * 4^(n-k)` combinations where k = number of digits with 3-letter mappings (2,3,4,5,6,8) and n-k = digits with 4-letter mappings (7,9).
- **`digits.length ≤ 4`** — the constraint makes even brute-force trivially fast here. Focus on demonstrating the backtracking pattern cleanly.

---

## Related Problems

- LeetCode #22 — Generate Parentheses (backtracking with validity constraint)
- LeetCode #46 — Permutations (backtracking over elements)
- LeetCode #77 — Combinations (backtracking over subsets of fixed size)
- LeetCode #784 — Letter Case Permutation (toggle case at each position)
