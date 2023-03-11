package main

import (
	"fmt"
	"strings"
	"unicode"
)

// возвращает слово в нижнем регистре
func toLowerUnicode(word string) []rune {
	runes := []rune{}
	for _, ch := range word {
		runes = append(runes, unicode.ToLower(ch))
	}

	return runes
}

// проверяет, является ли слово анаграммой
func isAnagram(runeKey []rune, word string) bool {
	for _, ch := range runeKey {
		if !strings.ContainsRune(word, ch) {
			return false
		}
	}
	return true
}

func mapAnagrams(words []string) *map[string]*[]string {

	result := make(map[string]*[]string) // map[слово]{1, 2, 3 ... анаграммы}
	used := make(map[string]struct{}, 0) // использованные слова
	for _, word := range words {

		// использованы все слова
		if len(used) >= len(words) {
			break
		}

		// word уже использовано
		if _, ok := used[word]; ok {
			continue
		}

		used[word] = struct{}{} // занесение word в список использованных слов
		runeKey := toLowerUnicode(word)
		result[string(runeKey)] = &[]string{} // result[word] = {owrd, rowd, dorw, ...}
		for _, word := range words {

			if _, ok := used[word]; ok { // word уже использовано
				continue
			}

			if len(runeKey) != len([]rune(word)) { // word и runeKey разной длины
				continue
			}

			if isAnagram(runeKey, word) { // word и runeKey одной длины
				used[word] = struct{}{}
				*result[string(runeKey)] = append(*result[string(runeKey)], string(toLowerUnicode(word)))
			}
		}
	}

	return &result
}

var words = []string{"пятак", "пятка", "тяпка", "листок", "слиток", "столик"}

func main() {
	result := mapAnagrams(words)
	for key, val := range *result {
		fmt.Printf("%q: %v\n", key, *val)
	}
}
