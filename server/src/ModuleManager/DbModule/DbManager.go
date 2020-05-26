/**
 数据库管理
 */
//使用golang查询mysql数据库时，如果取出来的其中一个字段的值为null，则会panic，进入defer，最后得到的也是空。
//暂时改变数据库的字段为不能为null来解决
package DbModule

import (
	//"database/sql"

	"fmt"
	_ "github.com/go-sql-driver/mysql"
	//相比database/sql方法还多了新语法，也就是实现将获取的数据直接转换结构体实现。
	//Get(dest interface{}, …) error 用于获取单个结果然后Scan
	//Select(dest interface{}, …) error Select用来获取结果切片
	"github.com/jmoiron/sqlx" //go get "github.com/jmoiron/sqlx"
)

type UserId int
type UserName string
type UserPwd string
type UseType string
type UseTime int
type Mon int32
type Year int32

/**
 db模块
 */
type DbModule interface {
	Insert(...interface{}) error // 插入函数
	Del(...interface{}) error
	query() error
	update() error
}

/**
 db
 */
type Db struct {
	Ip string  // ip
	MaxIdle int //
	MaxOpen int // 最大连接数
	User string // 数据库的用户
	Pwd string // 数据库的密码
	DbName string // 数据库名字
	Port string // 端口
	Pool *sqlx.DB // 连接池
}

/**
 db初始化
 */
func (b *Db) Init() (err error) {
	url := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8",
		b.User,
		b.Pwd,
		b.Ip,
		b.Port,
		b.DbName,
		)
	fmt.Println("连接Db：", url)
	b.Pool, err = sqlx.Open("mysql", url) // 全局只需要调用一次
	if err != nil {
		return
	}
	err = b.Pool.Ping() // 每次用时，都需要ping一下
	if err != nil {
		return
	}
	if b.MaxIdle == 0 {
		b.MaxIdle = 10
	}
	if b.MaxOpen == 0 {
		b.MaxOpen = 20
	}
	// 设置最大连接数
	b.Pool.SetMaxIdleConns(b.MaxIdle) // 设置最大空闲连接数
	b.Pool.SetMaxOpenConns(b.MaxOpen) // 设置最大打开的连接数
	return
}

