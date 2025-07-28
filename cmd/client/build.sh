#!/usr/bin/bash
oss=(darwin linux windows)
archs=(amd64 arm64)
ldatgs="-X main.buildVersion=v1.0.1 -X 'main.buildTime=$(date +'%Y/%m/%d %H:%M:%S')' -X 'main.buildCommit=$(git rev-parse --short HEAD)'"

for os in ${oss[@]}
do
    for arch in ${archs[@]}
    do
        env GOOS=${os} GOARCH=${arch} go build -ldflags "${ldatgs}" -o keeper_${os}_${arch}
    done
done

#env CGO_ENABLED=1 GOOS=ios GOARCH=${arch} go build -o keeper_ios_${arch}


        # env GOOS=${os} GOARCH=${arch} go build -ldflags "     \
        # -X main.buildVersion=v1.0.1                           \
        # -X 'main.buildTime=$$(date +'%Y/%m/%d %H:%M:%S')'     \
        # -X 'main.buildCommit=$$(git rev-parse --short HEAD)'" \
        # -o keeper_${os}_${arch}