package api

import (
	"github.com/Yara-Rules/yara-endpoint/server/context"
	"github.com/Yara-Rules/yara-endpoint/server/models"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func ShowErrors(ctx *context.Context) {
	errors := new([]models.Error)
	ctx.DB.C(models.Errors).Find(bson.M{"acknowledge": false}).All(errors)
	ctx.JSON(200, errors)
}

func ErrorDelete(ctx *context.Context) {
	id := ctx.Params(":id")
	if bson.IsObjectIdHex(id) {
		err := ctx.DB.C(models.Errors).Update(bson.ObjectIdHex(id), bson.M{"acknowledge": true})
		if err == mgo.ErrNotFound {
			ctx.JSON(400, Error{
				Error:    true,
				ErrorID:  1,
				ErrorMsg: "Rule not found",
			})
			return
		} else {
			ctx.JSON(200, Error{
				Error: false,
			})
			return
		}
	}
	ctx.JSON(400, Error{
		Error:    true,
		ErrorID:  1,
		ErrorMsg: "Not valid ID",
	})
}
