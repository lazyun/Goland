package mgoTools

import (
	"gopkg.in/mgo.v2"
	"net/url"
	"fmt"
	"reflect"
	"gopkg.in/mgo.v2/bson"
)


// mongo 的操作
// 不是线程安全的
// 操作条件最好使用 bson.M{}

type MgoTools struct {
	MgoUrl				string
	MgoPass				string
	MgoUrlStr 			string
	MgoWorkDB			string
	MgoPassParse		string

	MgoErr				error
	MgoUpdateInfo		*mgo.ChangeInfo

	// 当前 db 下所有的集合名称
	MgoCollections		map[string]int
	MgoCollectionsLen	int

	// 插入是否为重复
	MgoInsertDup		bool

	mongoDB				*mgo.Database
	mongoSession 		*mgo.Session

	mongoColMap			map[string]*mgo.Collection
	mongoLogHandle		func(args ...interface{})
}


func (this *MgoTools) SetConfig(mgoUrl , passswd, db string) bool {
	this.MgoUrl = mgoUrl
	this.MgoPass = passswd
	this.MgoPassParse = url.QueryEscape(passswd)
	this.MgoUrlStr = fmt.Sprintf(mgoUrl, this.MgoPassParse)
	this.MgoWorkDB = db

	this.MgoCollections = make(map[string]int)
	this.mongoColMap = make(map[string]*mgo.Collection)

	this.mongoSession, this.MgoErr = mgo.Dial(this.MgoUrlStr)
	if nil == this.MgoErr {
		this.mongoDB = this.mongoSession.DB(this.MgoWorkDB)
		this.mongoLogHandle = func(args ...interface{}) {

		}
		return true
	}

	return false
}


// 设置异常处理的接口
func (this *MgoTools) SetLogHandle(funcName interface{}) bool {
	refValue := reflect.ValueOf(funcName)
	if reflect.Func != refValue.Kind() {
		return false
	}

	f := func(args ...interface{}) {
		refArgs := []reflect.Value{}
		for _, value := range args {
			refArgs = append(refArgs, reflect.ValueOf(value))
		}

		refValue.Call(refArgs)
	}

	this.mongoLogHandle = f
	return true
}


// 查看所有 DB 需要在 admin 库执行命令的权限
// 返回 DB 列表，error 信息
func (this *MgoTools) MgoDBS(ret *[]string) (bool) {
	*ret, this.MgoErr = this.mongoSession.DatabaseNames()
	if nil == this.MgoErr {
		return true
	}

	this.mongoLogHandle(this.MgoErr)
	return false
}


// 选择工作的 DB，
func (this *MgoTools) SelectDB(db string) {
	this.mongoDB = this.mongoSession.DB(db)
	fmt.Println( this.mongoSession.DatabaseNames() )
}


// 查询当前 db 下所有的 collection
func (this *MgoTools) MgoCols(ret *[]string) (bool) {
	*ret, this.MgoErr = this.mongoDB.CollectionNames()
	if nil == this.MgoErr {
		return true
	}

	this.MgoCollectionsLen = len(*ret)
	for _, value := range *ret{
		this.MgoCollections[value] = 1
	}

	this.mongoLogHandle(this.MgoErr)
	return false
}


// 判断当前 db 是否存在 col
// 返回 true 存在、返回 false 并且 this.MgoErr == nil 时不存在、否则有 error
func (this *MgoTools) ColIsExist(col string) bool {
	var ret []string

	ret, this.MgoErr = this.mongoDB.CollectionNames()
	if nil != this.MgoErr {
		this.mongoLogHandle(this.MgoErr)
		return false
	}

	retLen := len(ret)
	if retLen == this.MgoCollectionsLen {
		goto JUDGE
	}

	this.MgoCollectionsLen = retLen
	for _, value := range ret{
		this.MgoCollections[value] = 1
	}

	JUDGE:
	_, ok := this.MgoCollections[col]
	return ok
}


