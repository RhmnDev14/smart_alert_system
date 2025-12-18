package scheduler

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/robfig/cron/v3"
	"smart_alert_system/internal/usecase"
)

type Scheduler struct {
	cron          *cron.Cron
	schedulerUC   *usecase.SchedulerUseCase
	morningTime   string
	eveningTime   string
	location      *time.Location
}

func NewScheduler(schedulerUC *usecase.SchedulerUseCase, morningTime, eveningTime string, location *time.Location) *Scheduler {
	c := cron.New(cron.WithLocation(location))
	return &Scheduler{
		cron:        c,
		schedulerUC: schedulerUC,
		morningTime: morningTime,
		eveningTime: eveningTime,
		location:    location,
	}
}

func (s *Scheduler) Start() error {
	// Schedule morning alert (format: "05:00" -> "0 5 * * *")
	morningCron := s.parseTimeToCron(s.morningTime)
	_, err := s.cron.AddFunc(morningCron, func() {
		log.Println("Running morning alert scheduler...")
		ctx := context.Background()
		if err := s.schedulerUC.SendMorningAlerts(ctx); err != nil {
			log.Printf("Error sending morning alerts: %v", err)
		}
	})
	if err != nil {
		return fmt.Errorf("failed to schedule morning alert: %w", err)
	}

	// Schedule evening summary (format: "22:00" -> "0 22 * * *")
	eveningCron := s.parseTimeToCron(s.eveningTime)
	_, err = s.cron.AddFunc(eveningCron, func() {
		log.Println("Running evening summary scheduler...")
		ctx := context.Background()
		if err := s.schedulerUC.SendEveningSummaries(ctx); err != nil {
			log.Printf("Error sending evening summaries: %v", err)
		}
	})
	if err != nil {
		return fmt.Errorf("failed to schedule evening summary: %w", err)
	}

	s.cron.Start()
	log.Printf("Scheduler started. Morning alert: %s (%s), Evening summary: %s (%s)", 
		s.morningTime, morningCron, s.eveningTime, eveningCron)
	return nil
}

func (s *Scheduler) parseTimeToCron(timeStr string) string {
	// Parse time string like "05:00" or "22:00" to cron format "0 H * * *"
	// For simplicity, assume format is "HH:MM"
	if len(timeStr) >= 5 {
		hour := timeStr[0:2]
		minute := timeStr[3:5]
		return minute + " " + hour + " * * *"
	}
	// Default to 5:00 AM if parsing fails
	return "0 5 * * *"
}

func (s *Scheduler) Stop() {
	s.cron.Stop()
	log.Println("Scheduler stopped")
}

