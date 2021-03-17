package model

// 创建 userFactory
// 参数说明： 传递空值，默认使用 配置文件选项：UseDbType（mysql）
func CreateFaceFactory(sqlType string) *FaceModel {
	return &FaceModel{BaseModel: BaseModel{DB: UseDbConn(sqlType)}}
}

type FaceModel struct {
	BaseModel
	ImgPath           string `json:"img_path"`
	EdgeRecognizeName string `json:"edge_recognize_name"`
	ResidentName      string `json:"resident_name"`
	Gender            string `json:"gender"`
	CardNo            string `json:"card_no"`
	Property          string `json:"property"`
	EventTime         int    `json:"event_time"`
	OutTime           int    `json:"out_time"`
	InTime            int    `json:"in_time"`
	ResidentId        int    `json:"resident_id"`
	CameraName        string `json:"camera_name"`
}

type FaceModelTwo struct {
	BaseModel
	ImgPath           string `json:"img_path"`
	EdgeRecognizeName string `json:"edge_recognize_name"`
	ResidentName      string `json:"resident_name"`
	Gender            string `json:"gender"`
	CardNo            string `json:"card_no"`
	Property          string `json:"property"`
	EventTime         string `json:"event_time"`
	OutTime           string `json:"out_time"`
	InTime            string `json:"in_time"`
	ResidentId        int    `json:"resident_id"`
	CameraName        string `json:"camera_name"`
}

// 表名
func (f *FaceModel) TableName() string {
	return "tb_face_record"
}

func (f *FaceModel) Add(data *FaceModel) (err error) {
	var db = UseDbConn("mysql")
	model := db.Model(f)
	err = model.Create(data).Error
	return
}

func (f *FaceModel) FindByResidentId(resident_id int) *FaceModel {
	face_data := &FaceModel{}
	var db = UseDbConn("mysql")
	model := db.Model(f)
	model.Where("resident_id = ?", resident_id).
		Order("event_time desc").First(face_data)
	return face_data
}

func (f *FaceModel) Edit(event_time int, resident_id int64, c_name string) (err error) {
	var db = UseDbConn("mysql")
	model := db.Model(f)
	err = model.Where("id = ?", resident_id).Update("event_time", event_time).Error
	err = model.Where("id = ?", resident_id).Update("out_time", event_time).Error
	err = model.Where("id = ?", resident_id).Update("camera_name", c_name).Error

	return
}

func (f *FaceModel) FindTen() *[]FaceModelTwo {
	face_data := &[]FaceModelTwo{}
	var db = UseDbConn("mysql")
	model := db.Model(f)
	model.Select("img_path,edge_recognize_name,resident_name," +
		"CASE in_time WHEN 0 THEN '' ELSE CONCAT(FROM_UNIXTIME(in_time,'%Y-%m-%d %H:%i:%s'),'') END in_time," +
		"CASE out_time WHEN 0 THEN '' ELSE CONCAT(FROM_UNIXTIME(out_time,'%Y-%m-%d %H:%i:%s'),'') END out_time," +
		"CASE event_time WHEN 0 THEN '' ELSE CONCAT(FROM_UNIXTIME(event_time,'%Y-%m-%d %H:%i:%s'),'') END event_time,gender,property,camera_name,card_no").
		Order("event_time desc").Limit(10).Scan(face_data)
	return face_data
}
