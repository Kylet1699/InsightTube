package main

// Import necessary packages
import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

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
}

type VideoStats struct {
	Title        string `json:"title"`
	ViewCount    uint64 `json:"viewCount"`
	LikeCount    uint64 `json:"likeCount"`
	CommentCount uint64 `json:"commentCount"`
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
	}

	// Fetch comments
	commentsResponse, err := service.CommentThreads.List([]string{"snippet"}).VideoId(videoID).MaxResults(100).Do()
	if err != nil {
		return nil, fmt.Errorf("error fetching comments: %v", err)
	}

	var comments []Comment

	for _, item := range commentsResponse.Items {
		sentiment, err := analyzeSentiment(item.Snippet.TopLevelComment.Snippet.TextDisplay)
		if err != nil {
			return nil, fmt.Errorf("error analyzing sentiment: %v", err)
		}
		comment := Comment{
			ID:        item.Snippet.TopLevelComment.Id,
			Text:      item.Snippet.TopLevelComment.Snippet.TextDisplay,
			Sentiment: sentiment,
		}
		comments = append(comments, comment)
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
