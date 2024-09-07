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

	cfg = detect.DetectorConfig{
		Camera:           cam,
		SaveImagePath:    "./logs/detected",
		ScheduleDuration: 10 * time.Second,
		SchedulerConfig: schedule.WorkerConfig{
			SchedulerOptions: []gocron.SchedulerOption{},
			Logger:           logger,
		},
	}
}

func main() {
	window := gocv.NewWindow("gohome")
	defer window.Close()

	detector, err := detect.NewDetector(cfg)
	if err != nil {
		logger.Fatal(err)
	}

	img := gocv.NewMat()
	err = detector.DetectWithWindow(window, &img)
	if err != nil {
		log.Fatalln(err)
	}
}
