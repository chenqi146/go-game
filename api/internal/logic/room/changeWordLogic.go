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

type ChangeWordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChangeWordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangeWordLogic {
	return &ChangeWordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChangeWordLogic) ChangeWord(req *types.ChangeWordReq) (resp *types.ChangeWordResp, err error) {
	
	roomInfoResp, err := NewRoomInfoLogic(l.ctx, l.svcCtx).RoomInfo(&types.RoomInfoReq{RoomId: req.RoomId})
	if err != nil {
		return nil, err
	}
	
	userId, err := ctx.GetUidFromCtx(l.ctx)
	if err != nil {
		return nil, errors.Wrap(err, "获取用户id失败")
	}
	if roomInfoResp.HomeOwnersUserId != userId {
		return nil, qerr.NewServerMessageError("只有房主才能更换词语")
	}
	
	gameInfo := NewStartGameLogic(l.ctx, l.svcCtx).BuildGame(roomInfoResp.Game.GameId, roomInfoResp)
	roomInfoResp.Game = gameInfo
	
	bytes, _ := json.Marshal(roomInfoResp)
	err = l.svcCtx.RedisClient.Hset(consts.RedisPrefix+"room_info", strconv.FormatInt(req.RoomId, 10), string(bytes))
	if err != nil {
		return nil, errors.Wrap(err, "设置房间信息失败")
	}
	
	//  ws推送
	connection.SendMessageByRoomId(req.RoomId, connection.WsMessage{
		Type:    connection.ChangeWord,
		Message: gameInfo,
	})
	
	return
}
