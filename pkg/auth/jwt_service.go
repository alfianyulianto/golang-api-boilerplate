package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/alfianyulianto/pds-service/internal/model"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"time"
)

type JWTConfig struct {
	AppName               string
	ExpireDuration        int
	RefreshExpireDuration int
	SecreteKey            string
}

type JWTService struct {
	*JWTConfig
	Redis *redis.Client
}

func NewJWTService(JWTConfig *JWTConfig, redis *redis.Client) *JWTService {
	return &JWTService{JWTConfig: JWTConfig, Redis: redis}
}

func (s *JWTService) CreateToken(ctx context.Context, claims *model.UserClaimToken) (*model.AuthResponse, error) {
	claims.Type = "access"
	claims.RegisteredClaims = jwt.RegisteredClaims{
		Issuer:    s.AppName,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(s.ExpireDuration) * time.Second)),
		ID:        uuid.NewString(), // tujuaannya untuk membedakan token satu dengan yang lain, sehingga bisa di revoke satu per satu, multi device bisa
	}

	mySigningKey := []byte(s.SecreteKey)
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedAccessToken, err := accessToken.SignedString(mySigningKey)
	if err != nil {
		return nil, err
	}
	key := fmt.Sprintf("access_token:%s", claims.RegisteredClaims.ID)
	if err = s.Redis.SetEx(ctx, key, claims.ID, time.Duration(s.ExpireDuration)*time.Second).Err(); err != nil {
		return nil, err
	}

	claims.Type = "refresh"
	claims.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Duration(s.RefreshExpireDuration) * time.Second))
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedRefreshToken, err := refreshToken.SignedString(mySigningKey)
	if err != nil {
		return nil, err
	}
	key = fmt.Sprintf("refresh_token:%s", claims.RegisteredClaims.ID)
	if err = s.Redis.SetEx(ctx, key, claims.ID, time.Duration(s.RefreshExpireDuration)*time.Second).Err(); err != nil {
		return nil, err
	}

	return &model.AuthResponse{
		AccessToken:  signedAccessToken,
		RefreshToken: signedRefreshToken,
	}, nil
}

func (s *JWTService) ParseAccessToken(ctx context.Context, accessToken string) (*model.UserClaimToken, error) {
	token, err := jwt.ParseWithClaims(accessToken, new(model.UserClaimToken), func(token *jwt.Token) (interface{}, error) {
		return []byte(s.SecreteKey), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*model.UserClaimToken)
	if !ok || !token.Valid || claims.Type != "access" {
		return nil, errors.New("Invalid access token")
	}

	key := fmt.Sprintf("access_token:%s", claims.RegisteredClaims.ID)
	result, err := s.Redis.Exists(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	if result == 0 {
		return nil, errors.New("Invalid access token, not found in redis")
	}

	return claims, nil
}

func (s *JWTService) ParseRefreshToken(ctx context.Context, refreshToken string) (*model.UserClaimToken, error) {
	token, err := jwt.ParseWithClaims(refreshToken, new(model.UserClaimToken), func(token *jwt.Token) (interface{}, error) {
		return []byte(s.SecreteKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*model.UserClaimToken)
	if !ok || !token.Valid || claims.Type != "refresh" {
		return nil, errors.New("Invalid refresh token")
	}

	key := fmt.Sprintf("refresh_token:%s", claims.RegisteredClaims.ID)
	result, err := s.Redis.Exists(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	if result == 0 {
		return nil, errors.New("Invalid refresh token, not found in redis")
	}

	key = fmt.Sprintf("access_token:%s", claims.RegisteredClaims.ID)
	if err = s.Redis.Del(ctx, key).Err(); err != nil {
		return nil, err
	}

	key = fmt.Sprintf("refresh_token:%s", claims.RegisteredClaims.ID)
	if err = s.Redis.Del(ctx, key).Err(); err != nil {
		return nil, err
	}

	return claims, nil
}

func (s *JWTService) RevokeToken(ctx context.Context, claims *model.UserClaimToken) error {
	keys := s.Redis.Keys(ctx, fmt.Sprintf("*_token:%s", claims.RegisteredClaims.ID)).Val()
	if len(keys) > 0 {
		if err := s.Redis.Del(ctx, keys...).Err(); err != nil {
			return err
		}
	}

	return nil
}
