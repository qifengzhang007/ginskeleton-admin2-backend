package test

import (
	"goskeleton/app/model/auth"
	"testing"
)

//tb_casbin_rule 表的单元测试

func TestInsert(t *testing.T) {
	res := auth.CreateCasbinRuleFactory("").InsertData("p", "1", "/admin/users/test/abcd", "GET", "", "", "")
	if res {
		t.Log("单元测试通过")
	} else {
		t.Errorf("单元测试失败")
	}
}

func TestUpdate(t *testing.T) {
	res := auth.CreateCasbinRuleFactory("").UpdateData(8, "p", "1", "/admin/users/test/def", "POST", "", "", "")
	if res {
		t.Log("单元测试通过")
	} else {
		t.Errorf("单元测试失败")
	}
}

func TestDelete(t *testing.T) {
	res := auth.CreateCasbinRuleFactory("").DeleteData(8)
	if res {
		t.Log("单元测试通过")
	} else {
		t.Errorf("单元测试失败")
	}
}
