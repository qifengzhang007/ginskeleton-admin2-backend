package data_bind

import (
	"errors"
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"reflect"
)

const (
	modelStructMustPtr = "modelStruct 必须传递一个指针"
)

// 绑定form表单验证器已经验证完成的参数到 model 结构体,绑定的依据为表单参数验证器和model拥有想用的字段名、数据类型，并且设置了相同的 json 标签
// mode 结构体支持匿名嵌套
// 数据绑定原则： model 定义的结构体字段和表单验证器结构体设置的json标签名称、数据类型一致，才可以绑定

func ShouldBindFormDataToModel(c *gin.Context, modelStrcut interface{}) error {
	mTypeOf := reflect.TypeOf(modelStrcut)
	if mTypeOf.Kind() != reflect.Ptr {
		return errors.New(modelStructMustPtr)
	}
	mValueOf := reflect.ValueOf(modelStrcut)

	//分析 modelStruct 字段
	mValueOfEle := mValueOf.Elem()
	mtf := mValueOf.Elem().Type()
	fieldNum := mtf.NumField()
	for i := 0; i < fieldNum; i++ {
		if !mtf.Field(i).Anonymous {
			fieldSetValue(c, mValueOfEle, mtf, i)
		} else if mtf.Field(i).Type.Kind() == reflect.Struct {
			//处理结构体(有名+匿名)
			mValueOfEle.Field(i).Set(analysisAnonymousStruct(c, mValueOfEle.Field(i)))
		}
	}
	return nil
}

// 分析结构体（包括匿名的）,并且获取结构体的值
func analysisAnonymousStruct(c *gin.Context, value reflect.Value) reflect.Value {

	typeOf := value.Type()
	fieldNum := typeOf.NumField()
	newStruct := reflect.New(typeOf)
	newStructElem := newStruct.Elem()
	for i := 0; i < fieldNum; i++ {
		fieldSetValue(c, newStructElem, typeOf, i)
	}
	return newStructElem
}

// 为结构体字段赋值
func fieldSetValue(c *gin.Context, valueOf reflect.Value, typeOf reflect.Type, colIndex int) {
	relaKey := typeOf.Field(colIndex).Tag.Get("json")
	if relaKey != "-" {
		relaKey = consts.ValidatorPrefix + typeOf.Field(colIndex).Tag.Get("json")
		switch typeOf.Field(colIndex).Type.Kind() {
		case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
			valueOf.Field(colIndex).SetInt(int64(c.GetFloat64(relaKey)))
		case reflect.Float32, reflect.Float64:
			valueOf.Field(colIndex).SetFloat(c.GetFloat64(relaKey))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			valueOf.Field(colIndex).SetUint(uint64(c.GetFloat64(relaKey)))
		case reflect.String:
			valueOf.Field(colIndex).SetString(c.GetString(relaKey))
		case reflect.Bool:
			valueOf.Field(colIndex).SetBool(c.GetBool(relaKey))
		default:
			// model 数据库模型日期时间字段请统一设置为 字符串即可
		}
	}
}
