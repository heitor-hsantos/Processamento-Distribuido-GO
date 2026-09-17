package consumer

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"JungleGaming-test/config"
	"JungleGaming-test/internal/application/port"
	"JungleGaming-test/internal/application/usecase"
	"JungleGaming-test/internal/domain/money"
	"JungleGaming-test/internal/domain/outbox"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/google/uuid"
)

type BetConsumer struct {
	sqsCli   *sqs.Client
	queueURL string
	process  *usecase.ProcessBetUseCase
	inbox    port.InboxRepository
	uow      port.UnitOfWork
	cancel   context.CancelFunc
	wg       sync.WaitGroup
}

func NewBetConsumer(cfg config.Config, process *usecase.ProcessBetUseCase, inbox port.InboxRepository, uow port.UnitOfWork) (*BetConsumer, error) {
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

	return &BetConsumer{
		sqsCli:   cli,
		queueURL: cfg.SQSQueueURL,
		process:  process,
		inbox:    inbox,
		uow:      uow,
	}, nil
}

func (c *BetConsumer) Start(ctx context.Context) error {
	consumeCtx, cancel := context.WithCancel(context.Background())
	c.cancel = cancel
	c.wg.Add(1)
	go c.poll(consumeCtx)
	log.Println("BetConsumer started")
	return nil
}

func (c *BetConsumer) Stop(ctx context.Context) error {
	if c.cancel != nil {
		c.cancel()
	}
	c.wg.Wait()
	log.Println("BetConsumer stopped")
	return nil
}

func (c *BetConsumer) poll(ctx context.Context) {
	defer c.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			c.processMessages(ctx)
		}
	}
}

type betMessageEnvelope struct {
	MessageID  string            `json:"messageId"`
	Type       string            `json:"type"`
	Data       betMessagePayload `json:"data"`
	Body       betMessagePayload `json:"-"`
	OccurredAt string            `json:"occurredAt"`
}

type betMessagePayload struct {
	TransactionID  string `json:"transaction_id"`
	WalletID       string `json:"wallet_id"`
	ProviderID     string `json:"provider_id"`
	Amount         int64  `json:"amount"`
	Currency       string `json:"currency"`
	IdempotencyKey string `json:"idempotency_key"`
	ExternalID     string `json:"externalTransactionId"`
}

func (c *BetConsumer) processMessages(ctx context.Context) {
	out, err := c.sqsCli.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:            &c.queueURL,
		MaxNumberOfMessages: 10,
		WaitTimeSeconds:     10,
	})
	if err != nil {
		log.Printf("bet consumer: failed to receive messages: %v", err)
		time.Sleep(2 * time.Second)
		return
	}

	for _, msg := range out.Messages {
		if msg.MessageId == nil || msg.ReceiptHandle == nil || msg.Body == nil {
			log.Printf("bet consumer: invalid message payload received")
			continue
		}
		if err := c.handleMessage(ctx, msg); err != nil {
			log.Printf("bet consumer: failed to handle msg %s: %v", *msg.MessageId, err)
			continue
		}

		_, err := c.sqsCli.DeleteMessage(ctx, &sqs.DeleteMessageInput{
			QueueUrl:      &c.queueURL,
			ReceiptHandle: msg.ReceiptHandle,
		})
		if err != nil {
			log.Printf("bet consumer: failed to delete msg %s: %v", *msg.MessageId, err)
		}
	}
}

func (c *BetConsumer) handleMessage(ctx context.Context, msg types.Message) error {
	return c.uow.Do(ctx, func(ctxWithTx context.Context) error {
		messageID := aws.ToString(msg.MessageId)
		if messageID == "" {
			return fmt.Errorf("bet consumer: empty message id")
		}
		body := aws.ToString(msg.Body)
		if body == "" {
			return fmt.Errorf("bet consumer: empty body")
		}

		exists, err := c.inbox.FindByMessageID(ctxWithTx, messageID)
		if err != nil {
			return err
		}
		if exists != nil {
			log.Printf("bet consumer: message %s already processed", messageID)
			return nil
		}

		var envelope betMessageEnvelope
		if err := json.Unmarshal([]byte(body), &envelope); err != nil {
			return fmt.Errorf("invalid payload: %w", err)
		}
		payload := envelope.Data
		if payload.TransactionID == "" {
			payload = envelope.Body
		}
		if payload.TransactionID == "" {
			return fmt.Errorf("invalid message: transaction_id required")
		}
		if payload.IdempotencyKey == "" {
			payload.IdempotencyKey = payload.ProviderID + ":" + payload.TransactionID
		}
		if payload.Currency == "" {
			payload.Currency = string(money.BRL)
		}

		payloadHash := sha256.Sum256([]byte(body))
		inboxRec := &outbox.InboxRecord{
			ID:          uuid.New().String(),
			MessageID:   messageID,
			PayloadHash: hex.EncodeToString(payloadHash[:]),
			Payload:     []byte(body),
			ConsumedAt:  time.Now().UTC(),
			Handled:     true,
		}
		if err := c.inbox.Save(ctxWithTx, inboxRec); err != nil {
			return err
		}

		input := usecase.ProcessBetInput{
			TransactionID:  payload.TransactionID,
			WalletID:       payload.WalletID,
			ProviderID:     payload.ProviderID,
			Amount:         payload.Amount,
			Currency:       money.Currency(payload.Currency),
			IdempotencyKey: payload.IdempotencyKey,
		}
		if _, err = c.process.Execute(ctxWithTx, input); err != nil {
			return fmt.Errorf("usecase failed: %w", err)
		}
		return nil
	})
}
