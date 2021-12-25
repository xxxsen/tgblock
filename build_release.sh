#!/bin/bash

output=./release 

rm ${output} -rf
mkdir -p ${output}

oss="windows linux darwin freebsd"
archs="amd64 386 arm"

function build {
    os=$1
    arch=$2
    name="${os}_${arch}"
    if [ "$os" == "windows" ]; then 
        nameext="${name}.exe"
    fi
    svr="${output}/tgblock_svr_${nameext}"
    cli="${output}/tgblock_cli_${nameext}"
    tarfile="${output}/tgblock_${name}.tar.xz"
    CGO_ENABLED=0 GOOS=$os GOARCH=$arch go build -a -tags netgo -ldflags '-w' -o ${svr} ./cmd/svr
    CGO_ENABLED=0 GOOS=$os GOARCH=$arch go build -a -tags netgo -ldflags '-w' -o ${cli} ./cmd/svr
    tar -cJf ${tarfile} ${svr} ${cli}
    rm ${svr} ${cli}
}

for os in $(echo ${oss})
do
    for arch in $(echo ${archs})
    do 
        echo "build binary os:${os}, arch:${arch}..."
        build $os $arch 
    done 
done 

