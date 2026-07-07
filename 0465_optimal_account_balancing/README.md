# 0465 — Optimal Account Balancing

> LeetCode #465 · Difficulty: Hard · 🔒 Premium
> **Categories:** Array, Backtracking, Dynamic Programming, Bit Manipulation, Bitmask

---

## Problem Statement

You are given an array of transactions `transactions` where `transactions[i] = [fromᵢ, toᵢ, amountᵢ]` indicates that the person with `ID = fromᵢ` gave `amountᵢ $` to the person with `ID = toᵢ`.

Return *the minimum number of transactions required to settle the debt*.

**Example 1:**

```
Input: transactions = [[0,1,10],[2,0,5]]
Output: 2
Explanation:
Person #0 gave person #1 $10.
Person #2 gave person #0 $5.
Two transactions are needed. One way to settle the debt is person #1 pays person #0 and #2 $5 each.
```

**Example 2:**

```
Input: transactions = [[0,1,10],[1,0,1],[1,2,5],[2,0,5]]
Output: 1
Explanation:
Person #0 gave person #1 $10.
Person #1 gave person #0 $1.
Person #1 gave person #2 $5.
Person #2 gave person #0 $5.
Therefore, person #1 only need to give person #0 $4, and all debt is settled.
```

**Constraints:**

