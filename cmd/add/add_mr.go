package add

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// mrCmd represents the mr command
var mrCmd = &cobra.Command{
	Use:   "mr",
	Short: "",
	Run: func(cmd *cobra.Command, args []string) {
		comment, _ := cmd.Flags().GetString("comment")
		link, _ := cmd.Flags().GetString("link")
		tag, _ := cmd.Flags().GetString("tag")

		i := strings.LastIndex(link, "/") + 1
		mr := link[i:]
		decoratedTag := ""
		if len(tag) > 0 {
			decoratedTag = fmt.Sprintf(" #%s", tag)
		}
		updateDailyNoteFile(comment, link, mr, decoratedTag)
	},
}

func updateDailyNoteFile(comment, link, mr, tag string) {
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
	mrOpenedFound := false

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Check if the current line is the "## MRs-Opened" header.
		if strings.Contains(line, "## MRs-Opened") {
			mrOpenedFound = true
		} else if mrOpenedFound && len(strings.TrimSpace(line)) == 0 {
			// Insert your new list item here in the same format.
			// Example: - [/] [new task MR#30](https://link/to/mr30)
			lines = append(lines, fmt.Sprintf("- [ ] [%s MR#%s](%s)%s", comment, mr, link, tag))
			mrOpenedFound = false // Reset the flag after inserting.
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

func GetDailyFileName() string {
	// Get the current date
	currentDate := time.Now()

	// Get the year, month, and day
	year := currentDate.Year()
	month := currentDate.Month()
	day := currentDate.Day()

	// Get the week of the year
	_, week := currentDate.ISOWeek()

	// Get the weekday
	weekday := currentDate.Weekday()

	// Define the suffix for the day
	// suffix := getDaySuffix(day)

	// Format the date as desired
	dateFormat := fmt.Sprintf("%d/%s/%s/%d-%02d-%02d-%s",
		year, month, ordinal(week), year, month, day, weekday)

	return dateFormat
}

// // getDaySuffix returns the ordinal suffix for the day
// func getDaySuffix(day int) string {
// 	switch day {
// 	case 1, 21, 31:
// 		return "st"
// 	case 2, 22:
// 		return "nd"
// 	case 3, 23:
// 		return "rd"
// 	default:
// 		return "th"
// 	}
// }

// ordinal returns the ordinal representation of a number
func ordinal(n int) string {
	switch n {
	case 1:
		return "1st"
	case 2:
		return "2nd"
	case 3:
		return "3rd"
	default:
		return fmt.Sprintf("%dth", n)
	}
}

func init() {
	addCmd.AddCommand(mrCmd)
	mrCmd.Flags().StringP("comment", "m", "", "Provide a short text comment/memo for the MR")
	mrCmd.Flags().StringP("link", "l", "", "Provide a URL/Link for the MR")
	mrCmd.Flags().StringP("tag", "t", "", "Provide a Tag for the MR")
	mrCmd.MarkFlagRequired("comment")
	mrCmd.MarkFlagRequired("link")
}
