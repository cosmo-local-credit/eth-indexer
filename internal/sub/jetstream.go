package sub

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"time"

	"github.com/cosmo-local-credit/eth-indexer/pkg/router"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type (
	JetStreamOpts struct {
		Endpoints   []string
		JetStreamID string
		Logg        *slog.Logger
		Router      *router.Router
	}

	JetStreamSub struct {
		jsIter    jetstream.MessagesContext
		logg      *slog.Logger
		natsConn  *nats.Conn
		router    *router.Router
		durableID string
	}
)

const (
	pullStream  = "TRACKER"
	pullSubject = "TRACKER.*"
)

func NewJetStreamSub(o JetStreamOpts) (*JetStreamSub, error) {
	natsConn, err := nats.Connect(
		strings.Join(o.Endpoints, ","),
		nats.MaxReconnects(-1),
		nats.ReconnectWait(2*time.Second),
	)
	if err != nil {
		return nil, err
	}

	js, err := jetstream.New(natsConn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stream, err := js.Stream(ctx, pullStream)
	if errors.Is(err, jetstream.ErrStreamNotFound) {
		stream, err = js.CreateStream(ctx, jetstream.StreamConfig{
			Name:     pullStream,
			Subjects: []string{pullSubject},
		})
	}
	if err != nil {
		return nil, err
	}

	consumer, err := stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Durable:       o.JetStreamID,
		AckPolicy:     jetstream.AckExplicitPolicy,
		FilterSubject: pullSubject,
	})
	if err != nil {
		return nil, err
	}
	o.Logg.Info("successfully connected to NATS server")

	iter, err := consumer.Messages(
		jetstream.WithMessagesErrOnMissingHeartbeat(false),
		jetstream.PullMaxMessages(10),
	)
	if err != nil {
		return nil, err
	}

	return &JetStreamSub{
		jsIter:    iter,
		router:    o.Router,
		natsConn:  natsConn,
		logg:      o.Logg,
		durableID: o.JetStreamID,
	}, nil
}

func (s *JetStreamSub) Close() {
	s.jsIter.Stop()
}

func (s *JetStreamSub) Process() {
	for {
		msg, err := s.jsIter.Next()
		if err != nil {
			if errors.Is(err, jetstream.ErrMsgIteratorClosed) {
				s.logg.Debug("jetstream: iterator closed")
				return
			} else {
				s.logg.Debug("jetstream: unknown iterator error", "error", err)
				continue
			}
		}

		s.logg.Debug("processing nats message", "subject", msg.Subject())
		if err := s.router.Handle(context.Background(), msg); err != nil {
			s.logg.Error("jetstream: router: error processing nats message", "error", err)
		}
	}
}
