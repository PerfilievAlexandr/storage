package fileTransactionLoggerService

import (
	"bufio"
	"context"
	"fmt"
	"github.com/PerfilievAlexandr/storage/internal/domain"
	"github.com/PerfilievAlexandr/storage/internal/domain/enum"
	"github.com/PerfilievAlexandr/storage/internal/service"
	"os"
)

type fileTransactionLogger struct {
	events       chan<- domain.Event
	errors       <-chan error
	lastSequence uint64
	file         *os.File
}

func New(_ context.Context, filename string) (service.TransactionLoggerService, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		return nil, fmt.Errorf("cannot open transaction log file: %w", err)
	}

	return &fileTransactionLogger{file: file}, nil
}

func (f *fileTransactionLogger) WriteDelete(key string) {
	f.events <- domain.Event{EventType: enum.EventDelete, Key: key}
}

func (f *fileTransactionLogger) WritePut(key, value string) {
	f.events <- domain.Event{EventType: enum.EventPut, Key: key, Value: value}
}

func (f *fileTransactionLogger) Err() <-chan error {
	return f.errors
}

func (f *fileTransactionLogger) Run() {
	events := make(chan domain.Event, 16)
	f.events = events
	errors := make(chan error, 1)
	f.errors = errors

	go func() {
		for event := range events {
			f.lastSequence++

			_, err := fmt.Fprintf(f.file, "%d\t%d\t%s\t%s\n", f.lastSequence, event.EventType, event.Key, event.Value)

			if err != nil {
				errors <- err
			}
		}
	}()
}

func (f *fileTransactionLogger) ReadEvents() (<-chan domain.Event, <-chan error) {
	scanner := bufio.NewScanner(f.file)
	outEvents := make(chan domain.Event)
	outErrors := make(chan error, 1)

	go func() {
		var event domain.Event

		defer close(outEvents)
		defer close(outErrors)

		for scanner.Scan() {
			line := scanner.Text()

			_, err := fmt.Sscanf(line, "%d\t%d\t%s\t%s", &event.Sequence, &event.EventType, &event.Key, &event.Value)
			if err != nil {
				outErrors <- fmt.Errorf("input parse error: %w", err)
				return
			}

			if f.lastSequence >= event.Sequence {
				outErrors <- fmt.Errorf("transaction numbers out of sequence")
				return
			}

			f.lastSequence = event.Sequence

			outEvents <- event
		}

		if err := scanner.Err(); err != nil {
			outErrors <- fmt.Errorf("transaction log read failure: %w", err)
			return
		}
	}()

	return outEvents, outErrors
}
