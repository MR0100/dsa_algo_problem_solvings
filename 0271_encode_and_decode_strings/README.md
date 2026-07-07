# 0271 — Encode and Decode Strings

> LeetCode #271 · Difficulty: Medium
> **Categories:** String, Design, Array

---

## Problem Statement

Design an algorithm to encode a **list of strings** to a **single string**. The encoded string is then sent over the network and is decoded back to the original list of strings.

Machine 1 (sender) has the function:

```
string encode(vector<string> strs) {
  // ... your code
  return encoded_string;
}
```

Machine 2 (receiver) has the function:

```
vector<string> decode(string s) {
  //... your code
  return strs;
}
```

So Machine 1 does:

```
string encoded_string = encode(strs);
```

and Machine 2 does:

```
vector<string> strs2 = decode(encoded_string);
```

`strs2` in Machine 2 should be the same as `strs` in Machine 1.

Implement the `encode` and `decode` methods.

You are not allowed to solve the problem using any serialize methods (such as `eval`).

**Example 1:**

```
Input: dummy_input = ["Hello","World"]
Output: ["Hello","World"]
Explanation:
Machine 1:
Codec encoder = new Codec();
String msg = encoder.encode(strs);
Machine 1 ---msg---> Machine 2

Machine 2:
Codec decoder = new Codec();
String[] strs = decoder.decode(msg);
```

**Example 2:**

```
Input: dummy_input = [""]
Output: [""]
```

**Constraints:**

- `1 <= strs.length <= 200`
- `0 <= strs[i].length <= 200`
- `strs[i]` contains any possible characters out of `256` valid ASCII characters.

