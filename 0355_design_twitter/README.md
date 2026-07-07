# 0355 — Design Twitter

> LeetCode #355 · Difficulty: Medium
> **Categories:** Design, Hash Table, Heap (Priority Queue), Linked List

---

## Problem Statement

Design a simplified version of Twitter where users can post tweets, follow/unfollow
another user, and is able to see the `10` most recent tweets in the user's news
feed.

Implement the `Twitter` class:

- `Twitter()` Initializes your twitter object.
- `void postTweet(int userId, int tweetId)` Composes a new tweet with ID
  `tweetId` by the user `userId`. Each call to this function will be made with a
  unique `tweetId`.
- `List<Integer> getNewsFeed(int userId)` Retrieves the `10` most recent tweet
  IDs in the user's news feed. Each item in the news feed must be posted by users
  who the user followed or by the user themself. Tweets must be **ordered from
  most recent to least recent**.
- `void follow(int followerId, int followeeId)` The user with ID `followerId`
  started following the user with ID `followeeId`.
- `void unfollow(int followerId, int followeeId)` The user with ID `followerId`
  started unfollowing the user with ID `followeeId`.

**Example 1:**

```
Input
["Twitter", "postTweet", "getNewsFeed", "follow", "postTweet", "getNewsFeed",
 "unfollow", "getNewsFeed"]
[[], [1, 5], [1], [1, 2], [2, 6], [1], [1, 2], [1]]
Output
[null, null, [5], null, null, [6, 5], null, [5]]

Explanation
Twitter twitter = new Twitter();
twitter.postTweet(1, 5); // User 1 posts a new tweet (id = 5).
twitter.getNewsFeed(1);  // User 1's news feed should return a list with 1 tweet
                         // id -> [5]. return [5]
twitter.follow(1, 2);    // User 1 follows user 2.
twitter.postTweet(2, 6); // User 2 posts a new tweet (id = 6).
twitter.getNewsFeed(1);  // User 1's news feed should return a list with 2 tweet
                         // ids -> [6, 5]. Tweet id 6 should precede tweet id 5
                         // because it is posted after tweet id 5.
twitter.unfollow(1, 2);  // User 1 unfollows user 2.
twitter.getNewsFeed(1);  // User 1's news feed should return a list with 1 tweet
                         // id -> [5], since user 1 is no longer following user 2.
```

**Constraints:**

- `1 <= userId, followerId, followeeId <= 500`
- `0 <= tweetId <= 10^4`
- All the tweets have **unique** IDs.
- At most `3 * 10^4` calls will be made to `postTweet`, `getNewsFeed`, `follow`,
  and `unfollow`.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★☆ High       | 2024          |
| Google     | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Twitter/X  | ★★★☆☆ Medium     | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Hash tables for state** — `userId → tweets` and `userId → followee-set`
  → see [`/dsa/hash_map.md`](/dsa/hash_map.md)
- **Heap / priority queue for k-way merge** — merge k sorted per-user tweet lists
  and pull only the 10 newest → see [`/dsa/heap_priority_queue.md`](/dsa/heap_priority_queue.md)
- **Global timestamp ordering** — a monotonically increasing clock makes recency
  comparable across users → see [`/dsa/design_data_structures.md`](/dsa/design_data_structures.md)

---

## Approaches Overview

| # | Approach | getNewsFeed | Space | When to use |
|---|----------|-------------|-------|-------------|
| 1 | Brute Force (collect all + sort) | O(T log T) | O(U + T) | Simple; few tweets per user |
| 2 | k-Way Merge Max-Heap (Optimal) | O(k + 10 log k) | O(k) | Many tweets; only 10 needed |

`T` = total tweets among the relevant users, `k` = number of source users,
`U` = number of users.

---

## Approach 1 — Brute Force (Collect All + Sort)

### Intuition
The feed is "the 10 newest tweets among {self} ∪ followees". The simplest correct
recipe: gather every tweet from those users into one list, sort by timestamp
descending, and take the first 10. A global `clock` tags every post so recency is
comparable across users.

### Algorithm
1. `postTweet`: `clock++`; append `{id, clock}` to the user's tweet slice.
2. `follow` / `unfollow`: add/remove followee in the follower's set.
3. `getNewsFeed`: build `sources = {userId} ∪ followees`; concatenate all their
   tweets; sort by `time` descending; return the first up-to-10 ids.

### Complexity
- **Time:** `getNewsFeed` O(T log T) to sort all T relevant tweets; other ops O(1).
- **Space:** O(U + T) for follow-sets and tweets.

