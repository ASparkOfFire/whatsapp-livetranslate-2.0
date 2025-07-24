package utility

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	framework "github.com/asparkoffire/whatsapp-livetranslate-go/internal/cmdframework"
	"github.com/lrstanley/go-ytdlp"
)

type DownloadCommand struct {
	mu              sync.Mutex
	isDownloading   bool
	lastDownloadTime time.Time
}

var downloadInstance = &DownloadCommand{}

func NewDownloadCommand() *DownloadCommand {
	return downloadInstance
}

func (c *DownloadCommand) Execute(ctx *framework.Context) error {
	if len(ctx.Args) == 0 {
		return ctx.Handler.SendResponse(ctx.MessageInfo, framework.Error("Please provide a URL to download"))
	}

	// Check rate limiting
	c.mu.Lock()
	
	// Check if download is already in progress
	if c.isDownloading {
		c.mu.Unlock()
		return ctx.Handler.SendResponse(ctx.MessageInfo, framework.Warning("⏳ A download is already in progress. Please wait for it to complete."))
	}
	
	// Check cooldown period (1 minute)
	if !c.lastDownloadTime.IsZero() {
		timeSinceLastDownload := time.Since(c.lastDownloadTime)
		if timeSinceLastDownload < 1*time.Minute {
			remainingTime := 1*time.Minute - timeSinceLastDownload
			c.mu.Unlock()
			return ctx.Handler.SendResponse(ctx.MessageInfo, framework.Warning(fmt.Sprintf("⏱️ Please wait %d seconds before downloading again.", int(remainingTime.Seconds()))))
		}
	}
	
	// Mark as downloading
	c.isDownloading = true
	c.mu.Unlock()
	
	// Ensure we mark as not downloading when done
	defer func() {
		c.mu.Lock()
		c.isDownloading = false
		c.lastDownloadTime = time.Now()
		c.mu.Unlock()
		fmt.Printf("[DOWNLOAD] Download completed, cooldown started\n")
	}()

	url := ctx.Args[0]
	fmt.Printf("[DOWNLOAD] Starting download for URL: %s\n", url)

	// Send initial status message
	ctx.Handler.SendResponse(ctx.MessageInfo, framework.Info("🔍 Analyzing URL..."))

	// Create temporary directory for downloads
	tempDir, err := os.MkdirTemp("", "whatsapp-download-*")
	if err != nil {
		fmt.Printf("[DOWNLOAD] Failed to create temp directory: %v\n", err)
		// Reset download state on early error
		c.mu.Lock()
		c.isDownloading = false
		c.mu.Unlock()
		return ctx.Handler.SendResponse(ctx.MessageInfo, framework.Error("Failed to create temp directory"))
	}
	defer func() {
		fmt.Printf("[DOWNLOAD] Cleaning up temp directory: %s\n", tempDir)
		os.RemoveAll(tempDir)
	}()
	fmt.Printf("[DOWNLOAD] Created temp directory: %s\n", tempDir)

	// Configure output template
	outputTemplate := filepath.Join(tempDir, "download.%(ext)s")
	fmt.Printf("[DOWNLOAD] Output template: %s\n", outputTemplate)

	// Initialize ytdlp with options
	dl := ytdlp.New().
		Format("best[height<=720]/best"). // Better format selection
		NoPlaylist().                     // Download only single video, not playlist
		RestrictFilenames().              // Safe filenames
		Output(outputTemplate).
		NoCheckCertificates(). // Skip certificate verification
		Verbose()              // Enable verbose logging

	// Add cookies if available
	if cookiesPath := os.Getenv("COOKIES_PATH"); cookiesPath != "" {
		fmt.Printf("[DOWNLOAD] Using cookies from: %s\n", cookiesPath)
		dl = dl.Cookies(cookiesPath)
	}

	// Update status
	ctx.Handler.SendResponse(ctx.MessageInfo, framework.Info("⬇️ Downloading media..."))

	// Download the media
	fmt.Printf("[DOWNLOAD] Running yt-dlp...\n")
	result, err := dl.Run(context.Background(), url)
	if err != nil {
		fmt.Printf("[DOWNLOAD] yt-dlp failed: %v\n", err)
		return ctx.Handler.SendResponse(ctx.MessageInfo, framework.Error(fmt.Sprintf("Download failed: %v", err)))
	}

	if result != nil {
		fmt.Printf("[DOWNLOAD] yt-dlp exit code: %d\n", result.ExitCode)
		if result.Stdout != "" {
			fmt.Printf("[DOWNLOAD] yt-dlp stdout:\n%s\n", result.Stdout)
		}
		if result.Stderr != "" {
			fmt.Printf("[DOWNLOAD] yt-dlp stderr:\n%s\n", result.Stderr)
		}
	}

	// Wait a moment for file to be fully written
	time.Sleep(500 * time.Millisecond)

	// Find the downloaded file
	files, err := filepath.Glob(filepath.Join(tempDir, "download.*"))
	if err != nil {
		fmt.Printf("[DOWNLOAD] Glob error: %v\n", err)
		return ctx.Handler.SendResponse(ctx.MessageInfo, framework.Error("Error finding downloaded file"))
	}

	fmt.Printf("[DOWNLOAD] Found %d files in temp directory\n", len(files))
	for i, f := range files {
		info, _ := os.Stat(f)
		if info != nil {
			fmt.Printf("[DOWNLOAD] File %d: %s (size: %d bytes)\n", i+1, f, info.Size())
		}
	}

	if len(files) == 0 {
		// Try alternative patterns
		alternativePatterns := []string{
			filepath.Join(tempDir, "*.mp4"),
			filepath.Join(tempDir, "*.webm"),
			filepath.Join(tempDir, "*.mkv"),
			filepath.Join(tempDir, "*"),
		}

		for _, pattern := range alternativePatterns {
			altFiles, _ := filepath.Glob(pattern)
			if len(altFiles) > 0 {
				files = altFiles
				fmt.Printf("[DOWNLOAD] Found files with pattern %s\n", pattern)
				break
			}
		}

		if len(files) == 0 {
			// List all files in temp directory for debugging
			allFiles, _ := os.ReadDir(tempDir)
			fmt.Printf("[DOWNLOAD] All files in temp directory:\n")
			for _, f := range allFiles {
				fmt.Printf("[DOWNLOAD]   - %s\n", f.Name())
			}
			return ctx.Handler.SendResponse(ctx.MessageInfo, framework.Error("Downloaded file not found"))
		}
	}
	outputFile := files[0]
	fmt.Printf("[DOWNLOAD] Selected file: %s\n", outputFile)

	// Get file info
	fileInfo, err := os.Stat(outputFile)
	if err != nil {
		fmt.Printf("[DOWNLOAD] Failed to stat file: %v\n", err)
		return ctx.Handler.SendResponse(ctx.MessageInfo, framework.Error("Failed to access downloaded file"))
	}
	fmt.Printf("[DOWNLOAD] File size: %d bytes (%.2f MB)\n", fileInfo.Size(), float64(fileInfo.Size())/(1024*1024))

	// Read the file
	data, err := os.ReadFile(outputFile)
	if err != nil {
		fmt.Printf("[DOWNLOAD] Failed to read file: %v\n", err)
		return ctx.Handler.SendResponse(ctx.MessageInfo, framework.Error("Failed to read downloaded file"))
	}
	fmt.Printf("[DOWNLOAD] Read %d bytes from file\n", len(data))

	// Check if file is actually empty
	if len(data) == 0 {
		fmt.Printf("[DOWNLOAD] ERROR: File is empty!\n")
		return ctx.Handler.SendResponse(ctx.MessageInfo, framework.Error("Downloaded file is empty"))
	}

	// Check file size (WhatsApp has limits)
	const maxSize = 16 * 1024 * 1024 // 16MB limit for WhatsApp
	if len(data) > maxSize {
		fmt.Printf("[DOWNLOAD] File too large: %d bytes\n", len(data))
		return ctx.Handler.SendResponse(ctx.MessageInfo, framework.Error(fmt.Sprintf("File too large (%.1f MB). Maximum size is 16 MB", float64(len(data))/(1024*1024))))
	}

	// Prepare caption
	caption := fmt.Sprintf("📥 Downloaded from: %s", url)

	// Upload based on type
	uploader := framework.NewMediaUploader(ctx.Handler.GetClient())
	extension := strings.ToLower(filepath.Ext(outputFile))
	fmt.Printf("[DOWNLOAD] File extension: %s\n", extension)

	if extension == ".mp4" || extension == ".webm" || extension == ".mkv" || extension == ".avi" {
		fmt.Printf("[DOWNLOAD] Uploading as video...\n")
		// Upload as video
		resp, err := uploader.UploadVideo(ctx.Context, data)
		if err != nil {
			fmt.Printf("[DOWNLOAD] Video upload failed: %v\n", err)
			return ctx.Handler.SendResponse(ctx.MessageInfo, framework.Error(fmt.Sprintf("Failed to upload video: %v", err)))
		}
		fmt.Printf("[DOWNLOAD] Video uploaded successfully: URL=%s, DirectPath=%s, FileLength=%d\n", resp.URL, resp.DirectPath, resp.FileLength)

		// Send video with caption
		err = ctx.Handler.SendVideo(ctx.MessageInfo, resp, caption)
		if err != nil {
			fmt.Printf("[DOWNLOAD] Failed to send video message: %v\n", err)
			return err
		}
		fmt.Printf("[DOWNLOAD] Video sent successfully\n")
		return nil
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
