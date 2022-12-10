package utils

import (
	"context"
	"fmt"
)

func GetUserIDByCtx(ctx context.Context) (int64, error) {
	userIDFromCtx := ctx.Value("user_id")
	if userIDFromCtx == nil {
		return 0, fmt.Errorf("get from ctx fail")
	}
	userID, ok := userIDFromCtx.(int64)
	if !ok {
		return 0, fmt.Errorf("convert user_id fail")
	}
	return userID, nil
}
