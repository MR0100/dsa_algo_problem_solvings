package main

import "fmt"

// buildBalances nets all transactions into a per-person balance, then returns
// only the NON-ZERO balances (people who owe or are owed). People whose net is
// zero are irrelevant and are dropped. The sum of the returned slice is always
// 0 (money is conserved), so positives and negatives cancel exactly.
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

// ── Approach 1: Backtracking / DFS Settlement (Optimal for n ≤ 8) ─────────────
//
// backtracking solves Optimal Account Balancing by netting everyone's balance,
// then recursively settling the first unsettled debt against every later debt
// of the opposite sign, minimising the transaction count.
//
// Intuition:
//
//	Only net balances matter — who paid whom en route is irrelevant. After
//	netting, we have a list of non-zero debts summing to 0. To zero the list
//	with the fewest transfers, take the first non-zero debt and "settle" it by
//	moving its whole amount onto some later debt of the opposite sign (one
//	transaction). That may zero one or both. Recurse from the next index and
//	keep the best. This greedy-branching DFS is exact for the tiny input
//	(≤ 8 transactions ⇒ ≤ 12 non-zero people).
//
// Algorithm:
//  1. Compute the non-zero balances `debts`.
//  2. dfs(start): skip indices whose debt is already 0; when start reaches the
//     end, 0 more transactions are needed. Otherwise, for each j > start with
//     debts[j] of opposite sign to debts[start]: add debts[start] to debts[j]
//     (settle in one txn), recurse dfs(start+1), take 1 + min over choices,
//     then undo.
//  3. Return dfs(0).
//
// Time:  O(n!) worst case over the non-zero debts (n = count of them), pruned
//
//	heavily in practice by the opposite-sign filter and same-value skip.
//
// Space: O(n) recursion depth (debts mutated in place).
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
				debts[j] += debts[start] // one transaction moves start's whole balance onto j
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

// ── Approach 2: Bitmask Subset-Sum DP (Max Zero-Sum Groups) ───────────────────
//
// bitmaskDP solves Optimal Account Balancing via the identity
// answer = (#non-zero people) − (max number of disjoint zero-sum groups they
// split into). A group of size k that sums to 0 settles internally in k−1
// transactions, so maximising the number of groups minimises total transfers.
//
// Intuition:
//
//	If the n non-zero balances partition into g groups that each sum to 0,
//	each group of size k needs k−1 internal transfers, for a grand total of
//	n − g. So minimising transactions ⇔ maximising the count of zero-sum
//	groups. Enumerate subsets with a bitmask: sum[mask] is the total of the
//	chosen balances; groups[mask] is the max zero-sum groups the set `mask`
//	can be cut into. For each mask, if its own sum is 0 it may itself be a
//	closing group, so try removing a zero-sum sub-piece and add 1.
//
// Algorithm:
//  1. Compute non-zero balances; n = len; answer 0 if n == 0.
//  2. sum[mask] = Σ balances selected by mask (built incrementally).
//  3. groups[mask]: for each mask with sum[mask] == 0, try every sub-mask
//     `sub` of mask with sum[sub] == 0: groups[mask] = max(groups[mask],
//     groups[mask ^ sub] + 1). (Fix one element to avoid double counting.)
//  4. Return n − groups[full].
//
// Time:  O(3^n) — iterating all sub-masks of all masks (Σ over masks 2^popcount).
// Space: O(2^n) — the sum and groups tables.
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
		low := mask & (-mask)               // lowest set bit
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

// bitsTrailingZeros returns the index of the lowest set bit of a power-of-two
// value (i.e. how many trailing zeros it has). Local helper to keep imports lean.
func bitsTrailingZeros(v uint) int {
	idx := 0
	for v&1 == 0 { // shift until the set bit reaches position 0
		v >>= 1
		idx++
	}
	return idx
}

func main() {
	ex1 := [][]int{{0, 1, 10}, {2, 0, 5}}
	ex2 := [][]int{{0, 1, 10}, {1, 0, 1}, {1, 2, 5}, {2, 0, 5}}
	ex3 := [][]int{{1, 2, 3}, {3, 4, 3}, {5, 6, 3}} // three independent debts → 3 txns

	fmt.Println("=== Approach 1: Backtracking / DFS Settlement (Optimal for n ≤ 8) ===")
	fmt.Printf("[[0,1,10],[2,0,5]]                  got=%d  expected 2\n", backtracking(ex1))
	fmt.Printf("[[0,1,10],[1,0,1],[1,2,5],[2,0,5]]  got=%d  expected 1\n", backtracking(ex2))
	fmt.Printf("[[1,2,3],[3,4,3],[5,6,3]]           got=%d  expected 3\n", backtracking(ex3))

	fmt.Println("=== Approach 2: Bitmask Subset-Sum DP (Max Zero-Sum Groups) ===")
	fmt.Printf("[[0,1,10],[2,0,5]]                  got=%d  expected 2\n", bitmaskDP(ex1))
	fmt.Printf("[[0,1,10],[1,0,1],[1,2,5],[2,0,5]]  got=%d  expected 1\n", bitmaskDP(ex2))
	fmt.Printf("[[1,2,3],[3,4,3],[5,6,3]]           got=%d  expected 3\n", bitmaskDP(ex3))
}
