package functions

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"

	"github.com/cantylv/hackathon-bmstu-2024-healthcheck/internal/entity/dto"
	mc "github.com/cantylv/hackathon-bmstu-2024-healthcheck/internal/utils/myconstants"
	"github.com/spf13/viper"
)

type NewJwtTokenProps struct {
	Username string
}

// HTTP Headers "Cookie"
func GetJWtToken(r *http.Request) (string, error) {
	jwtCookie, err := r.Cookie(mc.JwtToken)
	if err != nil {
		return "", err
	}
	return jwtCookie.Value, nil
}

// NewCsrfToken
// Generates jwt-token.
func NewJwtToken(props NewJwtTokenProps) (string, error) {
	// Encode header.
	h := dto.JwtTokenHeader{
		Exp: time.Now().Format("02.01.2006 15:04:05 UTC-07"),
	}
	rawDataHeader, err := json.Marshal(h)
	if err != nil {
		return "", err
	}
	hEncoded := hex.EncodeToString(rawDataHeader)
	// Encode payload.
	p := dto.JwtTokenPayload{
		Username: props.Username,
	}
	rawDataPayload, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	pEncoded := hex.EncodeToString(rawDataPayload)
	// concatenate header and payload
	hpEncoded := hEncoded + "." + pEncoded
	signatureHash, err := hashWithStatement(hpEncoded)
	if err != nil {
		return "", err
	}
	signature := hex.EncodeToString([]byte(signatureHash))
	return hpEncoded + "." + signature, nil
}

// HashWithStatement
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
