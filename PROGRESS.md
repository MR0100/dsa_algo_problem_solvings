# Progress Tracker

Tracks completion status for every LeetCode problem in this repo, per
[CLAUDE.md](CLAUDE.md)'s required structure (exactly `main.go` + `README.md`,
every approach brute-force → optimal, 10-section README, `go run` verified).

Legend: ✅ Complete &nbsp;·&nbsp; ⚠️ Partial (code exists, README missing/incomplete) &nbsp;·&nbsp; ❌ Missing entirely

Last audited: 2026-07-02 (after a prior session was killed mid-batch; this file
reflects a from-scratch re-verification of every folder in #131–210 —
`go run` executed and compared against expected values, all 10 README
sections grep-checked, `/dsa/` links resolved).

---

## #0001 – #0130 — ✅ Complete (manually verified by user)

All 130 folders confirmed present with exactly `main.go` + `README.md`.
Marked complete on the user's word — not re-verified by this audit.

> ⚠️ **Known issue, not yet fixed:** 11 of these READMEs link to `/dsa/`
> files that don't exist under the *current* naming convention established
> later in the project (`dynamic_programming.md` should be
> `dynamic_programming_1d.md` / `_2d.md`; `arrays.md`, `math.md`, `heap.md`
> were never created). Affected: #0001, #0002, #0010, #0023, #0087, #0091,
> #0096, #0097, #0115, #0120, #0123. Left untouched pending your decision
> since you already signed off on this range — see summary message.

| # | Title |
|---|---|
| 0001–0130 | (see repo — all present, all 2-file structure confirmed) |

---

## #0131 – #0210 — Batch re-verification results

72 of 80 problems are ✅ complete and pass every check: `go run main.go`
compiles, every printed result matches its expected-value comment, 2–6
approaches per problem (brute force → optimal) each with full
Intuition/Algorithm/Time/Space doc comments, all 10 README sections present
in order with genuine step-by-step Dry Run tables, and every `/dsa/*.md`
link resolves. 8 problems have gaps (7 missing README, 1 missing entirely).

