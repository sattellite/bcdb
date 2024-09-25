package network

import (
	"io"

	"github.com/sattellite/bcdb/compute/result"
)

func (n *Network) Print(w io.Writer, res result.Result) error {
	_, err := w.Write(append(res.Bytes(), '\n'))
	return err
}
