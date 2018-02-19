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

	var data publicDashboard

	eps := new([]models.Endpoint)
	rls := new([]models.Rule)

	ctx.DB.C(models.Endpoints).Find(nil).All(eps)
	ctx.DB.C(models.Rules).Find(nil).All(rls)

	data = publicDashboard{
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
	var endpoints []models.Endpoint
	ctx.DB.C(models.Endpoints).Find(nil).All(&endpoints)

	var rules []models.Rule
	ctx.DB.C(models.Rules).Find(nil).All(&rules)

	var reports []models.Report
	ctx.DB.C(models.Reports).Find(nil).All(&reports)

	var pReports []publicReports
	pReports = make([]publicReports, 0)

	for _, report := range reports {
		one := publicReports{
			ULID:     report.ULID,
			Hostname: lookupHostname(endpoints, report.ULID),
			Reports:  make([]reports_, 0),
			CreateAt: report.CreateAt,
			UpdateAt: report.UpdateAt,
		}

		for _, rep := range report.Reports {
			task := task_{
				TaskID:   rep.Task.TaskID,
				Command:  rep.Task.Command,
				Rules:    getRuleNames(ctx.DB, rep.Task.Rules),
				Params:   rep.Task.Params,
				When:     rep.Task.When,
				Status:   rep.Task.Status,
				CreateAt: rep.Task.CreateAt,
				UpdateAt: rep.Task.UpdateAt,
			}

			two := reports_{
				ReportID: rep.ReportID,
				Task:     task,
				Result:   make([]results_, 0),
				CreateAt: rep.CreateAt,
				UpdateAt: rep.UpdateAt,
			}

			for _, res := range rep.Result {
				three := results_{
					File:      res.File,
					RuleName:  res.RuleName,
					Namespace: res.Namespace,
					Tags:      res.Tags,
					Meta:      res.Meta,
					Strings:   make([]strings_, 0),
				}

				for _, str := range res.Strings {
					four := strings_{
						Name:   str.Name,
						Offset: str.Offset,
					}
					three.Strings = append(three.Strings, four)
				}
				two.Result = append(two.Result, three)
			}
			one.Reports = append(one.Reports, two)
		}
		pReports = append(pReports, one)
	}

	ctx.JSON(200, pReports)
}

func lookupHostname(endpoints []models.Endpoint, ulid string) string {
	for _, endpoint := range endpoints {
		if endpoint.ULID == ulid {
			return endpoint.Hostname
		}
	}
	return ""
}

func getRuleNames(db *mgo.Database, listRules []bson.ObjectId) []string {
	var rules []models.Rule
	var list_ []string
	db.C(models.Rules).Find(bson.M{"_id": bson.M{"$in": listRules}}).Select(bson.M{"name": 1}).All(&rules)
	list_ = make([]string, 0)
	for _, rule := range rules {
		list_ = append(list_, rule.Name)
	}
	return list_
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
