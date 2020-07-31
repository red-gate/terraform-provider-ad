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
try {
    "Creating terraform-provider-ad_test_network network"
    docker network create terraform-provider-ad_test_network

    "Starting openldap container"
    docker run `
        --name test-ldap-server `
        --network terraform-provider-ad_test_network `
        -v ${pwd}/scripts/openldap:/customizations `
        --env LDAP_SEED_INTERNAL_SCHEMA_PATH=/customizations/schema `
        --env LDAP_SEED_INTERNAL_LDIF_PATH=/customizations/ldif `
        --detach osixia/openldap:1.4.0 --loglevel trace

    Start-Sleep 10

    "Starting tests"
    docker run --rm `
        -v ${pwd}:/app `
        -w /app `
        --network terraform-provider-ad_test_network `
        -e TF_ACC=1 `
        -e AD_DOMAIN=example.org `
        -e AD_USER_DOMAIN=example.org `
        -e AD_COMPUTER_DOMAIN=example.org `
        -e AD_OU_DOMAIN=example.org `
        -e "AD_COMPUTER_OU_DISTINGUISHED_NAME=dc=example,dc=org" `
        -e "AD_GROUP_OU_DISTINGUISHED_NAME=dc=example,dc=org" `
        -e AD_URL=ldap://test-ldap-server:389 `
        -e AD_USER="cn=admin,dc=example,dc=org" `
        -e AD_PASSWORD=admin `
        golang:1.14 `
        /bin/sh -c "go test -short -v ./ad"

} finally {
    "Stopping openldap container"
    docker logs test-ldap-server
    docker rm test-ldap-server --force
    "Removing terraform-provider-ad_test_network network"
    docker network rm terraform-provider-ad_test_network
}
