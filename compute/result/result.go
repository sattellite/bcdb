package result

type Result struct {
	Value string
}

func (r *Result) Bytes() []byte {
	return []byte(r.Value)
}
