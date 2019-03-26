module marwan.io/authproxy

go 1.12

require (
	github.com/kelseyhightower/envconfig v1.3.0
	golang.org/x/oauth2 v0.0.0-20190130055435-99b60b757ec1
)

replace golang.org/x/oauth2 => github.com/wlhee/oauth2 v0.0.0-20190308230854-b33c8d1d8308
