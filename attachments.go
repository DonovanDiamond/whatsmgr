package whatsmgr

import (
	"fmt"
	"mime"

	"go.mau.fi/whatsmeow/types/events"
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