- `1 <= transactions.length <= 8`
- `transactions[i].length == 3`
- `0 <= fromᵢ, toᵢ <= 20`
- `fromᵢ != toᵢ`
- `1 <= amountᵢ <= 100`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★★☆ High       | 2024          |
| TikTok     | ★★★☆☆ Medium     | 2023          |
| Amazon     | ★★★☆☆ Medium     | 2023          |
| Uber       | ★★☆☆☆ Low        | 2022          |
| Meta       | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Backtracking / DFS** — after netting balances, recursively settle the first non-zero debt against each later opposite-sign debt, branching over choices and taking the minimum transaction count → see [`/dsa/backtracking.md`](/dsa/backtracking.md)
- **Bitmask Subset-Sum DP** — reframed as "partition the non-zero balances into the maximum number of disjoint zero-sum groups"; enumerate subsets with a bitmask, and the answer is `n − (max groups)` → see [`/dsa/bit_manipulation.md`](/dsa/bit_manipulation.md)
- **Hashing to net balances** — collapse the raw transactions into per-person net amounts (a map), discarding who-paid-whom, before any search → see [`/dsa/hash_map.md`](/dsa/hash_map.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Backtracking / DFS Settlement | O(n!) pruned | O(n) | The intuitive, interview-friendly answer; fast for n ≤ 12 debts |
| 2 | Bitmask Subset-Sum DP (max zero-sum groups) | O(3ⁿ) | O(2ⁿ) | Elegant, deterministic bound; showcases the partition identity |

*(n = number of people with a non-zero net balance, ≤ 12 given ≤ 8 transactions.)*

---

## Approach 1 — Backtracking / DFS Settlement

### Intuition

The individual transactions do not matter — only each person's **net** balance does. Net everything into a list of non-zero debts (positives = owed money, negatives = owe money) that sums to `0`. To zero this list with the fewest transfers: take the first non-zero debt and settle it in **one transaction** by moving its entire amount onto some later debt of the **opposite sign** (a debtor pays a creditor). That may zero one balance or offset partially. Recurse on the rest and keep the minimum over all choices. With ≤ 12 non-zero people, this pruned DFS is fast and exact.

### Algorithm

1. Net all transactions into `debts` (drop zero balances).
2. `dfs(start)`:
   - Skip indices whose debt is already `0`.
   - If `start` reaches the end → `0` transactions needed.
   - Otherwise, for each `j > start` with `debts[j]` of **opposite sign** to `debts[start]`: add `debts[start]` to `debts[j]` (settle in one txn), recurse `1 + dfs(start+1)`, track the min, then undo.
3. Return `dfs(0)`.

### Complexity

- **Time:** O(n!) worst case, but the opposite-sign filter and skip-zero pruning make it far smaller in practice (n ≤ 12).
- **Space:** O(n) recursion depth; `debts` is mutated in place and restored.

### Code

```go
func backtracking(transactions [][]int) int {
	debts := buildBalances(transactions)

	var dfs func(start int) int
	dfs = func(start int) int {
		// Advance past debts already settled to 0 (nothing to do for them).
		for start < len(debts) && debts[start] == 0 {
			start++
		}
		if start == len(debts) {
			return 0 // every balance is zero → no more transactions
		}

		best := 1 << 30 // large sentinel for the minimum
		for j := start + 1; j < len(debts); j++ {
			// Only settle against an opposite-sign debt (a debtor pays a creditor).
			if debts[j]*debts[start] < 0 {
				debts[j] += debts[start]        // one transaction moves start's whole balance onto j
				if c := 1 + dfs(start+1); c < best {
					best = c // keep the cheapest continuation
				}
				debts[j] -= debts[start] // undo for the next branch
			}
		}
		return best
	}
	return dfs(0)
}
```

Supporting netting step:

```go
func buildBalances(transactions [][]int) []int {
	balance := make(map[int]int) // person id → net amount (+ owed to them, − they owe)
	for _, t := range transactions {
		from, to, amt := t[0], t[1], t[2]
		balance[from] -= amt // giver's net goes down
		balance[to] += amt   // receiver's net goes up
	}
	debts := []int{}
	for _, v := range balance {
		if v != 0 {
			debts = append(debts, v) // keep only people who still owe/are owed
		}
	}
	return debts
}
```

### Dry Run

Example 1: `transactions = [[0,1,10],[2,0,5]]`.

Netting: `p0 = −10 + 5 = −5`, `p1 = +10`, `p2 = −5`. Non-zero `debts = [−5, +10, −5]` (order may vary; DP/DFS are order-agnostic).

`dfs(0)` with `debts[0] = −5`:

| Step | start | action | debts after | recurse | cost |
|------|-------|--------|-------------|---------|------|
| 1 | 0 (−5) | settle against j=1 (+10, opposite sign): debts[1] += −5 → +5 | `[−5, +5, −5]` | dfs(1) | 1 + dfs(1) |
| 2 | 1 (+5) | settle against j=2 (−5, opposite sign): debts[2] += +5 → 0 | `[−5, +5, 0]` | dfs(2) | 1 + dfs(2) |
| 3 | 2→3 | index 2 is 0, skip → start reaches end | — | — | 0 |

Back-substitute: `dfs(2) = 0` → `dfs(1) = 1` → `dfs(0) = 1 + 1 = 2`. (The alternative first settle j=2 gives the same `2`.)

Result: `2` ✔

---

## Approach 2 — Bitmask Subset-Sum DP (Max Zero-Sum Groups)

### Intuition

Key identity: if the `n` non-zero balances partition into `g` groups that each sum to `0`, then a group of size `k` settles **internally** in `k − 1` transactions, so the total is `Σ(kᵢ − 1) = n − g`. Therefore **minimising transactions ⇔ maximising the number of zero-sum groups**. Enumerate subsets with a bitmask: `sum[mask]` is the total of the selected balances, and `groups[mask]` is the maximum number of zero-sum groups the set `mask` can be split into. A mask can be fully grouped only if its own sum is `0`; then peel off one zero-sum sub-group (anchored on its lowest element to avoid double counting) and add `1`.

### Algorithm

1. Net to non-zero `debts`; `n = len`. If `n == 0`, return `0`.
2. `sum[mask]` = Σ of selected balances (built by adding the lowest set bit's value).
3. For every `mask` with `sum[mask] == 0`, iterate its sub-masks `sub` that contain the anchor bit; if `sum[sub] == 0`, `groups[mask] = max(groups[mask], groups[mask ^ sub] + 1)`.
4. Return `n − groups[full]`.

### Complexity

- **Time:** O(3ⁿ) — summing over all masks of `2^popcount(mask)` sub-masks equals `3ⁿ`.
- **Space:** O(2ⁿ) — the `sum` and `groups` tables.

### Code

```go
func bitmaskDP(transactions [][]int) int {
	debts := buildBalances(transactions)
	n := len(debts)
	if n == 0 {
		return 0 // everyone already balanced
	}

	full := 1 << uint(n)
	sum := make([]int, full)    // sum[mask] = total balance of the selected people
	groups := make([]int, full) // groups[mask] = max #zero-sum groups mask splits into

	// Precompute the sum of every subset by adding the lowest set bit's value.
	for mask := 1; mask < full; mask++ {
		low := mask & (-mask)              // lowest set bit
		idx := bitsTrailingZeros(uint(low)) // which person that bit is
		sum[mask] = sum[mask^low] + debts[idx]
	}

	// DP over masks: only zero-sum masks can be fully partitioned into zero-sum groups.
	for mask := 1; mask < full; mask++ {
		if sum[mask] != 0 {
			continue // a set that does not net to 0 cannot be all zero-sum groups
		}
		// Anchor on the lowest set bit so each group is counted once; that bit's
		// closing group is some zero-sum sub-mask containing it.
		low := mask & (-mask)
		for sub := mask; sub > 0; sub = (sub - 1) & mask {
			if sub&low == 0 {
				continue // the closing group must contain the anchor element
			}
			if sum[sub] == 0 { // sub is itself a valid zero-sum group
				if cand := groups[mask^sub] + 1; cand > groups[mask] {
					groups[mask] = cand // one more zero-sum group formed
				}
			}
		}
	}
	// n people minus the most zero-sum groups = fewest transactions.
	return n - groups[full-1]
}
```

### Dry Run

Example 1: `debts = [−5, +10, −5]`, `n = 3`, `full = 8`. Index the people as bit0 = `−5`, bit1 = `+10`, bit2 = `−5`.

Subset sums (`sum[mask]`):

| mask (bits) | selected | sum |
|-------------|----------|-----|
| 001 | {−5} | −5 |
| 010 | {+10} | +10 |
| 011 | {−5,+10} | +5 |
| 100 | {−5} | −5 |
| 101 | {−5,−5} | −10 |
| 110 | {+10,−5} | +5 |
| 111 | {−5,+10,−5} | 0 |

Only `mask = 111` has `sum == 0`. Anchor bit = `001`. Try sub-masks containing bit0:

| sub (bits) | sum[sub] | zero-sum? | groups[111^sub] + 1 |
|------------|----------|-----------|---------------------|
| 111 | 0 | yes | groups[000] + 1 = 0 + 1 = **1** |
| 101 | −10 | no | — |
| 011 | +5 | no | — |
| 001 | −5 | no | — |

Best `groups[111] = 1` (the whole set is one zero-sum group; it cannot be split further because no proper subset containing bit0 nets to 0).

Answer: `n − groups[full−1] = 3 − 1 = 2`. Result: `2` ✔

---

## Key Takeaways

- **Net first, then settle.** The raw edges are noise; per-person net balances (summing to 0) are the real state. This collapse is the crucial first move.
- **Min transactions = (non-zero people) − (max zero-sum groups).** A size-`k` self-settling group costs `k − 1`, so more groups ⇒ fewer transfers. Recognising this turns a fuzzy optimisation into a clean subset-partition DP.
- **Small `n` ⇒ exponential is fine.** `≤ 8` transactions bound the non-zero people to `≤ 12`; both `O(n!)`-pruned DFS and `O(3ⁿ)` bitmask DP run instantly. The constraint size is the signal to reach for subset enumeration.
- **`mask & (-mask)`** extracts the lowest set bit — used both to build subset sums incrementally and to anchor groups so each partition is counted exactly once.
- **Sub-mask enumeration `sub = (sub - 1) & mask`** walks every subset of `mask` in `O(2^popcount)`, the backbone of subset DP.

---

## Related Problems

- LeetCode #698 — Partition to K Equal Sum Subsets (bitmask subset partition)
- LeetCode #416 — Partition Equal Subset Sum (subset-sum DP)
- LeetCode #473 — Matchsticks to Square (backtracking into k equal groups)
- LeetCode #1723 — Find Minimum Time to Finish All Jobs (bitmask/backtracking assignment)
- LeetCode #2172 — Maximum AND Sum of Array (bitmask assignment DP)
