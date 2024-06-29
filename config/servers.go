package config

import (
	"fmt"
	"mime"
	"net/http"
	"path/filepath"
)

func StartFileServer() {
	// Directorio que se va a servir
	fileServer := http.FileServer(http.Dir("./storage"))
	// Aplicar middleware CORS
	corsFileServer := corsMiddleware(fileServer)
	// Ruteo de archivos
	http.Handle("/", corsFileServer)
	// Dirección y puerto del servidor
	addr := ":8000"
	fmt.Printf("Servidor de archivos escuchando en %s...", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		panic(err)
	}
}

// Middleware para agregar encabezados CORS
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			return
		}
		next.ServeHTTP(w, r)
	})
}

func fileHandler(w http.ResponseWriter, r *http.Request) {
	filePath := "." + r.URL.Path
	ext := filepath.Ext(filePath)
	mimeType := mime.TypeByExtension(ext)
	w.Header().Set("Content-Type", mimeType)
	http.ServeFile(w, r, filePath)
}
