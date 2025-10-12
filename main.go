package main

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	"image/color"
	"log"
	"os"

	"github.com/golang/freetype"

	// "image/draw"
	"image/png"
)

const imageSize = 1000
const imageSizeThird = imageSize / 3
const imageSizeTwoThird = imageSize / 3 * 2

const firstLinePosition = imageSizeThird / 3
const secondLinePosition = imageSizeTwoThird / 3

const cellSize = imageSize / 9

var white = color.RGBA{255, 255, 255, 255}
var black = color.RGBA{0, 0, 0, 255}
var gray = color.RGBA{0, 0, 0, 150}

func main() {
	f, err := os.Open("puzzle.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	file := bufio.NewReader(f)
	data, err := file.Peek(81)
	if err != nil {
		panic(err)
	}
	fmt.Printf("81 bytes: %s\n", string(data))

	OutputPuzzle(data)
}

func OutputPuzzle(puzzle []byte) {
	var output []byte
	counter := 1
	separator := "-------+-------+-------"
	for _, character := range puzzle {
		if (counter-1)%9 == 0 {
			output = append(output, ' ')
		}
		output = append(output, character, ' ')
		switch {
		case counter%81 == 0:
			output = append(output, '\n')
		case counter%27 == 0:
			output = append(output, '\n')
			output = append(output, []byte(separator)...)
			output = append(output, '\n')
		case counter%9 == 0:
			output = append(output, '\n')
		case counter%3 == 0:
			output = append(output, '|', ' ')
		}
		counter++
	}

	buf := bytes.NewBuffer(output)
	fmt.Println(buf.String())
}

func CreateImage() {
	rect := image.Rect(0, 0, imageSize, imageSize)
	baseImage := image.NewRGBA(rect)
	imageFile, err := os.Create("output.png")
	if err != nil {
		panic(err)
	}
	defer imageFile.Close()

	for y := range imageSize {
		for x := range imageSize {
			baseImage.Set(x, y, white)
		}
	}

	lines := []struct {
		Color color.RGBA
		X     int
	}{
		{
			Color: black,
			X:     imageSizeThird,
		},
		{
			Color: black,
			X:     imageSizeTwoThird,
		},
		{
			Color: gray,
			X:     firstLinePosition,
		},
		{
			Color: gray,
			X:     secondLinePosition,
		},
		{
			Color: gray,
			X:     firstLinePosition + imageSizeThird,
		},
		{
			Color: gray,
			X:     secondLinePosition + imageSizeThird,
		},
		{
			Color: gray,
			X:     firstLinePosition + imageSizeTwoThird,
		},
		{
			Color: gray,
			X:     secondLinePosition + imageSizeTwoThird,
		},
	}

	for _, line := range lines {
		for i := 0; i < imageSize; i++ {
			baseImage.Set(i, line.X, line.Color)
			baseImage.Set(line.X, i, line.Color)
		}
	}

	fontBytes, err := os.ReadFile("roboto.ttf")
	if err != nil {
		panic(err)
	}

	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Println(err)
		return
	}

	c := freetype.NewContext()
	c.SetDPI(100)
	c.SetFont(font)
	c.SetFontSize(100)
	c.SetClip(rect)
	c.SetDst(baseImage)
	c.SetSrc(image.Black)

	addLabel(c, 1, 1, "0")

	if err := png.Encode(imageFile, baseImage); err != nil {
		panic(err)
	}
}

func GetNumberPosition(x int, y int) (int, int) {
	if x > 9 || y > 9 || x <= 0 || y <= 0 {
		log.Fatalf("border exceeded, y:%d x:%d", y, x)
	}

	return (x-1)*imageSizeThird + 10, (y-1)*imageSizeThird + cellSize - 5
}

func addLabel(c *freetype.Context, x, y int, label string) {
	x, y = GetNumberPosition(x, y)
	pt := freetype.Pt(x, y)

	if _, err := c.DrawString(label, pt); err != nil {
		panic(err)
	}
}
