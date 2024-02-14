package services

import (
	"context"
	"time"
)

type (
	HealthService struct {
		clients []Readiness
	}

	Check struct {
		Status string
		Time   int64
	}

	Readiness interface {
		IsReady(context.Context) error
	}
)

const (
	HealthOkStatus = "OK"
	ReadyTimeout   = 1000 * time.Millisecond
)

func GetTime() int64 {
	return time.Now().UTC().UnixMilli()
}

func (hs *HealthService) GetReadiness(ctx context.Context) (*Check, error) {
	t, cancel := context.WithTimeout(ctx, ReadyTimeout)
	defer cancel()

	for _, c := range hs.clients {
		if err := c.IsReady(t); err != nil {
			return nil, err
		}
	}

	return &Check{
		Status: HealthOkStatus,
		Time:   GetTime(),
	}, nil
}

func (hs *HealthService) GetHealth(ctx context.Context) *Check {
	return &Check{
		Status: HealthOkStatus,
		Time:   GetTime(),
	}
}

func NewHealthService(c []Readiness) *HealthService {
	return &HealthService{
		clients: c,
	}
}
