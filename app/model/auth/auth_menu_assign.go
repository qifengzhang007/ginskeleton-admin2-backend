package auth

import (
	"go.uber.org/zap"
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
	"strconv"
	"strings"
)

func CreateAuthMenuAssignFactory(sqlType string) *AuthMenuAssignModel {
	return &AuthMenuAssignModel{BaseModel: model.BaseModel{DB: model.UseDbConn(sqlType)}}
}

type AuthMenuAssignModel struct {
	model.BaseModel
}

// 待分配的系统菜单、按钮 数据列表
func (a *AuthMenuAssignModel) GetSystemMenuButtonList() (counts int64, data []AuthSystemMenuButton) {
	sql := `
		SELECT a.id  AS  system_menu_id,a.fid  AS  system_menu_fid,a.title,
		IFNULL( b.fr_auth_system_menu_id,0) AS fr_auth_system_menu_id,
		(CASE WHEN  IFNULL(b.fr_auth_system_menu_id,0)=0 THEN 'menu' ELSE 'button' END  ) AS node_type,
		IFNULL( c.id,0) AS button_id,
		IFNULL(c.cn_name,'') AS button_name,(CASE WHEN a.fid=0 THEN 1 ELSE 0 END)  AS  expand
		FROM
		tb_auth_system_menu a LEFT   JOIN tb_auth_system_menu_button  b ON a.id=b.fr_auth_system_menu_id
		LEFT JOIN  tb_auth_button_cn_en  c  ON  b.fr_auth_button_cn_en_id=c.id
		ORDER   BY a.sort DESC, a.id ASC, a.fid ASC,c.id ASC
			`
	if res := a.Raw(sql).Find(&data); res.Error == nil && res.RowsAffected > 0 {
		return res.RowsAffected, data
	} else {
		variable.ZapLog.Error("查询系统待分配菜单出错：", zap.Error(res.Error))
	}

	return 0, nil
}

// 已分配给部门、岗位的系统菜单、按钮
func (a *AuthMenuAssignModel) GetAssignedMenuButtonList(orgPostId int) (counts int64, data []AssignedSystemMenuButton) {
	sql := `
		SELECT  
		a.fr_auth_orgnization_post_id  AS org_post_id,
		b.id AS  system_menu_id,b.fid AS system_menu_fid, b.title,
		a.id AS post_mount_has_menu_id  ,
		IFNULL(c.fr_auth_post_mount_has_menu_id,0) AS fr_mount_has_menu_id  , 
		IFNULL(c.id,0) AS  post_mount_has_menu_button_id,
		(CASE WHEN c.id>0 THEN 'button' ELSE 'menu' END) AS node_type,
		IFNULL(d.cn_name,'')  AS  button_name,
		0 as checked,1 as expand
		FROM 
		tb_auth_post_mount_has_menu  a  LEFT JOIN tb_auth_system_menu  b  ON  a.fr_auth_system_menu_id=b.id
		LEFT JOIN tb_auth_post_mount_has_menu_button  c  ON  a.id=c.fr_auth_post_mount_has_menu_id
		LEFT  JOIN tb_auth_button_cn_en  d  ON  c.fr_auth_button_cn_en_id=d.id
		WHERE  a.status=1
		AND  a.fr_auth_orgnization_post_id=?
		ORDER BY  b.sort DESC, b.id ASC, b.fid ASC,d.id ASC
			`
	if res := a.Raw(sql, orgPostId).Find(&data); res.Error == nil && res.RowsAffected > 0 {
		return res.RowsAffected, data
	}
	return 0, nil
}

