package router

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log"
	Dialogue "mianshi-helper/dialogue"
	"mianshi-helper/engine"
	InInitializer "mianshi-helper/initializer"
	"mianshi-helper/service"
	"mianshi-helper/tools"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Question struct {
	ContextStr string `json:"context"`
	SessionId  string `json:"sessionId"`
}

type LoginParam struct {
	Username         string `json:"username"`
	Password         string `json:"password"`
	VerificationCode string `json:"verificationCode"`
	SessionId        string `json:"sessionId"`
}

type UserNameRequest struct {
	UserName string `json:"userName" binding:"required"` // 使用json标签来指定字段名，并使用binding来确保字段是必需的
}

func SetupRouter() *gin.Engine {
	router := gin.Default()
	expiration := 1440 * time.Minute
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // 允许所有源
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,           // 允许携带凭证（cookies等）
		MaxAge:           12 * time.Hour, // 预检请求的有效期
	}))
	redisHelper := engine.GetRedisHelper()

	router.POST("/create", func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		sessionId := InInitializer.CreateDialogue()
		log.Println(tools.VerifyLogin(redisHelper, token))
		if tools.VerifyLogin(redisHelper, token) == "success" {
			c.JSON(200, gin.H{
				"sessionId": sessionId,
			})
			redisHelper.Set(context.Background(), "test-sessionId", sessionId, expiration)
		} else {
			c.JSON(403, gin.H{
				"message": "请重新登录",
			})
		}
	})

	router.POST("/answer", func(c *gin.Context) {
		var question Question
		token := c.GetHeader("Authorization")
		if tools.VerifyLogin(redisHelper, token) == "success" {
			// 解析 JSON 请求体到 user 结构体中
			err := json.NewDecoder(c.Request.Body).Decode(&question)
			redisHelper.Set(context.Background(), "test", question.SessionId, expiration)
			if err != nil {
				return
			}
			response := Dialogue.SendDialogueContent(question.ContextStr, question.SessionId)
			c.JSON(200, gin.H{
				"answer": response,
			})
		} else {
			c.JSON(403, gin.H{
				"message": "请重新登录",
			})
		}
	})

	router.GET("/getVerificationCode/:sessionId", func(c *gin.Context) {
		sessionId := c.Param("sessionId")
		captchaResult := tools.GenerateCaptcha(c)
		// 去除Base64字符串中可能的空格和换行符
		sanitizedBase64 := strings.ReplaceAll(captchaResult.Base64Blog, " ", "")
		sanitizedBase64 = strings.ReplaceAll(sanitizedBase64, "\n", "")
		sanitizedBase64 = strings.ReplaceAll(sanitizedBase64, "data:image/png;base64,", "")
		// 将Base64编码的图片数据解码为字节切片
		imageData, err := base64.StdEncoding.DecodeString(sanitizedBase64)
		if err != nil {
			c.String(500, "图片解码失败")
			log.Println(err)
			return
		}
		// 设置响应头为图片类型
		c.Header("Content-Type", "image/png")
		// 直接写入解码后的图片数据到响应中
		_, err = c.Writer.Write(imageData)
		if err != nil {
			c.JSON(500, gin.H{
				"message": "发送图片数据失败",
			})
			return
		}

		// 将验证码的验证值存储到Redis中
		redisHelper.Set(context.Background(), sessionId, captchaResult.VerifyValue, 1*time.Minute)
	})

	router.POST("/login", func(c *gin.Context) {
		var loginParam LoginParam
		err := json.NewDecoder(c.Request.Body).Decode(&loginParam)
		if err != nil {
			c.JSON(403, gin.H{
				"message": "未知错误!",
			})
			return
		}
		currentVerifyCode := redisHelper.Get(context.Background(), loginParam.SessionId).Val()
		if currentVerifyCode == "" {
			c.JSON(403, gin.H{
				"message": "验证码过期",
			})
		} else if currentVerifyCode != loginParam.VerificationCode {
			c.JSON(403, gin.H{
				"message": "验证码错误",
			})
		} else {
			redisHelper.Del(context.Background(), loginParam.Username)
			isRight := service.VerifyUser(loginParam.Username, loginParam.Password)
			if !isRight {
				c.JSON(403, gin.H{
					"message": "用户名或密码错误!",
				})
			} else {
				jwt, err := tools.GenToken(loginParam.Username)
				if err != nil {
					c.JSON(403, gin.H{
						"message": "登陆失败,请重试!",
					})
				}
				redisHelper.Set(context.Background(), jwt, loginParam.Username, expiration)
				c.JSON(200, gin.H{
					"message": "登陆成功!",
					"data": map[string]interface{}{
						"token": jwt,
						// 过期时间是当前时间+登陆有效时间
						"expirationTime": time.Now().Add(expiration).Unix(),
					},
				})
			}
		}
	})

	router.POST("/verifyUserNameIsExist", func(c *gin.Context) {
		var request UserNameRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			// 如果绑定失败，返回错误
			c.JSON(401, gin.H{})
			return
		}

		// 使用绑定后的数据
		userName := request.UserName
		isExist := service.CheckUserNameIsInDB(userName)
		log.Println(isExist)
		log.Println(userName)
		if !isExist {
			c.JSON(200, gin.H{
				"message": "用户名已存在",
				"data":    !isExist,
			})
		} else {
			c.JSON(200, gin.H{
				"message": "用户名可用",
				"data":    isExist,
			})
		}
	})

	router.GET("/verifyPhoneIsExist", func(c *gin.Context) {
		phone := c.Query("phone")
		isExist := service.CheckPhoneIsInDB(phone)
		if isExist {
			c.JSON(401, gin.H{
				"message": "手机号已存在",
				"data":    isExist,
			})
		} else {
			c.JSON(200, gin.H{
				"message": "手机号可用",
				"data":    isExist,
			})
		}
	})

	router.GET("/verifyEmailIsExist", func(c *gin.Context) {
		email := c.Query("email")
		isExist := service.CheckEmailIsInDB(email)
		if isExist {
			c.JSON(401, gin.H{
				"message": "邮箱已存在",
				"data":    isExist,
			})
		} else {
			c.JSON(200, gin.H{
				"message": "邮箱可用",
				"data":    isExist,
			})
		}
	})

	router.GET("/verifyAuth", func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		ctx.JSON(200, gin.H{
			"message": tools.VerifyLogin(redisHelper, token),
		})
	})

	router.POST("/register", func(c *gin.Context) {
		var user service.User
		err := json.NewDecoder(c.Request.Body).Decode(&user)
		if err != nil {
			log.Println(err)
			c.JSON(500, gin.H{
				"message": "未知错误!",
			})
		} else {
			service.CreateUser(user)
			c.JSON(200, gin.H{
				"message": "注册成功",
			})
		}
	})

	return router
}
