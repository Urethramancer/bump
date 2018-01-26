package main

import (
	"time"

	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

func NewTagger(name, email string) object.Signature {
	return object.Signature{
		Name:  name,
		Email: email,
		When:  time.Now(),
	}
}

func AddAnnotatedTag(repo *git.Repository, tagger object.Signature, text, targethash string) error {
	target := plumbing.NewHash(targethash)
	tag := object.Tag{
		Name:    text,
		Message: text, Tagger: tagger,
		Target:     target,
		TargetType: plumbing.CommitObject,
	}
	e := repo.Storer.NewEncodedObject()
	tag.Encode(e)
	hash, err := repo.Storer.SetEncodedObject(e)
	if err != nil {
		return err
	}

	cref := plumbing.NewReferenceFromStrings("refs/tags/"+text, hash.String())
	err = repo.Storer.SetReference(cref)
	if err != nil {
		return err
	}

	return nil
}
