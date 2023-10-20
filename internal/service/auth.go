package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/patrickmn/go-cache"
	"github.com/zelas91/gofermart/internal/entities"
	"github.com/zelas91/gofermart/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var secret = []byte("secret_key")

type AuthService struct {
	repo  repository.Authorization
	cache *cache.Cache
}

type Claims struct {
	jwt.RegisteredClaims
	Login string
}

func generateJwt(login string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 1)),
		},
		Login: login,
	})
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func newAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo, cache: cache.New(time.Minute*10, time.Minute*10)}
}

func (a *AuthService) CreateUser(ctx context.Context, user *entities.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("generate password hash err")
	}

	return a.repo.CreateUser(ctx, user.Login, string(hashedPassword))
}

func (a *AuthService) CreateToken(ctx context.Context, authUser *entities.User) (string, error) {
	user, err := a.repo.GetUser(ctx, authUser)
	if err != nil {
		return "", err
	}
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(authUser.Password)); err != nil {
		return "", err
	}
	token, err := generateJwt(user.Login)
	if err != nil {
		return "", err
	}
	a.cache.Set(token, user, cache.DefaultExpiration)
	return token, err
}

func (a *AuthService) ParserToken(ctx context.Context, tokenString string) (*entities.User, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("error parsing jwt")
		}
		return secret, nil
	})
	if err != nil && token == nil {
		return nil, err
	}
	if !token.Valid && !time.Now().Before(claims.ExpiresAt.Time) {
		return nil, errors.New("token not valid")
	}
	val, ok := a.cache.Get(tokenString)

	if ok {
		user := val.(entities.User)
		return &user, nil
	}
	user, err := a.repo.GetUser(ctx, &entities.User{Login: claims.Login})
	if err != nil {
		return nil, err
	}
	a.cache.Set(tokenString, user, cache.DefaultExpiration)
	return &user, nil
}
