/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log/slog"
	"os"

	"github.com/arnegoeteyn/gh-revpr/github"
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

		client, err := github.NewClient()
		if err != nil {
			ui.Error("could not connect to github: %s", err.Error())
			slog.Error("could not connect to github", "error", err)
			os.Exit(1)
		}

		event := ui.Select("PR review event", []ui.SelectOption{
			{Value: string(github.ReviewEventApprove), Description: "Approve the PR"},
			{Value: string(github.ReviewEventComment), Description: "Leave a comment without approval"},
			{Value: string(github.ReviewEventRequestChanges), Description: "Request changes to the PR"},
		})
		if event == "" {
			ui.Warn("review creation cancelled")
			return
		}

		ui.Info("submitting current PR with event %q", event)

		var githubComments []github.Comment
		for _, c := range comments {
			githubComments = append(githubComments, github.Comment{
				Line: c.LineNumber,
				Path: c.FilePath,
				Body: c.Content,
			})
		}

		body := ui.Ask("PR review body")

		if err := client.Review(currentPR, github.Review{
			Event:    github.ReviewEvent(event),
			Comments: githubComments,
			Body:     body,
		}); err != nil {
			ui.Error("could not create review")
			slog.Error("could not create review", "error", err)
			os.Exit(1)
		}

		ui.Success("Approved PR with %d comments", len(comments))
	},
}

func init() {
	rootCmd.AddCommand(commentCmd)
}
