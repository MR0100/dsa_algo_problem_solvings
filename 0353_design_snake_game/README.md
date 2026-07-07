# 0353 — Design Snake Game

> LeetCode #353 · Difficulty: Medium
> **Categories:** Design, Queue, Deque, Hash Set, Simulation

---

## Problem Statement

Design a Snake game that is played on a device with screen size
`height x width`. Play the game online if you are not familiar with the game.

The snake is initially positioned at the top left corner `(0, 0)` with a length
of `1` unit.

You are given an array `food` where `food[i] = (ri, ci)` is the row and column
position of a piece of food that the snake can eat. When a snake eats a piece of
food, its length and the game's score both increase by `1`.

Each piece of food appears one by one on the screen, meaning the second piece of
food will not appear until the snake eats the first piece of food.

When a piece of food appears on the screen, it is **guaranteed** that it will not
appear on a block occupied by the snake.

The game is over if the snake goes out of bounds (hits a wall) or if its head
occupies a space that its body occupies **after** moving (i.e. a snake of length
4 cannot run into itself).

Implement the `SnakeGame` class:

- `SnakeGame(int width, int height, int[][] food)` Initializes the object with a
  screen of size `height x width` and the positions of the `food`.
- `int move(String direction)` Returns the score of the game after applying one
  `direction` move by the snake. If the game is over, return `-1`.

`direction` is one of `"U"`, `"D"`, `"L"`, or `"R"` (up, down, left, right).

**Example 1:**

```
Input
["SnakeGame", "move", "move", "move", "move", "move", "move"]
[[3, 2, [[1, 2], [0, 1]]], ["R"], ["D"], ["R"], ["U"], ["L"], ["U"]]
Output
[null, 0, 0, 1, 1, 2, -1]

Explanation
SnakeGame snakeGame = new SnakeGame(3, 2, [[1, 2], [0, 1]]);
snakeGame.move("R"); // return 0
snakeGame.move("D"); // return 0
snakeGame.move("R"); // return 1, snake eats the first piece of food. The second
                     //           piece of food appears at (0, 1).
snakeGame.move("U"); // return 1
snakeGame.move("L"); // return 2, snake eats the second food. No more food appears.
snakeGame.move("U"); // return -1, game over because snake collides with border.
```

**Constraints:**

- `1 <= width, height <= 10^4`
- `1 <= food.length <= 50`
- `food[i].length == 2`
- `0 <= ri < height`
- `0 <= ci < width`
- `direction.length == 1`
- `direction` is `'U'`, `'D'`, `'L'`, or `'R'`.
- At most `10^4` calls will be made to `move`.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★☆ High       | 2023          |
| Google     | ★★★☆☆ Medium     | 2023          |
| Bloomberg  | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Deque (double-ended queue)** — the snake's body pushes at the head and pops
  at the tail every move → see [`/dsa/queue_deque.md`](/dsa/queue_deque.md)
- **Hash set for O(1) membership** — occupied cells stored as encoded keys so
  self-collision is a constant-time lookup → see [`/dsa/hash_map.md`](/dsa/hash_map.md)
- **Stateful simulation / design** — maintain body, score, and food pointer
  across a sequence of moves → see [`/dsa/design_data_structures.md`](/dsa/design_data_structures.md)

---

## Approaches Overview

| # | Approach | move() | Space | When to use |
|---|----------|--------|-------|-------------|
| 1 | Brute Force (body slice + linear scan) | O(L) | O(L) | Short snakes; clearest model |
| 2 | Deque + occupied set (Optimal) | O(1) amortized | O(L) | Many moves / long snakes |

`L` = current snake length.

---

## Approach 1 — Brute Force (Body Slice + Linear Collision Scan)

### Intuition
Model the snake literally as an ordered slice of `[row, col]` cells, tail first,
head last. Each move computes the new head, checks the wall, decides eat-vs-move
(eating keeps the tail so the body grows; otherwise the tail is dropped), then
linearly scans the remaining body for a self-hit.

### Algorithm
1. `head = body[last]`; compute `newHead` via the direction delta.
2. If `newHead` is out of the `height x width` bounds → return `-1`.
3. If `newHead` equals the active food cell → `score++`, advance the food index
   (keep the tail — the snake grew). Else drop `body[0]` (tail vacates).
4. Scan the remaining body: if any cell equals `newHead` → return `-1`.
5. Append `newHead`; return `score`.

### Complexity
- **Time:** O(L) per move — the self-collision scan walks the whole body.
- **Space:** O(L) for the body slice.

### Code
```go
func (g *bruteForceSnakeGame) Move(direction string) int {
	d := dirDelta[direction]
	head := g.body[len(g.body)-1]
	nr, nc := head[0]+d[0], head[1]+d[1] // new head cell

	if nr < 0 || nr >= g.height || nc < 0 || nc >= g.width {
		return -1
	}

	ate := g.foodIdx < len(g.food) && g.food[g.foodIdx][0] == nr && g.food[g.foodIdx][1] == nc
	if ate {
		g.score++
		g.foodIdx++
	} else {
		g.body = g.body[1:] // no food: tail vacates its cell
	}

	for _, cell := range g.body {
		if cell[0] == nr && cell[1] == nc {
			return -1
		}
	}

	g.body = append(g.body, []int{nr, nc}) // advance head
	return g.score
}
```

