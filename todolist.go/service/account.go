package service

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	database "todolist.go/db"
)

type UserInfo struct {
	Name string `form:"name"`
	Pass string `form:"password"`
}

// CheckCookieID get id of cookie
func CheckCookieId(ctx *gin.Context) int {
	id_str, err := ctx.Request.Cookie("id")
	if err != nil {
		return 0
	}
	id, err := strconv.Atoi(id_str.Value)
	if err != nil {
		return 0
	}
	return id
}

// SetCookie set cookie
func SetCookie(ctx *gin.Context, id int) {
	_, err := ctx.Request.Cookie("id")
	if err != nil {
	} else {
		ctx.SetCookie("id", "", -1, "/", "localhost", false, true)
		fmt.Print("Disable Cookie.", "\n")
	}
	ctx.SetCookie("id", strconv.Itoa(id), 3600, "/", "localhost", false, true)
	fmt.Print("Set Cookie: ", id, "\n")
}

// LoginPage renders a login page
func LoginPage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "login.html", gin.H{})
}

// RegisterPage renders a register page
func RegisterPage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "register.html", gin.H{})
}

// ChangeNamePage renders a change_name page
func ChangeNamePage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "change_name.html", gin.H{})
}

// ChangePassPage renders a change_pass page
func ChangePassPage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "change_pass.html", gin.H{})
}

// DeleteAccountPage renders a delete_account page
func DeleteAccountPage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "delete_account.html", gin.H{})
}

// get SHA256 []byte from string
func GetSHA256Binary(s string) []byte {
	r := sha256.Sum256([]byte(s))
	return r[:]
}

// Make hash
func Hash(str string) string {
	byte := GetSHA256Binary(str)
	hex := hex.EncodeToString(byte)
	return string(hex)
}

// Register register account info
func Register(ctx *gin.Context) {
	// Catch data from post
	var user_info UserInfo
	ctx.Bind(&user_info)

	// hash
	user_info.Pass = Hash(user_info.Pass)

	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	// Register user info
	data := map[string]interface{}{"name": user_info.Name, "pass": user_info.Pass}
	_, err = db.NamedExec("INSERT INTO users (name, password) VALUES (:name, :pass)", data)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	fmt.Println(data["name"], "registered.")

	// get user id
	var user database.User
	err = db.Get(&user, "SELECT * FROM users WHERE name=?", data["name"])
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	// set cookie
	SetCookie(ctx, int(user.ID))

	// Redirect home
	ctx.Redirect(http.StatusSeeOther, "/")
}

// Login login
func Login(ctx *gin.Context) {
	// Catch data from post
	var user_info UserInfo
	ctx.Bind(&user_info)

	// hash
	user_info.Pass = Hash(user_info.Pass)

	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	// get password
	var user database.User
	err = db.Get(&user, "SELECT * FROM users WHERE name=?", user_info.Name)
	if err != nil {
		ctx.Redirect(http.StatusSeeOther, "/login")
		return
	}

	if user.Password == user_info.Pass {
		// set cookie
		SetCookie(ctx, int(user.ID))

		// Redirect home
		ctx.Redirect(http.StatusSeeOther, "/")
	} else {
		// Redirect home
		ctx.Redirect(http.StatusSeeOther, "/login")
	}
}

// ChangeName change name
func ChangeName(ctx *gin.Context) {
	userid := CheckCookieId(ctx)
	if userid == 0 {
		ctx.Redirect(http.StatusInternalServerError, "/")
	}

	// Catch data from post
	var user_info UserInfo
	ctx.Bind(&user_info)

	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	// Register user info
	data := map[string]interface{}{"name": user_info.Name, "user_id": userid}
	_, err = db.NamedExec("UPDATE users SET name = (:name) WHERE id = (:user_id)", data)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	fmt.Println(data["name"], "changed name.")

	// Redirect home
	ctx.Redirect(http.StatusSeeOther, "/")
}

// ChangePass change password
func ChangePass(ctx *gin.Context) {
	userid := CheckCookieId(ctx)
	if userid == 0 {
		ctx.Redirect(http.StatusInternalServerError, "/")
	}

	// Catch data from post
	var user_info UserInfo
	ctx.Bind(&user_info)

	// hash
	user_info.Pass = Hash(user_info.Pass)

	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	// Register user info
	data := map[string]interface{}{"name": user_info.Name, "pass": user_info.Pass, "user_id": userid}
	_, err = db.NamedExec("UPDATE users SET password = (:pass) WHERE id = (:user_id)", data)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	fmt.Println("password changed.")

	// Redirect home
	ctx.Redirect(http.StatusSeeOther, "/")
}

// DeleteAccount delete account
func DeleteAccount(ctx *gin.Context) {
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

	// Register user info
	data := map[string]interface{}{"user_id": userid}
	_, err = db.NamedExec("DELETE FROM users WHERE id = (:user_id)", data)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	fmt.Println("Account delete.")

	SetCookie(ctx, 0)

	// Redirect home
	ctx.Redirect(http.StatusSeeOther, "/")
}

// Logout logout
func Logout(ctx *gin.Context) {
	SetCookie(ctx, 0)

	// Redirect home
	ctx.Redirect(http.StatusSeeOther, "/")
}
