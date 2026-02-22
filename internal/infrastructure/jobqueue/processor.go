package jobqueue

import (
	"context"
	"log"

	"smart_alert_system/internal/usecase"
)

const (
	TaskMorningAlert            = "scheduler:morning_alert"
	TaskEveningSummary          = "scheduler:evening_summary"
	TaskActivityReminder        = "scheduler:activity_reminder"
	TaskProcessActivityReminder = "scheduler:process_activity_reminder"
	TaskProcessMorningAlert     = "scheduler:process_morning_alert"
	TaskProcessEveningSummary   = "scheduler:process_evening_summary"
)

// Processor handles the processing of global task queue
type Processor struct {
	schedulerUC *usecase.SchedulerUseCase
}

// NewProcessor creates a new Processor instance
func NewProcessor(schedulerUC *usecase.SchedulerUseCase) *Processor {
	return &Processor{
		schedulerUC: schedulerUC,
	}
}

// HandleMorningAlert handles the Morning Alert task
func (p *Processor) HandleMorningAlert(ctx context.Context, payload []byte) error {
	log.Println("jobqueue: Processing Morning Alert...")
	if err := p.schedulerUC.SendMorningAlerts(ctx); err != nil {
		log.Printf("jobqueue: Error sending morning alerts: %v", err)
		return err
	}
	return nil
}

// HandleProcessMorningAlert handles processing a single morning alert per user
func (p *Processor) HandleProcessMorningAlert(ctx context.Context, payload []byte) error {
	userIDStr := string(payload)
	if err := p.schedulerUC.ProcessSingleMorningAlert(ctx, userIDStr); err != nil {
		log.Printf("jobqueue: Error processing single morning alert %s: %v", userIDStr, err)
		return err
	}
	return nil
}

// HandleEveningSummary handles the Evening Summary task
func (p *Processor) HandleEveningSummary(ctx context.Context, payload []byte) error {
	log.Println("jobqueue: Processing Evening Summary...")
	if err := p.schedulerUC.SendEveningSummaries(ctx); err != nil {
		log.Printf("jobqueue: Error sending evening summaries: %v", err)
		return err
	}
	return nil
}

// HandleProcessEveningSummary handles processing a single evening summary per user
func (p *Processor) HandleProcessEveningSummary(ctx context.Context, payload []byte) error {
	userIDStr := string(payload)
	if err := p.schedulerUC.ProcessSingleEveningSummary(ctx, userIDStr); err != nil {
		log.Printf("jobqueue: Error processing single evening summary %s: %v", userIDStr, err)
		return err
	}
	return nil
}

// HandleActivityReminder handles the Activity Reminder task
func (p *Processor) HandleActivityReminder(ctx context.Context, payload []byte) error {
	if err := p.schedulerUC.SendActivityReminders(ctx); err != nil {
		log.Printf("jobqueue: Error sending activity reminders: %v", err)
		return err
	}
	return nil
}

// HandleProcessActivityReminder handles processing a single activity reminder per user
func (p *Processor) HandleProcessActivityReminder(ctx context.Context, payload []byte) error {
	activityIDStr := string(payload)
	if err := p.schedulerUC.ProcessSingleActivityReminder(ctx, activityIDStr); err != nil {
		log.Printf("jobqueue: Error processing single activity reminder %s: %v", activityIDStr, err)
		return err
	}
	return nil
}