// 给组织机构（部门、岗位）分配菜单权限
func (a *AuthMenuAssignModel) AssginAuthForOrg(orgId, systemMenuId, systemMenuFid, buttonId int, nodeType string) (assginRes bool) {
	// 权限分配模块
	// 如果在前端界面一次性批量勾线上百条节点同时分配，前端会并发提交，后台sql执行时可能会遇见死锁状态发生（insert into 时发生了死锁）
	// 这里出现死锁时，需要尝试重新执行sql 《高性能mysql》这个本书上有介绍，死锁在并发高的场景下很难避免，尝试重新执行sql是一种解决方案，其他解决方式请自行百度了解
	var failTryTimes = 1

	assginRes = true
	sql := `INSERT  INTO tb_auth_post_mount_has_menu(fr_auth_orgnization_post_id,fr_auth_system_menu_id)
			SELECT ?,? FROM  DUAL  WHERE   NOT EXISTS(SELECT 1 FROM tb_auth_post_mount_has_menu a  WHERE  a.fr_auth_orgnization_post_id=? AND a.fr_auth_system_menu_id=? FOR UPDATE)
			`
	//1.如果fid>0 增加一条父级菜单分配记录
	if systemMenuFid > 0 {
		if res := a.Exec(sql, orgId, systemMenuFid, orgId, systemMenuFid); res.Error != nil {
			variable.ZapLog.Error("tb_auth_post_mount_has_menu  表分配菜单父级失败", zap.Error(res.Error))
			return false
		}
	}
	//2.当前菜单增加一条分配记录
	if res := a.Exec(sql, orgId, systemMenuId, orgId, systemMenuId); res.Error == nil {
		if nodeType == "button" {
			sql = "select id from tb_auth_post_mount_has_menu where fr_auth_orgnization_post_id=? AND fr_auth_system_menu_id=? AND   status=1 "
			var temId int
			if res = a.Raw(sql, orgId, systemMenuId).First(&temId); res.Error == nil && temId > 0 {
				sql = `
					INSERT  INTO tb_auth_post_mount_has_menu_button(fr_auth_post_mount_has_menu_id,fr_auth_button_cn_en_id)
					SELECT ?,? FROM  DUAL  WHERE   NOT EXISTS(SELECT 1 FROM tb_auth_post_mount_has_menu_button a  WHERE  a.fr_auth_post_mount_has_menu_id=? AND a.fr_auth_button_cn_en_id=? FOR UPDATE)
					`
				if buttonId > 0 {
				label1:
					if res = a.Exec(sql, temId, buttonId, temId, buttonId); res.Error == nil {
						// 继续分配接口的访问权限(casbin_rules写入相关数据)
						var lastID int
						sql = "SELECT id  FROM tb_auth_post_mount_has_menu_button where  fr_auth_post_mount_has_menu_id=?  AND fr_auth_button_cn_en_id=?"
						if res = a.Raw(sql, temId, buttonId).First(&lastID); res.Error == nil {
							assginRes = a.AssginCasbinAuthPolicyToOrg(lastID, nodeType)
						}
					} else {
						// insert into 执行时遇见死锁状态，尝试重新执行，最大允许五次尝试，否则就记录错误
						if failTryTimes <= 5 {
							failTryTimes++
							goto label1
						}
						variable.ZapLog.Error("tb_auth_post_mount_has_menu_button  表分配按钮失败", zap.Error(res.Error))
						assginRes = false
					}
				}
			} else {
				variable.ZapLog.Error("tb_auth_post_mount_has_menu_button  表分配按钮失败", zap.Error(res.Error))
			}
		}
	} else {
		variable.ZapLog.Error("tb_auth_post_mount_has_menu  表分配菜单失败", zap.Error(res.Error))
		assginRes = false
	}
	return assginRes
}

// 从组织机构（部门、岗位）删除权限
func (a *AuthMenuAssignModel) DeleteAuthFromOrg(postMountHasMenuId, postMountHasMenuButtonId int, nodeType string) bool {
	if nodeType == "menu" {
		sql := "DELETE   FROM tb_auth_post_mount_has_menu WHERE  id=?"
		if res := a.Exec(sql, postMountHasMenuId); res.Error == nil {
			return true
		}
	} else if nodeType == "button" {
		sql := "DELETE   FROM tb_auth_post_mount_has_menu_button WHERE  id=?"
		if res := a.Exec(sql, postMountHasMenuButtonId); res.Error == nil {

			return a.DeleteCasbibRules(postMountHasMenuButtonId, nodeType)
		}
	}
	return false
}

