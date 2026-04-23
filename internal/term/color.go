package term

import "os"

func SupportsColor() bool {
	noColor := os.Getenv("NO_COLOR")
	if noColor != "" {
		return false
	}
	info, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (info.Mode() & os.ModeCharDevice) != 0
}

func Colorize(text string, colorCode string) string {
	if !SupportsColor() {
		return text
	}
	return "\u001b[" + colorCode + "m" + text + "\u001b[0m"
}

func ColorStatus(status string, mark string) string {
	switch status {
	case "ok":
		return Colorize(mark, "32")
	case "fail":
		return Colorize(mark, "31")
	default:
		return Colorize(mark, "33")
	}
}
