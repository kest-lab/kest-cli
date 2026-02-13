package support

import (
	"fmt"

	"github.com/fatih/color"
)

// PrintBanner prints the Eogo ASCII banner to console
func PrintBanner(version string) {
	bannerColor := color.New(color.FgCyan, color.Bold)
	secondaryColor := color.New(color.FgHiBlue)

	banner := `
  ______   ____    ____  ____ 
 |  ____| / __ \  / ___|/ __ \
 | |__   | |  | || |  _| |  | |
 |  __|  | |  | || | |_ | |  | |
 | |____ | |__| || |__| | |__| |
 |______| \____/  \____|\____/ 
`
	bannerColor.Print(banner)
	fmt.Println()
	secondaryColor.Printf("  Eogo Framework %s - Elegant Go Development\n", version)
	fmt.Println()
}
