# Combinatorial Game Theory

> **Scope:** two-player, **perfect-information**, deterministic games where players **alternate
> moves** and the last legal move decides the outcome. No hidden state, no chance.
> **Core question:** "With both sides playing optimally, does the player *to move* win?"
> **Three tools, in order of power:** win/lose state DP → minimax with a payoff → Sprague–Grundy / Nim.

---

## What it is

Combinatorial game theory analyses games of the form: a shared position, two players moving
in turn, a fixed rule set, and a **win condition based on who makes the last move**. Under
the **normal play convention** (the player who *cannot* move loses), every position is either:

- a **winning position (N-position)** — the player about to move (**N**ext) can force a win; or
- a **losing position (P-position)** — every move hands the opponent a winning position, so
  the **P**revious player (the one who just moved) is winning.

This gives a recursive definition you can compute directly:

> A position is **winning** iff **at least one** move leads to a **losing** position for the
> opponent. It is **losing** iff **every** move leads to a **winning** position (or there are
> no moves at all).

That one sentence is the win/lose DP. When the game instead has a **numeric score** (not just
win/lose) — e.g. "minimise the money you might lose" — you upgrade to **minimax**: the mover
picks the child that is best *for them*, assuming the opponent will reply with the child that
is best for *them*. When a game **splits into independent sub-games** played side by side
(classic Nim heaps), you use the **Sprague–Grundy theorem**: each sub-game gets a *Grundy
number*, and the whole game is winning iff the **XOR** of all Grundy numbers is non-zero.

### The three levels, and when each applies

| Level | Game shape | State per position | Combine rule |
|-------|-----------|--------------------|--------------|
| **1. Win/Lose DP** | Single game, only "who moves last" matters | boolean `canWin` | OR over moves (win if any child loses) |
| **2. Minimax** | Single game, a **numeric payoff** to optimise | best achievable score | max on my turn, min on opponent's |
| **3. Sprague–Grundy / Nim** | Game **decomposes into independent parallel sub-games** | Grundy number `g` (a non-negative int) | XOR of all sub-games' Grundy numbers |

Level 3 generalises Level 1: a position is losing (P-position) **iff its Grundy number is 0**,
so `canWin = (grundy != 0)`. You only *need* Grundy numbers when the game is a sum of
independent components; for a single indivisible game the boolean DP suffices.

### The pattern-recognition shortcut

