package whatsmgr

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

func (conn *Connection) SendMessage(message Message, sendOnCallback bool) (Message, error) {
	var out waE2E.Message
	if len(message.Attachments) > 0 {
		attachment := message.Attachments[0]

		parts := strings.Split(attachment, ".")
		if len(parts) < 2 {
			return message, fmt.Errorf("invalid attachment with no extension: %s", attachment)
		}
		ext := strings.ToLower(parts[len(parts)-1])

		path := fmt.Sprintf("%s/%s", conn.MediaPath, attachment)
		raw, err := os.ReadFile(path)
		if err != nil {
			return message, fmt.Errorf("failed to read file to send: %w", err)
		}

		switch ext {
		case "jpg", "jpeg", "png":
			resp, err := conn.client.Upload(context.Background(), raw, whatsmeow.MediaImage)
			if err != nil {
				return message, fmt.Errorf("failed to upload image to send: %w", err)
			}
			out.ImageMessage = &waE2E.ImageMessage{
				Caption:       proto.String(*message.ContentBody),
				Mimetype:      proto.String("image/" + ext),
				URL:           &resp.URL,
				DirectPath:    &resp.DirectPath,
				MediaKey:      resp.MediaKey,
				FileEncSHA256: resp.FileEncSHA256,
				FileSHA256:    resp.FileSHA256,
				FileLength:    &resp.FileLength,
			}
		case "mp4":
			resp, err := conn.client.Upload(context.Background(), raw, whatsmeow.MediaVideo)
			if err != nil {
				return message, fmt.Errorf("failed to upload video to send: %w", err)
			}
			out.VideoMessage = &waE2E.VideoMessage{
				Caption:       proto.String(*message.ContentBody),
				Mimetype:      proto.String("video/" + ext),
				URL:           &resp.URL,
				DirectPath:    &resp.DirectPath,
				MediaKey:      resp.MediaKey,
				FileEncSHA256: resp.FileEncSHA256,
				FileSHA256:    resp.FileSHA256,
				FileLength:    &resp.FileLength,
			}
		case "aac", "amr", "mp3", "m4a", "ogg":
			resp, err := conn.client.Upload(context.Background(), raw, whatsmeow.MediaAudio)
			if err != nil {
				return message, fmt.Errorf("failed to upload audio to send: %w", err)
			}
			mimeType := "audio/" + ext
			switch ext {
			case "mp3":
				mimeType = "audio/mpeg"
			case "m4a":
				mimeType = "audio/mp4"
			}
			out.AudioMessage = &waE2E.AudioMessage{
				Mimetype:      proto.String(mimeType),
				URL:           &resp.URL,
				DirectPath:    &resp.DirectPath,
				MediaKey:      resp.MediaKey,
				FileEncSHA256: resp.FileEncSHA256,
				FileSHA256:    resp.FileSHA256,
				FileLength:    &resp.FileLength,
			}
		case "txt", "xls", "xlsx", "doc", "docx", "ppt", "pptx", "pdf":
			resp, err := conn.client.Upload(context.Background(), raw, whatsmeow.MediaDocument)
			if err != nil {
				return message, fmt.Errorf("failed to upload document to send: %w", err)
			}
			mimeType := ""
			switch ext {
			case "txt":
				mimeType = "text/plain"
			case "xls":
				mimeType = "application/vnd.ms-excel"
			case "xlsx":
				mimeType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
			case "doc":
				mimeType = "application/msword"
			case "docx":
				mimeType = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
			case "ppt":
				mimeType = "application/vnd.ms-powerpoint"
			case "pptx":
				mimeType = "application/vnd.openxmlformats-officedocument.presentationml.presentation"
			case "pdf":
				mimeType = "application/pdf"
			default:
				return message, fmt.Errorf("unknown attachment type (could not determine mime type): %s", ext)
			}
			out.DocumentMessage = &waE2E.DocumentMessage{
				Caption:       proto.String(*message.ContentBody),
				Mimetype:      proto.String(mimeType),
				URL:           &resp.URL,
				DirectPath:    &resp.DirectPath,
				MediaKey:      resp.MediaKey,
				FileEncSHA256: resp.FileEncSHA256,
				FileSHA256:    resp.FileSHA256,
				FileLength:    &resp.FileLength,
			}
		default:
			return message, fmt.Errorf("unknown attachment type: %s", ext)
		}
	} else {
		out = waE2E.Message{
			Conversation: message.ContentBody,
		}
	}
	jid, err := types.ParseJID(message.ChatJID)
	if err != nil {
		return message, fmt.Errorf("failed to parse ChatJID: %w", err)
	}
	resp, err := conn.client.SendMessage(context.Background(), jid, &out)
	if err != nil {
		return message, fmt.Errorf("failed to send message: %w", err)
	}
	message.MessageID = resp.ID
	senderJID := resp.Sender.String()
	message.SenderJID = &senderJID
	message.Timestamp = &resp.Timestamp
	if sendOnCallback {
		conn.Callbacks.Message(message)
	}
	return message, nil
}

func (conn *Connection) SendEdit(message Message) error {
	chat, err := types.ParseJID(message.ChatJID)
	if err != nil {
		return fmt.Errorf("failed to parse ChatJID: %w", err)
	}
	if message.MessageID == "" {
		return errors.New("missing message.MessageID")
	}
	out := waE2E.Message{
		Conversation: message.ContentBody,
	}
	_, err = conn.client.SendMessage(context.Background(), chat, conn.client.BuildEdit(chat, message.MessageID, &out))
	if err != nil {
		return fmt.Errorf("failed to send edit message: %w", err)
	}
	return nil
}

func (conn *Connection) SendRead(messageIDs []string, when time.Time, chatJID string, senderJID string, receiptTypeExtra ...types.ReceiptType) error {
	chat, err := types.ParseJID(chatJID)
	if err != nil {
		return fmt.Errorf("failed to parse chatJID ('%s'): %w", chatJID, err)
	}
	sender, err := types.ParseJID(senderJID)
	if err != nil {
		return fmt.Errorf("failed to parse senderJID ('%s'): %w", senderJID, err)
	}
	return conn.client.MarkRead(messageIDs, when, chat, sender, receiptTypeExtra...)
}

func (conn *Connection) SendPlayed(messageIDs []string, when time.Time, chatJID string, senderJID string) error {
	return conn.SendRead(messageIDs, when, chatJID, senderJID, types.ReceiptTypePlayed)
}
