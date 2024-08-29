package user

import (
	"go-game/common/response"
	"net/http"
	
	"github.com/zeromicro/go-zero/rest/httpx"
	"go-game/api/internal/logic/user"
	"go-game/api/internal/svc"
	"go-game/api/internal/types"
)

func UpdateUserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateUserReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		
		l := user.NewUpdateUserLogic(r.Context(), svcCtx)
		resp, err := l.UpdateUser(&req)
		
		response.HttpResult(r, w, resp, err)
	}
}
