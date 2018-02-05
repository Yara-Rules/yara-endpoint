package api

import (
	"github.com/Yara-Rules/yara-endpoint/server/context"
	"github.com/Yara-Rules/yara-endpoint/server/models"
)

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
}

func AddRules(ctx *context.Context) {
	/* TODO: Añadir una nueva regla
	 */
}

func DeleteRules(ctx *context.Context) {
	/* TODO: Eliminar una nueva regla
	 */
}

func ShowTasks(ctx *context.Context) {
	/* TODO: Consultar en la base de datos el listado de tareas pendientes de ejecución o en ejección
	 */
}

func TasksAdd(ctx *context.Context) {
	/* TODO: Añadir una tarea nueva
	 */
}

func TasksDelete(ctx *context.Context) {
	/* TODO: Eliminar una tarea que aún no ha sido iniciada
	 */
}

func TasksResults(ctx *context.Context) {
	/* TODO: Consultar en la base de datos el listado de todos los resultados
	 */
}

func TasksResult(ctx *context.Context) {
	/* TODO: Consular en la base de datos un resultado en concreto
	 */
}
