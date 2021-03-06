package mgo

import (
	"cetm/qapi/x/math"
	"http/web"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type IdMaker interface {
	Next() string
}

type Table struct {
	*UnsafeTable
}

var defaultIdMaker = &math.RandStringMaker{Prefix: "def", Length: 20}

func NewTable(db *Database, name string, IdMaker IdMaker) *Table {
	var t = &Table{UnsafeTable: NewUnsafeTable(db, name, IdMaker)}
	if t.IdMaker == nil {
		t.IdMaker = defaultIdMaker
	}
	return t
}

func (t *Table) Create(i IModel) error {
	i.BeforeCreate()
	i.SetID(t.IdMaker.Next())
	return t.UnsafeInsert(i)
}

func (t *Table) UpdateByID(id string, i IModel) error {
	i.BeforeUpdate()
	return t.UnsafeUpdateByID(id, i)
}

func (t *Table) MarkDelete(id string) error {
	var data = bson.M{
		"dtime": time.Now().Unix(),
	}
	return t.UnsafeUpdateByID(id, data)
}

func (t *Table) ReadAll(ptr interface{}) error {
	return t.UnsafeReadMany(bson.M{"dtime": 0}, ptr)
}

func (t *Table) ReadManyIn(key string, values []string, ptr interface{}) error {
	return t.UnsafeReadMany(bson.M{"dtime": 0, key: bson.M{"$in": values}}, ptr)
}

func (t *Table) ReadMany(where M, ptr interface{}) error {
	where["dtime"] = 0
	return t.UnsafeReadMany(where, ptr)
}

func (t *Table) ReadOne(where M, ptr interface{}) error {
	return t.UnsafeReadOne(where, ptr)
}

func (t *Table) ReadByID(id string, ptr interface{}) error {
	return t.UnsafeGetByID(id, ptr)
}

func (t *Table) NotExist(where M) error {
	where["dtime"] = 0
	var c, err = t.UnsafeTable.UnsafeCount(where)
	if err != nil {
		return err
	}
	if c > 0 {
		return web.BadRequest("already exist")
	}
	return nil
}

func (t *Table) ReadByArrID(ids []string, ptr interface{}) error {
	return t.UnsafeRunGetAll(bson.M{"_id": bson.M{"$in": ids}}, ptr)
}
