# 0225 — Implement Stack using Queues

> LeetCode #225 · Difficulty: Easy
> **Categories:** Stack, Queue, Design

---

## Problem Statement

Implement a last-in-first-out (LIFO) stack using only two queues. The implemented stack should support all the functions of a normal stack (`push`, `top`, `pop`, and `empty`).

Implement the `MyStack` class:

- `void push(int x)` Pushes element `x` to the top of the stack.
- `int pop()` Removes the element on the top of the stack and returns it.
- `int top()` Returns the element on the top of the stack.
- `boolean empty()` Returns `true` if the stack is empty, `false` otherwise.

**Notes:**

- You must use **only** standard operations of a queue, which means that only `push to back`, `peek/pop from front`, `size` and `is empty` operations are valid.
- Depending on your language, the queue may not be supported natively. You may simulate a queue using a list or deque (double-ended queue) as long as you use only a queue's standard operations.

**Example 1:**

```
Input
["MyStack", "push", "push", "top", "pop", "empty"]
[[], [1], [2], [], [], []]
Output
[null, null, null, 2, 2, false]

Explanation
MyStack myStack = new MyStack();
myStack.push(1);
myStack.push(2);
myStack.top();   // return 2
myStack.pop();   // return 2
myStack.empty(); // return False
```

**Constraints:**

- `1 <= x <= 9`
- At most `100` calls will be made to `push`, `pop`, `top`, and `empty`.
- All the calls to `pop` and `top` are valid.

**Follow-up:** Can you implement the stack using only one queue?

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★☆☆ Medium     | 2024          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Google     | ★★☆☆☆ Low        | 2023          |
| Bloomberg  | ★★☆☆☆ Low        | 2022          |
| Apple      | ★☆☆☆☆ Rare       | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Stack (LIFO) & Queue (FIFO)** — the exercise is to synthesize LIFO behaviour from FIFO-only primitives → see [`/dsa/stack.md`](/dsa/stack.md) and [`/dsa/queue_deque.md`](/dsa/queue_deque.md)
- **Data-Structure Design** — expose a clean class API (`Push/Pop/Top/Empty`) over an internal representation → see [`/dsa/design_data_structures.md`](/dsa/design_data_structures.md)

---

## Approaches Overview

| # | Approach | Push | Pop / Top / Empty | Space | When to use |
|---|----------|------|-------------------|-------|-------------|
| 1 | Two Queues, Costly Push | O(n) | O(1) | O(n) | Straightforward; classic "two queue" answer |
| 2 | One Queue, Costly Push (Optimal) | O(n) | O(1) | O(n) | Answers the follow-up; single container, less code |

> A "costly pop" variant also exists (Push O(1), Pop O(n)); the two shown make Pop/Top O(1), which is the more common ask.

---

## Approach 1 — Two Queues, Costly Push

### Intuition
A queue only lets you remove from the **front**, but a stack must remove the **newest** element. So arrange the newest element to always sit at the front of the main queue `q1`. On each Push: put the new value into the empty scratch queue `q2` first, then drain everything from `q1` behind it; swap the two queues. Now `q1`'s front is the just-pushed element, and older elements trail in the correct LIFO order.

### Algorithm
1. `Push(x)`: `q2.enqueue(x)`; while `q1` non-empty, move its front to `q2`; swap `q1` and `q2`.
2. `Pop()`: `q1.dequeue()` (front = newest).
3. `Top()`: `q1.front()`.
4. `Empty()`: `q1.size() == 0`.

### Complexity
- **Time:** Push O(n) — every existing element is moved once; Pop/Top/Empty O(1).
- **Space:** O(n) — the `n` elements, transiently spread across two queues.

