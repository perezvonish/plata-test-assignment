package quote

import "errors"

var (
	ErrorWhileFindingJob = errors.New("while finding job")
	ErrorNotFoundJob     = errors.New("job not found")
)