### Dry Run
Board `width=3, height=2`, food `[[1,2],[0,1]]`. Body starts `[[0,0]]`.

| move | newHead | wall? | ate? | body after tail step | self-hit? | body/head | score |
|------|---------|-------|------|----------------------|-----------|-----------|-------|
| R | (0,1) | no | no | drop tail → [] then... | no | [[0,1]] | 0 |
| D | (1,1) | no | no | [] | no | [[1,1]] | 0 |
| R | (1,2) | no | **yes** food[0] | keep tail [[1,1]] | no | [[1,1],[1,2]] | 1 |
| U | (0,2) | no | no | drop [1,1] → [[1,2]] | no | [[1,2],[0,2]] | 1 |
| L | (0,1) | no | **yes** food[1] | keep [[1,2],[0,2]] | no | [...,[0,1]] | 2 |
| U | (-1,1) | **yes** | — | — | — | — | -1 |

Output sequence `0, 0, 1, 1, 2, -1` — matches the expected answer.

---

## Approach 2 — Deque + Occupied Set (Optimal)

### Intuition
The bottleneck in Approach 1 is the linear self-collision scan. Replace it with a
hash set of occupied cells (each cell encoded as `row*width + col`) for O(1)
membership. Store the body as a deque so head-push and tail-pop are both O(1).
Crucially, when not eating, free the tail cell from the set **before** the
self-check — the head is allowed to move into the square the tail is vacating.

### Algorithm
1. Compute `newHead` and its key `nr*width + nc`.
2. Wall check → `-1`.
3. If not eating, pop the tail from the deque and `delete` its key from the set.
4. If the set still contains the new head key → self-collision → `-1`.
5. Push `newHead` to the deque front and add its key to the set. If eating,
   `score++` and advance the food index. Return `score`.

### Complexity
- **Time:** O(1) amortized per move — deque ends and set ops are constant.
- **Space:** O(L) for the deque and the occupied set.

### Code
```go
func (g *SnakeGame) Move(direction string) int {
	d := dirDelta[direction]
	head := g.body.Front().Value.([2]int)
	nr, nc := head[0]+d[0], head[1]+d[1]

	if nr < 0 || nr >= g.height || nc < 0 || nc >= g.width {
		return -1
	}

	ate := g.foodIdx < len(g.food) && g.food[g.foodIdx][0] == nr && g.food[g.foodIdx][1] == nc

	// Free the tail BEFORE the self-check so the head can step into it.
	if !ate {
		tailEl := g.body.Back()
		tail := tailEl.Value.([2]int)
		g.body.Remove(tailEl)
		delete(g.occupied, g.key(tail[0], tail[1]))
	}

	if g.occupied[g.key(nr, nc)] {
		return -1
	}

	g.body.PushFront([2]int{nr, nc})
	g.occupied[g.key(nr, nc)] = true
	if ate {
		g.score++
		g.foodIdx++
	}
	return g.score
}
```

### Dry Run
Same board `3x2`, food `[[1,2],[0,1]]`. Keys use `r*3 + c`. Deque front = head.

| move | newHead | key | wall? | ate? | tail freed | in set? | deque (front→back) | score |
|------|---------|-----|-------|------|-----------|---------|--------------------|-------|
| R | (0,1) | 1 | no | no | (0,0)k0 | no | [(0,1)] | 0 |
| D | (1,1) | 4 | no | no | (0,1)k1 | no | [(1,1)] | 0 |
| R | (1,2) | 5 | no | **yes** | none | no | [(1,2),(1,1)] | 1 |
| U | (0,2) | 2 | no | no | (1,1)k4 | no | [(0,2),(1,2)] | 1 |
| L | (0,1) | 1 | no | **yes** | none | no | [(0,1),(0,2),(1,2)] | 2 |
| U | (-1,1) | — | **yes** | — | — | — | — | -1 |

Output `0, 0, 1, 1, 2, -1` — matches.

---

## Key Takeaways

- **Free the tail before the self-collision check.** A subtle but critical rule:
  the cell the tail leaves this turn is legal for the head to enter, so the set
  must be updated first (only when not eating).
- **Encode 2-D cells as one integer** (`row*width + col`) to use a plain
  `map[int]bool` as an O(1) occupancy set.
- A deque (Go's `container/list`) gives O(1) push-head / pop-tail — the exact two
  operations a moving snake needs.
- Model food as a queue with a single pointer; only the head food is active, so
  no separate spawn logic is needed.

---

## Related Problems

- LeetCode #146 — LRU Cache (deque + hash map design)
- LeetCode #362 — Design Hit Counter (sliding queue design)
- LeetCode #359 — Logger Rate Limiter (stateful map design)
- LeetCode #348 — Design Tic-Tac-Toe (grid game simulation)
