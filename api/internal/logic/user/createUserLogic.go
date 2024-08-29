package user

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"go-game/common/consts"
	"strconv"
	
	"go-game/api/internal/svc"
	"go-game/api/internal/types"
	
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateUserLogic {
	return &CreateUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateUserLogic) CreateUser(req *types.CreateUserReq) (resp *types.CreateUserResp, err error) {
	userId, err := l.svcCtx.RedisClient.Incr(consts.RedisPrefix + "user_id_generator")
	if err != nil {
		return nil, errors.Wrap(err, "生成用户ID失败")
	}
	
	user := types.User{
		UserId:   userId,
		Nickname: req.Nickname,
		Avatar:   req.Avatar,
	}
	
	bytes, err := json.Marshal(user)
	if err != nil {
		return nil, errors.Wrap(err, "序列化用户信息失败")
	}
	err = l.svcCtx.RedisClient.Hset(consts.RedisPrefix+"user_info", strconv.FormatInt(userId, 10), string(bytes))
	if err != nil {
		return nil, errors.Wrap(err, "保存用户信息失败")
	}
	return &types.CreateUserResp{UserId: userId}, nil
}
