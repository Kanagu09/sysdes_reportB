package service

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	database "todolist.go/db"
)

// Home renders index.html
func Home(ctx *gin.Context) {
	userid := CheckCookieId(ctx)
	fmt.Println("UserID:", userid)

	name := ""

	if userid != 0 {
		// Get DB connection
		db, err := database.GetConnection()
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		// get user id
		var user database.User
		err = db.Get(&user, "SELECT * FROM users WHERE id=?", userid)
		if err != nil {
			ctx.String(http.StatusBadRequest, err.Error())
			return
		}
		name = user.Name
	}

	ctx.HTML(http.StatusOK, "index.html", gin.H{"Title": "HOME", "UserID": userid, "Name": name})
}
