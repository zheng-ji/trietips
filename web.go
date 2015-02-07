package main

import (
	"trietips/trie"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type ValueJson struct {
	Value string `json:"value"`
}

var globalTrie * trie.Node

func Init() {
	globalTrie = trie.Build("data.txt")
}


func simpleSuggest(w http.ResponseWriter, r *http.Request) {

	var valueList []ValueJson

	r.ParseForm() //解析参数，默认是不会解析的
	keyword := r.Form["keyword"]

	fmt.Println("keyword:", keyword)
	if len(keyword) == 0 {

	} else {
        nodes := trie.Search(globalTrie, keyword[0], 10)
        for _, node := range nodes {
			var value ValueJson
			value.Value = string(node.FullWord)
			valueList = append(valueList, value)
        }
	}

	if len(valueList) > 10 {
		valueList = valueList[:10]
	}

	fmt.Println("return", valueList)
	if len(valueList) > 0 {
		b, err := json.Marshal(valueList)
		if err != nil {
			fmt.Println("json err:", err)
		}
		fmt.Fprintf(w, string(b)) //这个写入到w的是输出到客户端的
	} else {
		fmt.Fprintf(w, "[]") //这个写入到w的是输出到客户端的
	}
	return
}

func index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/index.html")
}

func main() {

    Init()

	// index
	http.HandleFunc("/", index) //设置访问的路由

	// tips
	http.HandleFunc("/suggest/", simpleSuggest) //设置访问的路由

	// static
	http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})

	err := http.ListenAndServe("0.0.0.0:9090", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
