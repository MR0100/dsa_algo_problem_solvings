package main

import "fmt"

// ── Approach 1: Greedy with Value-Symbol Table ────────────────────────────────
//
// greedyTable uses a table of (value, symbol) pairs in descending order and
// repeatedly subtracts the largest fitting value from num.
//
// Intuition:
//   Roman numeral construction is greedy: always use the largest symbol that
//   fits into the remaining value. The table includes the 6 subtractive pairs
//   (IV=4, IX=9, XL=40, XC=90, CD=400, CM=900) alongside the 7 additive ones.
//   By sorting descending and greedily subtracting, we get the correct numeral.
//
// Time:  O(1) — the table has 13 fixed entries; num ≤ 3999 means at most
//               ~15 iterations.
// Space: O(1) — output length is bounded by a constant.
func greedyTable(num int) string {
	vals := []int{1000, 900, 500, 400, 100, 90, 50, 40, 10, 9, 5, 4, 1}
	syms := []string{"M", "CM", "D", "CD", "C", "XC", "L", "XL", "X", "IX", "V", "IV", "I"}

	result := ""
	for i, v := range vals {
		for num >= v {
			result += syms[i]
			num -= v
		}
	}
	return result
}

// ── Approach 2: Digit-by-Digit Table Lookup ───────────────────────────────────
//
// digitByDigit extracts each decimal digit (thousands, hundreds, tens, ones)
// and maps it to its Roman numeral string using four pre-built lookup arrays.
//
// Intuition:
//   For any number 1–3999, each decimal digit position has at most 9 distinct
//   Roman representations. Pre-compute them in a 2-D array indexed by
//   [position][digit]. Then extract digits and concatenate.
//
// Time:  O(1) — four table lookups and a concatenation.
// Space: O(1) — fixed-size lookup tables.
func digitByDigit(num int) string {
	// Lookup tables for each digit position (index = digit value 0..9).
	thousands := []string{"", "M", "MM", "MMM"}
	hundreds := []string{"", "C", "CC", "CCC", "CD", "D", "DC", "DCC", "DCCC", "CM"}
	tens := []string{"", "X", "XX", "XXX", "XL", "L", "LX", "LXX", "LXXX", "XC"}
	ones := []string{"", "I", "II", "III", "IV", "V", "VI", "VII", "VIII", "IX"}

	return thousands[num/1000] +
		hundreds[(num%1000)/100] +
		tens[(num%100)/10] +
		ones[num%10]
}

func main() {
	examples := []struct {
		num    int
		expect string
	}{
		{3, "III"},
		{58, "LVIII"},
		{1994, "MCMXCIV"},
		{3999, "MMMCMXCIX"},
		{4, "IV"},
		{9, "IX"},
		{40, "XL"},
		{90, "XC"},
		{400, "CD"},
		{900, "CM"},
	}

	approaches := []struct {
		name string
		fn   func(int) string
	}{
		{"Approach 1: Greedy table          O(1) T | O(1) S", greedyTable},
		{"Approach 2: Digit-by-digit table ✅ O(1) T | O(1) S", digitByDigit},
	}

	for _, ex := range examples {
		fmt.Printf("num=%-5d  expect=%q\n", ex.num, ex.expect)
		for _, ap := range approaches {
			fmt.Printf("  %-50s → %q\n", ap.name, ap.fn(ex.num))
		}
		fmt.Println()
	}
}
