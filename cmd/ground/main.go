package main

import (
	"fmt"
	"github.com/goccy/go-json"
)

func main() {
	v := make(map[string]interface{})
	//var v map[string]string
	cMap := make(map[string]interface{})
	cMap["a"] = "b"
	c, _ := json.Marshal(v)
	fmt.Println(string(c))
	pMap := make(map[string]interface{})
	pMap["bbb"] = "ccc"
	p, _ := json.Marshal(pMap)
	fmt.Println(string(p))
	v["const"] = cMap
	v["pass_down"] = pMap
	result, _ := json.Marshal(v)
	fmt.Println(string(result))
}
