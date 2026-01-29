package output

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1)

	statusOKStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00")).
			Bold(true)

	statusErrStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Bold(true)

	borderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")).
			Padding(1)

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888"))
)

func PrintResponse(method, url string, status int, duration string, body []byte, recordID int64) {
	statusStr := fmt.Sprintf("%d", status)
	var sStyle lipgloss.Style
	if status >= 200 && status < 300 {
		sStyle = statusOKStyle
	} else {
		sStyle = statusErrStyle
	}

	headerText := fmt.Sprintf(" %s %s ", method, url)
	statusText := sStyle.Render(statusStr)
	durationText := infoStyle.Render(duration)

	var formattedBody string
	var obj interface{}
	if err := json.Unmarshal(body, &obj); err == nil {
		if prettyBody, err := json.MarshalIndent(obj, "", "  "); err == nil {
			formattedBody = string(prettyBody)
		} else {
			formattedBody = string(body)
		}
	} else {
		formattedBody = string(body)
	}

	content := fmt.Sprintf("Status: %s    Duration: %s\n\n%s", statusText, durationText, formattedBody)

	doc := strings.Builder{}
	doc.WriteString(titleStyle.Render(headerText) + "\n")
	doc.WriteString(borderStyle.Render(content) + "\n")

	if recordID > 0 {
		doc.WriteString(fmt.Sprintf(" ✓ Recorded as #%d\n", recordID))
	}

	fmt.Println(doc.String())
}
