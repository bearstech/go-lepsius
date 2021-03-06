package input

import (
	"fmt"
	"strings"
	"time"

	"github.com/coreos/go-systemd/sdjournal"
	"github.com/factorysh/go-lepsius/model"
	"github.com/mitchellh/mapstructure"
)

func init() {
	register("journald", &Journald{})
}

type Journald struct {
	journal *sdjournal.Journal
}

func (j *Journald) loop(lines chan *model.Line) {
	for {
		cur, err := j.journal.Next()
		if err != nil {
			panic(err)
		}
		if cur == 0 {
			j.journal.Wait(1 * time.Second)
		}
		entry, err := j.journal.GetEntry()
		if err != nil {
			panic(err)
		}
		l := make(map[string]interface{})
		for k, v := range entry.Fields {
			l[k] = v
		}
		lines <- &model.Line{
			Values: l,
		}
	}
}

func (j *Journald) Lines() chan *model.Line {
	lines := make(chan *model.Line)
	go j.loop(lines)
	return lines
}

type JournaldConf struct {
	Since   time.Duration
	Matches map[string]string
}

func (j *Journald) Configure(conf map[string]interface{}) error {
	var jconf JournaldConf
	err := mapstructure.Decode(conf, &jconf)
	if err != nil {
		return err
	}
	j.journal, err = sdjournal.NewJournal()
	if err != nil {
		return err
	}
	for key, value := range jconf.Matches {
		m := fmt.Sprintf("%v=%v", strings.ToUpper(key), value)
		j.journal.AddMatch(m)
	}
	err = j.journal.SeekTail()
	if err != nil {
		return err
	}
	return nil
}
