package internal

type HTTPError struct {
	Status int
	Msg    string
}

func (e HTTPError) Error() string {
	return e.Msg
}