### Code
```go
func (t *bruteForceTwitter) GetNewsFeed(userId int) []int {
	sources := map[int]bool{userId: true}
	for f := range t.follows[userId] {
		sources[f] = true
	}

	var all []tweet
	for u := range sources {
		all = append(all, t.tweets[u]...) // gather every relevant tweet
	}
	sort.Slice(all, func(i, j int) bool { return all[i].time > all[j].time })

	res := []int{}
	for i := 0; i < len(all) && i < 10; i++ {
		res = append(res, all[i].id)
	}
	return res
}
```

### Dry Run
Operations from the example. `clock` starts at 0.

| op | effect | state |
|----|--------|-------|
| postTweet(1,5) | clock=1 | tweets[1]=[{5,1}] |
| getNewsFeed(1) | sources={1}; all=[{5,1}]; sorted → [5] | returns **[5]** |
| follow(1,2) | — | follows[1]={2} |
| postTweet(2,6) | clock=2 | tweets[2]=[{6,2}] |
| getNewsFeed(1) | sources={1,2}; all=[{5,1},{6,2}]; sort desc → [{6,2},{5,1}] | returns **[6,5]** |
| unfollow(1,2) | — | follows[1]={} |
| getNewsFeed(1) | sources={1}; all=[{5,1}] | returns **[5]** |

---

## Approach 2 — k-Way Merge with a Max-Heap (Optimal)

### Intuition
Each user's own tweet list is already sorted by time (append order), so the
newest is the last element. We only want the 10 newest across `k` such sorted
lists — a textbook **k-way merge**. Seed a max-heap (keyed by timestamp) with each
source's newest tweet, then pop 10 times, each time pushing the next-older tweet
from the same list. We never sort more than we consume.

### Algorithm
1. `sources = {userId} ∪ followees`.
2. For each source with tweets, push its **last** (newest) tweet into a max-heap,
   remembering the owner and its index in that list.
3. Repeat up to 10 times: pop the newest; append its id; if the popped tweet had
   an older neighbour (`idx > 0`), push `list[idx-1]`.

### Complexity
- **Time:** `getNewsFeed` O(k + 10 log k) — seed k sources, then ≤10 heap ops.
- **Space:** O(k) — the heap holds at most one entry per source.

### Code
```go
func (t *Twitter) GetNewsFeed(userId int) []int {
	sources := map[int]bool{userId: true}
	for f := range t.follows[userId] {
		sources[f] = true
	}

	h := &maxHeap{}
	heap.Init(h)
	for u := range sources {
		lst := t.tweets[u]
		if len(lst) > 0 {
			last := len(lst) - 1
			heap.Push(h, heapItem{t: lst[last], owner: u, idx: last})
		}
	}

	res := []int{}
	for h.Len() > 0 && len(res) < 10 {
		top := heap.Pop(h).(heapItem) // globally newest remaining tweet
		res = append(res, top.t.id)
		if top.idx > 0 {
			prev := top.idx - 1
			heap.Push(h, heapItem{t: t.tweets[top.owner][prev], owner: top.owner, idx: prev})
		}
	}
	return res
}
```

### Dry Run
`getNewsFeed(1)` after `postTweet(1,5)`, `follow(1,2)`, `postTweet(2,6)`.
`tweets[1]=[{5,1}]`, `tweets[2]=[{6,2}]`.

| Step | heap (by time) | pop | res | push older? |
|------|----------------|-----|-----|-------------|
| seed | {5,t1,own1,i0}, {6,t2,own2,i0} | — | [] | — |
| 1 | pop max time=2 → {6} | {6} | [6] | idx 0 → none |
| 2 | {5,t1} | {5} | [6,5] | idx 0 → none |
| 3 | empty | — | [6,5] | stop |

Returns **[6,5]**, newest first — matching the expected output.

---

## Key Takeaways

- **A global timestamp is the linchpin.** Per-user counters can't be compared;
  one monotonically increasing clock makes recency total-ordered.
- **Only fetch what you need.** The feed wants 10 items, so a k-way merge with a
  heap beats sorting every tweet — pop 10, not T.
- Storing tweets in append order keeps each user's list pre-sorted, so the heap
  merge needs no per-user sorting.
- Follow-sets as `map[int]map[int]bool` give O(1) follow/unfollow and easy union
  with the user themself.

---

## Related Problems

- LeetCode #23 — Merge k Sorted Lists (the k-way merge core)
- LeetCode #253 — Meeting Rooms II (heap-driven design)
- LeetCode #146 — LRU Cache (stateful map design)
- LeetCode #362 — Design Hit Counter (time-windowed design)
