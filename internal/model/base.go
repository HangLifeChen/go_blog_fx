package model

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"go_blog/pkg/utils"
	"reflect"

	"gorm.io/gorm/schema"
	"gorm.io/plugin/soft_delete"
)

type Base struct {
	Id uint32 `json:"id" gorm:"primaryKey;autoIncrement;comment:primary key id"`
}

type BaseDelete struct {
	Id        uint32                `json:"id" gorm:"primaryKey;autoIncrement;comment:primary key id"`
	CreatedAt int64                 `json:"created_at" gorm:"comment:created time"`
	UpdatedAt int64                 `json:"updated_at" gorm:"comment:lastest update time"`
	DeletedAt soft_delete.DeletedAt `json:"deleted_at" gorm:"index;default:0;comment:deleted time"`
}

type BaseNoDelete struct {
	Id        uint32 `json:"id" gorm:"primaryKey;autoIncrement;comment:primary key id"`
	CreatedAt int64  `json:"created_at" gorm:"comment:created time"`
	UpdatedAt int64  `json:"updated_at" gorm:"comment:lastest update time"`
}

// todo 复习es数据库操作再回来研究这部分
type EncryptedString string

// ctx: contains request-scoped values
// field: the field using the serializer, contains GORM settings, struct tags
// dst: current model value, `user` in the below example
// dbValue: current field's value in database
func (es *EncryptedString) Scan(ctx context.Context, field *schema.Field, dst reflect.Value, dbValue interface{}) (err error) {
	switch value := dbValue.(type) {
	case []byte:
		*es = EncryptedString(value)
	case string:
		*es = EncryptedString(value)
	default:
		return fmt.Errorf("unsupported data %#v", dbValue)
	}
	return nil
}

// ctx: contains request-scoped values
// field: the field using the serializer, contains GORM settings, struct tags
// dst: current model value, `user` in the below example
// fieldValue: current field's value of the dst

func (es EncryptedString) Value(ctx context.Context, field *schema.Field, dst reflect.Value, fieldValue interface{}) (interface{}, error) {
	if len(es) == 0 {
		return "", nil
	}
	return utils.Base58(utils.Sha256(string(es))), nil
}

// custom type for JSON encoding and decoding
type StringArray []string

// implement GormDataType interface 实现接口将数据以json的形式存储到数据库中
func (s StringArray) GormDataType() string {
	return "json"
}

// mplement the Scan method of the Scanner interface // 实现Scanner接口 扫描取出数据StirngArray类型的数据会先将他转换为json
func (s *StringArray) Scan(value interface{}) error {
	if value == nil {
		*s = StringArray{}
		return nil
	}
	var data []byte
	switch v := value.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		return fmt.Errorf("unexpected type %T", value)
	}

	return json.Unmarshal(data, s)
}

// implement the Value method of the driver.Valuer interface。// 实现driver.Valuer接口 插入数据StirngArray类型数据时，将数据转换成json
func (s StringArray) Value() (driver.Value, error) {
	if s == nil {
		return "[]", nil
	}
	b, err := json.Marshal(s)
	return string(b), err
}
