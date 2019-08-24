package consts

const (
	urlApiPrefix      = "/api"
	urlPrefixVersion1 = "/v1"
	urlPrefixService  = "/service"
	urlPrefixUser     = "/user"
	urlPrefixCfg      = "/cfg"

	urlLogin  = "/login"
	urlLogout = "/logout"

	UrlApiLoginV1 = urlPrefixVersion1 + urlApiPrefix + urlLogin
	UrlLogin      = urlLogin

	UrlServiceV1 = urlPrefixVersion1 + urlApiPrefix + urlPrefixService

	UrlUserV1 = urlPrefixVersion1 + urlApiPrefix + urlPrefixUser

	UrlCfgV1 = urlPrefixVersion1 + urlApiPrefix + urlPrefixCfg
)
