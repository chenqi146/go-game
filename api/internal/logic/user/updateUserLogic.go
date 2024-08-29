package user

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"go-game/common/consts"
	"go-game/common/ctx"
	"strconv"
	
	"go-game/api/internal/svc"
	"go-game/api/internal/types"
	
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserLogic {
	return &UpdateUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateUserLogic) UpdateUser(req *types.UpdateUserReq) (resp *types.UpdateUserResp, err error) {
	
	userId, err := ctx.GetUidFromCtx(l.ctx)
	if err != nil {
		return nil, errors.Wrap(err, "获取用户id失败")
	}
	
	hexists, err := l.svcCtx.RedisClient.Hexists(consts.RedisPrefix+"user_info", strconv.FormatInt(userId, 10))
	if err != nil {
		return nil, errors.Wrap(err, "查询用户信息失败")
	}
	if !hexists {
		return nil, errors.New("用户不存在")
	}
	
	hget, err := l.svcCtx.RedisClient.Hget(consts.RedisPrefix+"user_info", strconv.FormatInt(userId, 10))
	if err != nil {
		return nil, errors.Wrap(err, "查询用户信息失败")
	}
	user := types.User{}
	err = json.Unmarshal([]byte(hget), &user)
	if err != nil {
		return nil, errors.Wrap(err, "反序列化用户信息失败")
	}
	user.Nickname = req.Nickname
	
	bytes, err := json.Marshal(user)
	if err != nil {
		return nil, errors.Wrap(err, "序列化用户信息失败")
	}
	err = l.svcCtx.RedisClient.Hset(consts.RedisPrefix+"user_info", strconv.FormatInt(userId, 10), string(bytes))
	if err != nil {
		return nil, errors.Wrap(err, "保存用户信息失败")
	}
	
	return
}
