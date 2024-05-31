package core

import (
	"fmt"
	"time"
)

type Session struct {
	stopped   bool
	extended  bool
	timeChain <-chan time.Time
	Data      any
}

func (s *Session) Close() {
	s.stopped = true
}

func (s *Session) Extend() {
	s.extended = true
}

const sessionDuration = time.Hour * 24

func (s *Session) startAwake() {
	s.timeChain = time.After(sessionDuration)
}

func (s *Session) Check() bool {
	if s.stopped {
		return true
	}

	select {
	case _ = <-s.timeChain:
		if !s.extended {
			return true
		}
		s.startAwake()

	default:
	}

	return false
}

type SessionsManager struct {
	sessions  map[string]*Session
	container *ServiceContainer
}

func (s *SessionsManager) GetSession(userId string) (*Session, error) {
	session, ok := s.sessions[userId]

	if !ok {
		return nil, fmt.Errorf("sessions: service \"%v\" haven't session for user with id \"%v\"", s.container.Name, userId)
	}

	return session, nil
}

func (s *SessionsManager) NewSession(userId string, data any) (*Session, error) {
	session := new(Session)
	session.Data = data

	_, ok := s.sessions[userId]

	if ok {
		return nil, fmt.Errorf("sessions: service \"%v\" already have session for user with id \"%v\"", s.container.Name, userId)
	}

	session.startAwake()
	s.sessions[userId] = session
	return session, nil
}

func newSessionsManager(container *ServiceContainer) *SessionsManager {
	sessions := new(SessionsManager)
	sessions.container = container
	sessions.sessions = make(map[string]*Session)

	return sessions
}
