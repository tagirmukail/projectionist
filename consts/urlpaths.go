package consts

const (
	urlApiPrefix      = "/api"
	urlPrefixVersion1 = "/v1"
	urlPrefixService  = "/service"
	urlPrefixUser     = "/user"
	urlPrefixCfg      = "/cfg"

	urlLogin  = "/login"
	urlLogout = "/logout"

	UrlApiLogin = urlPrefixVersion1 + urlApiPrefix + urlLogin
	UrlLogin    = urlLogin

	UrlService = urlPrefixVersion1 + urlApiPrefix + urlPrefixService

	UrlUser = urlPrefixVersion1 + urlApiPrefix + urlPrefixUser

	UrlCfg = urlPrefixVersion1 + urlApiPrefix + urlPrefixCfg
)
