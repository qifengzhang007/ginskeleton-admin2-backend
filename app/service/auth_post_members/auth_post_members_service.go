package auth_post_members

import (
	"github.com/qifengzhang007/sql_res_to_tree"
	"goskeleton/app/model/auth"
)

type AuthPostMembersService struct {
}

//根据用户ID获取所属的组织机构
func (a AuthPostMembersService) FindOrgs(id int64) []auth.OrgTree {
	orgs := auth.CreateAuthMenuAssignFactory("").GetAuthByUserId(int(id))
	var dest = make([]auth.OrgTree, 0)
	sql_res_to_tree.CreateSqlResFormatFactory().ScanToTreeData(orgs, &dest)
	return dest
}
