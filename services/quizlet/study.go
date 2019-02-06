package quizlet

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/davecgh/go-spew/spew"
)

type updateHighScoreData struct {
	Score          int  `json:"score"`
	PreviousRecord int  `json:"previous_record"`
	TooSmall       int  `json:"too_small"`
	TimeStarted    int  `json:"time_started"`
	SelectedOnly   bool `json:"selectedOnly"`
}

// UpdateHighScore updates the logged-in user's score
func (c *Client) UpdateHighScore(id, mode string, score int) (err error) {
	if _, err = c.Request(http.MethodGet, EndpointStudy(id), nil); err != nil {
		return
	}

	data, err := json.Marshal(updateHighScoreData{
		Score:          score,
		PreviousRecord: 0,
		TooSmall:       0,
		TimeStarted:    int(time.Now().Unix()) - (score / 10),
		SelectedOnly:   false,
	})
	if err != nil {
		return
	}

	_, err = c.RequestJSON(http.MethodPost, EndpointStudyHighScore(id, mode), map[string]string{
		"data":     obfuscate(data),
		"response": "json",
	})
	return
}

type endSpellGameData struct {
	CSToken      string `json:"cstoken"`
	SessionID    int    `json:"sessionId"`
	WrongAnswers int    `json:"wrong_answers"`
}

// EndSpellGame starts and ends a spelling game
func (c *Client) EndSpellGame(id string, score int) (err error) {
	if _, err = c.Request(http.MethodGet, EndpointStudy(id), nil); err != nil {
		return
	}

	sessionID, err := c.SessionID(id, StudyModeSpell)
	if err != nil {
		return
	}

	res, err := c.RequestForm(http.MethodPost, EndpointStudyEndGame(id, "spell"), url.Values{
		"cstoken":       {c.Headers["CS-Token"]},
		"sessionId":     {strconv.Itoa(sessionID)},
		"wrong_answers": {strconv.Itoa(score)},
	})
	b, _ := ioutil.ReadAll(res.Body)
	spew.Dump(b)
	return
}
