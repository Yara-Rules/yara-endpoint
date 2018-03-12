package api

import (
	"fmt"
	"time"

	"github.com/Yara-Rules/yara-endpoint/common/command"
	"github.com/Yara-Rules/yara-endpoint/server/context"
	"github.com/Yara-Rules/yara-endpoint/server/models"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func Tasks(ctx *context.Context) {
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

func NewTask(ctx *context.Context, newTask NewTaskForm) {
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

func DeleteTask(ctx *context.Context) {
	id := ctx.Params(":id")
	task_id := ctx.Params(":task")

	var schedule models.Schedule
	selector := bson.M{"$and": []bson.M{bson.M{"ulid": id}, bson.M{"tasks.task_id": task_id}, bson.M{"tasks.status": models.Initial}}}
	err := ctx.DB.C(models.Schedules).Find(selector).One(&schedule)
	if err == mgo.ErrNotFound {
		ctx.JSON(400, Error{
			Error:    true,
			ErrorID:  1,
			ErrorMsg: "Task not found",
		})
		return
	}

	if len(schedule.Tasks) == 1 {
		// Remove the hole schedule due to no task remmaing
		err = ctx.DB.C(models.Schedules).RemoveId(schedule.ID)
		if err != nil {
			ctx.JSON(400, Error{
				Error:    true,
				ErrorID:  1,
				ErrorMsg: fmt.Sprintf("%s", err),
			})
			return
		}
	} else {
		update := bson.M{"$pull": bson.M{"tasks": bson.M{"task_id": task_id}}}
		err := ctx.DB.C(models.Schedules).Update(selector, update)
		if err != nil {
			ctx.JSON(400, Error{
				Error:    true,
				ErrorID:  1,
				ErrorMsg: fmt.Sprintf("%s", err),
			})
			return
		}
	}
	ctx.JSON(200, Error{
		Error:    false,
		ErrorID:  0,
		ErrorMsg: "",
	})
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
