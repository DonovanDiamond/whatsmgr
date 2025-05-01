package whatsmgr

import (
	"fmt"
	"mime"
	"time"

	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types/events"
)

type Message struct {
	Timestamp *time.Time `json:",omitempty"`

	MessageID string  `json:",omitempty"`
	ChatJID   string  `json:",omitempty"`
	SenderJID *string `json:",omitempty"`

	IsFromMe *bool          `json:",omitempty"`
	Type     *string        `json:",omitempty"` // not always set
	Status   *MessageStatus `json:",omitempty"` // use MessageStatus* constants, not always set

	Starred *bool `json:",omitempty"` // not always set
	Deleted *bool `json:",omitempty"` // not always set

	ContentBody *string `json:",omitempty"` // not always set

	InfoQuotedMessageID *string `json:",omitempty"` // not always set
	InfoParticipant     *string `json:",omitempty"` // not always set
	InfoRemoteJID       *string `json:",omitempty"` // not always set

	Attachments []string `json:",omitempty"` // not always set

	ContactVcard       *string `json:",omitempty"` // not always set
	ContactDisplayName *string `json:",omitempty"` // not always set

	LocationLat              *float64 `json:",omitempty"` // not always set
	LocationLon              *float64 `json:",omitempty"` // not always set
	LocationName             *string  `json:",omitempty"` // not always set
	LocationAddress          *string  `json:",omitempty"` // not always set
	LocationURL              *string  `json:",omitempty"` // not always set
	LocationIsLive           *bool    `json:",omitempty"` // not always set
	LocationAccuracyInMeters *uint32  `json:",omitempty"` // not always set
	LocationComment          *string  `json:",omitempty"` // not always set

	CallLogOutcome         *CallLogOutcome `json:",omitempty"` // use CallLogOutcome* constants, not always set
	CallLogDurationSeconds *int64          `json:",omitempty"` // not always set
	CallLogType            *CallLogType    `json:",omitempty"` // use CallLogType* constants, not always set
	CallLogParticipantJIDs []string        `json:",omitempty"` // they have their own call outcomes that are not recorded, not always set
}

type MessageStatus string
type CallLogOutcome string
type CallLogType string

const (
	MessageStatusSent        MessageStatus = "sent"
	MessageStatusDelivered   MessageStatus = "delivered"
	MessageStatusRead        MessageStatus = "read"
	MessageStatusServerError MessageStatus = "server-error"

	CallLogOutcomeConnected CallLogOutcome = "connected"
	CallLogOutcomeMissed    CallLogOutcome = "missed"
	CallLogOutcomeFailed    CallLogOutcome = "failed"
	CallLogOutcomeRejected  CallLogOutcome = "rejected"
	CallLogOutcomeAccepted  CallLogOutcome = "accepted"
	CallLogOutcomeOngoing   CallLogOutcome = "ongoing"
	CallLogOutcomeSilenced  CallLogOutcome = "silenced"

	CallLogTypeRegular   CallLogType = "regular"
	CallLogTypeScheduled CallLogType = "scheduled"
	CallLogTypeVoiceChat CallLogType = "voice-chat"
)

