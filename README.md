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
