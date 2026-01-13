package worker

import (
	"log/slog"
	"time"
)

// Task represents an asynchronous task
type Task struct {
	Type    string
	Payload map[string]interface{}
}

// Worker handles asynchronous background tasks
type Worker struct {
	taskQueue chan Task
}

var GlobalWorker *Worker

// InitWorker initializes the background worker
func InitWorker() {
	GlobalWorker = &Worker{
		taskQueue: make(chan Task, 100),
	}
	go GlobalWorker.start()
	go GlobalWorker.startReminderTicker()
}

func (w *Worker) startReminderTicker() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		w.checkDueTodos()
	}
}

func (w *Worker) checkDueTodos() {
	slog.Info("Checking for todos due soon...")
	// This would typically query the DB for todos due in the next X minutes
	// For this example, we just log the check.
}

// Enqueue adds a task to the queue
func (w *Worker) Enqueue(t Task) {
	w.taskQueue <- t
}

func (w *Worker) start() {
	slog.Info("Background worker started")
	for task := range w.taskQueue {
		w.processTask(task)
	}
}

func (w *Worker) processTask(t Task) {
	slog.Info("Processing background task", slog.String("type", t.Type))

	switch t.Type {
	case "SEND_WELCOME_EMAIL":
		// Mock email sending
		time.Sleep(1 * time.Second) // Simulate network delay
		slog.Info("WELCOME EMAIL SENT",
			slog.String("email", t.Payload["email"].(string)),
			slog.String("name", t.Payload["name"].(string)),
		)
	case "TODO_COMPLETED_NOTIFICATION":
		time.Sleep(500 * time.Millisecond)
		slog.Info("TODO COMPLETION LOGGED",
			slog.String("title", t.Payload["title"].(string)),
		)
	default:
		slog.Warn("Unknown task type", slog.String("type", t.Type))
	}
}
