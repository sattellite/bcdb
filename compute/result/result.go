package result

type Result struct {
	Value string
	Error error
}

func (r *Result) Bytes() []byte {
	if r.Error != nil {
		return []byte("ERR: " + r.Error.Error())
	}
	return []byte("RES: " + r.Value)
}
