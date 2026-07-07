package main

import (
	"fmt"
	"math"
	"sort"
)

// LeetCode 262 — Trips and Users.
//
// The original problem is a SQL one. Two tables:
//
//	Trips:  id, client_id, driver_id, city_id, status, request_at
//	          status ∈ {completed, cancelled_by_driver, cancelled_by_client}
//	          request_at is a date string "YYYY-MM-DD"
//	Users:  users_id, banned ("Yes"/"No"), role ∈ {client, driver, partner}
//
// The CANCELLATION RATE of a day = (# cancelled trips whose BOTH client and
// driver are unbanned) / (# trips that day whose BOTH client and driver are
// unbanned), rounded to 2 decimals. Report the rate for each day in
// [2013-10-01, 2013-10-03] that has at least one qualifying trip.
//
// The accepted SQL answer:
//
//	SELECT t.request_at AS Day,
//	       ROUND(
//	         SUM(CASE WHEN t.status != 'completed' THEN 1 ELSE 0 END) / COUNT(*),
//	         2) AS "Cancellation Rate"
//	FROM Trips t
//	JOIN Users c ON t.client_id = c.users_id AND c.banned = 'No'
//	JOIN Users d ON t.driver_id = d.users_id AND d.banned = 'No'
//	WHERE t.request_at BETWEEN '2013-10-01' AND '2013-10-03'
//	GROUP BY t.request_at;
//
// Below we model the tables as Go structs and re-implement the query as an
// in-memory group-by / aggregation, using two join algorithms.

// Trip mirrors one row of the Trips table.
type Trip struct {
	ID        int
	ClientID  int
	DriverID  int
	CityID    int
	Status    string // "completed" | "cancelled_by_driver" | "cancelled_by_client"
	RequestAt string // "YYYY-MM-DD"
}

// User mirrors one row of the Users table.
type User struct {
	UsersID int
	Banned  string // "Yes" | "No"
	Role    string // "client" | "driver" | "partner"
}

// Result is one output row: a day and its cancellation rate (2 decimals).
type Result struct {
	Day  string
	Rate float64
}

// round2 rounds x to two decimal places, matching SQL ROUND(x, 2).
func round2(x float64) float64 {
	return math.Round(x*100) / 100 // scale, round to nearest int, scale back
}

// ── Approach 1: Brute Force (Nested-Loop Join + Group-By) ────────────────────
//
// nestedLoopJoin solves Trips and Users by, for each in-range trip, scanning
// the Users table to check both client and driver are unbanned, then grouping
// the surviving trips by day and computing the cancellation ratio.
//
// Intuition:
//
//	Reproduce the two JOINs literally: for every trip in the date window, look
//	up its client and driver by linear scan and keep the trip only if BOTH are
//	unbanned. Bucket the kept trips by request_at; per bucket the rate is
//	(non-completed count)/(total count).
//
// Algorithm:
//  1. For each trip with request_at in [2013-10-01, 2013-10-03]:
//  2. Linear-scan Users to find the client and driver rows.
//  3. If either is missing or banned, skip the trip.
//  4. Otherwise increment the day's total and (if status != completed) its
//     cancelled count.
//  5. For each day, rate = cancelled/total; sort days ascending.
//
// Time:  O(T·U) — every kept trip scans the whole Users table.
// Space: O(D) — one bucket per distinct day.
func nestedLoopJoin(trips []Trip, users []User) []Result {
	total := map[string]int{}     // day → count of qualifying trips
	cancelled := map[string]int{} // day → count of qualifying non-completed trips
	inRange := func(d string) bool {
		return d >= "2013-10-01" && d <= "2013-10-03" // lexical works for ISO dates
	}
	bannedOf := func(id int) (string, bool) { // linear scan lookup: (banned, found)
		for _, u := range users {
			if u.UsersID == id {
				return u.Banned, true
			}
		}
		return "", false
	}
	for _, t := range trips {
		if !inRange(t.RequestAt) { // WHERE request_at BETWEEN ...
			continue
		}
		cb, cok := bannedOf(t.ClientID) // JOIN Users c ON client_id
		db, dok := bannedOf(t.DriverID) // JOIN Users d ON driver_id
		if !cok || !dok || cb != "No" || db != "No" {
			continue // drop trips where either party is missing or banned
		}
		total[t.RequestAt]++
		if t.Status != "completed" { // status != 'completed' ⇒ a cancellation
			cancelled[t.RequestAt]++
		}
	}
	return finalize(total, cancelled)
}

