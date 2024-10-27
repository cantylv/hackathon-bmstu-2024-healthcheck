package entity

type JwtTokenHeader struct {
	Exp string `json:"exp"`
}

type JwtTokenPayload struct {
	Id string `json:"id"`
}
