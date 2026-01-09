package tempdir

import "errors"

var (
	ErrAlreadyDeleted = errors.New("tempdir: directory already deleted")
)
