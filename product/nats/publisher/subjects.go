package main

import (
	"os"
	"strings"
)

type NatsSubject string

const (
	envTag = "<environment>"
)

func (sub NatsSubject) ToString() string {
	subj := string(sub)

	env := os.Getenv("ENVIRONMENT")

	return strings.Replace(subj, envTag, strings.ToLower(env), 1)
}

// subjects
const (
	UserSessionEvent NatsSubject = "project.<environment>.authorization.user.session"
	UserStructEvent  NatsSubject = "project.<environment>.authorization.user.struct"
)
