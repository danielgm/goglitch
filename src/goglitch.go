package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math/rand"
	"os"
	"time"

	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"sync"
)

var wg sync.WaitGroup

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Results are much more interesting if GOMAXPROCS != 1.
func main() {
	numGoRoutines := 2

	ch := make(chan uint32)

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

	wg.Add(numGoRoutines * 2)

	for i := 0; i < numGoRoutines; i++ {
		go reader(ch, inImg, b)
		go writer(ch, outImg, b)
	}

	wg.Wait()

	outFile, _ := os.Create("../output.png")
	defer outFile.Close()

	png.Encode(outFile, outImg)
}

func reader(palette chan uint32, inImg image.Image, b image.Rectangle) {
	defer wg.Done()

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			r, g, b, _ := inImg.At(x, y).RGBA()
			palette <- r
			palette <- g
			palette <- b
		}
	}

}

func writer(palette chan uint32, outImg image.Image, b image.Rectangle) {
	defer wg.Done()

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			r := <-palette
			g := <-palette
			b := <-palette
			value := color.RGBA{
				uint8(r / 255),
				uint8(g / 255),
				uint8(b / 255),
				255,
			}
			outImg.(draw.Image).Set(x, y, value)
		}
	}
}

func processColor(c color.Color, index int) color.Color {
	r, g, b, _ := c.RGBA()
	return color.RGBA{
		uint8(r / 255),
		uint8(g / 255),
		uint8(b / 255),
		255,
	}
}

func printRgba(c color.Color) {
	r, g, b, a := c.RGBA()
	fmt.Printf("%d, %d, %d, %d\n", r/255, g/255, b/255, a/255)
}
