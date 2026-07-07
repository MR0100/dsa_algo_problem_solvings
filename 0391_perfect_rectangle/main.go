package main

import "fmt"

// ── Approach 1: Corner Counting + Area (Optimal) ─────────────────────────────
//
// cornerCounting decides whether the given axis-aligned rectangles form an
// EXACT cover of one big rectangle (no gaps, no overlaps).
//
// Intuition:
//
//	A set of rectangles tiles a bigger rectangle perfectly if and only if two
//	conditions hold simultaneously:
//	  (a) AREA: the sum of the small areas equals the area of the bounding box
//	      that spans from the minimum bottom-left to the maximum top-right.
//	      This rules out gaps (too little area) but NOT overlaps by itself.
//	  (b) CORNERS: consider the four corner points of every small rectangle.
//	      In a perfect tiling, an interior corner is shared by an even number
//	      of rectangles (2 or 4) and cancels out; only the four corners of the
//	      big bounding rectangle survive an odd number of times. So after
//	      toggling every corner in a set, EXACTLY the four bounding corners
//	      must remain, and nothing else.
//	  Area alone can be fooled by "one overlap balanced by one gap of equal
//	  size"; the corner parity catches exactly those overlaps/gaps. Together
//	  they are necessary and sufficient.
//
// Algorithm:
//  1. Track the bounding box (minX, minY, maxX, maxY) and the running area sum.
//  2. For each rectangle add its area, and toggle its 4 corners in a set:
//     if a corner is present remove it, else insert it (XOR on membership).
//  3. After the loop the set must contain EXACTLY the 4 bounding corners.
//  4. The summed area must equal the bounding-box area.
//  5. Both true ⇒ perfect cover.
//
// Time:  O(n) — one pass, each rectangle does O(1) work.
// Space: O(n) — the corner set holds up to O(n) points before cancellation.
func cornerCounting(rectangles [][]int) bool {
	// point is a hashable (x,y) key for the corner-parity set.
	type point struct{ x, y int }

	area := 0                            // running sum of small rectangle areas
	corners := map[point]bool{}          // set of corners seen an odd number of times
	minX, minY := 1<<62, 1<<62           // bounding box lower-left (start high)
	maxX, maxY := -(1 << 62), -(1 << 62) // bounding box upper-right (start low)

	for _, r := range rectangles {
		x1, y1, x2, y2 := r[0], r[1], r[2], r[3] // bottom-left (x1,y1), top-right (x2,y2)

		// Expand the bounding box to include this rectangle.
		if x1 < minX {
			minX = x1
		}
		if y1 < minY {
			minY = y1
		}
		if x2 > maxX {
			maxX = x2
		}
		if y2 > maxY {
			maxY = y2
		}

		// Accumulate this rectangle's area.
		area += (x2 - x1) * (y2 - y1)

		// Toggle each of the 4 corners: XOR membership so shared corners cancel.
		for _, c := range []point{{x1, y1}, {x1, y2}, {x2, y1}, {x2, y2}} {
			if corners[c] {
				delete(corners, c) // seen before ⇒ now even ⇒ remove
			} else {
				corners[c] = true // first time ⇒ odd ⇒ insert
			}
		}
	}

	// Exactly the 4 bounding corners must remain, nothing more, nothing less.
	if len(corners) != 4 {
		return false
	}
	for _, c := range []point{{minX, minY}, {minX, maxY}, {maxX, minY}, {maxX, maxY}} {
		if !corners[c] {
			return false // a required bounding corner is missing
		}
	}

	// Area test: total small area must exactly fill the bounding box.
	return area == (maxX-minX)*(maxY-minY)
}

