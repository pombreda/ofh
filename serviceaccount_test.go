// Copyright 2013 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed by the Apache v2.0 license that can be
// found in the LICENSE file.

package ofh

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"net/http"
	"testing"

	"github.com/maruel/ut"
)

func TestServiceAccount(t *testing.T) {
	// Generate a unique private key. It cannot be lower than 512, because
	// otherwise there's too much data to sign for the key size.
	key, err := rsa.GenerateKey(rand.Reader, 512)
	ut.AssertEqual(t, nil, err)
	blob := x509.MarshalPKCS1PrivateKey(key)
	block := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: blob})
	s := &ServiceAccount{
		ProjectID:    "a",
		ClientID:     "c",
		EmailAddress: "e",
		PrivateKey:   string(block),
	}
	tokReply := `{"access_token":"a", "token_type": "r", "id_token": ""}`
	resp := []*http.Response{
		{StatusCode: 200, Body: asReader(tokReply)},
	}
	r := &roundTripperStub{[]*http.Request{}, resp}
	_, err = s.GetClient("scope", r)
	ut.AssertEqual(t, nil, err)
	ut.AssertEqual(t, 1, len(r.requests))
}
