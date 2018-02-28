package api

import (
	"fmt"
	"math/rand"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/Yara-Rules/yara-endpoint/common/command"
	"github.com/Yara-Rules/yara-endpoint/server/context"
	"github.com/Yara-Rules/yara-endpoint/server/models"
	"github.com/oklog/ulid"
	log "github.com/sirupsen/logrus"
)

func Index(ctx *context.Context) {
	/**/
	ctx.HTML(200, "index")
}

func Dashboard(ctx *context.Context) {
	/**/

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

func Assets(ctx *context.Context) {
	/* TODO: Consultar en la base de datos el listado de assets registrados
	 */
	var m []models.Endpoint
	ctx.DB.C(models.Endpoints).Find(nil).All(&m)
	ctx.JSON(200, m)
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

func Commands(ctx *context.Context) {
	/* TODO: Consultar en la base de datos el listado de commandos registrados
	 */
	ctx.JSON(200, struct {
		Commands []string `json:"commands"`
	}{Commands: []string{"ScanFile", "ScanDir", "ScanPID"}})
}

func ShowRules(ctx *context.Context) {
	/* TODO: Consultar en la base de datos el listado de reglas
	 */
	rules := new([]models.Rule)
	ctx.DB.C(models.Rules).Find(nil).All(rules)
	ctx.JSON(200, rules)
}

func AddRules(ctx *context.Context, newRule NewRuleForm) {
	/* TODO: Añadir una nueva regla
	 */

	rule := &models.Rule{
		Name:     newRule.Name,
		RuleID:   newULID().String(),
		Tags:     newRule.Tags,
		Data:     newRule.Data,
		CreateAt: time.Now(),
		UpdateAt: time.Now(),
	}
	err := ctx.DB.C(models.Rules).Insert(rule)
	if err != nil {
		ctx.JSON(400, Error{
			Error:    true,
			ErrorMsg: fmt.Sprintf("%v", err),
		})
	} else {
		ctx.JSON(200, Error{
			Error: false,
		})
	}
}

func DeleteRule(ctx *context.Context) {
	/* TODO: Eliminar una nueva regla
	 */

	ulid := ctx.Params(":id")
	rule := models.Rule{}
	err := ctx.DB.C(models.Rules).Find(bson.M{"rule_id": ulid}).One(&rule)
	if err == mgo.ErrNotFound {
		ctx.JSON(400, Error{
			Error:    true,
			ErrorID:  1,
			ErrorMsg: fmt.Sprintf("%v", err),
		})
		return
	} else {
		// db.schedules.update({"tasks.status": 0}, {"$pull": {"tasks.$.rules": {"$in": [ObjectId("5a920b4ede38e3bb2c4ccd63")]}}})
		selector := bson.M{"tasks.status": 0}
		update := bson.M{"$pull": bson.M{"tasks.$.rules": bson.M{"$in": [...]bson.ObjectId{rule.ID}}}}
		ctx.DB.C(models.Schedules).Update(selector, update)

		// Cleaning up empy tasks
		selector = bson.M{"tasks.rules": bson.M{"$size": 0}}
		ctx.DB.C(models.Schedules).Remove(selector)

		err := ctx.DB.C(models.Rules).Remove(bson.M{"rule_id": ulid})
		if err == mgo.ErrNotFound {
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
			ErrorMsg: "Rule has been removed",
		})
	}
}

func ShowTasks(ctx *context.Context) {
	/* TODO: Consultar en la base de datos el listado de tareas pendientes de ejecución o en ejección
	 */
	var pTask []PublicTasks
	var rulename map[string]string
	var hostname map[string]string

	// Used as buffers
	rulename = map[string]string{}
	hostname = map[string]string{}

	var schedules []models.Schedule
	ctx.DB.C(models.Schedules).Find(nil).All(&schedules)

	var aa []interface{}
	ctx.DB.C(models.Schedules).Find(nil).All(&aa)

	pTask = make([]PublicTasks, 0)

	for _, schedule := range schedules {
		t := PublicTasks{}
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
			t.Task.TaskID = task.TaskID
			t.Task.Command = task.Command
			t.Task.Rules = tmp_rules
			t.Task.Status = task.Status
			t.Task.When = task.When
			t.Task.CreateAt = task.CreateAt
			t.Task.UpdateAt = task.UpdateAt

			pTask = append(pTask, t)
		}
	}
	ctx.JSON(200, pTask)
}

