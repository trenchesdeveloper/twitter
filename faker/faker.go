package faker

import (
	"fmt"
	"math/rand"
)

//func init() {
//	r := rand.New(rand.NewSource(time.Now().UnixNano()))
//
//	r.Uint64()
//}

var Password = "$2a$04$O8sswZMTnqy.iqYhnzAvpuun9asG.EvWqxnooLJgT.YTdcWJ.blNC"

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)

	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	return string(b)
}

func RandStringRunesWithLowercase(n int) string {
	b := make([]rune, n)

	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes)/2)]
	}

	return string(b)
}

func RandInt(min, max int) int {
	return rand.Intn(max-min+1) + min
}

func Username() string {
	return RandStringRunes(RandInt(2, 10))
}

func Email() string {
	return RandStringRunesWithLowercase(RandInt(2, 10)) + "@gmail.com"
}

func ID() string {
	return fmt.Sprintf("%s-%s-%s-%s", RandStringRunes(4), RandStringRunes(4), RandStringRunes(4), RandStringRunes(4), RandStringRunes(12))
}
