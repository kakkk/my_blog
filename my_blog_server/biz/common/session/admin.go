package session

import (
	"context"
	"fmt"

	"github.com/hertz-contrib/sessions"
)

func SetUserID(ctx context.Context, uid int64) error {
	session, ok := ctx.Value("admin_session").(sessions.Session)
	if !ok {
		return fmt.Errorf("get session from ctx fail")
	}
	session.Set("user_id", uid)
	err := session.Save()
	if err != nil {
		return fmt.Errorf("set session error:[%v]", err)
	}
	return nil
}
