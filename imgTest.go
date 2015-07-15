package main
import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	"os"
	"strconv"
	"time"
)
var COLOR_SET = [...]string{
	"000000", "FFFFFF", //"808080", "c0c0c0",
	"800000", "FF0000", "808000", "FFFF00",
	"008000", "00FF00", "008080", "00FFFF",
	"000080", "0000FF", "800080", "FF00FF",
}
func main() {
	path := ""
	start := time.Now()
	if len(os.Args) >= 2 {
		path = os.Args[1]
	} else {
		fmt.Println("コマンドライン引数に画像ファイルを指定してください")
		return
	}
	img := getIMG(path)
	newImg := convertColor(img)
	saveImage(newImg)
	end := time.Now()
	fmt.Printf("complete! (%fs)\n", (end.Sub(start)).Seconds())
}
//画像を読み込む
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
//画像を保存する
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
//画像を変換する
func convertColor(img image.Image) image.Image {
	rect := img.Bounds()
	rgba := image.NewRGBA(rect)
	conv_colors := make([][3]uint8, len(COLOR_SET))
	for i := 0; i < len(COLOR_SET); i++ {
		r, _ := (strconv.ParseUint(COLOR_SET[i][0:2], 16, 0))
		g, _ := (strconv.ParseUint(COLOR_SET[i][2:4], 16, 0))
		b, _ := (strconv.ParseUint(COLOR_SET[i][4:6], 16, 0))
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
//近い色を返す
func nearColor(r0, g0, b0 uint8, conv_colors [][3]uint8) (uint8, uint8, uint8) {
	sel_r, sel_g, sel_b := uint8(0), uint8(0), uint8(0)
	sel_d := 999999.0
	for i := 0; i < len(conv_colors); i++ {
		rx := conv_colors[i][0]
		gx := conv_colors[i][1]
		bx := conv_colors[i][2]
		rd := (float64(r0) - float64(rx)) * (float64(r0) - float64(rx))
		gd := (float64(g0) - float64(gx)) * (float64(g0) - float64(gx))
		bd := (float64(b0) - float64(bx)) * (float64(b0) - float64(bx))
		d := (rd + gd + bd)
		if d <= sel_d {
			sel_r, sel_g, sel_b,sel_d = rx, gx, bx,d
		}
	}
	return sel_r, sel_g, sel_b
}