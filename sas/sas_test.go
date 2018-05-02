package sas

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

const (
	sas = "SharedAccessSignature"
)

type (
	sig struct {
		sr  string
		se  string
		skn string
		sig string
	}
)

func TestNewSigner(t *testing.T) {
	keyName, key := "foo", "superSecret"
	signer := NewSigner(keyName, key)
	before := time.Now().UTC().Add(-2 * time.Second)
	sigStr, expiry := signer.SignWithDuration("http://microsoft.com", 1*time.Hour)
	nixExpiry, err := strconv.ParseInt(expiry, 10, 64)
	assert.True(t, time.Now().UTC().Add(1*time.Hour).After(time.Unix(nixExpiry, 0)), "now + 1 hour is after")
	fmt.Println(before, time.Unix(nixExpiry, 0).UTC(), time.Now().UTC())
	assert.True(t, before.Add(1*time.Hour).Before(time.Unix(nixExpiry, 0)), "before signing + 1 hour is before")

	sig, err := parseSig(sigStr)
	assert.Nil(t, err)
	assert.Equal(t, "http%3a%2f%2fmicrosoft.com", sig.sr)
	assert.Equal(t, keyName, sig.skn)
	assert.Equal(t, expiry, sig.se)
	assert.NotNil(t, sig.sig)
}

func parseSig(sigStr string) (*sig, error) {
	if !strings.HasPrefix(sigStr, sas+" ") {
		return nil, errors.New("should start with " + sas)
	}
	sigStr = strings.TrimPrefix(sigStr, sas+" ")
	parts := strings.Split(sigStr, "&")
	parsed := new(sig)
	for _, part := range parts {
		keyValue := strings.Split(part, "=")
		if len(keyValue) != 2 {
			return nil, errors.New("key value is malformed")
		}
		switch keyValue[0] {
		case "sr":
			parsed.sr = keyValue[1]
		case "se":
			parsed.se = keyValue[1]
		case "sig":
			parsed.sig = keyValue[1]
		case "skn":
			parsed.skn = keyValue[1]
		default:
			return nil, errors.New(fmt.Sprintf("unknown key / value: %q", keyValue))
		}
	}
	return parsed, nil
}
