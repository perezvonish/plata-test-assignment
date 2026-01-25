package job

type Status string

var (
	StatusPending    Status = "pending"
	StatusProcessing Status = "processing"
	StatusSuccess    Status = "success"
	StatusFailure    Status = "failure"
)
