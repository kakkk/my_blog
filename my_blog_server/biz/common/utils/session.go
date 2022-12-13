package utils

import (
	"context"
	"fmt"
)

func GetUserIDByCtx(ctx context.Context) (int64, error) {
	userID, ok := ctx.Value("user_id").(int64)
	if !ok {
		return 0, fmt.Errorf("get user_id by context fail")
	}
	return userID, nil
}
