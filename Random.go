package goblackholes

import (
	"math"
	"time"
)

var (
	x_1 = uint64(time.Now().UnixNano())
	x_2 = uint64(time.Now().UnixNano() / 2)
	x_3 = uint64(time.Now().UnixNano() / 3)
	mod = uint64(math.Pow(2, 32)) - 5
)

//func init() {
//	x_1 = uint64(time.Now().UnixNano())
//	x_2 = uint64(time.Now().UnixNano() / 2)
//	x_3 = uint64(time.Now().UnixNano() / 3)
//	mod = uint64(math.Pow(2, 32)) - 5
//}

func count() (x uint64) {
	x = (1176*x_1 + 1476*x_2 + 1776*x_3) % mod
	x_3, x_2, x_1 = x_2, x_1, x
	//x_2 = x_1
	//x_1 = x
	return
}

func NextInt64() (x uint64) {
	return count()
}

func NextDouble() float64 {
	return float64(count()) / float64(mod)
}

//func main() {
//	fmt.Println("I'll be counting random numbers.")
//
//	//fmt.Println(x_1, x_2, x_3, mod)
//
//	for i := 0; i < 100; i++ {
//		fmt.Println(NextDouble())
//		//fmt.Println(count(), float64(count()), float64(mod))
//	}
//
//}
