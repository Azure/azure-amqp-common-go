// Package rpc provides functionality for request / reply messaging. It is used by package mgmt and cbs.
package rpc

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/Azure/azure-amqp-common-go"
	"github.com/Azure/azure-amqp-common-go/uuid"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"pack.ag/amqp"
)

const (
	replyPostfix   = "-reply-to-"
	statusCodeKey  = "status-code"
	descriptionKey = "status-description"
)

type (
	// Link is the bidirectional communication structure used for CBS negotiation
	Link struct {
		session       *amqp.Session
		receiver      *amqp.Receiver
		sender        *amqp.Sender
		clientAddress string
		rpcMu         sync.Mutex
		id            string
	}

	// Response is the simplified response structure from an RPC like call
	Response struct {
		Code        int
		Description string
		Message     *amqp.Message
	}
)

// NewLink will build a new request response link
func NewLink(conn *amqp.Client, address string) (*Link, error) {
	authSession, err := conn.NewSession()
	if err != nil {
		return nil, err
	}

	authSender, err := authSession.NewSender(
		amqp.LinkTargetAddress(address),
	)
	if err != nil {
		return nil, err
	}

	linkID, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	id := linkID.String()
	clientAddress := strings.Replace("$", "", address, -1) + replyPostfix + id
	authReceiver, err := authSession.NewReceiver(
		amqp.LinkSourceAddress(address),
		amqp.LinkTargetAddress(clientAddress),
	)
	if err != nil {
		return nil, err
	}

	return &Link{
		sender:        authSender,
		receiver:      authReceiver,
		session:       authSession,
		clientAddress: clientAddress,
		id:            id,
	}, nil
}

// RetryableRPC attempts to retry a request a number of times with delay
func (l *Link) RetryableRPC(ctx context.Context, times int, delay time.Duration, msg *amqp.Message) (*Response, error) {
	res, err := common.Retry(times, delay, func() (interface{}, error) {
		res, err := l.RPC(ctx, msg)
		if err != nil {
			log.Debugf("error in RPC via link %s: %v", l.id, err)
			return nil, err
		}

		switch {
		case res.Code >= 200 && res.Code < 300:
			log.Debugf("successful rpc on link %s: status code %d and description: %s", l.id, res.Code, res.Description)
			return res, nil
		case res.Code >= 500:
			errMessage := fmt.Sprintf("server error link %s: status code %d and description: %s", l.id, res.Code, res.Description)
			log.Debugln(errMessage)
			return nil, common.Retryable(errMessage)
		default:
			errMessage := fmt.Sprintf("unhandled error link %s: status code %d and description: %s", l.id, res.Code, res.Description)
			log.Debugln(errMessage)
			return nil, common.Retryable(errMessage)
		}
	})
	if err != nil {
		return nil, err
	}
	return res.(*Response), nil
}

// RPC sends a request and waits on a response for that request
func (l *Link) RPC(ctx context.Context, msg *amqp.Message) (*Response, error) {
	l.rpcMu.Lock()
	defer l.rpcMu.Unlock()

	if msg.Properties == nil {
		msg.Properties = &amqp.MessageProperties{}
	}
	msg.Properties.ReplyTo = l.clientAddress

	err := l.sender.Send(ctx, msg)
	if err != nil {
		return nil, err
	}

	res, err := l.receiver.Receive(ctx)
	if err != nil {
		return nil, err
	}

	statusCode, ok := res.ApplicationProperties[statusCodeKey].(int32)
	if !ok {
		return nil, errors.New("status codes was not found on rpc message")
	}

	description, ok := res.ApplicationProperties[descriptionKey].(string)
	if !ok {
		return nil, errors.New("description was not found on rpc message")
	}

	return &Response{
		Code:        int(statusCode),
		Description: description,
		Message:     res,
	}, err
}

// Close the link receiver, sender and session
func (l *Link) Close() error {
	if err := l.closeReceiver(); err != nil {
		_ = l.closeSender()
		_ = l.closeSession()
		return err
	}

	if err := l.closeSender(); err != nil {
		_ = l.closeSession()
		return err
	}

	return l.closeSession()
}

func (l *Link) closeReceiver() error {
	if l.receiver != nil {
		return l.receiver.Close()
	}
	return nil
}

func (l *Link) closeSender() error {
	if l.sender != nil {
		return l.sender.Close()
	}
	return nil
}

func (l *Link) closeSession() error {
	if l.session != nil {
		return l.session.Close()
	}
	return nil
}
