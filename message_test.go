package whatsmgr

import (
	"encoding/json"
	"testing"

	"go.mau.fi/whatsmeow/types/events"
)

var testMessages = []struct {
	note                string
	eventMessageJSON    []byte
	expectedMessageJSON []byte
}{
	{
		note:                "Normal incoming text message",
		eventMessageJSON:    []byte(`{"Info":{"AddressingMode":"","BroadcastListOwner":"","Category":"","Chat":"123456789@s.whatsapp.net","DeviceSentMeta":null,"Edit":"","ID":"12345678901234567890123456789012","IsFromMe":false,"IsGroup":false,"MediaType":"","MsgBotInfo":{"EditSenderTimestampMS":"0001-01-01T00:00:00Z","EditTargetID":"","EditType":""},"MsgMetaInfo":{"DeprecatedLIDSession":null,"TargetID":"","TargetSender":"","ThreadMessageID":"","ThreadMessageSenderJID":""},"Multicast":false,"PushName":"Jow Blow","RecipientAlt":"","Sender":"123456789@s.whatsapp.net","SenderAlt":"123456789012345@lid","ServerID":0,"Timestamp":"2025-05-02T11:13:28Z","Type":"text","VerifiedName":null},"IsDocumentWithCaption":false,"IsEdit":false,"IsEphemeral":false,"IsLottieSticker":false,"IsViewOnce":false,"IsViewOnceV2":false,"IsViewOnceV2Extension":false,"Message":{"conversation":"Hello world!","messageContextInfo":{"deviceListMetadata":{"recipientKeyHash":"U29tZUJhc2U2NEhhc2g=","recipientTimestamp":1745938978,"senderTimestamp":1743814938},"deviceListMetadataVersion":2,"messageSecret":"U29tZUJhc2U2NEhhc2g="}},"NewsletterMeta":null,"RawMessage":{"conversation":"Hello world!","messageContextInfo":{"deviceListMetadata":{"recipientKeyHash":"U29tZUJhc2U2NEhhc2g=","recipientTimestamp":1745938978,"senderTimestamp":1743814938},"deviceListMetadataVersion":2,"messageSecret":"U29tZUJhc2U2NEhhc2g="}},"RetryCount":0,"SourceWebMsg":null,"UnavailableRequestID":""}`),
		expectedMessageJSON: []byte(`{"Timestamp":"2025-05-02T11:13:28Z","MessageID":"12345678901234567890123456789012","ChatJID":"123456789@s.whatsapp.net","SenderJID":"123456789@s.whatsapp.net","IsFromMe":false,"Type":"text","ContentBody":"Hello world!"}`),
	},
	{
		note:                "Incoming text message quote",
		eventMessageJSON:    []byte(`{"Info":{"AddressingMode":"","BroadcastListOwner":"","Category":"","Chat":"123456789@s.whatsapp.net","DeviceSentMeta":null,"Edit":"","ID":"09876543210987654321098765432109","IsFromMe":false,"IsGroup":false,"MediaType":"","MsgBotInfo":{"EditSenderTimestampMS":"0001-01-01T00:00:00Z","EditTargetID":"","EditType":""},"MsgMetaInfo":{"DeprecatedLIDSession":null,"TargetID":"","TargetSender":"","ThreadMessageID":"","ThreadMessageSenderJID":""},"Multicast":false,"PushName":"Donovan Diamond","RecipientAlt":"","Sender":"123456789@s.whatsapp.net","SenderAlt":"123456789012345@lid","ServerID":0,"Timestamp":"2025-05-02T11:36:36Z","Type":"text","VerifiedName":null},"IsDocumentWithCaption":false,"IsEdit":false,"IsEphemeral":false,"IsLottieSticker":false,"IsViewOnce":false,"IsViewOnceV2":false,"IsViewOnceV2Extension":false,"Message":{"extendedTextMessage":{"contextInfo":{"participant":"123456789@s.whatsapp.net","quotedMessage":{"conversation":"Hello world!"},"stanzaID":"12345678901234567890123456789012"},"inviteLinkGroupTypeV2":0,"previewType":0,"text":"Well"},"messageContextInfo":{"deviceListMetadata":{"recipientKeyHash":"U29tZUJhc2U2NEhhc2g=","recipientTimestamp":1745938978,"senderTimestamp":1743814938},"deviceListMetadataVersion":2,"messageSecret":"U29tZUJhc2U2NEhhc2g="}},"NewsletterMeta":null,"RawMessage":{"extendedTextMessage":{"contextInfo":{"participant":"123456789@s.whatsapp.net","quotedMessage":{"conversation":"Hello world!"},"stanzaID":"12345678901234567890123456789012"},"inviteLinkGroupTypeV2":0,"previewType":0,"text":"Well"},"messageContextInfo":{"deviceListMetadata":{"recipientKeyHash":"U29tZUJhc2U2NEhhc2g=","recipientTimestamp":1745938978,"senderTimestamp":1743814938},"deviceListMetadataVersion":2,"messageSecret":"U29tZUJhc2U2NEhhc2g="}},"RetryCount":0,"SourceWebMsg":null,"UnavailableRequestID":""}`),
		expectedMessageJSON: []byte(`{"Timestamp":"2025-05-02T11:36:36Z","MessageID":"09876543210987654321098765432109","ChatJID":"123456789@s.whatsapp.net","SenderJID":"123456789@s.whatsapp.net","IsFromMe":false,"Type":"text","ContentBody":"Well","InfoQuotedMessageID":"12345678901234567890123456789012"}`),
	},
}

func TestParseEventMessage(t *testing.T) {
	for i, test := range testMessages {
		var eventMessage events.Message
		if err := json.Unmarshal(test.eventMessageJSON, &eventMessage); err != nil {
			panic(err)
		}
		outputMessage := (&Connection{}).parseEventMessage(eventMessage)
		outputMessage.Raw = nil
		outputMessageJSON, err := json.Marshal(outputMessage)
		if err != nil {
			t.Error(err)
		}
		if string(outputMessageJSON) != string(test.expectedMessageJSON) {
			t.Fatalf("Test #%d (%s): output message not equal to expected message:\n\nEXPECTED:\n%s\n\nGOT:\n%s\n\n", i, test.note, test.expectedMessageJSON, outputMessageJSON)
		}
	}
}
