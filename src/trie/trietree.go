/**
 * code by zheng-ji.info
 */

package trie

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Node struct {
	Link     map[string]*Node
	Key      string
	IsLeaf   bool
	Weight   float64
	LongWord string
}

func (n *Node) Init(key string) {
	n.Link = make(map[string]*Node)
	n.IsLeaf = false
	n.Weight = 0
	n.Key = key
	n.LongWord = ""
}

func (n *Node) Has_next() bool {
	return len(n.Link) > 0
}

func (n *Node) AddSubnode(keyword string, subnode *Node) {
	n.Link[keyword] = subnode
}

func (n *Node) GetSubnode(keyword string) *Node {
	return n.Link[keyword]
}

func (n *Node) GetLassNodeWithPrefix(prefix string) *Node {
	top := n
	for _, c := range prefix {
		top = top.GetSubnode(string(c))
		if top == nil {
			return nil
		}
	}
	return top
}

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

func Depth_walk(node *Node) map[string]*Node {
	result := make(map[string]*Node)
	if node.IsLeaf {
		result[""] = node
	}

	if node.Has_next() {
		for k, _ := range node.Link {
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

	node = node.GetLassNodeWithPrefix(prefix)

	result := make(map[string]*Node)

	if node == nil {
		return make(NodeList, 0)
	}

	if node.IsLeaf {
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

	//去重
	sort.Sort(last_result)
	if len(last_result) < limit {
		return last_result
	}
	return last_result[:limit]

}

func (node *Node) Add(keyword string, weight float64) {
	one_node := node

	total_chars := []rune(keyword)
	last_index := len(total_chars) - 1

	for current_index, c := range total_chars {
		char := string(c)

		if find_node, found := one_node.Link[char]; found {
			one_node = find_node
			if current_index == last_index {
				one_node.IsLeaf = true
				one_node.LongWord = keyword
				one_node.Weight = weight
			}
		} else {
			new_node := new(Node)
			new_node.Init(char)
			if current_index == last_index {
				new_node.IsLeaf = true
				new_node.LongWord = keyword
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

	top_node := node.GetLassNodeWithPrefix(keyword)
	if top_node == nil {
		return
	}

	//递归往上，对父节点做的判断  遇到节点是某个关键词节点时，要退出
	if judge_leaf {
		if top_node.IsLeaf {
			return
		}
	} else { //非递归，调用delete
		if !top_node.IsLeaf {
			return
		}
	}

	if top_node.Has_next() {
		top_node.IsLeaf = false
		return
	} else {
		this_node := top_node
		chars := []rune(keyword)
		prefix := string(chars[:len(chars)-1])
		top_node = node.GetLassNodeWithPrefix(prefix)
		delete(top_node.Link, this_node.Key)
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
