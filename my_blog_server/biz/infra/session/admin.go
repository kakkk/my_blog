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
	return nil
}

func GetUserIDByCtx(ctx context.Context) (int64, error) {
	userID, ok := ctx.Value("user_id").(int64)
	if !ok {
		return 0, fmt.Errorf("get user_id by context fail")
	}
	return userID, nil
}
