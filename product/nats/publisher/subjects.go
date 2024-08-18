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
	UserSetEvent    NatsSubject = "project.<environment>.authorization.user.setData"
	UserGetEvent    NatsSubject = "project.<environment>.authorization.user.getData"
	UserUpdateEvent NatsSubject = "project.<environment>.authorization.user.updateData"
	UserDeleteEvent NatsSubject = "project.<environment>.authorization.user.deleteData"
	//
	//
	//UserSessionEvent NatsSubject = "project.<environment>.authorization.user.session"
	//UserStructEvent  NatsSubject = "project.<environment>.authorization.user.struct"
)
