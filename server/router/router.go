package router

import (
	"github.com/Yara-Rules/yara-endpoint/server/router/api"
	"github.com/go-macaron/binding"
	macaron "gopkg.in/macaron.v1"
)

func RegisterRoutes(m *macaron.Macaron) {
	m.Get("/", api.Index).Name("index")
	m.Get("/dashboard", api.Dashboard).Name("dashboard")
	m.Group("/assets", func() {
		m.Get("/", api.Assets).Name("assets")
		m.Post("/", binding.Bind(api.NewAssetForm{}), api.NewAsset).Name("newAsset")
		m.Put("/:id", binding.Bind(api.EditAssetForm{}), api.EditAsset).Name("editAsset")
		m.Delete("/:id", api.DeleteAsset).Name("deleteAsset")
	})
	m.Get("/commands", api.Commands).Name("commands")
	m.Group("/rules", func() {
		m.Get("/", api.Rules).Name("Rules")
		m.Post("/", binding.Bind(api.NewRuleForm{}), api.NewRule).Name("addRules")
		m.Put("/:id", binding.Bind(api.EditRuleForm{}), api.EditRule).Name("editRules")
		m.Delete("/:id", api.DeleteRule).Name("deleteRule")
	})
	m.Group("/tasks", func() {
		m.Get("/", api.Tasks).Name("Tasks")
		m.Post("/", binding.Bind(api.NewTaskForm{}), api.NewTask).Name("addTasks")
		// m.Put("/:id", binding.Bind(api.EditTaskForm{}), api.EditTask).Name("editTasks")
		m.Delete("/:id/:task", api.DeleteTask).Name("deleteTask")
		m.Group("/results", func() {
			m.Get("/", api.TasksResults).Name("tasksResults")
			// m.Get("/:id", api.TasksResult).Name("taskResult")
			// m.Get("/:id/report/:id", api.TasksReport).Name("taskReport")
		})
	})
	// m.Group("/errors", func() {
	// 	m.Get("/", api.ShowErrors).Name("showErrors")
	// 	m.Delete("/:id", api.ErrorDelete).Name("deleteError")
	// })
}
