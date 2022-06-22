package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

type Words struct {
	Word  string
	Count int
}

func Top10(s string) []string {
	w := []Words{}
	wordMap := make(map[string]int)

	for _, word := range strings.Fields(s) {
		wordMap[word]++
	}

	for key, val := range wordMap {
		w = append(w, Words{key, val})
	}

	sort.Slice(w, func(i, j int) bool {
		return w[i].Count > w[j].Count || (w[i].Count == w[j].Count && (strings.Compare(w[i].Word, w[j].Word) < 0))
	})

	//sort.Slice(w, func(i, j int) bool {
	//	return w[i].Word < w[j].Word
	//})
	//
	//sort.SliceStable(w, func(i, j int) bool {
	//	return w[i].Count > w[j].Count
	//})

	res := []string{}

	for _, val := range w {
		res = append(res, val.Word)
	}

	if len(res) < 11 {
		return res
	}

	return res[:10]
}
