package middleware

import (
    "context"
    "convention.ninja/utilities"
    firebase "firebase.google.com/go"
    "firebase.google.com/go/auth"
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
)


type TokenMiddleware struct {
    KeyPath string
}

func (t *TokenMiddleware) getToken(r *http.Request) (string, error) {
    auth := r.Header.Get(authorizationHeader)
    token := ""
    if strings.HasPrefix(auth, "Bearer") {
        authSplit := strings.SplitN(auth, " ", 2)
        if len(authSplit) == 2 {
            token = authSplit[1]
            jwt, err := jws.ParseJWT([]byte(token))
            if err != nil {
                return "", err
            }
            err = jwt.Verify(kmspb.AsymmetricSignRequest{
                Name: t.KeyPath,
            }, utilities.SigningMethodKRsaPss)
            if err != nil {
                return "", err
            }
            return token, nil
        }
    }
    return "", nil
}

func (t *TokenMiddleware) Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token, err := t.getToken(r)
        if err != nil {
            w.WriteHeader(403)
            return
        }
        if len(token) > 0 {
            r = r.WithContext(context.WithValue(r.Context(), "token", token))
        }
        next.ServeHTTP(w, r)
    })
}
