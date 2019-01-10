# Pastory.me API (Golang on AWS Lambda)
Pasry.me is the diary web service. 
https://pastory.me/

Web project is here: https://github.com/hareku/pastory-web

## Commands
### Development
```
# install dependencies
$ dep ensure

# build
$ make build

# run dynamodb-local
$ docker-compose up -d
# create Diaries table
aws dynamodb create-table --cli-input-json file://./dynamodb/diaries.json --endpoint-url http://localhost:8000

# run dev server
sam local start-api --port 8080 --docker-network pastory-api_default --env-vars env.json
```

### Deployment
```
# make package
$ aws cloudformation package --template-file .\template.yaml --output-template-file packaged.yaml --s3-bucket pastory-api

# deploy
$ aws cloudformation deploy --template-file ./packaged.yaml --stack-name pastory-api --capabilities CAPABILITY_IAM
```
