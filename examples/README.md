## init

```
go mod init bintray
```

## run

```
go build && ./bintray

or

go run jfrog_client.go package_config_new.yaml
```

## debug

```
dlv debug --build-flags="-gcflags='-N -l'" jfrog_client.go 

config max-string-len 99999
config -list

break jfrog_client.go:103
break jfrog_client.go:67
break github.com/jfrog/jfrog-client-go/utils/log.Info
break github.com/jfrog/jfrog-client-go/bintray/services/packages.(*PackageService).Create

c
p <var>
n
r
```

## build latest deps into go.mod

```
rm go.mod go.sum
go mod init bintray
go get -u github.com/arcolife/jfrog-client-go@<latest commit on master branch>
```