| # | Title | Status | Notes |
|---|---|---|---|
| 0131 | Palindrome Partitioning | ✅ | 3 approaches, verified |
| 0132 | Palindrome Partitioning II | ✅ | 4 approaches, verified |
| 0133 | Clone Graph | ✅ | 3 approaches, verified |
| 0134 | Gas Station | ✅ | 3 approaches, verified |
| 0135 | Candy | ✅ | 3 approaches, verified |
| 0136 | Single Number | ✅ | 5 approaches, verified |
| 0137 | Single Number II | ✅ | 5 approaches, verified |
| 0138 | Copy List with Random Pointer | ✅ | 3 approaches, verified |
| 0139 | Word Break | ✅ | 5 approaches, verified |
| 0140 | Word Break II | ✅ | 3 approaches, verified |
| 0141 | Linked List Cycle | ✅ | 4 approaches, verified |
| 0142 | Linked List Cycle II | ✅ | 4 approaches, verified |
| 0143 | Reorder List | ✅ | 3 approaches, verified |
| 0144 | Binary Tree Preorder Traversal | ✅ | 3 approaches, verified |
| 0145 | Binary Tree Postorder Traversal | ✅ | 4 approaches, verified |
| 0146 | LRU Cache | ✅ | 3 approaches, verified |
| 0147 | Insertion Sort List | ✅ | 3 approaches, verified |
| 0148 | Sort List | ✅ | 3 approaches, verified |
| 0149 | Max Points on a Line | ✅ | 3 approaches, verified |
| 0150 | Evaluate Reverse Polish Notation | ✅ | 3 approaches, verified |
| 0151 | Reverse Words in a String | ✅ | 3 approaches, verified |
| 0152 | Maximum Product Subarray | ✅ | 3 approaches, verified |
| 0153 | Find Minimum in Rotated Sorted Array | ✅ | 3 approaches, verified |
| 0154 | Find Minimum in Rotated Sorted Array II | ✅ | 3 approaches, verified |
| 0155 | Min Stack | ✅ | 4 approaches, verified |
| 0156 | Binary Tree Upside Down | ✅ | 3 approaches; independently deep-audited, fully compliant |
| 0157 | Read N Characters Given Read4 | ✅ | 2 approaches; independently deep-audited, fully compliant |
| 0158 | Read N Characters Given Read4 II | ✅ | 2 approaches; independently deep-audited, fully compliant |
| 0159 | Longest Substring with At Most Two Distinct Characters | ✅ | 3 approaches; independently deep-audited, fully compliant |
| 0160 | Intersection of Two Linked Lists | ✅ | 4 approaches; independently deep-audited, fully compliant |
| 0161 | One Edit Distance | ✅ | 3 approaches, verified |
| 0162 | Find Peak Element | ✅ | 3 approaches, verified |
| 0163 | Missing Ranges | ✅ | 2 approaches, verified |
| 0164 | Maximum Gap | ✅ | 4 approaches, verified |
| 0165 | Compare Version Numbers | ✅ | 2 approaches, verified |
| 0166 | Fraction to Recurring Decimal | ✅ | 2 approaches, verified |
| 0167 | Two Sum II - Input Array Is Sorted | ✅ | 4 approaches, verified |
| 0168 | Excel Sheet Column Title | ✅ | 3 approaches, verified |
| 0169 | Majority Element | ✅ | 6 approaches, verified |
| 0170 | Two Sum III - Data Structure Design | ✅ | 4 approaches, verified |
| 0171 | Excel Sheet Column Number | ✅ | 3 approaches, verified |
| 0172 | Factorial Trailing Zeroes | ✅ | 3 approaches, verified |
| 0173 | Binary Search Tree Iterator | ✅ | 3 approaches, verified |
| 0174 | Dungeon Game | ✅ | 5 approaches, verified |
| 0175 | Combine Two Tables | ⚠️ | **README.md missing** — main.go present (3 approaches) and runs correctly |
| 0176 | Second Highest Salary | ✅ | 3 approaches, verified |
| 0177 | Nth Highest Salary | ✅ | 4 approaches, verified |
| 0178 | Rank Scores | ✅ | 3 approaches, verified |
| 0179 | Largest Number | ✅ | 3 approaches, verified |
| 0180 | Consecutive Numbers | ❌ | **Folder does not exist** — needs full solve from scratch |
| 0181 | Employees Earning More Than Their Managers | ✅ | 3 approaches, verified |
| 0182 | Duplicate Emails | ✅ | 3 approaches, verified |
| 0183 | Customers Who Never Order | ✅ | 3 approaches, verified |
| 0184 | Department Highest Salary | ✅ | 3 approaches, verified |
| 0185 | Department Top Three Salaries | ⚠️ | **README.md missing** — main.go present (4 approaches) and runs correctly |
| 0186 | Reverse Words in a String II | ✅ | 3 approaches, verified |
| 0187 | Repeated DNA Sequences | ✅ | 3 approaches, verified |
| 0188 | Best Time to Buy and Sell Stock IV | ✅ | 4 approaches, verified |
| 0189 | Rotate Array | ✅ | 4 approaches, verified |
| 0190 | Reverse Bits | ⚠️ | **README.md missing** — main.go present (4 approaches) and runs correctly |
| 0191 | Number of 1 Bits | ✅ | 5 approaches, verified |
| 0192 | Word Frequency | ✅ | 3 approaches, verified |
| 0193 | Valid Phone Numbers | ✅ | 3 approaches, verified |
| 0194 | Transpose File | ✅ | 3 approaches, verified |
| 0195 | Tenth Line | ✅ | 3 approaches, verified |
| 0196 | Delete Duplicate Emails | ✅ | 3 approaches, verified |
| 0197 | Rising Temperature | ✅ | 3 approaches, verified |
| 0198 | House Robber | ⚠️ | **README.md missing** — main.go present (4 approaches) and runs correctly |
| 0199 | Binary Tree Right Side View | ⚠️ | **README.md missing** — main.go present (3 approaches) and runs correctly |
| 0200 | Number of Islands | ⚠️ | **README.md missing** — main.go present (3 approaches) and runs correctly |
| 0201 | Bitwise AND of Numbers Range | ✅ | 3 approaches, verified |
| 0202 | Happy Number | ✅ | 3 approaches, verified |
| 0203 | Remove Linked List Elements | ✅ | 3 approaches, verified |
| 0204 | Count Primes | ✅ | 4 approaches, verified |
| 0205 | Isomorphic Strings | ✅ | 3 approaches, verified |
| 0206 | Reverse Linked List | ✅ | 4 approaches, verified |
| 0207 | Course Schedule | ✅ | 3 approaches, verified |
| 0208 | Implement Trie (Prefix Tree) | ✅ | 4 approaches, verified |
| 0209 | Minimum Size Subarray Sum | ✅ | 3 approaches, verified |
| 0210 | Course Schedule II | ⚠️ | **README.md missing** — main.go present (3 approaches) and runs correctly |

### Gap summary for #131–210

- **1 fully missing:** #0180 (Consecutive Numbers)
- **7 missing README only** (code exists and runs correctly): #0175, #0185,
  #0190, #0198, #0199, #0200, #0210

---

## `/dsa/` reference library

31 concept files present, all substantive (92–404 lines). All links from
#131–210 READMEs resolve correctly. The two oldest files (`hash_map.md`,
`two_pointers.md`, 92–117 lines) are noticeably thinner than the 28 written
in the most recent batch (250–400 lines) — worth a depth pass later if
consistency matters.

---

## Next problem to solve

**#0211** (first problem after this audited range) — pending your decision
on how to handle the 8 gaps above first.
