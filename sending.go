package whatsmgr

import (
	"context"
	"fmt"

	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
)

func (conn *Connection) SendMessage(message Message, sendOnCallback bool) (Message, error) {
	out := waE2E.Message{
		Conversation: message.ContentBody,
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
