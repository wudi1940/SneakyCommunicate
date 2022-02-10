cd /Users/kangnan.peng/GolandProjects/SneakyCommunicate/server/
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ../MySrv cmd/serverCmd.go

cd /Users/kangnan.peng/GolandProjects/SneakyCommunicate/client/
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ../cli.exe cmd/clientCmd.go
go build -o ../cli cmd/clientCmd.go

cd /Users/kangnan.peng/GolandProjects/SneakyCommunicate/web/
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ../WebSrv cliDownload.go