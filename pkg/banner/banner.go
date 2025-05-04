package banner

import (
	"fmt"
	"io"
	"os"
	"sync"
)

const (
	red   = "\033[31m"
	reset = "\033[0m"
)

const nstreamBanner = `


` + red + `
███╗   ██╗███████╗████████╗██████╗ ███████╗ █████╗ ███╗   ███╗     █████╗ ██╗
████╗  ██║██╔════╝╚══██╔══╝██╔══██╗██╔════╝██╔══██╗████╗ ████║    ██╔══██╗██║
██╔██╗ ██║███████╗   ██║   ██████╔╝█████╗  ███████║██╔████╔██║    ███████║██║
██║╚██╗██║╚════██║   ██║   ██╔══██╗██╔══╝  ██╔══██║██║╚██╔╝██║    ██╔══██║██║
██║ ╚████║███████║   ██║   ██║  ██║███████╗██║  ██║██║ ╚═╝ ██║    ██║  ██║██║
╚═╝  ╚═══╝╚══════╝   ╚═╝   ╚═╝  ╚═╝╚══════╝╚═╝  ╚═╝╚═╝     ╚═╝    ╚═╝  ╚═╝╚═╝` + reset + `

                     AI agents for real time actions
                     
                     © 2024 NStream AI
                     https://nstream.ai
`

var (
	once     sync.Once
	writer   io.Writer = os.Stdout
	printed  bool
	printMux sync.Mutex
)

// PrintBanner prints the banner to the default writer (os.Stdout)
func PrintBanner() {
	printMux.Lock()
	defer printMux.Unlock()

	if !printed {
		fmt.Fprint(writer, nstreamBanner)
		fmt.Print("\n\n") // Add extra spacing after banner
		printed = true
	}
}

// GetBanner returns the banner as a string
func GetBanner() string {
	return nstreamBanner
}

// ResetBanner resets the printed flag (useful for testing)
func ResetBanner() {
	printMux.Lock()
	printed = false
	printMux.Unlock()
}
