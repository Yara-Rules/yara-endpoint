package api

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/Yara-Rules/yara-endpoint/server/context"
	"github.com/Yara-Rules/yara-endpoint/server/models"
)

type Error struct {
	Error    bool   `json:"error"`
	ErrorID  int    `json:"error_id"`
	ErrorMsg string `json:"error_msg"`
}

func Dashboard(ctx *context.Context) {
	/**/

	eps := new([]models.Endpoint)
	ctx.DB.C(models.Endpoints).Find(nil).All(eps)
	ctx.JSON(200, eps)
}

func Assets(ctx *context.Context) {
	/* TODO: Consultar en la base de datos el listado de assets registrados
	 */
	var m []models.Endpoint
	ctx.DB.C(models.Endpoints).Find(nil).All(&m)
	ctx.JSON(200, m)
}

func ShowRules(ctx *context.Context) {
	/* TODO: Consultar en la base de datos el listado de reglas
	 */
	rules := new([]models.Rule)
	ctx.DB.C(models.Rules).Find(nil).All(rules)
	ctx.JSON(200, rules)
}

func AddRules(ctx *context.Context) {
	/* TODO: Añadir una nueva regla
	 */
}

func DeleteRules(ctx *context.Context) {
	/* TODO: Eliminar una nueva regla
	 */

	id := ctx.Params(":id")
	if bson.IsObjectIdHex(id) {
		err := ctx.DB.C(models.Rules).RemoveId(bson.ObjectIdHex(id))
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

func ShowTasks(ctx *context.Context) {
	/* TODO: Consultar en la base de datos el listado de tareas pendientes de ejecución o en ejección
	 */
	tasks := new([]models.Schedule)
	ctx.DB.C(models.Schedules).Find(nil).All(tasks)
	ctx.JSON(200, tasks)
}

func TasksAdd(ctx *context.Context) {
	/* TODO: Añadir una tarea nueva
	 */
}

func TasksDelete(ctx *context.Context) {
	/* TODO: Eliminar una tarea que aún no ha sido iniciada
	 */

	id := ctx.Params(":id")
	if bson.IsObjectIdHex(id) {
		err := ctx.DB.C(models.Schedules).RemoveId(bson.ObjectIdHex(id))
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

func TasksResults(ctx *context.Context) {
	/* TODO: Consultar en la base de datos el listado de todos los resultados
	 */
	report := new([]models.Report)
	ctx.DB.C(models.Reports).Find(nil).All(report)
	ctx.JSON(200, report)
}

func TasksResult(ctx *context.Context) {
	/* TODO: Consular en la base de datos un resultado en concreto
	 */
}

func TasksReport(ctx *context.Context) {
	/* TODO: Consular en la base de datos un resultado en concreto
	 */
}

func ShowErrors(ctx *context.Context) {
	/* TODO: Consular en la base de datos un resultado en concreto
	 */
	errors := new([]models.Error)
	ctx.DB.C(models.Errors).Find(bson.M{"acknowledge": false}).All(errors)
	ctx.JSON(200, errors)
}

func ErrorDelete(ctx *context.Context) {
	/* TODO: Consular en la base de datos un resultado en concreto
	 */
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