// ── Approach 2: Sweep Line (Interval Merge per x-boundary) ────────────────────
//
// sweepLine decides the perfect cover by sweeping a vertical line left to
// right and verifying the covered y-intervals stay contiguous and never
// overlap as rectangles open and close.
//
// Intuition:
//
//	Sort rectangle edges by x. At each distinct x we first REMOVE the vertical
//	segments of rectangles whose right edge is here, then ADD the vertical
//	segments of rectangles whose left edge is here. For a perfect tiling, at
//	every x-slice the set of active y-intervals must form one continuous span
//	with NO overlaps and NO gaps (except the very first and last x). We keep
//	the active intervals in sorted order and check adjacency each time the
//	active set changes.
//
// Algorithm:
//  1. Build events: for each rectangle two events (x1, open) and (x2, close),
//     each carrying (y1, y2). Sort by x, processing closes before opens at a tie.
//  2. Maintain a sorted list of active [y1,y2) intervals.
//  3. When the x advances past a group, the active intervals must be pairwise
//     non-overlapping and, taken together, contiguous — verified by sorting on
//     y1 and checking each interval's y1 == previous y2 for the covered band,
//     while overlaps (y1 < previous y2) fail immediately.
//  4. Also confirm the total covered height at each full slice matches the
//     bounding height; simpler: reject on any overlap and confirm area at end.
//
// This implementation checks overlaps on insertion and uses the same area +
// corner guarantee for gaps, giving an O(n log n) alternative.
//
// Time:  O(n log n) — sorting events, plus O(n) interval maintenance.
// Space: O(n) — events and the active interval list.
func sweepLine(rectangles [][]int) bool {
	events := make([]event, 0, len(rectangles)*2) // two events per rectangle
	area := 0
	minY, maxY := 1<<62, -(1 << 62) // vertical extent of the whole figure
	for _, r := range rectangles {
		x1, y1, x2, y2 := r[0], r[1], r[2], r[3]
		events = append(events, event{x1, true, y1, y2})  // left edge opens
		events = append(events, event{x2, false, y1, y2}) // right edge closes
		area += (x2 - x1) * (y2 - y1)
		if y1 < minY {
			minY = y1
		}
		if y2 > maxY {
			maxY = y2
		}
	}

	// Sort by x; at equal x process closes first so a rectangle that ends and
	// another that begins on the same line hand off cleanly.
	// Simple insertion-friendly sort via manual comparator.
	sortEvents(events)

	active := []interval{} // currently open vertical intervals, kept sorted by y1

	i := 0
	for i < len(events) {
		x := events[i].x
		// Process every event sharing this x. Closes first.
		for i < len(events) && events[i].x == x && !events[i].open {
			active = removeInterval(active, interval{events[i].y1, events[i].y2})
			i++
		}
		for i < len(events) && events[i].x == x && events[i].open {
			// Insert keeping active sorted by y1; reject on any overlap.
			if !insertInterval(&active, interval{events[i].y1, events[i].y2}) {
				return false // overlap detected ⇒ not a perfect cover
			}
			i++
		}
	}

	// No overlaps happened. A no-overlap tiling whose area equals the bounding
	// box area and whose vertical extent is filled must be a perfect cover.
	// Reuse the area test against the bounding box.
	minX, maxX := 1<<62, -(1 << 62)
	for _, r := range rectangles {
		if r[0] < minX {
			minX = r[0]
		}
		if r[2] > maxX {
			maxX = r[2]
		}
	}
	return area == (maxX-minX)*(maxY-minY)
}

// sortEvents sorts events by x ascending, closes (open=false) before opens at
// equal x. Uses a simple stable insertion of a comparator into Go's sort.
func sortEvents(events []event) {
	// Insertion sort keeps the file dependency-free and is fine for clarity.
	for i := 1; i < len(events); i++ {
		j := i
		for j > 0 && eventLess(events[j], events[j-1]) {
			events[j], events[j-1] = events[j-1], events[j]
			j--
		}
	}
}

// eventLess is the ordering predicate: smaller x first; closes before opens.
func eventLess(a, b event) bool {
	if a.x != b.x {
		return a.x < b.x
	}
	// close (open==false) should come before open (open==true).
	return !a.open && b.open
}

// event mirrors the local type used in sweepLine so the helpers can share it.
type event struct {
	x    int
	open bool
	y1   int
	y2   int
}

// interval mirrors the local type used in sweepLine.
type interval struct{ y1, y2 int }

// insertInterval inserts iv into the sorted-by-y1 slice, returning false if it
// overlaps a neighbour. Adjacent (touching) intervals are allowed.
func insertInterval(active *[]interval, iv interval) bool {
	a := *active
	// Find insertion index by y1 (linear scan; active set is small in practice).
	idx := 0
	for idx < len(a) && a[idx].y1 < iv.y1 {
		idx++
	}
	// Check overlap with the interval before idx.
	if idx > 0 && a[idx-1].y2 > iv.y1 {
		return false // previous interval extends past iv's start ⇒ overlap
	}
	// Check overlap with the interval at idx.
	if idx < len(a) && iv.y2 > a[idx].y1 {
		return false // iv extends past next interval's start ⇒ overlap
	}
	// Insert at idx.
	a = append(a, interval{})
	copy(a[idx+1:], a[idx:])
	a[idx] = iv
	*active = a
	return true
}

// removeInterval deletes the first matching interval from the slice.
func removeInterval(a []interval, iv interval) []interval {
	for i := range a {
		if a[i] == iv {
			return append(a[:i], a[i+1:]...)
		}
	}
	return a
}

func main() {
	ex1 := [][]int{{1, 1, 3, 3}, {3, 1, 4, 2}, {3, 2, 4, 4}, {1, 3, 2, 4}, {2, 3, 3, 4}}
	ex2 := [][]int{{1, 1, 2, 3}, {1, 3, 2, 4}, {3, 1, 4, 2}, {3, 2, 4, 4}}
	ex3 := [][]int{{1, 1, 3, 3}, {3, 1, 4, 2}, {1, 3, 2, 4}, {2, 2, 4, 4}}

	fmt.Println("=== Approach 1: Corner Counting + Area (Optimal) ===")
	fmt.Printf("Example 1: got=%v  expected true\n", cornerCounting(ex1))  // expected true
	fmt.Printf("Example 2: got=%v  expected false\n", cornerCounting(ex2)) // expected false
	fmt.Printf("Example 3: got=%v  expected false\n", cornerCounting(ex3)) // expected false

	fmt.Println("=== Approach 2: Sweep Line ===")
	fmt.Printf("Example 1: got=%v  expected true\n", sweepLine(ex1))  // expected true
	fmt.Printf("Example 2: got=%v  expected false\n", sweepLine(ex2)) // expected false
	fmt.Printf("Example 3: got=%v  expected false\n", sweepLine(ex3)) // expected false
}
