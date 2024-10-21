package main

// Import necessary packages
import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/jonreiter/govader"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

// Comment struct defines the structure for our comment data
// It will be used for JSON encoding/decoding
type Comment struct {
	ID        string `json:"id"`
	Text      string `json:"text"`
	Sentiment string `json:"sentiment"`
	UpdatedAt string `json:"updatedAt"`
}

type VideoStats struct {
	Title        string `json:"title"`
	ViewCount    uint64 `json:"viewCount"`
	LikeCount    uint64 `json:"likeCount"`
	CommentCount uint64 `json:"commentCount"`
	ThumbnailURL string `json:"thumbnailUrl"`
}

type VideoData struct {
	VideoStats VideoStats `json:"videoStats"`
	Comments   []Comment  `json:"comments"`
}

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Retrieve YouTube API key and port number from environment variables
	apiKey := os.Getenv("YOUTUBE_API_KEY")
	port := os.Getenv("PORT")

	// Set up HTTP handler for the /api/comments endpoint
	http.HandleFunc("/api/video", func(w http.ResponseWriter, r *http.Request) {
		videoID := r.URL.Query().Get("videoId")
		if videoID == "" {
			http.Error(w, "Missing videoId parameter", http.StatusBadRequest)
			return
		}

		videoData, err := getVideoData(videoID, apiKey)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(videoData)
	})

	fmt.Printf("Server is running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// getComments fetches comments for a given video ID and performs sentiment analysis
func getVideoData(videoID, apiKey string) (*VideoData, error) {
	ctx := context.Background()
	service, err := youtube.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("error creating YouTube client: %v", err)
	}

	// Fetch video statistics
	videoResponse, err := service.Videos.List([]string{"snippet", "statistics"}).Id(videoID).Do()
	if err != nil {
		return nil, fmt.Errorf("error fetching video data: %v", err)
	}

	if len(videoResponse.Items) == 0 {
		return nil, fmt.Errorf("no video found with ID %s", videoID)
	}

	video := videoResponse.Items[0]
	videoStats := VideoStats{
		Title:        video.Snippet.Title,
		ViewCount:    video.Statistics.ViewCount,
		LikeCount:    video.Statistics.LikeCount,
		CommentCount: video.Statistics.CommentCount,
		ThumbnailURL: video.Snippet.Thumbnails.High.Url,
	}

	// Fetch comments
	// commentsResponse, err := service.CommentThreads.List([]string{"snippet"}).VideoId(videoID).MaxResults(200).Do()
	// if err != nil {
	// 	return nil, fmt.Errorf("error fetching comments: %v", err)
	// }

	// var comments []Comment

	// for _, item := range commentsResponse.Items {
	// 	sentiment, err := analyzeSentiment(item.Snippet.TopLevelComment.Snippet.TextDisplay)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("error analyzing sentiment: %v", err)
	// 	}

	// 	// Parse time string into time object to format
	// 	t, err := time.Parse(time.RFC3339, item.Snippet.TopLevelComment.Snippet.UpdatedAt)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("failed to parse time: %v", err)
	// 	}

	// 	comment := Comment{
	// 		ID:        item.Snippet.TopLevelComment.Id,
	// 		Text:      item.Snippet.TopLevelComment.Snippet.TextDisplay,
	// 		Sentiment: sentiment,
	// 		UpdatedAt: t.Format("2006-01-02 15:04:05"),
	// 	}
	// 	comments = append(comments, comment)
	// }

	// Fetch comments with pagination
	var comments []Comment
	pageToken := ""
	for {
		commentsResponse, err := service.CommentThreads.List([]string{"snippet"}).
			VideoId(videoID).
			MaxResults(100). // Max allowed per page
			PageToken(pageToken).
			Do()
		if err != nil {
			return nil, fmt.Errorf("error fetching comments: %v", err)
		}

		for _, item := range commentsResponse.Items {
			sentiment, err := analyzeSentiment(item.Snippet.TopLevelComment.Snippet.TextDisplay)
			if err != nil {
				return nil, fmt.Errorf("error analyzing sentiment: %v", err)
			}

			// Parse time string into time object to format
			t, err := time.Parse(time.RFC3339, item.Snippet.TopLevelComment.Snippet.UpdatedAt)
			if err != nil {
				return nil, fmt.Errorf("failed to parse time: %v", err)
			}

			comment := Comment{
				ID:        item.Snippet.TopLevelComment.Id,
				Text:      item.Snippet.TopLevelComment.Snippet.TextDisplay,
				Sentiment: sentiment,
				UpdatedAt: t.Format("2006-01-02 15:04:05"),
			}
			comments = append(comments, comment)
		}

		// Check if there are more pages
		pageToken = commentsResponse.NextPageToken
		if pageToken == "" {
			break // No more pages
		}

		// Optional: Add a limit to the number of comments or pages fetched
		if len(comments) >= 500 { // Example: Limit to 1000 comments
			break
		}
	}

	return &VideoData{
		VideoStats: videoStats,
		Comments:   comments,
	}, nil
}

// analyzeSentiment performs sentiment analysis on the given text
// This is a placeholder function and should be replaced with actual sentiment analysis logic
func analyzeSentiment(text string) (string, error) {
	analyzer := govader.NewSentimentIntensityAnalyzer()
	sentiment := analyzer.PolarityScores(text)

	switch {
	case sentiment.Compound > 0.05:
		return "positive", nil
	case sentiment.Compound < -0.05:
		return "negative", nil
	default:
		return "neutral", nil
	}
}
