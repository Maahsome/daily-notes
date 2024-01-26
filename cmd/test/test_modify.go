package test

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"

	// "github.com/gomarkdown/markdown"
	// "github.com/gomarkdown/markdown/ast"
	// "github.com/gomarkdown/markdown/html"
	// mdrenderer "github.com/gomarkdown/markdown/md"
	// "github.com/gomarkdown/markdown/parser"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

// var mds = `[link](http://example.com)`

var mds = `
## What Did You Do Today?
* OpenZiti / Cloudflare Warp Discussions

## To-Do
- [/] Learn OpenZiti [confluence](https://alteryx.atlassian.net/wiki/spaces/CE/pages/2033516939/OpenZiti)

## MRs-Opened
- [x] [enable cdn cache in sandbox on /auto-insights MR#25](https://git.alteryx.com/futurama/bender/cloudflare/provisioning/lowers/-/merge_requests/25/diffs)
- [/] [enable cdn cache in sandbox on /other-stuff MR#29](https://git.alteryx.com/futurama/bender/cloudflare/provisioning/lowers/-/merge_requests/29/diffs)

<< [[2024-01-24-Wednesday|Yesterday]] | [[Daily_Notes/2024/January/2024-01-26-Friday|Tomorrow]] >>

`

// modifyCmd represents the modify command
var modifyCmd = &cobra.Command{
	Use:   "modify",
	Short: "",
	Run: func(cmd *cobra.Command, args []string) {
		_, err := doThis()
		if err != nil {
			logrus.WithError(err).Error("Failed to express the version")
		}
		// if !c.FormatOverridden {
		// 	c.OutputFormat = "text"
		// }
		// c.OutputData(&out)
	},
}

func doThis() (string, error) {
	directTextManipulation()
	return "", nil
}

// func modifyAst(doc ast.Node) ast.Node {
// 	currentHeader := ""
// 	ast.WalkFunc(doc, func(node ast.Node, entering bool) ast.WalkStatus {
// 		if img, ok := node.(*ast.Image); ok && entering {
// 			attr := img.Attribute
// 			if attr == nil {
// 				attr = &ast.Attribute{}
// 			}
// 			// TODO: might be duplicate
// 			attr.Classes = append(attr.Classes, []byte("blog-img"))
// 			img.Attribute = attr
// 		}

// 		if heading, ok := node.(*ast.Heading); ok && entering {
// 			fmt.Printf("heading level: %d = %s\n", heading.Level, heading.Children[0].AsLeaf().Literal)
// 			currentHeader = string(heading.Children[0].AsLeaf().Literal)
// 		}
// 		if list, ok := node.(*ast.ListItem); ok && entering {
// 			if currentHeader == "MRs-Opened" {
// 				listText := list.Children[0].(*ast.Paragraph).Children // [0].AsLeaf().Literal
// 				fmt.Printf("allchild: %#v\n", listText)
// 				for _, c := range listText {
// 					fmt.Printf("child: %#v\n", c)
// 				}
// 				// newItem := []ast.Node{
// 				// 	*ast.Text{}
// 				// }
// 				// for _, v := range listText {
// 				// 	newItem = append(newItem, v)
// 				// }
// 				// list.Children[0].(*ast.Paragraph).Children = append(list.Children[0].(*ast.Paragraph).Children, newItem)
// 				// newItem := listText
// 				// list.Children = append(list.Children, newItem)
// 				// fmt.Printf("text: ~%s~\n", listText)
// 				// fmt.Printf("list item: %#v\n", list.AsContainer().Children[0].(*ast.Paragraph).AsContainer().Content)
// 			}
// 		}

// 		if link, ok := node.(*ast.Link); ok && entering {
// 			isExternalURI := func(uri string) bool {
// 				return (strings.HasPrefix(uri, "https://") || strings.HasPrefix(uri, "http://")) && !strings.Contains(uri, "blog.kowalczyk.info")
// 			}
// 			if isExternalURI(string(link.Destination)) {
// 				link.Destination = []byte("myplace.com")
// 				link.AdditionalAttributes = append(link.AdditionalAttributes, `target="_blank"`)
// 			}
// 		}

// 		return ast.GoToNext
// 	})
// 	return doc
// }

// func modifyAstExample() {
// 	md := []byte(mds)

// 	extensions := parser.CommonExtensions
// 	p := parser.NewWithExtensions(extensions)
// 	doc := p.Parse(md)

// 	doc = modifyAst(doc)

// 	htmlFlags := html.CommonFlags
// 	opts := html.RendererOptions{Flags: htmlFlags}
// 	renderer := html.NewRenderer(opts)
// 	html := markdown.Render(doc, renderer)

// 	mdrender := mdrenderer.NewRenderer()
// 	mdDoc := markdown.Render(doc, mdrender)

// 	fmt.Printf("-- Markdown:\n%s\n\n--- HTML:\n%s\n", mdDoc, html)
// }

func astWalkTest() {
	// Read markdown file
	filename := "/Users/christopher.maahs/src/christopher-maahs/Daily_Notes/2024/January/4th/2024-01-25-Thursday-1.md"
	source, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// Parse the markdown file
	reader := text.NewReader(source)
	parser := goldmark.DefaultParser()
	document := parser.Parse(reader)

	// Traverse the AST and find the "## MRs-Opened" section
	var mrOpenedHeader ast.Node
	ast.Walk(document, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering {
			if header, ok := node.(*ast.Heading); ok && header.Level == 2 {
				// fmt.Printf("header: %#v\n", header)
				// fmt.Printf("text: %s", header.Text(source))
				headingText := string(header.Text(source))
				if headingText == "MRs-Opened" {
					mrOpenedHeader = node
				}
			}
		}
		return ast.WalkContinue, nil
	})

	// Insert new item if the section was found
	if mrOpenedHeader != nil {
		listItem := ast.NewListItem(0)
		paragraph := ast.NewParagraph()
		listItem.AppendChild(listItem, paragraph)
		textBlock := ast.NewString([]byte("- [/] [new task MR#30](https://link/to/mr30)"))
		paragraph.AppendChild(paragraph, textBlock)

		// Insert after the header
		mrOpenedHeader.Parent().InsertAfter(mrOpenedHeader.Parent(), mrOpenedHeader, listItem)
	}

	// Render the AST back to markdown
	var buffer bytes.Buffer
	if err := goldmark.DefaultRenderer().Render(&buffer, source, document); err != nil {
		fmt.Println("Error rendering markdown:", err)
		return
	}

	// Save the modified markdown back to file
	if err := os.WriteFile(filename, buffer.Bytes(), os.ModePerm); err != nil {
		fmt.Println("Error writing file:", err)
		return
	}
}

func directTextManipulation() {
	// Open the markdown file.
	filename := "/Users/christopher.maahs/src/christopher-maahs/Daily_Notes/2024/January/4th/2024-01-25-Thursday-1.md"
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
			lines = append(lines, "- [/] [new task MR#30](https://link/to/mr30)")
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

func init() {
	testCmd.AddCommand(modifyCmd)
}
