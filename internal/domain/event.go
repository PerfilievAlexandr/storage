package domain

import "github.com/PerfilievAlexandr/storage/internal/domain/enum"

type Event struct {
	Sequence  uint64
	EventType enum.EventType
	Key       string
	Value     string
}
