package provider

//	MIT License
//
//	Copyright (c) Microsoft Corporation. All rights reserved.
//
//	Permission is hereby granted, free of charge, to any person obtaining a copy
//	of this software and associated documentation files (the "Software"), to deal
//	in the Software without restriction, including without limitation the rights
//	to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
//	copies of the Software, and to permit persons to whom the Software is
//	furnished to do so, subject to the following conditions:
//
//	The above copyright notice and this permission notice shall be included in all
//	copies or substantial portions of the Software.
//
//	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
//	IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
//	FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
//	AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
//	LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
//	OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
//	SOFTWARE

import (
	"context"
	"runtime"

	"github.com/Azure/azure-amqp-common-go/auth"
	"github.com/Azure/azure-amqp-common-go/cbs"

	log "github.com/sirupsen/logrus"
	"pack.ag/amqp"
)

type (
	// Provider provides AMQP connections with authentication using the provided TokenProvider
	Provider struct {
		Endpoint      string
		TokenProvider auth.TokenProvider
		Logger        *log.Logger
	}
)

// New creates a new AMQP connection Provider that connects to the endpoint provided with authentication from the tokenProvider
func New(endpoint string, tokenProvider auth.TokenProvider) *Provider {
	p := &Provider{
		Endpoint:      endpoint,
		TokenProvider: tokenProvider,
		Logger:        log.New(),
	}
	p.Logger.SetLevel(log.WarnLevel)
	return p
}

// NewConnection generates a new AMQP connection
func (p *Provider) NewConnection() (*amqp.Client, error) {
	return amqp.Dial(p.GetAmqpHostURI(),
		amqp.ConnSASLAnonymous(),
		amqp.ConnMaxSessions(65535),
		amqp.ConnProperty("product", "MSGolangClient"),
		amqp.ConnProperty("version", "0.0.1"),
		amqp.ConnProperty("platform", runtime.GOOS),
		amqp.ConnProperty("framework", runtime.Version()),
	)
}

// NegotiateClaim authenticates the the connection over CBS
func (p *Provider) NegotiateClaim(ctx context.Context, conn *amqp.Client, entityPath string) error {
	audience := p.GetEntityAudience(entityPath)
	return cbs.NegotiateClaim(ctx, audience, conn, p.TokenProvider)
}

// GetAmqpHostURI gets the host URI of the AMQP connection
func (p *Provider) GetAmqpHostURI() string {
	return "amqps://" + p.Endpoint + "/"
}

// GetEntityAudience gets the audience for a given entity path
func (p *Provider) GetEntityAudience(entityPath string) string {
	return p.Endpoint + "/" + entityPath
}
