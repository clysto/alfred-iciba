package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	aw "github.com/deanishe/awgo"
)

var wf *aw.Workflow
var CIBA_URL = "https://dict-mobile.iciba.com/interface/index.php?c=word&m=getsuggest&nums=9&is_need_mean=1&word="

type CibaMean struct {
	Part  string   `json:"part"`
	Means []string `json:"means"`
}

type CibaMessage struct {
	Key        string     `json:"key"`
	Paraphrase string     `json:"paraphrase"`
	Value      int        `json:"value"`
	Means      []CibaMean `json:"means"`
}

type CibaResult struct {
	Message []CibaMessage `json:"message"`
	Status  int           `json:"status"`
}

func search(word string) (*CibaResult, error) {
	word = url.QueryEscape(word)
	client := http.DefaultClient
	req, err := http.NewRequest(http.MethodGet, CIBA_URL+word, nil)
	if err != nil {
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var result CibaResult
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func run() {
	words := wf.Args()

	r, err := search(words[0])

	if err != nil {
		wf.FatalError(err)
	}

	for _, msg := range r.Message {
		if len(msg.Means) == 0 {
			continue
		}
		for _, mean := range msg.Means {
			title := mean.Part + strings.Join(mean.Means, "；")
			wf.NewItem(title).Subtitle(msg.Key).Arg(msg.Key).Copytext(title).Valid(true).Cmd().Subtitle("🔈发音")
		}
	}

	wf.SendFeedback()

}

func main() {
	wf = aw.New()
	wf.Run(run)
}
