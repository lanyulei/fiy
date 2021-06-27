package resource

import (
	"encoding/json"
	"fiy/app/cmdb/models/model"
	"fiy/app/cmdb/models/resource"
	"fiy/common/actions"
	orm "fiy/common/global"
	"fiy/common/pagination"
	"fiy/tools/app"
	"fmt"
	"strconv"
	"strings"

	"gorm.io/datatypes"

	uuid "github.com/satori/go.uuid"

	"github.com/gin-gonic/gin"
)

/*
  @Author : lanyulei
*/

// 获取数据列表
func DataList(c *gin.Context) {
	var (
		err             error
		dataList        []*resource.Data
		result          interface{}
		modelID         string
		searchType      string
		searchClassify  string
		searchString    string
		searchList      []map[string]interface{}
		status          string
		nodeID          string
		dataIDs         []int
		searchWhereList []string
	)

	modelID = c.Param("id")

	db := orm.Eloquent.Model(&resource.Data{}).Where("info_id = ?", modelID)

	searchType = c.DefaultQuery("search_type", "1")
	searchClassify = c.DefaultQuery("search_classiy", "2")
	searchString = c.DefaultQuery("search_list", "[]")

	err = json.Unmarshal([]byte(searchString), &searchList)
	if err != nil {
		app.Error(c, -1, err, "参数异常")
		return
	}

	if len(searchList) > 0 {
		for _, v := range searchList {
			if v["identifies"] != "" && v["value"] != "" {
				if searchClassify == "1" {
					searchWhereList = append(searchWhereList, fmt.Sprintf("data->'$.%v' = '%v'", v["identifies"], v["value"]))
				} else if searchClassify == "2" {
					searchWhereList = append(searchWhereList, fmt.Sprintf("data->'$.%v' like '%%%v%%'", v["identifies"], v["value"]))
				}
			}
		}
		if searchType == "1" { // 与
			db = db.Where(strings.Join(searchWhereList, " and "))
		} else if searchType == "2" {
			db = db.Where(fmt.Sprintf("(%v)", strings.Join(searchWhereList, " or ")))
		}
	}

	status = c.DefaultQuery("status", "")
	nodeID = c.DefaultQuery("nodeID", "")

	if status != "" {
		db = db.Where("status = ?", status)
	}

	if nodeID != "" {
		err = orm.Eloquent.Model(&resource.DataRelated{}).
			Where("source = ? and target_info_id = ?", nodeID, modelID).
			Pluck("target", &dataIDs).Error
		if err != nil {
			app.Error(c, -1, err, "查询节点绑定的数据ID失败")
			return
		}

		db = db.Where("id in ?", dataIDs)
	}

	result, err = pagination.Paging(&pagination.Param{
		C:  c,
		DB: db,
	}, &dataList)
	if err != nil {
		app.Error(c, -1, err, "分页查询关联类型列表失败")
		return
	}

	app.OK(c, result, "")
}

// 新建数据
func CreateData(c *gin.Context) {
	var (
		err  error
		data resource.Data
	)

	data.Uuid = uuid.NewV4().String()
	err = c.ShouldBind(&data)
	if err != nil {
		app.Error(c, -1, err, "参数绑定失败")
		return
	}

	tx := orm.Eloquent.Begin()

	err = tx.Create(&data).Error
	if err != nil {
		tx.Rollback()
		app.Error(c, -1, err, "新建资源数据失败")
		return
	}

	// 添加操作审计
	err = actions.AddAudit(c,
		tx,
		"资源",
		"资源数据",
		"新建",
		fmt.Sprintf("新建资源数据 \"%d\"", data.Id),
		map[string]interface{}{},
		data)
	if err != nil {
		tx.Rollback()
		app.Error(c, -1, err, "添加操作审计失败")
		return
	}

	tx.Commit()

	app.OK(c, nil, "")
}

// 批量新建数据
func BatchCreateData(c *gin.Context) {
	var (
		err  error
		data []*resource.Data
	)

	err = c.ShouldBind(&data)
	if err != nil {
		app.Error(c, -1, err, "参数绑定失败")
		return
	}

	for _, d := range data {
		d.Uuid = uuid.NewV4().String()
	}

	tx := orm.Eloquent.Begin()

	err = tx.Create(&data).Error
	if err != nil {
		tx.Rollback()
		app.Error(c, -1, err, "批量新建资源数据失败")
		return
	}

	// 添加操作审计
	err = actions.AddAudit(c,
		tx,
		"资源",
		"资源数据",
		"新建",
		fmt.Sprintf("批量新建资源数据"),
		map[string]interface{}{},
		nil)
	if err != nil {
		tx.Rollback()
		app.Error(c, -1, err, "添加操作审计失败")
		return
	}

	tx.Commit()

	app.OK(c, nil, "")
}

