package database

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

func ReadHTMLFile(filename string) (string, error) {
	filePath := "database/content/blog/" + filename
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file %s: %w", filePath, err)
	}
	log.Printf("Successfully read HTML file: %s", filePath)
	return string(content), nil
}

type BlogMetadata struct {
	Title    string
	Excerpt  string
	ImageURL string
	Content  string
}

func ExtractBlogMetadata(htmlContent string) BlogMetadata {
	metadata := BlogMetadata{}
	titleRegex := regexp.MustCompile(`(?s)<h1[^>]*>(.*?)</h1>`)
	if titleMatch := titleRegex.FindStringSubmatch(htmlContent); len(titleMatch) > 1 {
		metadata.Title = strings.TrimSpace(stripHTMLTags(titleMatch[1]))
	}

	paragraphRegex := regexp.MustCompile(`(?s)<p[^>]*>(.*?)</p>`)
	paragraphMatches := paragraphRegex.FindAllStringSubmatch(htmlContent, -1)

	var excerptParagraphs []string
	for i, match := range paragraphMatches {
		if len(match) > 1 {
			rawText := match[1]
			cleanText := strings.TrimSpace(stripHTMLTags(rawText))
			cleanText = regexp.MustCompile(`\s+`).ReplaceAllString(cleanText, " ")

			if len(cleanText) > 20 && !strings.Contains(strings.ToLower(cleanText), "blog post") {
				excerptParagraphs = append(excerptParagraphs, cleanText)

				if len(cleanText) < 150 && strings.Contains(strings.ToLower(cleanText), "hi,") && i+1 < len(paragraphMatches) {
					nextMatch := paragraphMatches[i+1]
					if len(nextMatch) > 1 {
						nextText := strings.TrimSpace(stripHTMLTags(nextMatch[1]))
						nextText = regexp.MustCompile(`\s+`).ReplaceAllString(nextText, " ")
						if len(nextText) > 20 && !strings.Contains(strings.ToLower(nextText), "blog post") {
							excerptParagraphs = append(excerptParagraphs, nextText)
						}
					}
				}
				break
			}
		}
	}

	if len(excerptParagraphs) > 0 {
		metadata.Excerpt = strings.Join(excerptParagraphs, " ")
	}

	imgRegex := regexp.MustCompile(`<img[^>]+src=["']([^"']+)["'][^>]*>`)
	if imgMatch := imgRegex.FindStringSubmatch(htmlContent); len(imgMatch) > 1 {
		metadata.ImageURL = imgMatch[1]
	}

	content := htmlContent

	if metadata.Title != "" {
		content = titleRegex.ReplaceAllString(content, "")
	}

	for _, excerptText := range excerptParagraphs {
		for _, match := range paragraphMatches {
			if len(match) > 1 {
				cleanText := strings.TrimSpace(stripHTMLTags(match[1]))
				cleanText = regexp.MustCompile(`\s+`).ReplaceAllString(cleanText, " ")
				if cleanText == excerptText {
					content = strings.Replace(content, match[0], "", 1)
					break
				}
			}
		}
	}

	headerRegex := regexp.MustCompile(`<header>.*?</header>`)
	content = headerRegex.ReplaceAllString(content, "")

	blogPostRegex := regexp.MustCompile(`<p[^>]*><strong>Blog Post \d+</strong></p>`)
	content = blogPostRegex.ReplaceAllString(content, "")

	content = strings.TrimSpace(content)
	metadata.Content = content

	return metadata
}

func stripHTMLTags(input string) string {
	re := regexp.MustCompile(`<[^>]*>`)
	return re.ReplaceAllString(input, "")
}

func ScanForNewBlogFiles() ([]BlogMetadata, error) {
	blogDir := "database/content/blog/"
	files, err := os.ReadDir(blogDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read blog directory: %w", err)
	}

	var blogMetadataList []BlogMetadata
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".html") && file.Name() != "templateBlog.html" {
			content, err := ReadHTMLFile(file.Name())
			if err != nil {
				log.Printf("Failed to read %s: %v", file.Name(), err)
				continue
			}

			metadata := ExtractBlogMetadata(content)
			metadata.Title = file.Name() + " -> " + metadata.Title
			blogMetadataList = append(blogMetadataList, metadata)
		}
	}
	return blogMetadataList, nil
}

func ValidateBlogHTML(filename string) error {
	content, err := ReadHTMLFile(filename)
	if err != nil {
		return err
	}

	metadata := ExtractBlogMetadata(content)

	if metadata.Title == "" {
		return fmt.Errorf("no title found - ensure the HTML has an <h1> tag")
	}

	if metadata.Excerpt == "" {
		return fmt.Errorf("no excerpt found - ensure the HTML has meaningful <p> tags")
	}

	if len(metadata.Content) < 50 {
		return fmt.Errorf("content too short after extraction - check HTML structure")
	}

	// log.Printf("✓ Blog HTML validation passed for %s", filename)
	// log.Printf("  Title: %s", metadata.Title)
	// log.Printf("  Excerpt: %.100s...", metadata.Excerpt)
	// log.Printf("  Content length: %d characters", len(metadata.Content))

	return nil
}
