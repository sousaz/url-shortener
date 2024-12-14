package utils

import "math/rand/v2"

func Generate_shortener(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	res := make([]byte, n)
	for i := range res {
		res[i] = letters[rand.IntN(len(letters))]
	}
	return string(res)
}
