package main

import "strings"

func validateBadWords(body string) string {
	lowerBody := strings.ToLower(body)
	splitWords := strings.Split(lowerBody, " ")
	for i, word := range splitWords {
		if word == "kerfuffle" || word == "sharbert" || word == "fornax" {
			splitWords[i] = "****"
		}
	}
	return strings.Join(splitWords, " ")
}
