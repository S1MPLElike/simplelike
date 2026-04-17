package router

import (
	"GoWorkingProject/Back/controller"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	r.Use(CorsMiddleWare())

	api := r.Group("/api/v1")
	{
		userGroup := api.Group("/user")
		{
			userGroup.POST("/register", controller.UserRegister)
			userGroup.POST("/login", controller.UserLogin)
			userGroup.POST("/logout", controller.UserLogout)
			userGroup.GET("/auto-login", controller.AutoLogin)
			userGroup.GET("/info/:id", controller.GetUserInfo)

			// 上传头像路由，需要登录
			authUserGroup := userGroup.Group("/")
			authUserGroup.Use(controller.AuthMiddleware())
			{
				authUserGroup.POST("/avatar", controller.UploadAvatar)
				authUserGroup.POST("/follow/:id", controller.FollowUser)
				authUserGroup.POST("/unfollow/:id", controller.UnfollowUser)
				authUserGroup.GET("/friends/:id", controller.GetFriends)
				authUserGroup.GET("/check-friendship/:id", controller.CheckFriendship)
				authUserGroup.GET("/check-follow/:id", controller.CheckFollowStatus)
				authUserGroup.GET("/info", controller.GetCurrentUserInfo)
				authUserGroup.POST("/update", controller.UpdateUserInfo)
				authUserGroup.POST("/change-password", controller.ChangePassword)
			}
		}

		// 聊天相关路由
		chatGroup := api.Group("/chat")
		chatGroup.Use(controller.AuthMiddleware())
		{
			chatGroup.POST("/send", controller.SendMessage)
			chatGroup.GET("/history/:id", controller.GetMessageHistory)
			chatGroup.GET("/unread-count", controller.GetUnreadMessageCount)
			chatGroup.GET("/unread-count-by-friend", controller.GetUnreadMessageCountByFriend)
			chatGroup.POST("/upload-image", controller.UploadImage)
			chatGroup.GET("/ws", controller.WebSocketHandler)
			chatGroup.GET("/online/:id", controller.GetOnlineStatus)
		}

		postGroup := api.Group("/post")
		postGroup.Use(controller.AuthMiddleware())
		{
			postGroup.POST("/create", controller.PostCreate)
			postGroup.GET("/hot", controller.GetHotPosts)
			postGroup.GET("/detail/:id", controller.GetPostDetail)
			postGroup.POST("/like/:id", controller.LikePost)
			postGroup.POST("/collect/:id", controller.CollectPost)
			postGroup.GET("/search", controller.SearchPosts)
			postGroup.GET("/comments/:id", controller.GetComments)
			postGroup.POST("/comment", controller.CommentCreate)
			postGroup.GET("/check-like/:id", controller.CheckLikeStatus)
			postGroup.GET("/check-collect/:id", controller.CheckCollectStatus)
			postGroup.POST("/comment/like/:id", controller.LikeComment)
			postGroup.GET("/user/:id", controller.GetUserPosts)
			postGroup.GET("/collected/:id", controller.GetCollectedPosts)
		}
	}
}

func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
func CorsMiddleWare() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:5174", "http://localhost:5175", "http://127.0.0.1:5175", "http://127.0.0.1:5176"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "X-Requested-With", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}
