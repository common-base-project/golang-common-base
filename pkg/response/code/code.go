package code

/*
  @Author : Mustang Kong
*/

var (
	// 通用
	InternalServerError = &Errno{Errno: 60001, Errmsg: "内部服务器错误"}
	ExistError          = &Errno{Errno: 160003, Errmsg: "数据已存在"}
	NotExistError       = &Errno{Errno: 10002, Errmsg: "数据不存在"}
	ParamError          = &Errno{Errno: 10001, Errmsg: "参数不正确"}
	BindError           = &Errno{Errno: 160006, Errmsg: "绑定失败"}

	// 成功
	Success = &Errno{Errno: 0, Errmsg: "请求成功"}

	// 未知失败
	UnknownError = &Errno{Errno: 199999, Errmsg: "未知错误"}

	// auth 1607xx
	CreateUserError = &Errno{Errno: 170701, Errmsg: "创建用户失败"}
	UpdateUserError = &Errno{Errno: 170702, Errmsg: "更新用户失败"}
	DeleteUserError = &Errno{Errno: 170703, Errmsg: "删除用户失败"}
	SelectUserError = &Errno{Errno: 170704, Errmsg: "查询用户失败"}

	// help
	SelectHelpContentError = &Errno{Errno: 161001, Errmsg: "获取帮助文档失败"}
	UpdateHelpContentError = &Errno{Errno: 161002, Errmsg: "更新帮助文档失败"}

	// namespace
	CreateCommonError = &Errno{Errno: 160001, Errmsg: "创建失败"}
	UpdateCommonError = &Errno{Errno: 160002, Errmsg: "更新失败"}
	DeleteCommonError = &Errno{Errno: 160003, Errmsg: "删除失败"}
	SelectCommonError = &Errno{Errno: 160005, Errmsg: "查询失败"}
)
