# 0068 — Text Justification

> LeetCode #68 · Difficulty: Hard
> **Categories:** Array, String, Simulation

---

## Problem Statement

Given an array of strings `words` and a width `maxWidth`, format the text such that each line has exactly `maxWidth` characters and is **fully** (left and right) justified.

You should pack your words in a greedy approach; that is, pack as many words as you can in each line. Pad extra spaces `' '` when necessary so that each line has exactly `maxWidth` characters.

Extra spaces between words should be distributed as **evenly** as possible. If the number of spaces on a line does not divide evenly between words, the empty slots on the **left** will be assigned more spaces than the slots on the right.

For the **last line** of text, it should be **left-justified**, and no extra space is inserted between words.

**Note:** A word is defined as a character sequence consisting of non-space characters only. Each word's length is guaranteed to be greater than `0` and not exceed `maxWidth`. The input array `words` contains at least one word.

**Example 1**
```
Input:  words = ["This","is","an","example","of","text","justification."], maxWidth = 16
Output: ["This    is    an","example  of text","justification.  "]
```

**Example 2**
```
Input:  words = ["What","must","be","acknowledgment","shall","be"], maxWidth = 16
Output: ["What   must   be","acknowledgment  ","shall be        "]
```

**Example 3**
```
Input:  words = ["Science","is","what","we","understand","well","enough","to","explain","to","a","computer.","Art","is","everything","else","we","do"], maxWidth = 20
Output: ["Science  is  what we","understand      well","enough to explain to","a  computer.  Art is","everything  else  we","do                  "]
```

**Constraints**
- `1 <= words.length <= 300`
- `1 <= words[i].length <= 20`
- `words[i]` consists of only English letters and symbols.
- `1 <= maxWidth <= 100`
- `words[i].length <= maxWidth`

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Google    | ★★★★★ Very High | 2024          |
| Amazon    | ★★★★☆ High      | 2024          |
| Meta      | ★★★★☆ High      | 2024          |
| Microsoft | ★★★★☆ High      | 2023          |
| Bloomberg | ★★★★☆ High      | 2023          |
| Apple     | ★★★☆☆ Medium    | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Greedy Line Packing** — pack as many words as possible per line (greedy = always fill to maximum).
- **Arithmetic Space Distribution** — `spacePerGap = total / gaps`, `extra = total % gaps`; first `extra` gaps get one additional space.

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Greedy Pack + Space Math ✅ | O(n × maxWidth) | O(maxWidth) | The only correct approach; simulation |

---

## Approach 1 — Greedy Pack + Space Distribution (Recommended ✅)

### Intuition
**Phase 1 (Greedy packing):** Greedily add words to the current line while `len(currentWords) + (numWords - 1) * 1 ≤ maxWidth` (one space minimum between words). When adding the next word would exceed `maxWidth`, close the line.

**Phase 2 (Space distribution):** For a line with `k` words, there are `k-1` gaps. Extra spaces = `maxWidth - sum(word lengths)`. Distribute: `spacePerGap = extra / (k-1)`, `remainder = extra % (k-1)`. The first `remainder` gaps get `spacePerGap + 1` spaces; the rest get `spacePerGap`.

**Special cases:**
- **Single-word line:** left-justify (one word + spaces to pad to `maxWidth`).
- **Last line:** left-justify (single spaces between words + spaces to pad right).

### Algorithm
```
i = 0
while i < n:
  compute j = last word index that fits on the line
  if last line or single word: left-justify
  else: compute spacePerGap and extra; distribute left-to-right
  append line; i = j
```

### Complexity
- **Time:** O(n × maxWidth) — n words, each line build is O(maxWidth).
- **Space:** O(maxWidth) — one line buffer at a time (output not counted).

### Dry Run — `words = ["This","is","an","example","of","text","justification."]`, `maxWidth = 16`
```
Line 1: "This" (4), "is" (2), "an" (2) → 4+1+2+1+2=10 ≤ 16.
        "example" → 10+1+7=18 > 16. Stop at j=3 (words 0-2).
        Gaps=2. totalSpaces=16-4-2-2=8. spacePerGap=4, extra=0.
        "This" + "    " + "is" + "    " + "an" = "This    is    an" ✓

Line 2: "example" (7), "of" (2), "text" (4) → 7+1+2+1+4=15 ≤ 16.
        "justification." → 15+1+14=30 > 16. Stop.
        Gaps=2. totalSpaces=16-7-2-4=3. spacePerGap=1, extra=1.
        "example" + "  " + "of" + " " + "text" = "example  of text" ✓

Line 3: "justification." — last line → left-justify + pad.
        "justification." + "  " = "justification.  " ✓
```

---

## Key Takeaways

- **Three distinct cases to handle**: regular (full-justify), single-word (pad right), last line (left-justify).
- **Extra spaces go to leftmost gaps** — `remainder = total % gaps` first gaps get one extra space.
- **Last line must be left-justified** — a common source of bugs: don't apply full-justification to the last line.
- **Line length invariant** — every line must be exactly `maxWidth` characters. Verify with `len(line) == maxWidth` during testing.
- **Greedy packing condition** — track `lineLen = sum(word lengths) + (numWords - 1)` (1 space minimum); add next word if `lineLen + 1 + len(nextWord) ≤ maxWidth`.

---

## Related Problems

- LeetCode #6 — Zigzag Conversion (string formatting/layout)
- LeetCode #38 — Count and Say (iterative string building)
