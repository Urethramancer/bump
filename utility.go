package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func pr(f string, v ...interface{}) {
	fmt.Printf(f+"\n", v...)
}

func showHelp() {
	pr("%s increases the version number of a semver-conformant git repository, based on the most recent annotated tag.\n", os.Args[0])
	pr("Usage:")
	pr("%s <segment>", os.Args[0])
	pr("Possible keywords are 'major', 'minor' and 'patch' with the shortcuts 'ma', 'mi' or 'p'.")
}

func isSemVer(ver string) bool {
	pattern := `(v*)(\d+)(\.)(\d+)(\.)(\d+)(.*)`
	match, err := regexp.MatchString(pattern, ver)
	if err != nil {
		return false
	}

	return match
}

// BumpVersion takes a version string and increases a segment by 1, depending on the part argument.
// major, ma: increase the first segment
// minor, mi: increase the second segment
// patch, p: increase the third (last) segment
// If it starts with "v" it's preserved in the resulting version, to maintain consistency.
func BumpVersion(ver, part string) (string, error) {
	bump := ver
	v := false
	if strings.HasPrefix(bump, "v") {
		v = true
		bump = bump[1:]
	}

	a := strings.Split(bump, ".")
	major, err := strconv.Atoi(a[0])
	if err != nil {
		return "", err
	}

	minor, err := strconv.Atoi(a[1])
	if err != nil {
		return "", err
	}

	patch, err := strconv.Atoi(a[2])
	if err != nil {
		return "", err
	}

	switch part {
	case "major", "ma":
		major++
		minor = 0
		patch = 0
	case "minor", "mi":
		minor++
		patch = 0
	case "patch", "p":
		patch++
	default:
		pr("Unknown command. Bailing out.")
	}

	var version string
	if v {
		version = fmt.Sprintf("v%d.%d.%d", major, minor, patch)
	} else {
		version = fmt.Sprintf("%d.%d.%d", major, minor, patch)
	}

	return version, nil
}
