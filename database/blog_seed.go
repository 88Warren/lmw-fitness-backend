package database

import (
	"errors"
	"log"

	"github.com/88warren/lmw-fitness-backend/models"
	"gorm.io/gorm"
)

// A struct to hold a blog post's data and its corresponding file name.
type blogSeedData struct {
	Title      string
	Excerpt    string
	ImageURL   string
	IsFeatured bool
	Category   string
	FileName   string
}

// BlogSeed populates the database with initial blog posts from HTML files.
func BlogSeed(db *gorm.DB) {
	blogData := []blogSeedData{
		{
			Title:      "If Not Now, Then When? Welcome to LMW Fitness!",
			Excerpt:    "Hi, I’m Laura and I can’t tell you how excited I am to finally share my brand-new website with you!...",
			ImageURL:   "https://placehold.co/1200x600/69002e/FFF?text=Welcome+to+LMW+Fitness",
			IsFeatured: true,
			Category:   "Motivation",
			FileName:   "welcome.html",
		},
		{
			Title:      "5 Nutrition Tips That Actually Make a Difference",
			Excerpt:    "Hi, it’s Laura! Let’s chat about food – because training hard without fuelling properly is like putting cheap petrol in a Ferrari. Here are my top 5 nutrition tips for keeping your energy up, supporting your training and managing weight without the faff.",
			ImageURL:   "https://images.pexels.com/photos/8846349/pexels-photo-8846349.jpeg",
			IsFeatured: false,
			Category:   "Nutrition",
			FileName:   "nutrition_tips.html",
		},
		{
			Title:      "5 Top Tips to Stay Motivated With Your Training",
			Excerpt:    "Hi, it’s Laura! Let’s be honest – staying motivated with training isn’t always easy. Some days the sofa wins, the weather’s rubbish, or life just gets in the way. I’ve been there, and I’ve helped countless clients through it too.",
			ImageURL:   "https://www.pexels.com/photo/woman-doing-body-check-6697178/",
			IsFeatured: false,
			Category:   "Motivation",
			FileName:   "stay_motivated.html",
		},
		{
			Title:      "10 Fitness Tips You’re Probably Overlooking (But Shouldn’t!)",
			Excerpt:    "Hi, it’s Laura here! I thought I’d share some of the little things that often get forgotten when people start a new fitness journey.",
			ImageURL:   "https://www.pexels.com/photo/woman-wearing-oversized-jeans-7991930/",
			IsFeatured: false,
			Category:   "Fitness Tips",
			FileName:   "fitness_tips.html",
		},
		{
			Title:      "Why Recovery is Just as Important as Your Workout",
			Excerpt:    "Hi, it’s Laura! We all love a good sweat session, but here’s the thing: recovery is just as important as the training itself. You can’t build strength, lose fat, or feel energised if you’re constantly running on empty",
			ImageURL:   "https://www.pexels.com/photo/topless-man-in-black-shorts-doing-push-up-training-4803911/",
			IsFeatured: true,
			Category:   "Recovery",
			FileName:   "why_recovery_important.html",
		},
		{
			Title:      "Your Mindset Matters More Than You Think",
			Excerpt:    "Hi, it’s Laura! I want to talk about something that gets overlooked way too often in fitness: your mindset. You can follow the perfect programme, eat the right foods and train consistently – but if your head isn’t in the game, progress can stall",
			ImageURL:   "https://www.pexels.com/photo/woman-checking-weight-on-scales-in-studio-6975466/",
			IsFeatured: true,
			Category:   "Mindset",
			FileName:   "mindset_matters.html",
		},
		{
			Title:      "5 Quick Workouts You Can Actually Stick To",
			Excerpt:    "Hi, it’s Laura! Short on time but still want results? Here are my five go-to quick workouts that get your body moving and your energy up – no excuses!",
			ImageURL:   "https://www.pexels.com/photo/crop-sportswoman-exercising-with-gymnastic-hula-hoop-4498154/",
			IsFeatured: false,
			Category:   "Workouts",
			FileName:   "quick_workouts.html",
		},
		{
			Title:      "15 Quick Workouts to Keep Your Training Fun",
			Excerpt:    "Hi, it’s Laura! Have a look at these 15 x quick workouts to mix things up and keep your training inspiring:",
			ImageURL:   "https://www.pexels.com/photo/a-man-in-black-shirt-holding-a-hamburger-5714271/",
			IsFeatured: false,
			Category:   "Workouts",
			FileName:   "quick_workouts_cont.html",
		},
	}

	for _, data := range blogData {
		var existingBlog models.Blog
		log.Printf("Checking for existing blog post: '%s'", data.Title)
		if err := db.Where("title = ?", data.Title).First(&existingBlog).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				log.Printf("Blog post '%s' not found, attempting to create.", data.Title)
				fullContent, readErr := ReadHTMLFile(data.FileName)
				if readErr != nil {
					log.Printf("Failed to read content for blog post '%s' from 'database/content/blog/%s': %v", data.Title, data.FileName, readErr)
					continue
				}
				log.Printf("Successfully read content for '%s' from 'database/content/blog/%s'", data.Title, data.FileName)

				blog := models.Blog{
					Title:       data.Title,
					Excerpt:     data.Excerpt,
					ImageURL:    data.ImageURL,
					IsFeatured:  data.IsFeatured,
					Category:    data.Category,
					FullContent: fullContent,
				}

				if result := db.Create(&blog); result.Error != nil {
					log.Printf("Failed to seed blog post '%s': %v", blog.Title, result.Error)
				} else {
					log.Printf("Seeded blog post: %s", blog.Title)
				}
			} else {
				log.Printf("Database error while checking for blog post '%s': %v", data.Title, err)
			}
		} else {
			log.Printf("Blog post '%s' already exists, skipping.", data.Title)
		}
	}
	log.Println("Finished blog data seeding.")
}
