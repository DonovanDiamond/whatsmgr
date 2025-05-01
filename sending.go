package whatsmgr

import (
	"context"
	"errors"
	"fmt"
	"time"

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
