# A simple http server in golang

### "Do one thing and do it well." 

Http server in golang that lets you save documents, edit/create and view them. 

`go run main.go`

`localhost:8080/save/{some_title}`
* 302 with location to "/view/{some_title} if created with success

`localhost:8080/edit/{some_title}`
* 200

`localhost:8080/view/{some_title}`
* 302 with location to "/view/{some_title} if not found
* 200

### ToDo

- [ ] Store templates in tmpl/ and page data in data/.

- [ ] Add a handler to make the web root redirect to /view/FrontPage.

- [ ] Spruce up the page templates by making them valid HTML and adding some CSS rules.

- [ ] Implement inter-page linking by converting instances of [PageName] to 
`<a href="/view/PageName">PageName</a>`. (hint: regexp.ReplaceAllFunc)
