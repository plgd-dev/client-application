package grpc

import (
	"context"
	"errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/plgd-dev/client-application/pb"
	"github.com/plgd-dev/hub/v2/identity-store/events"
	plgdJwt "github.com/plgd-dev/hub/v2/pkg/security/jwt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type JSONWebKeyCache struct {
	owner uuid.UUID
	keys  jwk.Set
}

func NewJSONWebKeyCache(owner uuid.UUID, keys jwk.Set) *JSONWebKeyCache {
	return &JSONWebKeyCache{
		owner: owner,
		keys:  keys,
	}
}

func (c *JSONWebKeyCache) GetKey(token *jwt.Token) (interface{}, error) {
	key, err := c.LookupKey(token)
	if err != nil {
		return nil, err
	}
	var v interface{}
	return v, key.Raw(&v)
}

func (c *JSONWebKeyCache) LookupKey(token *jwt.Token) (jwk.Key, error) {
	id, ok := token.Header["kid"].(string)
	if !ok {
		return nil, errors.New("missing key id in token")
	}

	if c.keys == nil {
		return nil, errors.New("empty JWK cache")
	}
	if key, ok := c.keys.LookupKeyID(id); ok {
		if key.Algorithm().String() == token.Method.Alg() {
			return key, nil
		}
	}
	return nil, errors.New("could not find JWK")
}

func (s *ClientApplicationServer) ParseWithClaims(_ context.Context, token string, claims jwt.Claims) error {
	if token == "" {
		if !s.HasJWTAuthorizationEnabled() {
			return nil
		}
		return status.Errorf(codes.Unauthenticated, "missing token")
	}
	c := s.jwksCache.Load()
	if c == nil {
		return status.Errorf(codes.Unauthenticated, "cannot validate token: missing JWK cache")
	}
	_, err := jwt.ParseWithClaims(token, claims, c.GetKey)
	if err != nil {
		return status.Errorf(codes.Unauthenticated, "could not parse token: %v", err)
	}

	scopeClaims, ok := claims.(*plgdJwt.ScopeClaims)
	if !ok {
		return status.Errorf(codes.Unauthenticated, "invalid type of token claims %T", claims)
	}
	plgdClaims := plgdJwt.Claims(*scopeClaims)
	cfg := s.GetConfig()
	owner, err := plgdClaims.GetOwner(cfg.RemoteProvisioning.GetJwtOwnerClaim())
	if err != nil {
		return status.Errorf(codes.Unauthenticated, "cannot get owner claim: %v", err)
	}
	if owner == "" {
		return status.Errorf(codes.Unauthenticated, "owner claim is not set")
	}
	ownerID, err := uuid.Parse(events.OwnerToUUID(owner))
	if owner == "" {
		return status.Errorf(codes.Unauthenticated, "cannot parse owner claim to UUID: %v", err)
	}
	if ownerID != c.owner {
		return status.Errorf(codes.Unauthenticated, "unexpected owner('%v')", owner)
	}

	return nil
}

func (s *ClientApplicationServer) HasJWTAuthorizationEnabled() bool {
	devService := s.serviceDevice.Load()
	if devService == nil {
		return false
	}
	return devService.GetDeviceAuthenticationMode() == pb.GetConfigurationResponse_X509
}
