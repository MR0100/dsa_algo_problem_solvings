# Progress Tracker

Tracks completion status for every LeetCode problem in this repo, per
[CLAUDE.md](CLAUDE.md)'s required structure (exactly `main.go` + `README.md`,
every approach brute-force → optimal, 10-section README, `go run` verified).

Legend: ✅ Complete &nbsp;·&nbsp; ⚠️ Partial (code exists, README missing/incomplete) &nbsp;·&nbsp; ❌ Missing entirely

**Last audited: 2026-07-07** — full-repo verification of #0001–#0400.

---

## Summary

| Range | Status | Count | Notes |
|-------|--------|-------|-------|
| #0001 – #0130 | ✅ Complete | 130 | Per-approach depth backfilled (Dry Run / Complexity / Code blocks); 11 stale `/dsa` links repaired. |
| #0131 – #0210 | ✅ Complete | 80 | Prior 8 gaps closed: #0180 solved from scratch; READMEs added for #0175, #0185, #0190, #0198, #0199, #0200, #0210. |
| #0211 – #0300 | ✅ Complete | 90 | Solved this project; strict per-approach audit clean. |
| #0301 – #0400 | ✅ Complete | 100 | Solved this project (recovered from a mid-batch rate-limit failure; regenerated in small waves). |
| **Total** | **✅** | **400 / 400** | Contiguous #0001–#0400, no duplicates, no gaps. |

---

## Verification results (#0001 – #0400, 2026-07-07)

All checks pass across the full 400-problem range:

- ✅ **Presence & structure** — 400/400 folders, each with exactly `main.go` + `README.md`, both non-empty, correct 4-digit snake_case naming.
- ✅ **Compiles** — `go build` succeeds for all 400.
- ✅ **Runs** — `go run main.go` exits 0 for all 400; approach outputs were checked against their inline `// expected` comments by the solving/verifying agents.
- ✅ **Formatting** — all `main.go` are `gofmt`-clean.
- ✅ **Strict per-approach audit** — every `## Approach` section in all 400 READMEs contains Intuition + Algorithm + Complexity + a fenced `go` Code block + a step-by-step Dry Run. (Uses a per-section parser, so uneven distribution can't hide a gap.)
- ✅ **README sections** — all 10 required top-level sections present and in order; Company Frequency disclaimer present.
- ✅ **`/dsa/` links** — **0 broken links** repo-wide (the 11 stale links previously flagged in #0001–#0130 — `arrays.md`, `math.md`, `heap.md`, `dynamic_programming.md` — were repointed to the current file names).

---

## `/dsa/` reference library

**43 concept files** present, all substantive; every README link across #0001–#0400 resolves (0 broken links).

A concept-link audit (2026-07-07) found many bullets linked to a "closest" file rather than the right one, and repaired them: every "DSA Concepts Used" bullet now points to the semantically correct file. **12 new concept files** were authored (Go-first, full worked examples) and their bullets repointed:
`arrays`, `geometry`, `line_sweep`, `boyer_moore_voting`, `counting_sort`,
`k_way_merge`, `digit_dp`, `game_theory`, `interval_dp`, `tree_dp`,
`longest_increasing_subsequence`, `shuffle`.
Also repointed to existing files: Manacher's → `string_algorithms`, Monotonic Stack → `monotonic_stack`, `Array*` → `arrays` (previously mislinked to `hash_map`/`matrix_traversal`/`prefix_sum`).

---

## Remaining work

- **#0401 onward — not yet started.** Next problem to solve is **#0401 (Binary Watch)**.
- **Optional:** author dedicated `/dsa/` pages for the concepts listed above and repoint the links (currently they resolve to the nearest existing file, so this is a polish item, not a defect).
