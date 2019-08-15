package utils

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type Random struct {
	numeric [10]byte
	chars26 [26]byte
	chars36 [36]byte
}

func (random *Random) InitRandom() {
	random.numeric = [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	random.chars26 = [26]byte{'A', 'B', 'C', 'D', 'E', 'F', 'G',
		'H', 'I', 'J', 'K', 'L', 'M', 'N',
		'O', 'P', 'Q', 'R', 'S', 'T',
		'U', 'V', 'W', 'X', 'Y', 'Z'}
	random.chars36 = [36]byte{'A', 'B', 'C', 'D', 'E', 'F', 'G',
		'H', 'I', 'J', 'K', 'L', 'M', 'N',
		'O', 'P', 'Q', 'R', 'S', 'T',
		'U', 'V', 'W', 'X', 'Y', 'Z',
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}
}

func (random *Random) RandomNumber(width int) string {
	rand.Seed(time.Now().UnixNano())
	var sb strings.Builder
	for i := 0; i < width; i++ {
		fmt.Fprintf(&sb, "%d", random.numeric[rand.Intn(10)])
	}
	return sb.String()
}

func (random *Random) RandomString(width int) string {
	rand.Seed(time.Now().UnixNano())
	var sb strings.Builder
	for i := 0; i < width; i++ {
		fmt.Fprintf(&sb, "%c", random.chars26[rand.Intn(26)])
	}
	return sb.String()
}

func (random *Random) RandomStringAndNumber(width int) string {
	rand.Seed(time.Now().UnixNano())
	var sb strings.Builder
	for i := 0; i < width; i++ {
		fmt.Fprintf(&sb, "%c", random.chars36[rand.Intn(36)])
	}
	return sb.String()
}
