package cmdframework

import (
	"fmt"
	"strings"
)

type ResponseBuilder struct {
	parts []string
}

func NewResponseBuilder() *ResponseBuilder {
	return &ResponseBuilder{
		parts: make([]string, 0),
	}
}

func (r *ResponseBuilder) AddLine(text string) *ResponseBuilder {
	r.parts = append(r.parts, text)
	return r
}

func (r *ResponseBuilder) AddEmptyLine() *ResponseBuilder {
	r.parts = append(r.parts, "")
	return r
}

func (r *ResponseBuilder) AddBold(text string) *ResponseBuilder {
	r.parts = append(r.parts, fmt.Sprintf("*%s*", text))
	return r
}

func (r *ResponseBuilder) AddItalic(text string) *ResponseBuilder {
	r.parts = append(r.parts, fmt.Sprintf("_%s_", text))
	return r
}

func (r *ResponseBuilder) AddCode(text string) *ResponseBuilder {
	r.parts = append(r.parts, fmt.Sprintf("`%s`", text))
	return r
}

func (r *ResponseBuilder) AddCodeBlock(text string) *ResponseBuilder {
	r.parts = append(r.parts, fmt.Sprintf("```\n%s\n```", text))
	return r
}

func (r *ResponseBuilder) AddList(items ...string) *ResponseBuilder {
	for _, item := range items {
		r.parts = append(r.parts, fmt.Sprintf("• %s", item))
	}
	return r
}

func (r *ResponseBuilder) AddNumberedList(items ...string) *ResponseBuilder {
	for i, item := range items {
		r.parts = append(r.parts, fmt.Sprintf("%d. %s", i+1, item))
	}
	return r
}

func (r *ResponseBuilder) AddHeading(text string) *ResponseBuilder {
	return r.AddBold(text).AddEmptyLine()
}

func (r *ResponseBuilder) Build() string {
	return strings.Join(r.parts, "\n")
}

type ErrorResponse struct {
	Code    string
	Message string
	Details string
}

func (e ErrorResponse) String() string {
	builder := NewResponseBuilder()
	builder.AddLine(fmt.Sprintf("❌ *Error: %s*", e.Code))
	
	if e.Message != "" {
		builder.AddLine(e.Message)
	}
	
	if e.Details != "" {
		builder.AddEmptyLine()
		builder.AddLine(e.Details)
	}
	
	return builder.Build()
}

func Success(message string) string {
	return fmt.Sprintf("✅ %s", message)
}

func Error(message string) string {
	return fmt.Sprintf("❌ %s", message)
}

func Warning(message string) string {
	return fmt.Sprintf("⚠️ %s", message)
}

func Info(message string) string {
	return fmt.Sprintf("ℹ️ %s", message)
}

func Processing(message string) string {
	return fmt.Sprintf("⏳ %s", message)
}

type Template struct {
	name     string
	template string
}

var Templates = struct {
	InvalidCommand      Template
	MissingParameter    Template
	InvalidParameter    Template
	PermissionDenied    Template
	RateLimited         Template
	InternalError       Template
	FeatureNotAvailable Template
}{
	InvalidCommand: Template{
		name:     "InvalidCommand",
		template: "❌ Unknown command: `/%s`\nType `/help` for available commands.",
	},
	MissingParameter: Template{
		name:     "MissingParameter",
		template: "❌ Missing required parameter: %s\nUsage: `%s`",
	},
	InvalidParameter: Template{
		name:     "InvalidParameter",  
		template: "❌ Invalid %s: %s\nExpected: %s",
	},
	PermissionDenied: Template{
		name:     "PermissionDenied",
		template: "❌ Permission denied: %s",
	},
	RateLimited: Template{
		name:     "RateLimited",
		template: "⏱️ Please wait %d seconds before using this command again.",
	},
	InternalError: Template{
		name:     "InternalError",
		template: "❌ An internal error occurred. Please try again later.",
	},
	FeatureNotAvailable: Template{
		name:     "FeatureNotAvailable",
		template: "🚧 This feature is currently not available.",
	},
}

func (t Template) Format(args ...interface{}) string {
	return fmt.Sprintf(t.template, args...)
}