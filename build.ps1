[CmdletBinding()]
param(
    [Parameter(Mandatory = $true)]
    $Version
)

docker run --rm `
    -v ${pwd}:/app `
    -w /app `
    -e GOOS=windows `
    -e GOARCH=amd64 `
    golang:1.14 `
    /bin/sh -c "go build -o terraform-provider-ad_v${Version}.exe"