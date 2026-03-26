package define

import (
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

type UserClaim struct {
	Id       uint64
	Identity string
	Name     string
	Role     int
	jwt.StandardClaims
}

func getenv(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

var JwtKey = getenv("GOPAN_JWT_KEY", "change-me-in-env") // JWT密钥
var MailPassword = getenv("GOPAN_MAIL_PASSWORD", "")     // 邮箱授权码
var CodeLength = 6                                       // 验证码长度
var CodeExpire = 10 * 60                                 // 验证码过期时间，单位为秒
var FromMail = getenv("GOPAN_FROM_MAIL", "")             // 发件人邮箱
var SmtpHost = getenv("GOPAN_SMTP_HOST", "smtp.163.com") // SMTP主机地址
var SmtpPort = getenv("GOPAN_SMTP_PORT", "465")          // SMTP端口
var TencentSecretID = getenv("GOPAN_TENCENT_SECRET_ID", "")
var TencentSecretKey = getenv("GOPAN_TENCENT_SECRET_KEY", "")
var COSBucketURL = getenv("GOPAN_COS_BUCKET_URL", "") // 腾讯云COS桶URL
var PageSize = 20                                     // 分页默认每页数量
var DateFormat = "2006-01-02 15:04:05"                // 日期时间格式
var TokenExpire = 24 * 60 * 60                        // token过期时间，单位为秒
var RefreshTokenExpire = 7 * 24 * 60 * 60             // refresh token过期时间，单位为秒
