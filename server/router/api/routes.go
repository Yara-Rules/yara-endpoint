package api

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/Yara-Rules/yara-endpoint/server/context"
	"github.com/Yara-Rules/yara-endpoint/server/models"
)

func Index(ctx *context.Context) {
	/**/
	ctx.HTML(200, "index")
}

func Dashboard(ctx *context.Context) {
	/**/

	var data dashboard

	eps := new([]models.Endpoint)
	rls := new([]models.Rule)

	ctx.DB.C(models.Endpoints).Find(nil).All(eps)
	ctx.DB.C(models.Rules).Find(nil).All(rls)

	data = dashboard{
		Asset: (*eps),
		Rules: (*rls),
	}

	ctx.JSON(200, data)
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
			ctx.JSON(400, error_{
				Error:    true,
				ErrorID:  1,
				ErrorMsg: "Rule not found",
			})
			return
		} else {
			ctx.JSON(200, error_{
				Error: false,
			})
			return
		}
	}
	ctx.JSON(400, error_{
		Error:    true,
		ErrorID:  1,
		ErrorMsg: "Not valid ID",
	})
}

func ShowTasks(ctx *context.Context) {
	/* TODO: Consultar en la base de datos el listado de tareas pendientes de ejecución o en ejección
	 */
	var pTask []publicTasks
	var rulename map[string]string
	var hostname map[string]string

	rulename = map[string]string{}
	hostname = map[string]string{}

	var schedules []models.Schedule
	ctx.DB.C(models.Schedules).Find(nil).All(&schedules)

	pTask = []publicTasks{}

	for _, schedule := range schedules {
		t := publicTasks{}
		t.ULID = schedule.ULID

		if val, ok := hostname[schedule.ULID]; ok {
			t.Hostname = val
		} else {
			tmp := new(models.Endpoint)
			ctx.DB.C(models.Endpoints).Find(bson.M{"ulid": schedule.ULID}).One(tmp)
			hostname[schedule.ULID] = tmp.Hostname
			t.Hostname = tmp.Hostname
		}

		for _, task := range schedule.Tasks {
			tmp_rules := []string{}
			for _, ruleID := range task.Rules {
				if val, ok := rulename[ruleID.Hex()]; ok {
					tmp_rules = append(tmp_rules, val)
				} else {
					tmp := new(models.Rule)
					ctx.DB.C(models.Rules).Find(bson.M{"_id": ruleID}).One(tmp)
					rulename[tmp.ID.Hex()] = tmp.Name
					tmp_rules = append(tmp_rules, tmp.Name)
				}
			}
			t.Task.Command = task.Command
			t.Task.Rules = tmp_rules
			t.Task.Status = task.Status
			t.Task.When = task.When
			t.Task.UpdateAt = task.UpdateAt

			pTask = append(pTask, t)
		}
	}

	ctx.JSON(200, pTask)
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
			ctx.JSON(400, error_{
				Error:    true,
				ErrorID:  1,
				ErrorMsg: "Rule not found",
			})
			return
		} else {
			ctx.JSON(200, error_{
				Error: false,
			})
			return
		}
	}
	ctx.JSON(400, error_{
		Error:    true,
		ErrorID:  1,
		ErrorMsg: "Not valid ID",
	})
}

func TasksResults(ctx *context.Context) {
	/* TODO: Consultar en la base de datos el listado de todos los resultados
	 */
	reports := new([]models.Report)
	ctx.DB.C(models.Reports).Find(nil).All(reports)
	ctx.JSON(200, reports)
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
			ctx.JSON(400, error_{
				Error:    true,
				ErrorID:  1,
				ErrorMsg: "Rule not found",
			})
			return
		} else {
			ctx.JSON(200, error_{
				Error: false,
			})
			return
		}
	}
	ctx.JSON(400, error_{
		Error:    true,
		ErrorID:  1,
		ErrorMsg: "Not valid ID",
	})
}
