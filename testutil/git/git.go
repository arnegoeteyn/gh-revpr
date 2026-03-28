package git

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-git/go-billy/v6/osfs"
	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/config"
	"github.com/go-git/go-git/v6/plumbing"
	"github.com/go-git/go-git/v6/plumbing/object"
	xworktree "github.com/go-git/go-git/v6/x/plumbing/worktree"
	"github.com/stretchr/testify/require"
)

func SetupTestRepo(t *testing.T) string {
	t.Helper()

	tmpDir := t.TempDir()

	repo, err := git.PlainInit(tmpDir, false)
	require.NoError(t, err, "failed to init git repo")

	wt, err := repo.Worktree()
	require.NoError(t, err, "failed to get worktree")

	dummyFile := filepath.Join(tmpDir, "README.md")
	err = os.WriteFile(dummyFile, []byte("# Test Repo\n"), 0644)
	require.NoError(t, err, "failed to create dummy file")

	_, err = wt.Add("README.md")
	require.NoError(t, err, "failed to add file")

	_, err = wt.Commit("Initial commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Test",
			Email: "test@example.com",
			When:  time.Now(),
		},
	})
	require.NoError(t, err, "failed to create initial commit")

	return tmpDir
}

func Chdir(t *testing.T, dir string) {
	t.Helper()

	originalDir, err := os.Getwd()
	require.NoError(t, err, "failed to get current directory")

	err = os.Chdir(dir)
	require.NoError(t, err, "failed to change directory")

	t.Cleanup(func() {
		_ = os.Chdir(originalDir)
	})
}

// SetupTestRepoWithRemote creates a local repo with an "origin" remote.
// It returns the local repo path. The remote is a bare repo in a separate temp directory.
// The local repo has one commit and tracks origin/main, so refs/remotes/origin/main exists.
func SetupTestRepoWithRemote(t *testing.T) string {
	t.Helper()

	mainBranch := plumbing.NewBranchReferenceName("main")

	// Create bare "remote" repo with an initial commit
	remoteDir := t.TempDir()
	remoteRepo, err := git.PlainInit(remoteDir, true, git.WithDefaultBranch(mainBranch))
	require.NoError(t, err, "failed to init bare remote repo")

	// We need to create a commit in the bare repo. The easiest way is to
	// create a temp non-bare repo, make a commit, and push to the bare repo.
	tempDir := t.TempDir()
	tempRepo, err := git.PlainInit(tempDir, false, git.WithDefaultBranch(mainBranch))
	require.NoError(t, err, "failed to init temp repo")

	_, err = tempRepo.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{remoteDir},
	})
	require.NoError(t, err, "failed to add origin to temp repo")

	wt, err := tempRepo.Worktree()
	require.NoError(t, err, "failed to get temp worktree")

	dummyFile := filepath.Join(tempDir, "README.md")
	err = os.WriteFile(dummyFile, []byte("# Test Repo\n"), 0644)
	require.NoError(t, err, "failed to create dummy file")

	_, err = wt.Add("README.md")
	require.NoError(t, err, "failed to add file")

	_, err = wt.Commit("Initial commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Test",
			Email: "test@example.com",
			When:  time.Now(),
		},
	})
	require.NoError(t, err, "failed to create initial commit")

	err = tempRepo.Push(&git.PushOptions{})
	require.NoError(t, err, "failed to push to bare remote")

	// Set HEAD on bare repo to point to main branch
	head, err := tempRepo.Head()
	require.NoError(t, err, "failed to get temp repo HEAD")
	err = remoteRepo.Storer.SetReference(head)
	require.NoError(t, err, "failed to set HEAD on bare repo")

	// Now clone the bare repo to create the actual local repo with tracking
	localDir := t.TempDir()
	_, err = git.PlainClone(localDir, &git.CloneOptions{
		URL: remoteDir,
	})
	require.NoError(t, err, "failed to clone from remote")

	return localDir
}

func SetupBareRepoWithWorktrees(t *testing.T) (repoRoot, mainWorktreePath string) {
	t.Helper()

	tmpDir := t.TempDir()

	barePath := filepath.Join(tmpDir, ".bare")
	_, err := git.PlainInit(barePath, true)
	require.NoError(t, err, "failed to init bare repo")

	dotGitPath := filepath.Join(tmpDir, ".git")
	err = os.WriteFile(dotGitPath, []byte("gitdir: ./.bare"), 0644)
	require.NoError(t, err, "failed to write .git file")

	tempClonePath := filepath.Join(tmpDir, "temp")
	tempRepo, err := git.PlainInit(tempClonePath, false)
	require.NoError(t, err, "failed to init temp repo")

	_, err = tempRepo.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{barePath},
	})
	require.NoError(t, err, "failed to add origin remote")

	tempWt, err := tempRepo.Worktree()
	require.NoError(t, err, "failed to get temp worktree")

	_, err = tempWt.Commit("Initial commit", &git.CommitOptions{
		AllowEmptyCommits: true,
		Author: &object.Signature{
			Name:  "Test",
			Email: "test@example.com",
			When:  time.Now(),
		},
	})
	require.NoError(t, err, "failed to create initial commit")

	err = tempRepo.Push(&git.PushOptions{})
	require.NoError(t, err, "failed to push to bare repo")

	err = os.RemoveAll(tempClonePath)
	require.NoError(t, err, "failed to remove temp clone")

	repo, err := git.PlainOpen(tmpDir)
	require.NoError(t, err, "failed to open repo via .git file")

	worktreeManager, err := xworktree.New(repo.Storer)
	require.NoError(t, err, "failed to create worktree manager")

	mainPath := filepath.Join(tmpDir, "main")
	mainFs := osfs.New(mainPath)
	err = worktreeManager.Add(mainFs, "main")
	require.NoError(t, err, "failed to create main worktree")

	devPath := filepath.Join(tmpDir, "development")
	devFs := osfs.New(devPath)
	err = worktreeManager.Add(devFs, "development")
	require.NoError(t, err, "failed to create development worktree")

	return tmpDir, mainPath
}
