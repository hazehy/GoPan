package helper

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"gopan/gopan/define"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/smtp"
	"net/url"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/jordan-wright/email"
	uuid "github.com/satori/go.uuid"
	"github.com/tencentyun/cos-go-sdk-v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	usernameRegexp = regexp.MustCompile(`^[\p{Han}a-zA-Z0-9_]+$`)
	emailRegexp    = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)
)

func NormalizeInput(value string) string {
	return strings.TrimSpace(value)
}

func IsValidUsername(name string) bool {
	name = NormalizeInput(name)
	if len(name) < 2 || len(name) > 20 {
		return false
	}
	return usernameRegexp.MatchString(name)
}

func IsValidPassword(password string) bool {
	return len(password) >= 6 && len(password) <= 32
}

func IsValidEmail(email string) bool {
	email = NormalizeInput(email)
	return emailRegexp.MatchString(email)
}

func IsValidFileOrFolderName(name string) bool {
	name = NormalizeInput(name)
	if len(name) == 0 || len(name) > 100 {
		return false
	}
	if strings.ContainsAny(name, `\\/:*?"<>|`) {
		return false
	}
	return true
}

func IsValidPositiveDays(days int) bool {
	return days > 0 && days <= 3650
}

// Bcrypt 加密密码
func Bcrypt(s string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("bcrypt加密失败: %v", err)
		return ""
	}

	return string(hash)
}

// ComparePassword 比较加密后的密码和明文密码
func ComparePassword(hash string, plain string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
	if err != nil {
		return false
	}

	return true
}

// GenerateToken 生成JWT token
func GenerateToken(id uint64, identity string, name string, role int, second int) (string, error) {
	uc := define.UserClaim{
		Id:       id,
		Identity: identity,
		Name:     name,
		Role:     role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(second) * time.Second).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, uc)
	tokenString, err := token.SignedString([]byte(define.JwtKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// AnalyzeToken 解析JWT token
func AnalyzeToken(token string) (*define.UserClaim, error) {
	uc := new(define.UserClaim)
	claims, err := jwt.ParseWithClaims(token, uc, func(token *jwt.Token) (interface{}, error) {
		return []byte(define.JwtKey), nil
	})
	if err != nil {
		return nil, err
	}
	if !claims.Valid {
		return uc, errors.New("token无效")
	}
	return uc, nil
}

// MailCodeSend 发送邮件验证码
func MailCodeSend(mail, code string) error {
	if strings.TrimSpace(define.FromMail) == "" || strings.TrimSpace(define.MailPassword) == "" {
		return errors.New("邮件配置缺失，请设置 GOPAN_FROM_MAIL 和 GOPAN_MAIL_PASSWORD")
	}
	if strings.TrimSpace(define.SmtpHost) == "" || strings.TrimSpace(define.SmtpPort) == "" {
		return errors.New("SMTP配置缺失，请设置 GOPAN_SMTP_HOST 和 GOPAN_SMTP_PORT")
	}

	e := email.NewEmail()
	e.From = "GoPan <" + define.FromMail + ">"
	e.To = []string{strings.TrimSpace(mail)}
	e.Subject = "邮箱验证码"
	e.HTML = []byte("您的验证码是: <b>" + code + "</b>")
	err := e.SendWithTLS(define.SmtpHost+":"+define.SmtpPort, smtp.PlainAuth("", define.FromMail, define.MailPassword, define.SmtpHost),
		&tls.Config{InsecureSkipVerify: true, ServerName: define.SmtpHost})
	if err != nil {
		return err
	}
	return nil
}

// RandomCode 生成随机验证码
func RandomCode() string {
	s := "0123456789abcdefghijklmnopqrstuvwxyz"
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var code strings.Builder
	code.Grow(define.CodeLength)
	for i := 0; i < define.CodeLength; i++ {
		code.WriteByte(s[rng.Intn(len(s))])
	}
	return code.String()
}

// GenerateUUID 生成UUID
func GenerateUUID() string {
	return uuid.NewV4().String()
}

// COSUpLoad COS上传文件
func COSUpLoad(r *http.Request) (string, error) {
	client, err := buildCOSClient()
	if err != nil {
		return "", err
	}
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		return "", err
	}
	defer file.Close()
	key := "gopan/" + GenerateUUID() + path.Ext(fileHeader.Filename)

	_, err = client.Object.Put(
		context.Background(), key, file, nil,
	)
	if err != nil {
		return "", err
	}
	return define.COSBucketURL + "/" + key, nil
}

// CosChunkInit COS分片上传初始化
func CosChunkInit(ext string) (string, string, error) {
	client, err := buildCOSClient()
	if err != nil {
		return "", "", err
	}
	key := "gopan/" + GenerateUUID() + ext
	v, _, err := client.Object.InitiateMultipartUpload(
		context.Background(), key, nil,
	)
	if err != nil {
		return "", "", err
	}
	return key, v.UploadID, nil
}

// CosChunkUpload COS分片上传
func CosChunkUpload(r *http.Request) (string, error) {
	client, err := buildCOSClient()
	if err != nil {
		return "", err
	}
	uploadID := r.PostForm.Get("upload_id")
	key := r.PostForm.Get("key")
	partNumber, err := strconv.Atoi(r.PostForm.Get("part_number"))
	if err != nil {
		return "", err
	}
	f, _, err := r.FormFile("file")
	if err != nil {
		return "", err
	}
	defer f.Close()
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, f); err != nil {
		return "", err
	}
	resp, err := client.Object.UploadPart(
		context.Background(), key, uploadID, partNumber, bytes.NewReader(buf.Bytes()), nil,
	)
	if err != nil {
		return "", err
	}
	partETag := resp.Header.Get("ETag")
	return strings.Trim(partETag, "\""), nil
}

// CosChunkComplete COS分片上传完成
func CosChunkComplete(key, uploadID string, parts []cos.Object) error {
	client, err := buildCOSClient()
	if err != nil {
		return err
	}
	opt := &cos.CompleteMultipartUploadOptions{}
	opt.Parts = append(opt.Parts, parts...)
	_, _, err = client.Object.CompleteMultipartUpload(
		context.Background(), key, uploadID, opt,
	)
	if err != nil {
		return err
	}
	return nil
}

func buildCOSClient() (*cos.Client, error) {
	if err := validateCOSConfig(); err != nil {
		return nil, err
	}

	u, err := url.Parse(define.COSBucketURL)
	if err != nil {
		return nil, err
	}
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  define.TencentSecretID,
			SecretKey: define.TencentSecretKey,
		},
	})
	return client, nil
}

func validateCOSConfig() error {
	if strings.TrimSpace(define.COSBucketURL) == "" {
		return errors.New("COS配置缺失，请设置 GOPAN_COS_BUCKET_URL")
	}
	if strings.TrimSpace(define.TencentSecretID) == "" || strings.TrimSpace(define.TencentSecretKey) == "" {
		return errors.New("COS密钥缺失，请设置 GOPAN_TENCENT_SECRET_ID 和 GOPAN_TENCENT_SECRET_KEY")
	}
	return nil
}
