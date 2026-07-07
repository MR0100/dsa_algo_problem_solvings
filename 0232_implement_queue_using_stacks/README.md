# 0232 — Implement Queue using Stacks

> LeetCode #232 · Difficulty: Easy
> **Categories:** Stack, Queue, Design

---

## Problem Statement

Implement a first in first out (FIFO) queue using only two stacks. The implemented queue should support all the functions of a normal queue (`push`, `peek`, `pop`, and `empty`).

Implement the `MyQueue` class:

- `void push(int x)` Pushes element `x` to the back of the queue.
- `int pop()` Removes the element from the front of the queue and returns it.
- `int peek()` Returns the element at the front of the queue.
- `boolean empty()` Returns `true` if the queue is empty, `false` otherwise.

**Notes:**
- You must use **only** standard operations of a stack, which means only `push to top`, `peek/pop from top`, `size`, and `is empty` operations are valid.
- Depending on your language, the stack may not be supported natively. You may simulate a stack using a list or deque (double-ended queue), as long as you use only a stack's standard operations.

**Example 1:**
```
Input
["MyQueue", "push", "push", "peek", "pop", "empty"]
[[], [1], [2], [], [], []]
Output
[null, null, null, 1, 1, false]

Explanation
MyQueue myQueue = new MyQueue();
myQueue.push(1); // queue is: [1]
myQueue.push(2); // queue is: [1, 2] (leftmost is front of the queue)
myQueue.peek();  // return 1
myQueue.pop();   // return 1, queue is [2]
myQueue.empty(); // return false
```

**Constraints:**
- `1 <= x <= 9`
- At most `100` calls will be made to `push`, `pop`, `peek`, and `empty`.
- All the calls to `pop` and `peek` are valid.

**Follow-up:** Can you implement the queue such that each operation is **amortized** O(1) time complexity? In other words, performing `n` operations will take overall O(n) time even if one of those operations may take longer.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★☆ High       | 2024          |
| Microsoft  | ★★★☆☆ Medium     | 2024          |
| Google     | ★★★☆☆ Medium     | 2023          |
| Bloomberg  | ★★★☆☆ Medium     | 2023          |
| Meta       | ★★☆☆☆ Low        | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used
- **Stack** — the only primitive allowed; both approaches build a queue from LIFO stacks → see [`/dsa/stack.md`](/dsa/stack.md)
- **Queue** — the FIFO behaviour we must emulate → see [`/dsa/queue_deque.md`](/dsa/queue_deque.md)
- **Design (Data Structures)** — implementing an ADT under operation constraints → see [`/dsa/design_data_structures.md`](/dsa/design_data_structures.md)

---

## Approaches Overview
| # | Approach | push | pop / peek | empty | Space | When to use |
|---|----------|------|-----------|-------|-------|-------------|
| 1 | Two Stacks, Costly Push | O(n) | O(1) | O(1) | O(n) | When pops dominate pushes |
| 2 | Two Stacks, Amortized Push (Optimal) | O(1) | O(1)* | O(1) | O(n) | General case; answers the follow-up |

`*` amortized

---

## Approach 1 — Two Stacks, Costly Push

### Intuition
A stack reverses insertion order, but a queue must preserve it. If on every push we move everything to a helper stack, drop the new element at the bottom, then move it all back, the main stack always has the **oldest** element on top. That makes `pop`/`peek` O(1); we pay the reordering cost entirely at push time.

### Algorithm
1. **Push:** pop all of `in` into `out`; push `x` onto `in`; pop all of `out` back onto `in`. Now `in` has the front on top.
2. **Pop / Peek:** remove / read the top of `in`.
3. **Empty:** `in` is empty.

### Complexity
- **Time:** Push O(n) — every element is moved twice; Pop / Peek / Empty O(1).
- **Space:** O(n) — the two stacks together hold n elements.

### Code
```go
type CostlyPushQueue struct {
	in  []int // primary stack, maintained with the queue front on top
	out []int // scratch stack used only during Push
}

func NewCostlyPushQueue() *CostlyPushQueue {
	return &CostlyPushQueue{}
}

func (q *CostlyPushQueue) Push(x int) {
	for len(q.in) > 0 { // dump existing elements onto the scratch stack
		q.out = append(q.out, q.in[len(q.in)-1])
		q.in = q.in[:len(q.in)-1]
	}
	q.in = append(q.in, x) // new element goes in first → ends up at the bottom
	for len(q.out) > 0 {   // restore the old elements on top of it
		q.in = append(q.in, q.out[len(q.out)-1])
		q.out = q.out[:len(q.out)-1]
	}
}

func (q *CostlyPushQueue) Pop() int {
	top := q.in[len(q.in)-1] // front element sits on top by invariant
	q.in = q.in[:len(q.in)-1]
	return top
}

func (q *CostlyPushQueue) Peek() int {
	return q.in[len(q.in)-1] // front is always the top of `in`
}

func (q *CostlyPushQueue) Empty() bool {
	return len(q.in) == 0
}
```

