package progress

import (
	"fmt"
	"os"
	"runtime"

	"github.com/schollz/progressbar/v3"
)

func NewBar(fileSize int64, message string) *progressbar.ProgressBar {
	width := 50
	if runtime.GOOS == "windows" {
		width = 30
	}

	themeOption := progressbar.OptionSetTheme(progressbar.Theme{
		Saucer:        "=",
		SaucerHead:    ">",
		SaucerPadding: " ",
		BarStart:      "[",
		BarEnd:        "]",
	})

	return progressbar.NewOptions64(
		fileSize,
		progressbar.OptionSetDescription(message),
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(width),
		progressbar.OptionOnCompletion(func() {
			fmt.Fprint(os.Stderr, "\n")
		}),
		themeOption,
	)
}