// 创建 collection 默认的
func (this *MgoTools) CreateCol(col string) bool {
	mgoCol := this.SelectCol(col)

	colInfo := mgo.CollectionInfo{}
	this.MgoErr = mgoCol.Create(&colInfo)
	if nil == this.MgoErr {
		return true
	}

	this.mongoLogHandle(this.MgoErr)
	return false
}


// 选择操作的 Collection
func (this *MgoTools) SelectCol(col string) *mgo.Collection {
	c, ok := this.mongoColMap[col]
	if ok {
		return c
	}

	c = this.mongoDB.C(col)
	this.mongoColMap[col] = c
	return c
}


// **************************************************** 查询 **************************************************************
// 查询 Mongo 返回一条数据
// 参数 col：集合名称、query：查询条件、ret 返回结果结构（字典类型指针）、projections 返回字段
// 返回结果：true：操作成功、false 操作失败
func (this *MgoTools) FindOne(col string, query interface{}, ret interface{}, projections ...interface{}) (bool) {
	mgoCol := this.SelectCol(col)
	findRet := mgoCol.Find(query)
	if nil != projections {
		findRet.Select(projections[0])
	}

	this.MgoErr = findRet.One(ret)
	if nil == this.MgoErr {
		return true
	}

	this.mongoLogHandle(this.MgoErr)
	return false
}


// 查询 Mongo 数据
// 返回 *Query 其他工作需要自己做
func (this *MgoTools) Find(col string, query interface{}) (*mgo.Iter) {
	mgoCol := this.SelectCol(col)
	findRet := mgoCol.Find(query).Iter()
	return findRet
}


// 查询 Mongo 数据并修改返回
// 结果存在 this.MgoUpdateInfo 中
// 返回结果：true：操作成功、false 操作失败
func (this *MgoTools) FindAndModify(
	col string, query interface{}, update interface{}, ret interface{},
	isUpsert, isReturnNew bool, projections ...interface{}) bool {
		mgoCol := this.SelectCol(col)

		change := mgo.Change{
			Update: update,
			Upsert: isUpsert,
			ReturnNew: isReturnNew,
		}

		q := mgoCol.Find(query)
		if nil != projections {
			q.Select(projections[0])
		}

		this.MgoUpdateInfo, this.MgoErr = q.Apply(change, ret)
		if nil == this.MgoErr {
			return true
		}

		this.mongoLogHandle(this.MgoErr)
		return false
}


// **************************************************** 新增 **************************************************************
// 新增记录 一条或多条
// 返回结果：true：操作成功、false 操作失败
func (this *MgoTools) Insert(col string, values ...interface{}) bool {
	mgoCol := this.SelectCol(col)
	this.MgoErr = mgoCol.Insert(values ...)

	if nil == this.MgoErr {
		return true
	}

	this.mongoLogHandle(this.MgoErr)
	return false
}


// bson.D 确保顺序、bson.M 无法确保顺序
func (this *MgoTools) InsertOne(col string, values interface{}) bool {
	mgoCol := this.SelectCol(col)
	this.MgoErr = mgoCol.Insert(values)

	if nil == this.MgoErr {
		return true
	}

	this.mongoLogHandle(this.MgoErr)
	if mgo.IsDup(this.MgoErr) {
		this.MgoInsertDup = true
	} else {
		this.MgoInsertDup = false
	}

	return false
}


// **************************************************** 更新 **************************************************************
// 更新一条记录
// 返回结果：true：操作成功、false 操作失败
func (this *MgoTools) UpdateOne(col string, selector, update interface{}) bool {
	mgoCol := this.SelectCol(col)
	this.MgoErr = mgoCol.Update(selector, update)

	if nil == this.MgoErr {
		return true
	}

	this.mongoLogHandle(this.MgoErr)
	return false
}


// 更新全部记录
// 更新结果存在 this.MgoUpdateInfo 中
// 返回结果：true：操作成功、false 操作失败
func (this *MgoTools) UpdateAll(col string, selector, update interface{}) bool {
	mgoCol := this.SelectCol(col)
	this.MgoUpdateInfo, this.MgoErr = mgoCol.UpdateAll(selector, update)

	if nil == this.MgoErr {
		return true
	}

	this.mongoLogHandle(this.MgoErr)
	return false
}


