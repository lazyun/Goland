# Mongo 使用说明

## 获取第三方包

`go get gopkg.in/mgo.v2`

## 使用说明

在插入 MongoDB 时，如果要保证插入字段的有序性，使用以下结构

```
insertValues = bson.D{
    {"key1", "value1"},
    {"key2", "value2"},
    {"key3", "value3"},
    {
        "key4": bson.D{
            {"key4Sub1", "value.."},
            ...
        }
    }
    ...
}
```

`bson.M` 底层是 map 不保证顺序， `bson.D` 底层是列表保证顺序。