package consts

const (
	urlApiPrefix      = "/api"
	urlPrefixVersion1 = "/v1"
	urlPrefixService  = "/service"
	urlPrefixUser     = "/user"
	urlPrefixCfg      = "/cfg"

	urlLogin  = "/login"
	urlNew    = "/new"
	urlEdit   = "/edit"
	urlDelete = "/delete"
	urlGet    = "/get"

	UrlApiLogin = urlPrefixVersion1 + urlApiPrefix + urlLogin
	UrlLogin    = urlLogin

	UrlNewService    = urlPrefixVersion1 + urlApiPrefix + urlPrefixService + urlNew
	UrlEditService   = urlPrefixVersion1 + urlApiPrefix + urlPrefixService + urlEdit
	UrlDeleteService = urlPrefixVersion1 + urlApiPrefix + urlPrefixService + urlDelete

	UrlNewUser    = urlPrefixVersion1 + urlApiPrefix + urlPrefixUser + urlNew
	UrlEditUser   = urlPrefixVersion1 + urlApiPrefix + urlPrefixUser + urlEdit
	UrlDeleteUser = urlPrefixVersion1 + urlApiPrefix + urlPrefixUser + urlDelete

	UrlNewCfg    = urlPrefixVersion1 + urlApiPrefix + urlPrefixCfg + urlNew
	UrlEditCfg   = urlPrefixVersion1 + urlApiPrefix + urlPrefixCfg + urlEdit
	UrlDeleteCfg = urlPrefixVersion1 + urlApiPrefix + urlPrefixCfg + urlDelete
)
