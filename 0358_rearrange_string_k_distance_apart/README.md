# 0358 — Rearrange String k Distance Apart

> LeetCode #358 · Difficulty: Hard
> **Categories:** Hash Table, String, Greedy, Sorting, Heap (Priority Queue), Counting

---

## Problem Statement

Given a string `s` and an integer `k`, rearrange `s` such that the same characters are **at least** distance `k` from each other. If it is not possible to rearrange the string, return an empty string `""`.

**Example 1:**

```
Input: s = "aabbcc", k = 3
Output: "abcabc"
Explanation: The same letters are at least a distance of 3 from each other.
```

**Example 2:**

```
Input: s = "aaabc", k = 3
Output: ""
Explanation: It is not possible to rearrange the string.
```

**Example 3:**

```
Input: s = "aaadbbcc", k = 2
Output: "abacabcd"
Explanation: The same letters are at least a distance of 2 from each other.
```

**Constraints:**

- `1 <= s.length <= 3 * 10^5`
- `s` consists of only lowercase English letters.
- `0 <= k <= s.length`

> Note: any valid rearrangement is accepted — the expected outputs above are one valid answer each.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★★☆ High       | 2024          |
| Amazon     | ★★★☆☆ Medium     | 2023          |
| Facebook   | ★★★☆☆ Medium     | 2023          |
| Uber       | ★★☆☆☆ Low        | 2022          |
| Microsoft  | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Greedy** — always emit the most-frequent still-available character; the scarcest resource is placed first → see [`/dsa/greedy.md`](/dsa/greedy.md)
- **Heap / Priority Queue** — a max-heap keyed by remaining count supplies the next character in O(log a) → see [`/dsa/heap_priority_queue.md`](/dsa/heap_priority_queue.md)
- **Hash Map (Counting)** — frequency table drives both feasibility and placement → see [`/dsa/hash_map.md`](/dsa/hash_map.md)
- **Queue / Cooldown** — a FIFO queue of size `k` enforces the "at least k apart" constraint → see [`/dsa/queue_deque.md`](/dsa/queue_deque.md)
- **Sorting** — round-robin variant sorts characters by descending frequency → see [`/dsa/sorting.md`](/dsa/sorting.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Greedy Max-Heap + Cooldown | O(n log a) | O(a) | Canonical; generalises to "task scheduler" |
| 2 | Greedy Round-Robin (sort by count) | O(n + a log a) | O(n) | No heap; slot-filling by columns |

*(a = alphabet size ≤ 26.)*

---

## Approach 1 — Greedy Max-Heap + Cooldown

### Intuition

Spreading characters `k` apart is a scheduling problem. The character that is **hardest** to place is the most frequent one, so greedily emit whichever available character has the highest remaining count. After emitting a character it must wait out `k-1` more slots, so park it in a FIFO cooldown queue of size `k`. When the queue fills to size `k`, its front is eligible again and (if it still has count) returns to the heap. If the heap empties while characters are still cooling down, no valid arrangement exists → return `""`.

### Algorithm

1. Count frequencies; push every `(char, count)` into a max-heap keyed by count.
2. While the heap is non-empty: pop the top, append its char, decrement its count, and enqueue it onto the cooldown queue.
3. Once the cooldown queue has `k` entries, dequeue the front; if it still has `count > 0`, push it back to the heap.
4. If the result length equals `len(s)`, return it; else return `""`.

### Complexity

- **Time:** O(n log a) — `n` emissions, each with a heap operation over ≤ `a` distinct characters.
- **Space:** O(a) — heap plus cooldown queue.

### Code

```go
func maxHeapGreedy(s string, k int) string {
	if k <= 1 {
		return s // any arrangement (or the string itself) already satisfies k<=1
	}
	freq := map[byte]int{}
	for i := 0; i < len(s); i++ {
		freq[s[i]]++ // tally each character
	}

	h := &charHeap{}
	for c, f := range freq {
		heap.Push(h, charCount{c, f}) // seed heap with every distinct char
	}

	type cooling struct {
		char  byte
		count int
	}
	queue := []cooling{} // FIFO of characters waiting out their cooldown
	result := make([]byte, 0, len(s))

	for h.Len() > 0 {
		top := heap.Pop(h).(charCount) // most frequent available char
		result = append(result, top.char)
		top.count-- // one occurrence consumed
		queue = append(queue, cooling{top.char, top.count})

		if len(queue) >= k {
			front := queue[0]
			queue = queue[1:]
			if front.count > 0 {
				heap.Push(h, charCount{front.char, front.count})
			}
		}
	}

	if len(result) != len(s) {
		return "" // ran out of eligible chars before placing all — impossible
	}
	return string(result)
}
```

### Dry Run

Example 1: `s = "aabbcc", k = 3`. Frequencies `a:2, b:2, c:2`.

| Step | Heap top (emit) | result | cooldown queue (front→back) | queue size ≥ k? release |
|------|-----------------|--------|-----------------------------|-------------------------|
| 1 | a (2→1) | `a` | [a:1] | no |
| 2 | b (2→1) | `ab` | [a:1, b:1] | no |
| 3 | c (2→1) | `abc` | [a:1, b:1, c:1] | size 3 → release a:1 back to heap |
| 4 | a (1→0) | `aba` | [b:1, c:1, a:0] | size 3 → release b:1 to heap |
| 5 | b (1→0) | `abab` | [c:1, a:0, b:0] | size 3 → release c:1 to heap |
| 6 | c (1→0) | `ababc`… | [a:0, b:0, c:0] | (heap empties after) |

Result length 6 == input length ⇒ a valid arrangement like `abcabc` (this trace yields `ababcc`-style order depending on tie-breaks; all valid). Result: **valid** ✔

---

## Approach 2 — Greedy Round-Robin (Sort by Count)

### Intuition

Feasibility depends only on the most frequent character `m` with count `f`: it needs `(f-1)` gaps of size `≥ k`, i.e. at least `(f-1)*k + 1` positions. If `n < (f-1)*k + 1`, impossible. Otherwise, sort characters by descending count and pour them into the output array by **columns**: fill indices `0, k, 2k, …` first (consecutive characters), then wrap to `1, k+1, …`, then `2, k+2, …`. Because a character's occurrences are always placed `k` slots apart within a column stride, equal characters never land closer than `k`.

### Algorithm

1. Count characters and sort by descending frequency.
2. Feasibility: if `(maxFreq-1)*k + 1 > n`, return `""`.
3. Walk an index that jumps by `k` per placement; when it overflows `n`, wrap to `idx%k + 1` (next start column). Assign characters in frequency order.
4. Join the filled array.

### Complexity

- **Time:** O(n + a log a) — counting + small-alphabet sort + linear fill.
- **Space:** O(n) — output buffer.

### Code

```go
func roundRobinGreedy(s string, k int) string {
	if k <= 1 {
		return s
	}
	freq := map[byte]int{}
	for i := 0; i < len(s); i++ {
		freq[s[i]]++
	}
	pairs := make([]charCount, 0, len(freq))
	for c, f := range freq {
		pairs = append(pairs, charCount{c, f})
	}
	for i := 1; i < len(pairs); i++ { // insertion sort by count desc
		for j := i; j > 0 && pairs[j].count > pairs[j-1].count; j-- {
			pairs[j], pairs[j-1] = pairs[j-1], pairs[j]
		}
	}

	n := len(s)
	maxFreq := pairs[0].count
	if (maxFreq-1)*k+1 > n {
		return "" // cannot spread the top char far enough
	}

	res := make([]byte, n)
	idx := 0 // current slot to fill; jumps by k each time
	for _, p := range pairs {
		for c := 0; c < p.count; c++ {
			res[idx] = p.char // place one occurrence
			idx += k          // next occurrence is k away
			if idx >= n {
				idx = idx%k + 1 // wrap to the next start column (0→1→2…)
			}
		}
	}
	return string(res)
}
```

### Dry Run

Example 1: `s = "aabbcc", k = 3`, `n = 6`. Sorted (ties arbitrary): `a:2, b:2, c:2`. `maxFreq = 2`; feasibility `(2-1)*3+1 = 4 ≤ 6` ✔.

| Place | char | idx before | write | idx after | wrap? |
|-------|------|-----------|-------|-----------|-------|
| 1 | a | 0 | res[0]=a | 3 | no |
| 2 | a | 3 | res[3]=a | 6 → wrap `6%3+1 = 1` | yes |
| 3 | b | 1 | res[1]=b | 4 | no |
| 4 | b | 4 | res[4]=b | 7 → wrap `7%3+1 = 2` | yes |
| 5 | c | 2 | res[2]=c | 5 | no |
| 6 | c | 5 | res[5]=c | 8 → wrap | — |

`res = a b c a b c` → **"abcabc"** ✔ — every letter is exactly 3 apart.

---

## Key Takeaways

- **Most-frequent-first is the correct greedy** for spread/scheduling problems: place the scarcest resource before it becomes impossible to place. Same idea powers Task Scheduler (#621) and Reorganize String (#767).
- **Cooldown queue of size k** cleanly enforces "≥ k apart": a char re-enters the pool only after `k` other emissions.
- **Feasibility test up front**: `(maxFreq-1)*k + 1 ≤ n`. If the top character can't be spread, nothing can.
- **Round-robin column fill** is an elegant heap-free alternative: sort by count, drop into slots `0,k,2k,…` wrapping through start columns — guarantees the k-gap by construction.

---

## Related Problems

- LeetCode #767 — Reorganize String (special case `k = 2`)
- LeetCode #621 — Task Scheduler (cooldown scheduling, same greedy)
- LeetCode #1054 — Distant Barcodes (k = 2 barcode variant)
- LeetCode #451 — Sort Characters By Frequency (frequency + heap/sort)
