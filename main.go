package main

import (
	"context"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var apiKey = os.Getenv("YOUTUBE_API_KEY")

var log = logrus.New()

var db *gorm.DB

const (
	query    = "football"
	interval = 20 * time.Second
)

// Video structure to store video details
type Video struct {
	gorm.Model
	ID          int    `json:"id" gorm:"primary_key"`
	Title       string `json:"title"`
	Description string `json:"description"`
	PublishDate string `json:"publish_date"`
	Thumbnails  string `json:"thumbnails"`
}

func fetchVideos(ctx context.Context, db *gorm.DB) {
	log.Info("Starting fetchVideos")
	for {
		select {
		case <-ctx.Done():
			log.Fatal("Context cancelled")
			return
		default:
			log.Info("Fetching videos from YouTube API")
			client, err := youtube.NewService(ctx, option.WithAPIKey(apiKey))
			if err != nil {
				log.Errorf("Error creating YouTube client: %v", err)
			}

			// search arguments parts[]
			parts := []string{"id", "snippet"}

			searchResponse, err := client.Search.List(parts).
				Q(query).
				Type("video").
				MaxResults(10).
				Order("date").
				PublishedAfter("2022-01-01T00:00:00Z").
				Do()
			if err != nil {
				log.Fatalf("Error fetching videos from YouTube API: %v", err)
			}

			for _, item := range searchResponse.Items {
				video := Video{
					Title:       item.Snippet.Title,
					Description: item.Snippet.Description,
					PublishDate: item.Snippet.PublishedAt,
					Thumbnails:  item.Snippet.Thumbnails.High.Url,
				}
				db.Create(&video)
			}
		}
		time.Sleep(interval)
	}

}

func getVideos(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	offset := (page - 1) * limit
	log.Infof("Fetching videos from database page: %d, limit: %d", page, limit)
	var videos []Video
	db.Limit(limit).Offset(offset).Find(&videos).Order("publish_date DESC")

	var total int64
	db.Model(&Video{}).Count(&total)

	pages := int(total) / limit
	if int(total)%limit != 0 {
		pages++
	}

	next_page := page + 1
	if next_page > pages {
		next_page = 0
	}
	prev_page := page - 1
	if prev_page < 1 {
		prev_page = 0
	}

	// return response
	return c.JSON(fiber.Map{
		"status":      "success",
		"data":        videos,
		"next_page":   next_page,
		"prev_page":   prev_page,
		"total_pages": pages,
	})
}

func searchVideos(c *fiber.Ctx) error {
	query := c.Query("q")
	log.Infof("Searching videos from database with query: %s", query)
	var videos []Video
	db.Where("title LIKE ? OR description LIKE ?", "%"+query+"%", "%"+query+"%").Find(&videos).Limit(5).Order("publish_date DESC")

	// return response
	return c.JSON(fiber.Map{
		"status": "success",
		"data":   videos,
	})
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go fetchVideos(ctx, db)

	app := fiber.New()

	app.Get("/videos", getVideos)
	app.Get("/search", searchVideos)

	app.Listen(":8000")
}

func init() {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	env, _ := godotenv.Read(".env")
	apiKey = env["YOUTUBE_API_KEY"]
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:    true,
		TimestampFormat:  time.RFC3339Nano,
		PadLevelText:     true,
		QuoteEmptyFields: true,
	})
	log.SetLevel(logrus.InfoLevel)

	// Connect to the database
	database_url := env["DATABASE_URL"]

	db, err = gorm.Open(mysql.Open(database_url), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	log.Info("Connected to database")

	// Migrate the schema
	db.AutoMigrate(&Video{})
	log.Info("Database migrated")

}
