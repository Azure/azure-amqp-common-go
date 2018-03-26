package conn

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
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	namespace = "mynamespace"
	keyName   = "keyName"
	secret    = "superSecret="
	hubName   = "myhub"
	connStr1  = "Endpoint=sb://" + namespace + ".servicebus.windows.net/;SharedAccessKeyName=" + keyName + ";SharedAccessKey=" + secret + ";EntityPath=" + hubName
	connStr2  = "Endpoint=sb://" + namespace + ".servicebus.windows.net/;SharedAccessKeyName=" + keyName + ";SharedAccessKey=" + secret
)

func TestParsedConnectionFromStr(t *testing.T) {
	parsed, err := ParsedConnectionFromStr(connStr1)
	assert.Nil(t, err)
	assert.Equal(t, "amqps://"+namespace+".servicebus.windows.net", parsed.Host)
	assert.Equal(t, namespace, parsed.Namespace)
	assert.Equal(t, keyName, parsed.KeyName)
	assert.Equal(t, secret, parsed.Key)
	assert.Equal(t, hubName, parsed.HubName)
}

func TestParsedConnectionFromStrWithoutEntityPath(t *testing.T) {
	parsed, err := ParsedConnectionFromStr(connStr2)
	assert.Nil(t, err)
	assert.Equal(t, "amqps://"+namespace+".servicebus.windows.net", parsed.Host)
	assert.Equal(t, namespace, parsed.Namespace)
	assert.Equal(t, keyName, parsed.KeyName)
	assert.Equal(t, secret, parsed.Key)
	assert.Equal(t, "", parsed.HubName)
}
