package quizlet

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type leaderboardUpdate struct {
	Score          int  `json:"score"`
	PreviousRecord int  `json:"previous_record"`
	TooSmall       int  `json:"too_small"`
	TimeStarted    int  `json:"time_started"`
	SelectedOnly   bool `json:"selectedOnly"`
}

// UpdateLeaderboard updates the logged-in user's score
func (c *Client) UpdateLeaderboard(id string, score int) (err error) {
	if _, err = c.Request(http.MethodGet, "/"+id, nil); err != nil {
		return
	}

	data, err := json.Marshal(leaderboardUpdate{
		Score:          score,
		PreviousRecord: 0,
		TooSmall:       0,
		TimeStarted:    int(time.Now().Unix() - int64(score/10)),
		SelectedOnly:   false,
	})
	if err != nil {
		return
	}

	body, err := json.Marshal(map[string]string{
		"data": obfuscate(data),
	})
	if err != nil {
		return
	}

	res, err := c.Request(http.MethodPost, "/"+id+"/match/highscores", bytes.NewReader(body))
	io.Copy(os.Stdout, res.Body)
	fmt.Println("")
	return
}
