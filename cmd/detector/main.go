package main

import (
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"gocv.io/x/gocv"
	"image/color"
	"log"
	"path"
	"time"
)

var inputFilePath string
var outputPath string

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	flag.StringVar(&inputFilePath, "input", "", "Input file path")
	flag.StringVar(&inputFilePath, "i", "", "Input file path")
	flag.StringVar(&outputPath, "output", "", "Output file path")
	flag.StringVar(&outputPath, "o", "", "Output file path")
	flag.Parse()
}

func main() {
	// HOG를 이용한 사람 감지
	motionDetected := false
	hog := gocv.NewHOGDescriptor()
	defer hog.Close()
	err := hog.SetSVMDetector(gocv.HOGDefaultPeopleDetector())
	if err != nil {
		log.Fatalln(err)
	}
	img := gocv.IMRead(inputFilePath, gocv.IMReadColor)
	defer img.Close()

	if img.Empty() {
		log.Fatalln("Cannot read image")
	}

	rects := hog.DetectMultiScale(img)
	for _, rect := range rects {
		// 감지된 영역에 사각형 그리기
		gocv.Rectangle(&img, rect, color.RGBA{G: 255}, 3)
		// 감지된 결과를 콘솔에 표시
		log.Printf("Detected person at: %v\n", rect)
		motionDetected = true
	}

	if motionDetected {
		now := time.Now().Format("20060102150405")
		name := path.Join(outputPath, fmt.Sprintf("%s.jpg", now))
		isSuccess := gocv.IMWrite(name, img)
		if isSuccess {
			fmt.Printf("Saved image: %s\n", name)
		} else {
			fmt.Printf("Failed to save image: %s\n", name)
		}
	}
}
