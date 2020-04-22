all : generate

generate :
	go install github.com/selabhvl/gotestgen/testgen
	go generate
