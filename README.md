# go-kube-client
A Go CLI based client which can manage, install and queue workflows in Kubernetes

You can build the app by running in required folder:
```
go build
```

Then you can run the app by hitting:

```
cd cli-app
./cli-app
```

```
Welcome to the world of kube-client - A CLI based client app to call Kubernetes APIs !!
NAME:
   cli-app - Using client-go effectively with Kubernetes api

USAGE:
   cli-app [global options] command [command options] [arguments...]

VERSION:
   1.0

COMMANDS:
   crud       Run CRUD example
   lister     Run lister example
   informer   Run informer example
   workqueue  Run workqueue example
   help, h    Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```
