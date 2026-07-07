# 0406 — Queue Reconstruction by Height

> LeetCode #406 · Difficulty: Medium
> **Categories:** Array, Greedy, Sorting, Binary Indexed Tree

---

## Problem Statement

You are given an array of people, `people`, which are the attributes of some people in a queue (not necessarily in order). Each `people[i] = [hi, ki]` represents the `ith` person of height `hi` with **exactly** `ki` other people in front who have a height greater than or equal to `hi`.

Reconstruct and return *the queue that is represented by the input array* `people`. The returned queue should be formatted as an array `queue`, where `queue[j] = [hj, kj]` is the attributes of the `jth` person in the queue (`queue[0]` is the person at the front of the queue).

**Example 1:**

```
Input: people = [[7,0],[4,4],[7,1],[5,0],[6,1],[5,2]]
Output: [[5,0],[7,0],[5,2],[6,1],[4,4],[7,1]]
Explanation:
Person 0 has height 5 with no other people taller or the same height in front.
Person 1 has height 7 with no other people taller or the same height in front.
Person 2 has height 5 with two persons taller or the same height in front, which is person 0 and 1.
Person 3 has height 6 with one person taller or the same height in front, which is person 1.
Person 4 has height 4 with four people taller or the same height in front, which are people 0, 1, 2, and 3.
Person 5 has height 7 with one person taller or the same height in front, which is person 1.
Hence [[5,0],[7,0],[5,2],[6,1],[4,4],[7,1]] is the reconstructed queue.
```

**Example 2:**

```
Input: people = [[6,0],[5,0],[4,0],[3,2],[2,2],[1,4]]
Output: [[4,0],[5,0],[2,2],[3,2],[1,4],[6,0]]
```

**Constraints:**

