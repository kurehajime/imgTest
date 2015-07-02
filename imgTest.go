// img_test
package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	"os"
	"sort"
)

func main() {
	file, err := os.Open("image.png")
	defer file.Close()

	if err != nil {
		fmt.Println(err)
		return
	}
	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Println(err)
		return
	}

	newImg := changeImage(img)

	saveImage(newImg)

}
func changeImage(img image.Image) image.Image {

	rect := img.Bounds()
	rgba := image.NewRGBA(rect)

	arr := make([]int, 4)

	for y := 0; y < rect.Size().Y; y++ {
		for x := 0; x < rect.Size().X; x++ {
			r0, g0, b0, _ := img.At(x, y).RGBA()
			r, g, b := uint8(r0), uint8(g0), uint8(b0)
			arr = append(arr, int((r+g+b)/3))
		}
	}

	sort.Ints(arr)

	avg := uint8(arr[int(len(arr)/2)])

	for y := 0; y < rect.Size().Y; y++ {
		for x := 0; x < rect.Size().X; x++ {
			r0, g0, b0, _ := img.At(x, y).RGBA()
			r, g, b := uint8(r0), uint8(g0), uint8(b0)
			if (r+g+b)/3 > avg {
				r, g, b = 255, 255, 255
			} else {
				r, g, b = 0, 0, 0
			}

			rgba.Set(x, y, color.RGBA{r, g, b, 255})

		}
	}
	return rgba
}

func saveImage(img image.Image) {
	out, err := os.Create("output.png")
	defer out.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = png.Encode(out, img)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
