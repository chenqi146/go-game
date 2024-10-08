// Code generated by goctl. DO NOT EDIT.
package handler

import (
	"net/http"
	"time"

	room "go-game/api/internal/handler/room"
	user "go-game/api/internal/handler/user"
	ws "go-game/api/internal/handler/ws"
	"go-game/api/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/room/changeWord",
				Handler: room.ChangeWordHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/room/create",
				Handler: room.CreateRoomHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/room/eliminateUser",
				Handler: room.EliminateUserHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/room/gameOver",
				Handler: room.GameOverHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/room/info",
				Handler: room.RoomInfoHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/room/join",
				Handler: room.JoinRoomHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/room/startGame",
				Handler: room.StartGameHandler(serverCtx),
			},
		},
		rest.WithPrefix("/api/v1"),
		rest.WithTimeout(30000*time.Millisecond),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/user/create",
				Handler: user.CreateUserHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/user/query",
				Handler: user.QueryUserHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/user/update",
				Handler: user.UpdateUserHandler(serverCtx),
			},
		},
		rest.WithPrefix("/api/v1"),
		rest.WithTimeout(30000*time.Millisecond),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/ws",
				Handler: ws.WebsocketHandler(serverCtx),
			},
		},
		rest.WithPrefix("/api/v1"),
		rest.WithTimeout(30000*time.Millisecond),
	)
}
