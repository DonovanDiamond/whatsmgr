package whatsmgr

import (
	"context"
	"errors"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
)

type Connection struct {
	Number    string
	DBPath    string
	MediaPath string
	Log       Logger

	Callbacks Callbacks

	client *whatsmeow.Client
	ctx    context.Context
}

var ErrAlreadyConnected = errors.New("already connected")

func (conn *Connection) Connect(ctx context.Context) error {
	if conn.client.IsConnected() || conn.client.IsLoggedIn() {
		conn.Callbacks.ConnStatus(ConnStatusConnected)
		return ErrAlreadyConnected
	}
	conn.ctx = ctx
	if conn.Number == "" {
		return fmt.Errorf("missing connection number")
	}
	if conn.DBPath == "" {
		return fmt.Errorf("missing db path")
	}
	if conn.MediaPath == "" {
		return fmt.Errorf("missing media path")
	}
	storeConainter, err := sqlstore.New(
		ctx,
		"sqlite3",
		fmt.Sprintf("file:%s?_foreign_keys=on", conn.DBPath),
		conn.Log.Sub("sqlite"),
	)
	if err != nil {
		return fmt.Errorf("failed to create sqlstore: %w", err)
	}

	device, err := storeConainter.GetFirstDevice(ctx)
	if err != nil {
		return fmt.Errorf("failed to get device from sqlstore: %w", err)
	}

	conn.client = whatsmeow.NewClient(device, conn.Log.Sub("client"))
	conn.client.EnableAutoReconnect = true

	err = conn.runQRHandler()
	if err != nil {
		return fmt.Errorf("failed to start QR code handler: %w", err)
	}

	conn.client.AddEventHandler(conn.handleEvent)

	err = conn.client.Connect()
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}

	return nil
}

func (conn *Connection) Disconnect() {
	conn.client.Disconnect()
}

func (conn *Connection) IsConnected() bool {
	return conn.client.IsConnected()
}

func (conn *Connection) IsLoggedIn() bool {
	return conn.client.IsLoggedIn()
}

func (conn *Connection) runQRHandler() error {
	ch, err := conn.client.GetQRChannel(context.Background())
	if errors.Is(err, whatsmeow.ErrQRStoreContainsID) {
		return nil
	} else if err != nil {
		return fmt.Errorf("failed to get qr code channel: %w", err)
	}

	go func() {
		for evt := range ch {
			if evt.Event == "code" {
				conn.Callbacks.QRCode(evt.Code)
			}
		}
	}()
	return nil
}

func (conn *Connection) SyncAllContacts() error {
	contacts, err := conn.client.Store.Contacts.GetAllContacts(conn.ctx)
	if err != nil {
		return fmt.Errorf("failed to get all contacts: %w", err)
	}

	for jid, contactInfo := range contacts {
		contact := Contact{
			JID: jid.String(),
		}
		if contactInfo.FullName != "" {
			contact.ContactName = &contactInfo.FullName
		}
		if contactInfo.PushName != "" {
			contact.PushName = &contactInfo.PushName
		}
		photo, _ := conn.pullProfilePhoto(jid)
		if photo != "" {
			contact.ProfilePhoto = &photo
		}
		conn.Callbacks.Contact(contact)
	}
	return nil
}
