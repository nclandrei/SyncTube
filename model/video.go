package model

import (
"time"

"github.com/nclandrei/YTSync/shared/database"

"gopkg.in/mgo.v2/bson"
)

// *****************************************************************************
// Video
// *****************************************************************************

// Video table contains the information for each Video per user
type Video struct {
    ObjectID    bson.ObjectId   `bson:"_id"`
    ID          string          `db:"id" bson:"id,omitempty"`
    Title       string          `db:"content" bson:"content"`
    URL         string          `db:"url" bson:"url"`
    PlaylistID  bson.ObjectId   `bson:"user_id"`
    UID         uint32          `db:"user_id" bson:"userid,omitempty"`
}

// VideoID returns the video id
func (u *Video) VideoID() string {
    r := ""
    r = u.ObjectID.Hex()
    return r
}

// VideoByUserID gets all Videos for a user
func VideoByUserID(userID string) ([]Video, error) {
    var err error

    var result []Video

    if database.CheckConnection() {
        // Create a copy of mongo
        session := database.Mongo.Copy()
        defer session.Close()
        c := session.DB(database.ReadConfig().MongoDB.Database).C("video")

        // Validate the object id
        if bson.IsObjectIdHex(userID) {
            err = c.Find(bson.M{"user_id": bson.ObjectIdHex(userID)}).All(&result)
        } else {
            err = ErrNoResult
        }
    } else {
        err = ErrUnavailable
    }

    return result, standardizeError(err)
}

// NoteCreate creates a note
func VideoCreate(content string, userID string) error {
    var err error

    now := time.Now()

    if database.CheckConnection() {
        // Create a copy of mongo
        session := database.Mongo.Copy()
        defer session.Close()
        c := session.DB(database.ReadConfig().MongoDB.Database).C("video")

        Video := &Video{
            ObjectID:  bson.NewObjectId(),
            Content:   content,
            UserID:    bson.ObjectIdHex(userID),
            CreatedAt: now,
            UpdatedAt: now,
            Deleted:   0,
        }
        err = c.Insert(Video)
    } else {
        err = ErrUnavailable
    }

    return standardizeError(err)
}

// NoteUpdate updates a note
func VideoUpdate(content string, userID string, VideoID string) error {
    var err error

    now := time.Now()

    if database.CheckConnection() {
        // Create a copy of mongo
        session := database.Mongo.Copy()
        defer session.Close()
        c := session.DB(database.ReadConfig().MongoDB.Database).C("video")
        var Video Video
        Video, err = VideoByID(userID, VideoID)
        if err == nil {
            // Confirm the owner is attempting to modify the note
            if Video.UserID.Hex() == userID {
                Video.UpdatedAt = now
                Video.Content = content
                err = c.UpdateId(bson.ObjectIdHex(VideoID), &Video)
            } else {
                err = ErrUnauthorized
            }
        }
    } else {
        err = ErrUnavailable
    }

    return standardizeError(err)
}

// VideoDelete deletes a note
func VideoDelete(userID string, videoID string) error {
    var err error

    if database.CheckConnection() {
        // Create a copy of mongo
        session := database.Mongo.Copy()
        defer session.Close()
        c := session.DB(database.ReadConfig().MongoDB.Database).C("video")

        var video Video
        video, err = VideoByID(userID, videoID)
        if err == nil {
            // Confirm the owner is attempting to modify the note
            if video.UserID.Hex() == userID {
                err = c.RemoveId(bson.ObjectIdHex(videoID))
            } else {
                err = ErrUnauthorized
            }
        }
    } else {
        err = ErrUnavailable
    }

    return standardizeError(err)
}
