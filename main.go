package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

var apiKey = os.Getenv("YOUTUBE_API_KEY")

var logger = log.New(os.Stdout, "youtube: ", log.LstdFlags)

const (
	query    = "football"
	interval = 10 * time.Second
)

// Video structure to store video details
type Video struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	PublishDate string `json:"publish_date"`
	Thumbnails  string `json:"thumbnails"`
}

func fetchVideos() {
	ctx := context.Background()
	client, err := youtube.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		logger.Fatalf("Error creating YouTube client: %v", err)
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
		fmt.Println(video)
		fmt.Println("")
	}
}

func main() {
	fetchVideos()
}

func init() {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		logger.Fatalf("Error loading .env file")
	}
	env, _ := godotenv.Read(".env")
	apiKey = env["YOUTUBE_API_KEY"]
}
