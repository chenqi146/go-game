package ws

import (
	"go-game/api/internal/server"
	"net/http"
	
	"go-game/api/internal/svc"
)

func WebsocketHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wsServer := server.NewGatewayWsServer(r.Context(), svcCtx)
		_ = wsServer.ServeWs(w, r)
	}
}
