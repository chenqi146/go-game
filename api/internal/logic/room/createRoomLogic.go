package room

import (
	"context"
	"encoding/json"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"go-game/api/internal/logic/user"
	"go-game/common/consts"
	"go-game/common/ctx"
	"strconv"
	
	"go-game/api/internal/svc"
	"go-game/api/internal/types"
	
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateRoomLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateRoomLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateRoomLogic {
	return &CreateRoomLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateRoomLogic) CreateRoom(req *types.CreateRoomReq) (resp *types.CreateRoomResp, err error) {
	
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
	
	roomId, err := l.svcCtx.RedisClient.Incr(consts.RedisPrefix + "room_id_generator")
	if err != nil {
		return nil, errors.Wrap(err, "生成房间ID失败")
	}
	
	room := types.Room{
		RoomId:           roomId,
		PlayerNum:        req.PlayerNum,
		UndercoverNum:    req.UndercoverNum,
		CivilianNum:      req.CivilianNum,
		HomeOwnersUserId: userId,
		Status:           0,
		Users: []types.User{
			userInfo,
		},
		Game: types.GameInfo{},
	}
	
	bytes, _ := json.Marshal(room)
	err = l.svcCtx.RedisClient.Hset(consts.RedisPrefix+"room_info", strconv.FormatInt(roomId, 10), string(bytes))
	if err != nil {
		return nil, errors.Wrap(err, "保存房间信息失败")
	}
	
	return &types.CreateRoomResp{RoomId: roomId}, nil
}
