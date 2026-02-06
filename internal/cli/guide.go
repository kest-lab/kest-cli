package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var guideCmd = &cobra.Command{
	Use:   "guide",
	Short: "Show Kest Flow (.flow.md) tutorial and best practices",
	Long:  `Learn how to use Kest Flow to document and test your APIs simultaneously using Markdown.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(guideText)
	},
}

func init() {
	rootCmd.AddCommand(guideCmd)
}

const guideText = `
# Kest Flow (.flow.md) 🌊

Kest Flow is the core feature that enables "Document-as-Code" testing.
A .flow.md file is a standard Markdown document containing Kest code blocks.

## 📝 How to write a Flow

Create a file named "login.flow.md" and use code blocks like this:

` + "```" + `flow
@flow id=login-flow
@name Login Flow
` + "```" + `

` + "```" + `step
@id login
@name Login
# 1. First line: METHOD URL
POST /api/v1/auth/login

# 2. Headers (Optional)
Content-Type: application/json

# 3. Body (Leave an empty line after headers)
{
  "username": "admin",
  "password": "password123"
}

# 4. Capture variables for the next step
[Captures]
token = data.access_token

# 5. Logical assertions
[Asserts]
status == 200
duration < 500ms
` + "```" + `

## 🔗 Chaining Requests

Use captured variables with the {{variable_name}} syntax:

` + "```" + `kest
GET /api/v1/profile
Authorization: Bearer {{token}}

[Asserts]
status == 200
` + "```" + `

## 🚀 Running the flow

$ kest run login.flow.md

## 💡 Tips

- Use "http" or "json" if you prefer syntax highlighting in editors.
- Non-code content (text, images) is ignored by Kest, perfect for documentation.
- Use "kest history" to see the results of previous runs.

Keep Every Step Tested. 🦅
`
