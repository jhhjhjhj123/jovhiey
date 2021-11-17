#!/usr/bin/env bash

package="ariga.io/atlas/cmd/atlas"
package_split=(${package//\// })
package_name=${package_split[3]}

platforms=("windows/amd64" "darwin/amd64" "darwin/arm64" "linux/amd64")
TAG=$(git describe --tags --abbrev=0)
if [[ -z "$TAG" ]]; then
  echo "branch must be tagged!"
  exit 1
fi

for platform in "${platforms[@]}"
do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}
    output_name=$package_name'-'$GOOS'-'$GOARCH'-'$TAG
    if [ $GOOS = "windows" ]; then
        output_name+='.exe'
    fi

    env GOOS=$GOOS GOARCH=$GOARCH go build -o $output_name $package
    if [ $? -ne 0 ]; then
        echo 'An error has occurred building! Aborting the script execution...'
        exit 1
    fi
    aws s3 cp $output_name s3://release.ariga.io/atlas/$output_name
    if [ $? -ne 0 ]; then
        echo 'An error has occurred uploading to aws! Aborting the script execution...'
        exit 1
    fi
done
