package api

import (
	"github.com/Yara-Rules/yara-endpoint/server/context"
	"github.com/Yara-Rules/yara-endpoint/server/models"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func TasksResults(ctx *context.Context) {
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

func TasksResult(ctx *context.Context) {
	/* TODO: API for retrieving a task result
	 */
}

func TasksReport(ctx *context.Context) {
	/* TODO: API for retrieving a report result
	 */
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
