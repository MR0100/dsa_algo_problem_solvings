# Combinatorics (Counting)

## What it is

Combinatorics is the mathematics of **counting** the number of ways something can
happen, without enumerating them one by one. It underpins many "how many …?" and
"minimum resources to distinguish N outcomes" problems. The core building blocks:

| Principle | Statement | Example |
|-----------|-----------|---------|
| **Rule of product** | If step A has `a` choices and step B has `b`, together they have `a·b` | 3 shirts × 4 pants = 12 outfits |
| **Rule of sum** | Disjoint alternatives add: `a + b` | pick a vowel OR a digit |
| **Permutations** | Ordered arrangements: `P(n,k) = n!/(n−k)!` | seat 3 of 5 people in a row |
| **Combinations** | Unordered choices: `C(n,k) = n!/(k!(n−k)!)` | pick 3 of 5 toppings |
| **Multiset / stars-and-bars** | put `n` identical items in `k` bins: `C(n+k−1, k−1)` | non-negative solutions of `x1+…+xk=n` |
| **Information bound** | to distinguish `N` outcomes with `p` independent tests of `b` states each, need `b^p ≥ N` | Poor Pigs |

Combinatorics is a *closed-form* alternative to DP or brute-force enumeration: when
a count has structure, an O(1) or O(k) formula often replaces an O(2^n) search.

---

## When to recognise it

| Signal in the problem | Combinatorial tool |
|-----------------------|--------------------|
| "how many ways / paths / arrangements" | product/sum rule, C(n,k) |
| grid monotone paths (right/down only) | `C(rows+cols, rows)` |
| "count numbers/strings with property, length ≤ n" | per-position product (digit counting) |
| "minimum tests/pigs/weighings to identify one of N" | information bound `states^tests ≥ N` |
| distributing identical items into groups | stars and bars |
| Catalan structures (BST shapes, valid parens, triangulations) | Catalan number `C(2n,n)/(n+1)` |

If you find yourself about to enumerate exponentially many objects only to *count*
them, stop and look for a formula.

---

## General templates / pseudocode

### Binomial coefficient — Pascal's triangle (overflow-safe, no division)
```go
// nCr for small n via additive Pascal recurrence: C(n,k) = C(n-1,k-1)+C(n-1,k).
func binom(n, k int) int64 {
    if k < 0 || k > n {
        return 0
    }
    dp := make([]int64, k+1)
    dp[0] = 1
    for i := 1; i <= n; i++ {
        // iterate j downward so each dp[j] uses the previous row's values
        for j := min(i, k); j >= 1; j-- {
            dp[j] += dp[j-1]
        }
    }
    return dp[k]
}
```

### Multiplicative nCr (fewer allocations, keeps numbers small)
```go
func choose(n, k int) int64 {
    if k < 0 || k > n {
        return 0
    }
    if k > n-k {
        k = n - k // symmetry: C(n,k) == C(n,n-k)
    }
    res := int64(1)
    for i := 0; i < k; i++ {
        res = res * int64(n-i) / int64(i+1) // exact: product of i+1 consecutive ints is divisible by (i+1)!
    }
    return res
}
```

### Information-theoretic minimum tests: smallest `p` with `states^p ≥ n`
```go
// minTests returns the fewest p-state trials needed to distinguish n outcomes.
func minTests(states, n int) int {
    p := 0
    for pow := 1; pow < n; pow *= states { // pow = states^p
        p++
    }
    return p
}
```

---

## Worked example — #458 Poor Pigs

`buckets = 1000`, one test takes `minutesToDie = 15`, total `minutesToTest = 60`.

1. **Rounds available:** each pig can be reused across rounds; number of feeding
   rounds `t = minutesToTest / minutesToDie = 60 / 15 = 4`.
2. **States per pig:** a pig can die after round 1, 2, 3, 4, or survive all → `t + 1 = 5`
   distinguishable outcomes. Each pig is an independent base-5 "digit".
3. **Counting power:** `p` pigs encode `5^p` distinct outcomes. We need `5^p ≥ 1000`.

| `p` | `5^p` | ≥ 1000? |
|-----|-------|---------|
| 3 | 125 | no |
| 4 | 625 | no |
| 5 | 3125 | **yes** |

Answer: **5 pigs**. Closed form: `ceil(log_5(1000)) = 5` — no search over pig
assignments needed.

---

## Complexity

- Pascal `binom`: O(n·k) time, O(k) space.
- Multiplicative `choose`: O(k) time, O(1) space.
- Information bound `minTests`: O(log_states n) time, O(1) space.
- Direct formulas turn many exponential enumerations into O(1)/O(k) counts.

---

## Common pitfalls

- **Overflow.** `n!` explodes fast (21! > 2^63). Prefer Pascal or the multiplicative
  form (dividing as you go), or compute `nCr mod p` with modular inverses.
- **Order matters ≠ order doesn't.** Permutations vs combinations — dividing by `k!`
  (or not) is the single most common counting bug.
- **Off-by-one in the information bound.** You need `states^p ≥ n`, and *states* is
  `rounds + 1` (survival is a valid outcome), not `rounds`.
- **Double counting / missing disjointness.** The sum rule only applies to *disjoint*
  cases; overlapping cases need inclusion–exclusion.
- **Confusing "count" with "enumerate".** If the problem only wants the number, don't
  build the objects — that reintroduces the exponential blow-up you avoided.

---

## Problems in this repo that use it

- [0062 Unique Paths](/0062_unique_paths/README.md) — monotone grid paths = `C(m+n−2, m−1)`
- [0357 Count Numbers with Unique Digits](/0357_count_numbers_with_unique_digits/README.md) — per-position product `9·9·8·…`
- [0458 Poor Pigs](/0458_poor_pigs/README.md) — information bound `(t+1)^pigs ≥ buckets`
- [0096 Unique Binary Search Trees](/0096_unique_binary_search_trees/README.md) — Catalan numbers
- [0377 Combination Sum IV](/0377_combination_sum_iv/README.md) — counting ordered compositions (DP, contrast with closed form)

**Related:** see [`/dsa/math_number_theory.md`](/dsa/math_number_theory.md) for modular
arithmetic (needed for `nCr mod p`), [`/dsa/digit_dp.md`](/dsa/digit_dp.md) for
counting numbers by digit structure, and [`/dsa/dynamic_programming_1d.md`](/dsa/dynamic_programming_1d.md)
for when a count has overlapping subproblems rather than a closed form.
