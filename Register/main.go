package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	//使用sql包时必须注入（至少）一个数据库驱动。
	//进行注入mysql数据库驱动，执行它的init函数
	_ "github.com/go-sql-driver/mysql"
)

//DB 数据库对象
var DB *sql.DB

func init() {
	//连接数据库
	db, err := sql.Open("mysql", "root:password@tcp(localhost:3306)/databasename")
	if err != nil {
		log.Fatal(err)
		return
	}
	DB = db
}

func check(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		username := r.FormValue("username")
		var Uusername string
		err := DB.QueryRow("select Uusername from user_table where Uusername = ?", username).Scan(&Uusername)
		if err == sql.ErrNoRows {
			w.Write([]byte("用户名错误，不存在该用户"))
			return
		} else if err != nil {
			w.Write([]byte("服务器错误500"))
			return
		}
		w.Write([]byte(""))
	} else {
		w.Write([]byte("客户端错误，"))
	}
}
func login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, err := template.ParseFiles("login.html")
		if err != nil {
			log.Fatal(err)
			w.Write([]byte("服务器错误,500"))
			return
		}
		t.Execute(w, nil)
	} else if r.Method == "POST" {
		r.ParseForm()
		username := r.PostForm["username"][0]
		password := r.PostForm["password"][0]
		var Uusername string
		err := DB.QueryRow("select Uusername from user_table where Uusername =?  and Upassword =?", username, password).Scan(&Uusername)
		if err == sql.ErrNoRows {

			t, err := template.ParseFiles("login.html")
			if err != nil {
				log.Fatal(err)
				w.Write([]byte("服务器错误,500"))
				return
			}
			t.Execute(w, "登入失败，请重新登入")
		} else if err != nil {
			w.Write([]byte("服务器错误500"))
			log.Fatal(err)
			return
		} else {
			success(w, r)
		}
	}
}

func success(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.PostForm["username"][0]
	t, err := template.ParseFiles("success.html")
	if err != nil {
		log.Fatal(err)
		return
	}
	t.Execute(w, username)
}

func checkRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		username := r.FormValue("username")
		var Uusername string
		err := DB.QueryRow("select Uusername from user_table where Uusername = ?", username).Scan(&Uusername)
		if err != sql.ErrNoRows {
			w.Write([]byte("该账号已注册"))
		} else {
			w.Write([]byte(""))
		}
	}
}

func register(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, err := template.ParseFiles("register.html")
		if err != nil {
			log.Fatal(err)
			http.NotFound(w, r)
			return
		}
		t.Execute(w, nil)
	} else if r.Method == "POST" {
		r.ParseForm()
		username := r.FormValue("username")
		password := r.FormValue("password")
		var Uusername string
		err := DB.QueryRow("select Uusername from user_table where Uusername = ?", username).Scan(&Uusername)
		if err == sql.ErrNoRows {
			DB.Exec("insert into user_table values(?,?)", username, password)
			http.Redirect(w, r, "/login?msg=注册成功", 302)
		} else {
			http.Redirect(w, r, "/register?msg=注册失败", 302)
		}
	}
}

func main() {
	http.HandleFunc("/login", login)
	http.HandleFunc("/check", check)
	http.HandleFunc("/register", register)
	http.HandleFunc("/checkRegister", checkRegister)
	http.ListenAndServe("", nil)

}
