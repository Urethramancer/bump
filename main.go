package main

import (
	"os"

	git "gopkg.in/src-d/go-git.v4"
)

func init() {
	for _, x := range os.Args {
		if x == "-h" {
			showHelp()
			os.Exit(0)
		}
	}
}

func main() {
	repo, err := git.PlainOpen(".")
	if err != nil {
		if err == git.ErrRepositoryNotExists {
			pr("No repository in the current path.")
			os.Exit(1)
		} else {
			pr("Error: %s.", err.Error())
		}
		os.Exit(2)
	}

	head, err := repo.Head()
	if err != nil {
		pr("Error: %s. Have you made any commits?", err.Error())
		os.Exit(2)
	}

	tags, err := GetTags(repo)
	if err != nil {
		pr("Error: %s.", err.Error())
		os.Exit(2)
	}

	var last string
	if len(tags) < 1 {
		last = "v0.0.0"
	} else {
		last = tags[len(tags)-1]
	}

	if !isSemVer(last) {
		pr("%s doesn't conform to a valid semantic version structure.", last)
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

	pr("The new version is %s", version)
}
