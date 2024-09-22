package repl

import "github.com/sattellite/bcdb/compute/result"

var (
	prefixIn  = []byte("> ")
	prefixOut = []byte("< ")
)

func (r *REPL) Print(res result.Result) error {
	perr := r.prompt(prefixOut)
	if perr != nil {
		return perr
	}
	_, err := r.out.Write(append(res.Bytes(), '\n'))
	return err
}

func (r *REPL) prompt(p []byte) error {
	_, err := r.out.Write(p)
	return err
}
