package converters

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	markdown "github.com/gomarkdown/markdown"
	markdownHtml "github.com/gomarkdown/markdown/html"
	markdownParser "github.com/gomarkdown/markdown/parser"
)

// Convert markdown text to html
func mdToHTML(mdContent string, title string, htmlFullpath string) string {
	htmlArticle := mdToHTMLBody(mdContent)

	htmlContent := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, user-scalable=no, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0">
	<title>%s</title>

	<!-- Custom CSS -->
	<link rel="stylesheet" type="text/css" href="style.css">

</head>
	<body>
	%s
	</body>
</html>
	`, title, htmlArticle)

	if len(htmlFullpath) > 0 {
		fHtml, errCreateFile := os.Create(htmlFullpath)
		if errCreateFile != nil {
			panic(errCreateFile)
		}
		defer fHtml.Close()

		fHtml.WriteString(htmlContent)
		fHtml.Sync()
	}

	return htmlContent
}

// Get HTML body from markdown file
func mdToHTMLBody(md string) string {
	// create markdown parser with extensions
	extensions := markdownParser.CommonExtensions | markdownParser.AutoHeadingIDs | markdownParser.NoEmptyLineBeforeBlock
	p := markdownParser.NewWithExtensions(extensions)
	doc := p.Parse([]byte(md))

	// create HTML renderer with extensions
	htmlFlags := markdownHtml.CommonFlags | markdownHtml.HrefTargetBlank
	opts := markdownHtml.RendererOptions{Flags: htmlFlags}
	renderer := markdownHtml.NewRenderer(opts)

	return string(markdown.Render(doc, renderer))
}

// Converts HTML to PDF using weasyprint
// (requires to be installed with "pip install weasyprint")
func htmlToPdf(htmlFullpath string, pdfFullpath string) {
	cmd := exec.Command("weasyprint", "-p", htmlFullpath, pdfFullpath)
	cmd.Dir = filepath.Dir(htmlFullpath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(string(output))
		log.Fatal(err)
	}
}
