# 0060 — Permutation Sequence

> LeetCode #60 · Difficulty: Hard
> **Categories:** Math, Recursion

---

## Problem Statement

The set `[1, 2, 3, ..., n]` contains a total of `n!` unique permutations.

By listing and labeling all of the permutations in order, we get the following sequence for `n = 3`:
```
1. "123"
2. "132"
3. "213"
4. "231"
5. "312"
6. "321"
```

Given `n` and `k`, return the `k`th permutation sequence.

**Example 1**
```
Input:  n = 3, k = 3
Output: "213"
```

**Example 2**
```
Input:  n = 4, k = 9
Output: "2314"
```

**Example 3**
```
Input:  n = 3, k = 1
Output: "123"
```

**Constraints**
- `1 <= n <= 9`
- `1 <= k <= n!`

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Amazon    | ★★★★☆ High      | 2024          |
| Google    | ★★★★☆ High      | 2024          |
| Bloomberg | ★★★★☆ High      | 2023          |
| Microsoft | ★★★☆☆ Medium    | 2023          |
| Meta      | ★★★☆☆ Medium    | 2023          |
| Apple     | ★★☆☆☆ Low       | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Factorial Number System** — represent k-1 in the factorial base to directly compute which digit to pick at each position without generating all permutations.
- **Greedy Digit Picking** — at each step, determine the correct digit by integer division with the factorial of remaining choices.

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Generate All (Backtracking) | O(n! × n) | O(n! × n) | Reference; TLE for large n |
| 2 | Factorial Number System ✅ | O(n²) | O(n) | Optimal; direct O(n²) construction |

---

## Approach 1 — Generate All Permutations

### Intuition
Backtrack to generate all n! permutations in lexicographic order, collect them in a slice, return the (k-1)-th element (0-indexed).

### Complexity
- **Time:** O(n! × n) — generates all permutations.
- **Space:** O(n! × n).

---

## Approach 2 — Factorial Number System (Recommended ✅)

### Intuition
There are `n!` permutations of `n` digits. The `n!` permutations are grouped into `n` blocks of `(n-1)!` each, one block per leading digit. The first digit of the k-th permutation is determined by `(k-1) / (n-1)!`.

After selecting the first digit, we have `n-1` remaining digits and need to find the `((k-1) % (n-1)!)` th permutation among them. Repeat.

This is exactly representing `k-1` in the **factorial number system** (base = decreasing factorials):
```
k-1 = d[0]*(n-1)! + d[1]*(n-2)! + ... + d[n-2]*1! + d[n-1]*0!
```
`d[i]` = index into the remaining available digits at step `i`.

### Algorithm
```
digits = [1, 2, ..., n]
fact = [0!, 1!, ..., (n-1)!]  // fact[i] = i!
k -= 1  // convert to 0-indexed
for i = n downto 1:
  idx = k / fact[i-1]
  append digits[idx] to result
  remove digits[idx]
  k %= fact[i-1]
```

### Complexity
- **Time:** O(n²) — n iterations; each digit removal from the digits list is O(n).
- **Space:** O(n) — digits list.

### Code
```go
func factorialSystem(n, k int) string {
    fact := make([]int, n); fact[0] = 1
    for i := 1; i < n; i++ { fact[i] = fact[i-1] * i }

    digits := make([]int, n)
    for i := range digits { digits[i] = i + 1 }

    k--  // 0-indexed
    var sb strings.Builder
    for i := n; i >= 1; i-- {
        idx := k / fact[i-1]
        sb.WriteString(strconv.Itoa(digits[idx]))
        digits = append(digits[:idx], digits[idx+1:]...)
        k %= fact[i-1]
    }
    return sb.String()
}
```

### Dry Run — `n = 4, k = 9`
```
fact = [1, 1, 2, 6]  (0!=1, 1!=1, 2!=2, 3!=6)
digits = [1, 2, 3, 4]
k = 8 (0-indexed)

i=4: idx = 8/6 = 1 → pick digits[1]=2. digits=[1,3,4]. k=8%6=2.
i=3: idx = 2/2 = 1 → pick digits[1]=3. digits=[1,4]. k=2%2=0.
i=2: idx = 0/1 = 0 → pick digits[0]=1. digits=[4]. k=0%1=0.
i=1: idx = 0/1 = 0 → pick digits[0]=4. digits=[].

Result: "2314" ✓
```

---

## Key Takeaways

- **`k -= 1` converts to 0-indexed** — the key to making the `k / fact[i-1]` formula work. Without this, the first permutation would require special-casing.
- **Factorial number system** — a numeral system where position `i` from the right has base `i+1`. `k-1` written in this system directly gives the digit indices.
- **O(n²) is optimal for this problem** — you cannot do better than O(n) to write the answer; the removal step costs O(n) per position, giving O(n²) total. A linked-list or Fenwick tree can reduce removal to O(log n) giving O(n log n), but n ≤ 9 makes this irrelevant.
- **Compare with #31 and #46** — #31 generates the next permutation; #46 generates all; #60 jumps directly to the k-th.

---

## Related Problems

- LeetCode #31 — Next Permutation (generate exactly one successor)
- LeetCode #46 — Permutations (generate all permutations)
- LeetCode #47 — Permutations II (deduplicated permutations)
