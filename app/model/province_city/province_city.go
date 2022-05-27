package province_city

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
	"goskeleton/app/utils/data_bind"
)

func CreateProvinceCityFactory(sqlType string) *ProvinceCityModel {
	return &ProvinceCityModel{BaseModel: model.BaseModel{DB: model.UseDbConn(sqlType)}}
}

type ProvinceCityModel struct {
	model.BaseModel
	Fid       int64  `json:"fid"`
	Name      string `json:"name"`
	Status    int    `json:"status"`
	Sort      int    `json:"sort"`
	NodeLevel int    `json:"node_level"`
	Remark    string `json:"remark"`
}

// 表名
func (p *ProvinceCityModel) TableName() string {
	return "tb_province_city"
}

func (p *ProvinceCityModel) GetCount(fid int, name string) (count int64) {
	p.Model(p).Select("id").Where("fid=? AND name like ?", fid, "%"+name+"%").Count(&count)
	return
}

//查询
func (p *ProvinceCityModel) List(name string, fid, limitStart, limit int) (counts int64, list []ProvinceCityModel) {
	if counts = p.GetCount(fid, name); counts > 0 {
		sql := `
			SELECT
			a.id,  a.fid,a.name ,a.node_level ,a.status ,a.sort ,a.remark ,a.created_at , a.updated_at
			FROM tb_province_city a
			WHERE   a.fid= ? AND   a.name LIKE  ? ORDER  BY a.sort DESC, a.fid ASC ,a.id  ASC
			LIMIT ? , ?
	`
		_ = p.Raw(sql, fid, "%"+name+"%", limitStart, limit).Find(&list)
	}
	return
}

// 根据fid查询子级节点全部数据
func (p *ProvinceCityModel) GetSubListByfid(fid int) []ProvinceCityTree {
	sql := `
		SELECT
		a.id,  a.fid,a.name ,a.node_level ,a.status ,a.sort ,a.remark , a.created_at , a.updated_at,
		(SELECT  CASE  WHEN  COUNT(*) =0 THEN 1 ELSE  0 END  FROM tb_province_city  WHERE  fid=a.id ) AS  is_leaf
		FROM tb_province_city a
		WHERE   fid= ?
	`
	var inSlice []ProvinceCityTree
	if res := p.Raw(sql, fid).Find(&inSlice); res.Error == nil && len(inSlice) > 0 {
		return inSlice
	} else if res.Error != nil {
		variable.ZapLog.Error("ProvinceCityModel 根据fid查询子级出错:", zap.Error(res.Error))
	}
	return nil
}

// 新增
func (p *ProvinceCityModel) InsertData(c *gin.Context) bool {
	var tmp ProvinceCityModel
	if err := data_bind.ShouldBindFormDataToModel(c, &tmp); err == nil {
		var counts int64
		if res := p.Model(p).Where("fid=? and name=?", tmp.Fid, tmp.Name).Count(&counts); res.Error == nil && counts > 0 {
			return false
		} else {
			if res := p.Create(&tmp); res.Error == nil {
				_ = p.updatePathInfoNodeLevel(int(tmp.Id))
				return true
			} else {
				variable.ZapLog.Error("ProvinceCityModel 数据新增出错：", zap.Error(res.Error))
			}
		}
	} else {
		variable.ZapLog.Error("ProvinceCityModel 数据从验证器绑定到model出错：", zap.Error(err))
	}
	return false
}

// 更新path_info 、node_level 字段
func (p *ProvinceCityModel) updatePathInfoNodeLevel(curItemid int) bool {
	sql := `
		UPDATE tb_province_city a  LEFT JOIN tb_province_city  b
		ON  a.fid=b.id
		SET  a.node_level=b.node_level+1,  a.path_info=CONCAT(b.path_info,',',a.id)
		WHERE  a.id=?
		`
	if res := p.Exec(sql, curItemid); res.Error == nil && res.RowsAffected >= 0 {
		return true
	} else {
		variable.ZapLog.Error("tb_province_city 更新 node_level , path_info 失败", zap.Error(res.Error))
	}
	return false
}

// 更新
func (p *ProvinceCityModel) UpdateData(c *gin.Context) bool {
	var tmp ProvinceCityModel

	if err := data_bind.ShouldBindFormDataToModel(c, &tmp); err == nil {
		// Omit 表示忽略指定字段(CreatedAt)，其他字段全量更新
		if res := p.Omit("CreatedAt").Save(&tmp); res.Error == nil {
			_ = p.updatePathInfoNodeLevel(int(tmp.Id))
		}
		return true
	} else {
		variable.ZapLog.Error("ProvinceCityModel 数据更新失败，错误详情：", zap.Error(err))
	}
	return false

}

// 删除
func (p *ProvinceCityModel) DeleteData(id int) bool {
	if res := p.Delete(p, id); res.Error == nil {
		return true
	} else {
		variable.ZapLog.Error("ProvinceCityModel 删除数据出错：", zap.Error(res.Error))
	}
	return false
}

// 查询该 id 是否存在子节点
func (p *ProvinceCityModel) HasSubNode(id int) (count int64) {
	p.Model(p).Select("id").Where("fid=?", id).Count(&count)
	return count
}
