package user

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"go-game/api/internal/logic/user"
	"go-game/api/internal/svc"
	"go-game/api/internal/types"
)

func QueryUserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.QueryUserReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := user.NewQueryUserLogic(r.Context(), svcCtx)
		resp, err := l.QueryUser(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
