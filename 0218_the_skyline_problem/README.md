# 0218 — The Skyline Problem

> LeetCode #218 · Difficulty: Hard
> **Categories:** Array, Divide and Conquer, Binary Indexed Tree, Segment Tree, Line Sweep, Heap (Priority Queue), Ordered Set

---

## Problem Statement

A city's **skyline** is the outer contour of the silhouette formed by all the buildings in that city when viewed from a distance. Given the locations and heights of all the buildings, return *the skyline formed by these buildings collectively*.

The geometric information of each building is given in the array `buildings` where `buildings[i] = [lefti, righti, heighti]`:

- `lefti` is the x coordinate of the left edge of the `ith` building.
- `righti` is the x coordinate of the right edge of the `ith` building.
- `heighti` is the height of the `ith` building.

You may assume all buildings are perfect rectangles grounded on an absolutely flat surface at height `0`.

The skyline should be represented as a list of "key points" **sorted by their x-coordinate** in the form `[[x1,y1],[x2,y2],...]`. Each key point is the left endpoint of some horizontal segment in the skyline except the last point in the list, which always has a y-coordinate `0` and is used to mark the skyline's termination where the rightmost building ends. Any ground between the leftmost and rightmost buildings should be part of the skyline's contour.

**Note:** There must be no consecutive horizontal lines of equal height in the output skyline. For instance, `[...,[2 3],[4 5],[7 5],[11 5],[12 7],...]` is not acceptable; the three lines of height 5 should be merged into one in the final output as such: `[...,[2 3],[4 5],[12 7],...]`.

**Example 1:**

```
Input: buildings = [[2,9,10],[3,7,15],[5,12,12],[15,20,10],[19,24,8]]
Output: [[2,10],[3,15],[7,12],[12,0],[15,10],[20,8],[24,0]]
Explanation:
Figure A shows the buildings of the input.
Figure B shows the skyline formed by those buildings. The red points in figure B represent the key points in the output list.
```

**Example 2:**

```
Input: buildings = [[0,2,3],[2,5,3]]
Output: [[0,3],[5,0]]
```

**Constraints:**

