package main

import (
	"fmt"
	"strings"
)

// ── Approach 1: Greedy Line Packing + Space Distribution ─────────────────────
//
// fullJustify solves Text Justification by greedily packing words onto lines,
// then distributing spaces according to the justification rules.
//
// Intuition:
//   1. Greedy packing: add words to the current line while they fit (sum of
//      word lengths + at least one space between each pair ≤ maxWidth).
//   2. Space distribution: for a line with k words, there are k-1 gaps.
//      Total spaces = maxWidth - sum(word lengths). Distribute extra spaces
//      round-robin from left to right. Last line: left-justify (single spaces).
//
// Algorithm:
//   i = 0
//   while i < len(words):
//     pack words[i..j] onto line while possible; update j
//     compute spaces and build line string
//     i = j
//
// Time:  O(n × maxWidth) — n words, each line build is O(maxWidth).
// Space: O(maxWidth) per line; O(total output) overall.
func fullJustify(words []string, maxWidth int) []string {
	var result []string
	i := 0
	n := len(words)

	for i < n {
		// pack as many words as possible onto this line
		lineLen := len(words[i])
		j := i + 1
		for j < n && lineLen+1+len(words[j]) <= maxWidth {
			lineLen += 1 + len(words[j])
			j++
		}
		// words[i..j-1] go on this line
		numWords := j - i
		numGaps := numWords - 1

		var line strings.Builder
		line.WriteString(words[i])

		if j == n || numWords == 1 {
			// last line or single word: left-justify (single spaces + pad right)
			for k := i + 1; k < j; k++ {
				line.WriteByte(' ')
				line.WriteString(words[k])
			}
			// pad with spaces on the right
			for line.Len() < maxWidth {
				line.WriteByte(' ')
			}
		} else {
			// regular line: distribute spaces evenly
			totalSpaces := maxWidth
			for k := i; k < j; k++ {
				totalSpaces -= len(words[k])
			}
			spacePerGap := totalSpaces / numGaps
			extraSpaces := totalSpaces % numGaps // first extraSpaces gaps get one extra

			for k := 1; k < numWords; k++ {
				spaces := spacePerGap
				if k-1 < extraSpaces {
					spaces++ // distribute extra spaces left to right
				}
				for s := 0; s < spaces; s++ {
					line.WriteByte(' ')
				}
				line.WriteString(words[i+k])
			}
		}

		result = append(result, line.String())
		i = j
	}

	return result
}

func main() {
	fmt.Println("=== Text Justification ===")

	w1 := []string{"This", "is", "an", "example", "of", "text", "justification."}
	r1 := fullJustify(w1, 16)
	fmt.Println("Expected:")
	fmt.Println("  \"This    is    an\"")
	fmt.Println("  \"example  of text\"")
	fmt.Println("  \"justification.  \"")
	fmt.Println("Got:")
	for _, line := range r1 {
		fmt.Printf("  %q  (len=%d)\n", line, len(line))
	}

	fmt.Println()
	w2 := []string{"What", "must", "be", "acknowledgment", "shall", "be"}
	r2 := fullJustify(w2, 16)
	fmt.Println("Expected:")
	fmt.Println("  \"What   must   be\"")
	fmt.Println("  \"acknowledgment  \"")
	fmt.Println("  \"shall be        \"")
	fmt.Println("Got:")
	for _, line := range r2 {
		fmt.Printf("  %q  (len=%d)\n", line, len(line))
	}

	fmt.Println()
	w3 := []string{"Science", "is", "what", "we", "understand", "well", "enough", "to", "explain",
		"to", "a", "computer.", "Art", "is", "everything", "else", "we", "do"}
	r3 := fullJustify(w3, 20)
	fmt.Println("Got (maxWidth=20):")
	for _, line := range r3 {
		fmt.Printf("  %q  (len=%d)\n", line, len(line))
	}
}
