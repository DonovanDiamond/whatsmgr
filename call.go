package whatsmgr

import "time"

type Call struct {
	Timestamp   time.Time
	CallID      string
	From        string
	CallCreator string

	Status CallStatus // use CallStatus* constants

	Media *CallMedia // use CallMedia* constants, not always set
	Type  *CallType  // use CallType* constants, not always set

	TerminateReaason *string
}

type CallStatus string
type CallMedia string
type CallType string

const (
	CallStatusOffer     CallStatus = "offer"
	CallStatusAccept    CallStatus = "accept"
	CallStatusPreAccept CallStatus = "pre-accept"
	CallStatusTransport CallStatus = "transport"
	CallStatusTerminate CallStatus = "terminate"
	CallStatusReject    CallStatus = "reject"

	CallMediaAudio CallMedia = "audio"
	CallMediaVideo CallMedia = "video"

	CallTypeGroup CallType = "group"
)
