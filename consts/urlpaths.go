package consts

const (
	urlPrefixVersion = "/v1"
	urlPrefixService = "/service"
	urlPrefixUser    = "/user"
	urlPrefixCfg     = "/cfg"

	urlNew    = "/new"
	urlEdit   = "/edit"
	urlDelete = "/delete"
	urlGet    = "/get"

	UrlNewService    = urlPrefixVersion + urlPrefixService + urlNew
	UrlEditService   = urlPrefixVersion + urlPrefixService + urlEdit
	UrlDeleteService = urlPrefixVersion + urlPrefixService + urlDelete

	UrlNewUser    = urlPrefixVersion + urlPrefixUser + urlNew
	UrlEditUser   = urlPrefixVersion + urlPrefixUser + urlEdit
	UrlDeleteUser = urlPrefixVersion + urlPrefixUser + urlDelete

	UrlNewCfg    = urlPrefixVersion + urlPrefixCfg + urlNew
	UrlEditCfg   = urlPrefixVersion + urlPrefixCfg + urlEdit
	UrlDeleteCfg = urlPrefixVersion + urlPrefixCfg + urlDelete
)
