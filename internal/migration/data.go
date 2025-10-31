package migration

import (
	"go_blog/internal/model"

	"gorm.io/gorm"
)

func initData(db *gorm.DB) (err error) {
	if err = initRepositoriesTag(db); err != nil {
		return err
	}
	if err = initAdmin(db); err != nil {
		return err
	}
	if err = initVersion(db); err != nil {
		return err
	}
	return
}

func initRepositoriesTag(db *gorm.DB) (err error) {
	data := []model.RepositoriesTag{
		{Base: model.Base{Id: 1}, Name: "Multimodal", Pid: 0},
		{Base: model.Base{Id: 2}, Name: "Audio-Text-to-Text", Pid: 1},
		{Base: model.Base{Id: 3}, Name: "Image-Text-to-Text", Pid: 1},
		{Base: model.Base{Id: 4}, Name: "Visual Question Answering", Pid: 1},
		{Base: model.Base{Id: 5}, Name: "Document Question Answering", Pid: 1},
		{Base: model.Base{Id: 6}, Name: "Video-Text-to-Text", Pid: 1},
		{Base: model.Base{Id: 7}, Name: "Any-to-Any", Pid: 1},
		{Base: model.Base{Id: 8}, Name: "Computer Vision", Pid: 0},
		{Base: model.Base{Id: 9}, Name: "Depth Estimation", Pid: 8},
		{Base: model.Base{Id: 10}, Name: "Image Classification", Pid: 8},
		{Base: model.Base{Id: 11}, Name: "Object Detection", Pid: 8},
		{Base: model.Base{Id: 12}, Name: "Image Segmentation", Pid: 8},
		{Base: model.Base{Id: 13}, Name: "Text-to-Image", Pid: 8},
		{Base: model.Base{Id: 14}, Name: "Image-to-Text", Pid: 8},
		{Base: model.Base{Id: 15}, Name: "Image-to-Image", Pid: 8},
		{Base: model.Base{Id: 16}, Name: "Unconditional Image Generation", Pid: 8},
		{Base: model.Base{Id: 17}, Name: "Video Classification", Pid: 8},
		{Base: model.Base{Id: 18}, Name: "Text-to-Video", Pid: 8},
		{Base: model.Base{Id: 19}, Name: "Zero-Shot Image Classification", Pid: 8},
		{Base: model.Base{Id: 20}, Name: "Mask Generation", Pid: 8},
		{Base: model.Base{Id: 21}, Name: "Zero-Shot Object Detection", Pid: 8},
		{Base: model.Base{Id: 22}, Name: "Text-to-3D", Pid: 8},
		{Base: model.Base{Id: 23}, Name: "Image-to-3D", Pid: 8},
		{Base: model.Base{Id: 24}, Name: "Image Feature Extraction", Pid: 8},
		{Base: model.Base{Id: 25}, Name: "Keypoint Detection", Pid: 8},
		{Base: model.Base{Id: 26}, Name: "Natural Language Processing", Pid: 0},
		{Base: model.Base{Id: 27}, Name: "Text Classification", Pid: 26},
		{Base: model.Base{Id: 28}, Name: "Token Classification", Pid: 26},
		{Base: model.Base{Id: 29}, Name: "Table Question Answering", Pid: 26},
		{Base: model.Base{Id: 30}, Name: "Question Answering", Pid: 26},
		{Base: model.Base{Id: 31}, Name: "Zero-Shot Classification", Pid: 26},
		{Base: model.Base{Id: 32}, Name: "Translation", Pid: 26},
		{Base: model.Base{Id: 33}, Name: "Summarization", Pid: 26},
		{Base: model.Base{Id: 34}, Name: "Feature Extraction", Pid: 26},
		{Base: model.Base{Id: 35}, Name: "Text Generation", Pid: 26},
		{Base: model.Base{Id: 36}, Name: "Text2Text Generation", Pid: 26},
		{Base: model.Base{Id: 37}, Name: "Fill-Mask", Pid: 26},
		{Base: model.Base{Id: 38}, Name: "Sentence Similarity", Pid: 26},
		{Base: model.Base{Id: 39}, Name: "Audio", Pid: 0},
		{Base: model.Base{Id: 40}, Name: "Text-to-Speech", Pid: 39},
		{Base: model.Base{Id: 41}, Name: "Text-to-Audio", Pid: 39},
		{Base: model.Base{Id: 42}, Name: "Automatic Speech Recognition", Pid: 39},
		{Base: model.Base{Id: 43}, Name: "Audio-to-Audio", Pid: 39},
		{Base: model.Base{Id: 44}, Name: "Audio Classification", Pid: 39},
		{Base: model.Base{Id: 45}, Name: "Voice Activity Detection", Pid: 39},
		{Base: model.Base{Id: 46}, Name: "Tabular", Pid: 0},
		{Base: model.Base{Id: 47}, Name: "Tabular Classification", Pid: 46},
		{Base: model.Base{Id: 48}, Name: "Tabular Regression", Pid: 46},
		{Base: model.Base{Id: 49}, Name: "Time Series Forecasting", Pid: 46},
		{Base: model.Base{Id: 50}, Name: "Reinforcement Learning", Pid: 0},
		{Base: model.Base{Id: 51}, Name: "Reinforcement Learning", Pid: 50},
		{Base: model.Base{Id: 52}, Name: "Robotics", Pid: 50},
		{Base: model.Base{Id: 53}, Name: "Other", Pid: 0},
		{Base: model.Base{Id: 54}, Name: "Graph Machine Learning", Pid: 53},
	}
	for k := range data {
		//不使用where gorm 会直接使用 &data（结构体非零值字段）作为查询条件
		err = db.Model(&model.RepositoriesTag{}).FirstOrCreate(&data[k]).Error
		if err != nil {
			return
		}
	}
	return
}

// 仅新建时设置
// db.Where("name = ?", "Tom").
//     Attrs(User{Age: 18}).
//     FirstOrCreate(&user)

// 复制代码
// 不管是否存在都设置 Age
// db.Where("name = ?", "Tom").
//     Assign(User{Age: 18}).
//     FirstOrCreate(&user)

func initAdmin(db *gorm.DB) (err error) {
	data := model.Admin{
		Username:    "admin",
		Password:    model.EncryptedString("2n0e2b5u1l7a4i4ope5nc5om0put1e"),
		Nickname:    "administrator",
		Avatar:      "",
		LastLoginIp: "127.0.0.1",
	}
	err = db.Model(&model.Admin{}).FirstOrCreate(&data).Error
	if err != nil {
		return
	}
	return
}

func initVersion(db *gorm.DB) (err error) {
	data := model.Version{
		Base:        model.Base{Id: 1},
		Version:     "1.0.0",
		RefreshTime: 60,
		ClearLocal:  false,
	}
	err = db.Model(&model.Version{}).FirstOrCreate(&data).Error
	if err != nil {
		return
	}
	return
}