// 更新一条记录
// 更新结果存在 this.MgoUpdateInfo 中
// 返回结果：true：操作成功、false 操作失败
func (this *MgoTools) UpsertOne(col string, selector, update interface{}) bool {
	mgoCol := this.SelectCol(col)
	this.MgoUpdateInfo, this.MgoErr = mgoCol.Upsert(selector, update)

	if nil == this.MgoErr {
		return true
	}

	this.mongoLogHandle(this.MgoErr)
	return false
}

// **************************************************** 删除 **************************************************************
// 删除一条记录
// 返回结果：true：操作成功、false 操作失败
func (this * MgoTools) DeleteOne(col string, selector interface{}) bool {
	mgoCol := this.SelectCol(col)
	this.MgoErr = mgoCol.Remove(selector)

	if nil == this.MgoErr {
		return true
	}

	this.mongoLogHandle(this.MgoErr)
	return false
}


// 删除全部记录
// 更新结果存在 this.MgoUpdateInfo 中
// 返回结果：true：操作成功、false 操作失败
func (this *MgoTools) DeleteAll(col string, selector interface{}) bool {
	mgoCol := this.SelectCol(col)
	this.MgoUpdateInfo, this.MgoErr = mgoCol.RemoveAll(selector)

	if nil == this.MgoErr {
		return true
	}

	this.mongoLogHandle(this.MgoErr)
	return false
}


// **************************************************** 索引 **************************************************************
// 说明：创建普通索引
// 入参：col：集合名称、key：索引字段、isSparse：索引字段不存在的文档不会建立索引
// 		isUnique：是否为一、isDropDup：删除重复数据
// 输出：true：操作成功、false 操作失败
func (this *MgoTools) CreateIndexNormal(col, key string, isUnique, isDropDup, isSparse bool) bool {
	mgoCol := this.SelectCol(col)

	index := mgo.Index{
		Key: []string{key, },
		Unique: isUnique,
		DropDups: isDropDup,
		Background: true, // See notes.
		Sparse: isSparse,
	}
	this.MgoErr = mgoCol.EnsureIndex(index)
	if nil == this.MgoErr {
		return true
	}

	this.mongoLogHandle(this.MgoErr)
	return false
}


// 说明：创建 hash 索引
// 入参：col：集合名称、key：索引字段、isExist：索引字段不存在的文档不会建立索引
// 输出：true：操作成功、false 操作失败
func (this *MgoTools) CreateIndexHashed(col, key string, isSparse bool) bool {
	mgoCol := this.SelectCol(col)

	hashKey := "$hashed:" + key
	index := mgo.Index{
		Key: []string{hashKey, },
		Unique: false,
		DropDups: false,
		Background: true, // See notes.
		Sparse: isSparse,
	}
	this.MgoErr = mgoCol.EnsureIndex(index)
	if nil == this.MgoErr {
		return true
	}

	this.mongoLogHandle(this.MgoErr)
	return false
}


// 说明：对集合创建分片
// 入参：
func (this *MgoTools) ShardingCol(db, col, key string, isHashed bool, ret interface{}) bool {
	adminDB := this.mongoSession.DB("admin")

	shardingCol := db + "." + col

	var shardingType string
	if isHashed {
		shardingType = "hashed"
	} else {
		shardingType = "1"
	}

	cmd := bson.D{
			{"shardCollection",shardingCol},
			{
				"key",
				bson.M{
					key: shardingType,
				},
			},
		}


	this.MgoErr = adminDB.Run(cmd, ret)
	if nil == this.MgoErr {
		return true
	}

	this.mongoLogHandle(this.MgoErr)
	return false
}

// **************************************************** 执行命令 **************************************************************
// 简单封装、实际使用查看 ShardingCol 函数
func (this *MgoTools) RunCmd(db, col, cmd string, isHashed bool, ret interface{}) bool {
	runCmdDB := this.mongoSession.DB(db)
	this.MgoErr = runCmdDB.Run(cmd, ret)
	if nil == this.MgoErr {
		return true
	}

	this.mongoLogHandle(this.MgoErr)
	return false
}