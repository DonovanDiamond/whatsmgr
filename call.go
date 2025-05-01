package whatsmgr

import "time"

type Call struct {
	Timestamp   time.Time
	CallID      string `json:",omitempty"`
	From        string `json:",omitempty"`
	CallCreator string `json:",omitempty"`

	Status CallStatus `json:",omitempty"` // use CallStatus* constants

	Media *CallMedia `json:",omitempty"` // use CallMedia* constants, not always set
	Type  *CallType  `json:",omitempty"` // use CallType* constants, not always set

	TerminateReaason *string `json:",omitempty"`
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
