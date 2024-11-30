package xcomm

import (
	"context"
	"fmt"
	comatproto "github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/util/cliutil"
	"github.com/bluesky-social/indigo/xrpc"
	"github.com/go-logr/zapr"
	"github.com/jlewi/bsctl/pkg/config"
	//"github.com/golang-jwt/jwt/v5/"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"time"
)

// XRPCManager is a struct that manages XRPC connections and requests.
//
// It is primarily responsible for fetching and refreshing credentials.
// TODO(jeremy): Is there a better pattern for hanadling credential refreshing? Could we use
// a RoundTripper?
type XRPCManager struct {
	AuthManager    AuthManager
	Config         *config.Config
	client         *xrpc.Client
	expirationTime time.Time
}

// CreateClient creates a new XRPC client. It fetches credentials if needed.
// Credentials are persisted using the AuthManager.
func (m *XRPCManager) CreateClient(ctx context.Context) (*xrpc.Client, error) {
	log := zapr.NewLogger(zap.L())

	if m.Config.Host == "" {
		return nil, errors.New("host not set")
	}

	if m.Config.Handle == "" {
		return nil, errors.New("handle not set")
	}

	if m.Config.Password == "" {
		return nil, errors.New("password not set")
	}

	// If we already have a client just return it
	if m.client != nil && time.Now().Before(m.expirationTime) {
		return m.client, nil
	}

	log.Info("Creating XRPC Client")
	xrpcc := &xrpc.Client{
		Client: cliutil.NewHttpClient(),
		Host:   m.Config.Host,
		Auth:   &xrpc.AuthInfo{Handle: m.Config.Handle},
	}

	auth, err := m.AuthManager.ReadAuth()
	if err == nil && auth.AccessJwt != "" && auth.RefreshJwt != "" {
		log.Info("Auth found, attempting to refresh session")
		xrpcc.Auth = auth
		xrpcc.Auth.AccessJwt = xrpcc.Auth.RefreshJwt
		refresh, err2 := comatproto.ServerRefreshSession(context.TODO(), xrpcc)
		if err2 != nil {
			err = err2
		} else {
			xrpcc.Auth.Did = refresh.Did
			xrpcc.Auth.AccessJwt = refresh.AccessJwt
			xrpcc.Auth.RefreshJwt = refresh.RefreshJwt

			log.Info("Persisting auth information")
			if err := m.AuthManager.WriteAuth(xrpcc.Auth); err != nil {
				return nil, errors.Wrapf(err, "cannot persist authorization information")
			}
		}
	}

	if err != nil || (xrpcc.Auth.AccessJwt == "" || xrpcc.Auth.RefreshJwt == "") {
		log.Info("Auth not found, creating new session")
		auth, err := comatproto.ServerCreateSession(context.TODO(), xrpcc, &comatproto.ServerCreateSession_Input{
			Identifier: xrpcc.Auth.Handle,
			Password:   m.Config.Password,
		})
		if err != nil {
			return nil, fmt.Errorf("cannot create session: %w", err)
		}
		xrpcc.Auth.Did = auth.Did
		xrpcc.Auth.AccessJwt = auth.AccessJwt
		xrpcc.Auth.RefreshJwt = auth.RefreshJwt

		log.Info("New session created, persisting auth information")
		if err := m.AuthManager.WriteAuth(xrpcc.Auth); err != nil {
			return nil, errors.Wrapf(err, "cannot persist authorization information")
		}
	}

	eTime, err := getExpirationTime(xrpcc.Auth.AccessJwt)

	if err != nil {
		return nil, errors.Wrapf(err, "cannot get expiration time")
	}

	log.Info("Access token expiration time", "time", eTime)
	m.client = xrpcc
	// Subtract 1 minute from the duration to give us some time.
	m.expirationTime = eTime.Add(-1 * time.Minute)
	return xrpcc, nil
}

func getExpirationTime(accessJwt string) (time.Time, error) {
	eTime := time.Unix(0, 0)

	// Parse the JWT
	token, err := jwt.Parse([]byte(accessJwt), jwt.WithVerify(false))
	if err != nil {
		return eTime, errors.Wrapf(err, "cannot parse JWT")
	}

	// Extract claims
	return token.Expiration(), nil
}
