### TrieTips

#### 简介

使用前缀树实现的搜索框关键词提示功能, 数据结构用Go语言编写

### 编译和运行

```
cd src; make 
./trietips
```

### 访问

```
http://127.0.0.1:8080
```

### 界面效果

```
http://zheng-ji.info/blog/2015/02/08/trie-suggestion/
```

### 接口示范

```
搜索关键词：
curl "http://127.0.0.1:8080/tips?keyword=hello"

添加关键词：
curl -d "keyword=hello" "http://127.0.0.1:8080/add"

删除关键词：
curl -d "keyword=hello" "http://127.0.0.1:8080/del"
```

[Author](http://zheng-ji.info)
