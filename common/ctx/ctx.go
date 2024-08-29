package ctx

import (
	"context"
	"go-game/common/consts"
	"strconv"
)

func GetUidFromCtx(ctx context.Context) (int64, error) {
	i, err := strconv.ParseInt(ctx.Value(consts.UserId).(string), 10, 64)
	return i, err
}