// 删除数据
func DeleteData(c *gin.Context) {
	var (
		err    error
		dataId string
		data   resource.Data
	)

	dataId = c.Param("id")

	err = orm.Eloquent.Find(&data, dataId).Error
	if err != nil {
		app.Error(c, -1, err, "查询资源数据失败")
		return
	}

	tx := orm.Eloquent.Begin()

	err = tx.Delete(&resource.Data{}, dataId).Error
	if err != nil {
		tx.Rollback()
		app.Error(c, -1, err, "删除资源数据失败")
		return
	}

	// 添加操作审计
	err = actions.AddAudit(c,
		tx,
		"资源",
		"资源数据",
		"删除",
		fmt.Sprintf("删除资源数据 \"%s\"", dataId),
		data,
		map[string]interface{}{})
	if err != nil {
		tx.Rollback()
		app.Error(c, -1, err, "添加操作审计失败")
		return
	}

	tx.Commit()

	app.OK(c, nil, "")
}

// 编辑数据
func EditData(c *gin.Context) {
	var (
		err     error
		data    resource.Data
		oldData resource.Data
		dataId  string
	)

	dataId = c.Param("id")

	err = c.ShouldBind(&data)
	if err != nil {
		app.Error(c, -1, err, "参数绑定失败")
		return
	}

	err = orm.Eloquent.Find(&oldData, dataId).Error
	if err != nil {
		app.Error(c, -1, err, "查询资源数据失败")
		return
	}

	tx := orm.Eloquent.Begin()

	err = tx.Model(&data).Where("id = ?", dataId).Updates(map[string]interface{}{
		"info_id": data.InfoId,
		"data":    data.Data,
	}).Error
	if err != nil {
		tx.Rollback()
		app.Error(c, -1, err, "更新资源数据失败")
		return
	}

	// 添加操作审计
	err = actions.AddAudit(c,
		tx,
		"资源",
		"资源数据",
		"编辑",
		fmt.Sprintf("编辑资源数据 \"%s\"", dataId),
		map[string]interface{}{
			"info_id": oldData.InfoId,
			"data":    oldData.Data,
		},
		map[string]interface{}{
			"info_id": data.InfoId,
			"data":    data.Data,
		})
	if err != nil {
		tx.Rollback()
		app.Error(c, -1, err, "添加操作审计失败")
		return
	}

	tx.Commit()

	app.OK(c, nil, "")
}

// 获取数据详情
func GetDataDetails(c *gin.Context) {
	var (
		err     error
		details resource.Data
		dataId  string
	)

	dataId = c.Param("id")

	err = orm.Eloquent.Model(&details).Where("id = ?", dataId).Find(&details).Error
	if err != nil {
		app.Error(c, -1, err, "参数绑定失败")
		return
	}

	app.OK(c, details, "")
}

