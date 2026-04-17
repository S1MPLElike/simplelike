package controller

import (
	"GoWorkingProject/Back/config"
	"GoWorkingProject/Back/model"
	"GoWorkingProject/Back/utils"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone" binding:"required,len=11"`
}

func UserRegister(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, utils.Fail(utils.CodeParamError, "参数错误: "+err.Error()))
		return
	}

	user := &model.User{
		Username: req.Username,
		Password: req.Password, // 明文存储
		Email:    req.Email,
		Phone:    req.Phone,
	}

	if err := model.DB.Create(user).Error; err != nil {
		c.JSON(200, utils.Fail(utils.CodeDBError, "注册失败: "+err.Error()))
		return
	}

	c.JSON(200, utils.Success(gin.H{
		"user_id":  user.ID,
		"username": user.Username,
		"email":    user.Email,
	}))
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func UserLogin(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, utils.Fail(utils.CodeParamError, "参数错误: "+err.Error()))
		return
	}

	var user model.User
	if err := model.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(200, utils.Fail(utils.CodeUserNotFound, "用户不存在"))
		return
	}

	if user.Password != req.Password {
		c.JSON(200, utils.Fail(utils.CodePwdError, "密码错误"))
		return
	}

	token := utils.GenerateToken()
	if err := utils.SaveTokenToRedis(token, user.ID); err != nil {
		c.JSON(200, utils.Fail(utils.CodeRedisError, "登录失败: "+err.Error()))
		return
	}

	utils.SetTokenCookie(c, token)

	c.JSON(200, utils.Success(gin.H{
		"user_id":  user.ID,
		"username": user.Username,
	}))
}

func UserLogout(c *gin.Context) {
	token, err := c.Cookie(config.TokenConfig.CookieName)
	if err != nil {
		c.JSON(200, utils.Fail(utils.CodeNotLogin, "未登录，无需退出"))
		return
	}

	if err := utils.DeleteToken(token); err != nil {
		c.JSON(200, utils.Fail(utils.CodeRedisError, "退出登录失败: "+err.Error()))
		return
	}

	utils.ClearTokenCookie(c)
	c.JSON(200, utils.Success(nil))
}

// GetUserByID 根据用户ID查询用户信息
// GetUserByID 根据用户ID查询用户信息
func GetUserByID(userID uint64) (*model.User, error) {
	var user model.User
	if err := model.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// 新增：自动登录接口（仅验证Cookie，返回用户信息）
func AutoLogin(c *gin.Context) {
	// 1. 复用鉴权逻辑：读取并验证Cookie中的Token
	token, err := c.Cookie(config.TokenConfig.CookieName)
	if err != nil {
		c.JSON(200, utils.Fail(utils.CodeNotLogin, "未登录"))
		return
	}

	userID, err := utils.VerifyToken(token)
	if err != nil {
		c.JSON(200, utils.Fail(utils.CodeTokenInvalid, "无效的token"))
		return
	}

	// 2. 验证通过：查询用户信息（可选，根据需要返回）
	user, err := GetUserByID(userID)
	if err != nil {
		c.JSON(200, utils.Fail(utils.CodeUserNotFound, "用户不存在"))
		return
	}

	// 3. 返回用户信息（完成自动登录）
	c.JSON(200, utils.Success(gin.H{
		"user_id":  user.ID,
		"username": user.Username,
		"email":    user.Email, // 按需返回，敏感信息别返
	}))

}

// 获取用户信息接口
func GetUserInfo(c *gin.Context) {
	userIDStr := c.Param("id")

	// 转换用户ID为uint64
	var userID uint64
	_, err := fmt.Sscanf(userIDStr, "%d", &userID)
	if err != nil {
		c.JSON(200, utils.Fail(utils.CodeParamError, "无效的用户ID"))
		return
	}

	var user model.User
	if err := model.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(200, utils.Fail(utils.CodeUserNotFound, "用户不存在"))
		return
	}

	// 计算用户的点赞、评论、收藏总数
	var likeCount, commentCount, collectCount int64
	model.DB.Model(&model.Post{}).Where("user_id = ?", userID).Select("SUM(like_count)").Scan(&likeCount)
	model.DB.Model(&model.Comment{}).Where("user_id = ?", userID).Count(&commentCount)
	model.DB.Model(&model.Collect{}).Where("user_id = ?", userID).Count(&collectCount)

	// 计算用户的关注者数量和关注数量
	var followers, following int64
	model.DB.Model(&model.Follow{}).Where("following_id = ?", userID).Count(&followers)
	model.DB.Model(&model.Follow{}).Where("follower_id = ?", userID).Count(&following)

	c.JSON(200, utils.Success(gin.H{
		"user_id":     user.ID,
		"username":    user.Username,
		"email":       user.Email,
		"avatar":      user.Avatar,
		"followers":   followers,
		"following":   following,
		"likes":       likeCount,
		"comments":    commentCount,
		"collections": collectCount,
	}))
}
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie(config.TokenConfig.CookieName)
		if err != nil {
			c.JSON(200, utils.Fail(utils.CodeNotLogin, "未登录，请先登录"))
			c.Abort()
			return
		}

		userID, err := utils.VerifyToken(token)
		if err != nil {
			c.JSON(200, utils.Fail(utils.CodeTokenInvalid, err.Error()))
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}

