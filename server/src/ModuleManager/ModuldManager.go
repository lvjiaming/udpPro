package ModuleManager

import (
	"../Cfg"
	"./DbModule"
	"errors"
	"fmt"
	"sync"
)

var (
	moduleMge *moduleManager
)

/**
 管理模块
 */
type moduleManager struct {
	db map[string]*DbModule.Db
}

/**
 连接db
 */
func (m *moduleManager) ConnectDb (waitGroup *sync.WaitGroup, dbName string)  {
	defer waitGroup.Done()
	m.db[dbName] = &DbModule.Db{
		Ip: Cfg.DbIp,
		Port: Cfg.DbPort,
		MaxIdle: 10,
		MaxOpen: 20,
		User: Cfg.DbUser,
		Pwd: Cfg.DbPwd,
		DbName: dbName,
	}
	curDb := m.db[dbName]
	err := curDb.Init()
	if err != nil {
		fmt.Println("连接出错：", err.Error())
	} else {
		fmt.Println("连接成功")
		//stmt, err := db.Pool.Prepare("INSERT user SET id=?,name=?,pw=?")
		//defer stmt.Close()
		//if err != nil {
		//	fmt.Println("prepare出错了")
		//}
		//res, err := stmt.Exec(10, "张三", "1234")
		//if err != nil {
		//	fmt.Println("插入出错")
		//} else {
		//	fmt.Println("插入成功", res)
		//}
	}
}

/**
 按dbname获取db
 */
func (m *moduleManager) GetDb (dbName string) (*DbModule.Db, error) {
	var err error
	var dbM *DbModule.Db
	if m.db[dbName] == nil {
		err = errors.New("db not find")
	} else {
		dbM = m.db[dbName]
		//switch dbName {
		//case Cfg.UserDb:
		//	dbM, err = DbModule.GetUserDb(m.db[dbName])
		//case Cfg.InfoDb:
		//	dbM, err = DbModule.GetInfoDb(m.db[dbName])
		//}
	}
	return dbM, err
}

/**
 获取模块管理器
 */
func GetModuleManager() *moduleManager {
	if moduleMge == nil {
		moduleMge = &moduleManager{db: make(map[string]*DbModule.Db)}
	}
	return moduleMge
}