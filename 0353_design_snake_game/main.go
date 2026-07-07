package main

import (
	"container/list"
	"fmt"
)

// Design the classic Snake game on a width x height board.
//   - The snake starts length 1 at the top-left cell (0,0).
//   - Food is a queue of [row, col] cells; only the current head food is active,
//     the next appears after the current is eaten.
//   - move(dir) returns the score, or -1 when the game is over (the snake runs
//     into a wall or into its own body).
//   - Eating food grows the snake (tail stays) and bumps the score; otherwise the
//     tail cell is vacated as the head advances.
//
// Two designs:
//   1. bruteForceSnakeGame — body as a plain slice; collision test is a linear
//      scan of the body each move.
//   2. SnakeGame (optimal) — body as a deque (for O(1) head push / tail pop) plus
//      a hash-set of occupied cells for O(1) self-collision tests.

// dirDelta maps a direction letter to its (drow, dcol) step.
var dirDelta = map[string][2]int{
	"U": {-1, 0},
	"D": {1, 0},
	"L": {0, -1},
	"R": {0, 1},
}

// ── Approach 1: Brute Force (Body Slice + Linear Collision Scan) ──────────────
//
// bruteForceSnakeGame stores the body as a slice of [row,col] cells (index 0 =
// tail, last = head). Each move computes the new head and scans the body for a
// self-collision.
//
// Intuition:
//
//	Model the snake literally as the ordered list of cells it occupies. Moving
//	means: compute the next head; if it hits a wall → game over; if it is food,
//	keep the tail (grow) and advance food; else drop the tail. Then check the new
//	head against the body for a self-hit.
//
// Algorithm:
//
//	move(dir):
//	  1. head = body[last]; compute newHead by dirDelta.
//	  2. If newHead is out of [0,h)x[0,w) → return -1.
//	  3. If newHead == current food → score++, advance food index (keep tail).
//	     else → remove body[0] (tail moves up).
//	  4. If newHead collides with any remaining body cell → return -1.
//	  5. Append newHead; return score.
//
// Time:  O(L) per move — L is the snake length (linear collision scan).
// Space: O(L) body slice.
type bruteForceSnakeGame struct {
	width, height int
	food          [][]int // remaining food cells, in order
	foodIdx       int     // index of the currently active food
	score         int
	body          [][]int // body cells; body[0] = tail, body[len-1] = head
}

// newBruteForceSnakeGame initializes the board, food queue, and the snake at (0,0).
func newBruteForceSnakeGame(width, height int, food [][]int) *bruteForceSnakeGame {
	return &bruteForceSnakeGame{
		width:  width,
		height: height,
		food:   food,
		body:   [][]int{{0, 0}}, // start length 1 at top-left
	}
}

// Move applies one direction and returns the score, or -1 on game over.
func (g *bruteForceSnakeGame) Move(direction string) int {
	d := dirDelta[direction]
	head := g.body[len(g.body)-1]
	nr, nc := head[0]+d[0], head[1]+d[1] // new head cell

	// Wall collision: outside the height x width board.
	if nr < 0 || nr >= g.height || nc < 0 || nc >= g.width {
		return -1
	}

	// Is the new head the active food?
	ate := g.foodIdx < len(g.food) && g.food[g.foodIdx][0] == nr && g.food[g.foodIdx][1] == nc
	if ate {
		g.score++
		g.foodIdx++ // consume this food; next food (if any) becomes active
	} else {
		g.body = g.body[1:] // no food: tail vacates its cell
	}

	// Self-collision: new head lands on a cell still occupied by the body.
	for _, cell := range g.body {
		if cell[0] == nr && cell[1] == nc {
			return -1
		}
	}

	g.body = append(g.body, []int{nr, nc}) // advance head
	return g.score
}

