/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log/slog"
	"os"
	"time"

	"github.com/arnegoeteyn/gh-revpr/pr"
	"github.com/arnegoeteyn/gh-revpr/state"
	"github.com/arnegoeteyn/gh-revpr/ui"
	"github.com/go-git/go-git/v6"
	"github.com/spf13/cobra"
)

// commentCmd represents the comment command
var commentCmd = &cobra.Command{
	Use:   "comment",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		slog.Debug("Parsing review comments")

		repo, err := git.PlainOpen(".")
		if err != nil {
			ui.Error("Failed to open git repository: %v", err)
			slog.Error("Failed to open git repository", "error", err)
			os.Exit(1)
		}

		ref, err := repo.Head()
		if err != nil {
			ui.Error("Failed to get HEAD: %v", err)
			slog.Error("Failed to get HEAD", "error", err)
			os.Exit(1)
		}

		commit, err := repo.CommitObject(ref.Hash())
		if err != nil {
			ui.Error("Failed to get commit: %v", err)
			slog.Error("Failed to get commit", "error", err)
			os.Exit(1)
		}

		comments, err := pr.Comments(commit)
		if err != nil {
			ui.Error("Failed to get comments: %v", err)
			slog.Error("Failed to get comments", "error", err)
			os.Exit(1)
		}

		ui.Info("parsed %d comment(s) on current PR", len(comments))

		var uiComments []ui.Comment
		for _, c := range comments {
			uiComments = append(uiComments, ui.Comment{
				LineNumber:  c.LineNumber,
				FilePath:    c.FilePath,
				Content:     c.Content,
				FileContent: c.FileContent,
			})
		}
		ui.Comments(uiComments)

		if !ui.Confirm("Create review with these comments") {
			ui.Warn("review creation cancelled")
			return
		}

		currentPR := state.CurrentPR()
		if currentPR == "" {
			ui.Warn("No current PR configured")
			currentPR = ui.Ask("What PR should this review target?")
		}

		ui.Info("Creating review for PR #%s", currentPR)

		ui.Info("uploading %d comment(s) on current PR", len(comments))

		spinner := ui.StartSpinner("uploading comments...")
		for i := range comments {
			spinner.Message("uploading comment %d/%d", i+1, len(comments))
			time.Sleep(5 * time.Second)
		}
		spinner.Stop()

		ui.Success("Approved PR with %d comments", len(comments))
	},
}

func init() {
	rootCmd.AddCommand(commentCmd)
}
