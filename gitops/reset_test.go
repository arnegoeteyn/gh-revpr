package gitops_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/arnegoeteyn/gh-revpr/gitops"
	testgit "github.com/arnegoeteyn/gh-revpr/testutil/git"
	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing/object"
	"github.com/stretchr/testify/require"
)

func TestResetToRemoteBranch_Success(t *testing.T) {
	repoPath := testgit.SetupTestRepoWithRemote(t)
	testgit.Chdir(t, repoPath)

	repo, err := git.PlainOpen(repoPath)
	require.NoError(t, err)

	// Get the current HEAD (should match origin/main)
	remoteRef, err := repo.Reference("refs/remotes/origin/main", true)
	require.NoError(t, err)
	originalRemoteCommit := remoteRef.Hash()

	// Add a local-only commit
	wt, err := repo.Worktree()
	require.NoError(t, err)

	localFile := filepath.Join(repoPath, "local-only.txt")
	err = os.WriteFile(localFile, []byte("local content"), 0644)
	require.NoError(t, err)

	_, err = wt.Add("local-only.txt")
	require.NoError(t, err)

	_, err = wt.Commit("Local commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Test",
			Email: "test@example.com",
			When:  time.Now(),
		},
	})
	require.NoError(t, err)

	// Verify local file exists before reset
	_, err = os.Stat(localFile)
	require.NoError(t, err, "local file should exist before reset")

	// Execute reset
	err = gitops.ResetToRemoteBranch(repo, "main")
	require.NoError(t, err)

	// Verify HEAD now matches the remote commit
	head, err := repo.Head()
	require.NoError(t, err)
	require.Equal(t, originalRemoteCommit, head.Hash(), "HEAD should match origin/main after reset")

	// Verify local-only file is gone
	_, err = os.Stat(localFile)
	require.True(t, os.IsNotExist(err), "local file should be removed after hard reset")
}

func TestResetToRemoteBranch_RemoteRefNotFound(t *testing.T) {
	repoPath := testgit.SetupTestRepo(t)
	testgit.Chdir(t, repoPath)

	repo, err := git.PlainOpen(repoPath)
	require.NoError(t, err)

	err = gitops.ResetToRemoteBranch(repo, "nonexistent")
	require.Error(t, err)
}
