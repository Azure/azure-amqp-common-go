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
	"errors"
	"fmt"
	"regexp"
)

var (
	connStrRegex = regexp.MustCompile(`Endpoint=sb:\/\/(?P<Host>.+?);SharedAccessKeyName=(?P<KeyName>.+?);SharedAccessKey=(?P<Key>.+?);EntityPath=(?P<HubName>.+)`)
	hostStrRegex = regexp.MustCompile(`^(?P<Namespace>.+?)\.(.+?)\/`)
)

type (
	// ParsedConn is the structure of a parsed Service Bus or Event Hub connection string.
	ParsedConn struct {
		Host      string
		Suffix    string
		Namespace string
		HubName   string
		KeyName   string
		Key       string
	}
)

// newParsedConnection is a constructor for a parsedConn and verifies each of the inputs is non-null.
func newParsedConnection(host, suffix, namespace, hubName, keyName, key string) (*ParsedConn, error) {
	if host == "" || keyName == "" || key == "" {
		return nil, errors.New("connection string contains an empty entry")
	}
	return &ParsedConn{
		Host:      "amqps://" + host,
		Suffix:    suffix,
		Namespace: namespace,
		KeyName:   keyName,
		Key:       key,
		HubName:   hubName,
	}, nil
}

// ParsedConnectionFromStr takes a string connection string from the Azure portal and returns the parsed representation.
func ParsedConnectionFromStr(connStr string) (*ParsedConn, error) {
	matches := connStrRegex.FindStringSubmatch(connStr)
	namespaceMatches := hostStrRegex.FindStringSubmatch(matches[1])
	fmt.Println(matches[1], namespaceMatches[2], namespaceMatches[1], matches[2], matches[3])
	return newParsedConnection(matches[1], namespaceMatches[2], namespaceMatches[1], matches[4], matches[2], matches[3])
}
