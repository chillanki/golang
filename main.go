package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

var DB *gorm.DB

func main() {

	// Load Configurations from config.json using Viper
	// LoadAppConfig()

	// // Initialize Database
	// database.Connect(AppConfig.ConnectionString)
	// database.Migrate()
	////C code
	// fmt.Println("Golang going to call a C function!")
	// C.cHello()
	// fmt.Println("Golang going to call another C function!")
	// goMessage := C.CString("Sent from Golang!")
	// defer C.free(unsafe.Pointer(&goMessage))
	// C.printMessage(goMessage)
	// fmt.Println(C.printMessage(goMessage))

	//Python code

	// Initialize the router

	// Register Routes
	router := mux.NewRouter().StrictSlash(true)
	RegisterProductRoutes(router)
	srvr := &http.Server{
		Addr:    ":9090",
		Handler: router,
	}

	srvr.ListenAndServe()
	// Start the server
	// log.Println(fmt.Sprintf("Starting Server on port %s", AppConfig.Port))
	// log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", AppConfig.Port), router))

}

func RegisterProductRoutes(router *mux.Router) {
	router.HandleFunc("/", controllers.GetProducts).Methods("GET")
	router.HandleFunc("/{parameters}", controllers.GetParameters).Methods("GET")
}
