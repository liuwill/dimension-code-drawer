package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
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
	// for i := target.Bounds().Min.X + 10; i < target.Bounds().Max.X-10; i++ {
	// 	target.Set(i, target.Bounds().Max.Y/2, black) // to change a single pixel
	// }

	logoPos := size/2 - logo.Bounds().Max.Y/2
	draw.Draw(target, target.Bounds(), logo, image.Point{X: -logoPos, Y: -logoPos}, draw.Over)

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
