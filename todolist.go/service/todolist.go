package service

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	database "todolist.go/db"
)

type AddData struct {
	Title string `form:"title"`
}

type EditData struct {
	Title  string `form:"title"`
	IsDone string `form:"is_done"`
}

type Option struct {
	Filter string `form:"filter"`
	Search string `form:"search"`
}

// TaskList renders list of tasks in DB
func TaskList(ctx *gin.Context) {
	userid := CheckCookieId(ctx)
	if userid == 0 {
		ctx.Redirect(http.StatusInternalServerError, "/")
	}

	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	// Get tasks in DB
	var tasks []database.Task
	err = db.Select(&tasks, "SELECT * FROM tasks WHERE user_id =?", userid) // Use DB#Select for multiple entries
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	// Render tasks
	ctx.HTML(http.StatusOK, "task_list.html", gin.H{"Tasks": tasks})
}

// TaskList renders list of tasks in DB
func FilteredTaskList(ctx *gin.Context) {
	userid := CheckCookieId(ctx)
	if userid == 0 {
		ctx.Redirect(http.StatusInternalServerError, "/")
	}

	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	// get data from post
	var option Option
	ctx.Bind(&option)

	sql := fmt.Sprintf("SELECT * FROM tasks WHERE user_id =%d AND", userid)
	switch option.Filter {
	case "all":
		sql += " (is_done = 0 OR is_done = 1)"
	case "done":
		sql += " is_done = 1"
	case "undone":
		sql += " is_done = 0"
	}

	if option.Search != "" {
		sql += " AND title LIKE '%" + option.Search + "%'"
	}

	// Get tasks in DB
	var tasks []database.Task
	err = db.Select(&tasks, sql) // Use DB#Select for multiple entries
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	// Render tasks
	ctx.HTML(http.StatusOK, "task_list.html", gin.H{"Title": "Task list", "Tasks": tasks, "Option": option})
}

// ShowTask renders a task with given ID
func ShowTask(ctx *gin.Context) {
	userid := CheckCookieId(ctx)
	if userid == 0 {
		ctx.Redirect(http.StatusInternalServerError, "/")
	}

	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	// parse ID given as a parameter
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	// Get a task with given ID
	var task database.Task
	err = db.Get(&task, "SELECT * FROM tasks WHERE user_id=? AND id=?", userid, id) // Use DB#Get for one entry
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	// Render task
	ctx.HTML(http.StatusOK, "show_task.html", gin.H{"ID": task.ID, "Title": task.Title})
}

// AddTask add a task
func AddTask(ctx *gin.Context) {
	userid := CheckCookieId(ctx)
	if userid == 0 {
		ctx.Redirect(http.StatusInternalServerError, "/")
	}

	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	// Catch data from post
	var add_data AddData
	ctx.Bind(&add_data)

	// Add task
	if add_data.Title != "" {
		data := map[string]interface{}{"title": add_data.Title, "user_id": userid}
		_, err = db.NamedExec("INSERT INTO tasks (title, user_id) VALUES (:title, :user_id)", data)
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}
		fmt.Println(data["title"], "added.")
	}

	// Redirect list
	ctx.Redirect(http.StatusSeeOther, "/list")
}

// DoneTask done a task
func DoneTask(ctx *gin.Context) {
	userid := CheckCookieId(ctx)
	if userid == 0 {
		ctx.Redirect(http.StatusInternalServerError, "/")
	}

	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	// parse ID given as a parameter
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	// Add task
	data := map[string]interface{}{"id": id, "user_id": userid}
	_, err = db.NamedExec("UPDATE tasks SET is_done = 1 WHERE user_id = (:user_id) AND id = (:id)", data)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	fmt.Println("No.", id, "doned.")

	// Redirect list
	ctx.Redirect(http.StatusSeeOther, "/list")
}

// UndoneTask undone a task
func UndoneTask(ctx *gin.Context) {
	userid := CheckCookieId(ctx)
	if userid == 0 {
		ctx.Redirect(http.StatusInternalServerError, "/")
	}

	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	// parse ID given as a parameter
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	// Add task
	data := map[string]interface{}{"id": id, "user_id": userid}
	_, err = db.NamedExec("UPDATE tasks SET is_done = 0 WHERE user_id = (:user_id) AND id = (:id)", data)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	fmt.Println("No.", id, "undoned.")

	// Redirect list
	ctx.Redirect(http.StatusSeeOther, "/list")
}

// ShowEdit renders a edit page with given ID
func ShowEdit(ctx *gin.Context) {
	userid := CheckCookieId(ctx)
	if userid == 0 {
		ctx.Redirect(http.StatusInternalServerError, "/")
	}

	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	// parse ID given as a parameter
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	// Get a task with given ID
	var task database.Task
	err = db.Get(&task, "SELECT * FROM tasks WHERE user_id=? AND id=?", userid, id) // Use DB#Get for one entry
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	// Render task
	ctx.HTML(http.StatusOK, "task_edit.html", &task)
}

// EditTask change a task info
func EditTask(ctx *gin.Context) {
	userid := CheckCookieId(ctx)
	if userid == 0 {
		ctx.Redirect(http.StatusInternalServerError, "/")
	}

	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	// parse ID given as a parameter
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	// get data from post
	var edit_data EditData
	ctx.Bind(&edit_data)

	// Add task
	data := map[string]interface{}{"id": id, "title": edit_data.Title, "is_done": !(edit_data.IsDone == ""), "user_id": userid}
	_, err = db.NamedExec("UPDATE tasks SET title = (:title) WHERE user_id = (:user_id) AND id = (:id)", data)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	_, err = db.NamedExec("UPDATE tasks SET is_done = (:is_done) WHERE user_id = (:user_id) AND id = (:id)", data)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	fmt.Println("No.", id, "edited.")

	// Redirect list
	ctx.Redirect(http.StatusSeeOther, "/list")
}

// DeleteTask delete a task
func DeleteTask(ctx *gin.Context) {
	userid := CheckCookieId(ctx)
	if userid == 0 {
		ctx.Redirect(http.StatusInternalServerError, "/")
	}

	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	// parse ID given as a parameter
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	// Delete task
	data := map[string]interface{}{"id": id, "user_id": userid}
	_, err = db.NamedExec("DELETE FROM tasks WHERE user_id = (:user_id) AND id = (:id)", data)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	fmt.Println("No.", id, "deleted.")

	// Redirect list
	ctx.Redirect(http.StatusSeeOther, "/list")
}
