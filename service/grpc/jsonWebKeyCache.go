package grpc

import (
	"fmt"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/plgd-dev/client-application/pb"
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
	c := s.jwksCache.Load()
	if c == nil {
		return status.Errorf(codes.Unauthenticated, "cannot validate token: missing JWK cache")
	}
	if token == "" {
		return status.Errorf(codes.Unauthenticated, "missing token")
	}

	_, err := jwt.ParseWithClaims(token, claims, c.GetKey)
	if err != nil {
		return status.Errorf(codes.Unauthenticated, "could not parse token: %v", err)
	}
	return nil
}

func (s *ClientApplicationServer) HasJWTAuthorizationEnabled() bool {
	return s.serviceDevice.GetDeviceAuthenticationMode() == pb.GetConfigurationResponse_X509
}