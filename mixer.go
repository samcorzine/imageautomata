package main

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

type weightedAverage struct {
	up, down, left, right, here float64
}

func weightColor(theWeight float64, theColor color.Color) color.Color {
	r, g, b, _ := theColor.RGBA()
	newr := uint16(theWeight * float64(r))
	newg := uint16(theWeight * float64(g))
	newb := uint16(theWeight * float64(b))
	newColor := color.NRGBA64{newr, newg, newb, 65535}
	return newColor
}

func sumColors(colors []color.Color) color.Color {
	var newr uint32
	var newg uint32
	var newb uint32
	var newa uint32
	for _, col := range colors {
		theR, theG, theB, _ := col.RGBA()
		newr += theR
		newg += theG
		newb += theB
		newa = 65535
	}
	return color.NRGBA64{uint16(newr), uint16(newg), uint16(newb), uint16(newa)}
}

type colorAndCoordinates struct {
	x, y int
	col  color.Color
}

func (weights weightedAverage) blend(theImage image.Image) image.Image {
	theBounds := theImage.Bounds()
	xMax := theBounds.Max.X
	yMax := theBounds.Max.Y
	xMin := theBounds.Min.X
	yMin := theBounds.Min.Y

	theNewRectangle := image.Rect(xMin, yMin, xMax, yMax)
	theNewImage := image.NewNRGBA64(theNewRectangle)
	for x := xMin; x < xMax; x += 1 {
		for y := yMin; y < yMax; y += 1 {
			var up, down, left, right, val color.Color
			if y+1 < yMax {
				up = theImage.At(x, y+1)
			} else {
				up = color.NRGBA64{8000, 0, 0, 50000}
			}
			if y-1 >= 0 {
				down = theImage.At(x, y-1)
			} else {
				down = color.NRGBA64{0, 8000, 0, 50000}
			}
			if x-1 >= 0 {
				left = theImage.At(x-1, y)
			} else {
				left = color.NRGBA64{0, 0, 8000, 50000}
			}
			if x+1 < xMax {
				right = theImage.At(x+1, y)
			} else {
				right = color.NRGBA64{0, 0, 0, 0}
			}
			val = theImage.At(x, y)
			modUp := weightColor(weights.up, up)
			modDown := weightColor(weights.down, down)
			modLeft := weightColor(weights.left, left)
			modRight := weightColor(weights.right, right)
			modHere := weightColor(weights.here, val)
			newVal := sumColors([]color.Color{modUp, modDown, modLeft, modRight, modHere})
			theNewImage.Set(x, y, newVal)
		}
	}
	return theNewImage
}

func main() {
	theFile, _ := os.Open("test.png")
	defer theFile.Close()
	testImage, _, _ := image.Decode(theFile)
	theWeights := weightedAverage{0.2, 0.2, 0.2, 0.2, 0.2}
	blendedImage := testImage
	for i := 0; i < 100; i += 1 {
		blendedImage = theWeights.blend(blendedImage)
	}

	os.Remove("result.png")
	newFile, _ := os.Create("result.png")
	png.Encode(newFile, blendedImage)
}
