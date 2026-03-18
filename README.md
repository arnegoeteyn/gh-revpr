# gh-revpr

## Pre-requisites
Since this is a gh-cli extension, `gh` needs to be installed. 

## Installation
To install this script as a gh cli extension, clone this repo and run the following command.
```bash
go build . && gh extension install .
``` 

If there are ever any changes, simply run `go build .` in your clone directory again.

## Usage

In order to start a review, run the following command:
```bash
gh revpr init <PRNUMBER>
```

This will create a worktree where you can do your review.

Creating a review goes as follows:
1. go to the new worktree.
2. Aggregate all diffs from the pr target branch into a single commit, writing review comments as you go. You write a review comment by adding a comment starting with `PR:`.
3. run `gh revpr comment` to aggregate these comments and send the review to github.
