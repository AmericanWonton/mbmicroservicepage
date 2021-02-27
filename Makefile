run:
	go run *.go
gomod-exp:
	export GO111MODULE=on
gobuild:
	GOOS=linux GOARCH=amd64 go build -o mainwebpage
dockerbuild:
	docker build -t mainwebpage .
dockerbuildandpush:
	docker build -t mainwebpage .
	docker tag mainwebpage americanwonton/mainwebpage
	docker push americanwonton/mainwebpage
dockerrun:
	docker run -it -p 8080:8080 mainwebpage
dockerrundetached:
	docker run -d -p 8080:8080 mainwebpage
dockerrunitvolume:
	docker run -it -p 8080:8080 -v photo-images:/static/images mainwebpage
dockerrundetvolume:
	docker run -d -p 8080:8080 -v photo-images:/static/images mainwebpage
dockertagimage:
	docker tag mainwebpage americanwonton/mainwebpage
dockerimagepush:
	docker push americanwonton/mainwebpage
dockerallpush:
	docker tag mainwebpage americanwonton/mainwebpage
	docker push americanwonton/mainwebpage
dockerseeshell:
	docker run -it mainwebpage sh