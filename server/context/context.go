package context

import (
	"github.com/Yara-Rules/yara-endpoint/server/database"
	log "github.com/sirupsen/logrus"
	macaron "gopkg.in/macaron.v1"
	mgo "gopkg.in/mgo.v2"
)

type Context struct {
	*macaron.Context
	DBSess *mgo.Session
	DB     *mgo.Database
}

func (ctx *Context) Handle(status int, title string, err error) {
	if err != nil {
		if macaron.Env != macaron.PROD {
			ctx.Data["ErrorMsg"] = err
		}
	}

	switch status {
	case 404:
		ctx.Data["Title"] = "Page Not Found"
	case 500:
		ctx.Data["Title"] = "Internal Server Error"
	}
	// ctx.HTML(status, base.TplName(fmt.Sprintf("status/%d", status)))
}

func Contexter() macaron.Handler {
	return func(c *macaron.Context) {
		ctx := &Context{
			Context: c,
		}

		sess, err := database.GetSession()
		if err != nil {
			log.Fatalf("Unable to get the database session: %v", err)
			ctx.Handle(500, "Get DbSession", err)
			return
		}
		ctx.DBSess = sess
		ctx.DB = database.GetDb(sess)

		c.Map(ctx)
	}
}
