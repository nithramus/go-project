package main

import (
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

// func main() {
// 	fmt.Printf("hello, world\n")
// 	resp, err := http.Get("http://google.com")
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	fmt.Println(resp.Body)
// 	defer resp.Body.Close()
// 	body, err := ioutil.ReadAll(resp.Body)
// 	s := string(body)
// 	fmt.Println(s)
// }

type Data struct {
	Price int
	Name  string
}

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func middlewareTest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		fmt.Println("log")
		fmt.Println(reflect.TypeOf(req.Body))

		// body, _ := ioutil.ReadAll(req.Body)
		// // fmt.Println(string(body))
		// var data Data
		// err := json.Unmarshal([]byte(body), &data)
		// if err != nil || data.Price == 0 {
		// 	fmt.Print(err)
		next.ServeHTTP(res, req)
		// 	return
		// }
		// fmt.Printf("here")
		// test :=
		// req.Context()
		// req.Context.Set(r, "foo", "bar")
		// ctx := context.WithValue(req.Context(), "data", data)
		// rWithContext := req.WithContext(ctx)
		// fmt.Println(ctx)
		// // data2 := ctx.Value("data")
		// data2 := rWithContext.Context().Value("data")
		// fmt.Printf("final %+v\n", data2)
		// // fmt.Println(test)
		// next.ServeHTTP(res, rWithContext)
	})
}

func upload(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method)
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("uploadfile")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	fmt.Fprintf(w, "%v", handler.Header)
	f, err := os.OpenFile("./new"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	io.Copy(f, file)
}

func main() {
	db, err := gorm.Open("mysql", "root:1234@tcp(127.0.0.1:3306)/test?parseTime=true")
	if err != nil {
		fmt.Println(err)
		panic("failed to connect database")
	}
	defer db.Close()

	// Migrate the schema
	db.AutoMigrate(&Product{})
	db.AutoMigrate(&Data{})
	var product Product

	// Create
	db.Create(&Product{Code: "L1212", Price: 1000})
	db.First(&product, 1)
	fmt.Println(product.Price)
	db.Model(&product).Update("Price", 2000)

	mux := http.NewServeMux()
	mux.Handle("/nike", http.NotFoundHandler())
	mux.HandleFunc("/upload", upload)
	mux.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
		// var data Data
		fmt.Println("meh nan")
		var data = r.Context().Value("data").(Data)
		db.Create(&data)
		fmt.Printf("final2 %+v\n", data)
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})
	fmt.Println(mux)
	yolo := middlewareTest(mux)
	fmt.Println("Server launched")
	log.Fatal(http.ListenAndServe(":7000", yolo))
}
