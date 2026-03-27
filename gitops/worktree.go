package gitops

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/arnegoeteyn/gh-revpr/faker"
	"github.com/go-git/go-billy/v6/osfs"
	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/storage/filesystem"

	xworktree "github.com/go-git/go-git/v6/x/plumbing/worktree"
)

const worktreeSubdir = ".revpr/wt"

// CreateWorktree creates a new git worktree for reviewing.
// It searches upward from the current directory to find the git repository,
// then creates the worktree relative to the git directory.
// For a regular repo: <repo-root>/.revpr/wt/<random-name>
// For a bare repo: <bare-repo>/revpr/wt/<random-name>
func CreateReviewWorktree() (*git.Repository, error) {
	wt, err := addWorktree(faker.ReviewWorktreeName())
	if err != nil {
		return nil, err
	}

	return wt, nil
}

// GetReviewWorktree returns an existing worktree named "review".
// It searches upward from the current directory to find the git repository,
// then opens the worktree at the expected location.
func GetReviewWorktree() (*git.Repository, error) {
	return openWorktree("review")
}

// worktreeContext holds the common setup needed for worktree operations.
type worktreeContext struct {
	manager   *xworktree.Worktree
	commonDir string
}

func newWorktreeContext() (*worktreeContext, error) {
	repo, err := git.PlainOpenWithOptions(".", &git.PlainOpenOptions{
		DetectDotGit: true,
	})
	if err != nil {
		return nil, fmt.Errorf("could not open git repo: %w", err)
	}

	// Get the git directory path from the storage filesystem
	storage, ok := repo.Storer.(*filesystem.Storage)
	if !ok {
		return nil, fmt.Errorf("unsupported storage type")
	}

	gitDir := storage.Filesystem().Root()

	// Check if we're in a worktree by looking for a commondir file.
	// If it exists, resolve it to get the main git directory.
	// This is equivalent to `git rev-parse --git-common-dir`.
	commonDir := resolveCommonDir(gitDir)

	// If we're in a worktree (gitDir != commonDir), we need to open the main
	// repository to get a proper storer for the worktree manager.
	mainRepo := repo
	if gitDir != commonDir {
		// Open the main repository using the common directory's parent
		// (which is the repo root for bare repos, or the .git parent for regular repos)
		mainRepoPath := filepath.Join(commonDir, "..")
		mainRepo, err = git.PlainOpen(mainRepoPath)
		if err != nil {
			return nil, fmt.Errorf("could not open main repo: %w", err)
		}
	}

	manager, err := xworktree.New(mainRepo.Storer)
	if err != nil {
		return nil, fmt.Errorf("could not create worktree manager: %w", err)
	}

	return &worktreeContext{
		manager:   manager,
		commonDir: commonDir,
	}, nil
}

func (c *worktreeContext) worktreePath(name string) string {
	return filepath.Join(c.commonDir, "..", worktreeSubdir, name)
}

func (c *worktreeContext) worktreePathWithoutSubdir(name string) string {
	return filepath.Join(c.commonDir, "..", name)
}

func openWorktree(path string) (*git.Repository, error) {
	ctx, err := newWorktreeContext()
	if err != nil {
		return nil, err
	}

	worktreeDir := osfs.New(ctx.worktreePathWithoutSubdir(path))

	repository, err := ctx.manager.Open(worktreeDir)
	if err != nil {
		return nil, fmt.Errorf("could not open worktree: %w", err)
	}

	return repository, nil
}

func addWorktree(name string) (*git.Repository, error) {
	ctx, err := newWorktreeContext()
	if err != nil {
		return nil, err
	}

	worktreeDir := osfs.New(ctx.worktreePath(name))

	if err := ctx.manager.Add(worktreeDir, name); err != nil {
		return nil, fmt.Errorf("could not create worktree: %w", err)
	}

	repository, err := ctx.manager.Open(worktreeDir)
	if err != nil {
		return nil, fmt.Errorf("could not open worktree: %w", err)
	}

	return repository, nil
}

// resolveCommonDir returns the common git directory.
// If gitDir contains a "commondir" file (indicating we're in a worktree),
// it resolves the path to the main git directory.
// Otherwise, it returns gitDir unchanged.
func resolveCommonDir(gitDir string) string {
	commondirPath := filepath.Join(gitDir, "commondir")
	content, err := os.ReadFile(commondirPath)
	if err != nil {
		// No commondir file, we're in the main git directory
		return gitDir
	}

	// commondir contains a relative path like "../.."
	relativePath := strings.TrimSpace(string(content))
	return filepath.Join(gitDir, relativePath)
}
