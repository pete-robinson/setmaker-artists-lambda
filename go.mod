module github.com/pete-robinson/setmaker-artist-meta-lambda

go 1.18

require (
	github.com/aws/aws-lambda-go v1.34.1
	github.com/aws/aws-sdk-go-v2 v1.17.1
	github.com/aws/aws-sdk-go-v2/config v1.18.0
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue v1.10.3
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.17.4
	github.com/golang/protobuf v1.5.2
	github.com/google/uuid v1.3.0
	github.com/pete-robinson/setmaker-proto v1.0.4
	github.com/sirupsen/logrus v1.9.0
	github.com/zmb3/spotify/v2 v2.3.0
	golang.org/x/oauth2 v0.0.0-20210810183815-faf39c7919d5
	google.golang.org/grpc v1.50.1
)

require (
	github.com/aws/aws-sdk-go-v2/credentials v1.13.0 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.12.19 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.1.25 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.4.19 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.3.26 // indirect
	github.com/aws/aws-sdk-go-v2/service/dynamodbstreams v1.13.23 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.9.10 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.7.19 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.9.19 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.11.25 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.13.8 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.17.2 // indirect
	github.com/aws/smithy-go v1.13.4 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	golang.org/x/net v0.0.0-20210813160813-60bc85c4be6d // indirect
	golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8 // indirect
	golang.org/x/text v0.3.6 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20200825200019-8632dd797987 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
)

replace github.com/pete-robinson/set-maker-grpc => /Users/peterobinson/go/src/github.com/pete-robinson/set-maker-grpc
