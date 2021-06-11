package search

import (
	"fiy/pkg/es"
	"fiy/tools/app"
	"strconv"

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
		list         []interface{}
		searchValue  interface{}
		page         int
		limit        int
	)

	searchValue = c.DefaultQuery("value", "")
	if searchValue == "" {
		app.Error(c, -1, nil, "请输入搜索内容")
		return
	}

	page, _ = strconv.Atoi(c.DefaultQuery("page", "0"))
	limit, _ = strconv.Atoi(c.DefaultQuery("limit", "10"))

	searchResult, err = es.EsClient.Query(searchValue, page, limit)
	if err != nil {
		app.Error(c, -1, err, "搜索数据失败")
		return
	}

	if searchResult.TotalHits() > 0 {
		for _, item := range searchResult.Hits.Hits {
			list = append(list, item.Source)
		}
	}

	app.OK(c, map[string]interface{}{
		"list":  list,
		"total": searchResult.TotalHits(),
	}, "")
}
