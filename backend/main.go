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
	Text      string `json:"text"`      // The actual comment text
	Sentiment string `json:"sentiment"` // The analyzed sentiment of the comment
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
	http.HandleFunc("/api/comments", func(w http.ResponseWriter, r *http.Request) {
		// Extract videoId from the query parameters
		videoID := r.URL.Query().Get("videoId")
		if videoID == "" {
			http.Error(w, "Missing videoId parameter", http.StatusBadRequest)
			return
		}

		// Call getComments function to fetch and analyze comments
		comments, err := getComments(videoID, apiKey)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Set response header to JSON
		w.Header().Set("Content-Type", "application/json")
		// Encode comments as JSON and send the response
		json.NewEncoder(w).Encode(comments)
	})

	// Print server start message
	fmt.Printf("Server is running on port %s\n", port)
	// Start the HTTP server
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// getComments fetches comments for a given video ID and performs sentiment analysis
func getComments(videoID, apiKey string) ([]Comment, error) {
	// Create a new YouTube service client
	ctx := context.Background()
	service, err := youtube.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("error creating YouTube client: %v", err)
	}

	// Prepare the API call to fetch comment threads
	call := service.CommentThreads.List([]string{"snippet"}).VideoId(videoID).MaxResults(100)
	// Execute the API call
	response, err := call.Do()
	if err != nil {
		return nil, fmt.Errorf("error fetching comments: %v", err)
	}

	// Process the response and create Comment structs
	var comments []Comment
	for _, item := range response.Items {
		// Call analyzeSentiment and handle both the sentiment and error
		sentiment, err := analyzeSentiment(item.Snippet.TopLevelComment.Snippet.TextDisplay)
		if err != nil {
			return nil, fmt.Errorf("Error analyzing sentiment: %v", err)
		}

		comment := Comment{
			Text:      item.Snippet.TopLevelComment.Snippet.TextDisplay,
			Sentiment: sentiment,
		}
		comments = append(comments, comment)
	}

	return comments, nil
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
