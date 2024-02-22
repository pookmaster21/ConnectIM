package types

type Prefered uint8

const (
	TELEGRAM Prefered = iota
	WHATSAPP
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
