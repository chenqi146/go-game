package room

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"go-game/common/consts"
	"go-game/common/qerr"
	"strconv"
	
	"go-game/api/internal/svc"
	"go-game/api/internal/types"
	
	"github.com/zeromicro/go-zero/core/logx"
)

type RoomInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRoomInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoomInfoLogic {
	return &RoomInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RoomInfoLogic) RoomInfo(req *types.RoomInfoReq) (resp *types.RoomInfoResp, err error) {
	
	exist, err := l.svcCtx.RedisClient.Hexists(consts.RedisPrefix+"room_info", strconv.FormatInt(req.RoomId, 10))
	if err != nil {
		return nil, errors.Wrap(err, "查询房间信息失败")
	}
	if !exist {
		return nil, qerr.NewServerMessageError("房间不存在")
	}
	
	roomStr, err := l.svcCtx.RedisClient.Hget(consts.RedisPrefix+"room_info", strconv.FormatInt(req.RoomId, 10))
	if err != nil {
		return nil, errors.Wrap(err, "查询房间信息失败")
	}
	
	room := types.Room{}
	_ = json.Unmarshal([]byte(roomStr), &room)
	
	return &types.RoomInfoResp{Room: room}, nil
}