- `1 <= people.length <= 2000`
- `0 <= hi <= 10^6`
- `0 <= ki < people.length`
- It is guaranteed that the queue can be reconstructed.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★☆ High       | 2024          |
| Google     | ★★★☆☆ Medium     | 2024          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Meta       | ★★★☆☆ Medium     | 2023          |
| Bloomberg  | ★★☆☆☆ Low        | 2023          |
| Apple      | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Greedy** — the crux is a greedy ordering decision: fix people in a height order that makes each placement *locally final* (never invalidated by later placements). Sorting tallest-first turns `k` directly into an insertion index → see [`/dsa/greedy.md`](/dsa/greedy.md)
- **Custom Sorting** — both approaches hinge on a two-key comparator (height as the primary key, `k` as the tie-breaker) so equal-height people, which count toward each other, are ordered correctly → see [`/dsa/sorting.md`](/dsa/sorting.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (Slot Placement) | O(n²) | O(n) | Makes the "k-th empty slot" invariant explicit; good for intuition |
| 2 | Greedy Insertion (Tallest First) | O(n²) | O(n) | The standard interview answer; `k` becomes the insertion index |

Both are O(n²); the greedy insertion is the canonical solution. An O(n log n) variant exists using a Fenwick tree / order-statistics structure to find the k-th empty slot, but the quadratic greedy is what interviewers expect.

---

## Approach 1 — Brute Force (Slot Placement)

### Intuition

Process people **shortest first**. When we drop a short person into the queue, everyone already placed is shorter than (or, via tie-breaking, equal to) them — those people do **not** contribute to this person's `k`. The slots we leave empty will later be filled by **taller** people, and every taller person **does** contribute to `k`. Therefore the current person must sit in the empty slot preceded by exactly `k` empty slots: those `k` slots are reserved for the taller people that outrank them. Placing shortest-first guarantees each decision stays valid forever.

The tie-break matters: two people of equal height each count toward the other's `k`. Sorting equal heights by **larger `k` first** ensures the equal-height person who needs more people ahead is placed earlier and leaves the right room.

### Algorithm

1. Sort `people` ascending by height; break height ties by **descending** `k`.
2. Create `result`, a slice of `n` empty slots.
3. For each person `p` in sorted order, scan `result` left→right counting empty slots. When the number of empty slots passed equals `p.k`, place `p` in that empty slot.
4. Return `result`.

### Complexity

- **Time:** O(n²) — for each of the `n` people we may scan all `n` slots to find the `k`-th empty one.
- **Space:** O(n) — the result array (plus the in-place sort).

### Code

```go
func bruteForce(people [][]int) [][]int {
	n := len(people)
	// Sort ascending by height; ties broken by larger k first because two people
	// of the same height each count toward the other's k.
	sort.Slice(people, func(i, j int) bool {
		if people[i][0] == people[j][0] {
			return people[i][1] > people[j][1] // taller-or-equal tie → bigger k first
		}
		return people[i][0] < people[j][0] // shorter first
	})

	result := make([][]int, n) // fixed-size queue; nil entries = empty slots
	for _, p := range people {
		empties := 0 // how many empty slots we have passed so far
		for idx := 0; idx < n; idx++ {
			if result[idx] != nil {
				continue // occupied by a shorter person already placed — skip
			}
			// This slot is empty. If exactly k empties precede it, p belongs here:
			// its k taller people will land in those k earlier empty slots.
			if empties == p[1] {
				result[idx] = []int{p[0], p[1]} // claim the slot
				break
			}
			empties++ // count this empty slot and keep scanning
		}
	}
	return result
}
```

### Dry Run

Example 1: `people = [[7,0],[4,4],[7,1],[5,0],[6,1],[5,2]]`.

Sorted ascending by height, ties by larger `k`: `[[4,4],[5,2],[5,0],[6,1],[7,1],[7,0]]`.

`result` starts as `[_,_,_,_,_,_]` (six empty slots).

| Step | Person (h,k) | Need k-th empty slot | Placement | result after |
|------|--------------|----------------------|-----------|--------------|
| 1 | [4,4] | 4th empty → index 4 | slot 4 | `[_,_,_,_,(4,4),_]` |
| 2 | [5,2] | 2nd empty → index 2 | slot 2 (empties 0,1 at idx 0,1) | `[_,_,(5,2),_,(4,4),_]` |
| 3 | [5,0] | 0th empty → index 0 | slot 0 | `[(5,0),_,(5,2),_,(4,4),_]` |
| 4 | [6,1] | 1st empty → index 3 | empties: idx1(0), idx3(1) → slot 3 | `[(5,0),_,(5,2),(6,1),(4,4),_]` |
| 5 | [7,1] | 1st empty → index 5 | empties: idx1(0), idx5(1) → slot 5 | `[(5,0),_,(5,2),(6,1),(4,4),(7,1)]` |
| 6 | [7,0] | 0th empty → index 1 | slot 1 | `[(5,0),(7,0),(5,2),(6,1),(4,4),(7,1)]` |

Result: `[[5,0],[7,0],[5,2],[6,1],[4,4],[7,1]]` ✔

---

## Approach 2 — Greedy Insertion (Tallest First)

### Intuition

Reverse the order and place the **tallest** people first. Two facts make this clean:

1. Inserting a **shorter** person later can never change an already-placed taller person's `k`, because a shorter person is invisible to `k` (which only counts heights ≥ mine).
2. When we insert person `p`, **everyone already in the list is ≥ `p`'s height**. So the number of taller-or-equal people that must precede `p` is exactly `p.k` — meaning `p` goes at **index `k`** of the current list. Insert it there and the invariant holds automatically.

Tie-break: among equal heights, insert **smaller `k` first** so the person who belongs further front is put in before the equal-height person who belongs behind them.

### Algorithm

1. Sort `people` by **descending** height; break ties by **ascending** `k`.
2. Start with an empty `result`.
3. For each person `p` in that order, insert `p` into `result` at index `p.k`.
4. Return `result`.

### Complexity

- **Time:** O(n²) — each of the `n` insertions shifts up to `n` elements in the slice.
- **Space:** O(n) — the result slice.

### Code

```go
func greedyInsert(people [][]int) [][]int {
	// Tallest first; equal heights → smaller k first (smaller k must sit earlier).
	sort.Slice(people, func(i, j int) bool {
		if people[i][0] == people[j][0] {
			return people[i][1] < people[j][1] // equal height → smaller k first
		}
		return people[i][0] > people[j][0] // taller first
	})

	result := make([][]int, 0, len(people))
	for _, p := range people {
		k := p[1] // everyone already placed is >= p's height, so k IS the target index
		// Insert p at position k: grow by one, shift the tail right, drop p in.
		result = append(result, nil)   // extend length by one
		copy(result[k+1:], result[k:]) // shift elements from k rightward
		result[k] = []int{p[0], p[1]}  // place p exactly at index k
	}
	return result
}
```

### Dry Run

Example 1: `people = [[7,0],[4,4],[7,1],[5,0],[6,1],[5,2]]`.

Sorted descending by height, ties by smaller `k`: `[[7,0],[7,1],[6,1],[5,0],[5,2],[4,4]]`.

| Step | Person (h,k) | Insert at index k | result after |
|------|--------------|-------------------|--------------|
| 1 | [7,0] | 0 | `[[7,0]]` |
| 2 | [7,1] | 1 | `[[7,0],[7,1]]` |
| 3 | [6,1] | 1 | `[[7,0],[6,1],[7,1]]` |
| 4 | [5,0] | 0 | `[[5,0],[7,0],[6,1],[7,1]]` |
| 5 | [5,2] | 2 | `[[5,0],[7,0],[5,2],[6,1],[7,1]]` |
| 6 | [4,4] | 4 | `[[5,0],[7,0],[5,2],[6,1],[4,4],[7,1]]` |

Result: `[[5,0],[7,0],[5,2],[6,1],[4,4],[7,1]]` ✔

---

## Key Takeaways

- **Fix the "hard" dimension first.** When each element has two competing constraints (here: height and count-of-taller), sort by one so the other becomes a simple, order-independent decision. Placing tallest-first makes `k` collapse into an array index.
- **Monotone insertion invariant:** shorter people are invisible to taller people's counts, so once the tall skeleton is correct, later insertions never break it. Recognising this "later work can't disturb earlier work" property is the greedy proof.
- **Comparator tie-breaks carry real logic.** Equal-height people count toward each other, so the `k` tie-break direction (descending for shortest-first slot filling, ascending for tallest-first insertion) is not cosmetic — get it wrong and the reconstruction fails.
- The O(n²) slice insertion can be upgraded to **O(n log n)** with a Fenwick tree / order-statistics tree that finds the k-th empty position in O(log n), but interviewers almost always accept the quadratic greedy.

---

## Related Problems

- LeetCode #315 — Count of Smaller Numbers After Self (order statistics via BIT/merge)
- LeetCode #354 — Russian Doll Envelopes (sort one dimension, solve the other)
- LeetCode #452 — Minimum Number of Arrows to Burst Balloons (greedy after sorting)
- LeetCode #1030 — Matrix Cells in Distance Order (custom sort key)
