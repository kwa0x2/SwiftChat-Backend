package types

type RequestStatus string

const (
	Pending  RequestStatus = "pending"
	Rejected RequestStatus = "rejected"
	Accepted RequestStatus = "accepted"
)
