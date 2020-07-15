package fleet

import (
	"encoding/json"
	"fmt"
	"net/http"

	"../../pkg/database"
)

//GetAllShipsHandler Lets retrieve all the ships from the fleet, simple select all rows query
//Due to there's no easy way of mapping the mysql data as in Laravel, I re used a code from stackoverflow: https://stackoverflow.com/a/29164115
func GetAllShipsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DBCon.Query("SELECT * FROM ships")
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close() // Clean up
	columns, err := rows.Columns()
	if err != nil {
		fmt.Println(err)
	}
	count := len(columns)
	tableData := make([]map[string]interface{}, 0)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	for rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		entry := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			entry[col] = v
		}
		tableData = append(tableData, entry)
	}
	jsonData, err := json.Marshal(tableData)
	if err != nil {
		fmt.Println(err)
	}
	w.Write([]byte(string(jsonData)))

}

// CreateShipHandler will read post on API call and form the struct on Ship to then post to the mysql database
func CreateShipHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var post Ship
	_ = json.NewDecoder(r.Body).Decode(&post)
	stmt, err := database.DBCon.Prepare("INSERT INTO ships(name, image, class, crew, status, value) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		fmt.Fprintf(w, "Error:%+v ", err)
		w.WriteHeader(500)
	}
	defer stmt.Close() // Always close resource after no longer usage else memory leak..

	_, err = stmt.Exec(post.Name, post.Image, post.Class, post.Crew, post.Status, post.Value)
	if err != nil {
		fmt.Fprintf(w, "Error:%+v ", err)
		w.WriteHeader(500)
	}
	fmt.Fprintf(w, "success")
	w.WriteHeader(201)
}
