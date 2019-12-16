module github.com/noah-blockchain/noah-explorer-api

go 1.12

replace (
	github.com/ugorji/go v1.1.4 => github.com/ugorji/go v0.0.0-20190204201341-e444a5086c43
	mellium.im/sasl v0.2.1 => github.com/mellium/sasl v0.2.1
)

require (
	github.com/AlekSi/pointer v1.1.0
	github.com/gin-contrib/cors v1.3.0
	github.com/gin-gonic/gin v1.4.0
	github.com/go-pg/pg v8.0.5+incompatible
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jpillora/backoff v1.0.0 // indirect
	github.com/noah-blockchain/coinExplorer-tools v0.0.0-20191122122608-9900f7b55b3a
	github.com/noah-blockchain/noah-go-node v0.2.0
	github.com/ugorji/go v1.1.7 // indirect
	golang.org/x/time v0.0.0-20190921001708-c4c64cad1fd0
	gopkg.in/go-playground/validator.v8 v8.18.2
	gopkg.in/guregu/null.v3 v3.4.0
	mellium.im/sasl v0.2.1 // indirect
)
