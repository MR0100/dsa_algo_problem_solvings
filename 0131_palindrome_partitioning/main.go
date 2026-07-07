package main

import "fmt"

// ── Approach 1: Brute-Force Backtracking ─────────────────────────────────────
//
// bruteForceBacktracking solves Palindrome Partitioning using plain
// backtracking with an on-the-fly palindrome check.
//
// Intuition: every partition is defined by a set of cut positions. At each
// position `start` we try every possible next piece s[start..end]; if that
// piece is a palindrome we keep it and recurse on the rest of the string.
// When `start` reaches the end of the string, the current path is one valid
// partition.
//
// Algorithm:
//  1. dfs(start): if start == len(s), record a copy of the current path.
//  2. Otherwise, for every end in [start, len(s)-1]:
//     a. Check s[start..end] with a two-pointer palindrome scan.
//     b. If it is a palindrome, push it, recurse dfs(end+1), then pop (undo).
//
// Time:  O(n · 2^n) — up to 2^(n-1) partitions, each costing O(n) to check/copy.
// Space: O(n) auxiliary — recursion depth and current path (output excluded).
func bruteForceBacktracking(s string) [][]string {
	n := len(s)
	result := [][]string{} // all valid partitions collected here
	path := []string{}     // current partial partition being built

	// isPal checks s[i..j] inclusive with two converging pointers.
	isPal := func(i, j int) bool {
		for i < j {
			if s[i] != s[j] { // mismatch means not a palindrome
				return false
			}
			i++ // move both ends toward the middle
			j--
		}
		return true
	}

	var dfs func(start int)
	dfs = func(start int) {
		if start == n {
			// consumed the whole string: snapshot the path (copy, because
			// path's backing array keeps mutating during backtracking)
			part := make([]string, len(path))
			copy(part, path)
			result = append(result, part)
			return
		}
		for end := start; end < n; end++ {
			if isPal(start, end) { // only recurse on palindromic prefixes
				path = append(path, s[start:end+1]) // choose this piece
				dfs(end + 1)                        // partition the remainder
				path = path[:len(path)-1]           // undo the choice (backtrack)
			}
		}
	}

	dfs(0)
	return result
}

// ── Approach 2: Backtracking + DP Palindrome Table ───────────────────────────
//
// dpTableBacktracking solves Palindrome Partitioning using backtracking, but
// with all palindrome answers precomputed in an n×n DP table so each
// palindrome query is O(1).
//
// Intuition: the brute force re-scans the same substrings over and over.
// Palindromicity has optimal substructure — s[i..j] is a palindrome iff
// s[i] == s[j] AND s[i+1..j-1] is a palindrome — so a bottom-up table over
// increasing substring lengths answers every query once, in O(n²) total.
//
// Algorithm:
//  1. Build isPal[i][j] for all i <= j:
//     - length 1: always true;
//     - length 2: true iff the two characters match;
//     - length ≥ 3: s[i]==s[j] && isPal[i+1][j-1].
//  2. Run the same dfs as Approach 1, but test isPal[start][end] in O(1).
//
// Time:  O(n · 2^n) — partitions dominate; each check is now O(1) (table O(n²)).
// Space: O(n²) — the palindrome table (plus O(n) recursion depth).
func dpTableBacktracking(s string) [][]string {
	n := len(s)

	// isPal[i][j] == true iff s[i..j] is a palindrome.
	isPal := make([][]bool, n)
	for i := range isPal {
		isPal[i] = make([]bool, n)
		isPal[i][i] = true // every single character is a palindrome
	}
	// fill by increasing substring length so smaller answers already exist
	for length := 2; length <= n; length++ {
		for i := 0; i+length-1 < n; i++ {
			j := i + length - 1 // right end of the window
			if s[i] == s[j] {
				// inner part must be a palindrome too (or be empty, length 2)
				isPal[i][j] = length == 2 || isPal[i+1][j-1]
			}
		}
	}

	result := [][]string{}
	path := []string{}

	var dfs func(start int)
	dfs = func(start int) {
		if start == n {
			part := make([]string, len(path)) // snapshot the finished partition
			copy(part, path)
			result = append(result, part)
			return
		}
		for end := start; end < n; end++ {
			if isPal[start][end] { // O(1) table lookup instead of a scan
				path = append(path, s[start:end+1])
				dfs(end + 1)
				path = path[:len(path)-1] // backtrack
			}
		}
	}

	dfs(0)
	return result
}

// ── Approach 3: Memoized Suffix Partitions (Optimal reuse) ───────────────────
//
// memoizedSuffixes solves Palindrome Partitioning with top-down DP: for every
// start index it computes (once) the full list of partitions of the suffix
// s[start:], memoizing the answer so shared suffixes are never re-partitioned.
//
// Intuition: dfs(start) in the backtracking solutions is re-entered many
// times with the same argument along different paths. The set of partitions
// of s[start:] does not depend on how we got there, so it can be cached:
// partitions(start) = { [s[start..end]] + rest | s[start..end] palindrome,
// rest ∈ partitions(end+1) }.
//
// Algorithm:
//  1. Build the same O(n²) palindrome table as Approach 2.
//  2. solve(start): if cached, return it. Base case solve(n) = [ [] ]
//     (one empty partition). Otherwise for each palindromic prefix
//     s[start..end], prepend it to every partition of solve(end+1).
//  3. Answer is solve(0).
//
// Time:  O(n · 2^n) — output-size bound; each suffix is solved once and shared.
// Space: O(n · 2^n) — the memo stores complete partition lists per suffix.
func memoizedSuffixes(s string) [][]string {
	n := len(s)

	// palindrome table, identical construction to Approach 2
	isPal := make([][]bool, n)
	for i := range isPal {
		isPal[i] = make([]bool, n)
		isPal[i][i] = true
	}
	for length := 2; length <= n; length++ {
		for i := 0; i+length-1 < n; i++ {
			j := i + length - 1
			if s[i] == s[j] {
				isPal[i][j] = length == 2 || isPal[i+1][j-1]
			}
		}
	}

	memo := make(map[int][][]string) // start index → all partitions of s[start:]

	var solve func(start int) [][]string
	solve = func(start int) [][]string {
		if start == n {
			// exactly one way to partition the empty suffix: the empty list
			return [][]string{{}}
		}
		if cached, ok := memo[start]; ok {
			return cached // suffix already fully solved
		}
		res := [][]string{}
		for end := start; end < n; end++ {
			if !isPal[start][end] {
				continue // first piece must be a palindrome
			}
			piece := s[start : end+1]
			for _, rest := range solve(end + 1) {
				// build a fresh slice: piece followed by the cached tail
				part := make([]string, 0, len(rest)+1)
				part = append(part, piece)
				part = append(part, rest...)
				res = append(res, part)
			}
		}
		memo[start] = res // cache before returning
		return res
	}

	return solve(0)
}

func main() {
	fmt.Println("=== Approach 1: Brute-Force Backtracking ===")
	fmt.Println(bruteForceBacktracking("aab")) // [[a a b] [aa b]]
	fmt.Println(bruteForceBacktracking("a"))   // [[a]]

	fmt.Println("=== Approach 2: Backtracking + DP Palindrome Table ===")
	fmt.Println(dpTableBacktracking("aab")) // [[a a b] [aa b]]
	fmt.Println(dpTableBacktracking("a"))   // [[a]]

	fmt.Println("=== Approach 3: Memoized Suffix Partitions (Optimal reuse) ===")
	fmt.Println(memoizedSuffixes("aab")) // [[a a b] [aa b]]
	fmt.Println(memoizedSuffixes("a"))   // [[a]]
}