// GetUserIDFromContext 从上下文中获取用户ID
func GetUserIDFromContext(c *gin.Context) (uint64, error) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, fmt.Errorf("未登录")
	}
	return userID.(uint64), nil
}

// FollowUser 关注用户接口
func FollowUser(c *gin.Context) {
	userIDStr := c.Param("id")
	var targetUserID uint64
	_, err := fmt.Sscanf(userIDStr, "%d", &targetUserID)
	if err != nil {
		c.JSON(200, utils.Fail(utils.CodeParamError, "无效的用户ID"))
		return
	}

	// 获取当前登录用户ID
	currentUserID, err := GetUserIDFromContext(c)
	if err != nil {
		c.JSON(200, utils.Fail(utils.CodeNotLogin, "未登录"))
		return
	}

	// 检查是否关注自己
	if currentUserID == targetUserID {
		c.JSON(200, utils.Fail(utils.CodeOperationFailed, "不能关注自己"))
		return
	}

	// 检查是否已关注
	var count int64
	model.DB.Model(&model.Follow{}).Where("follower_id = ? AND following_id = ?", currentUserID, targetUserID).Count(&count)
	if count > 0 {
		c.JSON(200, utils.Fail(utils.CodeOperationFailed, "已关注该用户"))
		return
	}

	// 创建关注关系
	follow := model.Follow{
		FollowerID:  currentUserID,
		FollowingID: targetUserID,
	}
	if err := model.DB.Create(&follow).Error; err != nil {
		c.JSON(200, utils.Fail(utils.CodeOperationFailed, "关注失败: "+err.Error()))
		return
	}
	// 打印日志，确认关注关系创建成功
	fmt.Printf("关注关系创建成功: 关注者ID=%d, 被关注者ID=%d\n", currentUserID, targetUserID)

	// 检查对方是否已关注自己，判断是否成为好友
	var reverseCount int64
	model.DB.Model(&model.Follow{}).Where("follower_id = ? AND following_id = ?", targetUserID, currentUserID).Count(&reverseCount)
	isFriend := reverseCount > 0

	c.JSON(200, utils.Success(gin.H{
		"followed":  true,
		"is_friend": isFriend,
	}))
}

// UnfollowUser 取消关注用户接口
func UnfollowUser(c *gin.Context) {
	userIDStr := c.Param("id")
	var targetUserID uint64
	_, err := fmt.Sscanf(userIDStr, "%d", &targetUserID)
	if err != nil {
		c.JSON(200, utils.Fail(utils.CodeParamError, "无效的用户ID"))
		return
	}

	// 获取当前登录用户ID
	currentUserID, err := GetUserIDFromContext(c)
	if err != nil {
		c.JSON(200, utils.Fail(utils.CodeNotLogin, "未登录"))
		return
	}

	// 检查是否已关注
	var count int64
	model.DB.Model(&model.Follow{}).Where("follower_id = ? AND following_id = ?", currentUserID, targetUserID).Count(&count)
	if count == 0 {
		c.JSON(200, utils.Fail(utils.CodeOperationFailed, "未关注该用户"))
		return
	}

	// 取消关注
	if err := model.DB.Where("follower_id = ? AND following_id = ?", currentUserID, targetUserID).Delete(&model.Follow{}).Error; err != nil {
		c.JSON(200, utils.Fail(utils.CodeOperationFailed, "取消关注失败"))
		return
	}

	// 检查是否仍然是好友（对方是否还关注自己）
	var reverseCount int64
	model.DB.Model(&model.Follow{}).Where("follower_id = ? AND following_id = ?", targetUserID, currentUserID).Count(&reverseCount)
	isStillFriend := reverseCount > 0

	c.JSON(200, utils.Success(gin.H{
		"unfollowed":      true,
		"is_still_friend": isStillFriend,
	}))
}

