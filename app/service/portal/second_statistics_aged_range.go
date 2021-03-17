package portal

import (
	"goskeleton/app/model"
)

func CreateSecondPageFactory() *secondStatisticsAgedRange {

	return &secondStatisticsAgedRange{}
}

// 按照老年人年龄段统计数据，处理model 返回的数据

type secondStatisticsAgedRange struct {
}

// 1.基于年龄段统计老年人数据在街道的分布
func (s *secondStatisticsAgedRange) AgedRange(args []model.AgedRange) []StaticsRange {

	var ageRangeSlice []StaticsRange

	//两次循环，外层处理键（年龄标题），内层循环处理相关的 detail 数据
	for _, obj := range args {
		var ageRange = StaticsRange{
			Deatil: make([]ComunityStreet, 0),
		}
		if obj.AgedRanged != "" {
			ageRange.Title = obj.AgedRanged
			ageRange.OrderNo = obj.OrderNo
			for index2, obj2 := range args {
				if obj2.AgedRanged == obj.AgedRanged {
					var tmp ComunityStreet
					tmp.StreetName = obj2.CityAreaComunityStreetName
					tmp.Num = obj2.Num
					tmp.Unit = "人"
					ageRange.Deatil = append(ageRange.Deatil, tmp)
					// 不能用以下方法删除，会导致切片索引整体变化
					//args=append(args[:index2],args[index2+1:]...)
					args[index2].AgedRanged = ""
				}
			}
			// 获取一个年龄段对应的老人居住街道、数量明细
			ageRangeSlice = append(ageRangeSlice, ageRange)
		}
	}
	return ageRangeSlice
}

// 2.不同健康状态的老年人在各个街道的数据分布、统计
func (s *secondStatisticsAgedRange) HealthRange(args []model.HealthyStatics) []StaticsRange {

	var ageRangeSlice []StaticsRange

	//两次循环，外层处理键（年龄标题），内层循环处理相关的 detail 数据
	for _, obj := range args {
		var ageRange = StaticsRange{
			Deatil: make([]ComunityStreet, 0),
		}
		if obj.HealthyName != "" {
			ageRange.Title = obj.HealthyName
			ageRange.OrderNo = obj.OrderNo
			for index2, obj2 := range args {
				if obj2.HealthyName == obj.HealthyName {
					var tmp ComunityStreet
					tmp.StreetName = obj2.CityAreaComunityStreetName
					tmp.Num = obj2.Num
					tmp.Unit = "人"
					ageRange.Deatil = append(ageRange.Deatil, tmp)
					// 不能用以下方法删除，会导致切片索引整体变化
					//args=append(args[:index2],args[index2+1:]...)
					args[index2].HealthyName = ""
				}
			}
			// 获取一个年龄段对应的老人居住街道、数量明细
			ageRangeSlice = append(ageRangeSlice, ageRange)
		}
	}
	return ageRangeSlice
}

// 3.长护险按照时间段的长短分类、统计不同时间点的人数
func (s *secondStatisticsAgedRange) LongEnsurance(args []model.LongEnsurance) []LongEnsurance {
	// 将数据库数据处理成层次结构
	// 外层循环主要处理长时间段、内层循环处理该时间段对应的明细数据
	var longEnsuracneSlice []LongEnsurance
	for _, obj := range args {
		var longEnsuranceOne = LongEnsurance{
			Deatil: make([]Detail, 0),
		}
		longEnsuranceOne.DateCategory = obj.DateCategory
		longEnsuranceOne.OrderNo = obj.OrderNo
		if obj.DateCategory != "" {
			for index2, obj2 := range args {
				if obj.DateCategory == obj2.DateCategory {
					var detail Detail
					detail.PassAppDate = obj2.PassAppDatetime
					detail.Num = obj2.Num
					detail.Unit = "人"
					longEnsuranceOne.Deatil = append(longEnsuranceOne.Deatil, detail)
					args[index2].DateCategory = ""
				}
			}
			longEnsuracneSlice = append(longEnsuracneSlice, longEnsuranceOne)
		}
	}
	return longEnsuracneSlice
}
