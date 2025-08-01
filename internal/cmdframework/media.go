package cmdframework

import (
	"context"
	"fmt"
	"io"
	"net/http"
	
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
	waProto "go.mau.fi/whatsmeow/proto/waE2E"
	"google.golang.org/protobuf/proto"
)

type MediaUploader struct {
	client ClientInterface
}

func NewMediaUploader(client ClientInterface) *MediaUploader {
	return &MediaUploader{client: client}
}

// Upload methods (without sending)
func (m *MediaUploader) UploadImage(ctx context.Context, imageData []byte) (UploadResponse, error) {
	return m.client.Upload(ctx, imageData, MediaImage)
}

func (m *MediaUploader) UploadVideo(ctx context.Context, videoData []byte) (UploadResponse, error) {
	return m.client.Upload(ctx, videoData, MediaVideo)
}

func (m *MediaUploader) UploadDocument(ctx context.Context, docData []byte, filename string) (UploadResponse, error) {
	return m.client.Upload(ctx, docData, MediaDocument)
}

func (m *MediaUploader) UploadAndSendImage(ctx context.Context, to types.JID, imageData []byte, caption string) error {
	uploaded, err := m.client.Upload(ctx, imageData, MediaImage)
	if err != nil {
		return fmt.Errorf("failed to upload image: %w", err)
	}
	
	msg := &waProto.Message{
		ImageMessage: &waProto.ImageMessage{
			Caption:       proto.String(caption),
			Mimetype:      proto.String("image/jpeg"),
			URL:           proto.String(uploaded.URL),
			DirectPath:    proto.String(uploaded.DirectPath),
			MediaKey:      uploaded.MediaKey,
			FileEncSHA256: uploaded.FileEncSHA256,
			FileSHA256:    uploaded.FileSHA256,
			FileLength:    proto.Uint64(uploaded.FileLength),
		},
	}
	
	_, err = m.client.SendMessage(ctx, to, msg)
	return err
}

func (m *MediaUploader) UploadAndSendVideo(ctx context.Context, to types.JID, videoData []byte, caption string) error {
	uploaded, err := m.client.Upload(ctx, videoData, MediaVideo)
	if err != nil {
		return fmt.Errorf("failed to upload video: %w", err)
	}
	
	msg := &waProto.Message{
		VideoMessage: &waProto.VideoMessage{
			Caption:       proto.String(caption),
			Mimetype:      proto.String("video/mp4"),
			URL:           proto.String(uploaded.URL),
			DirectPath:    proto.String(uploaded.DirectPath),
			MediaKey:      uploaded.MediaKey,
			FileEncSHA256: uploaded.FileEncSHA256,
			FileSHA256:    uploaded.FileSHA256,
			FileLength:    proto.Uint64(uploaded.FileLength),
		},
	}
	
	_, err = m.client.SendMessage(ctx, to, msg)
	return err
}

func (m *MediaUploader) UploadAndSendDocument(ctx context.Context, to types.JID, docData []byte, filename, caption string) error {
	uploaded, err := m.client.Upload(ctx, docData, MediaDocument)
	if err != nil {
		return fmt.Errorf("failed to upload document: %w", err)
	}
	
	msg := &waProto.Message{
		DocumentMessage: &waProto.DocumentMessage{
			Caption:       proto.String(caption),
			FileName:      proto.String(filename),
			Mimetype:      proto.String("application/octet-stream"),
			URL:           proto.String(uploaded.URL),
			DirectPath:    proto.String(uploaded.DirectPath),
			MediaKey:      uploaded.MediaKey,
			FileEncSHA256: uploaded.FileEncSHA256,
			FileSHA256:    uploaded.FileSHA256,
			FileLength:    proto.Uint64(uploaded.FileLength),
		},
	}
	
	_, err = m.client.SendMessage(ctx, to, msg)
	return err
}

func DownloadMedia(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to download media: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download media: status %d", resp.StatusCode)
	}
	
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read media data: %w", err)
	}
	
	return data, nil
}

func ConvertMediaType(mediaType MediaType) whatsmeow.MediaType {
	switch mediaType {
	case MediaImage:
		return whatsmeow.MediaImage
	case MediaVideo:
		return whatsmeow.MediaVideo
	case MediaDocument:
		return whatsmeow.MediaDocument
	case MediaAudio:
		return whatsmeow.MediaAudio
	default:
		return whatsmeow.MediaImage
	}
}