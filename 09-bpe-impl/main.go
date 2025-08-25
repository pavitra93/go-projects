package main

import (
	"fmt"
	"sort"
	_ "sort"
)

// pretty-print bytes: turn spaces into ␠ so they’re visible
func vis(b []byte) string {
	s := ""
	for _, ch := range b {
		if ch == ' ' {
			s += "␠"
		} else if ch >= 32 && ch < 127 {
			s += string(rune(ch))
		} else {
			s += fmt.Sprintf("\\x%02X", ch)
		}
	}
	return s
}

type pairCount struct {
	pair  Pair
	count int
}

// returns pair counts sorted by count desc, then by (A,B) asc for deterministic ties
func (bpe *BPEEncoder) sortedPairCounts(tokens []int) []pairCount {
	m := bpe.getPairCounts(tokens)
	out := make([]pairCount, 0, len(m))
	for p, c := range m {
		out = append(out, pairCount{p, c})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].count != out[j].count {
			return out[i].count > out[j].count
		}
		if out[i].pair.A != out[j].pair.A {
			return out[i].pair.A < out[j].pair.A
		}
		return out[i].pair.B < out[j].pair.B
	})
	return out
}

// Pair represents a pair of adjacent tokens
type Pair struct {
	A, B int
}

// BPEEncoder holds the vocabulary and merge rules
type BPEEncoder struct {
	vocab  map[int][]byte // token ID -> byte sequence
	merges map[Pair]int   // pair -> new token ID
	nextID int            // next available token ID
}

// NewBPEEncoder initializes a BPE encoder with base byte vocabulary
func NewBPEEncoder() *BPEEncoder {
	encoder := &BPEEncoder{
		vocab:  make(map[int][]byte),
		merges: make(map[Pair]int),
		nextID: 256, // Start new tokens from 256
	}
	// Initialize vocabulary with all single bytes
	for i := 0; i < 256; i++ {
		encoder.vocab[i] = []byte{byte(i)}
	}
	return encoder
}

// Train trains the BPE encoder on a given text
func (bpe *BPEEncoder) Train(text []byte, numMerges int) {
	// start from raw bytes
	tokens := make([]int, len(text))
	for i, b := range text {
		tokens[i] = int(b)
	}

	fmt.Println("==== BPE TRAIN ====")
	fmt.Printf("Text: %q\n", string(text))
	fmt.Printf("Initial tokens (%d): %v\n", len(tokens), tokens)
	fmt.Printf("Initial decoded: %s\n\n", vis(text))

	for step := 1; step <= numMerges; step++ {
		counts := bpe.sortedPairCounts(tokens)
		if len(counts) == 0 {
			fmt.Println("No more pairs to merge.")
			break
		}

		// show top-N pairs for context
		topN := 5
		if len(counts) < topN {
			topN = len(counts)
		}
		fmt.Printf("Step %d: top pairs:\n", step)
		for i := 0; i < topN; i++ {
			A := bpe.vocab[counts[i].pair.A]
			B := bpe.vocab[counts[i].pair.B]
			fmt.Printf("  %d) '%s' + '%s'  → count=%d\n", i+1, vis(A), vis(B), counts[i].count)
		}

		// choose the most frequent (deterministic due to sorting)
		best := counts[0].pair
		newID := bpe.nextID
		bpe.nextID++

		// create merged token
		bpe.merges[best] = newID
		merged := append(append([]byte{}, bpe.vocab[best.A]...), bpe.vocab[best.B]...)
		bpe.vocab[newID] = merged

		fmt.Printf("Chosen merge: '%s' + '%s'  ==>  id=%d '%s'\n",
			vis(bpe.vocab[best.A]), vis(bpe.vocab[best.B]), newID, vis(merged))

		// apply merge to stream
		tokens = bpe.mergeTokens(tokens, best, newID)

		// show new token stream (both ids and decoded)
		decoded := bpe.Decode(tokens)
		fmt.Printf("Tokens after merge (%d): %v\n", len(tokens), tokens)
		fmt.Printf("Decoded after merge: %s\n\n", vis(decoded))
	}

	fmt.Println("==== END TRAIN ====\n")
}

// getPairCounts counts occurrences of adjacent pairs
func (bpe *BPEEncoder) getPairCounts(tokens []int) map[Pair]int {
	counts := make(map[Pair]int)
	for i := 0; i < len(tokens)-1; i++ {
		pair := Pair{tokens[i], tokens[i+1]}
		counts[pair]++
	}
	return counts
}

// findMostCommonPair finds the pair with the highest frequency
func (bpe *BPEEncoder) findMostCommonPair(counts map[Pair]int) Pair {
	var mostCommon Pair
	maxCount := -1
	for pair, count := range counts {
		if count > maxCount {
			maxCount = count
			mostCommon = pair
		}
	}
	return mostCommon
}

// mergeTokens replaces occurrences of a pair with a new token ID
func (bpe *BPEEncoder) mergeTokens(tokens []int, pair Pair, newID int) []int {
	var newTokens []int
	for i := 0; i < len(tokens); {
		if i+1 < len(tokens) && tokens[i] == pair.A && tokens[i+1] == pair.B {
			newTokens = append(newTokens, newID)
			i += 2
		} else {
			newTokens = append(newTokens, tokens[i])
			i++
		}
	}
	return newTokens
}

// Encode encodes a byte slice into a slice of token IDs
func (bpe *BPEEncoder) Encode(data []byte) []int {
	tokens := make([]int, len(data))
	for i, b := range data {
		tokens[i] = int(b)
	}

	// Apply merges based on learned rules
	for {
		changed := false
		newTokens := []int{}
		for i := 0; i < len(tokens); {
			if i+1 < len(tokens) {
				currentPair := Pair{tokens[i], tokens[i+1]}
				newID, ok := bpe.merges[currentPair]
				if ok {
					newTokens = append(newTokens, newID)
					i += 2
					changed = true
					continue
				}
			}
			newTokens = append(newTokens, tokens[i])
			i++
		}
		tokens = newTokens
		if !changed {
			break
		}
	}
	return tokens
}

// Decode decodes a slice of token IDs back into a byte slice
func (bpe *BPEEncoder) Decode(tokenIDs []int) []byte {
	var decodedBytes []byte
	for _, id := range tokenIDs {
		if b, ok := bpe.vocab[id]; ok {
			decodedBytes = append(decodedBytes, b...)
		} else {
			// Handle unknown token IDs if necessary (e.g., skip or error)
			fmt.Printf("Warning: Unknown token ID %d during decoding.\n", id)
		}
	}
	return decodedBytes
}

func main() {
	encoder := NewBPEEncoder()
	//fmt.Println(encoder.vocab)
	text := []byte("this is a test text for byte pair encoding encoding")
	encoder.Train(text, 10) // Train with 10 merges

	encoded := encoder.Encode(text)
	fmt.Printf("Encoded tokens: %v\n", encoded)

	decoded := encoder.Decode(encoded)
	fmt.Printf("Decoded text: %s\n", string(decoded))
	//fmt.Println(encoder.vocab)
}
