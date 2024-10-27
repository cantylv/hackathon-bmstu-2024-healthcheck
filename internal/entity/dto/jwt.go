package dto

type JwtTokenHeader struct {
	Exp string `json:"exp"`
}

type JwtTokenPayload struct {
	Username string `json:"username"`
}
