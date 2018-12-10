module github.com/shawnfeng/sutil

require (
	code.google.com/p/go-uuid v0.0.0
	code.google.com/p/goprotobuf v0.0.0
	github.com/BurntSushi/toml v0.3.1 // indirect

	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/fzzy/radix v0.4.9-0.20141113025130-a3a55de9c594
	github.com/julienschmidt/httprouter v1.0.1-0.20150106073633-b55664b9e920
	github.com/kaneshin/go-pkg v0.0.0-20150919125626-a8e1479186cf
	github.com/kr/pretty v0.1.0 // indirect
	github.com/pkg/errors v0.8.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/sdming/gosnow v0.0.0-20130403030620-3a05c415e886
	github.com/stretchr/testify v1.2.3-0.20181014000028-04af85275a5c // indirect
	github.com/stretchrcom/testify v1.2.2 // indirect
	github.com/vaughan0/go-ini v0.0.0-20130923145212-a98ad7ee00ec

	go.uber.org/atomic v1.3.2 // indirect
	go.uber.org/multierr v1.1.0 // indirect
	go.uber.org/zap v1.9.1
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect

	gopkg.in/mgo.v2 v2.0.0-20141107142503-e2e914857713
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	gopkg.in/yaml.v2 v2.2.2 // indirect

)

replace (
	code.google.com/p/go-uuid => github.com/shawnfeng/googleuuid v1.0.0
	code.google.com/p/goprotobuf => github.com/shawnfeng/googlpb v1.0.0
)
