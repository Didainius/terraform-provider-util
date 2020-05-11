default: install

install:
	sh -c "go build -o ~/.terraform.d/plugins/terraform-provider-util_v0.1.0"

testacc: install
	sh -c 'cd util && TF_ACC=1 go test -v'


# build: fmtcheck
# 	go install
	