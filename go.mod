module github.com/epiphany-platform/m-azure-kubernetes-service

go 1.15

//TODO fix
replace github.com/epiphany-platform/e-structures => ../../epiphany-platform/e-structures

require (
	github.com/epiphany-platform/e-structures v0.0.5
	github.com/google/go-cmp v0.5.3
	github.com/jinzhu/copier v0.2.3
	github.com/magiconair/properties v1.8.4 // indirect
	github.com/mkyc/go-terraform v0.0.7
	github.com/pelletier/go-toml v1.8.1 // indirect
	github.com/rs/zerolog v1.20.0
	github.com/spf13/afero v1.4.1 // indirect
	github.com/spf13/cast v1.3.1 // indirect
	github.com/spf13/cobra v1.1.1
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/viper v1.7.1
	golang.org/x/sys v0.0.0-20201106081118-db71ae66460a // indirect
	golang.org/x/text v0.3.4 // indirect
	gopkg.in/ini.v1 v1.62.0 // indirect
)
