package controller

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/nclandrei/synctube/model"
	"github.com/nclandrei/synctube/shared/session"
	"github.com/nclandrei/synctube/shared/youtube/auth"
	"github.com/nclandrei/synctube/shared/youtube/fetcher"
	"github.com/nclandrei/synctube/shared/youtube/file_manager"
	"google.golang.org/api/youtube/v3"
)

const (
	oauthStateString string = "random"
)

func YouTubeGET(w http.ResponseWriter, r *http.Request) {
	authURL := auth.GetAuthorizationURL()
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

func YouTubePOST(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	sess := session.Instance(r)

	if state != oauthStateString {
		fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")
	userID := fmt.Sprintf("%s", sess.Values["id"])

	// create this user's temporary folder where the zip will be created
	err := file_manager.CreateUserFolder(userID)

	if err != nil {
		log.Fatalf("Error in creating the user's temporary folder: %v", err.Error())
	}

	client := auth.GetClient(context.Background(), code, userID)

	service, err := youtube.New(client)
	if err != nil {
		log.Fatalf("Could not retrieve client - %v", err.Error())
	}

	likedVideos := fetcher.FetchLikes(userID, service)
	userPlaylistVideos := fetcher.FetchUserPlaylistVideos(userID, service)

	var toAddVideos []model.Video

	if !isPlaylistNew {
		storedVideos, err := model.VideosByPlaylistID(playlistId)
		if err != nil {
			log.Fatalf("Error when retrieving all videos in playlist: %v", err.Error())
		}
		toAddVideos = diffPlaylistVideos(videos, storedVideos)
		toDeleteVideos := diffPlaylistVideos(storedVideos, videos)
		for _, item := range toDeleteVideos {
			model.VideoDelete(item.ID, item.PlaylistID)
		}
	} else {
		toAddVideos = videos
	}

	for _, item := range toAddVideos {
		err := model.VideoCreate(item.ID, item.Title, item.PlaylistID)
		if err != nil {
			log.Fatalf("Error adding the video to the database: %v", err.Error())
		}
		log.Printf("adding item with title '%v' to mongo", item.Title)
	}

	// Finally, before redirecting to homepage, save the timestamp of the this sync
	err = model.UserUpdateLastSync(userID, time.Now())
	if err != nil {
		log.Fatalf("Error updating last sync timestamp for user: %v", err.Error())
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
