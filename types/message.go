package types

type Message struct {
	From *User
	Msg  string
	To   *User
}
