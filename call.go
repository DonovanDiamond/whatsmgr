package whatsmgr

import "time"

type Call struct {
	Timestamp   time.Time
	CallID      string
	From        string
	CallCreator string

	Status string // use CallStatus* constants

	Media *string // use CallMedia* constants, not always set
	Type  *string // use CallType* constants, not always set

	TerminateReaason *string
}

const (
	CallStatusOffer     = "offer"
	CallStatusAccept    = "accept"
	CallStatusPreAccept = "pre-accept"
	CallStatusTransport = "transport"
	CallStatusTerminate = "terminate"
	CallStatusReject    = "reject"

	CallMediaAudio = "audio"
	CallMediaVideo = "video"

	CallTypeGroup = "group"
)
