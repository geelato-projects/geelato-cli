package errors

type Code uint32

const (
	Success Code = 0

	ErrUnknown Code = 1000000

	ErrConfig Code = 1001000
	ErrConfigNotFound Code = 1001001
	ErrConfigInvalid Code = 1001002
	ErrConfigParse Code = 1001003

	ErrFile Code = 1002000
	ErrFileNotFound Code = 1002001
	ErrFilePermission Code = 1002002
	ErrFileRead Code = 1002003
	ErrFileWrite Code = 1002004

	ErrNetwork Code = 1003000
	ErrNetworkTimeout Code = 1003001
	ErrNetworkUnavailable Code = 1003002
	ErrNetworkResponse Code = 1003003

	ErrValidation Code = 1004000
	ErrValidationRequired Code = 1004001
	ErrValidationInvalid Code = 1004002
	ErrValidationNotMatch Code = 1004003

	ErrApp Code = 2001000
	ErrAppNotFound Code = 2001001
	ErrAppAlreadyExists Code = 2001002
	ErrAppInvalid Code = 2001003
	ErrAppInit Code = 2001004

	ErrModel Code = 2002000
	ErrModelNotFound Code = 2002001
	ErrModelAlreadyExists Code = 2002002
	ErrModelInvalid Code = 2002003
	ErrModelField Code = 2002004

	ErrAPI Code = 2003000
	ErrAPINotFound Code = 2003001
	ErrAPIAlreadyExists Code = 2003002
	ErrAPIInvalid Code = 2003003
	ErrAPIExecute Code = 2003004

	ErrGit Code = 2004000
	ErrGitNotRepo Code = 2004001
	ErrGitClone Code = 2004002
	ErrGitCommit Code = 2004003
	ErrGitPush Code = 2004004
	ErrGitPull Code = 2004005

	ErrSync Code = 2005000
	ErrSyncConflict Code = 2005001
	ErrSyncFailed Code = 2005002
	ErrSyncVersion Code = 2005003

	ErrPlatform Code = 3001000
	ErrPlatformAuth Code = 3001001
	ErrPlatformRequest Code = 3001002
	ErrPlatformResponse Code = 3001003
)

var messages = map[Code]string{
	Success: "成功",

	ErrUnknown: "未知错误",

	ErrConfig: "配置错误",
	ErrConfigNotFound: "配置文件不存在",
	ErrConfigInvalid: "配置文件格式无效",
	ErrConfigParse: "配置文件解析错误",

	ErrFile: "文件操作错误",
	ErrFileNotFound: "文件不存在",
	ErrFilePermission: "文件权限不足",
	ErrFileRead: "文件读取错误",
	ErrFileWrite: "文件写入错误",

	ErrNetwork: "网络错误",
	ErrNetworkTimeout: "网络请求超时",
	ErrNetworkUnavailable: "网络不可用",
	ErrNetworkResponse: "网络响应错误",

	ErrValidation: "数据验证错误",
	ErrValidationRequired: "必填字段为空",
	ErrValidationInvalid: "数据格式无效",
	ErrValidationNotMatch: "数据不匹配",

	ErrApp: "应用管理错误",
	ErrAppNotFound: "应用不存在",
	ErrAppAlreadyExists: "应用已存在",
	ErrAppInvalid: "应用配置无效",
	ErrAppInit: "应用初始化失败",

	ErrModel: "模型管理错误",
	ErrModelNotFound: "模型不存在",
	ErrModelAlreadyExists: "模型已存在",
	ErrModelInvalid: "模型定义无效",
	ErrModelField: "模型字段错误",

	ErrAPI: "API 管理错误",
	ErrAPINotFound: "API 不存在",
	ErrAPIAlreadyExists: "API 已存在",
	ErrAPIInvalid: "API 定义无效",
	ErrAPIExecute: "API 执行错误",

	ErrGit: "Git 操作错误",
	ErrGitNotRepo: "目录不是 Git 仓库",
	ErrGitClone: "Git 克隆失败",
	ErrGitCommit: "Git 提交失败",
	ErrGitPush: "Git 推送失败",
	ErrGitPull: "Git 拉取失败",

	ErrSync: "同步错误",
	ErrSyncConflict: "同步冲突",
	ErrSyncFailed: "同步失败",
	ErrSyncVersion: "版本不一致",

	ErrPlatform: "平台错误",
	ErrPlatformAuth: "平台认证失败",
	ErrPlatformRequest: "平台请求失败",
	ErrPlatformResponse: "平台响应错误",
}

func Message(code Code) string {
	if msg, ok := messages[code]; ok {
		return msg
	}
	return messages[ErrUnknown]
}
