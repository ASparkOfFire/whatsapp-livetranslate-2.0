package command

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type SimpleCommand struct {
	Meta     Metadata
	Response func(ctx *Context) string
}

func (c *SimpleCommand) Execute(ctx *Context) error {
	response := c.Response(ctx)
	return ctx.Handler.SendResponse(ctx.MessageInfo, response)
}

func (c *SimpleCommand) Metadata() *Metadata {
	return &c.Meta
}

type ParameterizedCommand struct {
	Meta      Metadata
	Handler   func(ctx *Context, params map[string]interface{}) error
}

func (c *ParameterizedCommand) Execute(ctx *Context) error {
	params, err := c.parseParameters(ctx)
	if err != nil {
		return ctx.Handler.SendResponse(ctx.MessageInfo, fmt.Sprintf("❌ %s", err.Error()))
	}
	
	return c.Handler(ctx, params)
}

func (c *ParameterizedCommand) Metadata() *Metadata {
	return &c.Meta
}

func (c *ParameterizedCommand) parseParameters(ctx *Context) (map[string]interface{}, error) {
	params := make(map[string]interface{})
	
	for i, param := range c.Meta.Parameters {
		var value string
		
		if i < len(ctx.Args) {
			value = ctx.Args[i]
		} else if param.Required {
			return nil, fmt.Errorf("missing required parameter: %s", param.Name)
		} else if param.Default != nil {
			params[param.Name] = param.Default
			continue
		} else {
			continue
		}
		
		if param.Validator != nil {
			if err := param.Validator(value); err != nil {
				return nil, fmt.Errorf("invalid %s: %v", param.Name, err)
			}
		}
		
		parsedValue, err := parseValue(value, param.Type)
		if err != nil {
			return nil, fmt.Errorf("invalid %s: %v", param.Name, err)
		}
		
		params[param.Name] = parsedValue
	}
	
	return params, nil
}

func parseValue(value string, typ ParameterType) (interface{}, error) {
	switch typ {
	case StringParam:
		return value, nil
	case IntParam:
		return strconv.Atoi(value)
	case FloatParam:
		return strconv.ParseFloat(value, 64)
	case BoolParam:
		lower := strings.ToLower(value)
		return lower == "true" || lower == "yes" || lower == "1", nil
	case DurationParam:
		return time.ParseDuration(value)
	default:
		return value, nil
	}
}

type MediaCommand struct {
	Meta           Metadata
	GenerateMedia  func(ctx *Context) ([]byte, MediaType, string, error)
}

func (c *MediaCommand) Execute(ctx *Context) error {
	data, mediaType, caption, err := c.GenerateMedia(ctx)
	if err != nil {
		return ctx.Handler.SendResponse(ctx.MessageInfo, fmt.Sprintf("❌ %s", err.Error()))
	}
	
	return ctx.Handler.SendMedia(ctx.MessageInfo, mediaType, data, caption)
}

func (c *MediaCommand) Metadata() *Metadata {
	return &c.Meta
}

type Middleware func(next Command) Command

type middlewareCommand struct {
	next Command
	fn   func(ctx *Context, next Command) error
}

func (m *middlewareCommand) Execute(ctx *Context) error {
	return m.fn(ctx, m.next)
}

func (m *middlewareCommand) Metadata() *Metadata {
	return m.next.Metadata()
}

func WithMiddleware(cmd Command, middleware ...Middleware) Command {
	for i := len(middleware) - 1; i >= 0; i-- {
		cmd = middleware[i](cmd)
	}
	return cmd
}

func RequireOwner() Middleware {
	return func(next Command) Command {
		return &middlewareCommand{
			next: next,
			fn: func(ctx *Context, next Command) error {
				if !ctx.MessageInfo.IsFromMe {
					return ctx.Handler.SendResponse(ctx.MessageInfo, "❌ This command requires owner permissions")
				}
				return next.Execute(ctx)
			},
		}
	}
}

func RateLimit(perMinute int) Middleware {
	lastUsed := make(map[string]time.Time)
	
	return func(next Command) Command {
		return &middlewareCommand{
			next: next,
			fn: func(ctx *Context, next Command) error {
				key := ctx.MessageInfo.Sender.String()
				now := time.Now()
				
				if last, exists := lastUsed[key]; exists {
					elapsed := now.Sub(last)
					minInterval := time.Minute / time.Duration(perMinute)
					
					if elapsed < minInterval {
						remaining := minInterval - elapsed
						return ctx.Handler.SendResponse(ctx.MessageInfo, 
							fmt.Sprintf("⏱️ Please wait %d seconds before using this command again", 
								int(remaining.Seconds())))
					}
				}
				
				lastUsed[key] = now
				return next.Execute(ctx)
			},
		}
	}
}