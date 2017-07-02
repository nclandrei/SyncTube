package controller

import (
	"fmt"
	"github.com/nclandrei/YTSync/shared/ytsync"
	"google.golang.org/api/youtube/v3"
	"net/http"
	"context"
)

const (
	oauthStateString string = "random"
)

func YouTubeGET(w http.ResponseWriter, r *http.Request) {
	authURL := ytsync.GetAuthorizationURL()
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

func YouTubePOST(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")

	if state != oauthStateString {
		fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")
	client := ytsync.GetClient(context.Background(), code)

	service, err := youtube.New(client)
	if err != nil {
		fmt.Errorf("Could not retrieve client - %v", err.Error())
	}

	// Start making YouTube API calls.
	// Call the channels.list method. Set the mine parameter to true to
	// retrieve the playlist ID for uploads to the authenticated user's
	// channel.
	//call := service.Channels.List("contentDetails").Mine(true)
	secondCall := service.Playlists.List("contentDetails").Mine(true)

	//response, err := call.Do()
	//if err != nil {
	//	// The channels.list method call returned an error.
	//	log.Fatalf("Error making API call to list channels: %v", err.Error())
	//}

	responseSecond, _ := secondCall.Do()
	for _, secondItem := range responseSecond.Items {
		fmt.Println(secondItem.Snippet.Title)
	}

	//for _, channel := range response.Items {
	//	playlistId := channel.ContentDetails.RelatedPlaylists.Likes
	//	// Print the playlist ID for the list of uploaded videos.
	//	fmt.Printf("Videos in list %s\r\n", playlistId)
	//
	//	nextPageToken := ""
	//	for {
	//		// Call the playlistItems.list method to retrieve the
	//		// list of uploaded videos. Each request retrieves 50
	//		// videos until all videos have been retrieved.
	//		playlistCall := service.PlaylistItems.List("snippet").
	//			PlaylistId(playlistId).
	//			MaxResults(50).
	//			PageToken(nextPageToken)
	//
	//		playlistResponse, err := playlistCall.Do()
	//
	//		if err != nil {
	//			// The playlistItems.list method call returned an error.
	//			log.Fatalf("Error fetching playlist items: %v", err.Error())
	//		}
	//
	//		for _, playlistItem := range playlistResponse.Items {
	//			title := playlistItem.Snippet.Title
	//			videoId := playlistItem.Snippet.ResourceId.VideoId
	//			fmt.Printf("%v, (%v)\r\n", title, videoId)
	//		}
	//
	//		// Set the token to retrieve the next page of results
	//		// or exit the loop if all results have been retrieved.
	//		nextPageToken = playlistResponse.NextPageToken
	//		if nextPageToken == "" {
	//			break
	//		}
	//		fmt.Println()
	//	}
	//}
	http.Redirect(w, r, "/", http.StatusFound)
}
