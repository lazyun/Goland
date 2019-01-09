package mysqlTools

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"strconv"
)


type MysqlConfig struct {
	Host       string `json:"host"`
	Port       int    `json:"port"`
	User       string `json:"user"`
	Passwd     string `json:"passwd"`
	Db         string `json:"db"`
	Charset    string `json:"charset"`
	Autocommit bool   `json:"autocommit"`
} //`json:"mysql-config"`


type SqlTools struct {
	// 配置信息
	Host		string
	User		string
	Passwd		string
	Db			string
	Port		int
	UrlStr		string

	// mysql 连接池
	Pool		*sql.DB

	// err 信息
	SqlErr		error
	logHandle	func (interface{})

	// sql Stmt 模版
	StmtMap		map[string]*sql.Stmt
}

// Set Mysql config and connect Mysql
// Return true or false. true mean connect success. SqlTools.SqlErr is fail to connect
func (this *SqlTools) SetConfig(host, user, passwd, db string, port int, charset string, autocommit bool) bool {
	this.Host = host
	this.User = user
	this.Passwd = passwd
	this.Db = db
	this.Port = port
	this.UrlStr = user + ":" + this.Passwd + "@tcp(" + host + ":" + strconv.Itoa(port) + ")/" + db + "?charset=" + charset
	if autocommit {
		this.UrlStr += "&autocommit=true"
	} else {
		this.UrlStr += "&autocommit=false"
	}

	this.StmtMap = make(map[string]*sql.Stmt)
	return this.Connect()
}

// Connect mysql
// Return true or false. SqlTools.SqlErr is fail to connect
func (this *SqlTools) Connect() (bool) {
	this.Pool, this.SqlErr = sql.Open("mysql", this.UrlStr)
	if nil != this.SqlErr {
		return false
	}

	this.SqlErr = this.Pool.Ping()
	if nil != this.SqlErr {
		return false
	}

	return true
}


// Close Mysql
func (this *SqlTools) Close() {
	this.Pool.Close()
}


func (this *SqlTools) SetLogHandle(funcName func (msg interface{})) {
	this.logHandle = funcName
}


func (this *SqlTools) SetErrLog(msg string, value ...interface{}) {
	if nil == this.logHandle {
		//fmt.Println("log handle is nil")
		//fmt.Printf(msg, value ...)
		return
	}

	ret := fmt.Sprintf(msg, value ...)
	this.logHandle(ret)
}


func (this *SqlTools) SetMaxConn(num int) {
	this.Pool.SetMaxOpenConns(num)
}

// Set Mysql sql template, parameter use ? replace
// Return true or false. SqlTools.SqlErr is fail to set
func (this *SqlTools) SetSqlStatement(key, value string) bool {
	temp, err := this.Pool.Prepare(value)
	if nil != err {
		this.SetErrLog("Set sql statement [%s] [%s]fail []\n", key, value, err)
		this.SqlErr = err
		return false
	}

	this.StmtMap[key] = temp
	return true
}


// Use sql template to query mysql
// Parameter key: SetSqlStatement the key, value: sql parameter
// Return true
func (this *SqlTools) QueryOneStmt(key string, value ...interface{}) (sql.Row, bool) {
	stmt, ok := this.StmtMap[key]
	if !ok {
		this.SetErrLog("Find sql statement [%s] fail\n", key)
		return sql.Row{}, false
	}

	row := new(sql.Row)
	if nil == value {
		row = stmt.QueryRow()
	} else {
		row = stmt.QueryRow(value ...)
	}

	// get value use row.Scan()
	return *row, true
}


func (this *SqlTools) ChangeStmt(key string, value ...interface{}) (sql.Result, bool) {
	stmt, ok := this.StmtMap[key]
	if !ok {
		this.SetErrLog("Find sql statement [%s] fail\n", key)
		return *new(sql.Result), false
	}

	this.SqlErr = nil
	var ret = new(sql.Result)

	if nil == value {
		*ret, this.SqlErr = stmt.Exec()
	} else {
		*ret, this.SqlErr = stmt.Exec(value...)
	}

	if nil != this.SqlErr {
		return *new(sql.Result), false
	}

	// ret.RowsAffected() ret.LastInsertId()
	return *ret, true
}


func (this *SqlTools) Query(sqlCmd string, args ...interface{}) ( func(isClose bool) ([]sql.RawBytes, bool) ) {
	rows, err := this.Pool.Query(sqlCmd, args...)
	if nil != err {
		this.SqlErr = err
		this.SetErrLog("Sql [%s] [%s]query fail [%s]\n", sqlCmd, args, err)
		return nil
	}

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		this.SqlErr = err
		this.SetErrLog("Sql [%s] [%s]query fail [%s]\n", sqlCmd, args, err)
		return nil
	}

	//fmt.Println("query field is ", columns)

	return func (isClose bool) ([]sql.RawBytes, bool) {
		if isClose {
			rows.Close()
		}

		values := make([]sql.RawBytes, len(columns))
		scanArgs := make([]interface{}, len(columns))
		for i := range values {
			scanArgs[i] = &values[i]
		}

		ok := rows.Next()
		if !ok {
			rows.Close()
			this.SqlErr = nil
			return nil, ok
		}

		err := rows.Scan(scanArgs...)
		if nil != err {
			rows.Close()
			this.SqlErr = err
			this.SetErrLog("Sql [%s] [%s]query fail [%s]\n", sqlCmd, args, err)
			return nil, false
		}

		// string( values[0] )
		return values, true
	}

}


func (this *SqlTools) ExecCmd(sql string, values ...interface{}) (sql.Result) {
	// insert、update、delete
	ok := true
	err := this.Pool.Ping()

	if nil != err {
		this.SqlErr = err
		this.SetErrLog("Sql ping fail [%s]\n", err)
		ok = this.Connect()
	}

	if !ok {
		return nil
	}

	ret, err := this.Pool.Exec(sql, values ...)
	if err != nil {
		this.SqlErr = err
		return nil
	}

	// ret.RowsAffected()、ret.LastInsertId()
	return ret
}

