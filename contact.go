package whatsmgr

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
)

type Contact struct {
	Timestamp *time.Time `json:",omitempty"` // not always set

	JID          string  `json:",omitempty"`
	PushName     *string `json:",omitempty"` // not always set
	ContactName  *string `json:",omitempty"` // not always set
	DisplayName  *string `json:",omitempty"` // not always set
	Username     *string `json:",omitempty"` // not always set
	ProfilePhoto *string `json:",omitempty"` // not always set

	Pinned           *bool  `json:",omitempty"` // not always set
	Muted            *bool  `json:",omitempty"` // not always set
	MuteEndTimestamp *int64 `json:",omitempty"` // not always set
	StatusMuted      *bool  `json:",omitempty"` // not always set
	Archived         *bool  `json:",omitempty"` // not always set

	Available *bool      `json:",omitempty"` // not always set
	LastSeen  *time.Time `json:",omitempty"` // not always set
	Typing    *string    `json:",omitempty"` // the JID who is typing, not always set
	Recording *string    `json:",omitempty"` // the JID who is recording, not always set

	IsGroup *bool `json:",omitempty"` // not always set

	Group
}

type Group struct {
	GroupName                    *string            `json:",omitempty"` // not always set
	GroupTopic                   *string            `json:",omitempty"` // not always set
	GroupInfoLockedToAdmins      *bool              `json:",omitempty"` // not always set
	GroupOnlyAdminCanMessage     *bool              `json:",omitempty"` // not always set
	GroupOnlyAdminsCanAddMembers *bool              `json:",omitempty"` // not always set
	GroupIsParent                *bool              `json:",omitempty"` // not always set
	GroupLinkedParentJID         *string            `json:",omitempty"` // not always set
	GroupIsDefaultSubGroup       *bool              `json:",omitempty"` // not always set
	GroupJoinApprovalRequired    *bool              `json:",omitempty"` // not always set
	GroupCreated                 *time.Time         `json:",omitempty"` // not always set
	GroupReplaceParticipants     []GroupParticipant `json:",omitempty"` // replace all existing participants with these, not always set
	GroupAddParticipants         []GroupParticipant `json:",omitempty"` // add to existing participants, not always set
	GroupRemovedParticipants     []GroupParticipant `json:",omitempty"` // remove from existing participants, not always set
	GroupInviteLink              *string            `json:",omitempty"` // not always set
}

type GroupParticipant struct {
	UserJID string
	Rank    GroupParticipantRank
}

type GroupParticipantRank string

const (
	GroupParticipantRankRegular    GroupParticipantRank = "regular"
	GroupParticipantRankAdmin      GroupParticipantRank = "admin"
	GroupParticipantRankSuperAdmin GroupParticipantRank = "superadmin"
)

func (conn *Connection) pullProfilePhoto(jid types.JID) (imageName string, err error) {
	existingID := conn.Callbacks.GetExistingProfilePhotoID(jid.String())

	info, err := conn.client.GetProfilePictureInfo(jid, &whatsmeow.GetProfilePictureParams{
		ExistingID: existingID,
	})
	if err != nil || info == nil {
		return "", nil
	}
	conn.Callbacks.PushNewProfilePhotoID(jid.String(), info.ID)

	resp, err := http.Get(info.URL)
	if err != nil {
		return "", fmt.Errorf("failed to download profile image: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error returned when downloading profile image: %w", err)
	}
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read profile image: %w", err)
	}
	imageName = conn.hashFile(raw) + ".jpeg"
	path := fmt.Sprintf("%s/%s", conn.MediaPath, imageName)
	if err = conn.writeFileIfNotExists(path, raw); err != nil {
		return "", fmt.Errorf("failed to save profile image: %w", err)
	}
	return
}
