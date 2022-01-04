package main

import (
	"io"

	log "github.com/sirupsen/logrus"
)

type noUI struct{}

func (*noUI) ReadLine(prompt string) (string, error) { return "", io.EOF }

func (*noUI) Print(args ...interface{}) { log.Debug(args...) }

func (*noUI) PrintErr(args ...interface{}) { log.Error(args...) }

func (*noUI) IsTerminal() bool { return false }

func (*noUI) WantBrowser() bool { return false }

func (*noUI) SetAutoComplete(complete func(string) string) {}
