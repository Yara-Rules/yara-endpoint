package api

import (
	"math/rand"
	"time"

	"github.com/Yara-Rules/yara-endpoint/server/context"
	"github.com/Yara-Rules/yara-endpoint/server/models"
	"github.com/oklog/ulid"
)

func Index(ctx *context.Context) {
	ctx.HTML(200, "index")
}

func Dashboard(ctx *context.Context) {
	var data PublicDashboard

	eps := new([]models.Endpoint)
	rls := new([]models.Rule)

	ctx.DB.C(models.Endpoints).Find(nil).All(eps)
	ctx.DB.C(models.Rules).Find(nil).All(rls)

	data = PublicDashboard{
		Asset: (*eps),
		Rules: (*rls),
	}

	ctx.JSON(200, data)
}

func Commands(ctx *context.Context) {
	/* TODO: Check in the DB what are the avaliable commands
	 */
	ctx.JSON(200, struct {
		Commands []string `json:"commands"`
	}{Commands: []string{"ScanFile", "ScanDir", "ScanPID"}})
}

func newULID() ulid.ULID {
	t := time.Now()
	entropy := rand.New(rand.NewSource(t.UnixNano()))
	return ulid.MustNew(ulid.Timestamp(t), entropy)
}