func (conn *Connection) pullAttachments(m events.Message) (attachments []string, caption string, err error) {
	if att := m.Message.GetImageMessage(); att != nil {
		caption = att.GetCaption()
		raw, err := conn.client.Download(att)
		if err != nil {
			return attachments, caption, fmt.Errorf("failed to download attachment: %w", err)
		}
		if len(raw) > 0 {
			fileName := conn.hashFile(raw)
			exts, _ := mime.ExtensionsByType(att.GetMimetype())
			if len(exts) > 0 {
				fileName += exts[0]
			}
			path := fmt.Sprintf("%s/%s", conn.MediaPath, fileName)
			if err := conn.writeFileIfNotExists(path, raw); err != nil {
				return attachments, caption, fmt.Errorf("failed to write attachment: %w", err)
			}
			attachments = append(attachments, fileName)
		}
	}
	if att := m.Message.GetAudioMessage(); att != nil {
		raw, err := conn.client.Download(att)
		if err != nil {
			return attachments, caption, fmt.Errorf("failed to download attachment: %w", err)
		}
		if len(raw) > 0 {
			fileName := conn.hashFile(raw)
			exts, _ := mime.ExtensionsByType(att.GetMimetype())
			if len(exts) > 0 {
				fileName += exts[0]
			} else {
				fileName += ".ogg"
			}
			path := fmt.Sprintf("%s/%s", conn.MediaPath, fileName)
			if err := conn.writeFileIfNotExists(path, raw); err != nil {
				return attachments, caption, fmt.Errorf("failed to write attachment: %w", err)
			}
			attachments = append(attachments, fileName)
		}
	}
	if att := m.Message.GetVideoMessage(); att != nil {
		caption = att.GetCaption()
		raw, err := conn.client.Download(att)
		if err != nil {
			return attachments, caption, fmt.Errorf("failed to download attachment: %w", err)
		}
		if len(raw) > 0 {
			fileName := conn.hashFile(raw)
			exts, _ := mime.ExtensionsByType(att.GetMimetype())
			if len(exts) > 0 {
				fileName += exts[0]
			} else {
				fileName += ".mp4"
			}
			path := fmt.Sprintf("%s/%s", conn.MediaPath, fileName)
			if err := conn.writeFileIfNotExists(path, raw); err != nil {
				return attachments, caption, fmt.Errorf("failed to write attachment: %w", err)
			}
			attachments = append(attachments, fileName)
		}
	}
	if att := m.Message.GetDocumentMessage(); att != nil {
		caption = att.GetCaption()
		raw, err := conn.client.Download(att)
		if err != nil {
			return attachments, caption, fmt.Errorf("failed to download attachment: %w", err)
		}
		if len(raw) > 0 {
			fileName := conn.hashFile(raw)
			exts, _ := mime.ExtensionsByType(att.GetMimetype())
			if len(exts) > 0 {
				fileName += exts[0]
			}
			path := fmt.Sprintf("%s/%s", conn.MediaPath, fileName)
			if err := conn.writeFileIfNotExists(path, raw); err != nil {
				return attachments, caption, fmt.Errorf("failed to write attachment: %w", err)
			}
			attachments = append(attachments, fileName)
		}
	}
	if att := m.Message.GetStickerMessage(); att != nil {
		raw, err := conn.client.Download(att)
		if err != nil {
			return attachments, caption, fmt.Errorf("failed to download attachment: %w", err)
		}
		if len(raw) > 0 {
			fileName := conn.hashFile(raw)
			exts, _ := mime.ExtensionsByType(att.GetMimetype())
			if len(exts) > 0 {
				fileName += exts[0]
			}
			path := fmt.Sprintf("%s/%s", conn.MediaPath, fileName)
			if err := conn.writeFileIfNotExists(path, raw); err != nil {
				return attachments, caption, fmt.Errorf("failed to write attachment: %w", err)
			}
			attachments = append(attachments, fileName)
		}
	}
	return
}

