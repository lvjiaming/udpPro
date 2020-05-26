package DbModule

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

var (
	infoDb *InfoDb
	infoCreTabStr = "create table if not exists %s " +
		"(id int primary key auto_increment," +
		"time int(20),usetype char(50)," +
		"val float(5,2));"
)

const ( // 统计信息的类型
	WEEK_INFO = iota
	MON_INFO
	YEAR_INFO
)

type StatisticalInfo struct { // 统计信息
	Week float32
	Mon float32
	Year float32
}

type InfoDb struct {
	db *Db
	userId UserId
	tableName string
}

/**
用于接收查询的结构体
*/
type Info struct {
	Id int `db:"id"`
	UserType string `db:"usetype"`
	Val float32 `db:"val"`
	Time int `db:"time"`
}

/**
 插入数据
 */
func (i *InfoDb) Insert (utype UseType, val float32) error {
	if utype == "" || val == 0 {
		return errors.New("信息错误")
	}
	var err error
	pingErr := i.db.Pool.Ping()
	if pingErr != nil {
		return pingErr
	}
	sql := fmt.Sprintf("INSERT %s SET usetype=?,val=?,time=?", i.getTableName())
	stmt, err := i.db.Pool.Prepare(sql)
	if err != nil {
		return err
	} else {
		defer stmt.Close()
	}
	timeTem := time.Now().Unix()
	//fmt.Println("当前时间戳：", timeTem)
	_, err = stmt.Exec(utype, val, timeTem)
	return err
}

/**
删除信息
 */
func (i *InfoDb) DelInfo (idList []int32) error {
	if len(idList) == 0 {
		return errors.New("未传入任何id")
	}
	var err error
	pingErr := i.db.Pool.Ping()
	if pingErr != nil {
		return pingErr
	}
	// 事务
	//一个Tx会在整个生命周期中保存一个连接，然后在调用commit或Rollback()的时候释放掉
	tx, txErr := i.db.Pool.Beginx()
	if txErr != nil {
		tx.Rollback()
		return txErr
	}
	sql := fmt.Sprintf("delete from %s where id = ?", i.getTableName())
	for _, id := range idList{
		_, e := i.db.Pool.Exec(sql, id)
		if e != nil {
			err = e
			tx.Rollback()
			break
		}
	}
	if err == nil {
		tx.Commit()
	}
	return err
}

/*
查找所有的
 */
func (i *InfoDb) QueryAllInfo () ([]*Info, error) {
	var err error
	err = i.db.Pool.Ping()
	if err != nil {
		return nil, err
	}
	var info []*Info
	// Select函数可以获取切片数据，并结构体化
	sql := fmt.Sprintf("SELECT id,usetype,val,time FROM %s", i.getTableName())
	err = i.db.Pool.Select(&info, sql)
	return info, err
}

/**
 根据筛选条件查询
 */
func (i *InfoDb) QueryInfoToTerm (args ...interface{}) ([]*Info, error) {
	if len(args) == 0 {
		return nil, errors.New("未传入任何参数")
	}
	var err error
	err = i.db.Pool.Ping()
	if err != nil {
		return nil, err
	}
	var uType interface{}
	var uTime []UseTime
	for _, val := range args {
		switch val.(type) {
		case UseType:
			uType = val
		case UseTime:
			if val.(UseTime) != 0 {
				uTime = append(uTime, val.(UseTime))
			}
		}
	}
	var infoList []*Info
	sql := ""
	if uType != nil && len(uTime) == 2 {
		if uTime[0] > uTime[1] {
			uTime[0], uTime[1] = uTime[1], uTime[0]
		}
		sql = "SELECT id,usetype,val,time FROM %s where time>=? and time <=? and usetype like ?"
	} else {
		if uType != nil {
			sql = "SELECT id,usetype,val,time FROM %s where usetype like ?"
		} else if len(uTime) == 2 {
			sql = "SELECT id,usetype,val,time FROM %s where time between ? and ?"
		} else {
			return infoList, errors.New("条件不对啊")
		}
	}
	//fmt.Println("参数：", uType, uTime)
	sqlStr := fmt.Sprintf(sql, i.getTableName())
	//fmt.Println("查询语句：", sqlStr)
	//通配符％，应该是参数字符串的一部分，也就是说%必须作为字符串写到参数里面去，而不能在sql语句
	if uType != nil && len(uTime) == 2 {
		err = i.db.Pool.Select(&infoList, sqlStr, uTime[0], uTime[1], "%" + uType.(UseType) + "%")
	} else {
		if uType != nil {
			err = i.db.Pool.Select(&infoList, sqlStr, "%" + uType.(UseType) + "%")
		} else if len(uTime) == 2 {
			err = i.db.Pool.Select(&infoList, sqlStr, uTime[0], uTime[1])
		}
	}
	return infoList, err
}

/**
更具infoid查找info
 */
func (i *InfoDb) QueryInfoById (id int32) (Info, error) {
	if id == 0 {
		return Info{}, errors.New("id 未传入")
	}
	var err error
	var info Info
	err = i.db.Pool.Ping()
	if err != nil {
		return info, err
	}
	err = i.db.Pool.Get(&info,
		fmt.Sprintf("SELECT id,usetype,val,time FROM %s where id = ?", i.getTableName()), id)
	return info, err
}

/**
 更新信息
 */