**Follow-up:** Could you write a generalized algorithm to work on any possible set of characters?

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★★★ Very High  | 2024          |
| Meta       | ★★★★☆ High       | 2024          |
| Amazon     | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Bloomberg  | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Design / Codec** — the task is to design a self-describing wire format and its inverse parser; the same skill powers serialize/deserialize problems → see [`/dsa/design_data_structures.md`](/dsa/design_data_structures.md)
- **String Parsing** — decoding is a linear scan that reads a length header, then consumes an exact byte span; classic hand-rolled parsing → see [`/dsa/string_algorithms.md`](/dsa/string_algorithms.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Length-Prefix (chunked) | O(N) | O(N) | The correct, delimiter-free answer; handles any characters |
| 2 | Escaping a Delimiter | O(N) | O(N) | When you want a human-readable stream and still be safe |
| 3 | Non-ASCII Sentinel | O(N) | O(N) | Only if you can guarantee the sentinel never appears — fragile |

---

## Approach 1 — Length-Prefix (Chunked)

### Intuition

Any delimiter you pick can also appear inside a payload, creating ambiguity. Sidestep the problem entirely: instead of *searching* for a boundary inside the bytes, *announce* each string's exact byte length before it. The decoder reads the number, skips one separator, then blindly consumes exactly that many bytes — the payload can contain `#`, digits, newlines, anything, and nothing is ambiguous. This is exactly how HTTP chunked transfer encoding works.

### Algorithm

1. **Encode:** for each string `s`, append `len(s)` (decimal), then `'#'`, then `s`. Concatenate all chunks.
2. **Decode:** set `i = 0`. While `i < len(encoded)`:
   1. Scan from `i` to the next `'#'`; the digits in between are the chunk length `L`.
   2. The payload is the `L` bytes immediately after that `'#'`.
   3. Append the payload; advance `i` to just past the payload; repeat.

### Complexity

- **Time:** O(N) — N = total bytes; every byte is written once (encode) and read once (decode).
- **Space:** O(N) — the encoded string / decoded list.

### Code

```go
func lengthPrefixEncode(strs []string) string {
	var b strings.Builder
	for _, s := range strs {
		// len(s) is the BYTE length (Go strings are byte slices); this is what
		// the decoder will consume, so byte length — not rune count — is correct.
		b.WriteString(strconv.Itoa(len(s))) // announce payload size
		b.WriteByte('#')                    // separator between size and payload
		b.WriteString(s)                    // the raw payload, unescaped
	}
	return b.String()
}

func lengthPrefixDecode(encoded string) []string {
	res := []string{}
	i := 0
	for i < len(encoded) {
		// Scan forward to the '#' that terminates the length header.
		j := i
		for encoded[j] != '#' {
			j++
		}
		// encoded[i:j] is the ASCII length; parse it.
		length, _ := strconv.Atoi(encoded[i:j])
		// Payload starts right after '#' (at j+1) and is exactly `length` bytes.
		start := j + 1
		res = append(res, encoded[start:start+length])
		// Jump past this whole chunk to the next length header.
		i = start + length
	}
	return res
}
```

### Dry Run

Encode `["Hello","World"]` → `"5#Hello5#World"`. Now decode `"5#Hello5#World"`:

| Step | i | scan to '#' at j | length L = encoded[i:j] | payload encoded[j+1:j+1+L] | res | next i |
|------|---|------------------|-------------------------|----------------------------|-----|--------|
| 1 | 0 | j=1 | `"5"` → 5 | `encoded[2:7]` = `"Hello"` | `["Hello"]` | 7 |
| 2 | 7 | j=8 | `"5"` → 5 | `encoded[9:14]` = `"World"` | `["Hello","World"]` | 14 |
| 3 | 14 | loop ends (i == len) | — | — | `["Hello","World"]` | — |

Result: `["Hello","World"]` ✔

---

## Approach 2 — Escaping a Delimiter

### Intuition

Keep a real delimiter between strings, but first *neutralise* any occurrence of it inside each payload. Introduce an escape marker `#`: rewrite every literal `#` in a payload as `#h`, then terminate each string with the sentinel `#:`. Because every payload `#` is now `#h`, a raw `#:` can only ever be a true boundary — the escaping guarantees the sentinel is unambiguous.

### Algorithm

1. **Encode:** for each `s`, replace every `#` with `#h`, append the escaped payload, then append the sentinel `#:`.
2. **Decode:** split the stream on `#:` and drop the trailing empty piece (there is one after the final sentinel). Un-escape each piece by turning `#h` back into `#`.

### Complexity

- **Time:** O(N) — one pass to escape, one pass to split/un-escape.
- **Space:** O(N) — escaping can at most double each payload; output is still linear.

### Code

```go
func escapeEncode(strs []string) string {
	var b strings.Builder
	for _, s := range strs {
		// Escape the escape-introducer '#' so no literal '#' can be mistaken
		// for the start of our sentinel.
		esc := strings.ReplaceAll(s, "#", "#h")
		b.WriteString(esc) // safe payload
		b.WriteString("#:") // sentinel: a '#' followed by ':' — never appears
		//                     inside an escaped payload (every payload '#' is "#h")
	}
	return b.String()
}

func escapeDecode(encoded string) []string {
	if encoded == "" {
		return []string{}
	}
	// Every original string ended with "#:"; splitting on it yields the pieces
	// plus one trailing "" (after the final sentinel), which we drop.
	parts := strings.Split(encoded, "#:")
	parts = parts[:len(parts)-1] // remove the trailing empty element
	res := make([]string, len(parts))
	for i, p := range parts {
		// Reverse the escaping: "#h" → "#".
		res[i] = strings.ReplaceAll(p, "#h", "#")
	}
	return res
}
```

### Dry Run

Encode `["Hello","World"]`:

| Step | s | escape `#`→`#h` | append + sentinel | stream so far |
|------|---|------------------|-------------------|---------------|
| 1 | `Hello` | `Hello` (no `#`) | `Hello#:` | `Hello#:` |
| 2 | `World` | `World` (no `#`) | `World#:` | `Hello#:World#:` |

Decode `"Hello#:World#:"`:

| Step | action | value |
|------|--------|-------|
| 1 | split on `#:` | `["Hello","World",""]` |
| 2 | drop trailing `""` | `["Hello","World"]` |
| 3 | un-escape `#h`→`#` each | `["Hello","World"]` |

Result: `["Hello","World"]` ✔

---

## Approach 3 — Non-ASCII Sentinel (Fragile)

### Intuition

If we *assume* input never contains some exotic character, we can just `Join`/`Split` on it. This is the tempting one-liner — and it is wrong the moment an adversary embeds the sentinel. Shown only for contrast: it violates the follow-up ("any possible set of characters").

### Algorithm

1. **Encode:** join the list with a rare Unicode sentinel rune (e.g. `␟`, U+241F).
2. **Decode:** split the string on that same rune.

### Complexity

- **Time:** O(N) — a single join / split.
- **Space:** O(N) — the joined string.

### Code

```go
func sentinelEncode(strs []string) string {
	// U+241F ("SYMBOL FOR UNIT SEPARATOR") stands in for a byte unlikely to
	// occur in normal text. This is the fragile assumption.
	return strings.Join(strs, "␟")
}

func sentinelDecode(encoded string) []string {
	if encoded == "" {
		return []string{}
	}
	return strings.Split(encoded, "␟")
}
```

### Dry Run

Encode `["Hello","World"]`:

| Step | action | value |
|------|--------|-------|
| 1 | `Join` with `␟` | `"Hello␟World"` |

Decode `"Hello␟World"`:

| Step | action | value |
|------|--------|-------|
| 1 | `Split` on `␟` | `["Hello","World"]` |

Result: `["Hello","World"]` ✔ — but if any payload contained `␟`, this splits it into extra pieces and corrupts the list.

---

## Key Takeaways

- **Length-prefixing beats delimiter-hunting.** When a payload can contain *any* byte, announce its size and consume exactly that many bytes — no character is ever special. This is the interview-correct answer and the one that satisfies the follow-up.
- **A self-describing format needs no reserved characters.** The header (`length#`) is unambiguous by construction; the payload is passed through verbatim.
- **Escaping works but is trickier:** you must escape both the delimiter *and* the escape character, or an adversary reconstructs your delimiter.
- **"Just pick a rare delimiter" is a bug, not a solution** — it only holds under an assumption the problem explicitly removes.
- The same length-prefix idea generalises to serializing trees, nested structures, and binary protocols.

---

## Related Problems

- LeetCode #297 — Serialize and Deserialize Binary Tree (length/marker-based codec)
- LeetCode #443 — String Compression (encode a run as count + char)
- LeetCode #394 — Decode String (parse a self-describing encoded string)
- LeetCode #535 — Encode and Decode TinyURL (design an invertible codec)
