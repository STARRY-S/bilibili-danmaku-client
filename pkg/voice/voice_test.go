package voice

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
}

func Test_Say(t *testing.T) {
	s := NewVoice("12345")
	err := s.Say()
	if err != nil {
		t.Error(err)
		return
	}

	s = NewVoice("67890")
	err = s.Say()
	if err != nil {
		t.Error(err)
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(3)
	for i := 0; i < 3; i++ {
		go func() {
			s = NewVoice(fmt.Sprintf("这是 %d", i))
			err := s.Say()
			if err != nil {
				t.Error(err)
			}
			wg.Done()
		}()
		time.Sleep(time.Millisecond * 200)
	}
	wg.Wait()
}
