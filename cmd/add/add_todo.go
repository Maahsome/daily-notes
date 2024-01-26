package add

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// todoCmd represents the mr command
var todoCmd = &cobra.Command{
	Use:   "todo",
	Short: "",
	Run: func(cmd *cobra.Command, args []string) {
		comment, _ := cmd.Flags().GetString("comment")
		link, _ := cmd.Flags().GetString("link")
		tag, _ := cmd.Flags().GetString("tag")

		decoratedTag := ""
		if len(tag) > 0 {
			decoratedTag = fmt.Sprintf(" #%s", tag)
		}
		updateDailyNoteFileToDo(comment, link, decoratedTag)
	},
}

func updateDailyNoteFileToDo(comment, link, tag string) {
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
	toDoFound := false

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Check if the current line is the "## MRs-Opened" header.
		if strings.Contains(line, "## To-Do") {
			toDoFound = true
		} else if toDoFound && len(strings.TrimSpace(line)) == 0 {
			// Insert your new list item here in the same format.
			// Example: - [/] [new task MR#30](https://link/to/mr30)
			if len(link) > 0 {
				lines = append(lines, fmt.Sprintf("- [ ] [%s](%s)%s", comment, link, tag))
			} else {
				lines = append(lines, fmt.Sprintf("- [ ] %s%s", comment, tag))
			}
			toDoFound = false // Reset the flag after inserting.
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
	addCmd.AddCommand(todoCmd)
	todoCmd.Flags().StringP("comment", "m", "", "Provide a short text comment/memo for the ToDo")
	todoCmd.Flags().StringP("link", "l", "", "Provide a URL/Link for the ToDo")
	todoCmd.Flags().StringP("tag", "t", "", "Provide a Tag for the ToDo")
	todoCmd.MarkFlagRequired("comment")
}
