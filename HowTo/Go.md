# Go Steps

#### Needed software
https://golang.org/

#### Installs & Setup

First: Setup go modules with the git repo:  
> go mod init github.com/AndreasRoither/Snippets

Then install needed modules:
> go get -u github.com/go-chi/chi  
> go get -u github.com/go-chi/render   
> go get -u github.com/lib/pq  
> go get -u github.com/go-chi/cors  
> go get -u github.com/go-swagger/go-swagger  
For debugging:  
> go get -v github.com/go-delve/delve/cmd/dlv  

#### Go Commands

Install needed packages when checking out the repo:
> go install

Build an .exe file:  
> go build  

Run:
> go run main.go strings.go domain.go

Run tests:  
> go test -v

#### Swagger

Important:  
Check if you are using tabs for indentation. YAML doesn't allow tabs; it requires spaces.

Generate docs:  
> swagger generate spec -o ./swagger.json  

Serve swagger:  
> swagger serve .\swagger.json