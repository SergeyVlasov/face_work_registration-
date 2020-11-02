package main


import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	_ "github.com/lib/pq"
	"github.com/gorilla/mux"
    //"time"
	"strings"
	

    
)


var database *sql.DB




func main() {

	db, err := sql.Open("postgres", "dbname=worktime user=postgres password=password host=localhost sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	database = db
	defer db.Close()


    router := mux.NewRouter()
	router.HandleFunc("/", MainPageHandler)
	router.HandleFunc("/{date}/{time}/{iduser}/{inout}", AddCheck)
    http.Handle("/",router)
    fmt.Println("Server is listening...")
	http.ListenAndServe(":8080", nil)


}


type User struct {
	Number int
	Name1  string
	Name2  string
	Name3  string
}

func MainPageHandler(w http.ResponseWriter, r *http.Request) { // находимся на главной странице  www.XXX.xx/

	rows, err := database.Query("select  number, name1, name2, name3 from public.users;")

	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	user_current := []User{}
	for rows.Next() {
		p := User{}
		err := rows.Scan(&p.Number, &p.Name1, &p.Name2, &p.Name3 )
		if err != nil {
			fmt.Println(err)
			continue
		}
		user_current = append(user_current, p)
	}

	tmpl, _ := template.ParseFiles("templates/main_page.html")
	w.Header().Set("Content-Type", "text/html")
	tmpl.Execute(w, user_current)

}





// добавляем запись в таблицу учета времени
func AddCheck(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
	date := str_handle(vars["date"])
	time := str_handle(vars["time"])
	iduser := str_handle(vars["iduser"])
	inout := str_handle(vars["inout"])
    rows, err := database.Query("INSERT INTO public.checktime(date, time, iduser, inout) VALUES ('" + date + "', '" + time + "', " + iduser + "," + inout + ");")
	if err != nil {
		log.Println(err)
	} 
	defer rows.Close()


}


// security block 
func str_handle(inpt string) (outpt string) {
	filter0 := strings.ToTitle(inpt)
	filter1 := strings.Replace(filter0, ";", "", -1)
	filter2 := strings.Replace(filter1, "'", "", -1)
	filter3 := strings.Replace(filter2, "%", "", -1)
	filter4 := strings.Replace(filter3, "&", "", -1)
	filter5 := strings.Replace(filter4, "?", "", -1)
	filter6 := strings.Replace(filter5, "drop", "", -1)
	filter7 := strings.Replace(filter6, "table", "", -1)
	filter8 := strings.Replace(filter7, "delete", "", -1)
	filter9 := strings.Replace(filter8, "alter", "", -1)
	outpt = filter9
	return outpt
}
