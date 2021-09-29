package model

import (
	"fmt"
	"reflect"
	"time"
	"web-base/pkg/setting"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

type Model struct {
	ID         uint32 `gorm:"primary_key" json:"id,omitempty"`
	CreatedBy  string `json:"created_by,omitempty"`
	ModifiedBy string `json:"modified_by,omitempty"`
	CreatedOn  uint32 `json:"created_on,omitempty"`
	ModifiedOn uint32 `json:"modified_on,omitempty"`
	DeletedOn  uint32 `json:"deleted_on,omitempty"`
	IsDel      uint8  `json:"is_del,omitempty"`
}

func NewDBEngine(databaseSetting *setting.DatabaseSettingS) (*gorm.DB, error) {

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=%t&loc=Local",
		databaseSetting.UserName,
		databaseSetting.Password,
		databaseSetting.Host,
		databaseSetting.DBName,
		databaseSetting.Charset,
		databaseSetting.ParseTime,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		fmt.Println("open db failed")
		return nil, err
	}

	// if global.ServerSetting.RunMode == "debug" {
	// 	db.LogMode(true)
	// }

	// 注册公共回调
	db.Callback().Create().Before("gorm:create").Register("gorm:update_time", updateTimeStampForCreateCallback)
	db.Callback().Update().Before("gorm:update").Register("gorm:update_time", updateTimeStampForUpdateCallback)

	// todo delete bad
	//db.Callback().Delete().Before("gorm:delete").Register("gorm:delete_time", deleteCallback)

	sqlDB.SetMaxIdleConns(databaseSetting.MaxIdleConns)
	sqlDB.SetMaxOpenConns(databaseSetting.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(databaseSetting.ConnMaxLifetime)
	return db, nil
}

func updateTimeStampForCreateCallback(db *gorm.DB) {
	if db.Error == nil && db.Statement.Schema != nil {
		nowTime := time.Now().Unix()
		switch db.Statement.ReflectValue.Kind() {
		case reflect.Slice, reflect.Array:
			for i := 0; i < db.Statement.ReflectValue.Len(); i++ {
				rv := reflect.Indirect(db.Statement.ReflectValue.Index(i))
				field1 := db.Statement.Schema.FieldsByDBName["created_on"]
				if field1 != nil {
					field1.Set(rv, nowTime)
				}
				field := db.Statement.Schema.FieldsByDBName["modified_on"]
				if field != nil {
					field.Set(rv, nowTime)
				}
			}
		case reflect.Struct:
			if db.Statement.Schema.FieldsByDBName["created_on"] != nil {
				db.Statement.SetColumn("created_on", nowTime, true)
			}
			if db.Statement.Schema.FieldsByDBName["modified_on"] != nil {
				db.Statement.SetColumn("modified_on", nowTime)
			}
		}

	}
}

func updateTimeStampForUpdateCallback(db *gorm.DB) {
	if db.Error == nil && db.Statement.Schema != nil {
		db.Statement.SetColumn("modified_on", time.Now().Unix())
		db.Statement.AddClause(clause.Where{
			Exprs: []clause.Expression{clause.Eq{Column: "is_del", Value: 0}},
		})
	}
}

func deleteCallback(db *gorm.DB) {
	if db.Error == nil {
		if db.Statement.Schema != nil && !db.Statement.Unscoped {
			for _, c := range db.Statement.Schema.DeleteClauses {
				db.Statement.AddClause(c)
			}
		}

		if db.Statement.SQL.String() == "" {
			db.Statement.SQL.Grow(100)
			db.Statement.AddClauseIfNotExists(clause.Delete{})

			if db.Statement.Schema != nil {
				_, queryValues := schema.GetIdentityFieldValuesMap(db.Statement.ReflectValue, db.Statement.Schema.PrimaryFields)
				column, values := schema.ToQueryValues(db.Statement.Table, db.Statement.Schema.PrimaryFieldDBNames, queryValues)

				if len(values) > 0 {
					db.Statement.AddClause(clause.Where{Exprs: []clause.Expression{clause.IN{Column: column, Values: values}}})
				}

				if db.Statement.ReflectValue.CanAddr() && db.Statement.Dest != db.Statement.Model && db.Statement.Model != nil {
					_, queryValues = schema.GetIdentityFieldValuesMap(reflect.ValueOf(db.Statement.Model), db.Statement.Schema.PrimaryFields)
					column, values = schema.ToQueryValues(db.Statement.Table, db.Statement.Schema.PrimaryFieldDBNames, queryValues)

					if len(values) > 0 {
						db.Statement.AddClause(clause.Where{Exprs: []clause.Expression{clause.IN{Column: column, Values: values}}})
					}
				}
			}
			fmt.Println("test")
			db.Statement.AddClauseIfNotExists(clause.From{})
			db.Statement.Build("DELETE", "FROM", "WHERE")
		}

		if _, ok := db.Statement.Clauses["WHERE"]; !db.AllowGlobalUpdate && !ok && db.Error == nil {
			db.AddError(gorm.ErrMissingWhereClause)
			return
		}

		if !db.DryRun && db.Error == nil {
			if !db.Statement.Unscoped {
				db.Statement.Vars[0] = time.Now().Format("2006-01-02 15:04:05")
			}
			result, err := db.Statement.ConnPool.ExecContext(db.Statement.Context, db.Statement.SQL.String(), db.Statement.Vars...)

			if err == nil {
				db.RowsAffected, _ = result.RowsAffected()
			} else {
				db.AddError(err)
			}
		}
	}
}
