package gitops

import (
	"fmt"
	"log/slog"

	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
)

type Client interface {
	GetPullRequestRef(pr string) (string, error)
}

// ResetToRemoteBranch hard resets the current branch in the worktree to origin/<branch>.
// It updates both the current branch reference and the working tree.
func ResetToRemoteBranch(repo *git.Repository, branch string) error {
	wt, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("could not get worktree: %w", err)
	}

	// Get the remote reference for origin/<branch>
	remoteRef := plumbing.NewRemoteReferenceName("origin", branch)
	remote, err := repo.Reference(remoteRef, true)
	if err != nil {
		return fmt.Errorf("could not get remote reference %q: %w", "origin/"+branch, err)
	}

	// Get the current branch reference (HEAD)
	head, err := repo.Head()
	if err != nil {
		return fmt.Errorf("could not get HEAD: %w", err)
	}

	// Update the current branch reference to point to the remote commit
	currentBranchRef := plumbing.NewHashReference(head.Name(), remote.Hash())
	if err := repo.Storer.SetReference(currentBranchRef); err != nil {
		return fmt.Errorf("could not update branch reference: %w", err)
	}

	// Hard reset the worktree to the remote commit
	if err := wt.Reset(&git.ResetOptions{
		Commit: remote.Hash(),
		Mode:   git.HardReset,
	}); err != nil {
		return fmt.Errorf("could not hard reset to %q: %w", "origin/"+branch, err)
	}

	slog.Debug("hard reset to remote branch", "remoteBranch", "origin/"+branch, "commit", remote.Hash().String())

	return nil
}
