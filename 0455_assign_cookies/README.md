# 0455 — Assign Cookies

> LeetCode #455 · Difficulty: Easy
> **Categories:** Array, Two Pointers, Greedy, Sorting

---

## Problem Statement

Assume you are an awesome parent and want to give your children some cookies. But, you should give each child at most one cookie.

Each child `i` has a greed factor `g[i]`, which is the minimum size of a cookie that the child will be content with; and each cookie `j` has a size `s[j]`. If `s[j] >= g[i]`, we can assign the cookie `j` to the child `i`, and the child `i` will be content. Your goal is to maximize the number of your content children and output the maximum number.

**Example 1:**

```
Input: g = [1,2,3], s = [1,1]
Output: 1
Explanation: You have 3 children and 2 cookies. The greed factors of 3 children are 1, 2, 3.
And even though you have 2 cookies, since their size is both 1, you could only make the child whose greed factor is 1 content.
You need to output 1.
```

**Example 2:**

```
Input: g = [1,2], s = [1,2,3]
Output: 2
Explanation: You have 2 children and 3 cookies. The greed factors of 2 children are 1, 2.
You have 3 cookies and their sizes are big enough to gratify all of the children,
You need to output 2.
```

**Constraints:**

- `1 <= g.length <= 3 * 10^4`
- `0 <= s.length <= 3 * 10^4`
- `1 <= g[i], s[j] <= 2^31 - 1`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★☆☆ Medium     | 2023          |
| Google     | ★★☆☆☆ Low        | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2022          |
| Apple      | ★★☆☆☆ Low        | 2022          |
| Uber       | ★☆☆☆☆ Rare       | 2021          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Greedy** — the core: content the least-greedy child with the smallest adequate cookie, never wasting a big cookie where a small one suffices. An exchange argument proves this local choice is globally optimal → see [`/dsa/greedy.md`](/dsa/greedy.md)
- **Two Pointers** — after sorting, one pointer walks children and the other walks cookies in a single synchronized sweep, each advancing only forward → see [`/dsa/two_pointers.md`](/dsa/two_pointers.md)
- **Sorting** — both greed factors and cookie sizes are sorted ascending so that "smallest fit" and "largest fit" become simple linear scans → see [`/dsa/sorting.md`](/dsa/sorting.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Greedy Two Pointers, Smallest-Fit (Optimal) | O(n log n + m log m) | O(1) extra | The canonical answer; assign each child the smallest cookie that fits |
| 2 | Greedy From the Largest (Biggest-Fit) | O(n log n + m log m) | O(1) extra | Mirror image; match the greediest child with the largest cookie |

---

## Approach 1 — Greedy Two Pointers, Smallest-Fit (Optimal)

### Intuition

To maximise contented children, spend cookies frugally: give a child the *smallest* cookie that still satisfies them, saving larger cookies for greedier children. Sort children by greed and cookies by size, both ascending. Walk both with two pointers. For the current least-greedy unserved child, skip cookies that are too small (a cookie too small for this child is too small for every greedier child left, so it is useless) until one fits — assign it, advance both pointers. Repeat until you run out of children or cookies. The exchange argument: if an optimal solution gave this child a larger cookie, swapping in the smallest adequate one frees the larger cookie without reducing the count — so greedy is optimal.

### Algorithm

1. Sort `g` and `s` ascending.
2. Set `child = 0`, `cookie = 0`, `count = 0`.
3. While `child < len(g)` and `cookie < len(s)`:
   - If `s[cookie] >= g[child]`: assign → `count++`, `child++`, `cookie++`.
   - Else: cookie too small → `cookie++`.
4. Return `count`.

### Complexity

- **Time:** O(n log n + m log m) — sorting both arrays dominates; the two-pointer sweep is O(n + m).
- **Space:** O(1) extra — sorting is in place; only pointers and a counter.

### Code

```go
func greedyTwoPointers(g []int, s []int) int {
	sort.Ints(g) // children by ascending greed (easiest to please first)
	sort.Ints(s) // cookies by ascending size (spend the smallest first)

	child := 0 // index into g: the current least-greedy unserved child
	cookie := 0 // index into s: the current smallest unused cookie
	count := 0 // contented children so far
	for child < len(g) && cookie < len(s) {
		if s[cookie] >= g[child] {
			// This cookie satisfies the current child: assign it and advance
			// both — this child is done, this cookie is spent.
			count++
			child++
			cookie++
		} else {
			// Cookie too small for the least greedy remaining child, hence too
			// small for everyone left; drop it and try the next larger cookie.
			cookie++
		}
	}
	return count
}
```

### Dry Run

Example 1: `g = [1,2,3]`, `s = [1,1]`. Already sorted: `g = [1,2,3]`, `s = [1,1]`.

| Step | child (g[child]) | cookie (s[cookie]) | s[cookie] >= g[child]? | Action | count |
|------|------------------|--------------------|------------------------|--------|-------|
| 1 | 0 (1) | 0 (1) | yes (1 ≥ 1) | assign; child→1, cookie→1 | 1 |
| 2 | 1 (2) | 1 (1) | no (1 < 2) | skip cookie; cookie→2 | 1 |
| 3 | 1 (2) | 2 | cookie out of range | loop ends | 1 |

Result: `1` ✔ — only the greed-1 child gets a cookie.

---

## Approach 2 — Greedy From the Largest (Biggest-Fit)

### Intuition

The symmetric greedy. Sort both ascending but walk from the back. Take the greediest remaining child and the largest remaining cookie. If that cookie contents the child, it is a match (consume both). If even the largest cookie cannot content this child, then *nothing* can (this is the greediest child and that was the biggest cookie), so abandon this child and move to the next-greediest — keeping the cookie for a less greedy child. Same optimal count; some find "if the biggest can't please the greediest, that child is hopeless" the clearer invariant.

### Algorithm

1. Sort `g` and `s` ascending.
2. Set `child = len(g)-1`, `cookie = len(s)-1`, `count = 0`.
3. While `child >= 0` and `cookie >= 0`:
   - If `s[cookie] >= g[child]`: match → `count++`, `child--`, `cookie--`.
   - Else: child unservable → `child--` (keep the cookie).
4. Return `count`.

### Complexity

- **Time:** O(n log n + m log m) — sorting dominates; the sweep is O(n + m).
- **Space:** O(1) extra + O(log) sort stack.

### Code

```go
func greedyLargestFirst(g []int, s []int) int {
	sort.Ints(g)
	sort.Ints(s)

	child := len(g) - 1  // greediest child
	cookie := len(s) - 1 // largest cookie
	count := 0
	for child >= 0 && cookie >= 0 {
		if s[cookie] >= g[child] {
			// Largest remaining cookie contents the greediest remaining child.
			count++
			child--
			cookie--
		} else {
			// Even the biggest cookie can't satisfy this (greediest) child, so
			// nothing can; this child stays unhappy, keep the cookie for a less
			// greedy child.
			child--
		}
	}
	return count
}
```

### Dry Run

Example 1: `g = [1,2,3]`, `s = [1,1]`. Sorted; start `child = 2 (g=3)`, `cookie = 1 (s=1)`.

| Step | child (g[child]) | cookie (s[cookie]) | s[cookie] >= g[child]? | Action | count |
|------|------------------|--------------------|------------------------|--------|-------|
| 1 | 2 (3) | 1 (1) | no (1 < 3) | child hopeless; child→1 | 0 |
| 2 | 1 (2) | 1 (1) | no (1 < 2) | child hopeless; child→0 | 0 |
| 3 | 0 (1) | 1 (1) | yes (1 ≥ 1) | match; child→-1, cookie→0 | 1 |
| 4 | -1 | 0 | child out of range | loop ends | 1 |

Result: `1` ✔.

---

## Key Takeaways

- **Greedy matching on two sorted sequences.** Sort both sides, then sweep with two pointers — the standard shape for "assign resources to requests to maximise satisfied requests."
- **Smallest-fit vs biggest-fit are dual and both optimal.** Assign the least-greedy child the smallest adequate cookie, *or* the greediest child the biggest cookie; the exchange argument works either direction.
- **Skip, don't backtrack.** A cookie too small for the current least-greedy child is useless for all remaining children, so discard it and move on — every pointer only ever advances, keeping the sweep linear.
- Guard the empty-`s` case naturally: the `cookie < len(s)` bound makes `s = []` return `0` with no special-casing.

---

## Related Problems

- LeetCode #870 — Advantage Shuffle (greedy assignment on sorted arrays)
- LeetCode #1005 — Maximize Sum Of Array After K Negations (greedy over sorted values)
- LeetCode #881 — Boats to Save People (two-pointer greedy pairing)
- LeetCode #452 — Minimum Number of Arrows to Burst Balloons (greedy + sorting on intervals)
- LeetCode #135 — Candy (greedy distribution with constraints)
