package add

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// commentCmd represents the mr command
var commentCmd = &cobra.Command{
	Use:     "comment",
	Aliases: []string{"note"},
	Short:   "",
	Run: func(cmd *cobra.Command, args []string) {
		comment, _ := cmd.Flags().GetString("comment")
		link, _ := cmd.Flags().GetString("link")
		tag, _ := cmd.Flags().GetString("tag")

		decoratedTag := ""
		if len(tag) > 0 {
			decoratedTag = fmt.Sprintf(" #%s", tag)
		}
		updateDailyNoteFileComment(comment, link, decoratedTag)
	},
}

func updateDailyNoteFileComment(comment, link, tag string) {
	// Open the markdown file.
	filename := fmt.Sprintf("/Users/christopher.maahs/src/christopher-maahs/Daily_Notes/%s.md", GetDailyFileName())
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Create a temporary slice to hold the lines.
	var lines []string
	commentsFound := false

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Check if the current line is the "## MRs-Opened" header.
		if strings.Contains(line, "## What Did You Do Today?") {
			commentsFound = true
		} else if commentsFound && len(strings.TrimSpace(line)) == 0 {
			// Insert your new list item here in the same format.
			// Example: - [/] [new task MR#30](https://link/to/mr30)
			if len(link) > 0 {
				lines = append(lines, fmt.Sprintf("* [%s](%s)%s", comment, link, tag))
			} else {
				lines = append(lines, fmt.Sprintf("* %s%s", comment, tag))
			}
			commentsFound = false // Reset the flag after inserting.
		}
		// Add the line to the slice.
		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// Open the file again for writing.
	file, err = os.Create(filename)
	if err != nil {
		fmt.Println("Error opening file for writing:", err)
		return
	}
	defer file.Close()

	// Write the updated lines back to the file.
	writer := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(writer, line)
	}
	writer.Flush()
}

func init() {
	addCmd.AddCommand(commentCmd)
	commentCmd.Flags().StringP("comment", "m", "", "Provide a short text comment/memo for the Comment")
	commentCmd.Flags().StringP("link", "l", "", "Provide a URL/Link for the Comment")
	commentCmd.Flags().StringP("tag", "t", "", "Provide a Tag for the Comment")
	commentCmd.MarkFlagRequired("comment")
}
