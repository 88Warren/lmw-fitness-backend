package database

import (
	"errors"
	"log"

	"github.com/88warren/lmw-fitness-backend/models"
	"gorm.io/gorm"
)

type blogSeedConfig struct {
	IsFeatured bool
	Category   string
	FileName   string
	ImageURL   string
}

func BlogSeed(db *gorm.DB) {
	blogConfigs := []blogSeedConfig{
		{
			IsFeatured: true,
			Category:   "Motivation",
			FileName:   "welcome.html",
			ImageURL:   "https://placehold.co/1200x600/ffcf00/FFF?text=Welcome+to+LMW+Fitness",
		},
		{
			IsFeatured: false,
			Category:   "Nutrition",
			FileName:   "nutrition_tips.html",
			ImageURL:   "https://images.pexels.com/photos/566566/pexels-photo-566566.jpeg?_gl=1*1xl9fus*_ga*MTM3OTk1NTgzOS4xNzU2MzI3MTY2*_ga_8JE65Q40S6*czE3NTY4NDU2MjkkbzQkZzEkdDE3NTY4NDU4NzQkajQ5JGwwJGgw",
		},
		{
			IsFeatured: false,
			Category:   "Motivation",
			FileName:   "stay_motivated.html",
			ImageURL:   "https://images.pexels.com/photos/5238670/pexels-photo-5238670.jpeg?_gl=1*yrfyei*_ga*MTM3OTk1NTgzOS4xNzU2MzI3MTY2*_ga_8JE65Q40S6*czE3NTY4OTY1MTgkbzUkZzEkdDE3NTY4OTY2NzYkajQ1JGwwJGgw",
		},
		{
			IsFeatured: false,
			Category:   "Fitness Tips",
			FileName:   "fitness_tips.html",
			ImageURL:   "https://images.pexels.com/photos/271897/pexels-photo-271897.jpeg?_gl=1*4jsn39*_ga*MTM3OTk1NTgzOS4xNzU2MzI3MTY2*_ga_8JE65Q40S6*czE3NTY4NDEwMTAkbzMkZzEkdDE3NTY4NDEwNjgkajIkbDAkaDA.",
		},
		{
			IsFeatured: true,
			Category:   "Recovery",
			FileName:   "why_recovery_important.html",
			ImageURL:   "https://images.pexels.com/photos/2821823/pexels-photo-2821823.jpeg?_gl=1*11c1f3v*_ga*MTM3OTk1NTgzOS4xNzU2MzI3MTY2*_ga_8JE65Q40S6*czE3NTY4OTY1MTgkbzUkZzEkdDE3NTY4OTY1NDckajMxJGwwJGgw",
		},
		{
			IsFeatured: true,
			Category:   "Mindset",
			FileName:   "mindset_matters.html",
			ImageURL:   "https://images.pexels.com/photos/6690237/pexels-photo-6690237.jpeg?_gl=1*ecyrbb*_ga*MTM3OTk1NTgzOS4xNzU2MzI3MTY2*_ga_8JE65Q40S6*czE3NTY4NDEwMTAkbzMkZzEkdDE3NTY4NDI3MjAkajU5JGwwJGgw",
		},
		{
			IsFeatured: false,
			Category:   "Workouts",
			FileName:   "quick_workouts.html",
			ImageURL:   "https://images.pexels.com/photos/5714271/pexels-photo-5714271.jpeg?_gl=1*1iqkjxu*_ga*MTM3OTk1NTgzOS4xNzU2MzI3MTY2*_ga_8JE65Q40S6*czE3NTY4Mjc5MzQkbzIkZzEkdDE3NTY4MjgwMTQkajQ3JGwwJGgw",
		},
	}

	for _, config := range blogConfigs {
		// log.Printf("Processing blog file: %s", config.FileName)

		htmlContent, readErr := ReadHTMLFile(config.FileName)
		if readErr != nil {
			log.Printf("Failed to read content for blog file '%s': %v", config.FileName, readErr)
			continue
		}
		metadata := ExtractBlogMetadata(htmlContent)

		if metadata.Title == "" {
			log.Printf("Could not extract title from '%s', skipping", config.FileName)
			continue
		}

		var existingBlog models.Blog
		// log.Printf("Checking for existing blog post: '%s'", metadata.Title)
		if err := db.Where("title = ?", metadata.Title).First(&existingBlog).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// log.Printf("Blog post '%s' not found, creating new entry.", metadata.Title)

				imageURL := config.ImageURL
				if imageURL == "" {
					imageURL = metadata.ImageURL
				}

				blog := models.Blog{
					Title:       metadata.Title,
					Excerpt:     metadata.Excerpt,
					ImageURL:    imageURL,
					IsFeatured:  config.IsFeatured,
					Category:    config.Category,
					FullContent: metadata.Content,
				}

				if result := db.Create(&blog); result.Error != nil {
					log.Printf("Failed to seed blog post '%s': %v", blog.Title, result.Error)
				} else {
					log.Printf("Successfully seeded blog post: %s", blog.Title)
					log.Printf("  - Extracted title: %s", metadata.Title)
					log.Printf("  - Extracted excerpt: %.100s...", metadata.Excerpt)
					log.Printf("  - Using image: %s", imageURL)
				}
			} else {
				log.Printf("Database error while checking for blog post '%s': %v", metadata.Title, err)
			}
		} else {
			log.Printf("Blog post '%s' already exists, skipping.", metadata.Title)
		}
	}
	log.Println("Finished blog data seeding.")
}
