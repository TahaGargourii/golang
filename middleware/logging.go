package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		end := time.Now()
		duration := end.Sub(start)
		logMessage := fmt.Sprintf("[%s] %s %s %v", end.Format("2006-01-02 15:04:05"), r.Method, r.URL.Path, duration)
		fmt.Println(logMessage)

		logToFile(logMessage)
	})
}

func logToFile(message string) {
	file, err := os.OpenFile("api.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error opening log file:", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(message + "\n")
	if err != nil {
		fmt.Println("Error writing to log file:", err)
	}
}
