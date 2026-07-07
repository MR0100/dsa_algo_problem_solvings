package main

import "fmt"

// file simulates the hidden file object that sits behind the read4 API.
// On LeetCode this is invisible; here we need it to drive the examples.
type file struct {
	data string // full file contents
	pos  int    // read4's own file pointer (persists across read4 calls)
}

// read4 is the provided API: it reads up to 4 consecutive characters from the
// file into buf4 and returns the number of characters actually read.
func (f *file) read4(buf4 []byte) int {
	count := 0
	// copy until we have 4 chars or the file is exhausted
	for count < 4 && f.pos < len(f.data) {
		buf4[count] = f.data[f.pos]
		count++
		f.pos++
	}
	return count
}

// ── Approach 1: Read Whole File (Brute Force) ────────────────────────────────
//
// bruteForceReadAll solves Read N Characters Given Read4 by slurping the
// entire file first.
//
// Intuition:
//
//	The simplest correct thing: keep calling read4 until it signals EOF
//	(returns < 4), collecting everything into one big internal buffer. Then
//	copy the first min(n, fileLength) characters into the destination.
//
// Algorithm:
//  1. Repeatedly call read4, appending buf4[:cnt] to an internal slice.
//  2. Stop when read4 returns fewer than 4 characters (EOF reached).
//  3. total = min(n, len(all)); copy all[:total] into buf; return total.
//
// Time:  O(F) where F = file length — reads the whole file even if n is tiny.
// Space: O(F) — the internal buffer holds the entire file.
func bruteForceReadAll(read4 func([]byte) int, buf []byte, n int) int {
	all := []byte{}         // internal buffer collecting the whole file
	buf4 := make([]byte, 4) // scratch buffer for the read4 API
	for {
		cnt := read4(buf4)
		all = append(all, buf4[:cnt]...) // keep only the valid characters
		if cnt < 4 {
			break // read4 returned short → end of file
		}
	}
	total := n
	if len(all) < n {
		total = len(all) // file shorter than requested → return what exists
	}
	copy(buf, all[:total]) // hand exactly total chars to the caller
	return total
}

// ── Approach 2: Direct Copy (Optimal) ────────────────────────────────────────
//
// directCopy solves Read N Characters Given Read4 by copying chunks straight
// into the destination buffer, stopping as soon as n is satisfied.
//
// Intuition:
//
//	We never need more than n characters, so read 4 at a time and write each
//	chunk directly into buf at offset total. Two stop conditions: we have n
//	characters, or read4 hits EOF (returns 0). If the last chunk overshoots n,
//	clamp it — read is called only once, so discarding the extras is safe.
//
// Algorithm:
//  1. total = 0.
//  2. While total < n:
//     a. cnt = read4(buf4); if cnt == 0 → EOF, stop.
//     b. Clamp: if total+cnt > n, cnt = n-total (don't overflow the request).
//     c. copy(buf[total:], buf4[:cnt]); total += cnt.
//  3. Return total.
//
// Time:  O(n) — at most ⌈n/4⌉ read4 calls, each O(1) work per character.
// Space: O(1) — a single fixed 4-byte scratch buffer.
func directCopy(read4 func([]byte) int, buf []byte, n int) int {
	buf4 := make([]byte, 4) // fixed scratch buffer for the API
	total := 0              // characters delivered to buf so far
	for total < n {
		cnt := read4(buf4)
		if cnt == 0 {
			break // EOF: file has no more characters
		}
		if total+cnt > n {
			cnt = n - total // clamp: caller asked for exactly n
		}
		copy(buf[total:], buf4[:cnt]) // write chunk directly at offset total
		total += cnt
	}
	return total
}

// runCase executes one solver against a fresh simulated file and prints result.
func runCase(solver func(func([]byte) int, []byte, int) int, data string, n int) {
	f := &file{data: data} // fresh file: read4's pointer starts at 0
	buf := make([]byte, n) // destination buffer sized for the request
	got := solver(f.read4, buf, n)
	fmt.Printf("file=%q n=%d  →  read=%d buf=%q\n", data, n, got, string(buf[:got]))
}

func main() {
	fmt.Println("=== Approach 1: Read Whole File (Brute Force) ===")
	runCase(bruteForceReadAll, "abc", 4)   // expected read=3 buf="abc"
	runCase(bruteForceReadAll, "abcde", 5) // expected read=5 buf="abcde"

	fmt.Println("=== Approach 2: Direct Copy (Optimal) ===")
	runCase(directCopy, "abc", 4)     // expected read=3 buf="abc"
	runCase(directCopy, "abcde", 5)   // expected read=5 buf="abcde"
	runCase(directCopy, "abcdefg", 3) // expected read=3 buf="abc" (n < file length)
}
