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
	"time"
)

func main() {
	path := ""
	start := time.Now()
	if len(os.Args) >= 2 {
		path = os.Args[1]
	} else {
		fmt.Println("Please specify the path to image file.")
		return
	}

	img := getIMG(path)

	newImg := convertMSPaint(img)
	newImg = replaceTexture(newImg, getIMG("FF0000.png"), "FF0000")

	saveImage(newImg)
	end := time.Now()
	fmt.Printf("complete! (%fs)\n", (end.Sub(start)).Seconds())

}
func getIMG(path string) image.Image {
	file, err := os.Open(path)
	defer file.Close()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return img
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
	conv_colors := make([][3]uint8, len(colorSet))
	for i := 0; i < len(colorSet); i++ {
		r, _ := (strconv.ParseUint(colorSet[i][0:2], 16, 0))
		g, _ := (strconv.ParseUint(colorSet[i][2:4], 16, 0))
		b, _ := (strconv.ParseUint(colorSet[i][4:6], 16, 0))
		conv_colors[i] = [3]uint8{uint8(r), uint8(g), uint8(b)}
	}

	for y := 0; y < rect.Size().Y; y++ {
		for x := 0; x < rect.Size().X; x++ {
			r0, g0, b0, _ := img.At(x, y).RGBA()
			r, g, b := uint8(r0), uint8(g0), uint8(b0)
			r, g, b = nearColor(r, g, b, conv_colors)
			rgba.Set(x, y, color.RGBA{r, g, b, 255})

		}
	}
	return rgba

}
func replaceTexture(img, bg image.Image, colorCode string) image.Image {
	rect := img.Bounds()
	rect2 := bg.Bounds()
	r, _ := (strconv.ParseUint(colorCode[0:2], 16, 0))
	g, _ := (strconv.ParseUint(colorCode[2:4], 16, 0))
	b, _ := (strconv.ParseUint(colorCode[4:6], 16, 0))

	rgba := image.NewRGBA(rect)
	for y := 0; y < rect.Size().Y; y++ {
		for x := 0; x < rect.Size().X; x++ {
			r0, g0, b0, _ := img.At(x, y).RGBA()
			if uint8(r) == uint8(r0) && uint8(g) == uint8(g0) && uint8(b) == uint8(b0) {
				x1 := int(math.Mod(float64(x), float64(rect2.Size().X)))
				y1 := int(math.Mod(float64(y), float64(rect2.Size().Y)))
				r1, g1, b1, _ := bg.At(x1, y1).RGBA()
				rgba.Set(x, y, color.RGBA{uint8(r1), uint8(g1), uint8(b1), 255})
			} else {
				rgba.Set(x, y, color.RGBA{uint8(r0), uint8(g0), uint8(b0), 255})
			}
		}
	}
	return rgba

}

func nearColor(r0, g0, b0 uint8, conv_colors [][3]uint8) (uint8, uint8, uint8) {

	sel_r, sel_g, sel_b := uint8(0), uint8(0), uint8(0)
	sel_d := 999999.0

	for i := 0; i < len(conv_colors); i++ {
		rx := conv_colors[i][0]
		gx := conv_colors[i][1]
		bx := conv_colors[i][2]
		rd := math.Pow(float64(r0)-float64(rx), 2)
		gd := math.Pow(float64(g0)-float64(gx), 2)
		bd := math.Pow(float64(b0)-float64(bx), 2)

		d := math.Sqrt(rd + gd + bd)

		if d <= sel_d {
			sel_r, sel_g, sel_b = rx, gx, bx
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
