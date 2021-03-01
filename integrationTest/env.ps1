$Env:BIND_DOMAIN = "example.com"
Write-Host "go test -v -verbose -provider BIND"

$Env:HEXONET_DOMAIN = "a-b-c-movies.com"
$Env:HEXONET_ENTITY = "OTE"
$Env:HEXONET_PW = "test.passw0rd"
$Env:HEXONET_UID = "test.user"
Write-Host "go test -v -verbose -provider HEXONET"

$Env:MSDNS_DOMAIN = "example.com"
$Env:MSDNS_PSSESSION = ""
$Env:MSDNS_DNSSERVER = ""
# $Env:MSDNS_DNSSERVER = "TOMDEV"
# $Env:MSDNS_PSSESSION = $env:computername
Write-Host "go test -v -verbose -provider MSDNS"
