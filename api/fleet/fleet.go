package fleet

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"../../pkg/database"
	"github.com/gorilla/mux"
)

//GetAllShipsHandler Lets retrieve all the ships from the fleet, simple select all rows query
//Due to there's no easy way of mapping the mysql data as in Laravel, I re used a code from stackoverflow: https://stackoverflow.com/a/29164115
func GetAllShipsHandler(w http.ResponseWriter, r *http.Request) {
	var buffer bytes.Buffer // buffer for json response
	rows, err := database.DBCon.Query("SELECT * FROM ships")
	if err != nil {
		fmt.Println(err)
		buffer.WriteString(`{success: "false"}`)
		json.NewEncoder(w).Encode(buffer.String())
	}
	defer rows.Close() // Clean up
	columns, err := rows.Columns()
	if err != nil {
		fmt.Println(err)
		buffer.WriteString(`{success: "false"}`)
		json.NewEncoder(w).Encode(buffer.String())
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
		buffer.WriteString(`{success: "false"}`)
		json.NewEncoder(w).Encode(buffer.String())
	}
	if string(jsonData) == "[]" { // if rows empty then return false
		buffer.WriteString(`{success: "false"}`)
		json.NewEncoder(w).Encode(buffer.String())
	} else {
		w.Write([]byte(string(jsonData)))
	}
}

