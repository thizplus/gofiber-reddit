package handlers

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gofiber-template/domain/repositories"
	"gofiber-template/domain/services"
	"gofiber-template/pkg/config"
)

type SEOHandler struct {
	postService services.PostService
	config      *config.Config
}

func NewSEOHandler(postService services.PostService, config *config.Config) *SEOHandler {
	return &SEOHandler{
		postService: postService,
		config:      config,
	}
}

// GetSitemap generates sitemap.xml for SEO
func (h *SEOHandler) GetSitemap(c *fiber.Ctx) error {
	baseURL := h.config.App.FrontendURL

	// Get all posts (non-deleted)
	posts, err := h.postService.ListPosts(c.Context(), 0, 10000, repositories.SortByNew, nil) // Get max 10k posts
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error generating sitemap")
	}

	// Build XML
	xml := `<?xml version="1.0" encoding="UTF-8"?>` + "\n"
	xml += `<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">` + "\n"

	// Homepage
	xml += "  <url>\n"
	xml += fmt.Sprintf("    <loc>%s/</loc>\n", baseURL)
	xml += fmt.Sprintf("    <lastmod>%s</lastmod>\n", time.Now().Format("2006-01-02"))
	xml += "    <changefreq>daily</changefreq>\n"
	xml += "    <priority>1.0</priority>\n"
	xml += "  </url>\n"

	// All posts
	for _, post := range posts.Posts {
		xml += "  <url>\n"
		xml += fmt.Sprintf("    <loc>%s/post/%s</loc>\n", baseURL, post.ID)
		xml += fmt.Sprintf("    <lastmod>%s</lastmod>\n", post.UpdatedAt.Format("2006-01-02"))
		xml += "    <changefreq>weekly</changefreq>\n"
		xml += "    <priority>0.8</priority>\n"
		xml += "  </url>\n"
	}

	xml += "</urlset>"

	c.Set("Content-Type", "application/xml")
	return c.SendString(xml)
}

// GetRobotsTxt generates robots.txt
func (h *SEOHandler) GetRobotsTxt(c *fiber.Ctx) error {
	baseURL := h.config.App.FrontendURL

	robotsTxt := "User-agent: *\n"
	robotsTxt += "Allow: /\n"
	robotsTxt += "Disallow: /api/\n"
	robotsTxt += "Disallow: /profile/edit\n"
	robotsTxt += "Disallow: /create-post\n"
	robotsTxt += "Disallow: /settings\n"
	robotsTxt += "\n"
	robotsTxt += fmt.Sprintf("Sitemap: %s/sitemap.xml\n", baseURL)

	c.Set("Content-Type", "text/plain")
	return c.SendString(robotsTxt)
}

// GetRSSFeed generates RSS feed
func (h *SEOHandler) GetRSSFeed(c *fiber.Ctx) error {
	baseURL := h.config.App.FrontendURL

	// Get latest 50 posts
	posts, err := h.postService.ListPosts(c.Context(), 0, 50, repositories.SortByNew, nil)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error generating RSS feed")
	}

	// Build RSS XML
	xml := `<?xml version="1.0" encoding="UTF-8"?>` + "\n"
	xml += `<rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom">` + "\n"
	xml += "  <channel>\n"
	xml += fmt.Sprintf("    <title>%s</title>\n", h.config.App.Name)
	xml += fmt.Sprintf("    <link>%s</link>\n", baseURL)
	xml += "    <description>Latest posts and updates</description>\n"
	xml += "    <language>th</language>\n"
	xml += fmt.Sprintf("    <lastBuildDate>%s</lastBuildDate>\n", time.Now().Format(time.RFC1123Z))
	xml += fmt.Sprintf("    <atom:link href=\"%s/rss.xml\" rel=\"self\" type=\"application/rss+xml\" />\n", baseURL)

	for _, post := range posts.Posts {
		xml += "    <item>\n"
		xml += fmt.Sprintf("      <title><![CDATA[%s]]></title>\n", post.Title)
		xml += fmt.Sprintf("      <link>%s/post/%s</link>\n", baseURL, post.ID)
		xml += fmt.Sprintf("      <guid>%s/post/%s</guid>\n", baseURL, post.ID)
		xml += fmt.Sprintf("      <pubDate>%s</pubDate>\n", post.CreatedAt.Format(time.RFC1123Z))

		// Content preview (first 300 chars)
		content := post.Content
		if len(content) > 300 {
			content = content[:300] + "..."
		}
		xml += fmt.Sprintf("      <description><![CDATA[%s]]></description>\n", content)

		// Author
		xml += fmt.Sprintf("      <author><![CDATA[%s]]></author>\n", post.Author.DisplayName)

		// Categories (tags)
		for _, tag := range post.Tags {
			xml += fmt.Sprintf("      <category><![CDATA[%s]]></category>\n", tag.Name)
		}

		xml += "    </item>\n"
	}

	xml += "  </channel>\n"
	xml += "</rss>"

	c.Set("Content-Type", "application/rss+xml")
	return c.SendString(xml)
}

// GetOpenGraphData returns post data for OG image generation
func (h *SEOHandler) GetOpenGraphData(c *fiber.Ctx) error {
	postID := c.Params("id")
	if postID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Post ID is required",
		})
	}

	// Parse UUID
	postUUID, err := uuid.Parse(postID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid post ID",
		})
	}

	// Get post
	post, err := h.postService.GetPost(c.Context(), postUUID, nil)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Post not found",
		})
	}

	// Build OG data
	ogData := fiber.Map{
		"title":       post.Title,
		"description": post.Content,
		"author":      post.Author.DisplayName,
		"createdAt":   post.CreatedAt,
		"votes":       post.Votes,
		"comments":    post.CommentCount,
		"url":         fmt.Sprintf("%s/post/%s", h.config.App.FrontendURL, post.ID),
	}

	// Add image if available
	if len(post.Media) > 0 {
		ogData["image"] = post.Media[0].URL
	}

	// Add tags
	tags := make([]string, len(post.Tags))
	for i, tag := range post.Tags {
		tags[i] = tag.Name
	}
	ogData["tags"] = tags

	return c.JSON(ogData)
}
