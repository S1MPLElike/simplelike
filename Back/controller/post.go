package controller

import (
	"GoWorkingProject/Back/model"
	"GoWorkingProject/Back/utils"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PostCreateRequest struct {
	Title   string `json:"title" binding:"required,min=3,max=100"`
	Content string `json:"content" binding:"required,min=10"`
}

func PostCreate(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req PostCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, utils.Fail(utils.CodeParamError, "参数错误: "+err.Error()))
		return
	}

	post := &model.Post{
		UserID:       userID.(uint64),
		Title:        req.Title,
		Content:      req.Content,
		LikeCount:    0,
		ReadCount:    0,
		CollectCount: 0,
	}

	if err := model.DB.Create(post).Error; err != nil {
		c.JSON(200, utils.Fail(utils.CodeDBError, "发布帖子失败: "+err.Error()))
		return
	}

	// 清除前5页按时间排序的缓存
	for i := 1; i <= 5; i++ {
		model.RedisDel(fmt.Sprintf("posts:page:time:%d", i))
	}

	// 清除默认排序（空字符串）的前5页缓存
	for i := 1; i <= 5; i++ {
		model.RedisDel(fmt.Sprintf("posts:page::%d", i))
	}

	c.JSON(200, utils.Success(gin.H{
		"post_id": post.ID,
		"title":   post.Title,
		"content": post.Content,
	}))
}

func GetHotPosts(c *gin.Context) {
	sort := c.Query("sort")
	pageStr := c.DefaultQuery("page", "1")
	page := 1
	fmt.Sscanf(pageStr, "%d", &page)
	if page < 1 {
		page = 1
	}

	pageSize := 5
	offset := (page - 1) * pageSize

	var orderBy string
	switch sort {
	case "likes":
		orderBy = "like_count desc"
	case "views":
		orderBy = "read_count desc"
	default:
		orderBy = "create_time desc"
	}

	// 尝试从Redis缓存获取
	cacheKey := fmt.Sprintf("posts:page:%s:%d", sort, page)
	cachedData, err := model.RedisGet(cacheKey)
	if err == nil {
		// 缓存命中
		var postsWithUser []struct {
			model.Post
			Username string `json:"username"`
		}
		if json.Unmarshal([]byte(cachedData), &postsWithUser) == nil {
			c.JSON(200, utils.Success(postsWithUser))
			return
		}
	}

	// 缓存未命中，尝试获取分布式锁
	lockKey := fmt.Sprintf("lock:posts:page:%s:%d", sort, page)
	lockSuccess, _ := model.RedisSetNX(lockKey, "1", 5*time.Second) // 5秒过期

	if lockSuccess {
		defer model.RedisDel(lockKey) // 确保锁释放

		var posts []model.Post
		if err := model.DB.Order(orderBy).Limit(pageSize).Offset(offset).Find(&posts).Error; err != nil {
			c.JSON(200, utils.Fail(utils.CodeDBError, "获取帖子失败"))
			return
		}

		// 为每个帖子添加用户名
		type PostWithUser struct {
			model.Post
			Username string `json:"username"`
		}

		var postsWithUser []PostWithUser
		for _, post := range posts {
			var user model.User
			model.DB.Select("username").Where("id = ?", post.UserID).First(&user)

			postsWithUser = append(postsWithUser, PostWithUser{
				Post:     post,
				Username: user.Username,
			})
		}

		// 缓存结果，5-10分钟随机过期
		if data, err := json.Marshal(postsWithUser); err == nil {
			model.RedisSet(cacheKey, string(data), getRandomExpiration())
		}

		c.JSON(200, utils.Success(postsWithUser))
	} else {
		// 未获得锁，等待后重试
		time.Sleep(100 * time.Millisecond)
		// 重新尝试获取缓存
		GetHotPosts(c)
	}
}

func GetPostDetail(c *gin.Context) {
	postID := c.Param("id")

	// 尝试从Redis缓存获取
	cacheKey := fmt.Sprintf("post:detail:%s", postID)
	cachedData, err := model.RedisGet(cacheKey)
	if err == nil {
		// 缓存命中
		var post model.Post
		if json.Unmarshal([]byte(cachedData), &post) == nil {
			// 异步增加阅读量，不影响返回速度
			go func() {
				model.DB.Model(&model.Post{}).Where("id = ?", postID).Update("read_count", gorm.Expr("read_count + ?", 1))
				// 更新缓存中的阅读量
				post.ReadCount++
				if data, err := json.Marshal(post); err == nil {
					model.RedisSet(cacheKey, string(data), getRandomExpiration())
				}
			}()
			c.JSON(200, utils.Success(post))
			return
		}
	}

	// 缓存未命中，尝试获取分布式锁
	lockKey := fmt.Sprintf("lock:post:detail:%s", postID)
	lockSuccess, _ := model.RedisSetNX(lockKey, "1", 5*time.Second) // 5秒过期

	if lockSuccess {
		defer model.RedisDel(lockKey) // 确保锁释放

		var post model.Post
		if err := model.DB.First(&post, postID).Error; err != nil {
			c.JSON(200, utils.Fail(utils.CodeDBError, "获取帖子详情失败"))
			return
		}

		// 增加阅读量
		model.DB.Model(&post).Update("read_count", post.ReadCount+1)
		post.ReadCount++

		// 缓存结果，5-10分钟随机过期
		if data, err := json.Marshal(post); err == nil {
			model.RedisSet(cacheKey, string(data), getRandomExpiration())
		}

		c.JSON(200, utils.Success(post))
	} else {
		// 未获得锁，等待后重试
		time.Sleep(100 * time.Millisecond)
		GetPostDetail(c)
	}
}

