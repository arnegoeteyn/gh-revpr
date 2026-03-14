package ui

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
)

var (
	success = color.New(color.FgGreen, color.Bold).SprintFunc()
	info    = color.New(color.FgCyan).SprintFunc()
	warn    = color.New(color.FgYellow, color.Bold).SprintFunc()
	error_  = color.New(color.FgRed, color.Bold).SprintFunc()

	fileHeader = color.New(color.FgBlue, color.Bold).SprintFunc()
	location   = color.New(color.FgMagenta).SprintFunc()
)

func Success(format string, args ...any) {
	fmt.Fprintln(os.Stdout, success("✓ "+fmt.Sprintf(format, args...)))
}

func Info(format string, args ...any) {
	fmt.Fprintln(os.Stdout, info("ℹ "+fmt.Sprintf(format, args...)))
}

func Warn(format string, args ...any) {
	fmt.Fprintln(os.Stdout, warn("⚠ "+fmt.Sprintf(format, args...)))
}

func Error(format string, args ...any) {
	fmt.Fprintln(os.Stderr, error_("✗ "+fmt.Sprintf(format, args...)))
}

func Confirm(message string) bool {
	fmt.Fprintf(os.Stdout, "%s %s [y/N]: ", info("?"), message)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	response := scanner.Text()
	return response == "y" || response == "Y"
}

func Ask(message string) string {
	fmt.Fprintf(os.Stdout, "%s %s: ", info("?"), message)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

type Comment struct {
	LineNumber  int
	FilePath    string
	Content     string
	FileContent string
}

func Comments(comments []Comment) {
	if len(comments) == 0 {
		Info("No comments found")
		return
	}

	grouped := make(map[string][]Comment)
	for _, c := range comments {
		grouped[c.FilePath] = append(grouped[c.FilePath], c)
	}

	var files []string
	for f := range grouped {
		files = append(files, f)
	}
	sort.Strings(files)

	for _, file := range files {
		fileComments := grouped[file]
		fmt.Fprintln(os.Stdout)
		fmt.Fprintf(os.Stdout, "%s %s\n", fileHeader("┌─"), fileHeader(file))
		for i, c := range fileComments {
			codeLine := location(fmt.Sprintf("│ %3d:", c.LineNumber))
			emptyLine := location("│     ")
			commentLine := color.New(color.FgYellow).SprintFunc()

			if i > 0 {
				fmt.Fprintln(os.Stdout, location("│"))
			}

			lines := strings.Split(c.Content, "\n")
			fmt.Fprintf(os.Stdout, "%s %s\n", emptyLine, commentLine(lines[0]))
			for _, line := range lines[1:] {
				fmt.Fprintf(os.Stdout, "%s %s\n", emptyLine, commentLine(line))
			}

			if c.FileContent != "" {
				codeLine_ := color.New(color.FgWhite).SprintFunc()
				fmt.Fprintf(os.Stdout, "%s %s\n", codeLine, codeLine_(c.FileContent))
			}
		}
		fmt.Fprintln(os.Stdout, fileHeader("└─"))
	}
}

var spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

type Spinner struct {
	message string
	stop    chan struct{}
	wg      sync.WaitGroup
	mu      sync.Mutex
}

func StartSpinner(message string) *Spinner {
	s := &Spinner{
		message: message,
		stop:    make(chan struct{}),
	}
	s.wg.Add(1)
	go s.run()
	return s
}

func (s *Spinner) Message(format string, args ...any) {
	s.mu.Lock()
	s.message = fmt.Sprintf(format, args...)
	s.mu.Unlock()
}

func (s *Spinner) run() {
	defer s.wg.Done()
	frame := 0
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.mu.Lock()
			msg := s.message
			s.mu.Unlock()
			fmt.Fprintf(os.Stderr, "\r%s %s", spinnerFrames[frame], msg)
			frame = (frame + 1) % len(spinnerFrames)
		case <-s.stop:
			s.mu.Lock()
			msg := s.message
			s.mu.Unlock()
			fmt.Fprintf(os.Stderr, "\r%s %s\n", success("✓"), msg)
			return
		}
	}
}

func (s *Spinner) Stop() {
	close(s.stop)
	s.wg.Wait()
}