// ── Approach 2: Deque + Occupied Set (Optimal) ───────────────────────────────
//
// SnakeGame stores the body as a doubly-linked list (deque) for O(1) head push
// and tail pop, and a set of occupied cells (encoded as row*width+col) for O(1)
// self-collision checks.
//
// Intuition:
//
//	The two costs in Approach 1 are the tail removal (fine on a slice) and the
//	linear self-collision scan (the bottleneck). Replace the scan with a hash-set
//	membership test, and use a deque so both ends are O(1). Encode each cell as a
//	single integer key row*width+col.
//
// Algorithm:
//
//	move(dir):
//	  1. Compute newHead cell + its key.
//	  2. Wall check → -1.
//	  3. If not eating, first pop the tail from the deque AND the set (so the
//	     snake can move into the cell its tail is leaving — that is legal).
//	  4. If newHead key already in the set → self-collision → -1.
//	  5. Push newHead to the deque front and the set; if eating, score++ and
//	     advance food. Return score.
//
// Time:  O(1) per move (amortized) — deque ops + set ops are constant.
// Space: O(L) for the deque and the occupied set.
type SnakeGame struct {
	width, height int
	food          [][]int
	foodIdx       int
	score         int
	body          *list.List   // front = head, back = tail; values are [2]int cells
	occupied      map[int]bool // set of occupied cell keys row*width+col
}

// Constructor initializes the game with the snake at (0,0).
func Constructor(width, height int, food [][]int) SnakeGame {
	body := list.New()
	body.PushFront([2]int{0, 0}) // head == tail at start
	return SnakeGame{
		width:    width,
		height:   height,
		food:     food,
		body:     body,
		occupied: map[int]bool{0: true}, // cell (0,0) → key 0
	}
}

// key encodes a cell into a single int for set membership.
func (g *SnakeGame) key(r, c int) int { return r*g.width + c }

// Move applies one direction, returning the score or -1 on game over.
func (g *SnakeGame) Move(direction string) int {
	d := dirDelta[direction]
	head := g.body.Front().Value.([2]int)
	nr, nc := head[0]+d[0], head[1]+d[1]

	// Wall collision.
	if nr < 0 || nr >= g.height || nc < 0 || nc >= g.width {
		return -1
	}

	ate := g.foodIdx < len(g.food) && g.food[g.foodIdx][0] == nr && g.food[g.foodIdx][1] == nc

	// If not eating, the tail moves — free its cell BEFORE the self-check so the
	// head may legally step into the square the tail is leaving this turn.
	if !ate {
		tailEl := g.body.Back()
		tail := tailEl.Value.([2]int)
		g.body.Remove(tailEl)
		delete(g.occupied, g.key(tail[0], tail[1]))
	}

	// Self-collision against the (possibly tail-freed) body.
	if g.occupied[g.key(nr, nc)] {
		return -1
	}

	// Advance the head.
	g.body.PushFront([2]int{nr, nc})
	g.occupied[g.key(nr, nc)] = true
	if ate {
		g.score++
		g.foodIdx++
	}
	return g.score
}

func main() {
	// Official example: width=3, height=2, food=[[1,2],[0,1]]
	//   move R -> 0, D -> 0, R -> 1 (eat food[0]), U -> 1, L -> 2 (eat food[1]),
	//   U -> -1 (hit the top wall).
	moves := []string{"R", "D", "R", "U", "L", "U"}
	expected := []int{0, 0, 1, 1, 2, -1}

	fmt.Println("=== Approach 1: Brute Force (Body Slice) ===")
	g1 := newBruteForceSnakeGame(3, 2, [][]int{{1, 2}, {0, 1}})
	for i, m := range moves {
		fmt.Printf("move(%q) = %d\texpected %d\n", m, g1.Move(m), expected[i])
	}

	fmt.Println("=== Approach 2: Deque + Occupied Set (Optimal) ===")
	g2 := Constructor(3, 2, [][]int{{1, 2}, {0, 1}})
	for i, m := range moves {
		fmt.Printf("move(%q) = %d\texpected %d\n", m, g2.Move(m), expected[i])
	}
}
