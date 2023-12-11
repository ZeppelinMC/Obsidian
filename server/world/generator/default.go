package generator

import (
	"fmt"
	"math"
	"obsidian/server/world/block"
	"strings"

	"github.com/aquilax/go-perlin"
)

var p = perlin.NewPerlin(2, 2, 2, 123)

type DefaultGenerator struct{}

func (d *DefaultGenerator) GenerateWorld(sizeX, sizeY, sizeZ int16) (blocks []int8) {
	blocks = make([]int8, int(sizeX)*int(sizeY)*int(sizeZ))

	var c int
	for x := int16(0); x < sizeX; x++ {
		for z := int16(0); z < sizeZ; z++ {
			y := int16(p.Noise2D(float64(x)/50, float64(z)/50) * 25)

			if y <= 0 {
				y = 0
			} else if y >= sizeY {
				y = sizeY
			}

			i := int(x) + int(sizeX)*(int(z)+int(sizeZ)*int(y))
			blocks[i] = block.BlockGrass

			l := len(fmt.Sprint(per(c-1, len(blocks)/int(sizeY)))) + 1

			fmt.Printf("%s%d%%", strings.Repeat("\b", l), per(c, len(blocks)/int(sizeY)))

			if c == (len(blocks)-1)/int(sizeY) {
				fmt.Println()
			}

			c++
		}
	}
	return
}

func per(c, i int) int {
	return int(math.Ceil(float64(c) / float64(i) * 100))
}
