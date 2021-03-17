package portal

import "goskeleton/app/model"

func CreateServerItemFactory() *serverItem {
	return &serverItem{}
}

type serverItem struct {
}

func (s *serverItem) GetLaoHuoBanTreeList(args []model.LaoWuLaoLists) []LaoWuLaoLists {

	var laoWuLaoLists = make([]LaoWuLaoLists, 0)
	for _, obj := range args {
		//逐条处理数据
		var laowulaoitem LaoWuLaoLists
		if obj.StreetName != "" {
			laowulaoitem.StreetName = obj.StreetName
			// 内循环处理明细数据
			for index2, item := range args {
				if item.StreetName == obj.StreetName {
					var tmp LaoWuLaoItem

					tmp.ItemName = item.ProvidedServerTitle
					tmp.ItemNum = item.ProvidedServerPersonNum
					tmp.ItemUnit = item.ProvidedServerPersonUnit
					laowulaoitem.Detail = append(laowulaoitem.Detail, tmp)

					tmp.ItemName = item.AssistedAgedTitle
					tmp.ItemNum = item.AssistedAgedNum
					tmp.ItemUnit = item.AssistedAgedUnit
					laowulaoitem.Detail = append(laowulaoitem.Detail, tmp)

					tmp.ItemName = item.JoinOrgTitle
					tmp.ItemNum = item.JoinOrgNum
					tmp.ItemUnit = item.JoinOrgUnit
					laowulaoitem.Detail = append(laowulaoitem.Detail, tmp)

					tmp.ItemName = item.ProvidedServerNumTitle
					tmp.ItemNum = item.ProvidedServerNum
					tmp.ItemUnit = item.ProvidedServerUnit
					laowulaoitem.Detail = append(laowulaoitem.Detail, tmp)

					tmp.ItemName = item.ServerItemTitle
					tmp.ItemNum = item.ServerItemNum
					tmp.ItemUnit = item.ServerItemUnit
					laowulaoitem.Detail = append(laowulaoitem.Detail, tmp)

					args[index2].StreetName = ""
				}
			}
			laoWuLaoLists = append(laoWuLaoLists, laowulaoitem)
		}
	}

	return laoWuLaoLists
}
