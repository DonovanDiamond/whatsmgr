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
