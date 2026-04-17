package model

import (
	"GoWorkingProject/Back/config"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 全局GORM DB对象（所有数据库操作都用这个）
var DB *gorm.DB

// User  用户数据模型（对应数据库users表，GORM自动映射）
type User struct {
	ID         uint64    `gorm:"primaryKey;autoIncrement" json:"user_id"`      // 主键，自增
	Username   string    `gorm:"size:20;uniqueIndex;not null" json:"username"` // 用户名，唯一索引，非空
	Password   string    `gorm:"size:18;not null" json:"-"`                    // 密码（加密存储），不返回给前端
	Email      string    `gorm:"size:20;uniqueIndex;not null" json:"email"`    // 邮箱，唯一索引，非空
	Phone      string    `gorm:"size:11;uniqueIndex;not null" json:"phone"`    // 电话，唯一索引，非空
	Avatar     string    `gorm:"size:255" json:"avatar"`                       // 头像URL
	Bio        string    `gorm:"size:500" json:"bio"`                          // 个人简介
	CreateTime time.Time `gorm:"autoCreateTime" json:"create_time"`            // 自动生成创建时间
	UpdateTime time.Time `gorm:"autoUpdateTime" json:"update_time"`            // 自动更新修改时间
}

// Post 帖子数据模型（对应posts表，后续完善功能用）
type Post struct {
	ID           uint64    `gorm:"primaryKey;autoIncrement" json:"post_id"`
	UserID       uint64    `gorm:"not null" json:"user_id"`               // 关联用户ID（外键）
	Title        string    `gorm:"size:30;not null" json:"title"`         // 帖子标题
	Content      string    `gorm:"type:longtext;not null" json:"content"` // 帖子内容
	LikeCount    int       `gorm:"default:0" json:"like_count"`           // 点赞数，默认0
	ReadCount    int       `gorm:"default:0" json:"read_count"`           // 阅读量，默认0
	CollectCount int       `gorm:"default:0" json:"collect_count"`        // 收藏量，默认0
	CreateTime   time.Time `gorm:"autoCreateTime" json:"create_time"`
	UpdateTime   time.Time `gorm:"autoUpdateTime" json:"update_time"`
}

// Comment 评论数据模型（对应comments表）
type Comment struct {
	ID         uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	PostID     uint64    `gorm:"not null" json:"post_id"`           // 关联帖子ID（外键）
	UserID     uint64    `gorm:"not null" json:"user_id"`           // 关联用户ID（外键）
	Content    string    `gorm:"type:text;not null" json:"content"` // 评论内容
	LikeCount  int       `gorm:"default:0" json:"like_count"`       // 点赞数，默认0
	CreateTime time.Time `gorm:"autoCreateTime" json:"create_time"`
	UpdateTime time.Time `gorm:"autoUpdateTime" json:"update_time"`
}

// Like 点赞数据模型（对应likes表）
type Like struct {
	ID         uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	PostID     uint64    `gorm:"not null;uniqueIndex:idx_user_target" json:"post_id"`     // 关联帖子ID
	CommentID  uint64    `gorm:"default:0;uniqueIndex:idx_user_target" json:"comment_id"` // 关联评论ID，0表示帖子点赞
	UserID     uint64    `gorm:"not null;uniqueIndex:idx_user_target" json:"user_id"`     // 关联用户ID
	CreateTime time.Time `gorm:"autoCreateTime" json:"create_time"`
}

// Collect 收藏数据模型（对应collects表）
type Collect struct {
	ID         uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	PostID     uint64    `gorm:"not null;uniqueIndex:idx_user_post" json:"post_id"` // 关联帖子ID，与user_id联合唯一
	UserID     uint64    `gorm:"not null;uniqueIndex:idx_user_post" json:"user_id"` // 关联用户ID，与post_id联合唯一
	CreateTime time.Time `gorm:"autoCreateTime" json:"create_time"`
}

// Follow 关注数据模型（对应follows表）
type Follow struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	FollowerID  uint64    `gorm:"not null;uniqueIndex:uk_follower_following" json:"follower_id"`  // 关注者ID
	FollowingID uint64    `gorm:"not null;uniqueIndex:uk_follower_following" json:"following_id"` // 被关注者ID
	CreateTime  time.Time `gorm:"autoCreateTime" json:"create_time"`
}

// Message 聊天消息数据模型（对应messages表）
type Message struct {
	ID         uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	SenderID   uint64    `gorm:"not null" json:"sender_id"`              // 发送者ID
	ReceiverID uint64    `gorm:"not null" json:"receiver_id"`            // 接收者ID
	Content    string    `gorm:"type:text;not null" json:"content"`      // 消息内容
	Type       string    `gorm:"size:20;default:'text'" json:"type"`     // 消息类型：text, image, emoji
	Status     string    `gorm:"size:20;default:'unread'" json:"status"` // 消息状态：unread, read
	CreateTime time.Time `gorm:"autoCreateTime" json:"create_time"`
}

// InitMySQL 初始化GORM连接MySQL（只执行一次）
func InitMySQL() error {
	// 拼接GORM连接DSN（格式：用户名:密码@tcp(地址:端口)/数据库名?参数）
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		config.MySQLConfig.Username,
		config.MySQLConfig.Password,
		config.MySQLConfig.Host,
		config.MySQLConfig.Port,
		config.MySQLConfig.Database,
		config.MySQLConfig.Charset,
	)

	// 配置GORM：开启SQL日志（开发时调试用），其他默认
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // 打印执行的SQL语句
	})
	if err != nil {
		return fmt.Errorf("GORM连接MySQL失败: %v", err)
	}

	// 获取底层SQL连接池，设置连接池参数（优化性能）
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取MySQL连接池失败: %v", err)
	}
	sqlDB.SetMaxOpenConns(100)  // 最大打开连接数
	sqlDB.SetMaxIdleConns(20)   // 最大空闲连接数
	sqlDB.SetConnMaxLifetime(0) // 连接存活时间（0表示不限制）

	// 自动迁移表：如果表不存在则创建，存在则不修改结构（安全）
	// 同时创建users、posts、comments、likes、collects、follows和messages七张表
	err = db.AutoMigrate(&User{}, &Post{}, &Comment{}, &Like{}, &Collect{}, &Follow{}, &Message{})
	if err != nil {
		return fmt.Errorf("自动迁移数据库表失败: %v", err)
	}

	// 赋值给全局DB对象，供其他模块使用
	DB = db

	fmt.Println("✅ GORM连接MySQL成功，已自动创建/迁移users、posts表")
	return nil
}
