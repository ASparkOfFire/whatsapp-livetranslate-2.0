# Adding New Commands to WhatsApp LiveTranslate

This guide explains how to add new commands to the WhatsApp LiveTranslate bot using the new command framework.

## Overview

The command system is designed to be modular and extensible. Each command is a self-contained unit that implements the `Command` interface.

## Step-by-Step Guide

### 1. Create Your Command File

Create a new file in the appropriate directory under `internal/handlers/`:
- `admin/` - Administrative commands
- `fun/` - Entertainment commands
- `translation/` - Translation-related commands
- `utility/` - General utility commands

### 2. Implement the Command Interface

Your command must implement the `Command` interface:

```go
type Command interface {
    Execute(ctx *Context) error
    Metadata() *Metadata
}
```

### 3. Basic Command Example

Here's a simple command that responds with a static message:

```go
package utility

import (
    framework "github.com/asparkoffire/whatsapp-livetranslate-go/internal/cmdframework"
)

type GreetCommand struct{}

func NewGreetCommand() *GreetCommand {
    return &GreetCommand{}
}

func (c *GreetCommand) Execute(ctx *framework.Context) error {
    greeting := "Hello! 👋 Welcome to WhatsApp LiveTranslate!"
    return ctx.Handler.SendResponse(ctx.MessageInfo, greeting)
}

func (c *GreetCommand) Metadata() *framework.Metadata {
    return &framework.Metadata{
        Name:        "greet",
        Description: "Send a greeting message",
        Category:    "Utility",
        Usage:       "/greet",
    }
}
```

### 4. Command with Parameters

For commands that accept arguments:

```go
package fun

import (
    "fmt"
    "strings"
    
    framework "github.com/asparkoffire/whatsapp-livetranslate-go/internal/cmdframework"
)

type EchoCommand struct{}

func NewEchoCommand() *EchoCommand {
    return &EchoCommand{}
}

func (c *EchoCommand) Execute(ctx *framework.Context) error {
    if len(ctx.Args) == 0 {
        return ctx.Handler.SendResponse(ctx.MessageInfo, 
            framework.Error("Please provide text to echo"))
    }
    
    // ctx.RawArgs contains the full argument string
    echoed := strings.ToUpper(ctx.RawArgs)
    return ctx.Handler.SendResponse(ctx.MessageInfo, echoed)
}

func (c *EchoCommand) Metadata() *framework.Metadata {
    return &framework.Metadata{
        Name:        "echo",
        Description: "Echo text in uppercase",
        Category:    "Fun",
        Usage:       "/echo <text>",
        Examples:    []string{"/echo hello world"},
        Parameters: []framework.Parameter{
            {
                Name:        "text",
                Type:        framework.StringParam,
                Description: "Text to echo",
                Required:    true,
            },
        },
    }
}
```

### 5. Media Command Example

For commands that send images, videos, or documents:

```go
func (c *MyMediaCommand) Execute(ctx *framework.Context) error {
    // Generate or fetch your media data
    imageData, err := generateImage()
    if err != nil {
        return ctx.Handler.SendResponse(ctx.MessageInfo, 
            framework.Error(fmt.Sprintf("Failed to generate image: %v", err)))
    }
    
    // Send the media
    caption := "Here's your generated image!"
    return ctx.Handler.SendMedia(ctx.MessageInfo, framework.MediaImage, imageData, caption)
}
```

### 6. Register Your Command

Add your command to the `InitializeCommands` function in `event_handler_new.go`:

```go
// Register your new command
if err := registry.Register(utility.NewGreetCommand()); err != nil {
    return fmt.Errorf("failed to register greet command: %w", err)
}
```

### 7. Add Middleware (Optional)

If your command requires owner permissions:

```go
// In InitializeCommands, after registering all commands
ownerCommands := []string{"ping", "setmodel", "greet"} // Add your command name
for _, cmdName := range ownerCommands {
    if cmd, exists := registry.Get(cmdName); exists {
        wrappedCmd := framework.WithMiddleware(cmd, framework.RequireOwner())
        registry.Register(wrappedCmd)
    }
}
```

