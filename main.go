package main

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"gopkg.in/src-d/go-git.v4/plumbing"

	git "gopkg.in/src-d/go-git.v4"
)

func main() {
	for _, x := range os.Args {
		if x == "-h" {
			showHelp()
			os.Exit(0)
		}
	}

	repo, err := git.PlainOpen(".")
	if err != nil {
		if err == git.ErrRepositoryNotExists {
			pr("No repository in the current path.")
			os.Exit(1)
		} else {
			pr("Error: %s", err.Error())
		}
		os.Exit(2)
	}

	head, err := repo.Head()
	if err != nil {
		pr("Error: %s. Have you made any commits?", err.Error())
		os.Exit(2)
	}

	ref, err := repo.Tags()
	if err != nil {
		pr("Error: %s", err.Error())
		os.Exit(2)
	}

	var tags []string
	ref.ForEach(func(p *plumbing.Reference) error {
		t := strings.Replace(p.Name().String(), "refs/tags/", "", 1)
		tags = append(tags, t)
		return nil
	})

	var last string
	if len(tags) < 1 {
		last = "v0.0.0"
	} else {
		sort.Strings(tags)
		last = tags[len(tags)-1]
	}

	pattern := `(v*)(\d+)(\.)(\d+)(\.)(\d+)(.*)`
	match, err := regexp.MatchString(pattern, last)
	if !match {
		pr("%s doesn't conform to a valid semantic version.", last)
		os.Exit(1)
	}

	user := os.Getenv("GITUSER")
	mail := os.Getenv("GITMAIL")
	pr("Latest version is %s", last)
	pr("Committing as %s <%s>", user, mail)
	if len(os.Args) == 1 {
		return
	}

	version, err := BumpVersion(last, os.Args[1])
	if err != nil {
		pr("Error bumping version: %s", err.Error())
		os.Exit(2)
	}

	tagger := NewTagger(user, mail)
	err = AddAnnotatedTag(repo, tagger, version, head.Strings()[1])
	if err != nil {
		pr("Error adding annotated tag: %s", err.Error())
		os.Exit(2)
	}

	pr("New version is %s", version)
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
