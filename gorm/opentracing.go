package gorm

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/manmanxing/go_center_common/gls"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"strings"
)

//自定义插件
//https://gorm.io/zh_CN/docs/write_plugins.html
//https://gorm.io/zh_CN/docs/settings.html

//给 gorm 增加 opentracing
//instance 必填
func InitCallback(db *gorm.DB, instance string) {
	if len(strings.TrimSpace(instance)) <= 0 {
		panic("instance is empty")
	}
	db.InstantSet("instance",instance)
	db.Callback().Create().Before("gorm:create").Register("before_create", genBeforeOperation("mysql-create"))
	db.Callback().Create().After("gorm:create").Register("after_create", afterOperation)
	db.Callback().Query().Before("gorm:query").Register("before_query", genBeforeOperation("mysql-query"))
	db.Callback().Query().After("gorm:query").Register("after_query", afterOperation)
	db.Callback().Update().Before("gorm:update").Register("before_update", genBeforeOperation("mysql-update"))
	db.Callback().Update().After("gorm:update").Register("after_update", afterOperation)
	db.Callback().Delete().Before("gorm:delete").Register("before_delete", genBeforeOperation("mysql-delete"))
	db.Callback().Delete().After("gorm:delete").Register("after_delete", afterOperation)
	db.Callback().RowQuery().Before("gorm:row_query").Register("before_row_query", genBeforeOperation("mysql-row-query"))
	db.Callback().RowQuery().After("gorm:row_query").Register("after_row_query", afterOperation)
}

func genBeforeOperation(opeation string) func(scope *gorm.Scope) {
	return func(scope *gorm.Scope) {
		ctx := gls.TraceCtx()
		span, _ := opentracing.StartSpanFromContext(ctx, opeation, ext.SpanKindRPCClient)
		scope.InstanceSet("gorm_span", span)
	}
}

func afterOperation(scope *gorm.Scope) {
	tmp, ok := scope.InstanceGet("gorm_span")
	if !ok || tmp == nil {
		return
	}
	span := tmp.(opentracing.Span)
	ext.DBType.Set(span, "mysql")
	instance, _ := scope.DB().Get("instance")
	ext.DBInstance.Set(span, fmt.Sprintf("%v", instance))
	ext.DBStatement.Set(span, fmt.Sprintf("sql=#%v#,args=#%v#", scope.SQL, scope.SQLVars))
	span.Finish()
}