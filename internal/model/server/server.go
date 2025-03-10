package server

type Server struct {

	/**
	 * CPU相关信息
	 */
	CPU CPU `json:"cpu"`

	/**
	 * 內存相关信息
	 */
	Mem Mem `json:"mem"`

	/**
	 * JVM相关信息
	 */
	Jvm JVM `json:"jvm"`

	/**
	 * 服务器相关信息
	 */
	Sys Sys `json:"system"`
	/**
	 * 磁盘相关信息
	 */
	SysFiles []SysFile `json:"sysFiles"`
}

type CPU struct {

	/**
	 * 核心数
	 */
	CpuNum int `json:"cpuNum"`

	/**
	 * CPU总的使用率
	 */
	Total float64 `json:"total"`

	/**
	 * CPU系统使用率
	 */
	Sys float64 `json:"system"`

	/**
	 * CPU用户使用率
	 */
	Used float64 `json:"used"`

	/**
	 * CPU当前等待率
	 */
	Wait float64 `json:"wait"`

	/**
	 * CPU当前空闲率
	 */
	Free float64 `json:"free"`
}

type Mem struct {

	/**
	 * 内存总量
	 */
	Total float64 `json:"total"`

	/**
	 * 已用内存
	 */
	Used float64 `json:"used"`

	/**
	 * 剩余内存
	 */
	Free float64 `json:"free"`
}

type JVM struct {

	/**
	 * 当前JVM占用的内存总数(M)
	 */
	Total float64 `json:"total"`

	/**
	 * JVM最大可用内存总数(M)
	 */
	Max float64 `json:"max"`

	/**
	 * JVM空闲内存(M)
	 */
	Free float64 `json:"free"`

	/**
	 * JDK版本
	 */
	Version string `json:"version"`

	/**
	 * JDK路径
	 */
	Home string `json:"home"`
}

type Sys struct {

	/**
	 * 服务器名称
	 */
	ComputerName string `json:"computerName"`

	/**
	 * 服务器Ip
	 */
	ComputerIp string `json:"computerIp"`

	/**
	 * 项目路径
	 */
	UserDir string `json:"userDir"`

	/**
	 * 操作系统
	 */
	OsName string `json:"osName"`

	/**
	 * 系统架构
	 */
	OsArch string `json:"osArch"`
}

type SysFile struct {

	/**
	 * 盘符路径
	 */
	DirName string `json:"dirName"`

	/**
	 * 盘符类型
	 */
	SysTypeName string `json:"sysTypeName"`

	/**
	 * 文件类型
	 */
	TypeName string `json:"typeName"`

	/**
	 * 总大小
	 */
	Total string `json:"total"`

	/**
	 * 剩余大小
	 */
	Free string `json:"free"`

	/**
	 * 已经使用量
	 */
	Used string `json:"used"`

	/**
	 * 资源的使用率
	 */
	Usage float64 `json:"usage"`
}
