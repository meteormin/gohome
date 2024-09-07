package detect

import (
	"fmt"
	"github.com/go-co-op/gocron/v2"
	"github.com/meteormin/gohome/pkg/schedule"
	"go.uber.org/zap"
	"gocv.io/x/gocv"
	"image/color"
	"log"
	"path"
	"time"
)

type Detector interface {
	Detect(img *gocv.Mat) error
	DetectWithWindow(window *gocv.Window, img *gocv.Mat) error
	CloseCamera() error
	StartSchedule()
	StopSchedule() error
	GetCurrentFrame() *gocv.Mat
}

type DetectorImpl struct {
	camera        *gocv.VideoCapture
	saveImagePath string
	scheduler     *schedule.Worker
	logger        *zap.SugaredLogger
	img           *gocv.Mat
}

func (d *DetectorImpl) Detect(img *gocv.Mat) error {
	hog := gocv.NewHOGDescriptor()
	defer hog.Close()
	err := hog.SetSVMDetector(gocv.HOGDefaultPeopleDetector())
	if err != nil {
		return err
	}

	frameCounter := 0
	motionDetected := false
	for frameCounter < 60 {
		if ok := d.camera.Read(img); !ok {
			return fmt.Errorf("cannot read device")
		}
		if img.Empty() {
			continue
		}

		// HOG를 이용한 사람 감지
		rects := hog.DetectMultiScale(*img)
		for _, rect := range rects {
			// 감지된 영역에 사각형 그리기
			gocv.Rectangle(img, rect, color.RGBA{G: 255}, 3)

			// 감지된 결과를 콘솔에 표시
			if d.logger != nil {
				d.logger.Infof("Detected person at: %v", rect)
			} else {
				log.Printf("Detected person at: %v\n", rect)
			}
			motionDetected = true
		}

		if motionDetected {
			if d.saveImagePath != "" {
				now := time.Now().Format("20060102150405")
				name := path.Join(d.saveImagePath,
					fmt.Sprintf("%s_%d.jpg", now, frameCounter))

				isSuccess := gocv.IMWrite(name, *img)
				if isSuccess {
					if d.logger != nil {
						d.logger.Infof("Saved image: %s", name)
					} else {
						log.Printf("Saved image: %s\n", name)
					}
				} else {
					if d.logger != nil {
						d.logger.Errorf("Failed to save image: %s", name)
					} else {
						log.Printf("Failed to save image: %s\n", name)
					}
				}
			}
			motionDetected = false
		}

		frameCounter++
	}

	return nil
}

func (d *DetectorImpl) DetectWithWindow(window *gocv.Window, img *gocv.Mat) error {
	hog := gocv.NewHOGDescriptor()
	defer hog.Close()
	err := hog.SetSVMDetector(gocv.HOGDefaultPeopleDetector())
	if err != nil {
		return err
	}

	frameCounter := 0
	motionDetected := false
	for {
		if d.logger != nil {
			d.logger.Infof("Read camera...")
		} else {
			log.Println("Read camera...")
		}
		if ok := d.camera.Read(img); !ok {
			return fmt.Errorf("cannot read device")
		}
		if img.Empty() {
			continue
		}

		// HOG를 이용한 사람 감지
		if d.logger != nil {
			d.logger.Infof("Detecting...")
		} else {
			log.Println("Detecting...")
		}
		rects := hog.DetectMultiScale(*img)
		for _, rect := range rects {
			// 감지된 영역에 사각형 그리기
			gocv.Rectangle(img, rect, color.RGBA{G: 255}, 3)

			// 감지된 결과를 콘솔에 표시
			if d.logger != nil {
				d.logger.Infof("Detected person at: %v", rect)
			} else {
				log.Printf("Detected person at: %v\n", rect)
			}
			motionDetected = true
		}

		if motionDetected && frameCounter%30 == 0 {
			if d.saveImagePath != "" {
				now := time.Now().Format("20060102150405")
				name := path.Join(d.saveImagePath,
					fmt.Sprintf("%s_%d.jpg", now, frameCounter))

				isSuccess := gocv.IMWrite(name, *img)
				if isSuccess {
					if d.logger != nil {
						d.logger.Infof("Saved image: %s", name)
					} else {
						log.Printf("Saved image: %s\n", name)
					}
				} else {
					if d.logger != nil {
						d.logger.Errorf("Failed to save image: %s", name)
					} else {
						log.Printf("Failed to save image: %s\n", name)
					}
				}
			}
		}

		if d.logger != nil {
			d.logger.Infof("Show frame: %d", frameCounter)
		} else {
			log.Printf("Show Frame: %d\n", frameCounter)
		}
		window.IMShow(*img)
		key := window.WaitKey(1)
		if key == 113 {
			log.Println("Quitting...")
			break
		}
		time.Sleep(30 * time.Millisecond)
		frameCounter++
		motionDetected = false
	}

	return nil
}

func (d *DetectorImpl) CloseCamera() error {
	if d.camera != nil && d.camera.IsOpened() {
		err := d.camera.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *DetectorImpl) StartSchedule() {
	d.scheduler.Run()
}

func (d *DetectorImpl) StopSchedule() error {
	return d.scheduler.Stop()
}

func (d *DetectorImpl) GetCurrentFrame() *gocv.Mat {
	return d.img
}

type DetectorConfig struct {
	Camera           *gocv.VideoCapture
	SaveImagePath    string
	ScheduleDuration time.Duration
	SchedulerConfig  schedule.WorkerConfig
}

func NewDetector(config DetectorConfig) (Detector, error) {
	var scheduler *schedule.Worker
	scheduler, err := schedule.NewWorker(config.SchedulerConfig)
	if err != nil {
		return nil, err
	}

	if config.Camera == nil {
		return nil, fmt.Errorf("camera is nil")
	}

	detector := &DetectorImpl{
		camera:        config.Camera,
		saveImagePath: config.SaveImagePath,
		scheduler:     scheduler,
		logger:        config.SchedulerConfig.Logger,
	}

	img := gocv.NewMat()
	detector.img = &img
	job, err := scheduler.NewJob(gocron.DurationJob(config.ScheduleDuration),
		gocron.NewTask(detector.Detect, &img),
		gocron.WithName("Detector"))

	if err != nil {
		return nil, err
	}

	if detector.logger != nil {
		detector.logger.Infof("camera is %v", config.Camera.IsOpened())
		detector.logger.Infof("saveImagePath: %s", config.SaveImagePath)
		detector.logger.Infof("scheduled job: ID=%s, Name=%s", job.ID(), job.Name())
	}

	return detector, nil
}
