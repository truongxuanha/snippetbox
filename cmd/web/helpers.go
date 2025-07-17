package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

// Trình trợ giúp serverError viết thông báo lỗi và dấu vết ngăn xếp vào errorLog
// Sau đó gửi phản hồi lỗi máy chủ nộ bộ 500 chung đến user
func (app *application) serverError(w http.ResponseWriter, err error){
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Println(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
func (app *application) clientError(w http.ResponseWriter, status int){
	http.Error(w, http.StatusText(status), status)
}
func (app *application) notFound(w http.ResponseWriter)  {
	app.clientError(w, http.StatusNotFound)
}
func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData){
	ts, ok := app.templateData[page]

	if !ok {
		err := fmt.Errorf("the template %s does not exits", page)
		app.serverError(w, err)
		return
	}
	w.WriteHeader(status)
	err := ts.ExecuteTemplate(w, "base", data)

	if err != nil {
		app.serverError(w, err)
	}
}