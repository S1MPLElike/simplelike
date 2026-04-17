package utils

import (
	"GoWorkingProject/Back/config"
	"GoWorkingProject/Back/model"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GenerateToken 生成随机唯一Token（32位16进制字符串，足够安全）
func GenerateToken() string {
	// 生成16字节随机数（转换为16进制后是32位）
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// SaveTokenToRedis 将Token存入Redis（key: token:xxx，value: userID，带过期时间）
// 参数：生成的Token，用户ID
// 返回：错误信息（nil表示成功）
func SaveTokenToRedis(token string, userID uint64) error {
	// Redis key格式：token:xxx（区分其他key，避免冲突）
	key := fmt.Sprintf("token:%s", token)
	// 存入Redis，设置过期时间（和Cookie过期时间一致

	err := model.RedisSet(key, strconv.FormatUint(userID, 10), config.TokenConfig.ExpireTime)
	if err != nil {
		return fmt.Errorf("Token存入Redis失败: %v", err)
	}
	return nil
}

// VerifyToken 验证Token有效性（从Redis查询，判断是否存在/过期）
// 参数：前端Cookie中的Token
// 返回：用户ID（Token有效），错误信息（Token无效/过期）
func VerifyToken(token string) (uint64, error) {
	// Token为空，直接返回无效
	if token == "" {
		return 0, fmt.Errorf("Token为空，请先登录")
	}

	// 拼接Redis key
	key := fmt.Sprintf("token:%s", token)
	// 从Redis获取用户ID（string类型）

	//有时候代码写对但是仍然报错的原因 ：
	//你中午改代码时，GoLand 的缓存 / 索引还 “记住” 了旧的正确状态（比如你当时 model.RedisGet 是存在的）；
	//晚上改了目录名（util→utils）、改了函数名大小写（redisGet→RedisGet）后，GoLand 不会 “实时” 更新索引 / 缓存（怕频繁更新拖慢速度），
	//仍然用旧的缓存去匹配新代码 —— 结果就是 “明明代码改对了，IDE 却认为找不到函数”；

	userIDStr, err := model.RedisGet(key)
	if err != nil {
		// Redis返回错误（一般是Token不存在/过期）
		return 0, fmt.Errorf("Token无效或已过期")
	}

	// 将string类型的userID转换为uint64
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("用户ID解析失败")
	}

	// 可选优化：刷新Token过期时间（用户活跃时延长有效期）
	err = model.RedisDel(key)
	if err != nil {
		return 0, fmt.Errorf("重设置token过期时间失败")
	}
	err = model.RedisSet(key, strconv.FormatUint(userID, 10), config.TokenConfig.ExpireTime)
	if err != nil {
		return 0, fmt.Errorf("重设置token过期时间失败")
	}

	return userID, nil
}

// DeleteToken 从Redis删除Token（退出登录时用）
func DeleteToken(token string) error {
	key := fmt.Sprintf("token:%s", token)
	return model.RedisDel(key)
}

// SetTokenCookie 将Token写入前端Cookie（HttpOnly，防止XSS攻击）
// 参数：Gin上下文，生成的Token
func SetTokenCookie(c *gin.Context, token string) {
	c.SetCookie(
		config.TokenConfig.CookieName,                // Cookie名称（全局配置）
		token,                                        // Cookie值（Token）
		int(config.TokenConfig.ExpireTime.Seconds()), // 过期时间（秒）
		"/",                                          // 作用路径（全站有效）
		"",                                           // 作用域名（本地开发留空）
		false,                                        // 是否只在HTTPS传输（本地HTTP设为false）
		true,                                         // 是否禁止JS读取（HttpOnly，防XSS）
	)
}

// ClearTokenCookie 清除前端Cookie（退出登录时用）
func ClearTokenCookie(c *gin.Context) {
	c.SetCookie(
		config.TokenConfig.CookieName,
		"", // 空值，清除Cookie
		-1, // 过期时间设为-1，立即删除
		"/",
		"",
		false,
		true,
	)
}
