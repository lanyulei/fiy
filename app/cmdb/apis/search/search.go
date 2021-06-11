package search

import (
	"fiy/app/cmdb/models/resource"
	"fiy/pkg/es"
	"fiy/tools/app"
	"reflect"

	"github.com/olivere/elastic/v7"

	"github.com/gin-gonic/gin"
)

/*
  @Author : lanyulei
*/

func GetData(c *gin.Context) {
	var (
		err          error
		searchResult *elastic.SearchResult
		list         []resource.Data
		searchValue  interface{}
	)

	searchValue = c.DefaultQuery("value", "")
	if searchValue == "" {
		app.Error(c, -1, nil, "请输入搜索内容")
		return
	}

	searchResult, err = es.EsClient.Query(searchValue)
	if err != nil {
		app.Error(c, -1, err, "分页查询关联类型列表失败")
		return
	}

	if searchResult.TotalHits() > 0 {
		for _, item := range searchResult.Each(reflect.TypeOf(resource.Data{})) {
			if t, ok := item.(resource.Data); ok {
				list = append(list, t)
			}
		}
	}

	app.OK(c, map[string]interface{}{
		"list":  list,
		"total": searchResult.TotalHits(),
	}, "")
}