//  删除 casbin 表接口已分配的权限
func (a *AuthMenuAssignModel) DeleteCasbibRules(authPostMountHasMenuButtonId int, nodeType string) (resBool bool) {
	resBool = true
	if nodeType == "button" {
		sql := "DELETE FROM tb_auth_casbin_rule  WHERE fr_auth_post_mount_has_menu_button_id=? AND ptype='p' "
		if res := a.Exec(sql, authPostMountHasMenuButtonId); res.Error != nil {
			// 角色继承关系暂时不删除，只要删除相关的节点权限即可
			variable.ZapLog.Error("AuthMenuAssignModel 删除casbin权限失败" + res.Error.Error())
			resBool = false
		}
	}
	return
}

// 给组织机构节点分配casbin的policy策略权限
func (a *AuthMenuAssignModel) AssginCasbinAuthPolicyToOrg(authPostMountHasMenuButtonId int, nodeType string) (resBool bool) {
	// 参见 69 行注释
	var failTryTimes = 1
	resBool = true
	// 分配了按钮，才可以同步分配按钮对应的路由接口
	if nodeType == "button" {
		// 首先给组织机构分配p权限(polic权限)
		sql := `
		SELECT   
		'p' as ptype,b.fr_auth_orgnization_post_id ,c.request_url,UPPER(c.request_method)  AS request_method ,
		a.id AS auth_post_mount_has_menu_button_id , b.id   AS   auth_post_mount_has_menu_id
		FROM  tb_auth_post_mount_has_menu_button  a ,tb_auth_post_mount_has_menu b   ,tb_auth_system_menu_button c
		WHERE   a.id=?  AND a.fr_auth_post_mount_has_menu_id = b.id
		AND  c.fr_auth_system_menu_id  = b.fr_auth_system_menu_id  AND  c.fr_auth_button_cn_en_id  =a.fr_auth_button_cn_en_id
		`
		var tmp struct {
			Ptype                        string
			FrAuthOrgnizationPostId      int
			RequestUrl                   string
			RequestMethod                string
			AuthPostMountHasMenuButtonId int
		}
		if res := a.Raw(sql, authPostMountHasMenuButtonId).First(&tmp); res.Error == nil {
			sql = `
			INSERT  INTO tb_auth_casbin_rule(ptype,v0,v1,v2,fr_auth_post_mount_has_menu_button_id,v3,v4,v5)
			SELECT  ?,?,?,?,?,'','',''  FROM   DUAL 
			WHERE NOT  EXISTS(SELECT 1 FROM tb_auth_casbin_rule a WHERE  a.ptype=? AND  a.v0=? AND  a.v1=? AND  a.v2=? FOR UPDATE)
			`
		label1:
			if res = a.Exec(sql, tmp.Ptype, tmp.FrAuthOrgnizationPostId, tmp.RequestUrl, tmp.RequestMethod, tmp.AuthPostMountHasMenuButtonId, tmp.Ptype, tmp.FrAuthOrgnizationPostId, tmp.RequestUrl, tmp.RequestMethod); res.Error == nil {
				// 为当前节点继续分配g(group权限，设置部门继承关系)
				return a.AssginCasbinAuthGroupToOrg(tmp.FrAuthOrgnizationPostId)
			} else {
				if failTryTimes <= 5 {
					failTryTimes++
					goto label1
				}
				resBool = false
				variable.ZapLog.Error("AuthMenuAssignModel 发生错误", zap.Error(res.Error))
			}
		}
	}
	return resBool
}

