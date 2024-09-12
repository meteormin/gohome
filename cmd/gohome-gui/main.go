package main

import (
	"github.com/go-co-op/gocron/v2"
	"github.com/joho/godotenv"
	"github.com/meteormin/gohome/internal/detect"
	"github.com/meteormin/gohome/pkg/schedule"
	"github.com/meteormin/gohome/pkg/zlogger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gocv.io/x/gocv"
	"log"
	"os"
	"strconv"
	"time"
)

var logger *zap.SugaredLogger
var cfg detect.DetectorConfig
var cam *gocv.VideoCapture

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	logger = zlogger.New(zlogger.Config{
		Name:       "gohome",
		LogLevel:   zapcore.DebugLevel,
		TimeFormat: "2006-01-02 15:04:05",
		FilePath:   "./logs",
		Filename:   "gohome.log",
		MaxSize:    10,
		MaxBackups: 30,
		MaxAge:     1,
		Compress:   false,
		TimeKey:    "timestamp",
		TimeZone:   "Asia/Seoul",
	})

	cam, err = gocv.VideoCaptureDevice(0)
	if err != nil {
		logger.Fatal(err)
	}

	scheduleDuration, err := strconv.Atoi(os.Getenv("SCHEDULE_DURATION"))
	if err != nil {
		scheduleDuration = 1
	}
	duration := time.Duration(scheduleDuration) * time.Second

	frameCount, err := strconv.Atoi(os.Getenv("FRAME_COUNT"))
	if err != nil {
		frameCount = 60
	}

	frameDelay, err := strconv.Atoi(os.Getenv("FRAME_DELAY"))
	if err != nil {
		frameDelay = 1
	}
	frameDelayDuration := time.Duration(frameDelay) * time.Millisecond

	cfg = detect.DetectorConfig{
		Camera:           cam,
		SaveImagePath:    "./logs/detected",
		ScheduleDuration: duration,
		SchedulerConfig: schedule.WorkerConfig{
			SchedulerOptions: []gocron.SchedulerOption{},
			Logger:           logger,
		},
		FrameCount: frameCount,
		FrameDelay: frameDelayDuration,
	}
}

func main() {
	window := gocv.NewWindow("gohome")
	defer window.Close()

	windowWidth, widthErr := strconv.Atoi(os.Getenv("WINDOW_WIDTH"))
	windowHeight, heightErr := strconv.Atoi(os.Getenv("WINDOW_HEIGHT"))
	if widthErr == nil && heightErr == nil {
		window.ResizeWindow(windowWidth, windowHeight)
	}

	img := gocv.NewMat()
	detector, err := detect.NewDetector(cfg)
	if err != nil {
		logger.Fatal(err)
	}

	err = detector.DetectWithWindow(window, &img)
	if err != nil {
		log.Fatalln(err)
	}
}