// ── Approach 2: Hash Join + Group-By (Optimal) ───────────────────────────────
//
// hashJoin solves Trips and Users by first indexing the Users table in a hash
// map (id → banned) so each trip's client/driver lookup is O(1), then grouping
// by day exactly as before.
//
// Intuition:
//
//	The nested-loop join wastes time rescanning Users for every trip. Build a
//	hash map users_id → banned once; then a single pass over the trips does two
//	O(1) lookups each, and accumulates the per-day counts on the fly.
//
// Algorithm:
//  1. Build map banned[users_id] = "Yes"/"No".
//  2. For each in-range trip, look up client and driver in O(1); keep only if
//     both are present and unbanned.
//  3. Accumulate total and cancelled counts per day, then compute rates.
//
// Time:  O(T + U) — one pass to index users, one pass over trips.
// Space: O(U + D) — the user index plus one bucket per day.
func hashJoin(trips []Trip, users []User) []Result {
	banned := make(map[int]string, len(users)) // users_id → "Yes"/"No"
	for _, u := range users {
		banned[u.UsersID] = u.Banned // index the whole Users table once
	}
	total := map[string]int{}
	cancelled := map[string]int{}
	for _, t := range trips {
		if t.RequestAt < "2013-10-01" || t.RequestAt > "2013-10-03" {
			continue // outside the reporting window
		}
		cb, cok := banned[t.ClientID] // O(1) client lookup
		db, dok := banned[t.DriverID] // O(1) driver lookup
		if !cok || !dok || cb != "No" || db != "No" {
			continue // either party unknown or banned ⇒ exclude
		}
		total[t.RequestAt]++
		if t.Status != "completed" {
			cancelled[t.RequestAt]++
		}
	}
	return finalize(total, cancelled)
}

// finalize turns the per-day counts into sorted Result rows with 2-dp rates.
func finalize(total, cancelled map[string]int) []Result {
	res := make([]Result, 0, len(total))
	for day, tot := range total {
		rate := round2(float64(cancelled[day]) / float64(tot)) // cancelled/total
		res = append(res, Result{Day: day, Rate: rate})
	}
	sort.Slice(res, func(i, j int) bool { return res[i].Day < res[j].Day }) // stable, ascending day
	return res
}

// printResults renders the result set like the LeetCode expected output.
func printResults(res []Result) {
	for _, r := range res {
		fmt.Printf("| %s | %.2f |\n", r.Day, r.Rate)
	}
}

func main() {
	// Official example tables.
	users := []User{
		{1, "No", "client"},
		{2, "Yes", "client"},
		{3, "No", "client"},
		{4, "No", "client"},
		{10, "No", "driver"},
		{11, "No", "driver"},
		{12, "No", "driver"},
		{13, "No", "driver"},
	}
	trips := []Trip{
		{1, 1, 10, 1, "completed", "2013-10-01"},
		{2, 2, 11, 1, "cancelled_by_driver", "2013-10-01"},
		{3, 3, 12, 6, "completed", "2013-10-01"},
		{4, 4, 13, 6, "cancelled_by_client", "2013-10-01"},
		{5, 1, 10, 1, "completed", "2013-10-02"},
		{6, 2, 11, 6, "completed", "2013-10-02"},
		{7, 3, 12, 6, "completed", "2013-10-02"},
		{8, 2, 12, 12, "completed", "2013-10-03"},
		{9, 3, 10, 12, "completed", "2013-10-03"},
		{10, 4, 13, 12, "cancelled_by_driver", "2013-10-03"},
	}

	// Expected output (trips with a banned party, e.g. user 2, are excluded):
	// | 2013-10-01 | 0.33 |
	// | 2013-10-02 | 0.00 |
	// | 2013-10-03 | 0.50 |

	fmt.Println("=== Approach 1: Nested-Loop Join + Group-By ===")
	printResults(nestedLoopJoin(trips, users)) // 2013-10-01:0.33, 10-02:0.00, 10-03:0.50

	fmt.Println("=== Approach 2: Hash Join + Group-By (Optimal) ===")
	printResults(hashJoin(trips, users)) // same three rows
}
