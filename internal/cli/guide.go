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

var docGuideCmd = &cobra.Command{
	Use:   "doc",
	Short: "Show API documentation generation guide",
	Long:  `Learn how to use Kest doc to generate and align API documentation from source code.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(docGuideText)
	},
}

func init() {
	rootCmd.AddCommand(guideCmd)
	guideCmd.AddCommand(docGuideCmd)
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

const docGuideText = `
# API Documentation Generation 📄

Kest CLI can scan your Go source code and generate high-quality API documentation automatically.

## 🚀 Basic Commands

### 📂 Scan a project
$ kest doc ./path/to/api -o ./docs

### 🧠 With AI Enhancement
Generates realistic examples, Mermaid diagrams, and permission summaries.
$ kest doc ./path/to/api -o ./docs --ai

### 🎯 Selective Scaning
Focus on a specific module to save time.
$ kest doc ./path/to/api -o ./docs -m module_name --ai

## 🌟 Superpowers (v0.6.1)

1. **Reality Seeding**: Use --ai to inject actual successful request/response data from your local Kest history into the documentation.
2. **Drift Detection**: Run with --verify in CI/CD to detect discrepancies between code and existing documentation.
3. **Interactive Portal**: Run with --serve to launch a beautiful local web portal to preview and test your APIs.
4. **Recursive Logic**: Scans Handler -> Service -> Repository to generate high-fidelity logic flow diagrams.

## 🧠 Anti-Hallucination

Kest ensures documentation perfectly matches your code:
- **No Ghost Wrapping**: Top-level arrays are not unnecessarily wrapped in objects.
- **Strict Key Adherence**: JSON keys match your source code tags exactly.
- **Deep Flow Analysis**: Sequence diagrams reflect implementation across layers.

鹰击长空，Keep Every Step Aligned. 🦅
`
