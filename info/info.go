package info

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// Possible options, to make output sorted independent of config/cmd
var possibleOptions = []string{"logo", "userline", "userunderline", "kernel",
	"uptime", "shell", "memory"}

// Regexp matching empty lines, useful to make output more pretty
var emptyLinesRegex = regexp.MustCompile(`\n$`)

// Helper function, chains fmt.Sprintf and os.Expand(..., ColorExpand)
func formatAndColor(format string, args ...any) string {
	return os.Expand(
		fmt.Sprintf(format, args...),
		ColorExpand,
	)
}

// Returns processed info for pretty output
func GetInfoString(options map[string]string) (string, error) {
	// out string
	var output string

	// offset for printing labels
	var offset int

	// Info lines count, for calibrating newlines at the end
	var logolines, lines int

	for _, possibleOption := range possibleOptions {
		value, exists := options[possibleOption]
		if !exists {
			continue
		}

		if value == "false" {
			continue
		}

		switch possibleOption {
		case "logo":
			logo, err := getLogo(value)
			if err != nil {
				return "", err
			}

			output += os.Expand(logo.Logo, ColorExpand) +
				strings.Repeat("\x1b[F", logo.Lines-1)
			offset = logo.MaxLength + 2
			Colors["caccent"] = Colors[logo.AccentColor]
			logolines = logo.Lines

		case "userline":
			username := getRawUser()
			hostname, err := getRawHostname()
			if err != nil {
				return "", nil
			}

			output += formatAndColor(
				"\x1b[%vG${caccent}%v${creset}@${caccent}%v${creset}\n",
				offset,
				username,
				hostname,
			)
			lines++

		case "userunderline":
			username := getRawUser()
			hostname, err := getRawHostname()
			if err != nil {
				return "", nil
			}

			output += formatAndColor(
				"\x1b[%vG%v\n",
				offset,
				strings.Repeat("-", len(username)+len(hostname)),
			)
			lines++

		case "kernel":
			kernel, err := getRawKernel()
			if err != nil {
				return "", err
			}

			output += formatAndColor(
				"\x1b[%vG${caccent}Kernel${creset}: %v\n",
				offset,
				kernel,
			)
			lines++

		case "uptime":
			uptime, err := getRawUptime()
			if err != nil {
				return "", err
			}

			output += formatAndColor(
				"\x1b[%vG${caccent}Uptime${creset}: %v minutes\n",
				offset,
				int(uptime/60),
			)
			lines++

		case "shell":
			shell := getRawShell()

			output += formatAndColor(
				"\x1b[%vG${caccent}Shell${creset}: %v\n",
				offset,
				shell,
			)
			lines++

		case "memory":
			used, total, err := getRawMemory()
			if err != nil {
				return "", err
			}

			output += formatAndColor(
				"\x1b[%vG${caccent}Memory${creset}: %v / %v Mb\n",
				offset,
				used,
				total,
			)
			lines++
		}
	}

	if lines < logolines {
		output += strings.Repeat("\n", logolines-lines)
	}

	// return output, nil
	return emptyLinesRegex.ReplaceAllString(output, ""), nil
}