package router

import (
	"github.com/Yara-Rules/yara-endpoint/server/router/api"
	macaron "gopkg.in/macaron.v1"
)

func RegisterRoutes(m *macaron.Macaron) {
	m.Get("/", api.Dashboard).Name("index")
	m.Get("/dashboard", api.Dashboard).Name("dashboard")
	m.Get("/assets", api.Assets).Name("assets")
	m.Group("/rules", func() {
		m.Get("/", api.ShowRules).Name("showRules")
		m.Post("/add", api.AddRules).Name("addRules")
		m.Delete("/delete/:id", api.DeleteRules).Name("deleteRules")
	})
	m.Group("/tasks", func() {
		m.Get("/", api.ShowTasks).Name("showTasks")
		m.Post("/add", api.TasksAdd).Name("addTasks")
		m.Delete("/delete/:id", api.TasksDelete).Name("deleteTask")
		m.Get("/results", api.TasksResults).Name("tasksResults")
		m.Get("/result/:id", api.TasksResult).Name("taskResult")
	})
}
