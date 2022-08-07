package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/mattn/go-colorable"
)

type Color string

const (
	Reset       Color = "\x1b[0m"
	Black       Color = "\x1b[30m"
	Red         Color = "\x1b[31m"
	Green       Color = "\x1b[32m"
	Yellow      Color = "\x1b[33m"
	Blue        Color = "\x1b[34m"
	Magenta     Color = "\x1b[35m"
	Cyan        Color = "\x1b[36m"
	White       Color = "\x1b[37m"
	BlackBold   Color = "\x1b[1;30m"
	RedBold     Color = "\x1b[1;31m"
	GreenBold   Color = "\x1b[1;32m"
	YellowBold  Color = "\x1b[1;33m"
	BlueBold    Color = "\x1b[1;34m"
	MagentaBold Color = "\x1b[1;35m"
	CyanBold    Color = "\x1b[1;36m"
	WhiteBold   Color = "\x1b[1;37m"
)

var (
	out                = colorable.NewColorableStdout()
	progressStart      time.Time
	progressMsg        string
	loadingTicker      *time.Ticker
	progressBarRunning bool
)

func BeginLoading(format string, a ...any) {
	FinishLoading()
	progressStart = time.Now()
	loadingTicker = time.NewTicker(time.Second / 2)
	progressMsg = fmt.Sprintf(format, a...)
	go func() {
		position := 0
		for {
			updateLoading(position)
			position++
			if loadingTicker == nil {
				return
			}
			_, ok := <-loadingTicker.C
			if !ok {
				return
			}
		}
	}()
}

func updateLoading(position int) {
	symbols := []string{"|", "/", "-", "\\"}
	fmt.Fprintf(out, "\r%s%s %s %ds%s", Cyan, symbols[position%4], progressMsg, int(time.Since(progressStart).Seconds()), Reset)
}

func CancelLoading() {
	if loadingTicker == nil {
		return
	}
	loadingTicker.Stop()
	loadingTicker = nil
	fmt.Fprintln(out)
}

func FinishLoading() {
	if loadingTicker == nil {
		return
	}
	loadingTicker.Stop()
	loadingTicker = nil

	totalDuration := time.Since(progressStart)
	var durationStr string
	if int(totalDuration.Seconds()) == 0 {
		durationStr = fmt.Sprintf("%dms", totalDuration.Milliseconds())
	} else {
		durationStr = fmt.Sprintf("%ds", int(totalDuration.Seconds()))
	}

	fmt.Fprintf(out, "\r%s√ %s %s   %s\n", Green, strings.TrimSuffix(progressMsg, "..."), durationStr, Reset)
}

func BeginProgressBar(format string, a ...any) {
	progressMsg = fmt.Sprintf(format, a...)
	progressStart = time.Now()
	fmt.Fprintf(out, "%s%s%s\n", Cyan, progressMsg, Reset)
	UpdateProgressBar(0)
}

func UpdateProgressBar(progress float64) {
	progressBarRunning = true
	stepSize := 0.03125
	fmt.Fprintf(out, "\r%s[", Cyan)
	for i := 0.0; i <= 1; i += stepSize {
		if i <= progress && i+stepSize > progress && progress > 0 && progress < 1 {
			fmt.Fprintf(out, ">")
		} else if i <= progress && progress > 0 {
			fmt.Fprintf(out, "=")
		} else {
			fmt.Fprintf(out, " ")
		}
	}
	fmt.Fprintf(out, "] %d%% %ds%s", int(progress*100), int(time.Since(progressStart).Seconds()), Reset)
}

func CancelProgressBar() {
	if progressBarRunning {
		fmt.Fprintln(out)
	}
	progressBarRunning = false
}

func FinishProgressBar() {
	if !progressBarRunning {
		return
	}
	fmt.Fprintf(out, "\x1b[1F%s%s%s\n", Green, strings.TrimSuffix(progressMsg, "...")+"   ", Reset)
	fmt.Fprintf(out, "%s[=================================] 100%% %ds%s\n", Green, int(time.Since(progressStart).Seconds()), Reset)
	progressBarRunning = false
}

func Print(format string, a ...any) {
	CancelLoading()
	CancelProgressBar()
	fmt.Fprintf(out, "%s\n", fmt.Sprintf(format, a...))
}

func PrintColor(color Color, format string, a ...any) {
	CancelLoading()
	CancelProgressBar()
	fmt.Fprintf(out, "%s%s%s\n", color, fmt.Sprintf(format, a...), Reset)
}

func Success(format string, a ...any) {
	PrintColor(Green, format, a...)
}

func Warn(format string, a ...any) {
	Print(string(Yellow)+"WARNING: "+string(Reset)+format, a...)
}

func Error(format string, a ...any) {
	Print(string(RedBold)+"ERROR: "+string(Reset)+format, a...)
}

func Clear() {
	fmt.Fprintf(out, "\033[H\033[2J")
}
