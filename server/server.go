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
	
    "os"
    "strconv"
)


var database *sql.DB




func main() {
	// read config
	file, err := os.Open("dbconfig.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	data := make([]byte, 100)
	config, err := file.Read(data)
	if err != nil {
		log.Fatal(err)
	}
    //fmt.Println(string(data[:config]))

	//db, err := sql.Open("postgres", "dbname=worktime user=postgres password=poilo777 host=localhost sslmode=disable")
	db, err := sql.Open("postgres", string(data[:config]))
	
	if err != nil {
		log.Fatal(err)
	}
	database = db
	defer db.Close()



	


    router := mux.NewRouter()
	router.HandleFunc("/", MainPageHandler)



	router.HandleFunc("/{iduser}", DatePageHandler)
	router.HandleFunc("/{iduser}/{date}", WorktimePage)
	router.HandleFunc("/{date}/{time}/{iduser}/{inout}", AddCheck)
    http.Handle("/",router)
    fmt.Println("Server is listening...")
	http.ListenAndServe(":8080", nil)


}



func MainPageHandler(w http.ResponseWriter, r *http.Request) { // находимся на главной странице  

	type User struct {
		Id int
		Number int
		Name1  string
		Name2  string
		Name3  string
	}

	rows, err := database.Query("select  id, number, name1, name2, name3 from public.users;")
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	user_current := []User{}
	for rows.Next() {
		p := User{}
		err := rows.Scan(&p.Id, &p.Number, &p.Name1, &p.Name2, &p.Name3 )
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



func DatePageHandler(w http.ResponseWriter, r *http.Request) { // находимся на странице  выбора даты
	tmpl, _ := template.ParseFiles("templates/page_data.html")
	//w.Header().Set("Content-Type", "text/html")
	tmpl.Execute(w, nil)
}



func WorktimePage(w http.ResponseWriter, r *http.Request) { // страница записей 2 таблицы вход/выход
	vars := mux.Vars(r)
	iduser := str_handle(vars["iduser"])
	date := str_handle(vars["date"])


	type TimeIN struct {
		Timein  string
	}

	type TimeOUT struct {
		Timeout  string
	}
	
	var Data struct {
        TimeIn []TimeIN
        TimeOut []TimeOUT
    }

    rows1, err := database.Query("select  time from public.checktime WHERE iduser = "+iduser+" AND date = '"+date+"' AND inout = 1;")
    if err != nil {
        log.Println(err)
        return
    }
    defer rows1.Close()

    for rows1.Next() {
        u := TimeIN{}
        if err := rows1.Scan(&u.Timein); err != nil {
            log.Println(err)
            continue
        }
        Data.TimeIn = append(Data.TimeIn, u)
    }

    rows2, err := database.Query("select  time from public.checktime WHERE iduser = "+iduser+" AND date = '"+date+"' AND inout = 0;")
    if err != nil {
        log.Println(err)
        return
    }
    defer rows2.Close()

    for rows2.Next() {
        a := TimeOUT{}
        if err := rows2.Scan(&a.Timeout); err != nil {
            log.Println(err)
            continue
        }
        Data.TimeOut = append(Data.TimeOut, a)
    }

	
    tmpl, err := template.ParseFiles("templates/work.html")
    if err != nil {
        log.Println(err)
        return
    }
    if err := tmpl.Execute(w, Data); err != nil {
        log.Println(err)
    }


}


// добавляем запись в таблицу учета времени
func AddCheck(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
	date := str_handle(vars["date"])
	time := str_handle(vars["time"])
	iduser := str_handle(vars["iduser"])
	inout := str_handle(vars["inout"])


	// проверка когда была сделана последняя запись (чтобы не дублировать близкие по времени)
	type Time struct {
		Value  string
	}
	row, err := database.Query("SELECT \"time\" FROM public.checktime WHERE iduser = "+iduser+" AND date = '"+date+"' ORDER BY  time DESC LIMIT 1;")  // запрос последнего времени
	if err != nil {
		log.Println(err)
	}
	defer row.Close()
	time_current := []Time{}
	for row.Next() {
		p := Time{}
		err := row.Scan(&p.Value)
		if err != nil {
			fmt.Println(err)
			continue
		}
		time_current = append(time_current, p)
	}

	//fmt.Println(time_current[0].Value)

	// сравнение последнего времени в базе с текущим
	if (len(time_current) != 0 ) {      
		time_in_base := time_current[0].Value
		
		if (time_in_base != "") {
			time_current := strings.Split(time, ":")
			time_compare := strings.Split(time_in_base, ":")
	
			a1, _ := strconv.Atoi(time_current[0]) 
			a2, _ := strconv.Atoi(time_current[1]) 
			a3, _ := strconv.Atoi(time_current[2]) 
			b1, _ := strconv.Atoi(time_compare[0])
			b2, _ := strconv.Atoi(time_compare[1])
			b3, _ := strconv.Atoi(time_compare[2])
	
	
			time_different := (a1*3600 + a2*60 + a3) - (b1*3600 + b2*60 + b3)  // подсчет сколько времени прошло с последней записи
	
			fmt.Println(time_different)
	
			if (time_different > 10) {
				// вносим запись если последняя была сделана не недавно
				rows, err := database.Query("INSERT INTO public.checktime(date, time, iduser, inout) VALUES ('" + date + "', '" + time + "', " + iduser + "," + inout + ");")
				fmt.Println("Распознан user с id=" + iduser + "  время: " + time + " дата: " + date)
				if err != nil {
		            log.Println(err)
	            }   
	            defer rows.Close()
			}
		}
	} else {
		rows, err := database.Query("INSERT INTO public.checktime(date, time, iduser, inout) VALUES ('" + date + "', '" + time + "', " + iduser + "," + inout + ");")
				fmt.Println("Распознан user с id=" + iduser + "  время: " + time + " дата: " + date)
				if err != nil {
		            log.Println(err)
	            }   
	            defer rows.Close()
	}
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


