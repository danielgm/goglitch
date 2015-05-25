package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"math"
	"math/rand"
	"os"
	"time"

	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"sync"
)

var numGoRoutines int

var wg sync.WaitGroup

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	numGoRoutines = 4

	fmt.Println("Startup.")

	inFile, _ := os.Open("../weyland.png")
	defer inFile.Close()

	inImg, format, err := image.Decode(bufio.NewReader(inFile))
	if err != nil {
		fmt.Println("Error: %v\n", err)
		return
	}
	fmt.Printf("Format: %s\n", format)
	b := inImg.Bounds()

	outImg := image.NewRGBA(b)

	wg.Add(numGoRoutines)

	for i := 0; i < numGoRoutines; i++ {
		go func(index int) {
			for y := b.Min.Y; y < b.Max.Y; y++ {
				for x := b.Min.X; x < b.Max.X; x++ {
					outImg.Set(x, y, processColor(inImg.At(x, y), index))
				}
			}

			wg.Done()
		}(i)
	}

	wg.Wait()

	outFile, _ := os.Create("../output.png")
	defer outFile.Close()

	png.Encode(outFile, outImg)
}

func processColor(c color.Color, index int) color.Color {
	_, g, b, _ := c.RGBA()
	return color.RGBA{
		uint8(math.Floor(float64(index) / float64(numGoRoutines) * 255)),
		uint8(g / 255),
		uint8(b / 255),
		255,
	}
}

func printRgba(c color.Color) {
	r, g, b, a := c.RGBA()
	fmt.Printf("%d, %d, %d, %d\n", r/255, g/255, b/255, a/255)
}
