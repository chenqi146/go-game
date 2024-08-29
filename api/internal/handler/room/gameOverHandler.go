package room

import (
	"go-game/common/response"
	"net/http"
	
	"github.com/zeromicro/go-zero/rest/httpx"
	"go-game/api/internal/logic/room"
	"go-game/api/internal/svc"
	"go-game/api/internal/types"
)

func GameOverHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GameOverReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		
		l := room.NewGameOverLogic(r.Context(), svcCtx)
		resp, err := l.GameOver(&req)
		
		response.HttpResult(r, w, resp, err)
	}
}