// CheckFriendship 检查两个用户是否是好友接口
func CheckFriendship(c *gin.Context) {
	userIDStr := c.Param("id")
	targetIDStr := c.Query("target_id")

	var userID, targetID uint64
	_, err := fmt.Sscanf(userIDStr, "%d", &userID)
	if err != nil {
		c.JSON(200, utils.Fail(utils.CodeParamError, "无效的用户ID"))
		return
	}

	_, err = fmt.Sscanf(targetIDStr, "%d", &targetID)
	if err != nil {
		c.JSON(200, utils.Fail(utils.CodeParamError, "无效的目标用户ID"))
		return
	}

	// 检查是否互相关注
	var count1, count2 int64
	model.DB.Model(&model.Follow{}).Where("follower_id = ? AND following_id = ?", userID, targetID).Count(&count1)
	model.DB.Model(&model.Follow{}).Where("follower_id = ? AND following_id = ?", targetID, userID).Count(&count2)

	isFriend := count1 > 0 && count2 > 0

	c.JSON(200, utils.Success(gin.H{
		"is_friend": isFriend,
	}))
}

// CheckFollowStatus 检查当前用户是否关注了目标用户
func CheckFollowStatus(c *gin.Context) {
	targetIDStr := c.Param("id")
	var targetID uint64
	_, err := fmt.Sscanf(targetIDStr, "%d", &targetID)
	if err != nil {
		c.JSON(200, utils.Fail(utils.CodeParamError, "无效的目标用户ID"))
		return
	}

	// 获取当前登录用户ID
	currentUserID, err := GetUserIDFromContext(c)
	if err != nil {
		c.JSON(200, utils.Fail(utils.CodeNotLogin, "未登录"))
		return
	}

	// 检查是否已关注
	var count int64
	model.DB.Model(&model.Follow{}).Where("follower_id = ? AND following_id = ?", currentUserID, targetID).Count(&count)

	isFollowing := count > 0

	c.JSON(200, utils.Success(gin.H{
		"following": isFollowing,
	}))
}

// GetFriends 获取好友列表接口
func GetFriends(c *gin.Context) {
	userIDStr := c.Param("id")
	var userID uint64
	_, err := fmt.Sscanf(userIDStr, "%d", &userID)
	if err != nil {
		c.JSON(200, utils.Fail(utils.CodeParamError, "无效的用户ID"))
		return
	}

	// 获取好友列表：用户A关注用户B，且用户B关注用户A
	var friends []model.User
	model.DB.Table("users").Joins("INNER JOIN follows f1 ON users.id = f1.following_id").Joins("INNER JOIN follows f2 ON users.id = f2.follower_id").Where("f1.follower_id = ? AND f2.following_id = ?", userID, userID).Find(&friends)

	// 构造返回数据
	var friendList []gin.H
	for _, friend := range friends {
		friendList = append(friendList, gin.H{
			"id":     friend.ID,
			"name":   friend.Username,
			"avatar": friend.Avatar,
		})
	}

	c.JSON(200, utils.Success(gin.H{
		"friends": friendList,
	}))
}

