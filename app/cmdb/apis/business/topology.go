package business

import (
	"encoding/json"
	"fiy/app/cmdb/models/business"
	"fiy/app/cmdb/models/model"
	"fiy/app/cmdb/models/resource"
	orm "fiy/common/global"
	"fiy/tools/app"

	uuid "github.com/satori/go.uuid"

	"github.com/gin-gonic/gin"
)

/*
  @Author : lanyulei
*/

// 业务树
func BusinessTree(c *gin.Context) {
	var (
		err      error
		biz      model.Info
		treeData []*struct {
			resource.Data
			Children []interface{} `json:"children"`
		}
	)

	// 查询业务模型数据
	err = orm.Eloquent.Model(&model.Info{}).
		Where("identifies = 'built_in_biz'").
		Find(&biz).Error
	if err != nil {
		app.Error(c, -1, err, "查询模型数据失败")
		return
	}
	// 所有业务ID
	bizIdList := make([]int, 0)
	err = orm.Eloquent.Model(&resource.Data{}).Where("info_id = ?", biz.Id).Pluck("id", &bizIdList).Error
	if err != nil {
		app.Error(c, -1, err, "查询业务数据列表失败")
		return
	}
	if len(bizIdList) > 0 {
		// 查询业务数据列表
		err = orm.Eloquent.Model(&resource.Data{}).Find(&treeData, bizIdList).Error
		if err != nil {
			app.Error(c, -1, err, "查询模块列表失败")
			return
		}
	}

	// 所有业务对应的集群ID
	setIdList := make([]int, 0)
	setList := make([]*struct {
		resource.Data
		Biz      int           `json:"biz"`
		Children []interface{} `json:"children"`
	}, 0)
	if len(bizIdList) > 0 {
		err = orm.Eloquent.Model(&resource.Data{}).
			Joins("left join cmdb_resource_data_related as dr on dr.target = cmdb_resource_data.id").
			Where("dr.source in ?", bizIdList).
			Select("cmdb_resource_data.*, dr.source as biz").
			Find(&setList).Error
		if err != nil {
			app.Error(c, -1, err, "查询集群数据失败")
			return
		}

		for _, s := range setList {
			setIdList = append(setIdList, s.Id)
		}
	}

	// 所有业务对应的集群ID
	moduleList := make([]*struct {
		resource.Data
		Cluster int `json:"cluster"`
	}, 0)
	if len(setIdList) > 0 {
		err = orm.Eloquent.Model(&resource.Data{}).
			Joins("left join cmdb_resource_data_related as dr on dr.target = cmdb_resource_data.id").
			Where("dr.source in ?", setIdList).
			Select("cmdb_resource_data.*, dr.source as cluster").
			Find(&moduleList).Error
		if err != nil {
			app.Error(c, -1, err, "查询模块数据失败")
			return
		}
	}

	for _, t := range treeData {
		for _, s := range setList {
			for _, m := range moduleList {
				if m.Cluster == s.Id {
					s.Children = append(s.Children, m)
					continue
				}
			}
			if s.Biz == t.Id {
				t.Children = append(t.Children, s)
				continue
			}
		}
	}

	app.OK(c, treeData, "")
}

// 添加节点
func AddBusinessNode(c *gin.Context) {
	var (
		err    error
		params struct {
			Name     string `json:"name"`
			Template int    `json:"template"`
			Classify int    `json:"classify"`
			Source   int    `json:"source"`
		}
		data      []byte
		modelInfo model.Info
	)

	err = c.ShouldBind(&params)
	if err != nil {
		app.Error(c, -1, err, "参数绑定失败")
		return
	}

	data, err = json.Marshal(map[string]interface{}{"built_in_field_name": params.Name})
	if err != nil {
		app.Error(c, -1, err, "Json序列化数据失败")
		return
	}

	tx := orm.Eloquent.Begin()

	if params.Classify == 1 { // 集群
		// 查询模型信息
		err = tx.Where("identifies = 'built_in_set'").Find(&modelInfo).Error
		if err != nil {
			tx.Rollback()
			app.Error(c, -1, err, "查询集群模型失败")
			return
		}

	} else if params.Classify == 2 { // 模块
		// 查询模型信息
		err = tx.Where("identifies = 'built_in_module'").Find(&modelInfo).Error
		if err != nil {
			tx.Rollback()
			app.Error(c, -1, err, "查询模块模型失败")
			return
		}
	}
	if modelInfo.Id == 0 {
		tx.Rollback()
		app.Error(c, -1, err, "模型ID不能为零值")
		return
	}

	// 新建数据
	resourceData := &resource.Data{
		Uuid:     uuid.NewV4().String(),
		InfoId:   modelInfo.Id,
		InfoName: modelInfo.Name,
		Status:   0,
		Data:     data,
	}
	err = tx.Create(resourceData).Error
	if err != nil {
		tx.Rollback()
		app.Error(c, -1, err, "新建数据失败")
		return
	}

	// 新建数据关联
	err = orm.Eloquent.Create(&resource.DataRelated{
		Source: params.Source,
		Target: resourceData.Id,
	}).Error
	if err != nil {
		tx.Rollback()
		app.Error(c, -1, err, "新建数据关联失败")
		return
	}

	// 模版关联
	if params.Template != 0 {
		err = tx.Create(&business.TemplateRelated{
			TplClassify: params.Classify,
			TplId:       params.Template,
			DataID:      resourceData.Id,
			InfoId:      modelInfo.Id,
		}).Error
		if err != nil {
			tx.Rollback()
			app.Error(c, -1, err, "创建模板关联失败")
			return
		}
	}

	tx.Commit()

	app.OK(c, nil, "")
}

// 删除节点
func DeleteBusinessNode(c *gin.Context) {
	var (
		err          error
		nodeID       string
		relatedCount int64
	)

	nodeID = c.Param("id")

	// 判断是否有关联数据
	err = orm.Eloquent.Model(&resource.DataRelated{}).
		Where("source = ?", nodeID).
		Count(&relatedCount).Error
	if err != nil {
		app.Error(c, -1, err, "查询树节点关联数据失败")
		return
	}

	if relatedCount > 0 {
		app.Error(c, -1, nil, "当前树节点存在关联关系，无法直接删除，请确认")
		return
	}

	err = orm.Eloquent.Model(&resource.Data{}).Delete(&resource.Data{}, nodeID).Error
	if err != nil {
		app.Error(c, -1, err, "删除节点失败")
		return
	}

	app.OK(c, nil, "")
}
