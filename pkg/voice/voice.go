package voice

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/STARRY-S/bilibili-danmaku-client/utils"
	"github.com/faiface/beep"
	_ "github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

type Voice struct {
	Content string
	URL     string
}

func NewVoice(s string) *Voice {
	return &Voice{
		Content: s,
	}
}

var (
	headerMap = map[string]string{
		"Accept":     "*/*",
		"Host":       "fanyi.sogou.com",
		"Referer":    "https://fanyi.sogou.com/",
		"User-Agent": "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36",
	}
	voiceApiURL        = "https://fanyi.sogou.com/reventondc/synthesis?text=%s&speed=1&from=translateweb&speaker=3"
	speakerInitialized = false
)

func (ss *Voice) Say() error {
	ss.URL = fmt.Sprintf(voiceApiURL, url.QueryEscape(ss.Content))
	res, err := utils.GetHttpData(ss.URL, headerMap)
	if err != nil {
		return err
	}
	if res.StatusCode == http.StatusServiceUnavailable {
		// retry once
		res, err = utils.GetHttpData(ss.URL, headerMap)
	}
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to get url %q: %v", ss.URL, res.Status)
	}

	streamer, format, err := mp3.Decode(res.Body)
	if err != nil {
		return fmt.Errorf("speaker.Play: %w", err)
	}
	defer streamer.Close()
	if !speakerInitialized {
		if err = speaker.Init(format.SampleRate,
			format.SampleRate.N(time.Second/10)); err != nil {
			return err
		}
		speakerInitialized = true
	}
	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))
	<-done

	return nil
}
