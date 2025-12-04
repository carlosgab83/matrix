package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/carlosgab83/matrix/go/internal/shared/integration/logging"
	"github.com/carlosgab83/matrix/go/internal/tank/domain"
	"github.com/carlosgab83/matrix/go/internal/tank/integration/notification"
	"github.com/carlosgab83/matrix/go/internal/tank/integration/reception"
)

type NotifierService struct {
	Config   domain.Config
	Receptor reception.Receptor
	Notifier notification.Notifier
	Logger   logging.Logger
}

const (
	DefaultNotifierWriteTimeout int    = 10000
	NotifierMaxRetries          int    = 3
	TmpChatID                   string = "5406106520" // for inital development purposes
)

var maxRetries int
var waitTimeout int

func NewNotifierService(cfg domain.Config, receptor reception.Receptor, notifier notification.Notifier, logger logging.Logger) (NotifierService, error) {
	waitTimeout = DefaultNotifierWriteTimeout
	nwt, err := strconv.Atoi(cfg.NotifierWriteTimeout)
	if err == nil && cfg.NotifierWriteTimeout != "" {
		waitTimeout = nwt
	}

	maxRetries = NotifierMaxRetries
	nmr, err := strconv.Atoi(cfg.NotifierMaxRetries)
	if err == nil && cfg.NotifierMaxRetries != "" {
		maxRetries = nmr
	}

	return NotifierService{
		Config:   cfg,
		Receptor: receptor,
		Notifier: notifier,
		Logger:   logger,
	}, nil
}

func (ns NotifierService) ListenAndNotify(ctx context.Context) error {
	go ns.Receptor.BeginConsumption()

	for {
		select {
		case notificationPayload, ok := <-ns.Receptor.ReceiveCh():
			if !ok {
				return fmt.Errorf("finalizing ListenAndNotify by closed channel")
			}

			go func() {
				payload := fmt.Sprintf("%s: %f", *notificationPayload.Symbol, *notificationPayload.Price)
				_ = ns.sendWithRetries(ctx, TmpChatID, payload)
				// TODO: Send this payload to a DLQ if sendWithRetries returns error
			}()

		case <-ctx.Done():
			return fmt.Errorf("finalizing ListenAndNotify: %v", ctx.Err())
		}
	}
}

func (ns NotifierService) sendWithRetries(ctx context.Context, chatID string, payload string) error {
	for i := 0; i < maxRetries; i++ {
		ctxTimeout, _ := context.WithTimeout(ctx, time.Duration(waitTimeout*(i+1))*time.Millisecond)
		done := make(chan bool)

		go func() {
			err := ns.Notifier.Notify(ctxTimeout, chatID, payload)
			if err != nil {
				done <- false
			}
			done <- true
		}()

		select {
		case res := <-done:
			if !res {
				continue
			}
			return nil
		case <-ctxTimeout.Done():
			continue
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return nil
}
