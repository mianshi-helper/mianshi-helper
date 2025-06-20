package tools

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
)

type CaptchaResult struct {
	Id          string `json:"id"`
	Base64Blog  string `json:"base_64_blog"`
	VerifyValue string `json:"code"`
}

// 默认存储10240个验证码，每个验证码10分钟过期
var store = base64Captcha.DefaultMemStore

// 生成图片验证码
func GenerateCaptcha(context *gin.Context) CaptchaResult {
	// 生成默认数字
	//driver := base64Captcha.DefaultDriverDigit
	driver := base64Captcha.NewDriverDigit(70, 130, 4, 0.8, 100)
	// 生成base64图片
	captcha := base64Captcha.NewCaptcha(driver, store)
	// 获取
	id, b64s, verifyValue, err := captcha.Generate()
	if err != nil {
		fmt.Println("Register GetCaptchaPhoto get base64Captcha has err:", err)
	}

	captchaResult := CaptchaResult{Id: id, Base64Blog: b64s, VerifyValue: verifyValue}
	return captchaResult
}

// 校验图片验证码,并清除内存空间
func VerifyCaptcha(id string, value string) bool {
	// TODO 只要id存在，就会校验并清除，无论校验的值是否成功, 所以同一id只能校验一次
	// 注意：id,b64s是空 也会返回true 需要在加判断
	verifyResult := store.Verify(id, value, true)
	return verifyResult
}
