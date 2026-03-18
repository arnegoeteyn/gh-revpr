package pr

import (
	"fmt"
	"log/slog"
	"regexp"
	"strings"

	"github.com/go-git/go-git/v6/plumbing/object"
)

type Comment struct {
	LineNumber  int
	FilePath    string
	Content     string
	FileContent string
}

var (
	startCommentRegex = regexp.MustCompile(`^\s*//\s*PR:\s*(?<PComment>.*)`)
	commentRegex      = regexp.MustCompile(`^\s*//\s?\s(?<PComment>.*)`)
)

func Comments(commit *object.Commit) ([]Comment, error) {
	patch, err := patchForCommit(commit)
	if err != nil {
		return nil, err
	}

	slog.Debug("Starting to extract comments from commit patch", "commit", commit.Hash.String(), "files", len(patch.FilePatches()))

	var allComments []Comment

	for _, fp := range patch.FilePatches() {
		from, to := fp.Files()

		if to == nil {
			slog.Debug("skipping deleted file", "file", from.Path())
			continue
		}

		slog.Debug("Processing file", "file", to.Path())

		if fp.IsBinary() {
			slog.Warn("skipping binary file", "path", to.Path())
			continue
		}

		f, err := commit.File(to.Path())
		if err != nil {
			panic(err)
		}

		comments, err := commentsFromFile(f)
		if err != nil {
			slog.Error("Failed to extract comments from file", "file", f.Name, "error", err)
			return nil, err
		}

		slog.Debug("Extracted comments from file", "file", f.Name, "count", len(comments))
		allComments = append(allComments, comments...)
	}

	slog.Debug("Finished extracting comments", "total", len(allComments), "files", len(patch.FilePatches()))
	return allComments, nil
}

func patchForCommit(commit *object.Commit) (*object.Patch, error) {
	parent, err := commit.Parent(0)
	if err != nil {
		return nil, fmt.Errorf("could not get parent of commit: %w", err)
	}

	patch, err := parent.Patch(commit)
	if err != nil {
		return nil, fmt.Errorf("could not get patch of commit: %w", err)
	}

	return patch, nil
}

func commentsFromFile(file *object.File) ([]Comment, error) {
	slog.Debug("Parsing file for comments", "file", file.Name)

	var comments []Comment

	lines, err := file.Lines()
	if err != nil {
		slog.Error("Failed to get lines from file", "file", file.Name, "error", err)
		return nil, fmt.Errorf("could not get lines from file: %w", err)
	}

	slog.Debug("Processing lines", "file", file.Name, "lineCount", len(lines))

	var currentComment strings.Builder
	for i, line := range lines {
		if currentComment.Len() != 0 {
			matches := commentRegex.FindStringSubmatch(line)
			if len(matches) > 1 {
				_, err := currentComment.WriteString("\n" + matches[1])
				if err != nil {
					slog.Error("Failed to write to comment builder", "file", file.Name, "line", i, "error", err)
					return nil, fmt.Errorf("could not write to comment builder: %w", err)
				}
			} else {
				comment := Comment{
					LineNumber:  i,
					Content:     currentComment.String(),
					FilePath:    file.Name,
					FileContent: line,
				}
				slog.Debug("Found comment end", "file", file.Name, "line", i, "content", comment.Content)
				comments = append(comments, comment)
				currentComment.Reset()
			}
		}

		matches := startCommentRegex.FindStringSubmatch(line)
		if len(matches) > 1 {
			slog.Debug("Found comment start", "file", file.Name, "line", i, "content", matches[1])
			_, err := currentComment.WriteString(matches[1])
			if err != nil {
				slog.Error("Failed to write to comment builder", "file", file.Name, "line", i, "error", err)
				return nil, fmt.Errorf("could not write to comment builder: %w", err)
			}
		}
	}

	slog.Debug("Finished parsing file", "file", file.Name, "commentsFound", len(comments))
	return comments, nil
}
