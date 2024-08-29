package response

import (
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"go-game/common/qerr"
	"google.golang.org/grpc/status"
	"net/http"
)

type R struct {
	Code    int64       `json:"code"`
	Message string      `json:"msg"`
	Data    interface{} `json:"data"`
}

func HttpResult(r *http.Request, w http.ResponseWriter, resp interface{}, err error) {
	
	if err == nil {
		response := R{
			Code:    0,
			Message: "success",
			Data:    resp,
		}
		httpx.WriteJson(w, http.StatusOK, response)
		return
	}
	
	serverError := qerr.NewServerError()
	
	causeErr := errors.Cause(err)
	if e, ok := causeErr.(*qerr.CodeError); ok {
		serverError.Code = e.Code
		serverError.Message = e.Message
	} else {
		//  grpc
		if grpcStatus, ok := status.FromError(causeErr); ok {
			serverError.Code = uint32(grpcStatus.Code())
			serverError.Message = grpcStatus.Message()
		}
	}
	
	logx.WithContext(r.Context()).Errorf("[api-error]: %+v ", err)
	httpx.WriteJson(w, http.StatusOK, &R{
		Code:    int64(serverError.Code),
		Message: serverError.Message,
		Data:    nil,
	})
	
}
