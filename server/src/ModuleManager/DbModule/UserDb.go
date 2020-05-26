package DbModule

import (
	"fmt"
)

var (
	userDb *UserDb
	createTableStr = "create table if not exists user " +
		"(id int primary key auto_increment," +
		"name char(50)," +
		"pwd char(50));"
	insertStr = "INSERT user SET name=?,pwd=?"
	delStr = "delete from user where id = ?"
	updateNameStr = "update user set name = ? where id = ?"
	updatePwdStr = "update user set pwd = ? where id = ?"
	updatePwdAndNameStr = "update user set name = ?,pwd = ? where id = ?"
	queryAllStr = "SELECT id,name,pwd FROM user"
	queryUserByNameStr = "SELECT id,name,pwd FROM user where name = ?"
	queryUserByIdStr = "SELECT id,name,pwd FROM user where id = ?"
)

/**
 用于接收查询的结构体
 */
type UserInfo struct {
	UserId int `db:"id"`
	UserName string `db:"name"`
	Password string `db:"pwd"`
}

type UserDb struct {
	db *Db
}
/**
 插入一个玩家（需提供名字和密码）
 */
func (u *UserDb) Insert (name UserName, pwd UserPwd) error {
	var err error
	pingErr := u.db.Pool.Ping()
	if pingErr != nil {
		return pingErr
	}
	stmt, err := u.db.Pool.Prepare(insertStr)
	if err != nil {
		return err
	} else {
		defer stmt.Close()
	}
	_, err = stmt.Exec(name, pwd)
	return err
}
/**
 删除一个数据（需提供id）一般不删除
 */
//func (u *UserDb) Del (id int) error {
//	var err error
//	pingErr := u.db.Pool.Ping()
//	if pingErr != nil {
//		return pingErr
//	}
//	stmt, err := u.db.Pool.Prepare(insertStr)
//	if err != nil {
//		return err
//	} else {
//		defer stmt.Close()
//	}
//	_, err = stmt.Exec(id)
//	return err
//}

/**
 查询所有玩家
 */
func (u *UserDb) QueryAllUser () ([]UserInfo, error) {
	var err error
	err = u.db.Pool.Ping()
	if err != nil {
		return nil, err
	}
	var userInfo []UserInfo
	// Select函数可以获取切片数据，并结构体化
	err = u.db.Pool.Select(&userInfo, queryAllStr)
	//fmt.Println(userInfo)
	return userInfo, err
}

/**
 通过筛选条件进行查询
 */
func (u *UserDb) QueryUserToTerm (term interface{}) (UserInfo, error) {
	var err error
	var userInfo UserInfo
	err = u.db.Pool.Ping()
	if err != nil {
		return userInfo, err
	}
	sql := ""
	switch term.(type) {
	case UserName:
		sql = queryUserByNameStr
	case UserId:
		sql = queryUserByIdStr
	}
	err = u.db.Pool.Get(&userInfo, sql, term)
	return userInfo, err
}

/**
 更新玩家信息
 */
func (u *UserDb) Update (args ...interface{}) error {
	var err error
	err = u.db.Pool.Ping()
	if err != nil {
		return err
	}
	var sqlStr = ""
	var name interface{}
	var pwd interface{}
	var id interface{}
	for _, val := range args{
		switch val.(type) {
		case UserName:
			sqlStr = updateNameStr
			name = val
		case UserPwd:
			sqlStr = updatePwdStr
			pwd = val
		case UserId:
			id = val
		}
	}
	if len(args) == 3 {
		sqlStr = updatePwdAndNameStr
	}
	fmt.Println("查询语句：", sqlStr)
	stmt, err := u.db.Pool.Prepare(sqlStr)
	if err != nil {
		return err
	} else {
		defer stmt.Close()
	}
	if len(args) == 3 {
		_, err = stmt.Exec(name, pwd, id)
	} else {
		if name != nil {
			_, err = stmt.Exec(name, id)
		} else {
			_, err = stmt.Exec(pwd, id)
		}
	}
	return err
}

/**
 获取userdb
 */
func GetUserDb(db * Db) (*UserDb, error) {
	pingErr := db.Pool.Ping()
	if pingErr != nil {
		return nil, pingErr
	}
	stmt, err := db.Pool.Prepare(createTableStr)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	stmt.Exec()
	if userDb == nil {
		userDb = &UserDb{db: db}
	}
	return userDb, nil
}