## Command Features

### Response Formatting

Use the response builder for consistent formatting:

```go
builder := framework.NewResponseBuilder()
builder.AddHeading("Results")
builder.AddList("Item 1", "Item 2", "Item 3")
builder.AddEmptyLine()
builder.AddCode("some code")

return ctx.Handler.SendResponse(ctx.MessageInfo, builder.Build())
```

### Error Handling

Use the built-in error formatting functions:

```go
framework.Success("Operation completed successfully")
framework.Error("Something went wrong")
framework.Warning("This might not work as expected")
framework.Info("Here's some information")
framework.Processing("Working on it...")
```

### Parameter Validation

Add validators to your parameters:

```go
Parameters: []framework.Parameter{
    {
        Name:        "count",
        Type:        framework.IntParam,
        Description: "Number of items (1-100)",
        Required:    true,
        Validator: func(value string) error {
            n, err := strconv.Atoi(value)
            if err != nil {
                return err
            }
            if n < 1 || n > 100 {
                return fmt.Errorf("count must be between 1 and 100")
            }
            return nil
        },
    },
}
```

### Using Base Command Types

The framework provides base types for common patterns:

1. **SimpleCommand** - For basic text responses
2. **ParameterizedCommand** - For commands with validated parameters
3. **MediaCommand** - For commands that generate and send media

## Best Practices

1. **Keep commands focused** - Each command should do one thing well
2. **Use appropriate categories** - Place commands in logical categories
3. **Provide good help text** - Include usage examples and parameter descriptions
4. **Handle errors gracefully** - Always provide user-friendly error messages
5. **Use consistent formatting** - Leverage the response builder utilities
6. **Test your commands** - Write unit tests for command logic

## Available Services

Commands have access to these services through the context:

- `ctx.Handler.GetTranslator()` - Translation service
- `ctx.Handler.GetImageGenerator()` - AI image generation
- `ctx.Handler.GetMemeGenerator()` - Meme fetching service
- `ctx.Handler.GetLangDetector()` - Language detection
- `ctx.Handler.GetClient()` - WhatsApp client for advanced operations

## Example: Complete Command

Here's a complete example of a dice rolling command:

```go
package fun

import (
    "fmt"
    "math/rand"
    "strconv"
    "time"
    
    framework "github.com/asparkoffire/whatsapp-livetranslate-go/internal/cmdframework"
)

type DiceCommand struct{}

func NewDiceCommand() *DiceCommand {
    return &DiceCommand{}
}

func (c *DiceCommand) Execute(ctx *framework.Context) error {
    sides := 6 // default
    
    if len(ctx.Args) > 0 {
        s, err := strconv.Atoi(ctx.Args[0])
        if err != nil {
            return ctx.Handler.SendResponse(ctx.MessageInfo, 
                framework.Error("Please provide a valid number"))
        }
        if s < 2 || s > 100 {
            return ctx.Handler.SendResponse(ctx.MessageInfo, 
                framework.Error("Dice sides must be between 2 and 100"))
        }
        sides = s
    }
    
    rand.Seed(time.Now().UnixNano())
    result := rand.Intn(sides) + 1
    
    response := fmt.Sprintf("🎲 You rolled a %d on a d%d!", result, sides)
    return ctx.Handler.SendResponse(ctx.MessageInfo, response)
}

func (c *DiceCommand) Metadata() *framework.Metadata {
    return &framework.Metadata{
        Name:        "dice",
        Aliases:     []string{"roll", "d"},
        Description: "Roll a dice",
        Category:    "Fun",
        Usage:       "/dice [sides]",
        Examples: []string{
            "/dice",
            "/dice 20",
            "/roll 100",
        },
        Parameters: []framework.Parameter{
            {
                Name:        "sides",
                Type:        framework.IntParam,
                Description: "Number of sides on the dice (2-100)",
                Required:    false,
                Default:     6,
            },
        },
    }
}
```

## Conclusion

The command framework makes it easy to add new functionality to the bot. Follow this guide, use the existing commands as examples, and your new command will integrate seamlessly with the bot's architecture.