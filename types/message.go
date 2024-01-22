package types

type Message struct {
	From string
	Msg  string `json:"msg"`
	To   string
}