// 查询节点对应的模型分组及模型信息
func GetNodeModelData(c *gin.Context) {
	var (
		err           error
		level         string
		nodeID        string
		modelNodeIDs  []int
		clusterData   model.Info
		modelData     model.Info
		targetInfoIDs []int
		modelList     []*struct {
			model.Group
			ModelList []model.Info `json:"model_list"`
		}
		modelGroupIDs []int
	)

	level = c.DefaultQuery("level", "")
	nodeID = c.DefaultQuery("nodeID", "")
	if level == "" || nodeID == "" {
		app.Error(c, -1, nil, "参数异常，level和nodeID必须传递")
		return
	}

	// 查询集群信息
	err = orm.Eloquent.Model(&clusterData).Select("id").Where("identifies = ?", "built_in_set").Find(&clusterData).Error
	if err != nil {
		app.Error(c, -1, err, "查询集群信息失败")
		return
	}

	// 查询模块信息
	err = orm.Eloquent.Model(&modelData).Select("id").Where("identifies = ?", "built_in_module").Find(&modelData).Error
	if err != nil {
		app.Error(c, -1, err, "查询模块信息失败")
		return
	}

	// 查询所有的模块节点
	if level == "1" {
		// 查询所有的集群ID
		clusterIDs := make([]int, 0)
		err = orm.Eloquent.Model(&resource.DataRelated{}).
			Where("source = ? and target_info_id = ?", nodeID, clusterData.Id).
			Pluck("target", &clusterIDs).Error
		if err != nil {
			app.Error(c, -1, err, "查询集群ID失败")
			return
		}

		// 查询所有的模块ID
		err = orm.Eloquent.Model(&resource.DataRelated{}).
			Where("source in ? and target_info_id = ?", clusterIDs, modelData.Id).
			Pluck("target", &modelNodeIDs).Error
		if err != nil {
			app.Error(c, -1, err, "查询模块ID失败")
			return
		}
	} else if level == "2" {
		// 查询所有的模块ID
		err = orm.Eloquent.Model(&resource.DataRelated{}).
			Where("source = ? and target_info_id = ?", nodeID, modelData.Id).
			Pluck("target", &modelNodeIDs).Error
		if err != nil {
			app.Error(c, -1, err, "查询模块ID失败")
			return
		}
	} else if level == "3" {
		modelID, _ := strconv.Atoi(nodeID)
		modelNodeIDs = append(modelNodeIDs, modelID)
	}

	// 查询节点关联的资产数据
	err = orm.Eloquent.Model(&resource.DataRelated{}).
		Select("distinct target_info_id").
		Where("source in ?", modelNodeIDs).
		Pluck("target_info_id", &targetInfoIDs).Error
	if err != nil {
		app.Error(c, -1, err, "查询绑定的数据模型ID失败")
		return
	}

	// 查询模型数据
	err = orm.Eloquent.Model(&model.Info{}).
		Select("distinct group_id").
		Where("id in ?", targetInfoIDs).
		Pluck("group_id", &modelGroupIDs).Error
	if err != nil {
		app.Error(c, -1, err, "查询模型数据失败")
		return
	}

	// 获取分组
	err = orm.Eloquent.Model(&model.Group{}).Where("id in ?", modelGroupIDs).Find(&modelList).Error
	if err != nil {
		app.Error(c, -1, err, "获取分组列表失败")
		return
	}

	// 获取分组对应的模型列表
	for _, group := range modelList {
		db := orm.Eloquent.Model(&model.Info{})
		err = db.Where("group_id = ? and is_usable = 1", group.Id).
			Find(&group.ModelList).Error
		if err != nil {
			app.Error(c, -1, err, "获取模型信息失败")
			return
		}
	}

	app.OK(c, modelList, "")
}

// 导出数据
func ExportData(c *gin.Context) {
	var (
		err       error
		infoID    string
		fields    []model.Fields
		tHeader   []string
		filterVal []string
		dataList  []datatypes.JSON
	)

	infoID = c.Param("id")
	err = orm.Eloquent.Model(&model.Fields{}).Where("info_id = ?", infoID).Find(&fields).Error
	if err != nil {
		app.Error(c, -1, err, "查询表头失败")
		return
	}

	for _, f := range fields {
		tHeader = append(tHeader, f.Name)
		filterVal = append(filterVal, f.Identifies)
	}

	err = orm.Eloquent.Model(&resource.Data{}).Where("info_id = ?", infoID).Pluck("data", &dataList).Error
	if err != nil {
		app.Error(c, -1, err, "查询表头标识失败")
		return
	}

	app.OK(c, map[string]interface{}{
		"tHeader":   tHeader,
		"filterVal": filterVal,
		"dataList":  dataList,
	}, "")
}

// 资源关联绑定
func AddRelatedData(c *gin.Context) {
	var (
		err    error
		params struct {
			Source       int   `json:"source"`
			Target       []int `json:"target"`
			SourceInfoId int   `json:"source_info_id"`
			TargetInfoId int   `json:"target_info_id"`
		}
		currentRelatedList []resource.DataRelated
		relatedParams      []resource.DataRelated
	)

	err = c.ShouldBind(&params)
	if err != nil {
		app.Error(c, -1, err, "参数绑定失败")
		return
	}

	// 查询现有的数据关联
	err = orm.Eloquent.Model(&resource.DataRelated{}).
		Where("source = ? and source_info_id = ?", params.Source, params.SourceInfoId).
		Find(&currentRelatedList).Error
	if err != nil {
		app.Error(c, -1, err, "查询现有的数据关联")
		return
	}
targetContinue:
	for _, target := range params.Target {
		for _, currentRelated := range currentRelatedList {
			if currentRelated.Source == params.Source &&
				currentRelated.Target == target &&
				currentRelated.SourceInfoId == params.SourceInfoId &&
				currentRelated.TargetInfoId == params.TargetInfoId {
				continue targetContinue
			}
		}
		relatedParams = append(relatedParams, resource.DataRelated{
			Source:       params.Source,
			Target:       target,
			SourceInfoId: params.SourceInfoId,
			TargetInfoId: params.TargetInfoId,
		})
	}

	if len(relatedParams) > 0 {
		// 若么有关联则直接新建关联
		err = orm.Eloquent.Create(&relatedParams).Error
		if err != nil {
			app.Error(c, -1, err, "创建资源关联失败")
			return
		}
	}

	app.OK(c, nil, "success")
}
