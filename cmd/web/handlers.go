package main

import (
	"errors"
	"fmt"
	"github.com/PossibleNPC/snippetbox/pkg/models"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"
)

// While this function receiver is only for logging at the moment, it is implied that we are adding additional
// functionality later on.
// I think this is an indication we should move that application struct out of main, but that might be later
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "home.page.tmpl", &templateData{
		Snippets: s,
	})

	//for _, snippet := range s {
	//	fmt.Fprintf(w, "%v\n", snippet)
	//}

	//data := &templateData{Snippets: s}

	//// For some reason, the template with the real content, NOT the layout, must come prior to that generalized template
	//// Order of these templates seem to only be for certain files, but not others. Unclear
	//files := []string {
	//	// 28Nov: Changed the following to reflect static server configuration
	//	app.staticPath + "../html/home.page.tmpl",
	//	app.staticPath + "../html/base.layout.tmpl",
	//	app.staticPath + "../html/footer.partial.tmpl",
	//	//"./ui/html/home.page.tmpl",
	//	//"./ui/html/base.layout.tmpl",
	//	//"./ui/html/footer.partial.tmpl",
	//}
	//
	//// While we can pass any number of strings, we are feeding in a slice, which is not of the desired underlying type;
	//// however, the types within that slice are of type string, so we can just unpack our container
	//ts, err := template.ParseFiles(files...)
	//if err != nil {
	//	app.serverError(w, err)
	//	return
	//}
	//
	//err = ts.Execute(w, data)
	//if err != nil {
	//	app.serverError(w, err)
	//	return
	//}
	//
	//w.Write([]byte("Hello from Snippetbox!"))
}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	s, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	app.render(w, r, "show.page.tmpl", &templateData{
		Snippet: s,
	})

	//data := &templateData{Snippet: s}
	//
	//files := []string{
	//	app.staticPath + "../html/show.page.tmpl",
	//	app.staticPath + "../html/base.layout.tmpl",
	//	app.staticPath + "../html/footer.partial.tmpl",
	//}
	//
	//ts, err := template.ParseFiles(files...)
	//if err != nil {
	//	app.serverError(w, err)
	//	return
	//}
	//
	//err = ts.Execute(w, data)
	//if err != nil {
	//	app.serverError(w, err)
	//}
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	expires := r.PostForm.Get("expires")

	errors := make(map[string]string)

	if strings.TrimSpace(title) == "" {
		errors["title"] = "This field cannot be blank"
	} else if utf8.RuneCountInString(title) > 100 {
		errors["title"] = "This field is too long (maximum is 100 characters)"
	}

	if strings.TrimSpace(content) == "" {
		errors["content"] = "This field cannot be blank"
	}

	if strings.TrimSpace(expires) == "" {
		errors["expires"] = "This field cannot be blank"
	} else if expires != "365" && expires != "7" && expires != "1" {
		errors["expires"] = "This field is invalid"
	}

	if len(errors) > 0 {
		fmt.Fprint(w, errors)
		return
	}

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}

func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tmpl", nil)
}