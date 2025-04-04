package controller

import (
	"book-management/configs"
	"book-management/internal/pkg/errcode"
	"book-management/internal/pkg/req"
	"book-management/internal/pkg/resp"
	"book-management/pkg/logger"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"io"
	"time"
)

const Key = "qazwsxedcrfvtgbyhnujmikolp123456"

type AuthCtrl struct {
	rdb *redis.Client
	up  map[string]string //user-password
}

func NewAuthCtrl(rdb *redis.Client, conf configs.AppConfig) *AuthCtrl {
	mp := make(map[string]string, len(conf.Users))
	for _, user := range conf.Users {
		mp[user.UserName] = user.Password
	}

	return &AuthCtrl{
		rdb: rdb,
		up:  mp,
	}
}

func (a *AuthCtrl) RegisterRoute(r *gin.Engine) {
	g := r.Group("/api/v1/auth")
	{
		g.GET("/get_verification_code", a.GenerateCaptcha)
		g.POST("/login", a.Login)
	}
}

// GenerateCaptcha 生成验证码
// @Summary 生成验证码
// @Description 生成验证码
// @Tags 用户
// @Accept json
// @Produce json
// @Param object query GenerateCaptchaReq true "请求参数"
// @Success 200 {object} GenerateCaptchaResp "成功"
func (a *AuthCtrl) GenerateCaptcha(c *gin.Context) {
	//这个不需要解析req
	// 生成验证码
	var captchaBuffer bytes.Buffer
	// 生成唯一Token（与验证码值分离）
	token := uuid.New().String()
	captchaId := captcha.New() // 生成验证码ID
	// 获取验证码的真实值（数字或字母）
	err := captcha.WriteImage(&captchaBuffer, captchaId, 120, 80) // 生成图片
	if err != nil {
		logger.LogPrinter.Errorf("captcha write image failed:%v", err)
		resp.SendResp(c, resp.NewRespFromErr(errcode.GenerateVerifyCodeError))
		return
	}

	// 将图片转换为Base64编码
	captchaBase64 := base64.StdEncoding.EncodeToString(captchaBuffer.Bytes())
	// 存储：Key=Token, Value=验证码值
	err = a.rdb.Set(c, "captcha:"+token, captchaId, 5*time.Minute).Err()
	if err != nil {
		logger.LogPrinter.Errorf("cache set verifycode failed:%v", err)
		resp.SendResp(c, resp.NewRespFromErr(errcode.GenerateVerifyCodeError))
		return
	}
	resp.SendResp(c, resp.WithData(resp.SuccessResp, map[string]interface{}{
		"image":                captchaBase64,
		"verification_code_id": token,
	}))
}

// Login 登录
// @Summary 登录
// @Description 登录
// @Tags 用户
// @Accept json
// @Produce json
// @Param object query LoginReq true "请求参数"
// @Success 200 {object} LoginResp "成功"
func (a *AuthCtrl) Login(c *gin.Context) {

	var loginReq LoginReq
	if err := req.ParseRequestBody(c, &loginReq); err != nil {
		resp.SendResp(c, resp.NewRespFromErr(err))
		return
	}

	var (
		token string
		err   error
		val   string
	)

	if _, ok := a.up[loginReq.UserID]; !ok {
		resp.SendResp(c, resp.NewRespFromErr(errcode.UserNotFoundError))
		goto DelCode
	}
	if a.up[loginReq.UserID] != loginReq.Password {
		resp.SendResp(c, resp.NewRespFromErr(errcode.PasswordError))
		goto DelCode
	}
	val, err = a.rdb.Get(c, "captcha:"+loginReq.VerificationCodeID).Result()
	if err != nil {
		resp.SendResp(c, resp.NewRespFromErr(errcode.VerificationCodeError))
		goto DelCode
	}
	if !captcha.VerifyString(val, loginReq.VerificationCode) {
		resp.SendResp(c, resp.NewRespFromErr(errcode.VerificationCodeError))
		goto DelCode
	}

	//生成token
	token, err = encryptString(loginReq.UserID, Key)
	if err != nil {
		resp.SendResp(c, resp.NewRespFromErr(errcode.GenerateTokenError))
		return
	}
	resp.SendResp(c, resp.WithData(resp.SuccessResp, map[string]interface{}{
		"token": token,
	}))

DelCode:
	a.rdb.Del(c, "captcha:"+loginReq.VerificationCodeID)
}

func (a *AuthCtrl) VerifyToken(token string) bool {
	if token == "" {
		return false
	}
	decryptedToken, err := decryptString(token, Key)
	if err != nil {
		return false
	}
	if _, ok := a.up[decryptedToken]; !ok {
		return false
	}
	return true
}

// 加密字符串
func encryptString(plainText string, key string) (string, error) {
	// 将密钥转换为32字节
	keyBytes := []byte(key)
	if len(keyBytes) < 32 {
		// 如果密钥不足32字节，用0填充
		paddedKey := make([]byte, 32)
		copy(paddedKey, keyBytes)
		keyBytes = paddedKey
	} else if len(keyBytes) > 32 {
		// 如果密钥超过32字节，截取前32字节
		keyBytes = keyBytes[:32]
	}

	// 创建加密块
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	// 将字符串转换为字节切片
	plainTextBytes := []byte(plainText)

	// 创建GCM模式的加密器
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// 创建随机nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// 加密数据
	cipherText := gcm.Seal(nonce, nonce, plainTextBytes, nil)

	// 返回Base64编码的加密结果
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// 解密字符串
func decryptString(encryptedText string, key string) (string, error) {
	// 将密钥转换为32字节
	keyBytes := []byte(key)
	if len(keyBytes) < 32 {
		paddedKey := make([]byte, 32)
		copy(paddedKey, keyBytes)
		keyBytes = paddedKey
	} else if len(keyBytes) > 32 {
		keyBytes = keyBytes[:32]
	}

	// 解码Base64字符串
	cipherText, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		return "", err
	}

	// 创建加密块
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	// 创建GCM模式的解密器
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// 检查密文长度是否足够
	if len(cipherText) < gcm.NonceSize() {
		return "", errors.New("密文太短")
	}

	// 分离nonce和实际密文
	nonce, cipherText := cipherText[:gcm.NonceSize()], cipherText[gcm.NonceSize():]

	// 解密数据
	plainTextBytes, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", err
	}

	// 返回解密后的字符串
	return string(plainTextBytes), nil
}
