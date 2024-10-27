/* 用户表 */
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,  -- 主键，自增ID
    username VARCHAR(50) NOT NULL UNIQUE,      -- 用户名，不允许为空且唯一
    phone VARCHAR(15) UNIQUE,                  -- 手机号，唯一
    email VARCHAR(100) UNIQUE,                 -- 邮箱，唯一
    password VARCHAR(255) NOT NULL,     -- 密码，不允许为空（实际应用中应存储加密后的密码）
    account_balance DECIMAL(10, 2),     -- 账号余额，十进制数，保留两位小数
    vip_level INT,                      -- VIP等级
    age INT,                           -- 年龄
    resume_url VARCHAR(255)            -- 简历地址（假设为URL链接）
);