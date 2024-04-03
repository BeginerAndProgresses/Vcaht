package web

import (
	"Vchat/internal/domain"
	"Vchat/internal/service"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"unicode/utf8"
)

const (
	// EmailRegexpPattern 邮箱通配符
	EmailRegexpPattern = `^[\w\-\.]+@([\w-]+\.)+[\w-]{2,}$`
	// PasswordRegexpPattern 至少包含一个数字一个字母和一个特殊字符
	PasswordRegexpPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[@$!%*#?&])[A-Za-z\d@$!%*#?&]{8,}$`
	DateRegexpPattern     = `^\d{4}\-\d{2}\-\d{2}$`
	PhoneRegexpPattern    = `^[^0-9]1[(38)|(55)|(86)|(52)][0-9]{9}[0-9]$`
)

type UserHandler struct {
	// 将正则分为两个是为了进行预编译，提升正则速度
	emailRegexp    *regexp.Regexp
	passwordRegexp *regexp.Regexp
	dateRegexp     *regexp.Regexp
	phoneRegexp    *regexp.Regexp
	svc            *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{
		emailRegexp:    regexp.MustCompile(EmailRegexpPattern, regexp.None),
		passwordRegexp: regexp.MustCompile(PasswordRegexpPattern, regexp.None),
		dateRegexp:     regexp.MustCompile(DateRegexpPattern, regexp.None),
		phoneRegexp:    regexp.MustCompile(PhoneRegexpPattern, regexp.None),
		svc:            svc,
	}
}

func (h *UserHandler) RegisterRouter(r *gin.Engine) {
	u := r.Group("/users")
	{
		u.POST("/signup", h.Signup)
		u.POST("/login", h.Login)
		u.GET("/profile", h.Profile)
		u.POST("/edit", h.Edit)
	}
}

func (h *UserHandler) Signup(c *gin.Context) {
	type Req struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}
	var req Req
	err := c.Bind(&req)
	if err != nil {
		return
	}
	if req.ConfirmPassword != req.Password {
		c.String(http.StatusOK, "两次输入密码不一致")
		return
	}
	h.ValidateEmail(c, req.Email)
	isPassword, err := h.passwordRegexp.MatchString(req.Password)
	if err != nil {
		c.String(http.StatusOK, "系统错误")
		return
	}
	if !isPassword {
		c.String(http.StatusOK, "密码格式应由字母、数字和特殊符号组成")
		return
	}
	err = h.svc.Signup(c, domain.UserDomain{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		c.String(http.StatusOK, "系统错误")
		return
	}
	c.String(http.StatusOK, "注册成功")
}

func (h *UserHandler) Login(c *gin.Context) {
	type Req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req Req
	if err := c.Bind(&req); err != nil {
		c.String(http.StatusOK, "系统错误")
		return
	}
	h.ValidateEmail(c, req.Email)
	u, err := h.svc.Login(c, req.Email, req.Password)
	switch err {
	case nil:
		sess := sessions.Default(c)
		sess.Set("uid", u.Id)
		sess.Options(sessions.Options{
			MaxAge: 60 * 60 * 24,
		})
		err = sess.Save()
		if err != nil {
			c.String(http.StatusOK, "系统错误")
			return
		}
		c.String(http.StatusOK, "登录成功")
	case service.ErrUserNotFound:
		c.String(http.StatusOK, "用户不存在")
	case service.ErrInvalidUserOrPassword:
		c.String(http.StatusOK, "用户名或密码错误")
	default:
		c.String(http.StatusOK, err.Error())
	}
}

func (h *UserHandler) ValidateEmail(c *gin.Context, email string) {
	isEmail, err := h.emailRegexp.MatchString(email)
	if err != nil {
		c.String(http.StatusOK, "系统错误")
		return
	}
	if !isEmail {
		c.String(http.StatusOK, "邮箱格式错误")
		return
	}

}

func (h *UserHandler) Edit(c *gin.Context) {
	type Req struct {
		Nickname string `json:"nickname"`
		Birthday string `json:"birthday"`
		AboutMe  string `json:"aboutMe"`
	}
	var req Req
	err := c.Bind(&req)
	if err != nil {
		c.String(http.StatusOK, "系统错误")
		return
	}
	n := utf8.RuneCountInString(req.Nickname)
	if n > 14 {
		c.String(http.StatusOK, "昵称不能超过14个字")
		return
	}
	a := utf8.RuneCountInString(req.AboutMe)
	if a > 255 {
		c.String(http.StatusOK, "个人简介不能超过255个字")
		return
	}
	isBd, err := h.dateRegexp.MatchString(req.Birthday)
	if err != nil {
		c.String(http.StatusOK, "系统错误")
		return
	}
	if !isBd {
		c.String(http.StatusOK, "生日格式错误")
		return
	}
	uid := h.getUidFromSession(c)
	if uid == -1 {
		c.String(http.StatusOK, "系统错误")
		return
	}
	err = h.svc.Edit(c, &domain.UserDomain{Id: uid,
		Nickname: req.Nickname,
		Birthday: req.Birthday,
		AboutMe:  req.AboutMe})
	if err != nil {
		c.String(http.StatusOK, "系统错误")
		return
	}
	c.String(http.StatusOK, "修改成功")
}

func (h *UserHandler) Profile(c *gin.Context) {
	uid := h.getUidFromSession(c)
	if uid == -1 {
		c.String(http.StatusOK, "系统错误")
		return
	}
	ud, err := h.svc.Profile(c, uid)
	if err != nil {
		c.String(http.StatusOK, "错误信息："+err.Error())
		return
	}
	c.JSON(http.StatusOK, ud)
}

// getUidFromSession 从Session中获取Uid如果获取不到，就返回-1
func (h *UserHandler) getUidFromSession(c *gin.Context) int64 {
	get := sessions.Default(c).Get("uid")
	if get == nil {
		return -1
	}
	uid, ok := get.(int64)
	if !ok {
		return -1
	}
	return uid
}
