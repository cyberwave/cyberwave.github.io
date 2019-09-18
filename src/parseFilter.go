/*
请求方法：GET
URL: 正常的URL后面参数携带的是复杂的过滤条件
查询条件
  格式为：Filter.N.查询条件，每个查询条件必含 name, values，选填 operator,默认为 eq, 可选 lt, le, gt, ge, ne,
        like, complex, complex预留用于描述更复杂的 过滤需求，complex的 value 格式使用google规范，
        参考：https://cloud.google.com/sdk/gcloud/reference/topic/filters
  tags.N tag过滤条件，格式为tags.N.查询条件，key,values。例如：tags.1.key=name1&tags.1.values.1=aaaa&tags.2.key=bbb&tags.2.values.1=xxx
*/

package main

type ReqFilter struct {
    Filter []Filter`json:"filter" form:"filter"`
    PageNumber int `json:"pageNumber" form:"pageNumber"`
    PageSize   int `json:"pageSize" form:"pageSize"`
}

type Filter struct {
    Name   string   `json:"name" form:"name"`
    Oper   string   `json:"operator" form:"operator"`
    Values []string `json:"values" form:"values"`
}

func ParseQueryFilter(c echo.Context, keyName string) *Filter {
    var filter *model.Filter
    var params = c.QueryParams()
    var keyNameRegex = regexp.MustCompile(`^filters.([0-9]+).name$`)
    for k, v := range params {
        // filter name的模式 filters.$n.name
        if match := keyNameRegex.FindStringSubmatch(k); len(match) > 1 && v[0] == keyName {
            filter = &model.Filter{Name: v[0]}

            // filter value的模式filters.$n.values.$d
            var valueRegex = regexp.MustCompile(fmt.Sprintf("^filters.%s.values.([0-9]+)$", match[1]))
            for k1, v1 := range params {
                if valueRegex.Match([]byte(k1)) {
                    filter.Values = append(filter.Values, v1[0])
                }
            }
            return filter
        }
    }

    return filter
}

// 自己修改后的
func ParseQueryFilter(ctx *http.Request, keyName string) *Filter {
    var filter *Filter
    var params = ctx.URL.Query()
    
    fmt.Printf("pageSize:%s\n",params["pageSize"])
    fmt.Printf("pageNumber:%s\n",params["pageNumber"])
    fmt.Printf("cluster_ids:%s\n",params["cluster_ids"])

    //var params = ctx.QueryParams()
    var keyNameRegex = regexp.MustCompile(`^filters.([0-9]+).name$`)
    for k, v := range params {
      // filter name的模式 filters.$n.name
      if match := keyNameRegex.FindStringSubmatch(k); len(match) > 1 && v[0] == keyName {
        fmt.Printf("len(v)=%d,v[0]=%s",len(v),v[0])

        filter = &model.Filter{Name: v[0]}

        // filter oper 的模式 filters.$n.operator
        var keyOperRegex = regexp.MustCompile(fmt.Sprintf("^filters.%s.operator$",match[1]))
        for kOper, vOper := range params {
          if keyOperRegex.Match([]byte(kOper)) {
            filter.Oper = vOper[0]
          }
        }

        // filter value的模式filters.$n.values.$d
        var valueRegex = regexp.MustCompile(fmt.Sprintf("^filters.%s.values.([0-9]+)$", match[1]))
        for k1, v1 := range params {
          if valueRegex.Match([]byte(k1)) {
            filter.Values = append(filter.Values, v1[0])
          }
        }
        //return filter
      }
    }

    return filter
}
