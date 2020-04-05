### Chat

### Dependencies
You need to have **mysql server** and **existing database**  
Also you need to fill `server/config.json` file with your data  

And get packages from github by `go get`
```
github.com/fanliao/go-promise
github.com/jinzhu/gorm
github.com/jinzhu/gorm/dialects/mysql
```

#### How to use
In first terminal start server
```
cd ./server
go run ./
```

In other terminals start clients
```
go run ./client/ username password
```