package internal

import (
	"math/rand"
	"strings"
	"time"
)

var randomWords = []string{
	"chocolate",
	"vanilla",
	"strawberry",
	"watermelon",
	"coffee",
	"grape",
	"coconut",
	"butterscotch",
	"pistachio",
	"banana",
	"mango",
}

func randomString() string {
	return strings.Join(
		[]string{
			randomStringShort(),
			randomStringShort(),
			randomStringShort(),
			randomStringShort(),
			randomStringShort(),
			randomStringShort(),
		},
		" ",
	)
}

func randomStringShort() string {
	rand.Seed(time.Now().UnixNano())
	return randomWords[rand.Intn(len(randomWords))]
}

func randomStringsShort(n int) []string {
	return shuffle(randomWords)[0:n]
}

func randomInt(n int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(n)
}

func shuffle(vals []string) []string {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	ret := make([]string, len(vals))
	perm := r.Perm(len(vals))
	for i, randIndex := range perm {
		ret[i] = vals[randIndex]
	}
	return ret
}
