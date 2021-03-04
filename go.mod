module github.com/TheCacophonyProject/fake-thermal-camera

require (
	github.com/TheCacophonyProject/go-config v1.4.0
	github.com/TheCacophonyProject/go-cptv v0.0.0-20200225002107-8095b1b6b929
	github.com/TheCacophonyProject/lepton3 v0.0.0-20200213011619-1934a9300bd3
	github.com/TheCacophonyProject/thermal-recorder v1.22.1-0.20200225033227-2090330c5c11
	github.com/alexflint/go-arg v1.1.0
	github.com/godbus/dbus v4.1.0+incompatible
	github.com/gorilla/mux v1.7.4
	golang.org/x/net v0.0.0-20201021035429-f5854403a974 // indirect
	golang.org/x/sys v0.0.0-20210119212857-b64e53b001e4 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/yaml.v1 v1.0.0-20140924161607-9f9df34309c0
)

replace periph.io/x/periph => github.com/TheCacophonyProject/periph v2.0.1-0.20171123021141-d06ef89e37e8+incompatible

go 1.15
