package main

import (
	"sort"
	"strings"
	"time"

	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

// GetTags returns a sorted list of all tags in the repository.
func GetTags(repo *git.Repository) ([]string, error) {
	ref, err := repo.Tags()
	if err != nil {
		return nil, err
	}

	var tags []string
	err = ref.ForEach(func(p *plumbing.Reference) error {
		t := strings.Replace(p.Name().String(), "refs/tags/", "", 1)
		tags = append(tags, t)
		return nil
	})
	if err != nil {
		return nil, err
	}

	sort.Strings(tags)
	return tags, nil
}

// NewTagger is a barely-useful function which creates the Signature required by some go-git objects.
func NewTagger(name, email string) object.Signature {
	return object.Signature{
		Name:  name,
		Email: email,
		When:  time.Now(),
	}
}

// AddAnnotatedTag is the meat of this program's funtionality.
// repo - an open repository acquired via go-git's Open() or other functions.
// tagger - an author signature created with NewTagger()
// text - this is the actual bare tag, which will have the string "refs/tags/" prepended
// targehash - the hash of the target commit to tag, typically HEAD (acquired with repo.Head())
func AddAnnotatedTag(repo *git.Repository, tagger object.Signature, text, targethash string) error {
	target := plumbing.NewHash(targethash)
	tag := object.Tag{
		Name:    text,
		Message: text, Tagger: tagger,
		Target:     target,
		TargetType: plumbing.CommitObject,
	}
	e := repo.Storer.NewEncodedObject()
	err := tag.Encode(e)
	if err != nil {
		return err
	}

	hash, err := repo.Storer.SetEncodedObject(e)
	if err != nil {
		return err
	}

	cref := plumbing.NewReferenceFromStrings("refs/tags/"+text, hash.String())
	return repo.Storer.SetReference(cref)
}