### Dry Run
Ops `["MyQueue","push","push","peek","pop","empty"]`, args `[[],[1],[2],[],[],[]]` (top of stack listed last):

| Op        | Action | `in` (bottom→top) | Return |
|-----------|--------|-------------------|--------|
| MyQueue   | init   | []                | null   |
| push(1)   | in=[1] | [1]               | null   |
| push(2)   | move 1→out, push 2, move 1 back | [2, 1] | null |
| peek()    | top of in | [2, 1]         | **1**  |
| pop()     | remove top | [2]           | **1**  |
| empty()   | len(in)=1 | [2]            | **false** |

Output: `[null, null, null, 1, 1, false]`.

---

## Approach 2 — Two Stacks, Amortized O(1) Push (Optimal)

### Intuition
Keep pushes cheap: just drop new elements on an `in` stack. When we need the front and the `out` stack is empty, pour all of `in` into `out`. Reversing a LIFO stack into another LIFO stack yields FIFO order, so the top of `out` is the oldest element. Each element moves from `in` to `out` exactly once over its lifetime → **amortized O(1)** per operation, answering the follow-up.

### Algorithm
1. **Push:** append `x` to `in`.
2. **transfer():** if `out` is empty, pop all of `in` onto `out`.
3. **Pop:** `transfer()`; pop the top of `out`.
4. **Peek:** `transfer()`; read the top of `out`.
5. **Empty:** both stacks empty.

### Complexity
- **Time:** Push O(1); Pop / Peek O(1) amortized (each element transferred once); Empty O(1).
- **Space:** O(n).

### Code
```go
type AmortizedQueue struct {
	in  []int // newest elements, back of the queue on top
	out []int // oldest elements reversed, front of the queue on top
}

func NewAmortizedQueue() *AmortizedQueue {
	return &AmortizedQueue{}
}

func (q *AmortizedQueue) Push(x int) {
	q.in = append(q.in, x) // O(1): reordering happens later, on demand
}

func (q *AmortizedQueue) transfer() {
	if len(q.out) == 0 { // only reverse when the front supply is exhausted
		for len(q.in) > 0 {
			q.out = append(q.out, q.in[len(q.in)-1]) // reverse in → out = FIFO
			q.in = q.in[:len(q.in)-1]
		}
	}
}

func (q *AmortizedQueue) Pop() int {
	q.transfer()               // ensure the front is on top of `out`
	top := q.out[len(q.out)-1] // oldest element
	q.out = q.out[:len(q.out)-1]
	return top
}

func (q *AmortizedQueue) Peek() int {
	q.transfer() // make sure `out` holds the front
	return q.out[len(q.out)-1]
}

func (q *AmortizedQueue) Empty() bool {
	return len(q.in) == 0 && len(q.out) == 0 // no element anywhere
}
```

### Dry Run
Same ops (top of each stack listed last):

| Op        | Action | `in` | `out` | Return |
|-----------|--------|------|-------|--------|
| MyQueue   | init   | []   | []    | null   |
| push(1)   | append 1 to in | [1] | [] | null |
| push(2)   | append 2 to in | [1, 2] | [] | null |
| peek()    | out empty → transfer: reverse [1,2] into out = [2,1]; top = 1 | [] | [2, 1] | **1** |
| pop()     | out non-empty; pop top | [] | [2] | **1** |
| empty()   | in empty, out=[2] → not empty | [] | [2] | **false** |

Output: `[null, null, null, 1, 1, false]`.

---

## Key Takeaways
- Reversing a stack into another stack converts LIFO to FIFO — the whole trick behind queue-from-stacks (and stack-from-queues in reverse).
- Choose *where* to pay the cost: costly-push keeps pop O(1) worst-case; lazy-transfer makes push O(1) and pop amortized O(1). Both do O(n) total work over n ops.
- **Amortized analysis:** because each element is transferred at most once, the total transfer work across n operations is O(n), even though a single pop can be O(n).
- Only expose `push top`, `pop top`, `peek top`, `size`, `empty` — respect the primitive constraints.

## Related Problems
- LeetCode #225 — Implement Stack using Queues (the mirror image)
- LeetCode #155 — Min Stack (auxiliary-stack design)
- LeetCode #622 — Design Circular Queue (queue as a ring buffer)
- LeetCode #933 — Number of Recent Calls (queue-based sliding window)
