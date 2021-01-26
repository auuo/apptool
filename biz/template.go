package biz

var confTmpl = `package conf

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

var Ins = conf{}

type DatabaseConf struct {
	Username string ` + "`yaml:\"username\"`" + `
	Password string ` + "`yaml:\"password\"`" + `
	Db       string ` + "`yaml:\"db\"`" + `
	Url      string ` + "`yaml:\"url\"`" + `
}

type RedisConf struct {
	Addr     string ` + "`yaml:\"addr\"`" + `
	Password string ` + "`yaml:\"password\"`" + `
}

type conf struct {
	Database DatabaseConf ` + "`yaml:\"database\"`" + `
	Redis    RedisConf    ` + "`yaml:\"redis\"`" + `
}

func Setup(path string) {
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		panic("read config file error, " + err.Error())
	}
	if err = yaml.Unmarshal(yamlFile, &Ins); err != nil {
		panic("init config error, " + err.Error())
	}
}
`

var yamlTmpl = `database:
  username: root
  password: 
  db: matrix
  url: localhost

redis:
  addr: localhost:6379
  password:`

var baseModelTmpl = `package model

type Model struct {
	ID        int64 ` + "`json:\"id\" gorm:\"primary_key;auto_increment\"`" + `
	CreatedAt int   ` + "`json:\"created_at\"`" + `
	UpdatedAt int   ` + "`json:\"updated_at\"`" + `
}

func ClearCreate(m *Model) {
	m.ID = 0
	ClearUpdate(m)
}

func ClearUpdate(m *Model) {
	m.CreatedAt = 0
	m.UpdatedAt = 0
}
`
var userModelTmpl = `package model

const (
	USER_NOT_COMPLETE = 0
	USER_REGULAR      = 1
)

type User struct {
	Model
	Mid      string ` + "`json:\"mid\" gorm:\"unique_index;not null\"`" + `
	Phone    string ` + "`json:\"phone\" gorm:\"unique_index;index;not null;size:24\"`" + `
	Nickname string ` + "`json:\"nickname\" gorm:\"size:24;not null\"`" + `
	Avatar   string ` + "`json:\"avatar\" gorm:\"size:1024;not null\"`" + `
	Status   int    ` + "`json:\"status\" gorm:\"not null\"`" + `
}

func (User) TableName() string {
	return "user"
}
`

var mainTmpl = `package main

import (
	"log"

	"{{.modName}}/conf"
	"{{.modName}}/dao"
	"{{.modName}}/handler"
	"{{.modName}}/middleware"
	"{{.modName}}/pkg/app"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func initLog() {
	logrus.SetReportCaller(true)
}

func initRouter() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.GenLogId)

	r.GET("ping", app.Wrap(handler.Pong))

	return r
}

func main() {
	conf.Setup("conf/conf.yaml")
	initLog()

	r := initRouter()
	pprof.Register(r)
	dao.InitMysql()

	log.Fatal(r.Run(":8080"))
}
`

var pingTmpl = `package handler

import "github.com/gin-gonic/gin"

func Pong(c *gin.Context) (interface{}, error) {
	return "pong", nil
}
`

var appResponseTmpl = `package app

import (
	"net/http"

	"{{.modName}}/pkg/e"

	"github.com/gin-gonic/gin"
)

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "ok",
		"data": data,
	})
}

func Abort(c *gin.Context, err error) {
	if dErr, ok := err.(*e.MatrixError); ok {
		c.AbortWithStatusJSON(http.StatusOK, gin.H{
			"code": dErr.Code,
			"msg":  dErr.Msg,
			"data": nil,
		})
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, gin.H{
		"code": 500,
		"msg":  err.Error(),
		"data": nil,
	})
}

func JSON(c *gin.Context, data interface{}, err error) {
	if err != nil {
		Abort(c, err)
		return
	}
	Success(c, data)
}

func SetUserId(c *gin.Context, id int64) {
	c.Set("userId", id)
}

func GetUserId(c *gin.Context) int64 {
	return c.GetInt64("userId")
}
`

