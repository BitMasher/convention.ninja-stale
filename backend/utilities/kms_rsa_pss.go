package utilities

import (
    "context"
    c "crypto"
    "crypto/rsa"
    "encoding/json"
    "errors"
    "github.com/SermoDigital/jose/crypto"

    kms "cloud.google.com/go/kms/apiv1"
    kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
)

type SigningMethodKmsRsaPss struct {
    *crypto.SigningMethodRSA
    Options *rsa.PSSOptions
}

var (
    SigningMethodKRsaPss = &SigningMethodKmsRsaPss{
        &crypto.SigningMethodRSA{
            Name: "PS256",
            Hash: c.SHA256,
        },
        &rsa.PSSOptions{
            SaltLength: rsa.PSSSaltLengthAuto,
            Hash:       c.SHA256,
        },
    }
)

// Verify implements the Verify method from SigningMethod.
// For this verify method, key must be an *rsa.PublicKey.
func (m *SigningMethodKmsRsaPss) Verify(raw []byte, signature crypto.Signature, key interface{}) error {
    rsaKey, ok := key.(*rsa.PublicKey)
    if !ok {
        return crypto.ErrInvalidKey
    }
    return rsa.VerifyPSS(rsaKey, m.Hash, m.sum(raw), signature, m.Options)
}

// Sign implements the Sign method from SigningMethod.
// For this signing method, key must be an *rsa.PrivateKey.
func (m *SigningMethodKmsRsaPss) Sign(raw []byte, key interface{}) (crypto.Signature, error) {
    ctx := context.Background()
    client, err := kms.NewKeyManagementClient(ctx)
    if err != nil {
        return nil, errors.New("kms failure")
    }
    defer client.Close()
    
    if req, ok := key.(kmspb.AsymmetricSignRequest); ok {
        resp, err := client.AsymmetricSign(ctx, &req)
        if err != nil {
            return nil, err
        }
        return resp.Signature, nil
    } else {
        return nil, crypto.ErrNotRSAPrivateKey
    }
}

func (m *SigningMethodKmsRsaPss) sum(b []byte) []byte {
    h := m.Hash.New()
    h.Write(b)
    return h.Sum(nil)
}

// Hasher implements the Hasher method from SigningMethod.
func (m *SigningMethodKmsRsaPss) Hasher() c.Hash { return m.Hash }

// MarshalJSON implements json.Marshaler.
// See SigningMethodECDSA.MarshalJSON() for information.
func (m *SigningMethodKmsRsaPss) MarshalJSON() ([]byte, error) {
    return []byte(`"` + m.Alg() + `"`), nil
}

var _ json.Marshaler = (*SigningMethodKmsRsaPss)(nil)