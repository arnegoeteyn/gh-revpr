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
			Debug("failed to get flag", "flag", flagCreateWorktree, "err", err)
			createWorktree = false
		}

		if !createWorktree {
			Debug("using review worktree")
			panic("not yet implemented")
		} else {
			Debug("creating worktree")
			newRepo, err := gitops.CreateReviewWorktree()
			if err != nil {
				ui.Error("Failed to create worktree: %v", err)
				Debug("Failed to create worktree", "error", err)
				os.Exit(1)
			}

			repo = newRepo
		}

		gh, err := github.NewClient()
		if err != nil {
			ui.Error("Failed to create GitHub client: %v", err)
			Debug("Failed to create GitHub client", "error", err)
			os.Exit(1)
		}

		pullRequest, err := gh.GetPullRequest(pr)
		if err != nil {
			ui.Error("Failed to get PR branch: %v", err)
			Debug("Failed to get PR branch", "error", err)
			os.Exit(1)
		}

		branch := pullRequest.Head.Ref
		Debug("Got PR branch", "branch", branch)

		if err := gitops.Checkout(repo, branch); err != nil {
			ui.Error("Failed to checkout branch: %v", err)
			Debug("Failed to checkout branch", "error", err)
			os.Exit(1)
		}
		Debug("Checked out branch", "branch", branch)

		ui.Success("Created worktree for review")

		wt, err := repo.Worktree()
		if err != nil {
			ui.Error("Failed to get worktree: %v", err)
			Debug("Failed to get worktree", "error", err)
			os.Exit(1)
		}
		ui.Info("Worktree at: %s", wt.Filesystem.Root())

		ui.PRBody(pullRequest.Body)

		slog.Debug("storing revpr state", "currentPR", pr)
		if err := os.Chdir(wt.Filesystem.Root()); err != nil {
			slog.Error("could not navigate to new worktree", "path", wt.Filesystem.Root())
			ui.Error("could not navigate to new worktree")
		}

		if err := state.SetCurrentPR(pr); err != nil {
			slog.Error("could not store config file", "error", err.Error())
			ui.Error("could not store config file")
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().BoolP("create-worktree", "c", false, "Generate a new worktree to review in. If this option is not set a default `review` worktree will be expected to exist and used.")
}