// 上传头像接口
func UploadAvatar(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(200, utils.Fail(utils.CodeNotLogin, "未登录，请先登录"))
		return
	}

	// 获取上传的文件
	file, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(200, utils.Fail(utils.CodeParamError, "获取文件失败: "+err.Error()))
		return
	}

	// 验证文件类型
	ext := filepath.Ext(file.Filename)
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" {
		c.JSON(200, utils.Fail(utils.CodeParamError, "只支持jpg、jpeg、png、gif格式的图片"))
		return
	}

	// 创建保存目录
	uploadDir := "./uploads/avatars"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(200, utils.Fail(utils.CodeServerError, "创建目录失败: "+err.Error()))
		return
	}

	// 生成文件名
	fileName := fmt.Sprintf("%d_%d%s", userID, time.Now().Unix(), ext)
	filePath := filepath.Join(uploadDir, fileName)

	// 保存文件
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(200, utils.Fail(utils.CodeServerError, "保存文件失败: "+err.Error()))
		return
	}

	// 更新用户头像URL
	avatarURL := "/uploads/avatars/" + fileName

	// 先查询用户当前的头像URL，以便删除旧头像
	var oldAvatar string
	if err := model.DB.Model(&model.User{}).Where("id = ?", userID).Select("avatar").Scan(&oldAvatar).Error; err == nil && oldAvatar != "" {
		// 解析旧头像的文件名
		oldFileName := filepath.Base(oldAvatar)
		oldFilePath := filepath.Join(uploadDir, oldFileName)
		// 删除旧头像文件
		if err := os.Remove(oldFilePath); err != nil {
			// 记录错误但不影响新头像的上传
			fmt.Printf("删除旧头像失败: %v\n", err)
		}
	}

	if err := model.DB.Model(&model.User{}).Where("id = ?", userID).Update("avatar", avatarURL).Error; err != nil {
		c.JSON(200, utils.Fail(utils.CodeDBError, "更新头像失败: "+err.Error()))
		return
	}

	c.JSON(200, utils.Success(gin.H{
		"avatar": avatarURL,
	}))
}

// 获取当前登录用户信息接口
func GetCurrentUserInfo(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(200, utils.Fail(utils.CodeNotLogin, "未登录，请先登录"))
		return
	}

	var user model.User
	if err := model.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(200, utils.Fail(utils.CodeUserNotFound, "用户不存在"))
		return
	}

	c.JSON(200, utils.Success(gin.H{
		"user_id":  user.ID,
		"username": user.Username,
		"email":    user.Email,
		"bio":      user.Bio,
		"avatar":   user.Avatar,
	}))
}

// 修改用户信息请求结构
type UpdateUserRequest struct {
	Username string `json:"username" binding:"omitempty,min=3,max=20"`
	Email    string `json:"email" binding:"omitempty,email"`
	Bio      string `json:"bio" binding:"omitempty,max=500"`
}

// 修改用户信息接口
func UpdateUserInfo(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(200, utils.Fail(utils.CodeNotLogin, "未登录，请先登录"))
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, utils.Fail(utils.CodeParamError, "参数错误: "+err.Error()))
		return
	}

	// 构建更新数据
	updateData := make(map[string]interface{})
	if req.Username != "" {
		updateData["username"] = req.Username
	}
	if req.Email != "" {
		updateData["email"] = req.Email
	}
	if req.Bio != "" {
		updateData["bio"] = req.Bio
	}

	// 如果没有要更新的字段，直接返回成功
	if len(updateData) == 0 {
		c.JSON(200, utils.Success(nil))
		return
	}

	// 执行更新
	if err := model.DB.Model(&model.User{}).Where("id = ?", userID).Updates(updateData).Error; err != nil {
		c.JSON(200, utils.Fail(utils.CodeDBError, "更新用户信息失败: "+err.Error()))
		return
	}

	c.JSON(200, utils.Success(nil))
}

// 修改密码请求结构
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// 修改密码接口
func ChangePassword(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(200, utils.Fail(utils.CodeNotLogin, "未登录，请先登录"))
		return
	}

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, utils.Fail(utils.CodeParamError, "参数错误: "+err.Error()))
		return
	}

	// 验证旧密码
	var user model.User
	if err := model.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(200, utils.Fail(utils.CodeUserNotFound, "用户不存在"))
		return
	}

	if user.Password != req.OldPassword {
		c.JSON(200, utils.Fail(utils.CodePwdError, "旧密码错误"))
		return
	}

	// 更新密码
	if err := model.DB.Model(&model.User{}).Where("id = ?", userID).Update("password", req.NewPassword).Error; err != nil {
		c.JSON(200, utils.Fail(utils.CodeDBError, "修改密码失败: "+err.Error()))
		return
	}

	c.JSON(200, utils.Success(nil))
}
