package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"math"
	"os"
	"os/exec"
)

var (
	white color.Color = color.RGBA{255, 255, 255, 255}
	black color.Color = color.RGBA{0, 0, 0, 255}
	blue  color.Color = color.RGBA{0, 0, 255, 255}
)

// ref) http://golang.org/doc/articles/image_draw.html
func loadImageObj() (img image.Image, err error) {
	filePath := "ic_launcher.png"
	fileSrc, err := os.Open(filePath)

	if err != nil {
		return nil, err
	}
	defer fileSrc.Close()

	img, err = png.Decode(fileSrc)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func buildCircle(d int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, d, d))

	radius := d / 2
	for i := 0.0; i < 360; i += 0.05 {
		x := float64(radius) + float64(radius)*math.Cos(i*math.Pi/180)
		y := float64(radius) + float64(radius)*math.Sin(i*math.Pi/180)

		x1 := float64(radius+1) + float64(radius+1)*math.Cos(i*math.Pi/180)
		y1 := float64(radius+1) + float64(radius+1)*math.Sin(i*math.Pi/180)
		img.Set(int(x), int(y), black)
		img.Set(int(x1), int(y1), black)
	}

	img.Set(int(d/2), int(d/2), black)

	return img
}

// func buildCircleMiddle(radius int) image.Image {
// 	img := image.NewRGBA(image.Rect(0, 0, radius*2, radius*2))

// 	x, y := 0, radius
// 	d := 3 - 2*radius
// 	for x < y {
// 		if d < 0 {
// 			d = d + 4*x + 6
// 		} else {
// 			d = d + 4*(x-y) + 10
// 			y--
// 		}
// 		x++
// 		img.Set(x, y, black)
// 	}
// 	return img
// }

func paintDataCircle(radius int, bits []int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, radius*2, radius*2))

	for i := 0.0; i < 360; i += 0.05 {
		x := float64(radius) + float64(radius)*math.Cos(float64(i)*math.Pi/180)
		y := float64(radius) + float64(radius)*math.Sin(float64(i)*math.Pi/180)

		color := black
		if bits != nil && len(bits) > 0 {
			pos := int(i * float64(len(bits)) / 360)
			if bits[pos] == 0 {
				color = white
			}
		}

		img.Set(int(x), int(y), color)

		x1 := float64(radius+1) + float64(radius+1)*math.Cos(i*math.Pi/180)
		y1 := float64(radius+1) + float64(radius+1)*math.Sin(i*math.Pi/180)
		img.Set(int(x1), int(y1), color)
	}
	return img
}

func main() {
	logo, err := loadImageObj()
	if err != nil {
		panic("logo not exist")
	}

	size := 430
	target := image.NewRGBA(image.Rect(0, 0, size, size)) //*NRGBA (image.Image interface)

	// fill m in blue
	draw.Draw(target, target.Bounds(), &image.Uniform{white}, image.ZP, draw.Src)

	// draw a line
	// for i := target.Bounds().Min.X + 10 i < target.Bounds().Max.X-10 i++ {
	// 	target.Set(i, target.Bounds().Max.Y/2, black) // to change a single pixel
	// }

	// angle := math.Sqrt(float64(size * size * 2))
	center := int(size / 2)
	logoPos := center - logo.Bounds().Max.Y/2
	draw.Draw(target, target.Bounds(), logo, image.Point{X: -logoPos, Y: -logoPos}, draw.Over)

	dataList := [][]int{
		nil,
		{1, 1, 0, 1, 0, 1, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 1, 1},
		{1, 1, 0, 0, 1, 1, 0, 1},
		{1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 1, 1, 0, 1, 0, 1, 0, 1},
	}
	for i, radius := range []int{90, 110, 130, 150, 185} {
		data := []int{}
		if i < len(dataList) {
			data = dataList[i]
		}

		line := paintDataCircle(radius, data)
		draw.Draw(target, target.Bounds(), line, image.Point{X: -(center - radius), Y: -(center - radius)}, draw.Over)
	}

	leftPos := -60
	rightPos := -size - leftPos + 40
	circle := buildCircle(40)
	draw.Draw(target, target.Bounds(), circle, image.Point{X: leftPos, Y: leftPos}, draw.Over)
	draw.Draw(target, target.Bounds(), circle, image.Point{X: leftPos, Y: rightPos}, draw.Over)
	draw.Draw(target, target.Bounds(), circle, image.Point{X: rightPos, Y: rightPos}, draw.Over)

	mark := buildCircle(80)
	draw.Draw(target, target.Bounds(), mark, image.Point{X: -size - leftPos + 70, Y: leftPos + 20}, draw.Over)

	w, _ := os.Create("qrcode.png")
	defer w.Close()
	png.Encode(w, target) //Encode writes the Image m to w in PNG format.

	Show(w.Name())

}

// show  a specified file by Preview.app for OS X(darwin)
func Show(name string) {
	command := "open"
	arg1 := "-a"
	arg2 := "/Applications/Preview.app"
	cmd := exec.Command(command, arg1, arg2, name)
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
