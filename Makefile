# This Makefile is meant to be used by people that do not usually work
# with Go source code. If you know what GOPATH is then you probably
# don't need to bother with make.

.PHONY: gmxt android ios gmxt-cross evm all test clean
.PHONY: gmxt-linux gmxt-linux-386 gmxt-linux-amd64 gmxt-linux-mips64 gmxt-linux-mips64le
.PHONY: gmxt-linux-arm gmxt-linux-arm-5 gmxt-linux-arm-6 gmxt-linux-arm-7 gmxt-linux-arm64
.PHONY: gmxt-darwin gmxt-darwin-386 gmxt-darwin-amd64
.PHONY: gmxt-windows gmxt-windows-386 gmxt-windows-amd64

GOBIN = ./build/bin
GO ?= latest
GORUN = env GO111MODULE=on go run

gmxt:
	$(GORUN) build/ci.go install ./cmd/gmxt
	@echo "Done building."
	@echo "Run \"$(GOBIN)/gmxt\" to launch gmxt."

all:
	$(GORUN) build/ci.go install

android:
	$(GORUN) build/ci.go aar --local
	@echo "Done building."
	@echo "Import \"$(GOBIN)/gmxt.aar\" to use the library."
	@echo "Import \"$(GOBIN)/gmxt-sources.jar\" to add javadocs"
	@echo "For more info see https://stackoverflow.com/questions/20994336/android-studio-how-to-attach-javadoc"
	
ios:
	$(GORUN) build/ci.go xcode --local
	@echo "Done building."
	@echo "Import \"$(GOBIN)/Gmxt.framework\" to use the library."

test: all
	$(GORUN) build/ci.go test

lint: ## Run linters.
	$(GORUN) build/ci.go lint

clean:
	env GO111MODULE=on go clean -cache
	rm -fr build/_workspace/pkg/ $(GOBIN)/*

# The devtools target installs tools required for 'go generate'.
# You need to put $GOBIN (or $GOPATH/bin) in your PATH to use 'go generate'.

devtools:
	env GOBIN= go get -u golang.org/x/tools/cmd/stringer
	env GOBIN= go get -u github.com/kevinburke/go-bindata/go-bindata
	env GOBIN= go get -u github.com/fjl/gencodec
	env GOBIN= go get -u github.com/golang/protobuf/protoc-gen-go
	env GOBIN= go install ./cmd/abigen
	@type "npm" 2> /dev/null || echo 'Please install node.js and npm'
	@type "solc" 2> /dev/null || echo 'Please install solc'
	@type "protoc" 2> /dev/null || echo 'Please install protoc'

# Cross Compilation Targets (xgo)

gmxt-cross: gmxt-linux gmxt-darwin gmxt-windows gmxt-android gmxt-ios
	@echo "Full cross compilation done:"
	@ls -ld $(GOBIN)/gmxt-*

gmxt-linux: gmxt-linux-386 gmxt-linux-amd64 gmxt-linux-arm gmxt-linux-mips64 gmxt-linux-mips64le
	@echo "Linux cross compilation done:"
	@ls -ld $(GOBIN)/gmxt-linux-*

gmxt-linux-386:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/386 -v ./cmd/gmxt
	@echo "Linux 386 cross compilation done:"
	@ls -ld $(GOBIN)/gmxt-linux-* | grep 386

gmxt-linux-amd64:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/amd64 -v ./cmd/gmxt
	@echo "Linux amd64 cross compilation done:"
	@ls -ld $(GOBIN)/gmxt-linux-* | grep amd64

gmxt-linux-arm: gmxt-linux-arm-5 gmxt-linux-arm-6 gmxt-linux-arm-7 gmxt-linux-arm64
	@echo "Linux ARM cross compilation done:"
	@ls -ld $(GOBIN)/gmxt-linux-* | grep arm

gmxt-linux-arm-5:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/arm-5 -v ./cmd/gmxt
	@echo "Linux ARMv5 cross compilation done:"
	@ls -ld $(GOBIN)/gmxt-linux-* | grep arm-5

gmxt-linux-arm-6:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/arm-6 -v ./cmd/gmxt
	@echo "Linux ARMv6 cross compilation done:"
	@ls -ld $(GOBIN)/gmxt-linux-* | grep arm-6

gmxt-linux-arm-7:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/arm-7 -v ./cmd/gmxt
	@echo "Linux ARMv7 cross compilation done:"
	@ls -ld $(GOBIN)/gmxt-linux-* | grep arm-7

gmxt-linux-arm64:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/arm64 -v ./cmd/gmxt
	@echo "Linux ARM64 cross compilation done:"
	@ls -ld $(GOBIN)/gmxt-linux-* | grep arm64

gmxt-linux-mips:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/mips --ldflags '-extldflags "-static"' -v ./cmd/gmxt
	@echo "Linux MIPS cross compilation done:"
	@ls -ld $(GOBIN)/gmxt-linux-* | grep mips

gmxt-linux-mipsle:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/mipsle --ldflags '-extldflags "-static"' -v ./cmd/gmxt
	@echo "Linux MIPSle cross compilation done:"
	@ls -ld $(GOBIN)/gmxt-linux-* | grep mipsle

gmxt-linux-mips64:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/mips64 --ldflags '-extldflags "-static"' -v ./cmd/gmxt
	@echo "Linux MIPS64 cross compilation done:"
	@ls -ld $(GOBIN)/gmxt-linux-* | grep mips64

gmxt-linux-mips64le:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/mips64le --ldflags '-extldflags "-static"' -v ./cmd/gmxt
	@echo "Linux MIPS64le cross compilation done:"
	@ls -ld $(GOBIN)/gmxt-linux-* | grep mips64le

gmxt-darwin: gmxt-darwin-386 gmxt-darwin-amd64
	@echo "Darwin cross compilation done:"
	@ls -ld $(GOBIN)/gmxt-darwin-*

gmxt-darwin-386:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=darwin/386 -v ./cmd/gmxt
	@echo "Darwin 386 cross compilation done:"
	@ls -ld $(GOBIN)/gmxt-darwin-* | grep 386

gmxt-darwin-amd64:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=darwin/amd64 -v ./cmd/gmxt
	@echo "Darwin amd64 cross compilation done:"
	@ls -ld $(GOBIN)/gmxt-darwin-* | grep amd64

gmxt-windows: gmxt-windows-386 gmxt-windows-amd64
	@echo "Windows cross compilation done:"
	@ls -ld $(GOBIN)/gmxt-windows-*

gmxt-windows-386:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=windows/386 -v ./cmd/gmxt
	@echo "Windows 386 cross compilation done:"
	@ls -ld $(GOBIN)/gmxt-windows-* | grep 386

gmxt-windows-amd64:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=windows/amd64 -v ./cmd/gmxt
	@echo "Windows amd64 cross compilation done:"
	@ls -ld $(GOBIN)/gmxt-windows-* | grep amd64
