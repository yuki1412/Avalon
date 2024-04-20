package router

import (
	"net/http"
	"wayne/service"
)

func HomeRouter() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/lobby", service.HomeHandler)
	mux.HandleFunc("/index", service.HomeHandler)
	mux.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("images"))))

	// Register AvalonRouter under /avalon prefix
	mux.Handle("/avalon/", http.StripPrefix("/avalon", AvalonRouter()))

	// // Fallback handler for unknown paths
	// mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	// Redirect all unknown paths to the "I'm a teapot" handler
	// 	http.Redirect(w, r, "/teapot", http.StatusFound)
	// })

	// // Handler for "I'm a teapot" path
	// mux.HandleFunc("/teapot", func(w http.ResponseWriter, r *http.Request) {
	// 	w.WriteHeader(http.StatusTeapot) // HTTP status code 418 - I'm a teapot
	// 	fmt.Fprintf(w, "I'm a teapot, would you like some tea or coffee?")
	// })

	return mux
}
