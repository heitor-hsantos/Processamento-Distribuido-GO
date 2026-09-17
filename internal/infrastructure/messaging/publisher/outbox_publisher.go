package publisher

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"JungleGaming-test/config"
	"JungleGaming-test/internal/application/port"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type OutboxPublisher struct {
	repo     port.OutboxRepository
	sqsCli   *sqs.Client
	queueURL string
	cancel   context.CancelFunc
	wg       sync.WaitGroup
}

func NewOutboxPublisher(cfg config.Config, repo port.OutboxRepository) (*OutboxPublisher, error) {
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{URL: "http://localhost:4566"}, nil
	})

	awscfg, err := awsconfig.LoadDefaultConfig(context.Background(),
		awsconfig.WithRegion("us-east-1"),
		awsconfig.WithEndpointResolverWithOptions(customResolver),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load aws config: %w", err)
	}

	cli := sqs.NewFromConfig(awscfg)

	return &OutboxPublisher{
		repo:     repo,
		sqsCli:   cli,
		queueURL: cfg.SQSQueueURL,
	}, nil
}

func (p *OutboxPublisher) Start(ctx context.Context) error {
	publishCtx, cancel := context.WithCancel(context.Background())
	p.cancel = cancel
	p.wg.Add(1)

	go p.poll(publishCtx)
	log.Println("OutboxPublisher started")
	return nil
}

func (p *OutboxPublisher) Stop(ctx context.Context) error {
	if p.cancel != nil {
		p.cancel()
	}
	p.wg.Wait()
	log.Println("OutboxPublisher stopped")
	return nil
}

func (p *OutboxPublisher) poll(ctx context.Context) {
	defer p.wg.Done()
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			p.processBatch(ctx)
		}
	}
}

func (p *OutboxPublisher) processBatch(ctx context.Context) {
	records, err := p.repo.FindPending(ctx, 50)
	if err != nil {
		log.Printf("outbox publisher: failed to fetch pending records: %v", err)
		return
	}

	for _, rec := range records {
		input := &sqs.SendMessageInput{
			QueueUrl:    &p.queueURL,
			MessageBody: aws.String(string(rec.Payload)),
		}
		if strings.HasSuffix(p.queueURL, ".fifo") {
			input.MessageGroupId = aws.String("wallet-events")
			input.MessageDeduplicationId = aws.String(rec.ID)
		}

		_, err := p.sqsCli.SendMessage(ctx, input)
		if err != nil {
			log.Printf("outbox publisher: failed to publish message %s: %v", rec.ID, err)
			continue
		}

		if err := p.repo.MarkAsPublished(ctx, rec.ID); err != nil {
			log.Printf("outbox publisher: failed to mark %s as published: %v", rec.ID, err)
		}
	}
}