Many contest/interview game problems have a **closed-form** answer that a small brute-force
reveals: compute win/lose for `n = 0,1,2,…`, eyeball the pattern, prove it. The archetype is
**Nim with one pile / Bash Game**: "take 1–3 stones, last to take wins" → the mover loses
**iff `n` is a multiple of 4** (#292). Always try this — a proven O(1) formula beats an O(n)
or exponential DP.

---

## When to recognise it

| Signal in the problem | Which tool |
|-----------------------|------------|
| "Two players take turns, **both play optimally**, who **wins**?" | Win/lose DP (or a pattern shortcut) |
| "Last player to **move / take / flip** wins (or loses)" | Normal/misère play → win/lose DP |
| Small move set on a single counter (take 1..k stones) | Try the **modulo pattern** first (Bash Game) |
| "Minimise the **maximum cost** you could be forced to pay" / a **score**, not just win/lose | **Minimax** interval/state DP (#375) |
| A **string/board that splits into independent regions** after a move (e.g. `"++"` → two shorter segments) | **Sprague–Grundy**: XOR the segments' Grundy numbers |
| "Several independent heaps, take from one per turn" | **Nim**: winning iff XOR of heap sizes ≠ 0 |
| "Can the **first player force a win** given a shared pool / set of choices?" (#464 Can I Win) | Win/lose DP over a **bitmask** of used choices + memo |

**When *not* to use it:** the game has hidden information or randomness (that is decision
theory / expectimax, not combinatorial game theory); or it is single-player (that is just
plain search/DP — no adversary). The defining features are *two players, alternating,
perfect information, deterministic*.

---

## General templates (Go)

### Template 1 — win/lose state DP (memoised)

The workhorse. State = whatever fully describes the position. Return "can the player to move
force a win from here?" You win if **any** move leaves the opponent in a losing state.

```go
// canWin reports whether the player to move wins from the given state.
// `moves(state)` yields every reachable next state.
func canWin(state State, memo map[State]bool) bool {
    if v, ok := memo[state]; ok {
        return v
    }
    win := false
    for _, next := range moves(state) {
        // If the opponent LOSES from `next`, then moving there wins for me.
        if !canWin(next, memo) {
            win = true
            break // one winning move is enough
        }
    }
    memo[state] = win
    return win
}
```

For **Can I Win (#464)** the state is a bitmask of already-chosen numbers (plus the running
total, which is determined by the mask, so the mask alone keys the memo):

```go
// used is a bitmask: bit i set => number (i+1) already taken.
func canWinFrom(used, remaining, maxChoosable int, memo map[int]bool) bool {
    if v, ok := memo[used]; ok {
        return v
    }
    win := false
    for i := 1; i <= maxChoosable; i++ {
        bit := 1 << uint(i)
        if used&bit != 0 {
            continue // number i already used
        }
        // Taking i wins immediately, OR forces the opponent into a losing state.
        if i >= remaining || !canWinFrom(used|bit, remaining-i, maxChoosable, memo) {
            win = true
            break
        }
    }
    memo[used] = win
    return win
}
```

### Template 2 — minimax with a numeric payoff

When the outcome is a score. The mover **maximises**; the opponent **minimises**. For the
common "the two players optimise the same quantity from opposite directions" LeetCode setup
(e.g. #375, where *you* try to minimise the worst-case cost the game imposes), it collapses to
a single min/max recurrence over sub-positions:

```go
// Guess Number Higher or Lower II (#375): pick a number 1..n; each wrong guess g
// costs g dollars and tells you higher/lower. Minimise the money you must guarantee
// to have to win no matter where the target is.
func getMoneyAmount(n int) int {
    // dp[i][j] = minimum guaranteed cost to find any target in the interval [i, j].
    dp := make([][]int, n+2)
    for i := range dp {
        dp[i] = make([]int, n+2)
    }
    // Iterate by increasing interval length so sub-intervals are ready (interval DP).
    for length := 2; length <= n; length++ {
        for i := 1; i+length-1 <= n; i++ {
            j := i + length - 1
            best := 1 << 30
            // Choose a guess k in [i, j]. Cost = k + worst of the two sides,
            // because the adversary places the target on whichever side is costlier.
            for k := i; k <= j; k++ {
                left, right := 0, 0
                if k-1 >= i {
                    left = dp[i][k-1]
                }
                if k+1 <= j {
                    right = dp[k+1][j]
                }
                cost := k + max(left, right) // pay k now, then the harder subproblem
                if cost < best {
                    best = cost
                }
            }
            dp[i][j] = best
        }
    }
    return dp[1][n]
}

func max(a, b int) int { if a > b { return a }; return b }
```

(This is where game theory meets **interval DP** — see `dsa/interval_dp.md`.)

### Template 3 — Sprague–Grundy / Nim

For a game that decomposes into independent sub-games. The **Grundy number** (a.k.a.
nimber / mex value) of a position is the **mex** — *minimum excludant*, the smallest
non-negative integer **not** among the Grundy numbers of its next positions.

```go
// grundy computes the Grundy number of a single position.
// A position with grundy 0 is losing (P-position); non-zero is winning.
func grundy(state State, memo map[State]int) int {
    if v, ok := memo[state]; ok {
        return v
    }
    reachable := map[int]bool{}
    for _, next := range moves(state) {
        reachable[grundy(next, memo)] = true
    }
    // mex: smallest non-negative integer not in `reachable`.
    g := 0
    for reachable[g] {
        g++
    }
    memo[state] = g
    return g
}

// For a game that is the SUM of independent sub-games, XOR their Grundy numbers.
// Overall position is winning iff the XOR is non-zero.
func nimSumWins(subgames []State, memo map[State]int) bool {
    x := 0
    for _, s := range subgames {
        x ^= grundy(s, memo)
    }
    return x != 0
}

// Classic Nim (heaps of stones, take any positive number from one heap):
// grundy(heap of size h) == h, so the game is just the XOR of heap sizes.
func nimWins(heaps []int) bool {
    x := 0
    for _, h := range heaps {
        x ^= h
    }
    return x != 0 // first player wins iff XOR != 0
}
```

### Template 4 — the pattern shortcut (Bash Game / #292)

Brute-force small `n`, spot the period, return O(1):

```go
// Nim Game (#292): remove 1..3 stones, last to remove wins.
// Brute force shows canWin[n] is false exactly when n % 4 == 0.
func canWinNim(n int) bool {
    return n%4 != 0
}
```

---

## Worked example — Flip Game II (#294), win/lose DP

Game: a string of `+` and `-`. A move flips **two adjacent `++` into `--`**. Players alternate;
the player who cannot move (no `++` left) loses. Does the **first** player have a forced win?

Take `s = "++++"` (four pluses). Enumerate from the start player's viewpoint —
`W` = winning for the mover, `L` = losing.

The moves from `"++++"` flip one of the `++` pairs at index 0, 1, or 2:

| Move (flip pair at) | Resulting string | Mover-to-move result of the child |
|---------------------|------------------|-----------------------------------|
| index 0 | `"--++"` | only move → `"----"` (no moves = **L** for that mover) ⇒ child is **W** |
| index 1 | `"+--+"` | no `++` left → **L** for the mover of the child |
| index 2 | `"++--"` | only move → `"----"` ⇒ child is **W** |

The first player picks the move that leaves the opponent **losing**. Flipping the **middle**
pair (index 1) yields `"+--+"`, which has no `++` — the opponent cannot move and **loses**.

So `"++++"` is a **first-player win**. ✓

Now `s = "++"`: the only move gives `"--"` (no moves for the opponent ⇒ opponent loses), so
`"++"` is also **W**. And `"+-+"` has no `++` at all → the mover immediately loses (**L**).

The DP just memoises this over string states:

```go
// canWinFlip reports whether the player to move can force a win.
// State: the current board string. Move: flip some "++" to "--".
func canWinFlip(s string, memo map[string]bool) bool {
    if v, ok := memo[s]; ok {
        return v
    }
    b := []byte(s)
    win := false
    for i := 0; i+1 < len(b); i++ {
        if b[i] == '+' && b[i+1] == '+' {
            b[i], b[i+1] = '-', '-' // make the move
            // If the opponent cannot win from the resulting board, this move wins.
            if !canWinFlip(string(b), memo) {
                win = true
            }
            b[i], b[i+1] = '+', '+' // undo (backtrack)
            if win {
                break
            }
        }
    }
    memo[s] = win
    return win
}
```

(#294 is *also* a textbook Sprague–Grundy problem: a run of `k` consecutive `+`s is an
independent sub-game, a move splits a run into two shorter runs, and the whole board's
Grundy number is the XOR of each run's Grundy number — the boolean DP above is the
simpler-to-code equivalent when you just need win/lose.)

---

## Complexity

Let `R` = number of reachable states, `M` = branching factor (moves per state).

| Tool | Time | Space | Reason |
|------|------|-------|--------|
| Win/lose DP (memoised) | **O(R · M)** | **O(R)** | Each state solved once; each tries up to M moves. |
| Can I Win (#464), bitmask state | **O(2ⁿ · n)** | **O(2ⁿ)** | `2ⁿ` subsets of choices, `n` moves each. |
| Flip Game II (#294), string state | Exponential in practice (state = string) | O(R) | Reachable boards; memoisation on the string keeps repeats down. |
| Minimax interval DP (#375) | **O(n³)** | **O(n²)** | `n²` intervals × `n` split points — see interval DP. |
| Sprague–Grundy | **O(R · M)** to fill Grundy table; **O(#subgames)** to XOR | O(R) | Same as win/lose plus the mex scan. |
| Pattern shortcut (#292) | **O(1)** | **O(1)** | Closed-form modulo test. |

The lesson: **look for the closed form first** (Nim/Bash patterns give O(1)); fall back to
memoised win/lose DP; reach for Grundy XOR only when the game genuinely splits into
independent parts; use minimax when there is a score rather than a binary outcome.

---

## Common pitfalls

1. **OR vs AND confusion in the win/lose rule.** You win if **some** move leaves the opponent
   losing (OR over children of `!canWin(child)`). You lose only if **every** move leaves the
   opponent winning. Flipping this quantifier inverts every answer.

2. **Normal play vs misère play.** Normal convention: *cannot move ⇒ you lose*. Misère:
   *the player who makes the last move loses*. They give different answers (and different
   Grundy analyses). Read which one the problem uses before writing the base case.

3. **Forgetting to memoise → exponential blow-up / TLE.** Game trees revisit states
   constantly (many move orders reach the same board). Always memoise on a key that captures
   the **full** state.

4. **Weak state key.** For Can I Win the running total is *implied* by the used-number mask,
   so key on the mask alone — but if you (wrongly) think two different masks with the same
   total are equivalent, you get wrong answers. Conversely, redundant keys waste memory.
   Include exactly the information that distinguishes positions.

5. **Grundy = XOR only for *independent* sub-games.** The Sprague–Grundy XOR rule requires
   the sub-games to be truly independent (a move touches exactly one of them). If a move can
   affect two regions at once, the decomposition is invalid.

6. **`mex` off-by-one.** The Grundy number is the *smallest non-negative integer absent* from
   the children's Grundy set — start the scan at 0, not 1, and a position with no moves has
   Grundy 0 (which correctly marks it losing under normal play).

7. **Assuming a pattern without proof.** A modulo pattern spotted for small `n` (e.g. "loses
   iff `n % 4 == 0`") should be *proven* (typically by induction: from a P-position every move
   reaches an N-position, and from an N-position some move reaches a P-position). Contest games
   sometimes break the obvious period.

8. **Minimax sign errors.** On your turn you optimise for yourself; on the opponent's turn they
   optimise for themselves (which is usually the opposite for you). In #375 the "adversary"
   is not a second player at all but the worst-case position of the hidden number — modelled as
   a `max` over the two sides. Be explicit about who maximises and who minimises at each level.

---

## Problems in this repo that use it

- [0292 — Nim Game](/0292_nim_game/README.md) — Bash-game pattern shortcut: the mover loses iff `n % 4 == 0`; O(1) after a brute-force reveals the period.
- [0294 — Flip Game II](/0294_flip_game_ii/README.md) — win/lose DP over the board string (and a Sprague–Grundy view: XOR of Grundy numbers of the independent `+`-runs).
- [0375 — Guess Number Higher or Lower II](/0375_guess_number_higher_or_lower_ii/README.md) — minimax over intervals: minimise the worst-case cost, `dp[i][j] = min over k of (k + max(dp[i][k-1], dp[k+1][j]))`.

### Related in this repo (setup / simpler variants)

- [0293 — Flip Game](/0293_flip_game/README.md) — the single-move (non-adversarial) version: just enumerate the resulting boards; the natural predecessor to #294.
- [0374 — Guess Number Higher or Lower](/0374_guess_number_higher_or_lower/README.md) — the interactive binary-search setup whose *adversarial-cost* variant becomes the #375 minimax.

### Related classics to know (not yet in repo)

- LeetCode #464 — Can I Win (bitmask win/lose DP over the shared number pool — Template 1)
- LeetCode #486 — Predict the Winner / #877 — Stone Game (minimax over a range, "score" games)
- LeetCode #1025 — Divisor Game (another modulo pattern: first player wins iff `n` is even)
- LeetCode #843-style Nim variants — Grundy numbers and nim-sum XOR
