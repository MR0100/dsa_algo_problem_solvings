# 0444 — Sequence Reconstruction

> LeetCode #444 · Difficulty: Medium (Premium)
> **Categories:** Graph, Topological Sort

---

## Problem Statement

You are given an integer array `nums` of length `n` where `nums` is a permutation of the integers in the range `[1, n]`. You are also given a 2D integer array `sequences` where `sequences[i]` is a subsequence of `nums`.

Check if `nums` is the shortest possible and the only **supersequence**. The shortest supersequence is a sequence **with the shortest length** and has all `sequences[i]` as subsequences. There could be multiple valid **supersequences** for the given array `sequences`.

- For example, for `sequences = [[1,2],[1,3]]`, there are two shortest supersequences, `[1,2,3]` and `[1,3,2]`.
- While for `sequences = [[1,2],[1,3],[1,2,3]]`, the only shortest supersequence possible is `[1,2,3]`. `[1,2,3,4]` is a possible supersequence but not the shortest.

Return `true` *if `nums` is the only shortest supersequence for `sequences`, or* `false` *otherwise*.

> Historical phrasing (original premium version): *Check whether the original sequence `org` can be uniquely reconstructed from the sequences in `seqs`. Reconstruction means building a shortest common supersequence of the sequences in `seqs` (a sequence so that all sequences in `seqs` are subsequences of it), and determining whether there is only one sequence that can be reconstructed and it equals `org`.*

**Example 1:**

```
Input: nums = [1,2,3], sequences = [[1,2],[1,3]]
Output: false
Explanation: There are two possible supersequences: [1,2,3] and [1,3,2].
The sequence [1,2] is a subsequence of both: [1,2,3] and [1,3,2].
The sequence [1,3] is a subsequence of both: [1,2,3] and [1,3,2].
Because nums is not the only shortest supersequence, we return false.
```

**Example 2:**

```
Input: nums = [1,2,3], sequences = [[1,2]]
Output: false
Explanation: The shortest possible supersequence is [1,2].
The sequence [1,2] is a subsequence of it: [1,2].
Because nums is not the shortest supersequence, we return false.
```

**Example 3:**

```
Input: nums = [1,2,3], sequences = [[1,2],[1,3],[2,3]]
Output: true
Explanation: The shortest possible supersequence is [1,2,3].
The sequence [1,2] is a subsequence of it: [1,2,3].
The sequence [1,3] is a subsequence of it: [1,2,3].
The sequence [2,3] is a subsequence of it: [1,2,3].
Because nums is the only shortest supersequence, we return true.
```

**(Historical) Example 4:**

```
Input: org = [4,1,5,2,6,3], seqs = [[5,2,6,3],[4,1,5,2]]
Output: true
```

**Constraints:**

