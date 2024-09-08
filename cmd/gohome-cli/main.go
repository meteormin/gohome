package main

import (
	"fmt"
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
	"os/signal"
	"strconv"
	"syscall"
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

	cfg = detect.DetectorConfig{
		Camera:           cam,
		SaveImagePath:    "./logs/detected",
		ScheduleDuration: duration,
		SchedulerConfig: schedule.WorkerConfig{
			SchedulerOptions: []gocron.SchedulerOption{},
			Logger:           logger,
		},
		FrameCount: frameCount,
	}
}

func main() {
	detector, err := detect.NewDetector(cfg)
	if err != nil {
		logger.Fatal(err)
	}
	defer func(detector detect.Detector) {
		log.Println("close camera")
		err = detector.CloseCamera()
		if err != nil {
			log.Fatalln(err)
		}
	}(detector)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-quit
		fmt.Println()
		log.Printf("received exit signal %s\n", sig)

		log.Println("stop scheduler")
		err = detector.StopSchedule()
		if err != nil {
			log.Fatalln(err)
		}

		log.Println("exit")
		os.Exit(0) // 프로그램 종료
	}()

	detector.StartSchedule()
}
