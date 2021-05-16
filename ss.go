package main

import (
	"fmt"
	"log"
	"math/big"
	"strconv"
	"strings"
	"time"
)

type Number = *big.Int

var (
	Num0 = big.NewInt(0)
	Num1 = big.NewInt(1)
)

func NumAdd(x Number, y Number) Number {
	return big.NewInt(0).Add(x, y)
}

func NumSub(x Number, y Number) Number {
	return big.NewInt(0).Sub(x, y)
}

func NumMul(x Number, y Number) Number {
	return big.NewInt(0).Mul(x, y)
}

type Matrix = [2][2]Number

var (
	Mat0 = Matrix{{Num1, Num0}, {Num0, Num1}}
	Mat1 = Matrix{{Num1, Num1}, {Num1, Num0}}
)

func MatAdd(x Matrix, y Matrix) (z Matrix) {
	z[0][0] = NumAdd(x[0][0], y[0][0])
	z[0][1] = NumAdd(x[0][1], y[0][1])
	z[1][0] = NumAdd(x[1][0], y[1][0])
	z[1][1] = NumAdd(x[1][1], y[1][1])
	return
}

func MatSub(x Matrix, y Matrix) (z Matrix) {
	z[0][0] = NumSub(x[0][0], y[0][0])
	z[0][1] = NumSub(x[0][1], y[0][1])
	z[1][0] = NumSub(x[1][0], y[1][0])
	z[1][1] = NumSub(x[1][1], y[1][1])
	return
}

func MatMul(x Matrix, y Matrix) (z Matrix) {
	z[0][0] = NumAdd(NumMul(x[0][0], y[0][0]), NumMul(x[0][1], y[1][0]))
	z[0][1] = NumAdd(NumMul(x[0][0], y[0][1]), NumMul(x[0][1], y[1][1]))
	z[1][0] = NumAdd(NumMul(x[1][0], y[0][0]), NumMul(x[1][1], y[1][0]))
	z[1][1] = NumAdd(NumMul(x[1][0], y[0][1]), NumMul(x[1][1], y[1][1]))
	return
}

func logElapsed(name string) func() {
	start := time.Now()
	return func() {
		log.Printf("%s: %v", name, time.Since(start).Seconds())
	}
}
func main() {
	fmt.Printf("a = %v ", fibV3Iterative(50000000))
}
func fibV3Iterative(n int) Number {
	defer logElapsed("fibV3Iterative")()
	x, y := Num1, Num0
	bs := strconv.FormatInt(int64(n), 2)
	bs = strings.TrimLeft(bs, "0")
	prev := '0'
	for _, curr := range bs {
		xx := NumMul(x, x)
		yy := NumMul(y, y)
		xy := NumSub(xx, yy)
		if prev != '0' {
			xy = NumAdd(xy, Num1)
		} else {
			xy = NumSub(xy, Num1)
		}
		if curr != '0' {
			x = NumAdd(NumAdd(xx, xy), xy)
			y = NumAdd(xx, yy)
		} else {
			x = NumAdd(xx, yy)
			y = NumAdd(xy, NumSub(xy, yy))
		}
		prev = curr
	}
	return y
}
