package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"os"

	"golang.org/x/image/draw"
)

func getImage(pic string, teaWidth int) image.Image {

	path := "./Assets/" + pic

	reader, err := os.Open(path)
	if err != nil {
		fmt.Println("os error")
	}

	img, _, err := image.Decode(reader)
	if err != nil {
		fmt.Println("Deocde Error")
	}

	return resizeImage(img, teaWidth)
}

func resizeImage(img image.Image, targetWidth int) image.Image {
	bounds := img.Bounds()
	targetHeight := (bounds.Dy() * targetWidth) / bounds.Dx()
	dst := image.NewRGBA((image.Rect(0, 0, targetWidth, targetHeight)))
	draw.BiLinear.Scale(dst, dst.Bounds(), img, bounds, draw.Over, nil)
	return dst
}

func renderANSIHalfBlocks(product string, teaWidth int) string {

	img := getImage(product, teaWidth)

	bounds := img.Bounds()
	var buf bytes.Buffer

	// Step by 2 vertically because each block represents two vertical pixels
	for y := bounds.Min.Y; y < bounds.Max.Y; y += 2 {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			top := img.At(x, y)
			var bottom color.Color = color.RGBA{0, 0, 0, 0}
			if y+1 < bounds.Max.Y {
				bottom = img.At(x, y+1)
			}

			tr, tg, tb, _ := top.RGBA()
			br, bg, bb, _ := bottom.RGBA()

			// 8-bit color shift (Go's RGBA returns 0-65535)
			// 38;2;R;G;B = foreground (top pixel)
			// 48;2;R;G;B = background (bottom pixel)
			fmt.Fprintf(&buf, "\x1b[38;2;%d;%d;%dm\x1b[48;2;%d;%d;%dm▀\x1b[0m",
				tr>>8, tg>>8, tb>>8,
				br>>8, bg>>8, bb>>8,
			)
		}
		buf.WriteByte('\n')
	}

	return buf.String()
}
