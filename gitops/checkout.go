package gitops

import (
	"fmt"

	"github.com/go-git/go-git/v6"
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

	branchRef := plumbing.NewBranchReferenceName(branch)
	if err := wt.Checkout(&git.CheckoutOptions{
		Branch: branchRef,
	}); err != nil {
		return fmt.Errorf("could not checkout branch %q: %w", branch, err)
	}

	return nil
}
