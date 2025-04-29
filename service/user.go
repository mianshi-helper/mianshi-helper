package service

import (
	"database/sql"
	"mianshi-helper/db"
)

type User struct {
	ID             int           `db:"id"`       // 主键，自增ID
	Username       string        `db:"username"` // 用户名，不允许为空
	Phone          string        `db:"phone"`    // 手机号，使用sql.NullString表示可能为空
	Email          string        `db:"email"`    // 邮箱，使用sql.NullString表示可能为空
	Password       string        `db:"password"` // 密码，不允许为空（实际应用中应存储加密后的密码）
	AccountBalance float64       `db:"account_balance"`
	VipLevel       sql.NullInt64 `db:"vip_level"`  // VIP等级，使用sql.NullInt64表示可能为空
	Age            string        `db:"age"`        // 年龄，使用sql.NullInt64表示可能为空
	ResumeURL      string        `db:"resume_url"` // 简历地址（假设为URL链接），使用sql.NullString表示可能为空
}

var dataBase = db.ConnectDB()

func VerifyUser(userName string, password string) bool {
	var user User
	result := dataBase.Where("userName = ? AND password = ?", userName, password).Take(&user)
	if result.Error != nil {
		return false
	}
	if result.RowsAffected == 0 {
		return false
	}
	return true
}

func CreateUser(user User) bool {
	result := dataBase.Create(&user)
	return result.Error == nil
}

func CheckValueIsInDB(columnName, value string) bool {
	var user User
	query := dataBase.Where(columnName+" = ?", value).Take(&user)
	return query.RowsAffected != 0
}

// 使用通用函数检查用户名、电话和电子邮件
func CheckUserNameIsInDB(userName string) bool {
	return CheckValueIsInDB("username", userName)
}

func CheckPhoneIsInDB(phone string) bool {
	return CheckValueIsInDB("phone", phone)
}

func CheckEmailIsInDB(email string) bool {
	return CheckValueIsInDB("email", email)
}

// GetUserByName 根据用户名查询用户信息（不包含密码）
func GetUserByName(userName string) (map[string]interface{}, error) {
    var user User
    if err := dataBase.Where("username = ?", userName).First(&user).Error; err != nil {
        return nil, err
    }
    return map[string]interface{}{
        "id":       user.ID,
        "username": user.Username,
        "email":    user.Email,
		"phone":	user.Phone,
		"resumeURL": user.ResumeURL,
        // 其他需要返回的字段
    }, nil
}

// UpdateUserResumeURL 更新用户的 resume_url 字段
func UpdateUserResumeURL(userName string, resumeURL string) error {
	result := dataBase.Model(&User{}).Where("username = ?", userName).Update("resume_url", resumeURL)
	return result.Error
}
