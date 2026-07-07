# 0382 — Linked List Random Node

> LeetCode #382 · Difficulty: Medium
> **Categories:** Linked List, Math, Reservoir Sampling, Randomized, Design

---

## Problem Statement

Given a singly linked list, return a random node's value from the linked list. Each node must have the **same probability** of being chosen.

Implement the `Solution` class:

- `Solution(ListNode head)` Initializes the object with the head of the singly linked list `head`.
- `int getRandom()` Chooses a node randomly from the list and returns its value. All the nodes of the list should be equally likely to be chosen.

**Example 1:**
```
Input
["Solution", "getRandom", "getRandom", "getRandom", "getRandom", "getRandom"]
[[[1, 2, 3]], [], [], [], [], []]
Output
[null, 1, 3, 2, 2, 3]

Explanation
Solution solution = new Solution([1, 2, 3]);
solution.getRandom(); // return 1
solution.getRandom(); // return 3
solution.getRandom(); // return 2
solution.getRandom(); // return 2
solution.getRandom(); // return 3
// getRandom() should return either 1, 2, or 3 randomly. Each element should have equal probability of returning.
```

**Constraints:**
- The number of nodes in the linked list will be in the range `[1, 10⁴]`.
- `-10⁴ <= Node.val <= 10⁴`
- At most `10⁴` calls will be made to `getRandom`.

**Follow up:** What if the linked list is extremely large and its length is unknown to you? Could you solve this efficiently without using extra space?

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★★☆ High       | 2024          |
| Amazon     | ★★★☆☆ Medium     | 2024          |
| Facebook   | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2023          |
| Bloomberg  | ★★☆☆☆ Low        | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Reservoir Sampling (size 1)** — pick a uniform element from a stream of unknown length in O(1) space; the intended follow-up answer → see [`/dsa/reservoir_sampling.md`](/dsa/reservoir_sampling.md)
- **Linked List Traversal** — single-pass walk over a singly linked list with no random access → see [`/dsa/linked_list.md`](/dsa/linked_list.md)
- **Design / Randomized Data Structure** — a stateful class exposing a randomized query, judged by distribution not a single value → see [`/dsa/design_data_structures.md`](/dsa/design_data_structures.md)
- **Probability / Math** — the 1/i replacement argument proving uniformity by induction → see [`/dsa/math_number_theory.md`](/dsa/math_number_theory.md)

---

## Approaches Overview

Let n = number of nodes in the list.

| # | Approach | Time (getRandom) | Space | When to use |
|---|----------|------------------|-------|-------------|
| 1 | Array Snapshot | O(1) | O(n) | Length known/bounded and many getRandom calls |
| 2 | Reservoir Sampling (Optimal) | O(n) | O(1) | Unknown/huge length; the follow-up answer |

---

## Approach 1 — Array Snapshot

### Intuition
A linked list has no O(1) random access, but a slice does. Pay one O(n) pass at construction to snapshot the values; afterwards each `getRandom` is a single `rand.Intn(n)` index — uniform because every index is equally likely. The cost is O(n) extra memory and needing the length up front.

### Algorithm
1. Constructor: walk the list, appending each `Val` to `vals`.
2. `getRandom`: return `vals[rand.Intn(len(vals))]`.

### Complexity
- **Time:** Constructor O(n); getRandom O(1) — direct index.
- **Space:** O(n) — the value snapshot.

### Code
```go
type ArraySolution struct {
	vals []int      // snapshot of all node values
	rng  *rand.Rand // deterministic RNG for reproducible demos
}

func NewArraySolution(head *ListNode, rng *rand.Rand) *ArraySolution {
	vals := []int{}
	for n := head; n != nil; n = n.Next {
		vals = append(vals, n.Val) // record every value once
	}
	return &ArraySolution{vals: vals, rng: rng}
}

func (s *ArraySolution) GetRandom() int {
	return s.vals[s.rng.Intn(len(s.vals))] // every index equally likely
}
```

### Dry Run
List `[1,2,3]`.

| Step | action | state |
|------|--------|-------|
| construct | walk 1→2→3 | `vals = [1,2,3]` |
| getRandom | `rand.Intn(3)` → 0 | return `vals[0] = 1` |
| getRandom | `rand.Intn(3)` → 2 | return `vals[2] = 3` |
| getRandom | `rand.Intn(3)` → 1 | return `vals[1] = 2` |

Each index has probability 1/3, so each value is uniform.

---

## Approach 2 — Reservoir Sampling (Optimal)

### Intuition
Walk the list once, keeping a single "chosen" value. When we meet the i-th node (1-indexed), replace `chosen` with its value with probability 1/i. By induction, after seeing k nodes every one is the current choice with probability exactly 1/k — so at the end each of the n nodes has probability 1/n. No length and no array needed: this answers the follow-up.

**Why 1/i works (induction):** suppose after i−1 nodes each is chosen with prob 1/(i−1). The i-th node becomes chosen with prob 1/i. Any earlier node stays chosen with prob (its 1/(i−1)) × (1 − 1/i) = 1/(i−1) × (i−1)/i = 1/i. So all i nodes are uniform at 1/i.

### Algorithm
1. `chosen = head.Val`, `i = 1`.
2. For each subsequent node: `i++`; with probability 1/i set `chosen = node.Val` (i.e. `rng.Intn(i) == 0`).
3. Return `chosen`.

### Complexity
- **Time:** Constructor O(1) (store head); getRandom O(n) — one pass.
- **Space:** O(1) — a single running candidate.

### Code
```go
type ReservoirSolution struct {
	head *ListNode
	rng  *rand.Rand
}

func NewReservoirSolution(head *ListNode, rng *rand.Rand) *ReservoirSolution {
	return &ReservoirSolution{head: head, rng: rng}
}

func (s *ReservoirSolution) GetRandom() int {
	chosen := s.head.Val // seed the reservoir with the first value
	i := 1               // number of nodes seen so far
	for n := s.head.Next; n != nil; n = n.Next {
		i++                     // now considering the i-th node
		if s.rng.Intn(i) == 0 { // pick it with probability 1/i
			chosen = n.Val
		}
	}
	return chosen
}
```

### Dry Run
List `[1,2,3]`, one `getRandom` call.

| Step | node | i | rng.Intn(i) | pick? | chosen |
|------|------|---|-------------|-------|--------|
| seed | 1 | 1 | — | — | 1 |
| iter | 2 | 2 | e.g. 1 (≠0) | no | 1 |
| iter | 3 | 3 | e.g. 0 (=0) | yes | 3 |
| end | — | — | — | — | return 3 |

Over many trials, node 1 keeps prob 1/3, node 2 gets 1/3, node 3 gets 1/3.

---

## Key Takeaways
- **Reservoir sampling size 1** is the go-to for "uniform element from a stream of unknown length in O(1) space". The rule: replace the kept item with the i-th item with probability 1/i.
- The array snapshot trades O(n) memory for O(1) queries — better when the list is small and queries are frequent.
- Randomized structures are verified by **distribution over many trials**, not a single return value.
- The 1/i uniformity proof generalizes to reservoir size k (replace with prob k/i).

---

## Related Problems
- LeetCode #398 — Random Pick Index (reservoir sampling over matching indices)
- LeetCode #384 — Shuffle an Array (Fisher–Yates, sibling randomized problem)
- LeetCode #528 — Random Pick with Weight (prefix sums + binary search)
- LeetCode #710 — Random Pick with Blacklist (remapping)
