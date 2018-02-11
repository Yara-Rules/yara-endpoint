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

		// Get user from session if logged in.
		// ctx.User, ctx.IsBasicAuth = auth.SignedInUser(ctx.Context, ctx.Session)

		// if ctx.User != nil {
		//  ctx.IsSigned = true
		//  ctx.Data["IsSigned"] = ctx.IsSigned
		//  ctx.Data["SignedUser"] = ctx.User
		//  ctx.Data["SignedUserID"] = ctx.User.ID
		//  ctx.Data["SignedUserName"] = ctx.User.Name
		//  ctx.Data["IsAdmin"] = ctx.User.IsAdmin
		// } else {
		//  ctx.Data["SignedUserID"] = 0
		//  ctx.Data["SignedUserName"] = ""
		// }

		// If request sends files, parse them here otherwise the Query() can't be parsed and the CsrfToken will be invalid.
		// if ctx.Req.Method == "POST" && strings.Contains(ctx.Req.Header.Get("Content-Type"), "multipart/form-data") {
		// 	if err := ctx.Req.ParseMultipartForm(config.AttachmentMaxSize << 20); err != nil && !strings.Contains(err.Error(), "EOF") { // 32MB max size
		// 		ctx.Handle(500, "ParseMultipartForm", err)
		// 		return
		// 	}
		// }

		c.Map(ctx)
	}
}
