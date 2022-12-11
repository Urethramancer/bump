package main

import (
	"fmt"
	"os"

	"github.com/Urethramancer/bump/semver"
	git "github.com/go-git/go-git/v5"
	"github.com/grimdork/climate/arg"
)

var o struct {
	// opt.DefaultHelp
	Part string `placeholder:"PART" help:"The part of the version to bump" choices:"major,minor,patch,ma,mi,p"`
}

func main() {
	opt := arg.New("bump")
	opt.SetDefaultHelp(true)
	opt.SetPositional("PART", "The part of the version to bump", "", true, arg.VarString)
	var err error
	err = opt.Parse(os.Args)
	if err != nil {
		if err == arg.ErrNoArgs {
			opt.PrintHelp()
			return
		}

		if err == arg.ErrRunCommand {
			return
		}

		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(2)
	}

	part := opt.GetPosString("PART")
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

	switch part {
	case "major", "ma":
		last.Bump(semver.Major)
	case "minor", "mi":
		last.Bump(semver.Minor)
	case "patch", "p":
		last.Bump(semver.Patch)
	default:
		pr("Unknown part: %s", o.Part)
	}

	tagger := NewTagger(user, mail)
	err = AddAnnotatedTag(repo, tagger, last.String(), head.Strings()[1])
	if err != nil {
		pr("Error adding annotated tag: %s", err.Error())
		os.Exit(2)
	}

	pr("The new version is %s", last.String())
}
