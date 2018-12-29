package main

import (
	SqlT "../mysqlTools"
	"fmt"
	"os"
	"bufio"
	//"github.com/derekparker/delve/pkg/dwarf/reader"
)


func main() {

	sqlT := new(SqlT.SqlTools)

	// 连接 Mysql
	ok := sqlT.SetConfig("127.0.0.1", "root", "", "test", 3306, "utf8", true)
	if !ok {
		fmt.Printf("Connect mysql fail error is %s !\n", sqlT.SqlErr.Error())
		os.Exit(1)
	}

	// 使用模版方式查询 只能返回一条数据
	ok = sqlT.SetSqlStatement("findId", "select id from goTest where created > ? limit 1")
	if !ok {
		fmt.Printf("Set mysql template findId fail error is %s !\n", sqlT.SqlErr.Error())
		os.Exit(1)
	}

	ok, QueryOneStmtRet := sqlT.QueryOneStmt("findId", 1531800702)
	if !ok {
		fmt.Printf("Query mysql template findId fail error is %s !\n", sqlT.SqlErr.Error())
		os.Exit(1)
	}

	// 取出模版数据
	var id int
	err := QueryOneStmtRet.Scan(&id)
	if nil != err {
		fmt.Printf("Get mysql template findId result fail error is %s !\n", err.Error())
		os.Exit(1)
	}

	fmt.Printf("Query id is %d\n", id)

	// 使用模版方式修改
	ok = sqlT.SetSqlStatement("insertGoTest", "insert into goTest(name, created) values(?, ?)")
	if !ok {
		fmt.Printf("Set mysql template insertGoTest fail error is %s !\n", sqlT.SqlErr.Error())
		os.Exit(1)
	}

	ok, ChangeStmtRet := sqlT.ChangeStmt("insertGoTest", "la~la~la~", 1531800702)
	if !ok {
		fmt.Printf("Change mysql template insertGoTest fail error is %s !\n", sqlT.SqlErr.Error())
		os.Exit(1)
	}

	idChangeStmt, err := ChangeStmtRet.LastInsertId()
	if nil != err {
		fmt.Printf("Get Change mysql template insertGoTest result fail error is %s !\n", sqlT.SqlErr.Error())
		os.Exit(1)
	}

	rowChangeStmt, err := ChangeStmtRet.RowsAffected()
	if nil != err {
		fmt.Printf("Get Change mysql template insertGoTest affected row fail error is %s !\n", sqlT.SqlErr.Error())
		os.Exit(1)
	}

	fmt.Printf("Insert data last id is %d, affected row %d\n", idChangeStmt, rowChangeStmt)

	// 查询多条数据
	queryRet := sqlT.Query("select id, name from goTest limit ?", 3)
	if nil == queryRet {
		fmt.Printf("Query fail error %s !\n", sqlT.SqlErr.Error())
		os.Exit(1)
	}
	//queryRet(true)
	//
	//queryRet = sqlT.Query("select id, name from goTest limit ?", 3)
	//queryRet(true)
	//
	//queryRet = sqlT.Query("select id, name from goTest limit ?", 3)
	//queryRet(true)
	//
	//queryRet = sqlT.Query("select id, name from goTest limit ?", 3)
	//queryRet(true)
	//
	//queryRet = sqlT.Query("select id, name from goTest limit ?", 3)
	//queryRet(true)
	//
	//queryRet = sqlT.Query("select id, name from goTest limit ?", 3)
	//fmt.Println(queryRet)
	//reader := bufio.NewReader(os.Stdin)
	//cc, _, err := reader.ReadRune()
	//fmt.Println(cc)
	//os.Exit(1)

	// 设置最大连接数目 通过 show processlist 查看数据库状态
	//sqlT.SetMaxConn(2)
	//fmt.Println(1)
	//queryRet = sqlT.Query("select id, name from goTest limit ?", 3)
	//fmt.Println(2)
	//queryRet = sqlT.Query("select id, name from goTest limit ?", 3)
	//fmt.Println(3)
	//queryRet = sqlT.Query("select id, name from goTest limit ?", 3)
	//fmt.Println(4)
	//reader := bufio.NewReader(os.Stdin)
	//cc, _, err := reader.ReadRune()
	//fmt.Println(cc)
	//os.Exit(1)

	for {
		ok, queryRetValue := queryRet(false)
		if !ok {
			if nil == sqlT.SqlErr {
				fmt.Printf("Get query result all!\n")
				break
			}

			fmt.Printf("Get query result fail error %s !\n", sqlT.SqlErr.Error())
			break
		}

		queryId := string( queryRetValue[0] )
		queryName := string( queryRetValue[1] )
		fmt.Printf("Query result id is %s, name is %s\n", queryId, queryName)
	}

	// 插入、更新、修改
	var execRetValue1 int64
	execRet := sqlT.ExecCmd("insert into goTest(name, created) values(?, ?)", "la~la~la~", 1531800702)
	if nil == execRet {
		fmt.Printf("ExecCmd fail error is %s\n", sqlT.SqlErr.Error() )
	} else {
		execRetValue1, _ = execRet.LastInsertId()
		execRetValue2, _ := execRet.RowsAffected()
		fmt.Printf("ExecCmd insert row %d, last id is %d\n", execRetValue2, execRetValue1)
	}

	execRet1 := sqlT.ExecCmd("update goTest set name = ? where id = ?", "biu~biu~biu~", 1)
	if nil == execRet1 {
		fmt.Printf("ExecCmd fail error is %s\n", sqlT.SqlErr.Error() )
	} else {
		execRetValue3, _ := execRet.LastInsertId()
		execRetValue4, _ := execRet.RowsAffected()
		fmt.Printf("ExecCmd update row %d, last id is %d\n", execRetValue4, execRetValue3)
	}

	execRet2 := sqlT.ExecCmd("delete from goTest where id = ?", execRetValue1)
	if nil == execRet2 {
		fmt.Printf("ExecCmd fail error is %s\n", sqlT.SqlErr.Error() )
	} else {
		execRetValue5, _ := execRet.LastInsertId()
		execRetValue6, _ := execRet.RowsAffected()
		fmt.Printf("ExecCmd delete row %d, last id is %d\n", execRetValue6, execRetValue5)
	}

}
