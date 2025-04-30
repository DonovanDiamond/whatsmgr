package whatsmgr

import (
	"context"
	"fmt"

	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
)

func (conn *Connection) SendMessage(message Message) error {
	out := waE2E.Message{
		Conversation: message.ContentBody,
	}
	jid, err := types.ParseJID(message.ChatJID)
	if err != nil {
		return fmt.Errorf("failed to parse ChatJID: %w", err)
	}
	resp, err := conn.client.SendMessage(context.Background(), jid, &out)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	message.MessageID = resp.ID
	senderJID := resp.Sender.String()
	message.SenderJID = &senderJID
	message.Timestamp = &resp.Timestamp
	conn.Callbacks.Message(message)
	return nil
}