//GetSingleShipHandler func will return single ship based on /api/fleet/id
func GetSingleShipHandler(w http.ResponseWriter, r *http.Request) {
	var buffer bytes.Buffer // buffer for json response
	params := mux.Vars(r)   // Get the /api/fleet/id parameter

	// Check if parameter is valid integer
	_, err := strconv.Atoi(params["id"])
	if err != nil {
		buffer.WriteString(`{error: "invalid id entered"}`)
		json.NewEncoder(w).Encode(buffer.String())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	rows, err := database.DBCon.Query("SELECT * FROM ships WHERE id=?", params["id"])
	if err != nil {
		fmt.Println(err)
		buffer.WriteString(`{success: "false"}`)
		json.NewEncoder(w).Encode(buffer.String())
	}
	defer rows.Close() // Clean up
	columns, err := rows.Columns()
	if err != nil {
		fmt.Println(err)
		buffer.WriteString(`{success: "false"}`)
		json.NewEncoder(w).Encode(buffer.String())
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
		buffer.WriteString(`{success: "false"}`)
		json.NewEncoder(w).Encode(buffer.String())
	}
	if string(jsonData) == "[]" { // if rows empty then return false
		buffer.WriteString(`{success: "false"}`)
		json.NewEncoder(w).Encode(buffer.String())
	} else {
		w.Write([]byte(string(jsonData)))
	}
}

// CreateShipHandler will read post on API call and form the struct on Ship to then post to the mysql database
func CreateShipHandler(w http.ResponseWriter, r *http.Request) {
	var buffer bytes.Buffer // buffer for json response
	w.Header().Set("Content-Type", "application/json")
	var post Ship
	_ = json.NewDecoder(r.Body).Decode(&post)
	stmt, err := database.DBCon.Prepare("INSERT INTO ships(name, image, class, crew, status, value) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		fmt.Println(err)
		buffer.WriteString(`{success: "false"}`)
		json.NewEncoder(w).Encode(buffer.String())
		w.WriteHeader(500)
	}
	defer stmt.Close() // Always close resource after no longer usage else memory leak..

	_, err = stmt.Exec(post.Name, post.Image, post.Class, post.Crew, post.Status, post.Value)
	if err != nil {
		fmt.Println(err)
		buffer.WriteString(`{success: "false"}`)
		json.NewEncoder(w).Encode(buffer.String())
		w.WriteHeader(500)
	}
	buffer.WriteString(`{success: "true"}`)
	json.NewEncoder(w).Encode(buffer.String())
	w.WriteHeader(201)
}

//DeleteShipHandler deletes ship based on id
func DeleteShipHandler(w http.ResponseWriter, r *http.Request) {
	var buffer bytes.Buffer // buffer for json response
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	// Check if parameter is valid integer
	_, err := strconv.Atoi(params["id"])
	if err != nil {
		buffer.WriteString(`{error: "invalid id entered"}`)
		json.NewEncoder(w).Encode(buffer.String())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	stmt, err := database.DBCon.Prepare("DELETE FROM ships WHERE id = ?")
	if err != nil {
		fmt.Println(err)
		buffer.WriteString(`{success: "false"}`)
		json.NewEncoder(w).Encode(buffer.String())
		w.WriteHeader(500)
	}
	defer stmt.Close() // Always close resource after no longer usage else memory leak..
	_, err = stmt.Exec(params["id"])
	if err != nil {
		fmt.Println(err)
		buffer.WriteString(`{success: "false"}`)
		json.NewEncoder(w).Encode(buffer.String())
		w.WriteHeader(500)
	}

	buffer.WriteString(`{success: "true"}`)
	json.NewEncoder(w).Encode(buffer.String())
	w.WriteHeader(201)
}

//GetAllShipsByClassHandler Filter by class ASC
func GetAllShipsByClassHandler(w http.ResponseWriter, r *http.Request) {
	var buffer bytes.Buffer // buffer for json response
	rows, err := database.DBCon.Query("SELECT * FROM ships ORDER BY class ASC")
	if err != nil {
		fmt.Println(err)
		buffer.WriteString(`{success: "false"}`)
		json.NewEncoder(w).Encode(buffer.String())
	}
	defer rows.Close() // Clean up
	columns, err := rows.Columns()
	if err != nil {
		fmt.Println(err)
		buffer.WriteString(`{success: "false"}`)
		json.NewEncoder(w).Encode(buffer.String())
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
		buffer.WriteString(`{success: "false"}`)
		json.NewEncoder(w).Encode(buffer.String())
	}
	if string(jsonData) == "[]" { // if rows empty then return false
		buffer.WriteString(`{success: "false"}`)
		json.NewEncoder(w).Encode(buffer.String())
	} else {
		w.Write([]byte(string(jsonData)))
	}

}

//GetAllShipsByNameHandler Filter by name ASC
func GetAllShipsByNameHandler(w http.ResponseWriter, r *http.Request) {
	var buffer bytes.Buffer // buffer for json response
	rows, err := database.DBCon.Query("SELECT * FROM ships ORDER BY name ASC")
	if err != nil {
		fmt.Println(err)
		buffer.WriteString(`{success: "false"}`)
		json.NewEncoder(w).Encode(buffer.String())
	}
	defer rows.Close() // Clean up
	columns, err := rows.Columns()
	if err != nil {
		fmt.Println(err)
		buffer.WriteString(`{success: "false"}`)
		json.NewEncoder(w).Encode(buffer.String())
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
		buffer.WriteString(`{success: "false"}`)
		json.NewEncoder(w).Encode(buffer.String())
	}
	if string(jsonData) == "[]" { // if rows empty then return false
		buffer.WriteString(`{success: "false"}`)
		json.NewEncoder(w).Encode(buffer.String())
	} else {
		w.Write([]byte(string(jsonData)))
	}
}

//GetAllShipsByStatusHandler Filter by status ASC
func GetAllShipsByStatusHandler(w http.ResponseWriter, r *http.Request) {
	var buffer bytes.Buffer // buffer for json response
	rows, err := database.DBCon.Query("SELECT * FROM ships ORDER BY status ASC")
	if err != nil {
		fmt.Println(err)
		buffer.WriteString(`{success: "false"}`)
		json.NewEncoder(w).Encode(buffer.String())
	}
	defer rows.Close() // Clean up
	columns, err := rows.Columns()
	if err != nil {
		fmt.Println(err)
		buffer.WriteString(`{success: "false"}`)
		json.NewEncoder(w).Encode(buffer.String())
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
		buffer.WriteString(`{success: "false"}`)
		json.NewEncoder(w).Encode(buffer.String())
	}
	if string(jsonData) == "[]" { // if rows empty then return false
		buffer.WriteString(`{success: "false"}`)
		json.NewEncoder(w).Encode(buffer.String())
	} else {
		w.Write([]byte(string(jsonData)))
	}
}

// UpdateShipHandler will read UPDATE on API call and update the entry in database
func UpdateShipHandler(w http.ResponseWriter, r *http.Request) {
	var buffer bytes.Buffer // buffer for json response
	params := mux.Vars(r)

	// Check if parameter is valid integer
	_, err := strconv.Atoi(params["id"])
	if err != nil {
		buffer.WriteString(`{error: "invalid id entered"}`)
		json.NewEncoder(w).Encode(buffer.String())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	var post Ship
	_ = json.NewDecoder(r.Body).Decode(&post)
	stmt, err := database.DBCon.Prepare("UPDATE ships SET name=?, image=?, class=?, crew=?, status=?, value=? WHERE id=?")
	if err != nil {
		fmt.Println(err)
		buffer.WriteString(`{success: "false"}`)
		json.NewEncoder(w).Encode(buffer.String())
		w.WriteHeader(http.StatusBadGateway)
	}
	defer stmt.Close() // Always close resource after no longer usage else memory leak..

	_, err = stmt.Exec(post.Name, post.Image, post.Class, post.Crew, post.Status, post.Value, params["id"])
	if err != nil {
		fmt.Println(err)
		buffer.WriteString(`{success: "false"}`)
		json.NewEncoder(w).Encode(buffer.String())
		w.WriteHeader(500)
	}
	buffer.WriteString(`{success: "true"}`)
	json.NewEncoder(w).Encode(buffer.String())
	w.WriteHeader(201)
}
