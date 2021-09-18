set GOARCH=amd64

set GOOS=linux
go build -o covid_data_ita .

set GOOS=windows
go build -o covid_data_ita.exe .