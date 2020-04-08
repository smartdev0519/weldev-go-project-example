package nsq

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/albertwidi/go-project-example/internal/pkg/nsq/fakensq"
)

func TestStartStop(t *testing.T) {
	t.Parallel()

	var (
		topic   = "test_topic"
		channel = "test_channel"
	)

	consumer, err := fakensq.NewFakeConsumer(topic, channel, nil)
	if err != nil {
		t.Error(err)
		return
	}
	wc, err := WrapConsumers([]string{"testing"}, consumer)
	if err != nil {
		t.Error(err)
		return
	}
	if err := wc.Start(); err != nil {
		t.Error(err)
		return
	}

	// give time for consumer to start the work
	time.Sleep(time.Millisecond * 100)

	for _, h := range wc.handlers {
		if h.workerNumber == 0 {
			t.Error("worker number should not be 0 because consumer is started")
			return
		}
	}

	if err := wc.Stop(); err != nil {
		t.Error(err)
		return
	}

	// give time for consumer to stop the work
	time.Sleep(time.Millisecond * 100)

	for _, h := range wc.handlers {
		if h.workerNumber != 0 {
			t.Error("worker number should be 0 because consumer is stopped")
			return
		}
	}
}

func TestMiddlewareChaining(t *testing.T) {
	t.Parallel()

	var (
		topic             = "test_topic"
		channel           = "test_channel"
		middlewareTestVal = "middleware_test"
		errChan           = make(chan error)

		// testing expect
		messageExpect = "testing middleware chaining"
		expectResult  = "test1:test2:test3"
		// to make sure that error is being sent back
		errNil = errors.New("error should be nil")
	)

	mw1 := func(handler HandlerFunc) HandlerFunc {
		return func(ctx context.Context, message *Message) error {
			ctx = context.WithValue(ctx, &middlewareTestVal, "test1")
			return handler(ctx, message)
		}
	}

	mw2 := func(handler HandlerFunc) HandlerFunc {
		return func(ctx context.Context, message *Message) error {
			val := ctx.Value(&middlewareTestVal).(string)
			val += ":test2"
			ctx = context.WithValue(ctx, &middlewareTestVal, val)
			return handler(ctx, message)
		}
	}

	mw3 := func(handler HandlerFunc) HandlerFunc {
		return func(ctx context.Context, message *Message) error {
			val := ctx.Value(&middlewareTestVal).(string)
			val += ":test3"
			ctx = context.WithValue(ctx, &middlewareTestVal, val)
			return handler(ctx, message)
		}
	}

	consumer, err := fakensq.NewFakeConsumer(topic, channel, nil)
	if err != nil {
		t.Error(err)
		return
	}
	producer := fakensq.NewFakeProducer(consumer)

	wc, err := WrapConsumers([]string{"test"}, consumer)
	if err != nil {
		t.Error(err)
		return
	}

	// chain from left to right or top to bottom
	wc.Use(
		mw1,
		mw2,
		mw3,
	)

	// handle message and check whether the middleware chaining is correct
	wc.Handle(topic, channel, func(ctx context.Context, message *Message) error {
		if string(message.Message.Body) != messageExpect {
			err := fmt.Errorf("epecting message %s but got %s", messageExpect, string(message.Message.Body))
			errChan <- err
			return err
		}
		val := ctx.Value(&middlewareTestVal).(string)
		if val != expectResult {
			err := fmt.Errorf("middleware chaining result is not as expected, expect %s but got %s", expectResult, val)
			errChan <- err
			return err
		}

		errChan <- errNil
		return err
	})

	if err := wc.Start(); err != nil {
		t.Error(err)
		return
	}

	if err := producer.Publish(topic, []byte("testing middleware chaining")); err != nil {
		t.Error(err)
		return
	}

	err = <-errChan
	if err != errNil {
		t.Error(err)
		return
	}

	if err := wc.Stop(); err != nil {
		t.Error(err)
		return
	}
}
