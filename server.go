package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

type Question struct {
	ID      int      `json:"id"`
	Text    string   `json:"text"`
	Options []string `json:"options"`
	Answer  int      `json:"answer"`
}

func main() {
	// ONLY ONE ROOT HANDLER
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		http.ServeFile(w, r, "index.html")
	})

	http.HandleFunc("/get-questions", func(w http.ResponseWriter, r *http.Request) {
		category := strings.ToLower(r.URL.Query().Get("type"))
		if category == "" { category = "jamb" }

		fileName := fmt.Sprintf("%s.json", category)
		fileData, err := os.ReadFile(fileName)
		if err != nil {
			http.Error(w, "File not found", 404)
			return
		}

		var subjectsData map[string][]Question
		if err := json.Unmarshal(fileData, &subjectsData); err != nil {
			log.Printf("JSON Error in %s: %v", fileName, err)
			http.Error(w, "JSON Error", 500)
			return
		}

		var pool []Question
		for _, qList := range subjectsData {
			pool = append(pool, qList...)
		}

		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(pool), func(i, j int) {
			pool[i], pool[j] = pool[j], pool[i]
		})

		if len(pool) > 100 { pool = pool[:100] }

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(pool)
	})

	http.HandleFunc("/sw.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "sw.js")
	})


	fmt.Println("Exam Hub Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