// GetComments 获取帖子评论列表
func GetComments(c *gin.Context) {
	postID := c.Param("id")

	// 尝试从Redis缓存获取
	cacheKey := fmt.Sprintf("post:comments:%s", postID)
	cachedData, err := model.RedisGet(cacheKey)
	if err == nil {
		// 缓存命中
		var commentsWithUser []struct {
			model.Comment
			Username string `json:"username"`
		}
		if json.Unmarshal([]byte(cachedData), &commentsWithUser) == nil {
			c.JSON(200, utils.Success(commentsWithUser))
			return
		}
	}

	var comments []model.Comment
	if err := model.DB.Where("post_id = ?", postID).Order("create_time desc").Find(&comments).Error; err != nil {
		c.JSON(200, utils.Fail(utils.CodeDBError, "获取评论列表失败"))
		return
	}

	// 为每个评论添加用户名和头像
	type CommentWithUser struct {
		model.Comment
		Username string `json:"username"`
		Avatar   string `json:"avatar"`
		IsLiked  bool   `json:"is_liked"`
	}

	var commentsWithUser []CommentWithUser
	userID, _ := c.Get("user_id") // 获取当前用户ID

	for _, comment := range comments {
		var user model.User
		model.DB.Select("username, avatar").Where("id = ?", comment.UserID).First(&user)

		// 检查当前用户是否点赞过该评论
		isLiked := false
		if userID != nil {
			var like model.Like
			result := model.DB.Where("comment_id = ? AND user_id = ?", comment.ID, userID).First(&like)
			isLiked = result.Error == nil
		}

		commentsWithUser = append(commentsWithUser, CommentWithUser{
			Comment:  comment,
			Username: user.Username,
			Avatar:   user.Avatar,
			IsLiked:  isLiked,
		})
	}

	// 缓存结果，3-5分钟随机过期
	if data, err := json.Marshal(commentsWithUser); err == nil {
		// 3-5分钟随机过期
		rand.Seed(time.Now().UnixNano())
		randomSeconds := 180 + rand.Intn(121) // 180秒(3分钟)到300秒(5分钟)
		model.RedisSet(cacheKey, string(data), time.Duration(randomSeconds)*time.Second)
	}

	c.JSON(200, utils.Success(commentsWithUser))
}

// CommentCreate 创建评论
type CommentCreateRequest struct {
	PostID  uint64 `json:"post_id" binding:"required"`
	Content string `json:"content" binding:"required,min=1,max=1000"`
}

func CommentCreate(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req CommentCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, utils.Fail(utils.CodeParamError, "参数错误: "+err.Error()))
		return
	}

	comment := &model.Comment{
		PostID:    req.PostID,
		UserID:    userID.(uint64),
		Content:   req.Content,
		LikeCount: 0,
	}

	if err := model.DB.Create(comment).Error; err != nil {
		c.JSON(200, utils.Fail(utils.CodeDBError, "创建评论失败: "+err.Error()))
		return
	}

	// 清除对应帖子的评论缓存
	model.RedisDel(fmt.Sprintf("post:comments:%d", req.PostID))

	c.JSON(200, utils.Success(gin.H{
		"comment_id": comment.ID,
		"content":    comment.Content,
	}))
}

