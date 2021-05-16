package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

const letterByte = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringByte(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterByte[rand.Intn(len(letterByte))]
	}
	return string(b)
}
func main() {

	firstweek, _ := time.Parse("20060102", "20200217")
	c := time.Now()
	fmt.Print(math.Floor(c.Sub(firstweek).Hours()/(24*7)) + 1)

}
