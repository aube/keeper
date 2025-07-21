#!/usr/bin/bash
oss=(darwin linux windows)
archs=(amd64 arm64)

for os in ${oss[@]}
do
    for arch in ${archs[@]}
    do
        env GOOS=${os} GOARCH=${arch} go build -o keeper_${os}_${arch}
    done
done

#env CGO_ENABLED=1 GOOS=ios GOARCH=${arch} go build -o keeper_ios_${arch}