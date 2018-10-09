package mysqlTools

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"strconv"
)


type SqlTools struct {
	// 配置信息
	Host	string
	User	string
	Passwd	string
	Db		string
	Port	int
	UrlStr	string

	// mysql 连接池
	Pool	*sql.DB

	// err 信息
	SqlErr	error

	// sql Stmt 模版
	StmtMap	map[string]*sql.Stmt
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


// Set Mysql sql template, parameter use ? replace
// Return true or false. SqlTools.SqlErr is fail to set
func (this *SqlTools) SetSqlStatement(key, value string) bool {
	temp, err := this.Pool.Prepare(value)
	if nil != err {
		this.SqlErr = err
		return false
	}

	this.StmtMap[key] = temp
	return true
}


// Use sql template to query mysql
// Parameter key: SetSqlStatement the key, value: sql parameter
// Return true
func (this *SqlTools) QueryOneStmt(key string, value ...interface{}) (bool, sql.Row) {
	stmt, ok := this.StmtMap[key]
	if !ok {
		fmt.Printf("not found %s sql statement!\n", key)
		return false, sql.Row{}
	}

	row := new(sql.Row)
	if nil == value {
		row = stmt.QueryRow()
	} else {
		row = stmt.QueryRow(value ...)
	}

	// get value use row.Scan()
	return true, *row
}


func (this *SqlTools) ChangeStmt(key string, value ...interface{}) (bool, sql.Result) {
	stmt, ok := this.StmtMap[key]
	if !ok {
		fmt.Printf("not found %s sql statement!\n", key)
		return false, *new(sql.Result)
	}

	this.SqlErr = nil
	var ret = new(sql.Result)

	if nil == value {
		*ret, this.SqlErr = stmt.Exec()
	} else {
		*ret, this.SqlErr = stmt.Exec(value...)
	}

	if nil != this.SqlErr {
		return false, *new(sql.Result)
	}

	// ret.RowsAffected() ret.LastInsertId()
	return true, *ret
}


func (this *SqlTools) Query(sqlCmd string, args ...interface{}) ( func() (bool, []sql.RawBytes) ) {
	rows, err := this.Pool.Query(sqlCmd, args...)
	if nil != err {
		this.SqlErr = err
		fmt.Printf("sql query fail! %s \n", err.Error())
		return nil
	}

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		this.SqlErr = err
		fmt.Printf("sql query get columns fail! %s \n", err.Error())
		return nil
	}

	//fmt.Println("query field is ", columns)

	return func () (bool, []sql.RawBytes) {
		values := make([]sql.RawBytes, len(columns))
		scanArgs := make([]interface{}, len(columns))
		for i := range values {
			scanArgs[i] = &values[i]
		}

		ok := rows.Next()
		if !ok {
			this.SqlErr = nil
			return ok, nil
		}

		err := rows.Scan(scanArgs...)
		if nil != err {
			this.SqlErr = err
			fmt.Printf("sql query fail! %s \n", err.Error())
			return false, nil
		}

		// string( values[0] )
		return true, values
	}

}


func (this *SqlTools) ExecCmd(sql string, values ...interface{}) (sql.Result) {
	// insert、update、delete
	ok := true
	err := this.Pool.Ping()

	if nil != err {
		this.SqlErr = err
		fmt.Printf("sql connect fail!\n", err)
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