var appWrapperTmpl = `package app

import (
	"github.com/gin-gonic/gin"
)

type Handler func(c *gin.Context) (interface{}, error)

func Wrap(handler Handler) func(*gin.Context) {
	return func(c *gin.Context) {
		resp, err := handler(c)
		JSON(c, resp, err)
	}
}
`

var appContextTmpl = `package app

import (
	"context"
	"{{.modName}}/pkg/logs"

	"github.com/gin-gonic/gin"
)

func NewCtx(c *gin.Context) context.Context {
	id := c.Value(logs.LogIdKey)
	if id == 0 {
		id = logs.NewLogId()
	}
	return context.WithValue(c.Copy(), logs.LogIdKey, id)
}
`

var errTmpl = `package e

import (
	"fmt"
)

type MatrixError struct {
	Code int
	Msg  string
	Err  error
}

func (m *MatrixError) Error() string {
	return fmt.Sprintf("%d: %s", m.Code, m.Msg)
}

func (m *MatrixError) Wrap(err error) error {
	return &MatrixError{
		Code: m.Code,
		Msg:  m.Msg,
		Err:  err,
	}
}

func (m *MatrixError) Unwrap() error {
	return m.Err
}

func New(code int, msg string) *MatrixError {
	return &MatrixError{
		Code: code,
		Msg:  msg,
	}
}
`

var errCommonTmpl = `package e

var (
	ErrorUnknown = New(1, "fail")

	ErrInvalidParams     = New(2, "请求参数错误")
	ErrForbidden         = New(3, "没有权限")
	ErrUnauthorized      = New(4, "未登录")
	ErrAuthorizedTimeout = New(5, "登录失效")
	ErrTokenCheckError   = New(6, "登录验证失败")

	ErrFind   = New(7, "获取失败")
	ErrCreate = New(8, "创建失败")
	ErrUpdate = New(9, "更新失败")
	ErrDelete = New(10, "删除失败")

	ErrUserNotFound = New(11, "用户不存在")
	ErrLogin        = New(12, "登录失败")
)
`

var genLogIdTmpl = `package middleware

import (
	"{{.modName}}/pkg/logs"
	"github.com/gin-gonic/gin"
)

func GenLogId(c *gin.Context) {
	id := logs.NewLogId()
	c.Set(logs.LogIdKey, id)
	c.Header("X-LOGID", id)
	c.Next()
}
`

var logTmpl = `package logs

import (
	"context"

	"github.com/sirupsen/logrus"
)

type F = logrus.Fields

func CtxError(ctx context.Context, err error) *logrus.Entry {
	entry := logrus.WithField("logId", ctx.Value(LogIdKey))
	if err != nil {
		entry = entry.WithError(err)
	}
	return entry
}
`

var logIdTmpl = `package logs

import (
	"fmt"
	"sync/atomic"
	"time"
)

var localId int64 = 0

const LogIdKey = "logId"

func NewLogId() string {
	return fmt.Sprintf("%d%d", getMSTimestamp(), atomic.AddInt64(&localId, 1))
}

func getMSTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
`

var modTmpl = `module {{.modName}}

require (
	github.com/gin-contrib/pprof v1.3.0
	github.com/gin-gonic/gin v1.6.2
	github.com/go-sql-driver/mysql v1.5.0
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/sirupsen/logrus v1.2.0
	golang.org/x/crypto v0.0.0-20190308221718-c2843e01d9a2 // indirect
	gopkg.in/yaml.v2 v2.2.8
	gorm.io/driver/mysql v1.0.1
	gorm.io/gorm v1.20.1
)

go 1.13
`

var daoTmpl = `package dao

import (
	"fmt"

	"{{.modName}}/conf"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitMysql() {
	d := conf.Ins.Database
	source := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=true&loc=Local", d.Username, d.Password, d.Url, d.Db)
	db, err := gorm.Open(mysql.Open(source), &gorm.Config{PrepareStmt: true})
	if err != nil {
		panic("init mysql error, " + err.Error())
	}
	DB = db
}
`

var gitIgnoreTmpl = `**/.DS_Store
.idea/
`