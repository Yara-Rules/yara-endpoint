package api

import (
	"fmt"
	"time"

	"github.com/Yara-Rules/yara-endpoint/server/context"
	"github.com/Yara-Rules/yara-endpoint/server/models"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func Assets(ctx *context.Context) {
	var m []models.Endpoint
	ctx.DB.C(models.Endpoints).Find(nil).All(&m)
	ctx.JSON(200, m)
}

func NewAsset(ctx *context.Context, newAsset NewAssetForm) {
	asset := models.Endpoint{
		ULID:          newULID().String(),
		Hostname:      newAsset.Hostname,
		ClientVersion: newAsset.ClientVersion,
		Tags:          newAsset.Tags,
		CreateAt:      time.Now(),
		UpdateAt:      time.Now(),
	}

	err := ctx.DB.C(models.Endpoints).Insert(&asset)
	if err != nil {
		ctx.JSON(400, Error{
			Error:    true,
			ErrorID:  1,
			ErrorMsg: fmt.Sprintf("%v", err),
		})
		return
	}
	ctx.JSON(200, Error{
		Error:    false,
		ErrorID:  0,
		ErrorMsg: "",
	})

}

func EditAsset(ctx *context.Context, newAsset EditAssetForm) {
	id := ctx.Params(":id")
	asset := models.Endpoint{}

	err := ctx.DB.C(models.Endpoints).Find(bson.M{"ulid": id}).One(&asset)
	if err == mgo.ErrNotFound {
		ctx.JSON(400, Error{
			Error:    true,
			ErrorID:  1,
			ErrorMsg: fmt.Sprintf("%v", err),
		})
		return
	}

	asset.Hostname = newAsset.Hostname
	asset.ClientVersion = newAsset.ClientVersion
	asset.Tags = newAsset.Tags
	asset.UpdateAt = time.Now()

	err = ctx.DB.C(models.Endpoints).Update(bson.M{"ulid": id}, &asset)
	if err != nil {
		ctx.JSON(400, Error{
			Error:    true,
			ErrorID:  1,
			ErrorMsg: fmt.Sprintf("%v", err),
		})
		return
	}
	ctx.JSON(200, Error{
		Error:    false,
		ErrorID:  0,
		ErrorMsg: "",
	})
}

func DeleteAsset(ctx *context.Context) {
	id := ctx.Params(":id")
	asset := models.Endpoint{}

	err := ctx.DB.C(models.Endpoints).Find(bson.M{"ulid": id}).One(&asset)
	if err == mgo.ErrNotFound {
		ctx.JSON(400, Error{
			Error:    true,
			ErrorID:  1,
			ErrorMsg: fmt.Sprintf("%v", err),
		})
		return
	}

	ctx.DB.C(models.Schedules).Remove(bson.M{"ulid": id})
	ctx.DB.C(models.Reports).Remove(bson.M{"ulid": id})
	ctx.DB.C(models.Endpoints).Remove(bson.M{"ulid": id})

	ctx.JSON(200, Error{
		Error:    false,
		ErrorID:  0,
		ErrorMsg: "Succesfully removed",
	})
}