- `n == nums.length`
- `1 <= n <= 10^4`
- `nums` is a permutation of all the integers in the range `[1, n]`.
- `1 <= sequences.length <= 10^4`
- `1 <= sequences[i].length <= 10^4`
- `1 <= sum(sequences[i].length) <= 10^5`
- `1 <= sequences[i][j] <= n`
- All the arrays of `sequences` are **unique**.
- `sequences[i]` is a subsequence of `nums`.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★☆☆ Medium     | 2023          |
| Amazon     | ★★★☆☆ Medium     | 2023          |
| Meta       | ★★☆☆☆ Low        | 2022          |
| Microsoft  | ★★☆☆☆ Low        | 2022          |
| Airbnb     | ★★☆☆☆ Low        | 2021          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Topological Sort (Kahn's BFS)** — the pairwise "a before b" constraints form a DAG; a *unique* topological order exists iff every step has exactly one in-degree-0 node → see [`/dsa/topological_sort.md`](/dsa/topological_sort.md)
- **Graph modelling from constraints** — turning consecutive pairs of each sequence into directed edges is the crux; the reconstruction question becomes a graph property → see [`/dsa/graph_bfs_dfs.md`](/dsa/graph_bfs_dfs.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Topological Sort — Uniqueness (Kahn BFS) | O(V + E) | O(V + E) | The canonical answer; directly checks "one order only" |
| 2 | Adjacent-Coverage Check | O(V + E) | O(n) | Slick O(n)-space alternative; no explicit graph build |

(`V = n` distinct values, `E = Σ(len(seq) − 1)` pairs.)

---

## Approach 1 — Topological Sort — Uniqueness (Kahn's BFS)

### Intuition

Every sequence lists its elements in the order they must appear, so each consecutive pair `(a, b)` becomes a directed edge `a → b`. `nums` is the *unique* shortest supersequence iff the resulting DAG has **exactly one** topological ordering **and** that ordering equals `nums`. In Kahn's algorithm, a topological order is unique precisely when, at every step, there is **exactly one** node with in-degree 0 — if two sources ever coexist, they could be output in either order, producing a second valid supersequence. We also require all `n` values to appear (otherwise the constraints don't pin `nums` at all).

### Algorithm

1. Build an adjacency list and in-degree map over the values in `sequences` (dedup parallel edges so in-degrees stay exact); validate each value is in `[1, n]` and count distinct values.
2. If the distinct count `≠ n`, return `false`.
3. Seed a queue with all in-degree-0 nodes.
4. Loop: if the queue holds `> 1` node, the order is ambiguous → `false`. Pop the single node; it must equal `nums[index]` (else contradiction → `false`); advance `index`. Decrement neighbours' in-degrees, enqueuing any that reach 0.
5. Return `index == n` (all values emitted in `nums`'s exact order).

### Complexity

- **Time:** O(V + E) — each node dequeued once, each edge relaxed once.
- **Space:** O(V + E) — adjacency list, in-degree map, and queue.

### Code

```go
func topoUniqueness(nums []int, sequences [][]int) bool {
	n := len(nums)

	adj := make(map[int]map[int]struct{}) // a -> set of b (dedup parallel edges)
	indeg := make(map[int]int)            // value -> number of prerequisites
	seen := make(map[int]bool)            // which values appear in seqs at all

	// Register a value node the first time we encounter it.
	touch := func(v int) {
		if !seen[v] {
			seen[v] = true
			indeg[v] = 0
			adj[v] = make(map[int]struct{})
		}
	}

	total := 0 // count of distinct values seen across seqs (must equal n)
	for _, seq := range sequences {
		for _, v := range seq {
			// Any value outside [1, n] means seqs can't reconstruct a
			// permutation of 1..n → impossible.
			if v < 1 || v > n {
				return false
			}
			if !seen[v] {
				touch(v)
				total++
			}
		}
		// Add an edge for every consecutive pair (dedup to keep in-degree exact).
		for i := 0; i+1 < len(seq); i++ {
			a, b := seq[i], seq[i+1]
			if _, exists := adj[a][b]; !exists {
				adj[a][b] = struct{}{}
				indeg[b]++ // b now has one more prerequisite
			}
		}
	}

	// If not every value 1..n showed up, nums cannot be reconstructed at all.
	if total != n {
		return false
	}

	// Seed queue with all sources (in-degree 0).
	queue := make([]int, 0, n)
	for v := range seen {
		if indeg[v] == 0 {
			queue = append(queue, v)
		}
	}

	index := 0 // position in nums we expect to emit next
	for len(queue) > 0 {
		// Unique next element required: more than one source ⇒ ambiguous order.
		if len(queue) > 1 {
			return false
		}
		node := queue[0]
		queue = queue[1:]

		// The forced next node must match nums at this position.
		if index >= n || node != nums[index] {
			return false
		}
		index++

		// Relax outgoing edges; enqueue neighbours that become sources.
		for nb := range adj[node] {
			indeg[nb]--
			if indeg[nb] == 0 {
				queue = append(queue, nb)
			}
		}
	}

	// Unique reconstruction iff we emitted all n values (order already checked).
	return index == n
}
```

### Dry Run

Example 1: `nums = [1,2,3]`, `sequences = [[1,2],[1,3]]`.

Edges: from `[1,2]` add `1→2`; from `[1,3]` add `1→3`. In-degrees: `1:0, 2:1, 3:1`. Distinct values = 3 = n.

| Step | queue (sources) | size | Action |
|------|-----------------|------|--------|
| seed | [1] | 1 | ok |
| 1 | pop 1 (== nums[0]) | — | relax: indeg[2]→0, indeg[3]→0; enqueue both → queue = [2,3] |
| 2 | [2,3] | **2** | size > 1 → ambiguous → **return false** |

Result: `false` ✔ — both `[1,2,3]` and `[1,3,2]` are valid, so `nums` isn't unique.

---

## Approach 2 — Adjacent-Coverage Check

### Intuition

`nums` is the unique shortest supersequence iff **every adjacent pair** `(nums[i], nums[i+1])` is directly forced by some sequence, and all values are accounted for. If even one adjacency `(nums[i], nums[i+1])` is never stated consecutively anywhere, those two neighbours could be swapped to form a different valid supersequence. Using `pos[v]` = index of `v` in `nums`, any consecutive pair `(a, b)` from a sequence must have `pos[a] < pos[b]` (else it contradicts `nums`); if it's *exactly* one apart (`pos[a] + 1 == pos[b]`), it "covers" that adjacency. Cover all `n − 1` adjacencies ⇒ unique.

### Algorithm

1. Build `pos` from `nums`.
2. For each consecutive pair `(a, b)` in `sequences`: validate values are in `[1, n]` and mark them seen. If `pos[a] > pos[b]`, contradiction → `false`. If `pos[a] + 1 == pos[b]`, mark that adjacency covered (decrement `toCover`).
3. Ensure every value `1..n` appeared; if any is missing → `false`.
4. Return `toCover == 0` (all adjacencies pinned).

### Complexity

- **Time:** O(V + E) — build `pos` in O(n), scan all pairs once.
- **Space:** O(n) — position map, coverage flags, and seen flags. No explicit graph, so lighter than Approach 1.

### Code

```go
func adjacentCoverage(nums []int, sequences [][]int) bool {
	n := len(nums)
	pos := make(map[int]int, n) // value -> its index in nums
	for i, v := range nums {
		pos[v] = i
	}

	covered := make([]bool, n)  // covered[i] = adjacency (nums[i], nums[i+1]) pinned
	toCover := n - 1            // number of adjacencies still needing a constraint
	sawValue := make([]bool, n+1) // sawValue[v] = value v appeared in seqs

	for _, seq := range sequences {
		for i := 0; i < len(seq); i++ {
			v := seq[i]
			// Value must index into nums (1..n); otherwise reconstruction fails.
			if v < 1 || v > n {
				return false
			}
			sawValue[v] = true

			if i+1 < len(seq) {
				a, b := seq[i], seq[i+1]
				if b < 1 || b > n {
					return false
				}
				// If a is supposed to come after b in nums, seqs contradict it.
				if pos[a] > pos[b] {
					return false
				}
				// Exactly adjacent in nums and not yet counted → cover it.
				if pos[a]+1 == pos[b] && !covered[pos[a]] {
					covered[pos[a]] = true
					toCover--
				}
			}
		}
	}

	// Every value 1..n must have appeared somewhere in seqs.
	for v := 1; v <= n; v++ {
		if !sawValue[v] {
			return false
		}
	}

	// Unique iff all n-1 adjacencies were pinned by some seq.
	return toCover == 0
}
```

### Dry Run

Example 1: `nums = [1,2,3]`, `sequences = [[1,2],[1,3]]`. `pos = {1:0, 2:1, 3:2}`; `toCover = 2` (adjacencies (1,2)@0 and (2,3)@1).

| seq | pair (a,b) | pos[a], pos[b] | pos[a] > pos[b]? | pos[a]+1 == pos[b]? | Action |
|-----|-----------|-----------------|-------------------|----------------------|--------|
| [1,2] | (1,2) | 0, 1 | no | 0+1 == 1 → yes | cover adj@0, toCover → 1 |
| [1,3] | (1,3) | 0, 2 | no | 0+1 == 2 → no | nothing (gap of 2) |

All values 1,2,3 seen. `toCover = 1 ≠ 0` → **return false** ✔ — adjacency (2,3) was never pinned, so `[1,3,2]` is also valid.

---

## Key Takeaways

- **"Unique topological order" ⇔ single source at every Kahn step.** This is the reusable test for *deterministic* orderings (also #269 Alien Dictionary uniqueness, task-scheduling determinism). The moment two in-degree-0 nodes coexist, the order is ambiguous.
- **Model ordering constraints as edges from consecutive pairs.** You do NOT need edges between *all* pairs of a sequence — consecutive pairs transitively imply the rest, and adding only them keeps `E` linear in the input size.
- **Dedup parallel edges** before counting in-degrees, or a repeated pair inflates an in-degree and corrupts the topo sort. Sequences here are unique, but a pair can still recur across different sequences.
- **Two lenses on the same fact:** either build the DAG and confirm a single order (Approach 1), or observe that uniqueness ⇔ every neighbour pair of `nums` is explicitly constrained (Approach 2). The second trades the graph for an O(n)-space coverage count.
- **Guard the totality condition:** even a perfectly consistent set of edges fails if some value `1..n` never appears — the supersequence would then not be pinned everywhere.

---

## Related Problems

- LeetCode #207 — Course Schedule (cycle detection via Kahn's / DFS)
- LeetCode #210 — Course Schedule II (produce a topological order)
- LeetCode #269 — Alien Dictionary (build a DAG from ordering hints, detect uniqueness/cycles)
- LeetCode #310 — Minimum Height Trees (peeling leaves, an in-degree-style BFS)
