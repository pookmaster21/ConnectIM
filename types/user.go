package types

type Prefered int

const (
	TELEGRAM Prefered = 1 << iota
	DISCORD
)

type User struct {
	Username string   `json:"username"`
	Password string   `json:"password"`
	Telegram string   `json:"telegram"`
	Whatsapp string   `json:"whatsapp"`
	Discord  string   `json:"discord"`
	Prefered Prefered `json:"prefered"`
}
