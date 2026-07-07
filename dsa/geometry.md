# Computational Geometry (basics)

> **The toolkit:** points, vectors, the **cross product** (orientation), areas of
> rectangles and polygons, axis-aligned overlap/union, segment intersection, and a
> handful of case-analysis tricks (self-crossing, reflection, "median beats mean").
> **The one weapon to memorise:** the 2-D cross product — it tells you *left turn,
> right turn, or straight* without a single trig call or floating-point division.

---

## What it is

Computational geometry on LeetCode is not the heavy academic field — it's a small
set of **integer, division-free** primitives applied to points on a grid. The wins
come from turning a geometric question ("do these turn left?", "do these rectangles
overlap?", "are these points collinear?") into **integer arithmetic** so you avoid
floating-point error entirely.

A **point** / **vector** is a pair `(x, y)`. Treat "point" and "vector" as the same
type; a vector from `A` to `B` is `B - A = (B.x-A.x, B.y-A.y)`.

```go
type Point struct{ X, Y int } // integer coordinates ⇒ exact arithmetic
```

The two products you build everything from:

```go
// cross(O->A, O->B): z-component of the 3-D cross product of the two vectors.
// Sign tells you the TURN direction from OA to OB; magnitude is 2× the triangle area.
func cross(O, A, B Point) int {
    return (A.X-O.X)*(B.Y-O.Y) - (A.Y-O.Y)*(B.X-O.X)
}

// dot(O->A, O->B): projection-based; >0 same-ish direction, <0 opposite, 0 perpendicular.
func dot(O, A, B Point) int {
    return (A.X-O.X)*(B.X-O.X) + (A.Y-O.Y)*(B.Y-O.Y)
}
```

**Orientation** is `sign(cross(O, A, B))`:

| `cross(O,A,B)` | meaning (going O→A→B) |
|----------------|-----------------------|
| `> 0`          | **counter-clockwise** (left turn) |
| `< 0`          | **clockwise** (right turn) |
| `= 0`          | **collinear** (A, B, O on one line) |

That single sign is the backbone of convex-hull, "max points on a line", polygon
area, segment intersection, and self-crossing detection.

---

## When to recognise it

| Signal in the problem | Geometry tool |
|-----------------------|---------------|
| "how many points lie on the same straight line" | **collinearity** via `cross == 0` (slope without division) |
| "is it a left turn / convex / which side of a line" | **orientation** = sign of the cross product |
| "area covered by rectangles / of a polygon" | **shoelace** formula, or width×height for axis-aligned rects |
| "do two rectangles overlap? total area of two rects" | **axis-aligned interval overlap** per dimension (inclusion–exclusion) |
| "do these two line segments cross" | **orientation test** on the 4 endpoint triples |
| "does this spiral path cross itself" | **case analysis** on consecutive move lengths |
| "can all points be mirrored across a vertical line" | **reflection**: pair up `x` around `(min+max)/2` |
| "minimise total travel / meeting point on a grid" | **the median** minimises the sum of absolute distances |
| "exact answer, coordinates are integers" | keep everything **integer** — cross/dot/area, never slope/`float` |

**The golden rule:** if you're about to compute a *slope* `(y2-y1)/(x2-x1)`, stop —
a division introduces float error and a divide-by-zero on vertical lines. The cross
product answers the same collinearity/turn question with integer multiplication.

---

## General templates / pseudocode

### 1. Orientation & collinearity

```go
// orientation returns +1 (CCW), -1 (CW), or 0 (collinear) for the ordered triple.
func orientation(o, a, b Point) int {
    v := cross(o, a, b)
    switch {
    case v > 0:
        return +1
    case v < 0:
        return -1
    default:
        return 0
    }
}

// collinear reports whether three points lie on a single straight line.
func collinear(a, b, c Point) bool { return cross(a, b, c) == 0 }
```

### 2. Max points on a line (fix an anchor, group by direction)

For each anchor point, every other point defines a *direction*. Points sharing a
**reduced** direction vector are collinear with the anchor. Reduce `(dx, dy)` by its
gcd (and canonicalise the sign) so parallel directions hash to the same key — no
floats, no slope, verticals handled for free.

```go
func maxPoints(points []Point) int {
    n := len(points)
    if n <= 2 {
        return n
    }
    best := 0
    for i := 0; i < n; i++ {
        // direction key -> how many points share it, relative to anchor i
        dirs := map[[2]int]int{}
        for j := 0; j < n; j++ {
            if j == i {
                continue
            }
            dx := points[j].X - points[i].X
            dy := points[j].Y - points[i].Y
            g := gcd(abs(dx), abs(dy))
            if g != 0 {
                dx /= g
                dy /= g
            }
            // canonicalise sign so (1,2) and (-1,-2) map together
            if dx < 0 || (dx == 0 && dy < 0) {
                dx, dy = -dx, -dy
            }
            dirs[[2]int{dx, dy}]++
        }
        for _, c := range dirs {
            if c+1 > best { // +1 to count the anchor itself
                best = c + 1
            }
        }
    }
    return best
}

func gcd(a, b int) int { for b != 0 { a, b = b, a%b }; return a }
func abs(x int) int    { if x < 0 { return -x }; return x }
```

### 3. Polygon area — the shoelace formula

```go
// shoelaceArea returns the area of a simple polygon given its vertices in order.
// Sum of cross products of consecutive edge vectors, halved and absolute-valued.
func shoelaceArea(p []Point) float64 {
    n := len(p)
    sum := 0
    for i := 0; i < n; i++ {
        j := (i + 1) % n
        sum += p[i].X*p[j].Y - p[j].X*p[i].Y
    }
    if sum < 0 {
        sum = -sum
    }
    return float64(sum) / 2.0
}
```

### 4. Axis-aligned rectangle area & overlap (inclusion–exclusion)

A rectangle is two intervals: `[ax1, ax2]` on x and `[ay1, ay2]` on y. Two
axis-aligned rectangles overlap **iff** their x-intervals overlap **and** their
y-intervals overlap. The overlap is itself a rectangle whose sides are the
per-dimension interval intersections.

```go
// totalCover returns the area covered by two axis-aligned rectangles (union),
// A = (ax1,ay1,ax2,ay2), B = (bx1,by1,bx2,by2), lower-left / upper-right corners.
func totalCover(ax1, ay1, ax2, ay2, bx1, by1, bx2, by2 int) int {
    areaA := (ax2 - ax1) * (ay2 - ay1)
    areaB := (bx2 - bx1) * (by2 - by1)

    // overlap width = intersection of the two x-intervals (0 if disjoint)
    ow := min(ax2, bx2) - max(ax1, bx1)
    oh := min(ay2, by2) - max(ay1, by1)
    overlap := 0
    if ow > 0 && oh > 0 {
        overlap = ow * oh
    }
    return areaA + areaB - overlap // inclusion–exclusion: don't double-count the middle
}
```

### 5. Segment intersection (orientation test)

Two segments `p1p2` and `p3p4` **properly** cross when the two endpoints of each
segment lie on opposite sides of the other segment — i.e. the orientations differ:

```go
func segmentsIntersect(p1, p2, p3, p4 Point) bool {
    d1 := orientation(p3, p4, p1)
    d2 := orientation(p3, p4, p2)
    d3 := orientation(p1, p2, p3)
    d4 := orientation(p1, p2, p4)
    // general case: each segment straddles the line of the other
    if d1 != d2 && d3 != d4 {
        return true
    }
    // (collinear-touching cases would add on-segment checks here)
    return false
}
```

### 6. "Median minimises sum of absolute distances" (Best Meeting Point)

On a line, the point minimising `Σ |x_i - p|` is the **median** of the `x_i`
(any value between the two middle points when the count is even). In 2-D with
**Manhattan** distance the x and y axes are independent, so solve each 1-D problem
separately and add. Proof sketch: below the median, moving right decreases more
distances than it increases; the balance flips exactly at the median.

```go
// minTotal1D = sum of distances from every collected coordinate to their median.
// Sort, then pair the smallest with the largest: their span is forced regardless
// of where the meeting point lands between them.
func minTotal1D(xs []int) int {
    sort.Ints(xs)
    total, i, j := 0, 0, len(xs)-1
    for i < j {
        total += xs[j] - xs[i] // outermost pair contributes its full span
        i, j = i+1, j-1
    }
    return total
}
```

### 7. Self-crossing — case analysis, not simulation

For a spiral of moves `d[0], d[1], …` (N, W, S, E, repeating), a crossing can only
happen against edges 2, 3, or 4 steps back. Compare the current move length to
earlier ones — three geometric cases cover every crossing:

```go
func isSelfCrossing(d []int) bool {
    for i := 3; i < len(d); i++ {
        // Case 1: current edge crosses the edge 3 back (│ shape closes)
        if d[i] >= d[i-2] && d[i-1] <= d[i-3] {
            return true
        }
        // Case 2: current edge touches the edge 4 back
        if i >= 4 && d[i-1] == d[i-3] && d[i]+d[i-4] >= d[i-2] {
            return true
        }
        // Case 3: current edge crosses the edge 5 back (inward-then-out spiral)
        if i >= 5 && d[i-2] >= d[i-4] && d[i] >= d[i-2]-d[i-4] &&
            d[i-1] <= d[i-3] && d[i-1] >= d[i-3]-d[i-5] {
            return true
        }
    }
    return false
}
```

---

## Worked example — orientation & collinearity

Points `A=(1,1)`, `B=(2,2)`, and two candidates. Is each collinear with A→B?

`cross(A, B, C) = (B.X-A.X)*(C.Y-A.Y) - (B.Y-A.Y)*(C.X-A.X) = 1*(C.Y-1) - 1*(C.X-1)`.

| C        | `1*(C.Y-1) - 1*(C.X-1)` | value | verdict |
|----------|--------------------------|-------|---------|
| `(3,3)`  | `1*2 - 1*2`              | `0`   | collinear (on the line y=x) |
| `(4,4)`  | `1*3 - 1*3`              | `0`   | collinear |
| `(2,3)`  | `1*2 - 1*1`              | `+1`  | left turn (above the line) |
| `(3,2)`  | `1*1 - 1*2`              | `-1`  | right turn (below the line) |

So `A, B, (3,3), (4,4)` all lie on one line → 4 points; the sign of `cross`
immediately classifies the off-line points without ever computing the slope `1`.

### Worked example — rectangle union (#223 shape)

`A = (0,0,2,2)` (area 4), `B = (1,1,3,3)` (area 4).

- overlap width `ow = min(2,3) - max(0,1) = 2 - 1 = 1`
- overlap height `oh = min(2,3) - max(0,1) = 2 - 1 = 1`
- overlap `= 1*1 = 1`
- union `= 4 + 4 - 1 = 7`.

The subtracted `1` is the shared centre square `[1,2]×[1,2]`; inclusion–exclusion
stops it being counted in both rectangles.

---

## Complexity

| Primitive | Time | Space | Note |
|-----------|------|-------|------|
| cross / dot / orientation | O(1) | O(1) | pure integer arithmetic |
| collinear / segment-intersect | O(1) | O(1) | a few cross products |
| rectangle overlap / union (two rects) | O(1) | O(1) | per-dimension interval math |
| shoelace polygon area | O(n) | O(1) | one pass over vertices |
| max points on a line | O(n²) | O(n) | anchor × direction-map per anchor |
| best meeting point (median) | O(n log n) or O(n) | O(n) | sort, or median via quickselect |
| self-crossing | O(n) | O(1) | one pass, constant look-back |

---

## Common pitfalls

1. **Using slope instead of the cross product.** `(y2-y1)/(x2-x1)` divides by zero on vertical lines and loses precision as `float`. The integer cross product answers every collinearity/turn question exactly — reach for it first.

2. **Not reducing direction vectors.** In "points on a line", `(2,4)` and `(1,2)` are the *same* direction but hash differently unless you divide by `gcd` **and** fix the sign (e.g. force `dx > 0`, or `dx==0 ⇒ dy>0`). Forgetting sign-canonicalisation splits one line into two.

3. **Overflow in cross/area.** `(A.X-O.X)*(B.Y-O.Y)` multiplies two coordinate *differences*; with large coordinates this overflows 32-bit. In Go, `int` is 64-bit on modern platforms — fine — but be deliberate if you narrow to `int32`.

4. **Inclusive vs. exclusive rectangle edges.** Decide whether rectangles are `[x1,x2)` or `[x1,x2]`. For **area** it rarely matters; for **coverage / perfect-rectangle** tiling, touching edges must *not* count as overlap — use strict `>` for the overlap test (`ow > 0`), not `>=`.

5. **Zero/negative overlap treated as positive.** If the intervals are disjoint, `min(hi) - max(lo)` is negative. Clamp to 0 before multiplying, or you subtract a bogus "overlap" and get too large a union.

6. **Even-count median ambiguity.** When there's an even number of coordinates, *any* point between the two central values is optimal for the sum of absolute distances — don't over-think picking one; pairing outermost coordinates (`xs[j]-xs[i]`) sidesteps the choice entirely.

7. **Assuming polygon vertices are ordered / simple.** The shoelace formula needs vertices in sequence (CW or CCW) around a **non-self-intersecting** polygon. A shuffled vertex list gives garbage.

8. **Self-crossing off-by-one.** The look-back indices (`d[i-2]`, `d[i-3]`, …) must be guarded (`i >= 4`, `i >= 5`) or you index out of range on the first few moves.

---

## Problems in this repo that use it

- [0149 — Max Points on a Line](/0149_max_points_on_a_line/README.md) — collinearity via cross product / gcd-reduced direction map (template 2)
- [0223 — Rectangle Area](/0223_rectangle_area/README.md) — axis-aligned union by inclusion–exclusion (template 4, traced above)
- [0296 — Best Meeting Point](/0296_best_meeting_point/README.md) — median minimises sum of Manhattan distances (template 6)
- [0335 — Self Crossing](/0335_self_crossing/README.md) — spiral self-intersection case analysis (template 7)
- [0356 — Line Reflection](/0356_line_reflection/README.md) — reflect points across the vertical axis `x = (min+max)/2`; every point needs its mirror
- [0391 — Perfect Rectangle](/0391_perfect_rectangle/README.md) — exact tiling: total sub-area must equal the bounding rectangle **and** only the 4 outer corners appear an odd number of times

### Related and adjacent techniques

- [`/dsa/line_sweep.md`](/dsa/line_sweep.md) — many rectangle/interval geometry problems (skyline, perfect rectangle) are solved by sweeping a line across the plane
- [`/dsa/intervals.md`](/dsa/intervals.md) — the 1-D projection of axis-aligned rectangle overlap
- [`/dsa/math_number_theory.md`](/dsa/math_number_theory.md) — gcd (used to reduce direction vectors) and integer arithmetic
- [`/dsa/sorting.md`](/dsa/sorting.md) — ordering points/coordinates is the first step of most geometry sweeps

### Classics worth knowing (not necessarily in repo yet)

- LeetCode #587 — Erect the Fence (convex hull via orientation / Andrew's monotone chain)
- LeetCode #963 / #939 — Minimum Area Rectangle (hash corners, check axis-aligned pairs)
- LeetCode #1266 — Minimum Time Visiting All Points (Chebyshev distance = `max(|dx|,|dy|)`)