func TasksAdd(ctx *context.Context, newTask NewTaskForm) {
	/* TODO: Añadir una tarea nueva
	 */

	log.Debug("Adding Task")
	var scs []models.Schedule
	scs = make([]models.Schedule, 0)

	for _, ep := range newTask.Assets {
		var tasks []models.Task
		tasks = make([]models.Task, 0)

		dtime, err := parseDateTime(newTask.When)
		if err != nil {
			ctx.JSON(400, Error{
				Error:    true,
				ErrorID:  1,
				ErrorMsg: "Unable to parse the datetime",
			})
			return
		}

		ts := models.Task{
			TaskID:   newULID().String(),
			Command:  command.RevAlias[newTask.Command],
			Rules:    getRulesId(ctx, newTask.Rules),
			Params:   newTask.Target,
			When:     dtime,
			Status:   models.Initial,
			CreateAt: time.Now(),
			UpdateAt: time.Now(),
		}

		tasks = append(tasks, ts)

		sc := models.Schedule{
			ULID:     ep,
			Tasks:    tasks,
			CreateAt: time.Now(),
			UpdateAt: time.Now(),
		}
		scs = append(scs, sc)
	}
	err := updateOrInsert(ctx, scs)
	if err != nil {
		ctx.JSON(400, Error{
			Error:    true,
			ErrorID:  1,
			ErrorMsg: "Unable to insert the task",
		})
		return
	}
	ctx.JSON(200, Error{
		Error: false,
	})
	return
}

func getRulesId(ctx *context.Context, ruleList []string) []bson.ObjectId {
	var rules []models.Rule
	ctx.DB.C(models.Rules).Find(bson.M{"rule_id": bson.M{"$in": ruleList}}).All(&rules)
	var rids []bson.ObjectId
	rids = make([]bson.ObjectId, 0)
	for _, r := range rules {
		rids = append(rids, r.ID)
	}
	return rids
}

func updateOrInsert(ctx *context.Context, scs []models.Schedule) error {
	coll := ctx.DB.C(models.Schedules)
	bulk := coll.Bulk()
	var aux models.Schedule

	for _, sc := range scs {
		err := coll.Find(bson.M{"ulid": sc.ULID}).One(&aux)
		if err == mgo.ErrNotFound {
			coll.Insert(sc)
		} else {
			for _, task := range sc.Tasks {
				selector := bson.M{"ulid": sc.ULID}
				update := bson.M{"$push": bson.M{"tasks": task}}
				bulk.Upsert(selector, update)
			}
		}
	}
	_, err := bulk.Run()
	if err != nil {
		return err
	}
	return nil
}

func parseDateTime(dt string) (time.Time, error) {
	return time.Parse("2006-01-02T15:04:05.000Z", dt)
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
	var endpoints []models.Endpoint
	ctx.DB.C(models.Endpoints).Find(nil).All(&endpoints)

	var rules []models.Rule
	ctx.DB.C(models.Rules).Find(nil).All(&rules)

	var reports []models.Report
	ctx.DB.C(models.Reports).Find(nil).All(&reports)

	var pReports []PublicReports
	pReports = make([]PublicReports, 0)

	for _, report := range reports {
		one := PublicReports{
			ULID:     report.ULID,
			Hostname: lookupHostname(endpoints, report.ULID),
			Reports:  make([]Report, 0),
			CreateAt: report.CreateAt,
			UpdateAt: report.UpdateAt,
		}

		for _, rep := range report.Reports {
			task := Task{
				TaskID:   rep.Task.TaskID,
				Command:  rep.Task.Command,
				Rules:    getRuleNames(ctx.DB, rep.Task.Rules),
				Params:   rep.Task.Params,
				When:     rep.Task.When,
				Status:   rep.Task.Status,
				CreateAt: rep.Task.CreateAt,
				UpdateAt: rep.Task.UpdateAt,
			}

			two := Report{
				ReportID: rep.ReportID,
				Task:     task,
				Result:   make([]Result, 0),
				CreateAt: rep.CreateAt,
				UpdateAt: rep.UpdateAt,
			}

			for _, res := range rep.Result {
				three := Result{
					File:      res.File,
					RuleName:  res.RuleName,
					Namespace: res.Namespace,
					Tags:      res.Tags,
					Meta:      res.Meta,
					Strings:   make([]YString, 0),
				}

				for _, str := range res.Strings {
					four := YString{
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

func newULID() ulid.ULID {
	t := time.Now()
	entropy := rand.New(rand.NewSource(t.UnixNano()))
	return ulid.MustNew(ulid.Timestamp(t), entropy)
}
