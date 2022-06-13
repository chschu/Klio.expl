package util

import (
	"github.com/hashicorp/go-multierror"
	"io"
)

func CloseAndAppendError(c io.Closer, err *error) {
	appendError(c.Close(), err)
}

func appendError(errToAppend error, err *error) {
	if errToAppend != nil {
		if *err != nil {
			*err = multierror.Append(*err, errToAppend)
		} else {
			*err = errToAppend
		}
	}
}
