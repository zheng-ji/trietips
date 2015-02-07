/**
 * zhengji@youmi.net
 * 2015-02-07 09:32:00
 */

package main

import (
	"trietips/src/trie"
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
	//globalTrie = trie.Build("../data/ouwan_app.csv")
	globalTrie = trie.Build("../data/ad_app.csv")
}


func simpleSuggest(w http.ResponseWriter, r *http.Request) {

	var valueList []ValueJson

	r.ParseForm()
    keyword := r.Form["query"]
	fmt.Println("query:", keyword)
	if len(keyword) == 0 {

	} else {
        nodes := trie.Search(globalTrie, keyword[0], 20)
        for _, node := range nodes {
			var value ValueJson
			value.Value = string(node.FullWord)
			valueList = append(valueList, value)
        }
	}
	//if len(valueList) > 10 {
		//valueList = valueList[:10]
	//}

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

func index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/index.html")
}

func main() {

    Init()
	// route
	http.HandleFunc("/", index)
    http.HandleFunc("/suggest/", simpleSuggest)
	http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})
	err := http.ListenAndServe("0.0.0.0:8080", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}