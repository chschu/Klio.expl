package settings

import (
	"github.com/sirupsen/logrus"
	"time"
)

var Instance = &settings{
	entryToStringLocation: mustLoadLocation("Europe/Berlin"),
}

func mustLoadLocation(name string) *time.Location {
	loc, err := time.LoadLocation(name)
	if err != nil {
		logrus.Fatal(err)
	}
	return loc
}

type settings struct {
	entryToStringLocation *time.Location
}

func (s *settings) MaxUTF16LengthForKey() int {
	return 50
}

func (s *settings) MaxUTF16LengthForValue() int {
	return 500
}

func (s *settings) EntryToStringLocation() *time.Location {
	return s.entryToStringLocation
}

func (s *settings) EntryToStringTimeFormat() string {
	return "02.01.2006 15:04"
}

func (s *settings) MaxExplCount() int {
	return 20
}

func (s *settings) ExplTokenValidity() time.Duration {
	return time.Minute * 5
}

func (s *settings) MaxFindCount() int {
	return 20
}

func (s *settings) FindTokenValidity() time.Duration {
	return time.Minute * 5
}

func (s *settings) MaxTopCount() int {
	return 100
}

func (s *settings) HandlerTimeout() time.Duration {
	return time.Second * 2
}
