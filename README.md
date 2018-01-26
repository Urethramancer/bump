# Bump
Bumps the semantic version number for the git repository in the current path.

## What
Just a simple little tool made while figuring out how to use src-d's go-git package.

## How
To bump the major version:

```sh
bump major
```

or

```sh
bump ma
```

The other keywords are `minor` and `patch`, with the shortcuts `mi` and `p`.

Running it without any arguments will show the latest detected version and author variables.

## About you
Due to not having figured out a way to get git globals via go-git, the program requires two environment variables: `GITUSER` and `GITMAIL`. Set these to your name and preferred git e-mail to have tags annotated properly with you as the author.
