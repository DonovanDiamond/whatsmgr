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
	Timestamp *time.Time // not always set

	JID          string
	PushName     *string // not always set
	ContactName  *string // not always set
	DisplayName  *string // not always set
	Username     *string // not always set
	ProfilePhoto *string // not always set

	Pinned           *bool  // not always set
	Muted            *bool  // not always set
	MuteEndTimestamp *int64 // not always set
	StatusMuted      *bool  // not always set
	Archived         *bool  // not always set

	Available *bool      // not always set
	LastSeen  *time.Time // not always set
	Typing    *string    // the JID who is typing, not always set
	Recording *string    // the JID who is recording, not always set

	IsGroup *bool // not always set

	Group
}

type Group struct {
	GroupName                    *string            // not always set
	GroupTopic                   *string            // not always set
	GroupInfoLockedToAdmins      *bool              // not always set
	GroupOnlyAdminCanMessage     *bool              // not always set
	GroupOnlyAdminsCanAddMembers *bool              // not always set
	GroupIsParent                *bool              // not always set
	GroupLinkedParentJID         *string            // not always set
	GroupIsDefaultSubGroup       *bool              // not always set
	GroupJoinApprovalRequired    *bool              // not always set
	GroupCreated                 *time.Time         // not always set
	GroupReplaceParticipants     []GroupParticipant // replace all existing participants with these, not always set
	GroupAddParticipants         []GroupParticipant // add to existing participants, not always set
	GroupRemovedParticipants     []GroupParticipant // remove from existing participants, not always set
	GroupInviteLink              *string            // not always set
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
	existingID := conn.ProfilePhotoCache[jid.String()]

	info, err := conn.client.GetProfilePictureInfo(jid, &whatsmeow.GetProfilePictureParams{
		ExistingID: existingID,
	})
	if err != nil || info == nil {
		return "", nil
	}
	conn.ProfilePhotoCache[jid.String()] = info.ID

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
	imageName = conn.hashFile(raw)
	path := fmt.Sprintf("%s/%s", conn.MediaPath, imageName)
	if err = conn.writeFileIfNotExists(path, raw); err != nil {
		return "", fmt.Errorf("failed to save profile image: %w", err)
	}
	return
}