func LikePost(c *gin.Context) {
	userID, _ := c.Get("user_id")
	postID := c.Param("id")

	// 检查帖子是否存在
	var post model.Post
	if err := model.DB.First(&post, postID).Error; err != nil {
		c.JSON(200, utils.Fail(utils.CodeDBError, "帖子不存在"))
		return
	}

	// 检查是否已经点赞过
	var like model.Like
	result := model.DB.Where("post_id = ? AND comment_id = 0 AND user_id = ?", post.ID, userID.(uint64)).First(&like)

	if result.Error == nil {
		// 已经点赞过，取消点赞
		model.DB.Model(&post).UpdateColumn("like_count", gorm.Expr("like_count - ?", 1))
		model.DB.Delete(&like)

		// 只清除帖子详情缓存，不清除分页列表缓存
		// 原因：点赞/收藏操作对分页列表的影响较小，为了性能考虑，不清除分页缓存
		model.RedisDel(fmt.Sprintf("post:detail:%s", postID))

		c.JSON(200, utils.Success(gin.H{
			"post_id":    post.ID,
			"like_count": post.LikeCount - 1,
			"liked":      false,
		}))
	} else {
		// 未点赞，添加点赞
		model.DB.Model(&post).UpdateColumn("like_count", gorm.Expr("like_count + ?", 1))

		like = model.Like{
			PostID:    post.ID,
			CommentID: 0,
			UserID:    userID.(uint64),
		}
		model.DB.Create(&like)

		// 只清除帖子详情缓存，不清除分页列表缓存
		// 原因：点赞/收藏操作对分页列表的影响较小，为了性能考虑，不清除分页缓存
		model.RedisDel(fmt.Sprintf("post:detail:%s", postID))

		c.JSON(200, utils.Success(gin.H{
			"post_id":    post.ID,
			"like_count": post.LikeCount + 1,
			"liked":      true,
		}))
	}
}

func CollectPost(c *gin.Context) {
	userID, _ := c.Get("user_id")
	postID := c.Param("id")

	// 检查帖子是否存在
	var post model.Post
	if err := model.DB.First(&post, postID).Error; err != nil {
		c.JSON(200, utils.Fail(utils.CodeDBError, "帖子不存在"))
		return
	}

	// 检查是否已经收藏过
	var collect model.Collect
	result := model.DB.Where("post_id = ? AND user_id = ?", post.ID, userID.(uint64)).First(&collect)

	if result.Error == nil {
		// 已经收藏过，取消收藏
		model.DB.Model(&post).UpdateColumn("collect_count", gorm.Expr("collect_count - ?", 1))
		model.DB.Delete(&collect)

		// 只清除帖子详情缓存，不清除分页列表缓存
		// 原因：点赞/收藏操作对分页列表的影响较小，为了性能考虑，不清除分页缓存
		model.RedisDel(fmt.Sprintf("post:detail:%s", postID))

		c.JSON(200, utils.Success(gin.H{
			"post_id":       post.ID,
			"collect_count": post.CollectCount - 1,
			"collected":     false,
		}))
	} else {
		// 未收藏，添加收藏
		model.DB.Model(&post).UpdateColumn("collect_count", gorm.Expr("collect_count + ?", 1))

		collect = model.Collect{
			PostID: post.ID,
			UserID: userID.(uint64),
		}
		model.DB.Create(&collect)

		// 只清除帖子详情缓存，不清除分页列表缓存
		// 原因：点赞/收藏操作对分页列表的影响较小，为了性能考虑，不清除分页缓存
		model.RedisDel(fmt.Sprintf("post:detail:%s", postID))

		c.JSON(200, utils.Success(gin.H{
			"post_id":       post.ID,
			"collect_count": post.CollectCount + 1,
			"collected":     true,
		}))
	}
}

func SearchPosts(c *gin.Context) {
	keyword := c.Query("keyword")
	sort := c.Query("sort")
	var orderBy string

	switch sort {
	case "likes":
		orderBy = "like_count desc"
	case "views":
		orderBy = "read_count desc"
	default:
		orderBy = "create_time desc"
	}

	// 尝试从Redis缓存获取
	cacheKey := fmt.Sprintf("posts:search:%s:%s", keyword, sort)
	cachedData, err := model.RedisGet(cacheKey)
	if err == nil {
		// 缓存命中
		var postsWithUser []struct {
			model.Post
			Username string `json:"username"`
		}
		if json.Unmarshal([]byte(cachedData), &postsWithUser) == nil {
			c.JSON(200, utils.Success(postsWithUser))
			return
		}
	}

	var posts []model.Post
	query := model.DB

	if keyword != "" {
		query = query.Where("title LIKE ?", "%"+keyword+"%")
	}

	if err := query.Order(orderBy).Limit(20).Find(&posts).Error; err != nil {
		c.JSON(200, utils.Fail(utils.CodeDBError, "搜索帖子失败"))
		return
	}

	// 为每个帖子添加用户名
	type PostWithUser struct {
		model.Post
		Username string `json:"username"`
	}

	var postsWithUser []PostWithUser
	for _, post := range posts {
		var user model.User
		model.DB.Select("username").Where("id = ?", post.UserID).First(&user)

		postsWithUser = append(postsWithUser, PostWithUser{
			Post:     post,
			Username: user.Username,
		})
	}

	// 缓存结果，5-10分钟随机过期
	if data, err := json.Marshal(postsWithUser); err == nil {
		model.RedisSet(cacheKey, string(data), getRandomExpiration())
	}

	c.JSON(200, utils.Success(postsWithUser))
}

