/**
 * code by zheng-ji.info
 */

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"trietips/src/trie"
)

type ValueJson struct {
	Value string `json:"value"`
}

var globalTrie *trie.Node

func Init() {
	globalTrie = trie.Build("../data/ad_app.csv")
}

//获取提示
func gettips(w http.ResponseWriter, r *http.Request) {

	var valueList []ValueJson

	r.ParseForm()
	keyword := r.Form["query"]
	fmt.Println("query:", keyword)
	if len(keyword) == 0 {
		fmt.Fprintf(w, "[]")
	} else {
		nodes := trie.Search(globalTrie, keyword[0], 20)
		for _, node := range nodes {
			var value ValueJson
			value.Value = string(node.LongWord)
			valueList = append(valueList, value)
		}
	}

	fmt.Println("return", valueList)
	if len(valueList) > 0 {
		b, err := json.Marshal(valueList)
		if err != nil {
			fmt.Println("json err:", err)
		}
		fmt.Fprintf(w, string(b))
	} else {
		fmt.Fprintf(w, "[]")
	}
	return
}

//添加词条
func addentry(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	keyword := r.Form["keyword"]
	if len(keyword) == 0 {
		fmt.Fprintf(w, "fail")
	} else {
		fmt.Println("add entry keword:", keyword)
		globalTrie.Add(keyword[0], 0)
		fmt.Fprintf(w, "ok")
	}
}

//删除词条
func delentry(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	keyword := r.Form["keyword"]
	if len(keyword) == 0 {
		fmt.Fprintf(w, "fail")
	} else {
		fmt.Println("del entry keword:", keyword)
		globalTrie.Delete(keyword[0], false)
		fmt.Fprintf(w, "ok")
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/index.html")
}

func main() {

	Init()
	// route
	http.HandleFunc("/", index)
	http.HandleFunc("/tips/", gettips)
	http.HandleFunc("/add/", addentry)
	http.HandleFunc("/del/", delentry)
	http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})
	err := http.ListenAndServe("0.0.0.0:8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
