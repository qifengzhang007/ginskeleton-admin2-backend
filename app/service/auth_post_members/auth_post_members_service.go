package auth_post_members

import (
	"github.com/qifengzhang007/sql_res_to_tree"
	"go.uber.org/zap"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/auth"
)

type AuthPostMembersService struct {
}

//根据用户ID获取所属的组织机构
func (a AuthPostMembersService) FindOrgs(id int64) []auth.OrgTree {
	orgs := auth.CreateAuthMenuAssignFactory("").GetAuthByUserId(int(id))
	var dest = make([]auth.OrgTree, 0)
	err := sql_res_to_tree.CreateSqlResFormatFactory().ScanToTreeData(orgs, &dest)
	if err != nil {
		variable.ZapLog.Error("权限分析结果数据树形化出错：", zap.Error(err))
		return nil
	}
	return dest
}
