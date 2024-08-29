package room

import (
	"go-game/common/response"
	"net/http"
	
	"github.com/zeromicro/go-zero/rest/httpx"
	"go-game/api/internal/logic/room"
	"go-game/api/internal/svc"
	"go-game/api/internal/types"
)

func EliminateUserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.EliminateUserReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		
		l := room.NewEliminateUserLogic(r.Context(), svcCtx)
		resp, err := l.EliminateUser(&req)
		response.HttpResult(r, w, resp, err)
	}
}
