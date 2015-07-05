// img_test
package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	"math"
	"os"
	"sort"
	"strconv"
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

	newImg := convertMSPaint(img)

	saveImage(newImg)

}
func convertMono(img image.Image) image.Image {

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

	avg1 := uint8(arr[int(len(arr)/3)])
	avg2 := uint8(arr[int(len(arr)/3)*2])

	for y := 0; y < rect.Size().Y; y++ {
		for x := 0; x < rect.Size().X; x++ {
			r0, g0, b0, _ := img.At(x, y).RGBA()
			r, g, b := uint8(r0), uint8(g0), uint8(b0)
			if (r+g+b)/3 > avg2 {
				r, g, b = 255, 255, 255
			} else if (r+g+b)/3 > avg1 {
				r, g, b = 128, 128, 128
			} else {
				r, g, b = 0, 0, 0
			}

			rgba.Set(x, y, color.RGBA{r, g, b, 255})

		}
	}
	return rgba
}
func convertMSPaint(img image.Image) image.Image {
	rect := img.Bounds()
	rgba := image.NewRGBA(rect)
	colorSet := [...]string{
		"000000", "FFFFFF", "808080", "c0c0c0",
		"800000", "FF0000", "808000", "FFFF00",
		"008000", "00FF00", "008080", "00FFFF",
		"000080", "0000FF", "800080", "FF00FF",
	}
	for y := 0; y < rect.Size().Y; y++ {
		for x := 0; x < rect.Size().X; x++ {
			r0, g0, b0, _ := img.At(x, y).RGBA()
			r, g, b := uint8(r0), uint8(g0), uint8(b0)

			r, g, b = nearColor(r, g, b, colorSet[:])
			rgba.Set(x, y, color.RGBA{r, g, b, 255})

		}
	}
	return rgba

}

func nearColor(r0, g0, b0 uint8, conv_colors []string) (uint8, uint8, uint8) {

	sel_r, sel_g, sel_b := uint8(0), uint8(0), uint8(9)
	sel_d := 999999.0

	for i := 0; i < len(conv_colors); i++ {
		rx, _ := (strconv.ParseUint(conv_colors[i][0:2], 16, 0))
		gx, _ := (strconv.ParseUint(conv_colors[i][2:4], 16, 0))
		bx, _ := (strconv.ParseUint(conv_colors[i][4:6], 16, 0))
		r, g, b := int(rx), int(gx), int(bx)
		rd := math.Pow(float64(int(r0)-r), 2)
		gd := math.Pow(float64(int(g0)-g), 2)
		bd := math.Pow(float64(int(b0)-b), 2)

		d := math.Sqrt(rd + gd + bd)

		if d <= sel_d {
			sel_r, sel_g, sel_b = uint8(r), uint8(g), uint8(b)
			sel_d = d
		}
	}
	return sel_r, sel_g, sel_b

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
