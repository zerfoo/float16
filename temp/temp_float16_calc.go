package main

import (
	"fmt"
	"github.com/zerfoo/float16"
	"math"
)

func main() {
	piDiv2 := float32(math.Pi / 2.0)
	f16PiDiv2 := float16.ToFloat16(piDiv2)
	fmt.Printf("Float16(math.Pi / 2.0) = 0x%04X (%.10f)\n", f16PiDiv2.Bits(), f16PiDiv2.ToFloat32())

	negPiDiv2 := float32(-(math.Pi / 2.0))
	f16NegPiDiv2 := float16.ToFloat16(negPiDiv2)
	fmt.Printf("Float16(-(math.Pi / 2.0)) = 0x%04X (%.10f)\n", f16NegPiDiv2.Bits(), f16NegPiDiv2.ToFloat32())
}
