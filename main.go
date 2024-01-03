package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

// Task represents a task in the database
type Task struct {
	taskname     string
	AssigneeName string
	DueDate      string
	Status       string
}

// PageVariables struct to pass data to the HTML template
type PageVariables struct {
	Tasks []Task
}

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("postgres", "user=postgres password=123 host=localhost dbname=task port=5432 sslmode=disable")

	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}
	// running
}

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/tasks", tasksHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	fmt.Println("Server is running on http://localhost:3000")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal("Error starting the server:", err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index.html", nil)
	fmt.Println("I am homehandler function")
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Handle form submission
		r.ParseForm()

		taskname := r.FormValue("taskName")
		assignee := r.FormValue("assignee")
		duedate := r.FormValue("duedate")
		status := r.FormValue("status")
		fmt.Println("I am function taskHandler")

		// Insert the task into the database
		err := insertTask(taskname, assignee, duedate, status)

		if err != nil {
			http.Error(w, "Error inserting task: "+err.Error(), http.StatusInternalServerError)
			log.Println("Error inserting task:", err)
			return
		}
	}

	// Retrieve tasks from the database
	tasks, err := getTasks()
	if err != nil {
		http.Error(w, "Error retrieving tasks", http.StatusInternalServerError)
		log.Println("Error retrieving tasks:", err)
		return
	}

	// Render the task list
	p := PageVariables{Tasks: tasks}
	renderTemplate(w, "index.html", &p)
}

func renderTemplate(w http.ResponseWriter, tmplFile string, p *PageVariables) {
	tmpl, err := template.ParseFiles(tmplFile)
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		log.Println("Error parsing template:", err)
		return
	}
	log.Println("i am renderTemplate") //code running up to here.

	err = tmpl.Execute(w, p)
	if err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		log.Println("Error executing template:", err)
	}
}

func insertTask(taskname, assignee, duedate, status string) error {
	query := "INSERT INTO taskinfo (taskname, assigneename, duedate, status) VALUES ($1, $2, $3, $4)"

	// Log the query and parameters for debugging
	log.Printf("Executing query: %s with parameters: %s, %s, %s, %s", query, taskname, assignee, duedate, status)

	_, err := db.Exec(query, taskname, assignee, duedate, status)
	if err != nil {
		log.Println("Error inserting task:", err)
		return err
	}

	log.Println("Task inserted successfully")
	fmt.Println("I am insertTask")
	return nil
}

func getTasks() ([]Task, error) {
	rows, err := db.Query("SELECT taskname, assigneename, duedate, status FROM taskinfo")
	if err != nil {
		log.Println("Error querying tasks:", err)
		return nil, err
	}
	defer rows.Close()

	var tasks []Task

	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.taskname, &task.AssigneeName, &task.DueDate, &task.Status); err != nil {
			log.Println("Error scanning task:", err)
			continue
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error iterating over rows:", err)
		return nil, err
	}

	return tasks, nil
}
