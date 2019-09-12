//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package api

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/gojuno/minimock"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/logicrunner/builtin/foundation"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/api/requester"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/testutils"
)

const CallUrl = "http://localhost:19192/api/rpc"

func TestTimeoutSuite(t *testing.T) {

	ctx, _ := inslogger.WithTraceField(context.Background(), "APItests")
	mc := minimock.NewController(t)
	defer mc.Wait(64 * time.Second)
	defer mc.Finish()

	ks := platformpolicy.NewKeyProcessor()
	sKey, err := ks.GeneratePrivateKey()
	require.NoError(t, err)
	sKeyString, err := ks.ExportPrivateKeyPEM(sKey)
	require.NoError(t, err)
	pKey := ks.ExtractPublicKey(sKey)
	pKeyString, err := ks.ExportPublicKeyPEM(pKey)
	require.NoError(t, err)

	userRef := gen.Reference().String()
	user, err := requester.CreateUserConfig(userRef, string(sKeyString), string(pKeyString))

	http.DefaultServeMux = new(http.ServeMux)
	cfg := configuration.NewAPIRunner(false)
	cfg.Address = "localhost:19192"
	api, err := NewRunner(&cfg, nil, nil, nil, nil, nil, nil, nil, nil)
	require.NoError(t, err)

	cr := testutils.NewContractRequesterMock(mc)
	cr.CallMock.Set(func(p context.Context, p1 *insolar.Reference, method string, p3 []interface{}, p4 insolar.PulseNumber) (insolar.Reply, *insolar.Reference, error) {
		requestReference, _ := insolar.NewReferenceFromBase58("14K3NiGuqYGqKPnYp6XeGd2kdN4P9veL6rYcWkLKWXZCu.14FFB8zfQoGznSmzDxwv4njX1aR9ioL8GHSH17QXH2AFa")
		switch method {
		case "Call":
			var result = "OK"
			data, _ := foundation.MarshalMethodResult(result, nil)
			return &reply.CallMethod{
				Result: data,
			}, requestReference, nil
		default:
			return nil, nil, errors.New("Unknown method: " + method)
		}
	})

	api.ContractRequester = cr
	api.Start(ctx)
	defer api.Stop(ctx)

	seed, err := api.SeedGenerator.Next()
	require.NoError(t, err)

	api.SeedManager.Add(*seed, 0)

	seedString := base64.StdEncoding.EncodeToString(seed[:])

	requester.SetTimeout(25)
	resp, err := requester.SendWithSeed(
		ctx,
		CallUrl,
		user,
		&requester.Params{CallSite: "member.create", CallParams: map[string]interface{}{}, PublicKey: user.PublicKey},
		seedString,
	)
	require.NoError(t, err)

	var result requester.ContractResponse
	err = json.Unmarshal(resp, &result)
	require.NoError(t, err)
	require.Nil(t, result.Error)
	require.Equal(t, "OK", result.Result.CallResult)
}

func TestDigestParser(t *testing.T) {
	invalidDigest := ""
	_, err := parseDigest(invalidDigest)
	require.Error(t, err)

	validDigest := "SHA-256=foo"
	_, err = parseDigest(validDigest)
	require.NoError(t, err)
}

func TestSignatureParser(t *testing.T) {
	invalidSignature := ""
	_, err := parseSignature(invalidSignature)

	validSignature := `keyId="member-pub-key", algorithm="ecdsa", headers="digest", signature=bar`
	_, err = parseSignature(validSignature)
	require.NoError(t, err)
}

func TestValidateRequestHeaders(t *testing.T) {
	body := []byte("foobar")
	h := sha256.New()
	_, err := h.Write(body)
	require.NoError(t, err)

	digest := h.Sum(nil)
	calculatedDigest := `SHA-256=` + base64.URLEncoding.EncodeToString(digest)
	signature := `keyId="member-pub-key", algorithm="ecdsa", headers="digest", signature=bar`
	sig, err := validateRequestHeaders(calculatedDigest, signature, body)
	require.NoError(t, err)
	require.Equal(t, "bar", sig)
}
