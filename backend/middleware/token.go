package middleware

import (
    kms "cloud.google.com/go/kms/apiv1"
    "context"
    "convention.ninja/utilities"
    "crypto/rsa"
    "errors"
    firebase "firebase.google.com/go"
    "firebase.google.com/go/auth"
    "fmt"
    "github.com/SermoDigital/jose/crypto"
    "github.com/SermoDigital/jose/jws"
    kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
    "log"
    "net/http"
    "strings"
)

var firebaseApp *firebase.App
var firebaseAuth *auth.Client

func init() {
    var err error
    firebaseApp, err = firebase.NewApp(context.Background(), nil)
    if err != nil {
        log.Fatalf("error initializing app: %v\n", err)
    }
    
    firebaseAuth, err = firebaseApp.Auth(context.Background())
    if err != nil {
        log.Fatalf("error initializing app: %v\n", err)
    }
    
}

var (
    authorizationHeader = "Authorization"
    identityHeader = "X-Identity"
    publicKeyCache = make(map[string]*rsa.PublicKey)
    errInvalidToken = errors.New("token provided is invalid")
)


type TokenMiddleware struct {
    KeyPath string
}

func (t *TokenMiddleware) getToken(r *http.Request) (string, error) {
    authHeader := r.Header.Get(authorizationHeader)
    token := ""
    if strings.HasPrefix(authHeader, "Bearer") {
        authSplit := strings.SplitN(authHeader, " ", 2)
        if len(authSplit) == 2 {
            token = authSplit[1]
            tokenJws, err := jws.Parse([]byte(token))
            if err != nil {
                return "", errInvalidToken
            }
            // TODO: verify this value is actually a string, and exists
            keyIdO := tokenJws.GetProtected("kid")
            if keyIdO == nil {
                return "", errInvalidToken
            }
            if keyId, ok := keyIdO.(string); ok {
                if key, ok := publicKeyCache[keyId]; ok {
                    err = tokenJws.Verify(key, utilities.SigningMethodKRsaPss)
                    if err != nil {
                        return "", errInvalidToken
                    }
                    return token, nil
                } else {
                    ctx := context.Background()
                    c, err := kms.NewKeyManagementClient(ctx)
                    if err != nil {
                        return "", errors.New("kms failure")
                    }
                    pubKey, err := c.GetPublicKey(ctx, &kmspb.GetPublicKeyRequest{Name: t.KeyPath + keyId})
                    if err != nil {
                        return "", errInvalidToken
                    }
                    publicKeyCache[keyId], err = crypto.ParseRSAPublicKeyFromPEM([]byte(pubKey.Pem))
                    // TODO: fix code duplication
                    err = tokenJws.Verify(pubKey.Pem, utilities.SigningMethodKRsaPss)
                    if err != nil {
                        return "", errInvalidToken
                    }
                    return token, nil
                }
            }
            return "", errInvalidToken
        }
    }
    return "", nil
}

func (t *TokenMiddleware) Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token, err := t.getToken(r)
        if err != nil {
            fmt.Println("Error in token middleware")
            fmt.Println(err)
            w.WriteHeader(403)
            return
        }
        if len(token) > 0 {
            r = r.WithContext(context.WithValue(r.Context(), "token", token))
        }
        next.ServeHTTP(w, r)
    })
}

type IdMiddleware struct {
    
}

func (t *IdMiddleware) Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        idToken := r.Header.Get(identityHeader)
        if len(idToken) > 0 {
            token, err := firebaseAuth.VerifyIDTokenAndCheckRevoked(context.Background(), idToken)
            if err != nil {
                fmt.Println("Error in id middleware")
                fmt.Println(err)
                w.WriteHeader(403)
                return
            }
            if token != nil {
                r = r.WithContext(context.WithValue(r.Context(), "idtoken", token))
            }
        }
        next.ServeHTTP(w, r)
    })
}