// 给组织机构节点分配casbin的group（角色继承关系权限）
func (a *AuthMenuAssignModel) AssginCasbinAuthGroupToOrg(orgId int) (resBool bool) {
	// 参见 69 行注释
	var failTryTimes = 1
	resBool = true
	sql := "SELECT path_info  FROM  tb_auth_organization_post WHERE   id =?"
	var pathInfo string
	if res := a.Raw(sql, orgId).First(&pathInfo); res.Error == nil {
		if len(pathInfo) > 0 {
			orgIdArray := strings.Split(pathInfo, ",")
			orgLen := len(orgIdArray)
			sql = `
				INSERT   INTO tb_auth_casbin_rule (ptype,v0,v1,v2,v3,v4,v5) 
				SELECT   'g',?,?,'','','',''  FROM   DUAL   
				WHERE   NOT  EXISTS(SELECT 1 FROM tb_auth_casbin_rule a WHERE a.ptype='g' AND v0=? AND  v1=? FOR UPDATE )
				`
			var lastId = 0
			var id = 0
			var err error
			for i := 1; i <= orgLen; i++ {
				// 遍历组织机构id
				if id, err = strconv.Atoi(orgIdArray[orgLen-i]); err == nil && i > 1 && id > 0 {
				label:
					if res = a.Exec(sql, lastId, id, lastId, id); res.Error != nil {
						if failTryTimes <= 5 {
							failTryTimes++
							goto label
						}
						variable.ZapLog.Error("AuthMenuAssignModel 批量插入角色继承关系时出错", zap.Error(res.Error))
						resBool = false
					}
				}
				lastId = id
			}
		}
	} else {
		resBool = false
	}
	return resBool
}

// 根据用户id查询已经分配的菜单
func (a *AuthMenuAssignModel) GetAuthByUserId(userId int) (OrgTree []OrgTree) {
	sql := `
		SELECT GROUP_CONCAT(b.path_info) id
		FROM
		tb_auth_post_members a
		LEFT JOIN tb_auth_organization_post b
		ON a.fr_auth_organization_post_id = b.id
		WHERE a.status=1 AND a.fr_user_id = ?
		GROUP BY a.fr_user_id
	`
	var orgPathInfo string
	if res := a.Raw(sql, userId).First(&orgPathInfo); res.Error != nil {
		variable.ZapLog.Error("查询用户所在的岗位、部门、公司节点出错：" + res.Error.Error())
		return nil
	} else if len(orgPathInfo) == 0 {
		return nil
	}
	sql = `
		SELECT   
			c.id,c.fid AS org_fid,
			c.title AS  org_title,  'dept' AS   node_type,
			1 AS expand
			FROM tb_auth_organization_post c
			WHERE FIND_IN_SET(c.id,? ) AND c.status=1
			UNION  
			SELECT DISTINCT  
			d.id *100, 
			-- e.id  as   menu_id,
			(CASE WHEN e.fid=0 THEN d.fr_auth_orgnization_post_id  ELSE  
			(SELECT  id*100  FROM  tb_auth_post_mount_has_menu  WHERE  fr_auth_orgnization_post_id=d.fr_auth_orgnization_post_id  AND  fr_auth_system_menu_id=e.fid  LIMIT 1 ) 
			END)  AS   fid ,
			e.title ,
			'menu' AS   node_type,
			(CASE WHEN  e.fid =0  THEN 1  ELSE   0 END) AS expand
			FROM  tb_auth_post_mount_has_menu  d ,tb_auth_system_menu  e 
			WHERE  
			FIND_IN_SET(d.fr_auth_orgnization_post_id,?)
			AND
			d.fr_auth_system_menu_id=e.id 
			UNION
			SELECT
			100000 AS  button_id ,
			f.fr_auth_post_mount_has_menu_id*100 AS  fid ,
			g.cn_name AS  button_name,
			'button' AS  node_type,
			0 AS expand
			FROM  
			tb_auth_post_mount_has_menu_button f  ,
			tb_auth_post_mount_has_menu d ,
			tb_auth_button_cn_en  g
			WHERE  f.fr_auth_post_mount_has_menu_id=d.id
			AND  d.status=1 AND f.status=1 AND 
			g.id=f.fr_auth_button_cn_en_id
			AND  FIND_IN_SET(d.fr_auth_orgnization_post_id,?)
			`
	a.Raw(sql, orgPathInfo, orgPathInfo, orgPathInfo).Scan(&OrgTree)
	return
}
