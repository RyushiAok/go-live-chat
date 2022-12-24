package handler

import (
	"RealTime_Group_Chat/presen/ws"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetRoomId(ctx *gin.Context) {
	roomId := ctx.Param("roomId")
	ctx.JSON(http.StatusOK, gin.H{"data": roomId})
}

func GetChatRoom(h *ws.Hub) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		roomId := ctx.Param("roomId")
		ws.ServeWs(h)(ctx.Writer, ctx.Request, roomId)
	}
}
