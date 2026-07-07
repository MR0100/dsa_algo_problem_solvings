# CLAUDE.md

Project guidance for Claude Code when working in this repository.

## About this repo

A personal DSA grinding repo: solving every LeetCode problem in **Go**, with
every possible approach (brute force → optimal), detailed explanations, dry
runs, and a companion `/dsa/` reference library. Goal: complete knowledge of
every algorithm and data structure pattern that appears on LeetCode.

---

## File structure — REQUIRED for every problem

```
NNNN_problem_name_in_snake_case/
├── main.go
└── README.md
```

- Four-digit zero-padded number prefix, e.g. `0001_two_sum/`.
- Exactly **two files** per folder — nothing else.
- Always run `go run main.go` and confirm every example's output before reporting done.

---

## `main.go` rules — REQUIRED

```go
package main

import "fmt"

// ── Approach 1: Brute Force ──────────────────────────────────────────────────
//
// bruteForce solves [problem name] using [approach name].
//
// Intuition: ...
// Algorithm: ...
//
// Time:  O(...)
// Space: O(...)
func bruteForce(...) ... {
    // inline comments on every non-obvious line
}

// ── Approach 2: [Name] ───────────────────────────────────────────────────────
//
// twoPointers solves [problem name] using [approach name].
// ...
func twoPointers(...) ... { ... }

// ── Approach N: [Name] (Optimal) ─────────────────────────────────────────────
func optimalApproach(...) ... { ... }

func main() {
    fmt.Println("=== Approach 1: Brute Force ===")
    fmt.Println(bruteForce(...))   // expected

    fmt.Println("=== Approach 2: [Name] ===")
    fmt.Println(twoPointers(...))  // expected

    fmt.Println("=== Approach N: [Name] (Optimal) ===")
    fmt.Println(optimalApproach(...)) // expected
}
```

Rules:
- Function names describe the approach: `bruteForce`, `twoPointers`, `hashMap`,
  `dpBottomUp`, `dpTopDown`, `divideAndConquer`, `binarySearch`, `greedyApproach`, etc.
- Every solution function has a doc-comment block: what it does, intuition,
  algorithm steps, Time complexity, Space complexity.
- Inline comments on every non-obvious line explaining *what* and *why*.
- `main()` calls **every** solution against **every** problem example, with
  labelled section headers and the expected output as an inline comment.
- Language: **Go only**. No extra files.

---

## `README.md` structure — REQUIRED, 10 sections in order

```
# NNNN — Problem Title

> LeetCode #N · Difficulty: Easy / Medium / Hard
> **Categories:** Tag1, Tag2, Tag3

---

## Problem Statement
(verbatim from LeetCode: description, all examples with explanations, constraints, follow-up)

---

## Company Frequency
| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★★★ Very High  | 2024          |
| Amazon     | ★★★★☆ High       | 2024          |
| ...        | ...              | ...           |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used
- **Concept Name** — why it applies here → see [`/dsa/concept.md`](/dsa/concept.md)

---

## Approaches Overview
| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force | O(...) | O(...) | ... |
| 2 | [Name] | O(...) | O(...) | ... |
| N | [Name] (Optimal) | O(...) | O(...) | ... |

---

## Approach 1 — Brute Force

### Intuition
### Algorithm (step-by-step numbered list)
### Complexity
- **Time:** O(...) — reason
- **Space:** O(...) — reason
### Code
(fenced Go block with the full function)
### Dry Run
(trace through Example 1 step by step, showing variable states)

---

## Approach 2 — [Name]
(same sub-sections)

---

## Approach N — [Name] (Optimal)
(same sub-sections)

---

## Key Takeaways
- Bullet points: reusable patterns, tricks, intuitions from this problem.

---

## Related Problems
- LeetCode #N — Name (same pattern)
```

---

## `/dsa/` reference library

One Markdown file per concept at the project root `/dsa/`. Created the first
time a problem introduces that concept. Each file covers:
- What the concept is and when to recognise it
- General template / pseudocode
- A worked example
- Common pitfalls
- Links to problems in this repo that use it

---

## General conventions

- Language: **Go only** for all implementations.
- Problem numbering: four-digit zero-padded, matching LeetCode's number.
- Never skip sections in README.md.
- Never leave a solution function without Time/Space complexity in its doc comment.
- The Dry Run section must trace the first example step by step with a table or
  numbered variable-state list — not a vague prose summary.
