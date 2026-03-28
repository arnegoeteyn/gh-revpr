package gitops

import (
	"fmt"
	"log/slog"

	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
)

func ResetToRemoteBranch(repo *git.Repository, branch string) error {
	wt, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("could not get worktree: %w", err)
	}

	remoteRef := plumbing.NewRemoteReferenceName("origin", branch)
	remote, err := repo.Reference(remoteRef, true)
	if err != nil {
		return fmt.Errorf("could not get remote reference %q: %w", remoteRef.String(), err)
	}

	slog.Debug("hard reset to remote branch", "remoteBranch", remoteRef.String(), "commit", remote.Hash().String())

	if err := wt.Reset(&git.ResetOptions{
		Commit: remote.Hash(),
		Mode:   git.HardReset,
	}); err != nil {
		return fmt.Errorf("could not hard reset to %q: %w", remoteRef.String(), err)
	}

	return nil
}
