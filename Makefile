# GOOS：darwin、freebsd、linux、windows
# GOARCH：386、amd64、arm、s390x

# parameters
GORUN=go run
GOBUILD=go build -o
OUTPUTDIR=./bin
APPNAME=poker
ENTRY=main.go

.PHONY: darwin linux windows publish clean server client

all: publish

publish:
	@CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) $(OUTPUTDIR)/$(APPNAME)-mac $(ENTRY)
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(OUTPUTDIR)/$(APPNAME) $(ENTRY)
	@CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) $(OUTPUTDIR)/$(APPNAME).exe $(ENTRY)

darwin:
	@CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) $(OUTPUTDIR)/$(APPNAME) $(ENTRY)

linux:
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(OUTPUTDIR)/$(APPNAME) $(ENTRY)

windows:
	@CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) $(OUTPUTDIR)/$(APPNAME).exe $(ENTRY)

clean:
	@rm -rf ./bin

server:
	@$(GORUN) $(ENTRY) --service=server

client:
	@$(GORUN) $(ENTRY) --service=client
