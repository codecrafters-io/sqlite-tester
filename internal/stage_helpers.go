package internal

import (
	"math/rand"
	"strings"
	"time"
)

var randomWords = []string{
	"humpty",
	"dumpty",
	"horsey",
	"donkey",
	"yikes",
	"monkey",
	"doo",
	"scooby",
	"dooby",
	"vanilla",
}

func randomString() string {
	rand.Seed(time.Now().UnixNano())

	return strings.Join(
		[]string{
			randomWords[rand.Intn(10)],
			randomWords[rand.Intn(10)],
			randomWords[rand.Intn(10)],
			randomWords[rand.Intn(10)],
			randomWords[rand.Intn(10)],
			randomWords[rand.Intn(10)],
		},
		" ",
	)
}

func randomStringShort() string {
	rand.Seed(time.Now().UnixNano())
	return randomWords[rand.Intn(10)]
}

func randomStringsShort(n int) []string {
	return shuffle(randomWords)[0:n]
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
