module github.com/cardiacsociety/web-services

go 1.12

// +heroku goVersion go1.11
// +heroku install ./...

require (
	github.com/34South/envr v0.0.0-20170706023707-e57a7716f427
	github.com/360EntSecGroup-Skylar/excelize v1.4.1
	github.com/8o8/email v0.1.0
	github.com/aws/aws-sdk-go v1.19.11
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-playground/locales v0.12.1 // indirect
	github.com/go-playground/universal-translator v0.16.0 // indirect
	github.com/go-sql-driver/mysql v1.4.1
	github.com/gorilla/mux v1.7.1
	github.com/graphql-go/graphql v0.7.8
	github.com/graphql-go/handler v0.2.3
	github.com/hashicorp/go-uuid v1.0.1
	github.com/imdario/mergo v0.3.7
	github.com/jung-kurt/gofpdf v1.0.2
	github.com/leodido/go-urn v1.1.0 // indirect
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/pkg/errors v0.8.1
	github.com/rs/cors v1.6.0
	github.com/sendgrid/sendgrid-go v3.4.1+incompatible
	github.com/urfave/negroni v1.0.0
	gopkg.in/go-playground/validator.v9 v9.28.0
	gopkg.in/mgo.v2 v2.0.0-20180705113604-9856a29383ce
)
