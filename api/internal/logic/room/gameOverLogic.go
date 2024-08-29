package room

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"go-game/api/internal/connection"
	"go-game/common/consts"
	"go-game/common/ctx"
	"go-game/common/qerr"
	"strconv"
	
	"go-game/api/internal/svc"
	"go-game/api/internal/types"
	
	"github.com/zeromicro/go-zero/core/logx"
)

type GameOverLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGameOverLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GameOverLogic {
	return &GameOverLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GameOverLogic) GameOver(req *types.GameOverReq) (resp *types.GameOverResp, err error) {
	
	roomInfoResp, err := NewRoomInfoLogic(l.ctx, l.svcCtx).RoomInfo(&types.RoomInfoReq{RoomId: req.RoomId})
	if err != nil {
		return nil, err
	}
	if roomInfoResp.Status != 1 {
		return nil, qerr.NewServerMessageError("游戏未在进行中")
	}
	
	userId, err := ctx.GetUidFromCtx(l.ctx)
	if err != nil {
		return nil, errors.Wrap(err, "获取用户id失败")
	}
	if roomInfoResp.HomeOwnersUserId != userId {
		return nil, qerr.NewServerMessageError("只有房主才能结束游戏")
	}
	
	roomInfoResp.Game = types.GameInfo{}
	roomInfoResp.Status = 0
	
	bytes, _ := json.Marshal(roomInfoResp)
	err = l.svcCtx.RedisClient.Hset(consts.RedisPrefix+"room_info", strconv.FormatInt(req.RoomId, 10), string(bytes))
	if err != nil {
		return nil, errors.Wrap(err, "设置房间信息失败")
	}
	//  ws推送
	connection.SendMessageByRoomId(req.RoomId, connection.WsMessage{
		Type:    connection.EliminateUser,
		Message: roomInfoResp.Game,
	})
	return
}
