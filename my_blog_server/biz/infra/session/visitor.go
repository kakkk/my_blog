package session

import (
	"context"
	"fmt"

	"github.com/hertz-contrib/sessions"
)

type VisitorSession struct {
	session sessions.Session
}

func NewVisitorSessionFromCtx(ctx context.Context) (*VisitorSession, error) {
	s, ok := ctx.Value("session").(sessions.Session)
	if !ok {
		return nil, fmt.Errorf("get session from ctx fail")
	}
	return &VisitorSession{
		session: s,
	}, nil
}

func (s *VisitorSession) CheckPostHasVisited(id int64) bool {
	flag, ok := s.session.Get(getPostVisitedKey(id)).(bool)
	if ok && flag == true {
		return true
	}
	return false
}

func (s *VisitorSession) SetPostHasVisited(id int64) {
	s.session.Set(getPostVisitedKey(id), true)
}

func getPostVisitedKey(id int64) string {
	return fmt.Sprintf("v_%v", id)
}
