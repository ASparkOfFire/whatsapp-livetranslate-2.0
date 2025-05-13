# WhatsApp Live Translation Bot using Gemini and Whatsmeow library

## Q. How to add more languages?
- A. Add it into the `SupportedLanguages` map inside `internal/constants/constants.go` file.

## Q. How to add new commands?
- A. Look into `handleMessage` function inside of `internal/services/messagehandler/event_handler.go` file.