// 生成5-10分钟之间的随机过期时间
func getRandomExpiration() time.Duration {
	// 初始化随机数生成器
	rand.Seed(time.Now().UnixNano())
	// 生成5-10分钟之间的随机秒数
	randomSeconds := 300 + rand.Intn(301) // 300秒(5分钟)到600秒(10分钟)
	return time.Duration(randomSeconds) * time.Second
}

// 清除帖子相关缓存
func clearPostCache(postID string) {
	// 清除帖子详情缓存
	model.RedisDel(fmt.Sprintf("post:detail:%s", postID))

	// 清除第一页的缓存（所有排序方式）
	sorts := []string{"", "likes", "views", "time"}
	for _, sort := range sorts {
		model.RedisDel(fmt.Sprintf("posts:page:%s:1", sort))
	}

	// 清除搜索缓存（所有可能的关键词和排序方式）
	// 注意：这里简化处理，实际项目中可能需要更精细的缓存管理
	// 或者使用Redis的键模式匹配来删除相关缓存
}

// 检查用户是否已点赞
func CheckLikeStatus(c *gin.Context) {
	userID, _ := c.Get("user_id")
	postID := c.Param("id")

	var like model.Like
	result := model.DB.Where("post_id = ? AND user_id = ?", postID, userID.(uint64)).First(&like)

	c.JSON(200, utils.Success(gin.H{
		"liked": result.Error == nil,
	}))
}

// 检查用户是否已收藏
func CheckCollectStatus(c *gin.Context) {
	userID, _ := c.Get("user_id")
	postID := c.Param("id")

	var collect model.Collect
	result := model.DB.Where("post_id = ? AND user_id = ?", postID, userID.(uint64)).First(&collect)

	c.JSON(200, utils.Success(gin.H{
		"collected": result.Error == nil,
	}))
}

// 评论点赞功能（无上限）
func LikeComment(c *gin.Context) {
	commentID := c.Param("id")

	var comment model.Comment
	if err := model.DB.First(&comment, commentID).Error; err != nil {
		c.JSON(200, utils.Fail(utils.CodeDBError, "评论不存在"))
		return
	}

	// 直接增加点赞数
	model.DB.Model(&comment).UpdateColumn("like_count", gorm.Expr("like_count + ?", 1))

	// 清除评论缓存
	model.RedisDel(fmt.Sprintf("post:comments:%d", comment.PostID))

	c.JSON(200, utils.Success(gin.H{
		"comment_id": comment.ID,
		"like_count": comment.LikeCount + 1,
	}))
}

// GetUserPosts 获取用户的帖子列表
func GetUserPosts(c *gin.Context) {
	userIDStr := c.Param("id")
	var userID uint64
	_, err := fmt.Sscanf(userIDStr, "%d", &userID)
	if err != nil {
		c.JSON(200, utils.Fail(utils.CodeParamError, "无效的用户ID"))
		return
	}

	// 获取分页参数
	pageStr := c.DefaultQuery("page", "1")
	var page int
	_, err = fmt.Sscanf(pageStr, "%d", &page)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize := 5
	offset := (page - 1) * pageSize

	var posts []model.Post
	if err := model.DB.Where("user_id = ?", userID).Order("create_time desc").Limit(pageSize).Offset(offset).Find(&posts).Error; err != nil {
		c.JSON(200, utils.Fail(utils.CodeDBError, "获取用户帖子失败"))
		return
	}

	c.JSON(200, utils.Success(posts))
}

// GetCollectedPosts 获取用户收藏的帖子列表
func GetCollectedPosts(c *gin.Context) {
	userIDStr := c.Param("id")
	var userID uint64
	_, err := fmt.Sscanf(userIDStr, "%d", &userID)
	if err != nil {
		c.JSON(200, utils.Fail(utils.CodeParamError, "无效的用户ID"))
		return
	}

	// 获取分页参数
	pageStr := c.DefaultQuery("page", "1")
	var page int
	_, err = fmt.Sscanf(pageStr, "%d", &page)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize := 5
	offset := (page - 1) * pageSize

	var collectedPosts []model.Post
	if err := model.DB.Table("posts").Joins("JOIN collects ON posts.id = collects.post_id").Where("collects.user_id = ?", userID).Order("collects.create_time desc").Limit(pageSize).Offset(offset).Find(&collectedPosts).Error; err != nil {
		c.JSON(200, utils.Fail(utils.CodeDBError, "获取收藏帖子失败"))
		return
	}

	c.JSON(200, utils.Success(collectedPosts))
}
