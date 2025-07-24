package utility

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	framework "github.com/asparkoffire/whatsapp-livetranslate-go/internal/cmdframework"
	"github.com/lrstanley/go-ytdlp"
)

type DownloadCommand struct{}

func NewDownloadCommand() *DownloadCommand {
	return &DownloadCommand{}
}

func (c *DownloadCommand) Execute(ctx *framework.Context) error {
	if len(ctx.Args) == 0 {
		return ctx.Handler.SendResponse(ctx.MessageInfo, framework.Error("Please provide a URL to download"))
	}

	url := ctx.Args[0]

	// Send initial status message
	ctx.Handler.SendResponse(ctx.MessageInfo, framework.Info("🔍 Analyzing URL..."))

	// Create temporary directory for downloads
	tempDir, err := os.MkdirTemp("", "whatsapp-download-*")
	if err != nil {
		return ctx.Handler.SendResponse(ctx.MessageInfo, framework.Error("Failed to create temp directory"))
	}
	defer os.RemoveAll(tempDir)

	// Configure output template
	outputTemplate := filepath.Join(tempDir, "download.%(ext)s")

	// Initialize ytdlp with options
	dl := ytdlp.New().
		FormatSort("res:720"). // Prefer 720p
		NoPlaylist().          // Download only single video, not playlist
		RestrictFilenames().   // Safe filenames
		Output(outputTemplate).
		Cookies(os.Getenv("COOKIES_PATH"))

	// Update status
	ctx.Handler.SendResponse(ctx.MessageInfo, framework.Info("⬇️ Downloading media..."))

	// Download the media
	_, err = dl.Run(context.Background(), url)
	if err != nil {
		return ctx.Handler.SendResponse(ctx.MessageInfo, framework.Error(fmt.Sprintf("Download failed: %v", err)))
	}

	// Find the downloaded file
	files, err := filepath.Glob(filepath.Join(tempDir, "download.*"))
	if err != nil || len(files) == 0 {
		return ctx.Handler.SendResponse(ctx.MessageInfo, framework.Error("Downloaded file not found"))
	}
	outputFile := files[0]

	// Read the file
	data, err := os.ReadFile(outputFile)
	if err != nil {
		return ctx.Handler.SendResponse(ctx.MessageInfo, framework.Error("Failed to read downloaded file"))
	}

	// Check file size (WhatsApp has limits)
	const maxSize = 16 * 1024 * 1024 // 16MB limit for WhatsApp
	if len(data) > maxSize {
		return ctx.Handler.SendResponse(ctx.MessageInfo, framework.Error(fmt.Sprintf("File too large (%.1f MB). Maximum size is 16 MB", float64(len(data))/(1024*1024))))
	}

	// Prepare caption
	caption := fmt.Sprintf("📥 Downloaded from: %s", url)

	// Upload based on type
	uploader := framework.NewMediaUploader(ctx.Handler.GetClient())

	if strings.HasSuffix(strings.ToLower(outputFile), ".mp4") || strings.HasSuffix(strings.ToLower(outputFile), ".webm") {
		// Upload as video
		resp, err := uploader.UploadVideo(ctx.Context, data)
		if err != nil {
			return ctx.Handler.SendResponse(ctx.MessageInfo, framework.Error(fmt.Sprintf("Failed to upload video: %v", err)))
		}

		// Send video with caption
		return ctx.Handler.SendVideo(ctx.MessageInfo, resp, caption)
	} else {
		// Upload as image or document
		if isImageFile(outputFile) {
			resp, err := uploader.UploadImage(ctx.Context, data)
			if err != nil {
				return ctx.Handler.SendResponse(ctx.MessageInfo, framework.Error(fmt.Sprintf("Failed to upload image: %v", err)))
			}
			return ctx.Handler.SendImage(ctx.MessageInfo, resp, caption)
		} else {
			// Send as document if not recognized media type
			resp, err := uploader.UploadDocument(ctx.Context, data, filepath.Base(outputFile))
			if err != nil {
				return ctx.Handler.SendResponse(ctx.MessageInfo, framework.Error(fmt.Sprintf("Failed to upload document: %v", err)))
			}
			return ctx.Handler.SendDocument(ctx.MessageInfo, resp, caption)
		}
	}
}

func (c *DownloadCommand) Metadata() *framework.Metadata {
	return &framework.Metadata{
		Name:        "download",
		Aliases:     []string{"dl", "ytdl"},
		Description: "Download media from various platforms",
		Category:    "Utility",
		Usage:       "/download <url>",
		Examples: []string{
			"/download https://www.youtube.com/watch?v=...",
			"/dl https://www.instagram.com/p/...",
			"/dl https://twitter.com/user/status/...",
		},
	}
}

func isImageFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	imageExts := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
	for _, imgExt := range imageExts {
		if ext == imgExt {
			return true
		}
	}
	return false
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// downloadWithHTTP is a fallback for direct image/video URLs
func downloadWithHTTP(url string) ([]byte, string, error) {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("bad status: %s", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	// Get filename from URL or content disposition
	filename := filepath.Base(url)
	if cd := resp.Header.Get("Content-Disposition"); cd != "" {
		if idx := strings.Index(cd, "filename="); idx != -1 {
			filename = strings.Trim(cd[idx+9:], "\"")
		}
	}

	return data, filename, nil
}
