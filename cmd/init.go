/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/arnegoeteyn/gh-revpr/github"
	"github.com/arnegoeteyn/gh-revpr/gitops"
	"github.com/arnegoeteyn/gh-revpr/ui"
	"github.com/spf13/cobra"
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

		repo, err := gitops.CreateReviewWorktree()
		if err != nil {
			ui.Error("Failed to create worktree: %v", err)
			Debug("Failed to create worktree", "error", err)
			os.Exit(1)
		}
		Debug("Created worktree")

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
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
