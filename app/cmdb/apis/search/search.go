package search

import (
	"fiy/app/cmdb/models/resource"
	"fiy/pkg/es"
	"fiy/tools/app"
	"fmt"
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
	)

	searchResult, err = es.EsClient.Query()
	if err != nil {
		app.Error(c, -1, err, "分页查询关联类型列表失败")
		return
	}

	fmt.Println(searchResult.TotalHits())
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
