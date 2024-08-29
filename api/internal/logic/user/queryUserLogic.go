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

type QueryUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewQueryUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryUserLogic {
	return &QueryUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryUserLogic) QueryUser(req *types.QueryUserReq) (resp *types.QueryUserResp, err error) {
	userId := req.UserId
	userStr, err := l.svcCtx.RedisClient.Hget(consts.RedisPrefix+"user_info", strconv.FormatInt(userId, 10))
	
	if err != nil {
		return nil, errors.Wrap(err, "查询用户信息失败")
	}
	
	user := types.User{}
	err = json.Unmarshal([]byte(userStr), &user)
	if err != nil {
		return nil, errors.Wrap(err, "反序列化用户信息失败")
	}
	
	resp = &types.QueryUserResp{
		UserId:   user.UserId,
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
	}
	return resp, nil
}
