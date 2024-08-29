package room

import (
	"context"
	"encoding/json"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"go-game/api/internal/connection"
	"go-game/api/internal/logic/user"
	"go-game/common/consts"
	"go-game/common/ctx"
	"go-game/common/qerr"
	"strconv"
	
	"go-game/api/internal/svc"
	"go-game/api/internal/types"
	
	"github.com/zeromicro/go-zero/core/logx"
)

type JoinRoomLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewJoinRoomLogic(ctx context.Context, svcCtx *svc.ServiceContext) *JoinRoomLogic {
	return &JoinRoomLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *JoinRoomLogic) JoinRoom(req *types.JoinRoomReq) (resp *types.JoinRoomResp, err error) {
	
	roomId := req.RoomId
	
	roomInfo, err := NewRoomInfoLogic(l.ctx, l.svcCtx).RoomInfo(&types.RoomInfoReq{RoomId: roomId})
	if err != nil {
		return nil, err
	}
	
	if int(roomInfo.PlayerNum) == len(roomInfo.Users) {
		return nil, qerr.NewErrCodeMsg(qerr.ServerError, "房间已满")
	}
	
	userId, err := ctx.GetUidFromCtx(l.ctx)
	if err != nil {
		return nil, errors.Wrap(err, "获取用户id失败")
	}
	
	queryUserResp, err := user.NewQueryUserLogic(l.ctx, l.svcCtx).QueryUser(&types.QueryUserReq{UserId: userId})
	if err != nil {
		return nil, errors.Wrap(err, "查询用户失败")
	}
	
	userInfo := types.User{}
	_ = copier.Copy(&userInfo, queryUserResp)
	roomInfo.Users = append(roomInfo.Users, userInfo)
	
	roomStr, _ := json.Marshal(roomInfo)
	err = l.svcCtx.RedisClient.Hset(consts.RedisPrefix+"room_info", strconv.FormatInt(roomId, 10), string(roomStr))
	if err != nil {
		return nil, errors.Wrap(err, "设置房间序列化信息失败")
	}
	
	// 推送ws
	connection.SendMessageByRoomId(req.RoomId, connection.WsMessage{
		Type:    connection.EliminateUser,
		Message: roomInfo,
	})
	
	return &types.JoinRoomResp{}, nil
}
