package gitops

import (
	"fmt"
	"log/slog"

	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/config"
	"github.com/go-git/go-git/v6/plumbing"
)

type Client interface {
	GetPullRequestRef(pr string) (string, error)
}

// Checkout checks out the branch associated with the given PR number.
func Checkout(repo *git.Repository, branch string) error {
	// Get the worktree and checkout the branch
	wt, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("could not get worktree: %w", err)
	}

	if err := createIfNotExists(repo, branch); err != nil {
		return err
	}

	branchRef := plumbing.NewBranchReferenceName(branch)
	if err := wt.Checkout(&git.CheckoutOptions{
		Branch: branchRef,
	}); err != nil {
		return fmt.Errorf("could not checkout branch %q: %w", branch, err)
	}

	return nil
}

func createIfNotExists(repo *git.Repository, branch string) error {
	branchRef := plumbing.NewBranchReferenceName(branch)

	_, err := repo.Reference(branchRef, false)
	if err == nil {
		slog.Debug("branch already exists, not creating it", "branch", branch)
		return nil
	}

	slog.Debug("branch does not exist, creating it", "branch", branch)

	remoteRef := plumbing.NewRemoteReferenceName("origin", branch)
	remote, err := repo.Reference(remoteRef, true)
	if err != nil {
		return fmt.Errorf("remote branch %q does not exist: %w", "origin/"+branch, err)
	}

	ref := plumbing.NewHashReference(branchRef, remote.Hash())
	if err := repo.Storer.SetReference(ref); err != nil {
		return fmt.Errorf("could not create branch reference: %w", err)
	}

	if err := repo.CreateBranch(&config.Branch{
		Name:   branch,
		Remote: "origin",
		Merge:  branchRef,
	}); err != nil {
		return fmt.Errorf("could not create branch config: %w", err)
	}
	return nil
}