func (i *InfoDb) UpdateInfo (id int, args ...interface{}) error {
	if id == 0 || len(args) == 0 {
		return errors.New("参数错误")
	}
	var err error
	err = i.db.Pool.Ping()
	if err != nil {
		return err
	}
	var uType interface{}
	var uVal interface{}
	for _, val := range args {
		switch val.(type) {
		case UseType:
			uType = val
		case float32:
			uVal = val
		}
	}
	sql := ""
	if uType != nil && uVal != nil {
		sql = fmt.Sprintf("update %s set usetype = ?,val = ? where id = ?", i.getTableName())
	} else {
		if uType != nil {
			sql = fmt.Sprintf("update %s set usetype = ? where id = ?", i.getTableName())
		} else if uVal != nil {
			sql = fmt.Sprintf("update %s set val = ? where id = ?", i.getTableName())
		} else {
			return errors.New("请传入需要修改的参数")
		}
	}
	stmt, err := i.db.Pool.Prepare(sql)
	if err != nil {
		return err
	} else {
		defer stmt.Close()
	}
	if uType != nil && uVal != nil {
		_, err = stmt.Exec(uType, uVal, id)
	} else {
		if uType != nil {
			_, err = stmt.Exec(uType, id)
		} else if uVal != nil {
			_, err = stmt.Exec(uVal, id)
		}
	}
	return err
}

/**
 查询统计信息列表
 */
func (i *InfoDb) QueryStatisticalInfoList (infoType int, times ...interface{}) ([]*Info, error) {
	var err error
	err = i.db.Pool.Ping()
	if err != nil {
		return nil, err
	}
	var info []*Info
	sql := i.getStatisticalInfoSql("id,usetype,val,time", infoType, times)
	err = i.db.Pool.Select(&info, sql)
	return info, err
}


/**
 查找统计信息
 */
func (i *InfoDb) QueryStatisticalInfo (infoType int, times ...interface{}) (float32, error) {
	var err error
	var res float32
	err = i.db.Pool.Ping()
	if err != nil {
		return res, err
	}
	sql := i.getStatisticalInfoSql("SUM(val)", infoType, times)
	stem, paperErr := i.db.Pool.Prepare(sql)
	if paperErr != nil {
		return res, paperErr
	}
	defer stem.Close()
	err = stem.QueryRow().Scan(&res)
	return res, err
}

/**
 获取当前的统计信息
 */
func (i *InfoDb) GetCurStatisticalInfo () *StatisticalInfo {
	week, _ := i.QueryStatisticalInfo(WEEK_INFO)
	mon, _  := i.QueryStatisticalInfo(MON_INFO)
	year, _ := i.QueryStatisticalInfo(YEAR_INFO)
	sInfo := &StatisticalInfo{
		Week: week,
		Mon: mon,
		Year: year,
	}
	return sInfo
}

/**
 获取统计信息的sql语句
 */
func (i *InfoDb)getStatisticalInfoSql(queryInfo string, infoType int, times ...interface{}) string {
	var year string
	var mon string
	for _, val := range times{
		switch val.(type) {
		case Year:
			year = strconv.Itoa(int(val.(Year)))
		case Mon:
			mon = strconv.Itoa(int(val.(Mon)))
		}
	}
	if year == "" {
		year = "year(curdate())"
	}
	if mon == "" {
		mon = "month(curdate())"
	}
	sql := ""

	// Select函数可以获取切片数据，并结构体化
	//FROM_UNIXTIME可以将时间戳转换为时间
	// curdate获取当前时间
	switch infoType {
	case WEEK_INFO:
		sql = fmt.Sprintf("SELECT %s FROM %s " +
			"where month(FROM_UNIXTIME(%s.time)) = " +
			"%s and " +
			"week(FROM_UNIXTIME(%s.time)) = week(curdate())",queryInfo , i.getTableName(),
			i.getTableName(), mon, i.getTableName())
	case MON_INFO:
		sql = fmt.Sprintf("SELECT %s FROM %s " +
			"where month(FROM_UNIXTIME(%s.time)) = " +
			"%s and " +
			"year(FROM_UNIXTIME(%s.time)) = %s",queryInfo , i.getTableName(),
			i.getTableName(),mon, i.getTableName(), year)
	case YEAR_INFO:
		sql = fmt.Sprintf("SELECT %s FROM %s " +
			"where year(FROM_UNIXTIME(%s.time)) = " +
			"%s", queryInfo, i.getTableName(),
			i.getTableName(), year)
	}
	return sql
}

/*
获取表名
 */
func (i *InfoDb) getTableName () string {
	//string转成int：
	//int, err := strconv.Atoi(string)
	//string转成int64：
	//int64, err := strconv.ParseInt(string, 10, 64)
	//int转成string：
	//string := strconv.Itoa(int)
	//int64转成string：
	//string := strconv.FormatInt(int64,10)
	return "tab_info" + strconv.Itoa(int(i.userId))
}

func GetInfoDb(db *Db, uid UserId) (*InfoDb, error) {
	pingErr := db.Pool.Ping()
	if pingErr != nil {
		return nil, pingErr
	}
	infoDb = &InfoDb{
		db: db,
		userId: uid,
	}
	tabName := infoDb.getTableName()
	sql := fmt.Sprintf(infoCreTabStr, tabName)
	stmt, err := db.Pool.Prepare(sql)
	if err != nil {
		infoDb = nil
		return nil, err
	}
	defer stmt.Close()
	stmt.Exec()
	return infoDb, err
}