- `1 <= buildings.length <= 10^4`
- `0 <= lefti < righti <= 2^31 - 1`
- `1 <= heighti <= 2^31 - 1`
- `buildings` is sorted by `lefti` in non-decreasing order.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★★☆ High       | 2024          |
| Amazon     | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Facebook   | ★★★☆☆ Medium     | 2023          |
| Bloomberg  | ★★☆☆☆ Low        | 2022          |
| Uber       | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Line Sweep** — process building edges left-to-right as events; the skyline only changes at edges → see [`/dsa/line_sweep.md`](/dsa/line_sweep.md)
- **Heap / Priority Queue** — a max-heap tracks the tallest currently-active building at the sweep line (with lazy deletion) → see [`/dsa/heap_priority_queue.md`](/dsa/heap_priority_queue.md)
- **Divide and Conquer** — split buildings, recurse, and merge two skylines like merge sort → see [`/dsa/divide_and_conquer.md`](/dsa/divide_and_conquer.md)
- **Hash Map** — supports lazy deletion by counting pending height removals → see [`/dsa/hash_map.md`](/dsa/hash_map.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Sweep Line + Max-Heap | O(n log n) | O(n) | The standard interview answer; event-driven and intuitive |
| 2 | Divide and Conquer | O(n log n) | O(n) | Elegant merge-sort analogue; great when you like recursion |

---

## Approach 1 — Sweep Line + Max-Heap

### Intuition

The silhouette can only step up or down at a building's **left or right edge**. So collect all edges as events and sweep a vertical line left-to-right. At each x we want the height of the tallest building currently straddling that x — that is a running maximum over a set that gains a member at each left edge and loses one at each right edge, i.e. a max-heap. A new key point appears exactly when this maximum changes.

The subtlety is tie-breaking at equal x. Encoding start heights as **negative** and sorting by `(x, code)` makes: starts precede ends at the same x (so a building beginning where another ends keeps the line from dipping); among simultaneous starts the taller is added first; among simultaneous ends the shorter is removed first. This is what suppresses the spurious point at `x=2` in Example 2. Go's heap has no random-access delete, so we use **lazy deletion**: an ended height is only counted as "to remove", and popped when it surfaces at the top.

### Algorithm

1. For each building `[L,R,H]` emit events `(L,-H)` (start) and `(R,+H)` (end).
2. Sort events by `x`, breaking ties by the height code.
3. Keep a max-heap of active heights (seed it with ground `0`) and a `removed` count map.
4. For each event: if a start, push `H`; if an end, increment `removed[H]`.
5. While the heap's top has a pending removal, pop it (lazy deletion).
6. Let `curMax` be the heap top; if `curMax != prevMax`, append `[x, curMax]` and update `prevMax`.

### Complexity

- **Time:** O(n log n) — sorting `2n` events plus O(log n) heap work per event.
- **Space:** O(n) — events, heap, and removal map.

### Code

```go
func sweepLineHeap(buildings [][]int) [][]int {
	// Each event is [x, heightCode]: negative code = building start (height H),
	// positive code = building end (height H).
	events := make([][2]int, 0, len(buildings)*2)
	for _, b := range buildings {
		L, R, H := b[0], b[1], b[2]
		events = append(events, [2]int{L, -H}) // start: negative height
		events = append(events, [2]int{R, H})  // end: positive height
	}
	// Sort by x; ties broken by heightCode so that at the same x:
	//   - starts (negative) precede ends (positive),
	//   - among starts, the taller (more negative) is processed first,
	//   - among ends, the shorter is processed first.
	sort.Slice(events, func(i, j int) bool {
		if events[i][0] != events[j][0] {
			return events[i][0] < events[j][0]
		}
		return events[i][1] < events[j][1]
	})

	result := make([][]int, 0)
	pq := &maxHeap{0}                    // active heights; ground 0 always present
	removed := make(map[int]int)         // lazy-deletion counts: height → pending removals
	prevMax := 0                         // skyline height before this event

	for _, e := range events {
		x, code := e[0], e[1]
		if code < 0 {
			h := -code             // a building starts here at height h
			heap.Push(pq, h)       // it is now active
		} else {
			removed[code]++        // schedule removal of this ended height
		}
		// Pop off any heights whose removal is pending and that sit on top.
		for pq.Len() > 0 {
			top := (*pq)[0]
			if removed[top] > 0 { // this top height has actually ended
				removed[top]--
				heap.Pop(pq)
			} else {
				break // the real current max is on top
			}
		}
		curMax := (*pq)[0] // tallest active building (0 if none)
		if curMax != prevMax {
			// Height changed at x → this is a key point of the skyline.
			result = append(result, []int{x, curMax})
			prevMax = curMax
		}
	}
	return result
}
```

### Dry Run

Example 1: `buildings = [[2,9,10],[3,7,15],[5,12,12],[15,20,10],[19,24,8]]`.

Events sorted: `(2,-10) (3,-15) (5,-12) (7,+15) (9,+10) (12,+12) (15,-10) (19,-8) (20,+10) (24,+8)`.

| x | event | heap (active heights) after cleanup | curMax | prevMax | emit? |
|---|-------|--------------------------------------|--------|---------|-------|
| 2 | start 10 | {10,0} | 10 | 0 | **[2,10]** |
| 3 | start 15 | {15,10,0} | 15 | 10 | **[3,15]** |
| 5 | start 12 | {15,12,10,0} | 15 | 15 | no |
| 7 | end 15 | {12,10,0} | 12 | 15 | **[7,12]** |
| 9 | end 10 | {12,0} | 12 | 12 | no |
| 12 | end 12 | {0} | 0 | 12 | **[12,0]** |
| 15 | start 10 | {10,0} | 10 | 0 | **[15,10]** |
| 19 | start 8 | {10,8,0} | 10 | 10 | no |
| 20 | end 10 | {8,0} | 8 | 10 | **[20,8]** |
| 24 | end 8 | {0} | 0 | 8 | **[24,0]** |

Result: `[[2,10],[3,15],[7,12],[12,0],[15,10],[20,8],[24,0]]` ✔

---

## Approach 2 — Divide and Conquer

### Intuition

One building's skyline is trivially `[[L,H],[R,0]]`. Two skylines merge exactly like the merge step of merge sort: sweep both key-point lists together by x, and at each x the combined height is `max(leftHeight, rightHeight)`, emitting a point only when that max changes. Recurse on halves and merge — `T(n) = 2T(n/2) + O(n)`.

### Algorithm

1. Base cases: 0 buildings → empty; 1 building → `[[L,H],[R,0]]`.
2. Split the buildings in half; recurse to get `left` and `right` skylines.
3. Merge with two pointers: track `leftH`/`rightH` (each side's current height). At the smaller x, advance that side and update its height; on a tie advance both. The merged height is `max(leftH, rightH)`; append `[x, h]` only when `h` differs from the last emitted height.

### Complexity

- **Time:** O(n log n) — `log n` levels of recursion, O(n) merge work per level.
- **Space:** O(n) — recursion stack and the skyline lists.

### Code

```go
func divideAndConquer(buildings [][]int) [][]int {
	if len(buildings) == 0 {
		return [][]int{}
	}
	if len(buildings) == 1 {
		b := buildings[0]
		// One building becomes two key points: rise at L, fall to 0 at R.
		return [][]int{{b[0], b[2]}, {b[1], 0}}
	}
	mid := len(buildings) / 2
	left := divideAndConquer(buildings[:mid])  // skyline of the left half
	right := divideAndConquer(buildings[mid:]) // skyline of the right half
	return mergeSkylines(left, right)
}

func mergeSkylines(left, right [][]int) [][]int {
	merged := make([][]int, 0, len(left)+len(right))
	i, j := 0, 0
	leftH, rightH := 0, 0 // current height contributed by each skyline
	for i < len(left) && j < len(right) {
		var x, h int
		if left[i][0] < right[j][0] { // left key point comes first
			x = left[i][0]
			leftH = left[i][1] // update left's current height
			i++
		} else if left[i][0] > right[j][0] { // right key point comes first
			x = right[j][0]
			rightH = right[j][1]
			j++
		} else { // same x: consume both
			x = left[i][0]
			leftH = left[i][1]
			rightH = right[j][1]
			i++
			j++
		}
		h = max(leftH, rightH)                       // combined skyline height at x
		if len(merged) == 0 || merged[len(merged)-1][1] != h {
			merged = append(merged, []int{x, h})     // only emit real changes
		}
	}
	// Append the leftovers (only one list can have remaining points).
	for i < len(left) {
		x, h := left[i][0], left[i][1]
		if len(merged) == 0 || merged[len(merged)-1][1] != h {
			merged = append(merged, []int{x, h})
		}
		i++
	}
	for j < len(right) {
		x, h := right[j][0], right[j][1]
		if len(merged) == 0 || merged[len(merged)-1][1] != h {
			merged = append(merged, []int{x, h})
		}
		j++
	}
	return merged
}
```

### Dry Run

Example 1, illustrating one merge. Say the left half `[[2,9,10],[3,7,15]]` produced skyline `L = [[2,10],[3,15],[7,10],[9,0]]` and the right half `[[5,12,12],[15,20,10],[19,24,8]]` produced `R = [[5,12],[12,0],[15,10],[20,8],[24,0]]`. Merging L and R:

| step | x | leftH | rightH | max | emit? (vs last) |
|------|---|-------|--------|-----|-----------------|
| 1 | 2 | 10 | 0 | 10 | **[2,10]** |
| 2 | 3 | 15 | 0 | 15 | **[3,15]** |
| 3 | 5 | 15 | 12 | 15 | no (still 15) |
| 4 | 7 | 10 | 12 | 12 | **[7,12]** |
| 5 | 9 | 0 | 12 | 12 | no |
| 6 | 12 | 0 | 0 | 0 | **[12,0]** |
| 7 | 15 | 0 | 10 | 10 | **[15,10]** |
| 8 | 20 | 0 | 8 | 8 | **[20,8]** |
| 9 | 24 | 0 | 0 | 0 | **[24,0]** |

Result: `[[2,10],[3,15],[7,12],[12,0],[15,10],[20,8],[24,0]]` ✔

---

## Key Takeaways

- **The skyline changes only at edges** — reduce a continuous silhouette to a discrete set of events at building sides. This "only endpoints matter" idea powers interval and sweep-line problems generally.
- **Tie-breaking is the whole game.** Encoding starts as negative heights and sorting `(x, code)` guarantees starts-before-ends, taller-start-first, shorter-end-first — which prevents both false dips (Example 2) and duplicate points.
- **Lazy deletion** lets a plain binary heap simulate a delete-arbitrary priority queue: mark for removal, and discard only when the element surfaces at the top.
- **Merging skylines == merge step of merge sort** on `[x, height]` points, taking `max` of the two live heights — a clean recursive alternative with the same O(n log n).
- **Always suppress consecutive equal heights** in the output — both approaches emit a point only when the running height actually changes.

---

## Related Problems

- LeetCode #56 — Merge Intervals (sweep / interval merging)
- LeetCode #253 — Meeting Rooms II (sweep line + heap of active intervals)
- LeetCode #1272 — Remove Interval
- LeetCode #732 — My Calendar III (max concurrent intervals, sweep)
- LeetCode #23 — Merge k Sorted Lists (divide and conquer merge pattern)
