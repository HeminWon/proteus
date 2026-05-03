package core

import "sort"

func levenshtein(a, b string) int {
	ar := []rune(a)
	br := []rune(b)

	if len(ar) == 0 {
		return len(br)
	}
	if len(br) == 0 {
		return len(ar)
	}

	dp := make([][]int, len(ar)+1)
	for i := range dp {
		dp[i] = make([]int, len(br)+1)
	}

	for i := 0; i <= len(ar); i++ {
		dp[i][0] = i
	}
	for j := 0; j <= len(br); j++ {
		dp[0][j] = j
	}

	for i := 1; i <= len(ar); i++ {
		for j := 1; j <= len(br); j++ {
			cost := 0
			if ar[i-1] != br[j-1] {
				cost = 1
			}
			del := dp[i-1][j] + 1
			ins := dp[i][j-1] + 1
			sub := dp[i-1][j-1] + cost
			dp[i][j] = minInt(del, ins, sub)
		}
	}

	return dp[len(ar)][len(br)]
}

func minInt(values ...int) int {
	best := values[0]
	for _, v := range values[1:] {
		if v < best {
			best = v
		}
	}
	return best
}

func bestSuggestion(input string, candidates []string, threshold int) string {
	if input == "" || len(candidates) == 0 {
		return ""
	}

	type scored struct {
		candidate string
		distance  int
	}

	scoredCandidates := make([]scored, 0, len(candidates))
	for _, candidate := range candidates {
		d := levenshtein(input, candidate)
		if d <= threshold {
			scoredCandidates = append(scoredCandidates, scored{candidate: candidate, distance: d})
		}
	}

	if len(scoredCandidates) == 0 {
		return ""
	}

	sort.Slice(scoredCandidates, func(i, j int) bool {
		if scoredCandidates[i].distance == scoredCandidates[j].distance {
			return scoredCandidates[i].candidate < scoredCandidates[j].candidate
		}
		return scoredCandidates[i].distance < scoredCandidates[j].distance
	})

	return scoredCandidates[0].candidate
}

func SuggestCommand(input string, candidates []string) string {
	s := bestSuggestion(input, candidates, 3)
	if s == "" {
		return ""
	}
	return ". Did you mean `" + s + "`?"
}

func SuggestFlag(input string, candidates []string) string {
	s := bestSuggestion(input, candidates, 3)
	if s == "" {
		return ""
	}
	return ". Did you mean `" + s + "`?"
}

func SuggestProvider(input string, candidates []string) string {
	s := bestSuggestion(input, candidates, 3)
	if s == "" {
		return ""
	}
	return ". Did you mean `" + s + "`?"
}

func SuggestProfile(input string, candidates []string) string {
	s := bestSuggestion(input, candidates, 3)
	if s == "" {
		return ""
	}
	return ". Did you mean `" + s + "`?"
}
