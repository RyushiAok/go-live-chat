package router

import (
	"RealTime_Group_Chat/presen/handler"
	"RealTime_Group_Chat/presen/middleware"
	"RealTime_Group_Chat/presen/ws"
	"github.com/gin-gonic/gin"
)

func Init() *gin.Engine {
	router := gin.Default()
	middleware.Cors(router)

	hub := ws.InitHub()
	go hub.Run()

	router.GET("/", handler.Health)
	router.GET("/room/:roomId", handler.GetRoomId)
	router.GET("/ws/:roomId", handler.GetChatRoom(hub))

	return router
}
