package middleware

import (
	"context"
	"go-game/common/consts"
	"net/http"
)

func UserIdMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从请求头中获取用户id
		userId := r.Header.Get(consts.UserId)
		// 将用户id放入请求上下文中
		ctx := r.Context()
		ctx = context.WithValue(ctx, consts.UserId, userId)
		r = r.WithContext(ctx)
		next(w, r)
	}
}
