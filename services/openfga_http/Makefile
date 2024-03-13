build:
	cd src/ext;GOOS=linux GOARCH=amd64 go build -o bin/extensions/lambda-cache-layer main.go
	cd src/functions/fgahandler;GOOS=linux GOARCH=amd64 go build -o bin/bootstrap handler.go
	cd src/functions/fgadbinit;GOOS=linux GOARCH=amd64 go build -o bin/bootstrap dbinit.go
package: build
	cd src/ext/bin;zip -r extension.zip extensions/ 
	cd src/functions/fgahandler/bin;zip -r handler.zip . 
	cd src/functions/fgadbinit/bin;zip -r handler.zip . 

#deploy: build package
#	cd src/ext/bin;aws lambda publish-layer-version  --layer-name 'lambda-cache-layer' --region eu-central-1 --zip-file 'fileb://extension.zip' --profile=storer_mfa
