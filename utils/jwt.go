package utils

import (
	"time"
	"strconv"
	"github.com/golang-jwt/jwt/v5"
	pgModel "go_clean/app/models/postgresql"
    mongoModel "go_clean/app/models/mongodb"

	"go_clean/config"
)

var MockGenerateToken func(u pgModel.User) (string, error)

func GenerateToken(u pgModel.User) (string, error) {

    if MockGenerateToken != nil {
        return MockGenerateToken(u)
    }

    jwtCfg := config.LoadJWT()
    claims := pgModel.JWTClaims{
        UserID:   strconv.Itoa(u.ID),   // Convert int → string
        Username: u.Username,
        Role:     u.Role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(jwtCfg.TTLHours) * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }

    tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return tok.SignedString(jwtCfg.Secret)
}

var MockGenerateTokenMongo func(u mongoModel.LoginMongo) (string, error)

func GenerateTokenMongo(u mongoModel.LoginMongo) (string, error) {

     if MockGenerateTokenMongo != nil {
        return MockGenerateTokenMongo(u)
    }

    jwtCfg := config.LoadJWT()

    claims := pgModel.JWTClaims{
        UserID:   u.ID.Hex(),
        Username: u.Username,
        Role:     u.Role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(jwtCfg.TTLHours) * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }

    tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return tok.SignedString(jwtCfg.Secret)
}

func ValidateToken(tokenStr string) (*pgModel.JWTClaims, error) {
	jwtCfg := config.LoadJWT()
	tok, err := jwt.ParseWithClaims(tokenStr, &pgModel.JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
		return jwtCfg.Secret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := tok.Claims.(*pgModel.JWTClaims); ok && tok.Valid {
		return claims, nil
	}
	return nil, jwt.ErrTokenInvalidClaims
}




