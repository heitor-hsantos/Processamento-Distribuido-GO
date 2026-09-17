package consumer

import "context"

type InboxWorker struct{}

func NewInboxWorker() *InboxWorker { return &InboxWorker{} }

func (w *InboxWorker) Start(ctx context.Context) error { return nil }
func (w *InboxWorker) Stop(ctx context.Context) error  { return nil }
