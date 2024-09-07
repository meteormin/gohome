package schedule

import (
	"encoding/json"
	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"time"
)

var logTag = "scheduler"
var worker *Worker

type Worker struct {
	scheduler gocron.Scheduler
	startAt   time.Time
	quitChan  chan bool
	logger    *zap.SugaredLogger
}

type WorkerConfig struct {
	SchedulerOptions []gocron.SchedulerOption
	Logger           *zap.SugaredLogger
}

type JobStats struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Error   error     `json:"error"`
	LastRun time.Time `json:"last_run"`
	Tags    []string  `json:"tags"`
}

func (js JobStats) Marshal() (string, error) {
	marshal, err := json.Marshal(js)
	if err != nil {
		return "null", err
	}

	return string(marshal), nil
}

type LoggerAdapter struct {
	core *zap.SugaredLogger
}

func (la *LoggerAdapter) Info(format string, args ...interface{}) {
	la.core.Infof(format, args...)
}

func (la *LoggerAdapter) Error(format string, args ...interface{}) {
	la.core.Errorf(format, args...)
}

func (la *LoggerAdapter) Debug(format string, args ...interface{}) {
	la.core.Debugf(format, args...)
}

func (la *LoggerAdapter) Warn(format string, args ...interface{}) {
	la.core.Warnf(format, args...)
}

func GetWorker() *Worker {
	return worker
}

func NewWorker(cfg WorkerConfig) (*Worker, error) {
	if cfg.Logger != nil {
		cfg.SchedulerOptions = append(cfg.SchedulerOptions,
			gocron.WithLogger(&LoggerAdapter{core: cfg.Logger}))
	}

	s, err := gocron.NewScheduler(cfg.SchedulerOptions...)
	if err != nil {
		return nil, err
	}
	worker = &Worker{
		scheduler: s,
		quitChan:  make(chan bool),
		logger:    cfg.Logger,
	}

	return worker, nil
}

func (w *Worker) FindJobByID(id uuid.UUID) gocron.Job {
	for _, job := range w.scheduler.Jobs() {
		if job.ID() == id {
			return job
		}
	}
	return nil
}

func (w *Worker) FindJobByName(name string) gocron.Job {
	for _, job := range w.scheduler.Jobs() {
		if job.Name() == name {
			return job
		}
	}
	return nil
}

func (w *Worker) Jobs() []gocron.Job {
	return w.scheduler.Jobs()
}

func (w *Worker) NewJob(jd gocron.JobDefinition, task gocron.Task, opts ...gocron.JobOption) (gocron.Job, error) {
	var err error
	job, err := w.scheduler.NewJob(jd, task, opts...)
	if err != nil {
		return nil, err
	}

	return job, err
}

func (w *Worker) RemoveJobByName(id uuid.UUID) error {
	err := w.scheduler.RemoveJob(id)
	if err != nil {
		return err
	}
	return nil
}

func (w *Worker) Stats() []JobStats {
	stats := make([]JobStats, 0)
	for _, job := range w.Jobs() {
		lastRun, err := job.LastRun()
		stats = append(stats, JobStats{
			ID:      job.ID(),
			Name:    job.Name(),
			Error:   err,
			LastRun: lastRun,
			Tags:    job.Tags(),
		})
	}

	return stats
}

func (w *Worker) statsJob() {
	stats := w.Stats()
	runningTime := time.Now().Sub(w.startAt)

	if w.logger != nil {
		w.logger.Infof("[%s] running time: %f sec", logTag, runningTime.Seconds())
	}

	for _, s := range stats {
		marshal, err := s.Marshal()
		if err != nil {
			if w.logger != nil {
				w.logger.Infof("[%s] error: %s", logTag, err)
			}
			continue
		}
		if w.logger != nil {
			w.logger.Infof("[%s] %s", logTag, marshal)
		}
	}
}

func (w *Worker) Run() {
	job, err := w.scheduler.NewJob(
		gocron.CronJob("* * * * *", false),
		gocron.NewTask(func() {
			w.statsJob()
		}),
		gocron.WithName("stats"),
		gocron.WithTags("stats"),
	)

	w.scheduler.Start()
	w.startAt = time.Now()
	if w.logger != nil {
		w.logger.Infof("[%s] Start... %s", logTag, w.startAt)
	}

	err = job.RunNow()
	if err != nil {
		if w.logger != nil {
			w.logger.Errorf("[%s] %v", logTag, err)
		}
	}

	for {
		select {
		case <-w.quitChan:
			if w.logger != nil {
				w.logger.Infof("[%s] Stop", logTag)
			}
			return
		default:
		}
		time.Sleep(time.Second)
	}
}

func (w *Worker) Stop() error {
	w.quitChan <- true
	err := w.scheduler.Shutdown()
	if err != nil {
		return err
	}
	return nil
}
