package util

import "strings"

func SplitWordsEvery(s string, chunkSize int) []string {
	if len(s) == 0 {
		return nil
	}

	var chunks []string

	for len(s) > 0 {
		// If remaining string is smaller than chunk, take the rest
		if len(s) <= chunkSize {
			chunks = append(chunks, s)
			break
		}

		// Look at the character at the boundary
		end := chunkSize

		// Backtrack to find the last space
		for end > 0 && s[end] != ' ' {
			end--
		}

		// If no space was found (a single word > chunkSize), hard cut at chunkSize
		if end == 0 {
			end = chunkSize
		}

		chunks = append(chunks, strings.TrimSpace(s[:end]))

		// Advance the string slice, trimming leading spaces for the next iteration
		s = strings.TrimSpace(s[end:])
	}

	return chunks
}
