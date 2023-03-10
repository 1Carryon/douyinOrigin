package controller

import (
	"douyinOrigin/dao"
	"douyinOrigin/middleware/jwt"
	"douyinOrigin/service"

	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Response 基础返回响应信息
type Response struct {
	//状态码0成功，其他值失败
	StatusCode int    `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}
type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id"`
	Token  string `json:"token"`
}
type UserResponse struct {
	Response `binding:"required"`
	User     service.User `json:"user"`
}

// Register 用户注册 /user/register/
func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	usi := service.UserServiceImpl{}

	u := usi.GetTableUserByUsername(username)
	if username == u.Name {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "User already exist",
			},
		})
	} else { //插入新用户
		newUser := dao.TableUser{
			Name: username,
			//调用加密算法先加密再存入数据库中
			Password: jwt.EnCoder(password),
		}
		if usi.InsertTableUser(&newUser) != true {
			println("Insert newUser Fail")
		}
		u := usi.GetTableUserByUsername(username) //返回user结构体
		token := jwt.GenerateToken(username)      //生成鉴权
		//创建成功后返回用户 id 和权限token
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{
				StatusCode: 0,
			},
			UserId: u.Id,
			Token:  token,
		})

	}
}

// Login 用户登录 /user/login/
func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	//将用户输入密码加密后与数据库中的密码比较，提高安全性
	encoderPassword := jwt.EnCoder(password)
	println("encoderPassword: ", encoderPassword) //标准错误输出encoderPassword
	usi := service.UserServiceImpl{}
	u := usi.GetTableUserByUsername(username) //查询数据库中是否有
	if encoderPassword == u.Password {
		token := jwt.GenerateToken(username) //对用户名加鉴权
		//	返回响应
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{
				StatusCode: 0,
				StatusMsg:  "登录成功",
			},
			UserId: u.Id,
			Token:  token,
		})
	} else {
		//	返回响应
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "账号或密码错误",
		})
	}

}

// UserInfo 用户信息 /user/
func UserInfo(c *gin.Context) {
	userId := c.Query("user_id") //从请求头中拿出数据
	//token := c.Query("token")
	//将字符串转换为数字 参数2：转换的进制；参数3 返回结果的bit大小
	id, _ := strconv.ParseInt(userId, 10, 64)
	usi := service.UserServiceImpl{}
	//var u userService.User
	//var err error
	//if token != "" {
	//	u, err = usi.GetUserByCurId(id, id)
	//} else {
	//	u, err = usi.GetUserById(id)
	//}
	u, err := usi.GetUserByCurId(id, id)
	if err == nil {
		fmt.Println("用户信息不为空")
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 0},
			User:     u,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status_code": 1,
			"status_msg":  "获取用户信息失败",
			"user":        nil,
		})
	}
}
