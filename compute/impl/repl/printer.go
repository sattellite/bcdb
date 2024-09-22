package repl

import (
	"io"

	"github.com/sattellite/bcdb/compute/result"
)

var (
	prefixIn  = []byte("> ")
	prefixOut = []byte("< ")
)

func (r *REPL) Print(w io.Writer, res result.Result) error {
	perr := r.prompt(w, prefixOut)
	if perr != nil {
		return perr
	}
	_, err := w.Write(append(res.Bytes(), '\n'))
	return err
}

func (r *REPL) prompt(w io.Writer, p []byte) error {
	_, err := w.Write(p)
	return err
}