func (conn *Connection) handleMessage(m events.Message) {
	log := conn.Log
	if m.Message == nil {
		return
	}
	sender := m.Info.Sender.String()
	message := Message{
		Timestamp: &m.Info.Timestamp,

		MessageID: m.Info.ID,
		ChatJID:   m.Info.Chat.String(),
		SenderJID: &sender,

		IsFromMe: &m.Info.IsFromMe,

		Type: &m.Info.Type,
	}

	attachments, caption, err := conn.pullAttachments(m)
	if err != nil {
		conn.Log.Error().Err(err).Msg("failed to pull attachment")
	}
	message.ContentBody = &caption
	message.Attachments = attachments

	if x := m.Message.GetConversation(); x != "" {
		message.ContentBody = &x
	}
	if x := m.Message.GetSenderKeyDistributionMessage(); x != nil {
		log.Warn().Any("x", x).Msg("NOT IMPLEMENTED: Message.GetSenderKeyDistributionMessage()")
	}
	if x := m.Message.GetContactMessage(); x != nil {
		message.ContactVcard = x.Vcard
		message.ContactDisplayName = x.DisplayName
	}
	if x := m.Message.GetLocationMessage(); x != nil {
		message.LocationAccuracyInMeters = x.AccuracyInMeters
		message.LocationAddress = x.Address
		message.LocationComment = x.Comment
		message.LocationIsLive = x.IsLive
		message.LocationLat = x.DegreesLatitude
		message.LocationLon = x.DegreesLongitude
		message.LocationName = x.Name
		message.LocationURL = x.URL
	}
	if x := m.Message.GetExtendedTextMessage(); x != nil {
		if len(x.JPEGThumbnail) > 0 {
			fileName := conn.hashFile(x.JPEGThumbnail) + ".jpeg"
			path := fmt.Sprintf("%s/%s", conn.MediaPath, fileName)
			err := conn.writeFileIfNotExists(path, x.JPEGThumbnail)
			if err == nil {
				message.Attachments = append(message.Attachments, fileName)
			}
		}
		if x.Text != nil {
			message.ContentBody = x.Text
		}
		if x.ContextInfo != nil {
			// StanzaID is the ID of the message being quoted as per https://github.com/tulir/whatsmeow/issues/88
			if x.ContextInfo.StanzaID != nil {
				message.InfoQuotedMessageID = x.ContextInfo.StanzaID
			}
		}
		// TODO: there is a lot more that could be implemented here...
	}
	if x := m.Message.GetLiveLocationMessage(); x != nil {
		message.LocationAccuracyInMeters = x.AccuracyInMeters
		message.LocationComment = x.Caption
		t := true
		message.LocationIsLive = &t
		message.LocationLat = x.DegreesLatitude
		message.LocationLon = x.DegreesLongitude
	}
	if x := m.Message.GetStickerMessage(); x != nil {
		if len(x.PngThumbnail) > 0 {
			fileName := conn.hashFile(x.PngThumbnail) + ".png"
			path := fmt.Sprintf("%s/%s", conn.MediaPath, fileName)
			err := conn.writeFileIfNotExists(path, x.PngThumbnail)
			if err == nil {
				message.Attachments = append(message.Attachments, fileName)
			}
		}
	}
	if x := m.Message.GetGroupInviteMessage(); x != nil {
		invite := "**Group Invite**\n"
		if x.GroupJID != nil {
			invite += "\nGroupJID: " + *x.GroupJID
		}
		if x.InviteCode != nil {
			invite += "\nInviteCode: " + *x.InviteCode
		}
		if x.InviteExpiration != nil {
			invite += "\nInviteExpiration: " + fmt.Sprint(*x.InviteExpiration)
		}
		if x.GroupName != nil {
			invite += "\nGroupName: " + *x.GroupName
		}
		if x.Caption != nil {
			invite += "\nCaption: " + *x.Caption
		}
		message.ContentBody = &invite
		if len(x.JPEGThumbnail) > 0 {
			fileName := conn.hashFile(x.JPEGThumbnail) + ".jpeg"
			path := fmt.Sprintf("%s/%s", conn.MediaPath, fileName)
			err := conn.writeFileIfNotExists(path, x.JPEGThumbnail)
			if err == nil {
				message.Attachments = append(message.Attachments, fileName)
			}
		}
	}
	if x := m.Message.GetReactionMessage(); x != nil {
		message.ContentBody = x.Text
		message.InfoQuotedMessageID = x.Key.ID
		message.InfoRemoteJID = x.Key.RemoteJID
		message.InfoParticipant = x.Key.Participant
	}
	if x := m.Message.GetCallLogMesssage(); x != nil {
		if x.CallOutcome != nil {
			var outcome CallLogOutcome
			switch *x.CallOutcome {
			case waE2E.CallLogMessage_CONNECTED:
				outcome = CallLogOutcomeConnected
			case waE2E.CallLogMessage_MISSED:
				outcome = CallLogOutcomeMissed
			case waE2E.CallLogMessage_FAILED:
				outcome = CallLogOutcomeFailed
			case waE2E.CallLogMessage_REJECTED:
				outcome = CallLogOutcomeRejected
			case waE2E.CallLogMessage_ACCEPTED_ELSEWHERE:
				outcome = CallLogOutcomeAccepted
			case waE2E.CallLogMessage_ONGOING:
				outcome = CallLogOutcomeOngoing
			case waE2E.CallLogMessage_SILENCED_BY_DND, waE2E.CallLogMessage_SILENCED_UNKNOWN_CALLER:
				outcome = CallLogOutcomeSilenced
			}
			message.CallLogOutcome = &outcome
		}
		if x.CallType != nil {
			var callType CallLogType
			switch *x.CallType {
			case waE2E.CallLogMessage_REGULAR:
				callType = CallLogTypeRegular
			case waE2E.CallLogMessage_SCHEDULED_CALL:
				callType = CallLogTypeScheduled
			case waE2E.CallLogMessage_VOICE_CHAT:
				callType = CallLogTypeVoiceChat
			}
			message.CallLogType = &callType
		}
		message.CallLogDurationSeconds = x.DurationSecs
		jids := []string{}
		for _, participant := range x.Participants {
			if participant == nil {
				continue
			}
			if participant.JID == nil {
				continue
			}
			jids = append(jids, *participant.JID)
		}
		message.CallLogParticipantJIDs = jids
	}
	if x := m.Message.GetEditedMessage(); x != nil {
		log.Warn().Any("x", x).Msg("NOT IMPLEMENTED: Message.GetEditedMessage()")
	}
	if x := m.Message.GetMessageContextInfo(); x != nil {
		//log.Warn().Any("x", x).Msg("NOT IMPLEMENTED: Message.GetMessageContextInfo()")
	}
	if x := m.Message.GetCall(); x != nil {
		log.Warn().Any("x", x).Msg("NOT IMPLEMENTED: Message.GetCall()")
	}
	if x := m.Message.GetChat(); x != nil {
		log.Warn().Any("x", x).Msg("NOT IMPLEMENTED: Message.GetChat()")
	}
	if x := m.Message.GetProtocolMessage(); x != nil {
		log.Warn().Any("x", x).Msg("NOT IMPLEMENTED: Message.GetProtocolMessage()")
	}
	if x := m.Message.GetContactsArrayMessage(); x != nil {
		log.Warn().Any("x", x).Msg("NOT IMPLEMENTED: Message.GetContactsArrayMessage()")
	}
	if x := m.Message.GetHighlyStructuredMessage(); x != nil {
		log.Warn().Any("x", x).Msg("NOT IMPLEMENTED: Message.GetHighlyStructuredMessage()")
	}
	if x := m.Message.GetFastRatchetKeySenderKeyDistributionMessage(); x != nil {
		log.Warn().Any("x", x).Msg("NOT IMPLEMENTED: Message.GetFastRatchetKeySenderKeyDistributionMessage()")
	}
	if x := m.Message.GetSendPaymentMessage(); x != nil {
		log.Warn().Any("x", x).Msg("NOT IMPLEMENTED: Message.GetSendPaymentMessage()")
	}
	if x := m.Message.GetRequestPaymentMessage(); x != nil {
		log.Warn().Any("x", x).Msg("NOT IMPLEMENTED: Message.GetRequestPaymentMessage()")
	}
	if x := m.Message.GetDeclinePaymentRequestMessage(); x != nil {
		log.Warn().Any("x", x).Msg("NOT IMPLEMENTED: Message.GetDeclinePaymentRequestMessage()")
	}
	if x := m.Message.GetCancelPaymentRequestMessage(); x != nil {
		log.Warn().Any("x", x).Msg("NOT IMPLEMENTED: Message.GetCancelPaymentRequestMessage()")
	}
	if x := m.Message.GetTemplateMessage(); x != nil {
		log.Warn().Any("x", x).Msg("NOT IMPLEMENTED: Message.GetTemplateMessage()")
	}
	if x := m.Message.GetTemplateButtonReplyMessage(); x != nil {
		log.Warn().Any("x", x).Msg("NOT IMPLEMENTED: Message.GetTemplateButtonReplyMessage()")
	}
	if x := m.Message.GetProductMessage(); x != nil {
		log.Warn().Any("x", x).Msg("NOT IMPLEMENTED: Message.GetProductMessage()")
	}
	if x := m.Message.GetDeviceSentMessage(); x != nil {
		log.Warn().Any("x", x).Msg("NOT IMPLEMENTED: Message.GetDeviceSentMessage()")
	}
	if x := m.Message.GetListMessage(); x != nil {
		log.Warn().Any("x", x).Msg("NOT IMPLEMENTED: Message.GetListMessage()")
	}
	if x := m.Message.GetViewOnceMessage(); x != nil {
		log.Warn().Any("x", x).Msg("NOT IMPLEMENTED: Message.GetViewOnceMessage()")
	}
	if x := m.Message.GetOrderMessage(); x != nil {
		log.Warn().Any("x", x).Msg("NOT IMPLEMENTED: Message.GetOrderMessage()")
	}
	if x := m.Message.GetListResponseMessage(); x != nil {
		log.Warn().Any("x", x).Msg("NOT IMPLEMENTED: Message.GetListResponseMessage()")
	}
	if x := m.Message.GetEphemeralMessage(); x != nil {
		log.Warn().Any("x", x).Msg("NOT IMPLEMENTED: Message.GetEphemeralMessage()")
	}
	if x := m.Message.GetInvoiceMessage(); x != nil {
		log.Warn().Any("x", x).Msg("NOT IMPLEMENTED: Message.GetInvoiceMessage()")
	}
	if x := m.Message.GetButtonsMessage(); x != nil {
		log.Warn().Any("x", x).Msg("NOT IMPLEMENTED: Message.GetButtonsMessage()")
	}
	if x := m.Message.GetButtonsResponseMessage(); x != nil {
		log.Warn().Any("x", x).Msg("NOT IMPLEMENTED: Message.GetButtonsResponseMessage()")
	}
	if x := m.Message.GetPaymentInviteMessage(); x != nil {
		log.Warn().Any("x", x).Msg("NOT IMPLEMENTED: Message.GetPaymentInviteMessage()")
	}
	if x := m.Message.GetInteractiveMessage(); x != nil {
		log.Warn().Any("x", x).Msg("NOT IMPLEMENTED: Message.GetInteractiveMessage()")
	}
	if x := m.Message.GetStickerSyncRmrMessage(); x != nil {
		log.Warn().Any("x", x).Msg("NOT IMPLEMENTED: Message.GetStickerSyncRmrMessage()")
	}
	if x := m.Message.GetInteractiveResponseMessage(); x != nil {
		log.Warn().Any("x", x).Msg("NOT IMPLEMENTED: Message.GetInteractiveResponseMessage()")
	}
	if x := m.Message.GetPollCreationMessage(); x != nil {
		log.Warn().Any("x", x).Msg("NOT IMPLEMENTED: Message.GetPollCreationMessage()")
	}
	if x := m.Message.GetPollUpdateMessage(); x != nil {
		log.Warn().Any("x", x).Msg("NOT IMPLEMENTED: Message.GetPollUpdateMessage()")
	}
	if x := m.Message.GetKeepInChatMessage(); x != nil {
		log.Warn().Any("x", x).Msg("NOT IMPLEMENTED: Message.GetKeepInChatMessage()")
	}
	if x := m.Message.GetDocumentWithCaptionMessage(); x != nil {
		log.Warn().Any("x", x).Msg("NOT IMPLEMENTED: Message.GetDocumentWithCaptionMessage()")
	}
	if x := m.Message.GetRequestPhoneNumberMessage(); x != nil {
		log.Warn().Any("x", x).Msg("NOT IMPLEMENTED: Message.GetRequestPhoneNumberMessage()")
	}
	if x := m.Message.GetViewOnceMessageV2(); x != nil {
		log.Warn().Any("x", x).Msg("NOT IMPLEMENTED: Message.GetViewOnceMessageV2()")
	}
	if x := m.Message.GetEncReactionMessage(); x != nil {
		log.Warn().Any("x", x).Msg("NOT IMPLEMENTED: Message.GetEncReactionMessage()")
	}
	if x := m.Message.GetViewOnceMessageV2Extension(); x != nil {
		log.Warn().Any("x", x).Msg("NOT IMPLEMENTED: Message.GetViewOnceMessageV2Extension()")
	}
	if x := m.Message.GetPollCreationMessageV2(); x != nil {
		log.Warn().Any("x", x).Msg("NOT IMPLEMENTED: Message.GetPollCreationMessageV2()")
	}
	if x := m.Message.GetScheduledCallCreationMessage(); x != nil {
		log.Warn().Any("x", x).Msg("NOT IMPLEMENTED: Message.GetScheduledCallCreationMessage()")
	}
	if x := m.Message.GetGroupMentionedMessage(); x != nil {
		log.Warn().Any("x", x).Msg("NOT IMPLEMENTED: Message.GetGroupMentionedMessage()")
	}
	if x := m.Message.GetPinInChatMessage(); x != nil {
		log.Warn().Any("x", x).Msg("NOT IMPLEMENTED: Message.GetPinInChatMessage()")
	}
	if x := m.Message.GetPollCreationMessageV3(); x != nil {
		log.Warn().Any("x", x).Msg("NOT IMPLEMENTED: Message.GetPollCreationMessageV3()")
	}
	if x := m.Message.GetScheduledCallEditMessage(); x != nil {
		log.Warn().Any("x", x).Msg("NOT IMPLEMENTED: Message.GetScheduledCallEditMessage()")
	}
	if x := m.Message.GetPtvMessage(); x != nil {
		log.Warn().Any("x", x).Msg("NOT IMPLEMENTED: Message.GetPtvMessage()")
	}
	if x := m.Message.GetBotInvokeMessage(); x != nil {
		log.Warn().Any("x", x).Msg("NOT IMPLEMENTED: Message.GetBotInvokeMessage()")
	}
	if x := m.Message.GetMessageHistoryBundle(); x != nil {
		log.Warn().Any("x", x).Msg("NOT IMPLEMENTED: Message.GetMessageHistoryBundle()")
	}
	if x := m.Message.GetEncCommentMessage(); x != nil {
		log.Warn().Any("x", x).Msg("NOT IMPLEMENTED: Message.GetEncCommentMessage()")
	}
	if x := m.Message.GetBcallMessage(); x != nil {
		log.Warn().Any("x", x).Msg("NOT IMPLEMENTED: Message.GetBcallMessage()")
	}
	if x := m.Message.GetLottieStickerMessage(); x != nil {
		log.Warn().Any("x", x).Msg("NOT IMPLEMENTED: Message.GetLottieStickerMessage()")
	}
	if x := m.Message.GetEventMessage(); x != nil {
		log.Warn().Any("x", x).Msg("NOT IMPLEMENTED: Message.GetEventMessage()")
	}
	if x := m.Message.GetEncEventResponseMessage(); x != nil {
		log.Warn().Any("x", x).Msg("NOT IMPLEMENTED: Message.GetEncEventResponseMessage()")
	}
	if x := m.Message.GetCommentMessage(); x != nil {
		log.Warn().Any("x", x).Msg("NOT IMPLEMENTED: Message.GetCommentMessage()")
	}
	if x := m.Message.GetNewsletterAdminInviteMessage(); x != nil {
		log.Warn().Any("x", x).Msg("NOT IMPLEMENTED: Message.GetNewsletterAdminInviteMessage()")
	}
	if x := m.Message.GetPlaceholderMessage(); x != nil {
		log.Warn().Any("x", x).Msg("NOT IMPLEMENTED: Message.GetPlaceholderMessage()")
	}
	if x := m.Message.GetSecretEncryptedMessage(); x != nil {
		log.Warn().Any("x", x).Msg("NOT IMPLEMENTED: Message.GetSecretEncryptedMessage()")
	}
	if x := m.Message.GetAlbumMessage(); x != nil {
		log.Warn().Any("x", x).Msg("NOT IMPLEMENTED: Message.GetAlbumMessage()")
	}
	if x := m.Message.GetEventCoverImage(); x != nil {
		log.Warn().Any("x", x).Msg("NOT IMPLEMENTED: Message.GetEventCoverImage()")
	}
	if x := m.Message.GetStickerPackMessage(); x != nil {
		log.Warn().Any("x", x).Msg("NOT IMPLEMENTED: Message.GetStickerPackMessage()")
	}
	if x := m.Message.GetStatusMentionMessage(); x != nil {
		log.Warn().Any("x", x).Msg("NOT IMPLEMENTED: Message.GetStatusMentionMessage()")
	}
	if x := m.Message.GetPollResultSnapshotMessage(); x != nil {
		log.Warn().Any("x", x).Msg("NOT IMPLEMENTED: Message.GetPollResultSnapshotMessage()")
	}
	if x := m.Message.GetPollCreationOptionImageMessage(); x != nil {
		log.Warn().Any("x", x).Msg("NOT IMPLEMENTED: Message.GetPollCreationOptionImageMessage()")
	}
	if x := m.Message.GetAssociatedChildMessage(); x != nil {
		log.Warn().Any("x", x).Msg("NOT IMPLEMENTED: Message.GetAssociatedChildMessage()")
	}
	if x := m.Message.GetGroupStatusMentionMessage(); x != nil {
		log.Warn().Any("x", x).Msg("NOT IMPLEMENTED: Message.GetGroupStatusMentionMessage()")
	}
	if x := m.Message.GetPollCreationMessageV4(); x != nil {
		log.Warn().Any("x", x).Msg("NOT IMPLEMENTED: Message.GetPollCreationMessageV4()")
	}
	if x := m.Message.GetPollCreationMessageV5(); x != nil {
		log.Warn().Any("x", x).Msg("NOT IMPLEMENTED: Message.GetPollCreationMessageV5()")
	}
	if x := m.Message.GetStatusAddYours(); x != nil {
		log.Warn().Any("x", x).Msg("NOT IMPLEMENTED: Message.GetStatusAddYours()")
	}
	if x := m.Message.GetGroupStatusMessage(); x != nil {
		log.Warn().Any("x", x).Msg("NOT IMPLEMENTED: Message.GetGroupStatusMessage()")
	}
	if x := m.Message.GetRichResponseMessage(); x != nil {
		log.Warn().Any("x", x).Msg("NOT IMPLEMENTED: Message.GetRichResponseMessage()")
	}
	if x := m.Message.GetStatusNotificationMessage(); x != nil {
		log.Warn().Any("x", x).Msg("NOT IMPLEMENTED: Message.GetStatusNotificationMessage()")
	}
	if x := m.Message.GetLimitSharingMessage(); x != nil {
		log.Warn().Any("x", x).Msg("NOT IMPLEMENTED: Message.GetLimitSharingMessage()")
	}
	if x := m.Message.GetBotTaskMessage(); x != nil {
		log.Warn().Any("x", x).Msg("NOT IMPLEMENTED: Message.GetBotTaskMessage()")
	}
	if x := m.Message.GetQuestionMessage(); x != nil {
		log.Warn().Any("x", x).Msg("NOT IMPLEMENTED: Message.GetQuestionMessage()")
	}

	conn.Callbacks.Message(message)
}
