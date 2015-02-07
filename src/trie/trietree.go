/**
*  Author: zhengji@youmi.net
*  2015-02-07 07:30:29
 */

package trie

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// search_count
type Node struct {
	Data     map[string]*Node
	Key      string
	Is_leaf  bool
	Weight   float64
	FullWord string
}

func (n *Node) Init(key string) {
	n.Data = make(map[string]*Node)
	n.Is_leaf = false
	n.Weight = 0
	n.Key = key
	n.FullWord = ""
}

// api for operation
func (n *Node) Has_next() bool {
	return len(n.Data) > 0
}

func (n *Node) AddSubnode(keyword string, subnode *Node) {
	n.Data[keyword] = subnode
}

func (n *Node) GetSubnode(keyword string) *Node {
	return n.Data[keyword]
}

func (n *Node) Get_the_top_node(prefix string) *Node {
	top := n
	for _, c := range prefix {
		top = top.GetSubnode(string(c))
		if top == nil {
			return nil
		}
	}
	return top
}

// api for sort

type NodeList []*Node

func (nl NodeList) Len() int {
	return len(nl)
}

func (nl NodeList) Less(i, j int) bool {
	return nl[i].Weight > nl[j].Weight
}

func (nl NodeList) Swap(i, j int) {
	nl[i], nl[j] = nl[j], nl[i]
}

// api for search
func Depth_walk(node *Node) map[string]*Node {
	result := make(map[string]*Node)
	if node.Is_leaf {
		result[""] = node
	}

	if node.Has_next() {
		for k, _ := range node.Data {
			s := Depth_walk(node.GetSubnode(k))
			for sk, sv := range s {
				result[k+sk] = sv
			}
		}
	} else {
		result[""] = node
	}
	return result
}
func Search(node *Node, prefix string, limit int) NodeList {

	node = node.Get_the_top_node(prefix)

	result := make(map[string]*Node)

	if node == nil {
		return make(NodeList, 0)
	}

	if node.Is_leaf {
		result[prefix] = node
	}

	d := Depth_walk(node)
	for suffix, sub_node := range d {
		result[prefix+suffix] = sub_node
	}

	last_result := make(NodeList, 0, len(result))
	for _, n := range result {
		last_result = append(last_result, n)
	}

	//关键词要去重 根据full_word

	sort.Sort(last_result)
	if len(last_result) < limit {
		return last_result
	}
	return last_result[:limit]

}

func (node *Node) Str() string {
	return "<Node>"
}

func (node *Node) Add(keyword string, weight float64) {
	one_node := node

	total_chars := []rune(keyword)
	last_index := len(total_chars) - 1

	for current_index, c := range total_chars {
		char := string(c)

		if find_node, found := one_node.Data[char]; found {
			one_node = find_node
			if current_index == last_index {
				one_node.Is_leaf = true
				one_node.FullWord = keyword
				one_node.Weight = weight
			}
		} else {
			new_node := new(Node)
			new_node.Init(char)
			if current_index == last_index {
				new_node.Is_leaf = true
				new_node.FullWord = keyword
				new_node.Weight = weight
			}
			one_node.AddSubnode(char, new_node)
			one_node = new_node
		}
	}
}

func (node *Node) Delete(keyword string, judge_leaf bool) {
	if len(keyword) == 0 {
		return
	}

	top_node := node.Get_the_top_node(keyword)
	if top_node == nil {
		return
	}

	//递归往上，对父节点做的判断  遇到节点是某个关键词节点时，要退出
	if judge_leaf {
		if top_node.Is_leaf {
			return
		}
	} else { //非递归，调用delete
		if !top_node.Is_leaf {
			return
		}
	}

	if top_node.Has_next() {
		top_node.Is_leaf = false
		return
	} else {
		this_node := top_node
		chars := []rune(keyword)
		prefix := string(chars[:len(chars)-1])
		top_node = node.Get_the_top_node(prefix)
		delete(top_node.Data, this_node.Key)
		node.Delete(prefix, true)
	}
}

func Build(file_path string) *Node {
	n := new(Node)
	n.Init("")

	f, err := os.Open(file_path)
	if err != nil {
		fmt.Println("Error opening data file ")
		os.Exit(1)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if err == nil {
			n.Add(line, 0)
		}
	}
	return n
}
