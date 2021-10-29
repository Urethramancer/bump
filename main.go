package main

import (
	"os"

	"github.com/Urethramancer/signor/opt"
	git "github.com/go-git/go-git/v5"
)

var o struct {
	opt.DefaultHelp
	Part string `placeholder:"PART" help:"The part of the version to bump"`
}

func main() {
	a := opt.Parse(&o)
	if o.Help {
		a.Usage()
		return
	}

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

	print(head.Name())
	tagger := NewTagger(user, mail)
	err = AddAnnotatedTag(repo, tagger, version, head.Strings()[1])
	if err != nil {
		pr("Error adding annotated tag: %s", err.Error())
		os.Exit(2)
	}

	pr("The new version is %s", version)
}
