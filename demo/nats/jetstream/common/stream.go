package common

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"math/rand"
	"strconv"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

const (
	StreamName     = "EVENTS"
	SubjectForward = "event.>"
	Subject        = "event.processing"
)

var (
	nc *nats.Conn
)

func InitNatsConnection() (func(), error) {
	var err error
	logger := slog.Default().With(slog.String("nats", *FlagNatsAddr))
	var isSafeClose = make(chan struct{})
	nc, err = nats.Connect(*FlagNatsAddr, nats.ClosedHandler(func(c *nats.Conn) {
		logger.Info("NATS connection closed", slog.Any("reason", c.LastError()))
		close(isSafeClose)
	}))
	if err == nil {
		return func() {
			innerErr := nc.Drain()
			if innerErr != nil {
				logger.Error("Failed to drain NATS connection", slog.Any("error", innerErr))
			}
			<-isSafeClose
		}, nil
	}
	return nil, err
}

func defaultJetStreamConfig() jetstream.StreamConfig {
	return jetstream.StreamConfig{
		Name:     StreamName,
		Subjects: []string{SubjectForward},
		Storage:  jetstream.FileStorage,

		Retention: jetstream.WorkQueuePolicy,
	}
}

func RunProducer(ctx context.Context) {
	idx, _ := RunnerIdxFromContext(ctx)
	logger := slog.Default().With(slog.Int("producer", idx))
	js, err := jetstream.New(nc)
	if err != nil {
		logger.Error("Failed to create JetStream context", slog.Any("error", err))
		return
	}

	stream, err := js.CreateOrUpdateStream(ctx, defaultJetStreamConfig())
	if err != nil {
		logger.Error("Failed to create JetStream stream", slog.Any("error", err))
		return
	}
	LogStreamInfo(ctx, logger, stream)

	strIdx := strconv.Itoa(idx)
	for i := range 10 {
		// 随机睡眠一会
		time.Sleep(time.Duration(rand.Intn(4)) * time.Millisecond * 500)
		msg := "hello:" + strIdx + "-" + strconv.Itoa(i)
		ack, err := js.Publish(ctx, Subject, []byte(msg))
		if err != nil {
			logger.Error("Failed to publish message", slog.Any("error", err))
			continue
		}
		logger.Info("Published message", slog.Any("ack", ack), slog.String("data", msg))
	}

}

func RunConsumer(ctx context.Context) {
	idx, _ := RunnerIdxFromContext(ctx)
	logger := slog.Default().With(slog.Int("consumer", idx))

	js, err := nc.JetStream(
		nats.PublishAsyncMaxPending(256),
		nats.Context(ctx),
	)
	if err != nil {
		logger.Error("Failed to create JetStream context", slog.Any("error", err))
		return
	}

	js.AddStream(&nats.StreamConfig{
		Name:     StreamName,
		Subjects: []string{SubjectForward},
	})

	sub, err := js.QueueSubscribeSync(Subject, "consumer")
	if err != nil {
		logger.Error("Failed to subscribe to JetStream", slog.Any("error", err))
		return
	}

	for {
		msg, err := sub.NextMsgWithContext(ctx)
		if err != nil {
			if errors.Is(err, nats.ErrTimeout) {
				continue
			}
			logger.Error("Failed to get message", slog.Any("error", err))
			return
		}
		logger.Info("Received message", slog.String("data", string(msg.Data)))
		msg.AckSync(nats.Context(ctx))
	}

}

func RunConsumer2(ctx context.Context) {
	idx, _ := RunnerIdxFromContext(ctx)
	logger := slog.Default().With(slog.Int("consumer", idx))
	js, err := jetstream.New(nc)
	if err != nil {
		logger.Error("Failed to create JetStream context", slog.Any("error", err))
		return
	}

	stream, err := js.CreateOrUpdateStream(ctx, defaultJetStreamConfig())
	if err != nil {
		logger.Error("Failed to create JetStream stream", slog.Any("error", err))
		return
	}
	LogStreamInfo(ctx, logger, stream)

	consumer, err := stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Name:          "processing",
		Durable:       "processing",
		Description:   "this is a consumer",
		DeliverPolicy: jetstream.DeliverAllPolicy,
		// AckPolicy:          0,
		// AckWait:            0,
		// MaxDeliver:         0,
		// BackOff:            []time.Duration{},
		// FilterSubject: Subject,
		ReplayPolicy: jetstream.ReplayInstantPolicy,
		// RateLimit:          0,
		// SampleFrequency:    "",
		// MaxWaiting:         0,
		// MaxAckPending:      0,
		// HeadersOnly:        false,
		// MaxRequestBatch:    0,
		// MaxRequestExpires:  0,
		// MaxRequestMaxBytes: 0,
		// InactiveThreshold:  0,
		// Replicas:           0,
		// MemoryStorage:      false,
	})

	if err != nil {
		logger.Error("Failed to create JetStream consumer", slog.Any("error", err))
		return
	}

	consumer.Consume(func(msg jetstream.Msg) {
		metaData, _ := msg.Metadata()
		logger.Info("Received message", slog.String("data", string(msg.Data())), slog.Any("seq", metaData.Sequence))
		msg.DoubleAck(context.Background())
	}, jetstream.PullMaxMessages(10), jetstream.PullExpiry(time.Second))
}

func LogStreamInfo(ctx context.Context, logger *slog.Logger, stream jetstream.Stream) {
	info, err := stream.Info(ctx)
	if err != nil {
		logger.Error("Failed to get stream info", slog.Any("error", err))
		return
	}
	b, _ := json.MarshalIndent(info.State, "", "  ")
	logger.Info("Stream info", slog.String("state", string(b)))
}