### Code
```go
type TwoQueueStack struct {
	q1 *queue // holds elements with newest at the front
	q2 *queue // scratch queue used during Push
}

func NewTwoQueueStack() *TwoQueueStack {
	return &TwoQueueStack{q1: &queue{}, q2: &queue{}}
}

func (s *TwoQueueStack) Push(x int) {
	s.q2.enqueue(x) // new element goes in first, so it ends up at the front
	for s.q1.size() > 0 {
		s.q2.enqueue(s.q1.dequeue()) // older elements queued behind it, order preserved
	}
	s.q1, s.q2 = s.q2, s.q1 // q2 now has the desired order; make it the main queue
}

func (s *TwoQueueStack) Pop() int    { return s.q1.dequeue() }
func (s *TwoQueueStack) Top() int    { return s.q1.front() }
func (s *TwoQueueStack) Empty() bool { return s.q1.size() == 0 }
```

### Dry Run
Example 1 operations `push(1), push(2), top, pop, empty`. Front of `q1` shown leftmost.

| Op | q2 build | after swap q1 | Returns |
|----|----------|---------------|---------|
| push(1) | q2=[1]; drain q1 (empty) | q1=[1] | null |
| push(2) | q2=[2]; drain q1=[1] → q2=[2,1] | q1=[2,1] | null |
| top() | — | q1=[2,1] | **2** (front) |
| pop() | — | q1=[1] | **2** |
| empty() | — | q1=[1] | **false** |

Output `[null, null, null, 2, 2, false]`. ✔

---

## Approach 2 — One Queue, Costly Push (Optimal)

### Intuition
Two queues aren't necessary. Enqueue `x` at the back of the single queue, then **rotate**: dequeue and re-enqueue the `size − 1` elements that were already there. They circle around behind `x`, so `x` ends up at the front. The invariant "front = current top" is thus maintained after every Push, making Pop/Top trivial front operations.

### Algorithm
1. `Push(x)`: `q.enqueue(x)`; repeat `size − 1` times: `q.enqueue(q.dequeue())`.
2. `Pop()`: `q.dequeue()`.
3. `Top()`: `q.front()`.
4. `Empty()`: `q.size() == 0`.

### Complexity
- **Time:** Push O(n) — rotates the prior `n − 1` elements; Pop/Top/Empty O(1).
- **Space:** O(n) — a single queue holding all elements.

### Code
```go
type OneQueueStack struct {
	q *queue // single queue; front is always the current stack top
}

func NewOneQueueStack() *OneQueueStack {
	return &OneQueueStack{q: &queue{}}
}

func (s *OneQueueStack) Push(x int) {
	s.q.enqueue(x)                      // x goes to the back...
	for i := 0; i < s.q.size()-1; i++ { // ...rotate the OLD elements around behind it
		s.q.enqueue(s.q.dequeue())
	}
}

func (s *OneQueueStack) Pop() int    { return s.q.dequeue() }
func (s *OneQueueStack) Top() int    { return s.q.front() }
func (s *OneQueueStack) Empty() bool { return s.q.size() == 0 }
```

### Dry Run
Example 1. Front leftmost; `size − 1` rotations after each enqueue.

| Op | enqueue | rotations | q (front→back) | Returns |
|----|---------|-----------|----------------|---------|
| push(1) | [1] | size−1 = 0 | [1] | null |
| push(2) | [1,2] | size−1 = 1: dequeue 1, enqueue 1 → [2,1] | [2,1] | null |
| top() | — | — | [2,1] | **2** |
| pop() | — | — | [1] | **2** |
| empty() | — | — | [1] | **false** |

Output `[null, null, null, 2, 2, false]`. ✔

---

## Key Takeaways
- To turn a FIFO queue into a LIFO stack, maintain the invariant **"the newest element is at the front"**; the cost lives entirely in Push (a rotation), keeping Pop/Top O(1).
- The single-queue rotation (`enqueue x`, then rotate the old `size−1` elements around) is the elegant answer to the follow-up and generalizes the two-queue swap.
- This is a design/invariant problem: pick which operation eats the O(n) cost based on the expected op mix (costly-push vs costly-pop).

---

## Related Problems
- LeetCode #232 — Implement Queue using Stacks (the dual problem)
- LeetCode #155 — Min Stack (stack design with an extra invariant)
- LeetCode #622 — Design Circular Queue (queue design)
- LeetCode #146 — LRU Cache (data-structure design with combined structures)
