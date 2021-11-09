package main

import (
	"os"

	"github.com/Urethramancer/bump/semver"
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
	println(head.String())
	println("")

	tags, err := GetTags(repo)
	if err != nil {
		pr("Error: %s.", err.Error())
		os.Exit(2)
	}

	list := semver.SemVerList{}
	for _, t := range tags {
		list = append(list, semver.New(t))
	}

	user := os.Getenv("GITUSER")
	mail := os.Getenv("GITMAIL")
	last := list.Last()
	pr("Latest version is %s", last)
	pr("Committing as %s <%s>", user, mail)
	if len(os.Args) == 1 {
		return
	}

	switch os.Args[1] {
	case "major", "ma":
		last.Bump(semver.Major)
	case "minor", "mi":
		last.Bump(semver.Minor)
	case "patch", "p":
		last.Bump(semver.Patch)
	default:
		pr("Unknown part: %s", os.Args[1])
		os.Exit(2)
	}

	tagger := NewTagger(user, mail)
	err = AddAnnotatedTag(repo, tagger, last.String(), head.Strings()[1])
	if err != nil {
		pr("Error adding annotated tag: %s", err.Error())
		os.Exit(2)
	}

	pr("The new version is %s", last.String())
}
