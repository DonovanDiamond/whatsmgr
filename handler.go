package whatsmgr

import (
	"fmt"
	"time"

	"go.mau.fi/whatsmeow/proto/waHistorySync"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
)

type Callbacks struct {
	QRCode func(string)

	ConnStatus func(ConnStatus)
	Error      func(error)

	Contact func(Contact)
	Message func(Message)
	Call    func(Call)
	User    func(User)

	GetExistingProfilePhotoID func(jid string) (photoID string)
	PushNewProfilePhotoID     func(jid, photoID string)
}

type ConnStatus string

const (
	ConnStatusConnected    ConnStatus = "connected"
	ConnStatusDisconnected ConnStatus = "disconnected"
	ConnStatusQRCodeScan   ConnStatus = "qr code scan"
	ConnStatusLoggedOut    ConnStatus = "logged out"
	ConnStatusError        ConnStatus = "error"
)

type User struct {
	Timestamp time.Time

	Name *string // not always set
}

func (conn *Connection) handleEvent(rawEvt any) {
	log := conn.Log.With().Str("_module", "events").Logger()
	log.Trace().Type("event", rawEvt).Send()
	switch evt := rawEvt.(type) {
	case *events.Contact:
		contact := Contact{
			Timestamp: &evt.Timestamp,
			JID:       evt.JID.String(),
		}
		if evt.Action.FullName != nil {
			contact.ContactName = evt.Action.FullName
		}
		imageName, err := conn.pullProfilePhoto(evt.JID)
		if err != nil {
			log.Warn().Str("jid", contact.JID).Err(err).Msg("failed to pull profile photo")
		}
		if imageName != "" {
			contact.ProfilePhoto = &imageName
		}
		conn.Callbacks.Contact(contact)
	case *events.PushName:
		contact := Contact{
			JID:      evt.JID.String(),
			PushName: &evt.NewPushName,
		}
		imageName, err := conn.pullProfilePhoto(evt.JID)
		if err != nil {
			log.Warn().Str("jid", contact.JID).Err(err).Msg("failed to pull profile photo")
		}
		if imageName != "" {
			contact.ProfilePhoto = &imageName
		}
		conn.Callbacks.Contact(contact)
	case *events.BusinessName:
		contact := Contact{
			JID:      evt.JID.String(),
			PushName: &evt.NewBusinessName,
		}
		imageName, err := conn.pullProfilePhoto(evt.JID)
		if err != nil {
			log.Warn().Str("jid", contact.JID).Err(err).Msg("failed to pull profile photo")
		}
		if imageName != "" {
			contact.ProfilePhoto = &imageName
		}
		conn.Callbacks.Contact(contact)
	case *events.Pin:
		contact := Contact{
			Timestamp: &evt.Timestamp,
			JID:       evt.JID.String(),
		}
		if evt.Action != nil {
			contact.Pinned = evt.Action.Pinned
		}
		conn.Callbacks.Contact(contact)
	case *events.Star:
		msg := Message{
			// Timestamp: evt.Timestamp, // we don't want to change the messages date because someone starred it

			MessageID: evt.MessageID,
			ChatJID:   evt.ChatJID.String(),
			// SenderJID: evt.SenderJID.String(), // we don't want to change the sender if I starred it

			// IsFromMe: evt.IsFromMe, // this may not be correct, maybe its because the event is from me, and the message itself is not...
		}
		if evt.Action != nil {
			msg.Starred = evt.Action.Starred
		}
		conn.Callbacks.Message(msg)
	case *events.DeleteForMe:
		msg := Message{
			// Timestamp: evt.Timestamp, // we don't want to change the messages date because someone deleted it

			MessageID: evt.MessageID,
			ChatJID:   evt.ChatJID.String(),
			// SenderJID: evt.SenderJID.String(), // we don't want to change the sender if I deleted it

			// IsFromMe: evt.IsFromMe, // this may not be correct, maybe its because the event is from me, and the message itself is not...
		}
		if evt.Action != nil {
			t := true
			msg.Deleted = &t
		}
		conn.Callbacks.Message(msg)
	case *events.Mute:
		contact := Contact{
			Timestamp: &evt.Timestamp,
			JID:       evt.JID.String(),
		}
		if evt.Action != nil {
			contact.Muted = evt.Action.Muted
			contact.MuteEndTimestamp = evt.Action.MuteEndTimestamp
		}
		conn.Callbacks.Contact(contact)
	case *events.Archive:
		contact := Contact{
			Timestamp: &evt.Timestamp,
			JID:       evt.JID.String(),
		}
		if evt.Action != nil {
			contact.Archived = evt.Action.Archived
			// We are not including message range here
		}
		conn.Callbacks.Contact(contact)
	case *events.MarkChatAsRead:
		log.Warn().Any("evt", evt).Type("type", evt).Msg("NOT IMPLEMENTED")
	case *events.ClearChat:
		log.Warn().Any("evt", evt).Type("type", evt).Msg("NOT IMPLEMENTED")
	case *events.DeleteChat:
		log.Warn().Any("evt", evt).Type("type", evt).Msg("NOT IMPLEMENTED")
	case *events.PushNameSetting:
		user := User{
			Timestamp: evt.Timestamp,
		}
		if evt.Action != nil {
			user.Name = evt.Action.Name
		}
		conn.Callbacks.User(user)
	case *events.UnarchiveChatsSetting:
		log.Warn().Any("evt", evt).Type("type", evt).Msg("NOT IMPLEMENTED")
	case *events.UserStatusMute:
		contact := Contact{
			Timestamp: &evt.Timestamp,
			JID:       evt.JID.String(),
		}
		if evt.Action != nil {
			contact.StatusMuted = evt.Action.Muted
		}
		conn.Callbacks.Contact(contact)
	case *events.LabelEdit:
		log.Warn().Any("evt", evt).Type("type", evt).Msg("NOT IMPLEMENTED")
	case *events.LabelAssociationChat:
		log.Warn().Any("evt", evt).Type("type", evt).Msg("NOT IMPLEMENTED")
	case *events.LabelAssociationMessage:
		log.Warn().Any("evt", evt).Type("type", evt).Msg("NOT IMPLEMENTED")
	case *events.AppStateSyncComplete:
		log.Info().Str("name", string(evt.Name)).Msg("Completed App State Sync")
	case *events.CallOffer:
		call := Call{
			Timestamp:   evt.Timestamp,
			CallID:      evt.CallID,
			From:        evt.From.String(),
			CallCreator: evt.CallCreator.String(),
			Status:      CallStatusOffer,
		}
		conn.Callbacks.Call(call)
	case *events.CallAccept:
		call := Call{
			Timestamp:   evt.Timestamp,
			CallID:      evt.CallID,
			From:        evt.From.String(),
			CallCreator: evt.CallCreator.String(),
			Status:      CallStatusAccept,
		}
		conn.Callbacks.Call(call)
	case *events.CallPreAccept:
		call := Call{
			Timestamp:   evt.Timestamp,
			CallID:      evt.CallID,
			From:        evt.From.String(),
			CallCreator: evt.CallCreator.String(),
			Status:      CallStatusPreAccept,
		}
		conn.Callbacks.Call(call)
	case *events.CallTransport:
		call := Call{
			Timestamp:   evt.Timestamp,
			CallID:      evt.CallID,
			From:        evt.From.String(),
			CallCreator: evt.CallCreator.String(),
			Status:      CallStatusTransport,
		}
		conn.Callbacks.Call(call)
	case *events.CallOfferNotice:
		var media CallMedia
		if evt.Media == string(CallMediaAudio) {
			media = CallMediaAudio
		} else if evt.Media == string(CallMediaVideo) {
			media = CallMediaVideo
		}
		var typ CallType
		if evt.Type == string(CallTypeGroup) {
			typ = CallTypeGroup
		}
		call := Call{
			Timestamp:   evt.Timestamp,
			CallID:      evt.CallID,
			From:        evt.From.String(),
			CallCreator: evt.CallCreator.String(),
			Status:      CallStatusOffer,
			Media:       &media,
			Type:        &typ,
		}
		conn.Callbacks.Call(call)
	case *events.CallRelayLatency:
		log.Warn().Any("evt", evt).Type("type", evt).Msg("NOT IMPLEMENTED")
	case *events.CallTerminate:
		call := Call{
			Timestamp:        evt.Timestamp,
			CallID:           evt.CallID,
			From:             evt.From.String(),
			CallCreator:      evt.CallCreator.String(),
			Status:           CallStatusTerminate,
			TerminateReaason: &evt.Reason,
		}
		conn.Callbacks.Call(call)
	case *events.CallReject:
		call := Call{
			Timestamp:   evt.Timestamp,
			CallID:      evt.CallID,
			From:        evt.From.String(),
			CallCreator: evt.CallCreator.String(),
			Status:      CallStatusReject,
		}
		conn.Callbacks.Call(call)
	case *events.UnknownCallEvent:
		log.Warn().Any("evt", evt).Type("type", evt).Msg("NOT IMPLEMENTED")
	case *events.QR:
		log.Warn().Any("evt", evt).Type("type", evt).Msg("NOT IMPLEMENTED")
	case *events.PairSuccess:
		log.Warn().Any("evt", evt).Type("type", evt).Msg("NOT IMPLEMENTED")
	case *events.PairError:
		log.Warn().Any("evt", evt).Type("type", evt).Msg("NOT IMPLEMENTED")
	case *events.QRScannedWithoutMultidevice:
		log.Warn().Any("evt", evt).Type("type", evt).Msg("NOT IMPLEMENTED")
	case *events.Connected:
		log.Info().Msg("Client is Connected, sending PresenceAvailable")
		conn.client.SendPresence(types.PresenceAvailable)
		conn.Callbacks.ConnStatus(ConnStatusConnected)
	case *events.KeepAliveTimeout:
		log.Warn().Any("evt", evt).Type("type", evt).Msg("NOT IMPLEMENTED")
	case *events.KeepAliveRestored:
		log.Warn().Any("evt", evt).Type("type", evt).Msg("NOT IMPLEMENTED")
	case *events.PermanentDisconnect:
		log.Info().Msg("Client is PermanentDisconnect")
		conn.Callbacks.ConnStatus(ConnStatusDisconnected)
	case *events.LoggedOut:
		log.Warn().Msg("Client is Logged Out (Disconnected)")
		conn.Callbacks.ConnStatus(ConnStatusLoggedOut)
	case *events.StreamReplaced:
		log.Error().Msg("Stream Replaced")
		conn.client.Disconnect()
		conn.Callbacks.ConnStatus(ConnStatusError)
		reason := evt.PermanentDisconnectDescription()
		if reason != "" {
			conn.Callbacks.Error(fmt.Errorf("stream has been replaced: %s", reason))
		}
	case *events.ManualLoginReconnect:
		// TODO: maybe implement?
	case *events.TempBanReason:
		reason := evt.String()
		log.Error().Str("reason", reason).Msg("WARNING: Client is Temporary Banned")
		conn.Callbacks.ConnStatus(ConnStatusError)
		if reason != "" {
			conn.Callbacks.Error(fmt.Errorf("client has been temporarily banned: %s", reason))
		}
	case *events.TemporaryBan:
		reason := evt.String()
		log.Error().Str("reason", reason).Dur("expires", evt.Expire).Msg("WARNING: Client is Temporary Banned")
		conn.Callbacks.ConnStatus(ConnStatusError)
		if reason != "" {
			conn.Callbacks.Error(fmt.Errorf("client has been temporarily banned: %s (expires in %s)", reason, evt.Expire.String()))
		}
	case *events.ConnectFailureReason:
		reason := evt.String()
		log.Warn().Str("reason", reason).Bool("isLoggedOut", evt.IsLoggedOut()).Msg("Connect failure")
		if evt.IsLoggedOut() {
			conn.Callbacks.ConnStatus(ConnStatusLoggedOut)
		} else {
			conn.Callbacks.ConnStatus(ConnStatusError)
			conn.Callbacks.Error(fmt.Errorf("failed to connect: %s", reason))
		}
	case *events.ConnectFailure:
		log.Warn().Str("message", evt.Message).Any("reason", evt.Reason).Msg("Connect failure")
		conn.Callbacks.ConnStatus(ConnStatusError)
		conn.Callbacks.Error(fmt.Errorf("failed to connect: %s (%v)", evt.Message, evt.Reason))
	case *events.ClientOutdated:
		log.Error().Msg("Client is outdated")
		conn.Callbacks.ConnStatus(ConnStatusError)
		conn.Callbacks.Error(fmt.Errorf("update is required"))
	case *events.CATRefreshError:
		log.Warn().Err(evt.Error).Msg("CATRefreshError")
		conn.Callbacks.ConnStatus(ConnStatusError)
		conn.Callbacks.Error(evt.Error)
	case *events.StreamError:
		log.Error().Msg("Stream Error")
		log.Warn().Str("code", evt.Code).Msg("Stream error")
		conn.Callbacks.ConnStatus(ConnStatusError)
		conn.Callbacks.Error(fmt.Errorf("stream error: %s", evt.Code))
	case *events.Disconnected:
		log.Info().Msg("Client is Disconnected")
		conn.Callbacks.ConnStatus(ConnStatusDisconnected)
	case *events.HistorySync:
		conn.Callbacks.ConnStatus(ConnStatusConnected)
		for _, conv := range evt.Data.Conversations {
			if conv.ID == nil {
				continue
			}
			contact := Contact{
				JID:         conv.GetID(),
				PushName:    conv.Name,
				DisplayName: conv.DisplayName,
				Username:    conv.Username,
				Archived:    conv.Archived,

				Group: Group{
					GroupName:               conv.Name,
					GroupTopic:              conv.Description,
					GroupInfoLockedToAdmins: conv.Locked,
					GroupIsParent:           conv.IsParentGroup,
					GroupLinkedParentJID:    conv.ParentGroupID,
					GroupIsDefaultSubGroup:  conv.IsDefaultSubgroup,
				},
			}
			jid, err := types.ParseJID(contact.JID)
			if err != nil {
				log.Warn().Str("jid", contact.JID).Err(err).Msg("failed to parse jid during history sync")
			} else {
				imageName, err := conn.pullProfilePhoto(jid)
				if err != nil {
					log.Warn().Str("jid", contact.JID).Err(err).Msg("failed to pull profile photo")
				}
				if imageName != "" {
					contact.ProfilePhoto = &imageName
				}
			}
			if conv.Pinned != nil {
				pinned := conv.GetPinned() != 0
				contact.Pinned = &pinned
			}
			if conv.MuteEndTime != nil {
				muteEndTime := int64(conv.GetMuteEndTime())
				muted := muteEndTime > 0
				contact.Muted = &muted
				contact.MuteEndTimestamp = &muteEndTime
			}
			groupParticipants := []GroupParticipant{}
			for _, par := range conv.Participant {
				if par == nil {
					continue
				}
				groupParticipant := GroupParticipant{
					UserJID: par.GetUserJID(),
				}
				switch par.GetRank() {
				case waHistorySync.GroupParticipant_REGULAR:
					groupParticipant.Rank = GroupParticipantRankRegular
				case waHistorySync.GroupParticipant_ADMIN:
					groupParticipant.Rank = GroupParticipantRankAdmin
				case waHistorySync.GroupParticipant_SUPERADMIN:
					groupParticipant.Rank = GroupParticipantRankSuperAdmin
				}
				groupParticipants = append(groupParticipants, groupParticipant)
			}
			if len(groupParticipants) > 0 {
				contact.GroupReplaceParticipants = groupParticipants
			}
			conn.Callbacks.Contact(contact)

			chatJID, _ := types.ParseJID(conv.GetID())
			for _, syncMessage := range conv.GetMessages() {
				m, err := conn.client.ParseWebMessage(chatJID, syncMessage.GetMessage())
				if err != nil {
					log.Error().Err(err).Msg("error parsing web message")
					continue
				}
				if m == nil {
					continue
				}
				conn.handleMessage(*m)
			}
		}
	case *events.DecryptFailMode:
		log.Warn().Any("evt", evt).Type("type", evt).Msg("NOT IMPLEMENTED")
	case *events.UndecryptableMessage:
		log.Warn().Any("evt", evt).Type("type", evt).Msg("NOT IMPLEMENTED")
	case *events.NewsletterMessageMeta:
		log.Warn().Any("evt", evt).Type("type", evt).Msg("NOT IMPLEMENTED")
	case *events.Message:
		if evt != nil {
			conn.handleMessage(*evt)
		}
	case *events.FBMessage:
		log.Warn().Any("evt", evt).Type("type", evt).Msg("NOT IMPLEMENTED")
	case *events.Receipt:
		var status MessageStatus
		switch evt.Type {
		case types.ReceiptTypeSender:
			status = MessageStatusSent
		case types.ReceiptTypeDelivered:
			status = MessageStatusDelivered
		case types.ReceiptTypeRead, types.ReceiptTypePlayed, types.ReceiptTypeReadSelf:
			status = MessageStatusRead
		case types.ReceiptTypeServerError:
			status = MessageStatusServerError
		case types.ReceiptTypeInactive:
			// this means that the chat is "inactive", so we should just update the contact and leave it at that
			available := false
			conn.Callbacks.Contact(Contact{
				JID: evt.Chat.String(),

				Available: &available,
			})
			return
		default:
			log.Error().Any("evt", evt).Type("type", evt.Type).Msg("Unknown receipt type received")
			return
		}
		for _, id := range evt.MessageIDs {
			message := Message{
				// Timestamp: evt.Timestamp, // we don't want to change the messages date because someone read it...

				MessageID: id,
				ChatJID:   evt.Chat.String(),
				// SenderJID: evt.Sender.String(), // we don't want to change the sender if a message was read by myself...

				// IsFromMe: evt.IsFromMe, // this may not be correct, maybe its because the event is from me, and the message itself is not...
				Status: &status,
			}
			conn.Callbacks.Message(message)
		}
	case *events.ChatPresence:
		contact := Contact{
			JID: evt.Chat.String(),
		}
		typing := ""
		recording := ""
		available := true
		lastSeen := time.Now()
		if evt.State == types.ChatPresenceComposing {
			if evt.Media == types.ChatPresenceMediaAudio {
				recording = evt.Sender.String()
			} else {
				typing = evt.Sender.String()
			}
			contact.LastSeen = &lastSeen
			contact.Available = &available
		}
		contact.Typing = &typing
		contact.Recording = &recording
		conn.Callbacks.Contact(contact)
	case *events.Presence:
		available := !evt.Unavailable
		contact := Contact{
			JID:       evt.From.String(),
			Available: &available,
			LastSeen:  &evt.LastSeen,
		}
		imageName, err := conn.pullProfilePhoto(evt.From)
		if err != nil {
			log.Warn().Str("jid", contact.JID).Err(err).Msg("failed to pull profile photo")
		}
		if imageName != "" {
			contact.ProfilePhoto = &imageName
		}
		conn.Callbacks.Contact(contact)
	case *events.JoinedGroup:
		parentJID := evt.GroupInfo.LinkedParentJID.String()
		onlyAdminsCanAddMembers := evt.GroupInfo.MemberAddMode == types.GroupMemberAddModeAdmin
		isGroup := true
		contact := Contact{
			JID:     evt.JID.String(),
			IsGroup: &isGroup,
			Group: Group{
				GroupName:                    &evt.GroupInfo.Name,
				GroupTopic:                   &evt.GroupInfo.Topic,
				GroupInfoLockedToAdmins:      &evt.GroupInfo.IsLocked,
				GroupOnlyAdminCanMessage:     &evt.GroupInfo.IsAnnounce,
				GroupIsParent:                &evt.GroupInfo.IsParent,
				GroupLinkedParentJID:         &parentJID,
				GroupIsDefaultSubGroup:       &evt.GroupInfo.IsDefaultSubGroup,
				GroupJoinApprovalRequired:    &evt.GroupInfo.IsJoinApprovalRequired,
				GroupCreated:                 &evt.GroupInfo.GroupCreated,
				GroupOnlyAdminsCanAddMembers: &onlyAdminsCanAddMembers,
			},
		}
		imageName, err := conn.pullProfilePhoto(evt.JID)
		if err != nil {
			log.Warn().Str("jid", contact.JID).Err(err).Msg("failed to pull profile photo")
		}
		if imageName != "" {
			contact.ProfilePhoto = &imageName
		}
		groupParticipants := []GroupParticipant{}
		for _, par := range evt.GroupInfo.Participants {
			groupParticipant := GroupParticipant{
				UserJID: par.JID.String(),
			}
			groupParticipant.Rank = GroupParticipantRankRegular
			if par.IsAdmin {
				groupParticipant.Rank = GroupParticipantRankAdmin
			} else if par.IsSuperAdmin {
				groupParticipant.Rank = GroupParticipantRankSuperAdmin
			}
			groupParticipants = append(groupParticipants, groupParticipant)
		}
		if len(groupParticipants) > 0 {
			contact.GroupReplaceParticipants = groupParticipants
		}
		conn.Callbacks.Contact(contact)
	case *events.GroupInfo:
		isGroup := true
		contact := Contact{
			JID:     evt.JID.String(),
			IsGroup: &isGroup,
			Group: Group{
				GroupInviteLink: evt.NewInviteLink,
			},
		}
		if evt.Name != nil {
			contact.GroupName = &evt.Name.Name
		}
		if evt.Topic != nil {
			contact.GroupTopic = &evt.Topic.Topic
		}
		if evt.Locked != nil {
			contact.GroupInfoLockedToAdmins = &evt.Locked.IsLocked
		}
		if evt.Announce != nil {
			contact.GroupOnlyAdminCanMessage = &evt.Announce.IsAnnounce
		}
		if evt.MembershipApprovalMode != nil {
			contact.GroupJoinApprovalRequired = &evt.MembershipApprovalMode.IsJoinApprovalRequired
		}
		imageName, err := conn.pullProfilePhoto(evt.JID)
		if err != nil {
			log.Warn().Str("jid", contact.JID).Err(err).Msg("failed to pull profile photo")
		}
		if imageName != "" {
			contact.ProfilePhoto = &imageName
		}

		contact.GroupAddParticipants = []GroupParticipant{}
		contact.GroupRemovedParticipants = []GroupParticipant{}
		for _, jid := range evt.Join {
			contact.GroupAddParticipants = append(contact.GroupAddParticipants, GroupParticipant{
				UserJID: jid.String(),
			})
		}
		for _, jid := range evt.Leave {
			contact.GroupRemovedParticipants = append(contact.GroupRemovedParticipants, GroupParticipant{
				UserJID: jid.String(),
			})
		}
		for _, jid := range evt.Promote {
			contact.GroupAddParticipants = append(contact.GroupAddParticipants, GroupParticipant{
				UserJID: jid.String(),
				Rank:    GroupParticipantRankAdmin,
			})
		}
		for _, jid := range evt.Demote {
			contact.GroupAddParticipants = append(contact.GroupAddParticipants, GroupParticipant{
				UserJID: jid.String(),
				Rank:    GroupParticipantRankRegular,
			})
		}
		conn.Callbacks.Contact(contact)
	case *events.Picture:
		contact := Contact{
			JID: evt.JID.String(),
		}
		imageName, err := conn.pullProfilePhoto(evt.JID)
		if err != nil {
			log.Warn().Str("jid", contact.JID).Err(err).Msg("failed to pull profile photo")
		}
		if imageName != "" {
			contact.ProfilePhoto = &imageName
		}
		conn.Callbacks.Contact(contact)
	case *events.UserAbout:
		log.Warn().Any("evt", evt).Type("type", evt).Msg("NOT IMPLEMENTED")
	case *events.IdentityChange:
		log.Warn().Any("evt", evt).Type("type", evt).Msg("NOT IMPLEMENTED")
	case *events.PrivacySettings:
		log.Warn().Any("evt", evt).Type("type", evt).Msg("NOT IMPLEMENTED")
	case *events.OfflineSyncPreview:
		log.Info().Any("evt", evt).
			Int("Total", evt.Total).
			Int("AppDataChanges", evt.AppDataChanges).
			Int("Messages", evt.Messages).
			Int("Notifications", evt.Notifications).
			Int("Receipts", evt.Receipts).
			Msg("Received sync preview")
	case *events.OfflineSyncCompleted:
		log.Info().Int("count", evt.Count).Msg("Completed Offline Sync")
	case *events.MediaRetryError:
		log.Warn().Any("evt", evt).Type("type", evt).Msg("NOT IMPLEMENTED")
	case *events.MediaRetry:
		log.Warn().Any("evt", evt).Type("type", evt).Msg("NOT IMPLEMENTED")
	case *events.BlocklistAction:
		log.Warn().Any("evt", evt).Type("type", evt).Msg("NOT IMPLEMENTED")
	case *events.Blocklist:
		log.Warn().Any("evt", evt).Type("type", evt).Msg("NOT IMPLEMENTED")
	case *events.BlocklistChangeAction:
		log.Warn().Any("evt", evt).Type("type", evt).Msg("NOT IMPLEMENTED")
	case *events.BlocklistChange:
		log.Warn().Any("evt", evt).Type("type", evt).Msg("NOT IMPLEMENTED")
	case *events.NewsletterJoin:
		log.Warn().Any("evt", evt).Type("type", evt).Msg("NOT IMPLEMENTED")
	case *events.NewsletterLeave:
		log.Warn().Any("evt", evt).Type("type", evt).Msg("NOT IMPLEMENTED")
	case *events.NewsletterMuteChange:
		log.Warn().Any("evt", evt).Type("type", evt).Msg("NOT IMPLEMENTED")
	case *events.NewsletterLiveUpdate:
		log.Warn().Any("evt", evt).Type("type", evt).Msg("NOT IMPLEMENTED")
	default:
		log.Error().Any("evt", evt).Type("type", evt).Msg("UNKNOWN EVENT TYPE")
	}
}
