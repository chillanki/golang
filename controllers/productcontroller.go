package controllers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"encoding/json"

	python3 "github.com/go-python/cpy3"

	"github.com/gorilla/mux"
)

var ErrProductNotFound = fmt.Errorf("Product not found")

var productList = []*entities.Product{
	&entities.Product{
		ID:          "1",
		Name:        "Latte",
		Description: "Frothy milky coffee",
		Price:       2.45,
	},
	&entities.Product{
		ID:          "2",
		Name:        "Esspresso",
		Description: "Short and strong coffee without milk",
		Price:       1.99,
	},
}

func GetProductById(w http.ResponseWriter, r *http.Request) {
	productId := mux.Vars(r)["parameters"]
	i := findIndexByProductID(productId)
	np := *productList[i]
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(np)

}

func GetProducts(w http.ResponseWriter, r *http.Request) {
	var products []entities.Product
	database.Instance.Find(&products)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)
}

func findIndexByProductID(id string) int {
	for i, p := range productList {
		if p.ID == id {
			return i
		}
	}

	return -1
}

func GetParameters(w http.ResponseWriter, r *http.Request) {

	data, err := ioutil.ReadAll(r.Body)
	fmt.Println(string(data), err)
	querystring, err := BuildQuery()
	if err != nil {
		log.Fatalf("Error building query: %s", err)
	}
	json.NewEncoder(w).Encode(querystring)

}

func BuildQuery() (string, error) {
	defer python3.Py_Finalize()
	python3.Py_Initialize()
	if !python3.Py_IsInitialized() {
		fmt.Println("Error initializing the python interpreter")
		os.Exit(1)
	}

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	// we could also use PySys_GetObject("path") + PySys_SetPath,
	//but this is easier (at the cost of less flexible error handling)
	fmt.Println(dir)
	dira := "C:/Users/rc255085/Projs/dummy/product-api"
	dirb := "C:/Users/rc255085/Projs/dummy/dist"
	dirc := "C:/Users/rc255085/AppData/Local/Programs/Python/Python37/Lib/site-packages"
	ret := python3.PyRun_SimpleString("import sys\nsys.path.append(\"" + dira + "\")\nsys.path.append(r\"" + dir + "\")\nsys.path.append(\"" + dirb + "\")\nsys.path.append(\"" + dirc + "\")")
	if ret != 0 {
		log.Fatalf("error appending '%s' to python sys.path", dir)
	}

	oImport := python3.PyImport_ImportModule("teradata_analytic_lib") //ret val: new ref
	if !(oImport != nil && python3.PyErr_Occurred() == nil) {
		python3.PyErr_Print()
		log.Fatal("failed to import module 'teradata_analytic_lib'")
	}

	defer oImport.DecRef()

	oModule := python3.PyImport_AddModule("teradata_analytic_lib") //ret val: borrowed ref (from oImport)

	if !(oModule != nil && python3.PyErr_Occurred() == nil) {
		python3.PyErr_Print()
		log.Fatal("failed to add module 'teradata_analytic_lib'")
	}

	oDict := python3.PyModule_GetDict(oModule) //ret val: Borrowed
	if !(oDict != nil && python3.PyErr_Occurred() == nil) {
		python3.PyErr_Print()
		return "null", fmt.Errorf("could not get dict for module")
	}
	buildQuery := python3.PyDict_GetItemString(oDict, "buildQuery") //retval: Borrowed
	if !(buildQuery != nil && python3.PyCallable_Check(buildQuery)) {
		return "null", fmt.Errorf("could not find function 'buildQuery'")
	}
	testdataPy := buildQuery.CallObject(nil) //retval: New reference
	if !(testdataPy != nil && python3.PyErr_Occurred() == nil) {
		python3.PyErr_Print()
		return "null", fmt.Errorf("error calling function buildQuery")
	}
	defer testdataPy.DecRef()
	testdataGo := python3.PyUnicode_AsUTF8(testdataPy)
	return testdataGo, nil
}
