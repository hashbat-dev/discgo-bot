package ping

import (
	"fmt"
	"io"
	"net/http"
)

func Run() {
	http.HandleFunc("/", getRoot)
	err := http.ListenAndServe(":3333", nil)
	if err != nil {
		fmt.Printf("\nERROR: %v", err)
	}
}

func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Ping: got / request\n")
	_, err := io.WriteString(w, "Ping!\n")
	if err != nil {
		fmt.Println(fmt.Printf("Error Loading Ping: %v", err))
	}
}
