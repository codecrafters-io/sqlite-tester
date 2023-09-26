package internal

import (
	"math/rand"
	"os"
	"strconv"
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

func initRandom() {
	if seed := os.Getenv("CODECRAFTERS_RANDOM_SEED"); seed != "" {
		seedInt, err := strconv.Atoi(seed)

		if err != nil {
			panic(err)
		}

		rand.Seed(int64(seedInt))
	} else {
		rand.Seed(time.Now().UnixNano())
	}
}

func randomWord() string {
	return randomWords[rand.Intn(len(randomWords))]
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
	return randomWord()
}

func randomStringsShort(n int) []string {
	return shuffle(randomWords)[0:n]
}

func randomInt(n int) int {
	return rand.Intn(n)
}

func shuffle(vals []string) []string {
	ret := make([]string, len(vals))
	perm := rand.Perm(len(vals))
	for i, randIndex := range perm {
		ret[i] = vals[randIndex]
	}
	return ret
}
