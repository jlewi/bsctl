package application

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
type XRPCManager struct {
	AuthManager AuthManager
	Config      *config.Config
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
	return xrpcc, nil
}

func getExpirationTime(accessJwt string) (time.Time, error) {
	eTime := time.Unix(0, 0)

	// Parse the JWT
	token, err := jwt.Parse([]byte(accessJwt), jwt.WithVerify(false))
	if err != nil {
		return eTime, errors.Wrapf(err, "cannot parse JWT")
	}
	//
	//parser := jwt.NewParser(jwt.WithoutClaimsValidation(), jwt.WithJSONNumber())
	//
	//parser.Parse(accessJwt)
	//// Parse the token
	//token, _, err := parser.ParseUnverified(accessJwt, jwt.MapClaims{})
	//if err != nil {
	//	return eTime, err
	//}

	// Extract claims
	return token.Expiration(), nil
	//claims, ok := token.Claims.(jwt.MapClaims)
	//if !ok {
	//	return eTime, errors.New("invalid token claims")
	//}
	//
	//// Get expiration time
	//exp, ok := claims["exp"].(float64)
	//if !ok {
	//	return eTime, errors.New("invalid expiration claim")
	//}
	//
	//// Compare expiration time with current time
	//expirationTime := time.Unix(int64(exp), 0)
	//return expirationTime, nil
}

//
//func isAccessJwtExpired(accessJwt string) (bool, error) {
//	// Parse the token
//	token, _, err := new(jwt.Parser).ParseUnverified(accessJwt, jwt.MapClaims{})
//	if err != nil {
//		return false, err
//	}
//
//	// Extract claims
//	claims, ok := token.Claims.(jwt.MapClaims)
//	if !ok {
//		return false, fmt.Errorf("invalid token claims")
//	}
//
//	// Get expiration time
//	exp, ok := claims["exp"].(float64)
//	if !ok {
//		return false, fmt.Errorf("invalid expiration claim")
//	}
//
//	// Compare expiration time with current time
//	expirationTime := time.Unix(int64(exp), 0)
//	return time.Now().After(expirationTime), nil
//}
