package gitops_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/arnegoeteyn/gh-revpr/gitops"
	"github.com/arnegoeteyn/gh-revpr/testutil/git"
	"github.com/stretchr/testify/require"
)

const wtPath = ".revpr/wt"

func TestCreateWorktree_FromRoot(t *testing.T) {
	repoPath := git.SetupTestRepo(t)
	git.Chdir(t, repoPath)

	_, err := gitops.CreateReviewWorktree()
	require.NoError(t, err)

	wtDir := filepath.Join(repoPath, wtPath)
	entries, err := os.ReadDir(wtDir)
	require.NoError(t, err, "worktree directory should exist")
	require.Len(t, entries, 1, "should have created one worktree")
}

func TestCreateWorktree_FromSubdirectory(t *testing.T) {
	repoPath := git.SetupTestRepo(t)

	// Create a subdirectory and change into it
	subDir := filepath.Join(repoPath, "some", "nested", "directory")
	err := os.MkdirAll(subDir, 0755)
	require.NoError(t, err, "failed to create subdirectory")

	git.Chdir(t, subDir)

	_, err = gitops.CreateReviewWorktree()
	require.NoError(t, err)

	// Verify worktree was created in .git directory, not in subdirectory
	wtDir := filepath.Join(repoPath, wtPath)
	entries, err := os.ReadDir(wtDir)
	require.NoError(t, err, "worktree directory should exist in .git")
	require.Len(t, entries, 1, "should have created one worktree")

	// Verify no worktree was created in the subdirectory
	_, err = os.Stat(filepath.Join(subDir, wtDir))
	require.True(t, os.IsNotExist(err), "worktree should not be created in subdirectory")
}

func TestCreateWorktree_NotAGitRepo(t *testing.T) {
	// Create a temp directory that is NOT a git repo
	tmpDir := t.TempDir()
	git.Chdir(t, tmpDir)

	_, err := gitops.CreateReviewWorktree()
	require.Error(t, err)
}

func TestCreateWorktree_MultipleWorktrees(t *testing.T) {
	repoPath := git.SetupTestRepo(t)
	git.Chdir(t, repoPath)

	// Create multiple worktrees
	_, err := gitops.CreateReviewWorktree()
	require.NoError(t, err)

	_, err = gitops.CreateReviewWorktree()
	require.NoError(t, err)

	// Verify both worktrees were created
	wtDir := filepath.Join(repoPath, wtPath)
	entries, err := os.ReadDir(wtDir)
	require.NoError(t, err)
	require.Len(t, entries, 2, "should have created two worktrees")

	// Verify they have different names
	require.NotEqual(t, entries[0].Name(), entries[1].Name(),
		"worktrees should have unique names")
}

func TestCreateWorktree_FromBareRepoWorktree(t *testing.T) {
	repoRoot, worktreePath := git.SetupBareRepoWithWorktrees(t)
	git.Chdir(t, worktreePath)

	_, err := gitops.CreateReviewWorktree()
	require.NoError(t, err)

	// Verify worktree was created at the repo root, not in the existing worktree
	wtDir := filepath.Join(repoRoot, wtPath)
	entries, err := os.ReadDir(wtDir)
	require.NoError(t, err, "worktree directory should exist at repo root")
	require.Len(t, entries, 1, "should have created one worktree")
}
