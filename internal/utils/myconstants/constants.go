package myconstants

type AccessKey string

// Частые переменные
const (
	RequestID = "request_id"
	XRealIP   = "X-Real-IP"
	JwtToken  = "jwt-token"
)

// Настройка хэширования с помощью Argon2
const (
	HashTime    = 1
	HashMemory  = 2 * 1024
	HashThreads = 2
	HashKeylen  = 56
	HashLetters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
)

var AllowedActivities = map[string]float32{
	"NFA": 1.2,
	"LA":  1.375,
	"MA":  1.55,
	"HA":  1.725,
	"EA":  1.9,
}

var AllowedHumanSex = map[string]struct{}{
	"F": {},
	"M": {},
}
