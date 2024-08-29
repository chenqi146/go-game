package server

import (
	"context"
	"fmt"
	"github.com/golang-module/carbon"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
	"go-game/api/internal/connection"
	"go-game/api/internal/svc"
	"net/http"
	"strconv"
	"time"
)

type GatewayWsServer struct {
	logx.Logger
	ctx      context.Context
	svcCtx   *svc.ServiceContext
	upgrader *websocket.Upgrader
}

func NewGatewayWsServer(ctx context.Context, svcCtx *svc.ServiceContext) *GatewayWsServer {
	return &GatewayWsServer{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		upgrader: &websocket.Upgrader{
			// ws握手过程中允许跨域
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

func (s *GatewayWsServer) ServeWs(w http.ResponseWriter, r *http.Request) error {
	
	defer func() {
		if err := recover(); err != nil {
			s.Logger.Infof("ws server panic: %v\n", err)
		}
	}()
	
	// 获取ws的url参数
	query := r.URL.Query()
	userId := query.Get("userId")
	if userId == "" {
		http.Error(w, "websocket url param userId is empty", 400)
		return nil
	}
	
	roomId := query.Get("roomId")
	
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		_, _ = w.Write([]byte("websocket upgrade err: " + err.Error()))
		return err
	}
	
	// 分配id
	clientId, err := s.initClientId(conn, userId, roomId)
	if err != nil {
		_, _ = w.Write([]byte("websocket init client id err: " + err.Error()))
		return err
	}
	
	defer s.handlerClientDisconnect(clientId)
	
	// 心跳
	go s.handlerClientHeartbeat(clientId)
	
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			logx.Errorf("read message err: %s", err.Error())
			return err
		} else {
			s.handlerClientMessage(conn, clientId, messageType, message)
		}
		
	}
}

func (s *GatewayWsServer) handlerClientMessage(conn *websocket.Conn, clientId string, messageType int, message []byte) {
	
	// ping  回复pong
	value, ok := connection.Clients.Load(clientId)
	if !ok {
		s.handlerClientDisconnect(clientId)
		return
	}
	s.Logger.Info("[ws]: 收到消息: " + string(message))
	
	client := value.(connection.WebsocketClient)
	
	if messageType != websocket.TextMessage {
		return
	}
	
	if string(message) == "ping" {
		
		if err := conn.WriteMessage(messageType, []byte("pong")); err != nil {
			s.handlerClientDisconnect(clientId)
			return
		}
		
		connection.Clients.Store(clientId, connection.WebsocketClient{
			Id:                clientId,
			Conn:              conn,
			LastHeartbeatTime: carbon.Now().Timestamp(),
			UserId:            client.UserId,
			RoomId:            client.RoomId,
		})
	}
	
}

func (s *GatewayWsServer) handlerClientHeartbeat(clientId string) {
	heartbeatPeriod := s.svcCtx.Config.HeartbeatPeriod
	
	for {
		select {
		case <-time.After(5 * time.Second):
			value, ok := connection.Clients.Load(clientId)
			if !ok {
				break
			}
			client := value.(connection.WebsocketClient)
			
			fmt.Printf("%s: 当前: %d, 客户端: %d, diff: %d\n", clientId, carbon.Now().Timestamp(), client.LastHeartbeatTime, heartbeatPeriod)
			if carbon.Now().Timestamp()-client.LastHeartbeatTime > heartbeatPeriod {
				s.Logger.Infof("[ws]: client: %s heartbeat timeout", clientId)
				s.handlerClientDisconnect(clientId)
				break
			}
		}
	}
}

func (s *GatewayWsServer) initClientId(conn *websocket.Conn, userIdStr string, roomIdStr string) (string, error) {
	clientId := uuid.New().String()
	
	userId, _ := strconv.Atoi(userIdStr)
	roomId := 0
	if len(roomIdStr) != 0 {
		roomId, _ = strconv.Atoi(roomIdStr)
	}
	
	client := connection.WebsocketClient{
		Id:                clientId,
		Conn:              conn,
		LastHeartbeatTime: carbon.Now().Timestamp(),
		UserId:            int64(userId),
		RoomId:            int64(roomId),
	}
	connection.Clients.Store(clientId, client)
	if err := conn.WriteMessage(websocket.TextMessage, []byte(clientId)); err != nil {
		logx.Infof("[ws]: 发生clientId(%s)异常: %v", clientId, err)
		s.handlerClientDisconnect(clientId)
		return "", err
	}
	value, ok := connection.Clients.Load(clientId)
	if ok {
		s.Logger.Info("初始化", value)
	}
	return clientId, nil
}

func (s *GatewayWsServer) handlerClientDisconnect(clientId string) {
	connection.Clients.Delete(clientId)
}
