package room

import (
	"context"
	_ "embed"
	"encoding/json"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"go-game/api/internal/connection"
	"go-game/common/consts"
	"go-game/common/ctx"
	"go-game/common/qerr"
	"math/rand"
	"slices"
	"strconv"
	"strings"
	
	"go-game/api/internal/svc"
	"go-game/api/internal/types"
	
	"github.com/zeromicro/go-zero/core/logx"
)

type StartGameLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStartGameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StartGameLogic {
	return &StartGameLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *StartGameLogic) StartGame(req *types.StartGameReq) (resp *types.StartGameResp, err error) {
	
	roomInfoResp, err := NewRoomInfoLogic(l.ctx, l.svcCtx).RoomInfo(&types.RoomInfoReq{RoomId: req.RoomId})
	if err != nil {
		return nil, err
	}
	if roomInfoResp.Status == 1 {
		return nil, qerr.NewServerMessageError("游戏已开始")
	}
	
	userId, err := ctx.GetUidFromCtx(l.ctx)
	if err != nil {
		return nil, errors.Wrap(err, "获取用户id失败")
	}
	if roomInfoResp.HomeOwnersUserId != userId {
		return nil, qerr.NewServerMessageError("只有房主才能开始游戏")
	}
	
	if int(roomInfoResp.PlayerNum) != len(roomInfoResp.Users) {
		return nil, qerr.NewServerMessageError("房间未满")
	}
	
	roomInfoResp.Status = 1
	
	gameId, err := l.svcCtx.RedisClient.Incr(consts.RedisPrefix + "game_id_generator")
	if err != nil {
		return nil, errors.Wrap(err, "生成游戏ID失败")
	}
	
	gameInfo := l.BuildGame(gameId, roomInfoResp)
	
	var r types.StartGameResp
	_ = copier.Copy(&r, &gameInfo)
	
	roomInfoResp.Game = gameInfo
	roomStr, _ := json.Marshal(roomInfoResp)
	err = l.svcCtx.RedisClient.Hset(consts.RedisPrefix+"room_info", strconv.FormatInt(roomInfoResp.RoomId, 10), string(roomStr))
	if err != nil {
		return nil, errors.Wrap(err, "设置房间序列化信息失败")
	}
	// ws推送
	connection.SendMessageByRoomId(roomInfoResp.RoomId, connection.WsMessage{
		Type:    connection.StartGame,
		Message: r,
	})
	
	return &r, nil
}

func (l *StartGameLogic) BuildGame(gameId int64, roomInfoResp *types.RoomInfoResp) types.GameInfo {
	pickedN := rand.Intn(len(l.svcCtx.WordPairs))
	pickedWords := strings.Split(l.svcCtx.WordPairs[pickedN], "——")
	
	gameInfo := types.GameInfo{
		GameId:            gameId,
		UndercoverWord:    pickedWords[pickedN%2],
		CivilianWord:      pickedWords[1-pickedN%2],
		UndercoverUserIds: nil,
		CivilianUserIds:   nil,
		EliminatedUserIds: nil,
	}
	
	gameInfo.CivilianUserIds = make([]int64, 0, roomInfoResp.CivilianNum)
	gameInfo.UndercoverUserIds = make([]int64, 0, roomInfoResp.UndercoverNum)
	gameInfo.EliminatedUserIds = make([]int64, 0, roomInfoResp.PlayerNum)
	perm := rand.Perm(len(roomInfoResp.Users))
	for i := 0; i < int(roomInfoResp.UndercoverNum); i++ {
		gameInfo.UndercoverUserIds = append(gameInfo.UndercoverUserIds, roomInfoResp.Users[perm[i]].UserId)
	}
	
	for i := range roomInfoResp.Users {
		player := roomInfoResp.Users[i]
		if !slices.Contains(gameInfo.UndercoverUserIds, player.UserId) {
			gameInfo.CivilianUserIds = append(gameInfo.CivilianUserIds, player.UserId)
		}
	}
	return gameInfo
}
