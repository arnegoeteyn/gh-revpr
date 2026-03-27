/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log/slog"
	"os"

	"github.com/arnegoeteyn/gh-revpr/github"
	"github.com/arnegoeteyn/gh-revpr/gitops"
	"github.com/arnegoeteyn/gh-revpr/state"
	"github.com/arnegoeteyn/gh-revpr/ui"
	"github.com/go-git/go-git/v6"
	"github.com/spf13/cobra"
)

const (
	flagCreateWorktree = "create-worktree"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init <PR>",
	Short: "Initialize a PR review.",
	Long: `Start with reviewing a PR. Init will create a new worktree that can be used for reviewing.

By default it will start an interactive rebase and do a mixed reset of the first commit.
Later commits can be reviewed by running 'revpr continue'.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pr := args[0]

		ui.Info("Starting init for PR %s", pr)

		var repo *git.Repository
		createWorktree, err := cmd.Flags().GetBool(flagCreateWorktree)
		if err != nil {
			slog.Debug("failed to get flag", "flag", flagCreateWorktree, "err", err)
			createWorktree = false
		}

		if !createWorktree {
			slog.Debug("using review worktree")
			existingRepo, err := gitops.GetReviewWorktree()
			handleErr(err, "failed to use review worktree")
			repo = existingRepo
		} else {
			slog.Debug("creating worktree")
			newRepo, err := gitops.CreateReviewWorktree()
			handleErr(err, "failed to create worktree")
			repo = newRepo
		}

		gh, err := github.NewClient()
		handleErr(err, "failed to create GitHub client")

		pullRequest, err := gh.GetPullRequest(pr)
		handleErr(err, "failed to get PR branch")

		branch := pullRequest.Head.Ref
		slog.Debug("Got PR branch", "branch", branch)

		err = gitops.ResetToRemoteBranch(repo, branch)
		handleErr(err, "failed to reset to branch")
		slog.Debug("Reset to branch", "branch", branch)

		ui.Success("Created worktree for review")

		wt, err := repo.Worktree()
		handleErr(err, "failed to get worktree")
		ui.Info("Worktree at: %s", wt.Filesystem.Root())

		ui.PRBody(pullRequest.Body)

		slog.Debug("storing revpr state", "currentPR", pr)
		err = os.Chdir(wt.Filesystem.Root())
		handleErr(err, "could not navigate to new worktree")

		err = state.SetCurrentPR(pr)
		handleErr(err, "could not store config file")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().BoolP("create-worktree", "c", false, "Generate a new worktree to review in. If this option is not set a default `review` worktree will be expected to exist and used.")
}
