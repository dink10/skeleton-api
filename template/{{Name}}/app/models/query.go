package models

import (
	"time"
)

type Pagination struct {
	Page   int `json:"page" example:"1" binding:"required"`
	Limit  int `json:"limit" example:"20" binding:"required"`
	Count  int `json:"count" example:"1456" binding:"required"`
	Offset int `json:"-" example:"20"`
}

type DateRange struct {
	Gt time.Time `json:"gt" example:"2019-08-08T10:47:43.084943Z"`
	Lt time.Time `json:"lt" example:"2019-08-08T10:47:43.084943Z"`
}

type QParams struct {
	order       string
	pagination  Pagination
	filter      map[string]interface{}
	rangeFilter map[string]DateRange
}

func (qp QParams) GetOrder() string {
	return qp.order
}

func (qp QParams) GetPagination() Pagination {
	return qp.pagination
}

func (qp QParams) GetFilter() map[string]interface{} {
	return qp.filter
}

func (qp QParams) GetRange() map[string]DateRange {
	return qp.rangeFilter
}

func (qp *QParams) SetFilter(filter map[string]interface{}) {
	qp.filter = filter
}

func (qp *QParams) SetOrder(order string) {
	qp.order = order
}

func (qp *QParams) SetLimit(limit int) {
	qp.pagination.Limit = limit
}

func (qp *QParams) SetPagination(p Pagination) {
	qp.pagination = p
}

func (qp *QParams) SetRangeFilter(f map[string]DateRange) {
	qp.rangeFilter = f
}

type Query interface {
	GetOrder() string
	GetPagination() Pagination
	GetFilter() map[string]interface{}
	GetRange() map[string]DateRange
	SetFilter(filter map[string]interface{})
}

type ModelQuery interface {
	TableName() string
	FilteredFields() []string
}
