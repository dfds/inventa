[CmdletBinding()]

param(    
  [Parameter(Mandatory = $False, Position = 0, ValueFromPipeline = $false)]
  [string]
  $TenantId = "",

  [Parameter(Mandatory = $False, Position = 1, ValueFromPipeline = $false)]
  [string]
  $ClientId = "",

  [Parameter(Mandatory = $False, Position = 2, ValueFromPipeline = $false)]
  [string]
  $ListenAddress = "0.0.0.0",

  [Parameter(Mandatory = $False, Position = 3, ValueFromPipeline = $false)]
  [string]
  $EnableServiceProxyController = "false",

  [Parameter(Mandatory = $False, Position = 4, ValueFromPipeline = $false)]
  [string]
  $EnableIngressProxyAnnotationController = "true",

  [Parameter(Mandatory = $False, Position = 5, ValueFromPipeline = $false)]
  [string]
  $EnableServiceProxyAnnotationController = "true",

  [Parameter(Mandatory = $False, Position = 6, ValueFromPipeline = $false)]
  [string]
  $EnableHttpApi = "true",

  [Parameter(Mandatory = $False, Position = 7, ValueFromPipeline = $false)]
  [string]
  $ApiEnableAuth = "false",

  [Parameter(Mandatory = $False, Position = 8, ValueFromPipeline = $false)]
  [string]
  $ApiEnableCrossplane = "false",

  [Parameter(Mandatory = $False, Position = 9, ValueFromPipeline = $false)]
  [string]
  $CrossplaneEndpoint = "",

  [Parameter(Mandatory = $False, Position = 10, ValueFromPipeline = $false)]
  [string]
  $CrossplaneK8sToken = $null
)
$env:INVENTA_OPERATOR_AUTH_TENANT_ID="$($TenantId)"
$env:INVENTA_OPERATOR_AUTH_CLIENT_ID="$($ClientId)"
$env:INVENTA_OPERATOR_LISTEN_ADDRESS="$($ListenAddress)"
$env:INVENTA_OPERATOR_ENABLE_SERVICEPROXY_CONTROLLER="$($EnableServiceProxyController)"
$env:INVENTA_OPERATOR_ENABLE_INGRESSPROXY_ANNOTATION_CONTROLLER="$($EnableIngressProxyAnnotationController)"
$env:INVENTA_OPERATOR_ENABLE_SERVICEPROXY_ANNOTATION_CONTROLLER ="$($EnableServiceProxyAnnotationController)"
$env:INVENTA_OPERATOR_ENABLE_HTTP_API="$($EnableHttpApi)"
$env:INVENTA_OPERATOR_API_ENABLE_AUTH="$($ApiEnableAuth)"
$env:INVENTA_OPERATOR_API_ENABLE_CROSSPLANE="$($ApiEnableCrossplane)"
$env:INVENTA_CROSSPLANE_ENDPOINT="$($CrossplaneEndpoint)"
$env:INVENTA_CROSSPLANE_TOKEN="$($CrossplaneK8sToken)"

go run main.go

$env:INVENTA_OPERATOR_AUTH_TENANT_ID=""
$env:INVENTA_OPERATOR_AUTH_CLIENT_ID=""
$env:INVENTA_OPERATOR_LISTEN_ADDRESS=""
$env:INVENTA_OPERATOR_ENABLE_SERVICEPROXY_CONTROLLER=""
$env:INVENTA_OPERATOR_ENABLE_INGRESSPROXY_ANNOTATION_CONTROLLER=""
$env:INVENTA_OPERATOR_ENABLE_SERVICEPROXY_ANNOTATION_CONTROLLER =""
$env:INVENTA_OPERATOR_ENABLE_HTTP_API=""
$env:INVENTA_OPERATOR_API_ENABLE_AUTH=""
$env:INVENTA_OPERATOR_API_ENABLE_CROSSPLANE=""
$env:INVENTA_CROSSPLANE_ENDPOINT=""
$env:INVENTA_CROSSPLANE_TOKEN=""