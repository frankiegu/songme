# songme

A Go web app that recommends songs.
  
- Shows random songs to user.
- Supports lots of different providers such as youtube, soundcloud etc. for adding songs.
- Has authentication for users with different roles such as moderate, admin.
- Has admin dashboard to confirm or delete songs added by users.
- Deployed to heroku.

https://songtome.herokuapp.com

## Installing

In order to download, run the following command:

```$ go get github.com/emredir/songme```

After that you may want to check **config.go** file. Here, you see some configuration variables that corresponds to environment variables. You may want to update database url by modifying following line: 

`postgres://username:password@localhost:5432/dbname?sslmode=disable` 

## Running

After your configurations are ready, issue the following commands in order to run:

```
	$ cd $GOPATH/src/github.com/emredir/songme
	$ go run cmd/songme/main.go
```