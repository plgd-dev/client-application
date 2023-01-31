package grpc

import (
	"fmt"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/jwk"
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
		return nil, fmt.Errorf("missing key id in token")
	}

	if c.keys == nil {
		return nil, fmt.Errorf("empty JWK cache")
	}
	if key, ok := c.keys.LookupKeyID(id); ok {
		if key.Algorithm() == token.Method.Alg() {
			return key, nil
		}
	}
	return nil, fmt.Errorf("could not find JWK")
}

func (s *ClientApplicationServer) ParseWithClaims(token string, claims jwt.Claims) error {
	if token == "" {
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
	owner := plgdClaims.Owner(cfg.RemoteProvisioning.GetJwtOwnerClaim())
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
	return s.serviceDevice.GetDeviceAuthenticationMode() == pb.GetConfigurationResponse_X509
}
