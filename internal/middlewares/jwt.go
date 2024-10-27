// Copyright © ivanlobanov. All rights reserved.
package middlewares

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/cantylv/hackathon-bmstu-2024-healthcheck/internal/entity/dto"
	f "github.com/cantylv/hackathon-bmstu-2024-healthcheck/internal/utils/functions"
	mc "github.com/cantylv/hackathon-bmstu-2024-healthcheck/internal/utils/myconstants"
	"github.com/cantylv/hackathon-bmstu-2024-healthcheck/internal/utils/myerrors"
	me "github.com/cantylv/hackathon-bmstu-2024-healthcheck/internal/utils/myerrors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// JWT --> header.payload.signature
// header --> base64(meta_information)
// payload --> base64(payload_data)
// signature --> hmacsha256(header + . + payload + secret)

//// e.g. header
// {
// 	"exp": "02.01.2006 15:04:05 UTC-07"
// }
//// e.g. payload
// {
// 	"username": "66b89cea43ad0d6f8cf3f54e",
// }

// JwtVerification
// Needed for authentication.
func JwtVerification(h http.Handler, logger *zap.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID, err := f.GetCtxRequestID(r)
		if err != nil {
			logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
		}

		jwtToken, err := f.GetJWtToken(r)
		if err != nil && !errors.Is(err, http.ErrNoCookie) {
			logger.Error(fmt.Sprintf("error while jwt getting: %v", err), zap.String(mc.RequestID, requestID))
			f.Response(w, dto.ResponseError{Error: me.ErrInternal.Error()}, http.StatusInternalServerError)
			return
		}
		if jwtToken != "" {
			username, err := jwtTokenIsValid(jwtToken)
			if err != nil {
				logger.Error(fmt.Sprintf("error while jwt verification: %v", err), zap.String(mc.RequestID, requestID))
				f.Response(w, dto.ResponseError{Error: me.ErrInvalidJwt.Error()}, http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), "username", username)
			r = r.WithContext(ctx)
		}
		// Decode payload and use data.
		h.ServeHTTP(w, r)
	})
}

// jwtTokenIsValid
// Needed for validation jwt-token.
func jwtTokenIsValid(token string) (string, error) {
	// check time validation of token
	// if all is okey, return true
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return "", myerrors.ErrInvalidJwt
	}
	signatureHash, err := hashWithStatement(parts[0] + "." + parts[1]) // header + "." + payload)
	if err != nil {
		return "", err
	}
	signature := hex.EncodeToString([]byte(signatureHash))
	if signature != parts[2] {
		return "", myerrors.ErrInvalidJwt
	}

	dataHeader, err := hex.DecodeString(parts[0])
	if err != nil {
		return "", err
	}
	var h dto.JwtTokenHeader
	err = json.Unmarshal(dataHeader, &h)
	if err != nil {
		return "", err
	}

	dataPayload, err := hex.DecodeString(parts[1])
	if err != nil {
		return "", err
	}
	var p dto.JwtTokenPayload
	err = json.Unmarshal(dataPayload, &p)
	if err != nil {
		return "", err
	}

	// "02.01.2006 15:04:05 UTC-07" template
	jwtDate, err := time.Parse("02.01.2006 15:04:05 UTC-07", h.Exp)
	if err != nil {
		return "", err
	}
	dateNow := time.Now()
	if jwtDate.Equal(dateNow) || dateNow.After(jwtDate) {
		return "", myerrors.ErrInvalidJwt
	}
	return p.Username, nil
}

// hashWithStatement
// Returns hash that is transmitted in the client-server model by custom header.
func hashWithStatement(statement string) (string, error) {
	secretKey := viper.GetString("secret_key")
	mac := hmac.New(sha256.New, []byte(secretKey))
	_, err := mac.Write([]byte(statement))
	if err != nil {
		return "", err
	}
	return string(mac.Sum(nil)), nil
}
