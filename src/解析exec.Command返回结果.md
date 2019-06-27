```go
package main

import(
  "fmt"
  "strings"
)

type ItemList struct {
  Node string
  UUID string
}
func main(){
  srcStr := `"node1" {f6fa9977-fe1e-44c9-ae14-35450b52898f}
"node2" {49f147c1-757d-48cd-938b-658389d30209}
"node3" {7f4f6365-bcfa-4ac3-b788-673e4be0aa00}
"node6" {7c90466e-0917-4b81-a66a-296084fc6401}
"node7" {0bebb20f-e0c8-4f62-a602-d43e13b28859}
"node8" {e112785a-307a-4a6a-b2d6-c860f1c80ee9}
"k8s-node1" {0775edab-9ccb-457c-b843-ee9904faaf41}
"k8s-node2" {ce8f2c06-2d94-4292-b963-9951a77953c7}
"k8s-node3" {2bc521a1-0454-4c23-b33c-936b9fab44de}
"ubuntu" {2afeb17a-29b9-4f29-a1da-efc46cac676b}
"k8s-node4" {72f2e60b-f33e-438c-8bba-49523a2ee61c}
`

  var result []string
  result = strings.Split(srcStr,"\n")
  
  // 此处只考首尾是空格，没有考 tab 制表符，也没有对所获得的数据是否为空进行判断，一切以
  // 数据是完美的方式来解决
  var itemList = make([]ItemList, len(result))
  for i, v := range result {
    val := strings.Split(strings.Trim(v," ")," ")
    itemList[i].Node = val[0][1:len(val[0])-1]
    itemList[i].Node = val[1][1:len(val[1])-1]
  }
}
```
