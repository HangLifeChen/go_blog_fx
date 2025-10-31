package config

import (
	"fmt"
	"io"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server             ServerConfig       `json:"server" yaml:"server"`                             // server config
	Cipher             Cipher             `json:"cipher" yaml:"cipher"`                             // cipher config
	I18n               I18n               `json:"i18n" yaml:"i18n"`                                 // i18n config
	Database           DatabaseConfig     `json:"database" yaml:"database"`                         // database config
	Jwt                Jwt                `json:"jwt" yaml:"jwt"`                                   // jwt config
	ThirdPartyServices ThirdPartyServices `json:"third_party_services" yaml:"third_party_services"` // third party services config
}

type ServerConfig struct {
	Domain           string          `json:"domain" yaml:"domain"`                       // domain
	Port             uint            `json:"port" yaml:"port"`                           // web server port
	ReadTimeout      time.Duration   `json:"read_timeout" yaml:"read_timeout"`           // read timeout
	WriteTimeout     time.Duration   `json:"write_timeout" yaml:"write_timeout"`         // write timeout
	GracefulShutdown time.Duration   `json:"graceful_shutdown" yaml:"graceful_shutdown"` // graceful shutdown time
	Mode             string          `json:"mode" yaml:"mode"`                           // environment(local/dev/pre/prod)
	Whitelist        map[string]bool `json:"whitelist" yaml:"whitelist"`                 // whitelist
}

type Cipher struct {
	AesGcm AesGcm `json:"aes_gcm" yaml:"aes_gcm"`
}

type AesGcm struct {
	Key string `json:"key" yaml:"key"`
}

type I18n struct {
	LanguageFilePath   string   `json:"language_file_path" yaml:"language_file_path"`     // language file path
	LanguageFileFormat string   `json:"language_file_format" yaml:"language_file_format"` // language file format
	DefaultLanguage    string   `json:"default_language" yaml:"default_language"`         // default language
	AcceptLanguages    []string `json:"accept_languages" yaml:"accept_languages"`         // accept languages
}

type DatabaseConfig struct {
	Db    DbSplittingConfig `json:"db" yaml:"db"`       // db config
	Redis RedisConfig       `json:"redis" yaml:"redis"` // redis config
	ES    ES                `json:"es" yaml:"es"`
}

type DbSplittingConfig struct {
	IsCluster bool       `json:"is_cluster" yaml:"is_cluster"` // is cluster
	Source    []DbConfig `json:"source" yaml:"source"`         // source
	Replica   []DbConfig `json:"replica" yaml:"replica"`       // replica
}

// Mysql 数据库配置
type DbConfig struct {
	Host         string `json:"host" yaml:"host"`                     // 数据库服务器的地址
	Port         int    `json:"port" yaml:"port"`                     // 数据库服务器的端口号
	Config       string `json:"config" yaml:"config"`                 // 数据库连接的配置参数，如驱动、字符集等
	DBName       string `json:"db_name" yaml:"db_name"`               // 要连接的数据库名称
	Username     string `json:"username" yaml:"username"`             // 用于连接数据库的用户名
	Password     string `json:"password" yaml:"password"`             // 用于连接数据库的密码
	MaxIdleConns int    `json:"max_idle_conns" yaml:"max_idle_conns"` // 最大空闲连接数，控制连接池中的空闲连接数量
	MaxOpenConns int    `json:"max_open_conns" yaml:"max_open_conns"` // 最大打开连接数，限制同时打开的数据库连接数量
	LogMode      string `json:"log_mode" yaml:"log_mode"`             // 日志模式，例如 "info" 或 "silent"，用于控制日志输出
}

// type DbConfig struct {
// 	Host     string `json:"host" yaml:"host"`         // host
// 	Port     int    `json:"port" yaml:"port"`         // port
// 	DbName   string `json:"db_name" yaml:"db_name"`   // database name
// 	Username string `json:"username" yaml:"username"` // username
// 	Password string `json:"password" yaml:"password"` // password
// }

type RedisConfig struct {
	Host      string `json:"host" yaml:"host"`             // host
	Port      int    `json:"port" yaml:"port"`             // port
	Password  string `json:"password" yaml:"password"`     // password
	Db        int    `json:"db" yaml:"db"`                 // db index
	EnableTLS bool   `json:"enable_tls" yaml:"enable_tls"` // enable tls or not
}

// ES ElasticSearch 配置
type ES struct {
	URL            string `json:"url" yaml:"url"`                           // Elasticsearch 服务的 URL，例如 http://localhost:9200
	Username       string `json:"username" yaml:"username"`                 // 用于连接 Elasticsearch 的用户名
	Password       string `json:"password" yaml:"password"`                 // 用于连接 Elasticsearch 的密码
	IsConsolePrint bool   `json:"is_console_print" yaml:"is_console_print"` // 是否在控制台打印 Elasticsearch 语句，true 表示打印，false 表示不打印
}

// Jwt jwt 配置
type Jwt struct {
	AccessTokenSecret      string `json:"access_token_secret" yaml:"access_token_secret"`             // 用于生成和验证访问令牌的密钥
	RefreshTokenSecret     string `json:"refresh_token_secret" yaml:"refresh_token_secret"`           // 用于生成和验证刷新令牌的密钥
	AccessTokenExpiryTime  string `json:"access_token_expiry_time" yaml:"access_token_expiry_time"`   // 访问令牌的过期时间，例如 "15m" 表示 15 分钟
	RefreshTokenExpiryTime string `json:"refresh_token_expiry_time" yaml:"refresh_token_expiry_time"` // 刷新令牌的过期时间，例如 "30d" 表示 30 天
	Issuer                 string `json:"issuer" yaml:"issuer"`                                       // JWT 的签发者信息，通常是应用或服务的名称
}

type ThirdPartyServices struct {
	Captcha Captcha `json:"captcha" yaml:"captcha"`
	Email   Email   `json:"email" yaml:"email"`
	Gaode   Gaode   `json:"gaode" yaml:"gaode"`
	Qiniu   Qiniu   `json:"qiniu" yaml:"qiniu"`
	QQ      QQ      `json:"qq" yaml:"qq"`
	System  System  `json:"system" yaml:"system"`
	Upload  Upload  `json:"upload" yaml:"upload"`
	Website Website `json:"website" yaml:"website"`
	Zap     Zap     `json:"zap" yaml:"zap"`
	Aliyun  Aliyun  `json:"aliyun" yaml:"aliyun"`
}

// 阿里云 OSS 配置
type Aliyun struct {
	Region        string `json:"region" yaml:"region"`                   // 存储区域
	Bucket        string `json:"bucket" yaml:"bucket"`                   // 空间名称
	ImgPath       string `json:"img_path" yaml:"img_path"`               // CDN 加速域名
	AccessKey     string `json:"access_key" yaml:"access_key"`           // 秘钥 AK
	SecretKey     string `json:"secret_key" yaml:"secret_key"`           // 秘钥 SK
	UseHTTPS      bool   `json:"use_https" yaml:"use_https"`             // 是否使用 https
	UseCdnDomains bool   `json:"use_cdn_domains" yaml:"use_cdn_domains"` // 上传是否使用 CDN 上传加速
}

// Captcha 验证码配置 mapstructure 定义数据序列化为map后key的值
type Captcha struct {
	Length   int     `mapstructure:"length" json:"length" yaml:"length"`
	Width    int     `mapstructure:"width" json:"width" yaml:"width"`
	Height   int     `mapstructure:"height" json:"height" yaml:"height"`
	MaxSkew  float64 `mapstructure:"max_skew" json:"max_skew" yaml:"max_skew"`
	DotCount int     `mapstructure:"dot_count" json:"dot_count" yaml:"dot_count"`
}

// Email 邮箱配置，可以登录 QQ 邮箱，打开设置，开启第三方服务服务，详情请见 https://mail.qq.com/
type Email struct {
	Host     string `json:"host" yaml:"host"`         // 邮件服务器地址，例如 smtp.qq.com
	Port     int    `json:"port" yaml:"port"`         // 邮件服务器端口，常见的如 587 (TLS) 或 465 (SSL)
	From     string `json:"from" yaml:"from"`         // 发件人邮箱地址
	Nickname string `json:"nickname" yaml:"nickname"` // 发件人昵称，用于显示在邮件中的发件人信息
	Secret   string `json:"secret" yaml:"secret"`     // 发件人邮箱的密码或应用专用密码，用于身份验证
	IsSSL    bool   `json:"is_ssl" yaml:"is_ssl"`     // 是否使用 SSL 加密连接，true 表示使用，false 表示不使用
}

// Gaode 高德服务配置，详情请见 https://lbs.amap.com/
type Gaode struct {
	Enable bool   `json:"enable" yaml:"enable"` // 是否开启高德服务，true 表示启用，false 表示禁用
	Key    string `json:"key" yaml:"key"`       // 高德服务的应用密钥，用于身份验证和服务访问
}

// Qiniu 七牛云配置，详情请见 https://www.qiniu.com/
type Qiniu struct {
	Zone          string `json:"zone" yaml:"zone"`                       // 存储区域
	Bucket        string `json:"bucket" yaml:"bucket"`                   // 空间名称
	ImgPath       string `json:"img_path" yaml:"img_path"`               // CDN 加速域名
	AccessKey     string `json:"access_key" yaml:"access_key"`           // 秘钥 AK
	SecretKey     string `json:"secret_key" yaml:"secret_key"`           // 秘钥 SK
	UseHTTPS      bool   `json:"use_https" yaml:"use_https"`             // 是否使用 https
	UseCdnDomains bool   `json:"use_cdn_domains" yaml:"use_cdn_domains"` // 上传是否使用 CDN 上传加速
}

// QQ qq 登录配置，详情请见 https://connect.qq.com/
type QQ struct {
	Enable      bool   `json:"enable" yaml:"enable"`             // 是否启用 qq 登录，true 表示启用，false 表示禁用
	AppID       string `json:"app_id" yaml:"app_id"`             // 应用 ID
	AppKey      string `json:"app_key" yaml:"app_key"`           // 应用密钥
	RedirectURI string `json:"redirect_uri" yaml:"redirect_uri"` // 网站回调域
}

func (qq QQ) QQLoginURL() string {
	return "https://graph.qq.com/oauth2.0/authorize?" +
		"response_type=code&" +
		"client_id=" + qq.AppID + "&" +
		"redirect_uri=" + qq.RedirectURI
}

// System 系统配置
type System struct {
	Host           string `json:"-" yaml:"host"`                          // 服务器绑定的主机地址，通常为 0.0.0.0 表示监听所有可用地址
	Port           int    `json:"-" yaml:"port"`                          // 服务器监听的端口号，通常用于 HTTP 服务
	Env            string `json:"-" yaml:"env"`                           // Gin 的环境类型，例如 "debug"、"release" 或 "test"
	RouterPrefix   string `json:"-" yaml:"router_prefix"`                 // API 路由前缀，用于构建 API 路径
	UseMultipoint  bool   `json:"use_multipoint" yaml:"use_multipoint"`   // 是否启用多点登录拦截，防止同一账户在多个地方同时登录
	SessionsSecret string `json:"sessions_secret" yaml:"sessions_secret"` // 用于加密会话的密钥，确保会话数据的安全性
	OssType        string `json:"oss_type" yaml:"oss_type"`               // 对应的对象存储服务类型，如 "local" 或 "qiniu"
}

func (s System) Addr() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

// Upload 图片上传配置
type Upload struct {
	Size int    `json:"size" yaml:"size"` // 图片上传的大小，单位 MB
	Path string `json:"path" yaml:"path"` // 图片上传的目录
}

// Website 网站信息
type Website struct {
	Logo                 string `json:"logo" yaml:"logo"`
	FullLogo             string `json:"full_logo" yaml:"full_logo"`
	Title                string `json:"title" yaml:"title"`                                   // 网站标题
	Slogan               string `json:"slogan" yaml:"slogan"`                                 // 网站标语
	SloganEn             string `json:"slogan_en" yaml:"slogan_en"`                           // 英文标语
	Description          string `json:"description" yaml:"description"`                       // 网站描述
	Version              string `json:"version" yaml:"version"`                               // 网站版本
	CreatedAt            string `json:"created_at" yaml:"created_at"`                         // 创建时间
	IcpFiling            string `json:"icp_filing" yaml:"icp_filing"`                         // ICP 备案
	PublicSecurityFiling string `json:"public_security_filing" yaml:"public_security_filing"` // 公安备案
	BilibiliURL          string `json:"bilibili_url" yaml:"bilibili_url"`                     // Bilibili 链接
	GiteeURL             string `json:"gitee_url" yaml:"gitee_url"`                           // Gitee 链接
	GithubURL            string `json:"github_url" yaml:"github_url"`                         // GitHub 链接
	Name                 string `json:"name" yaml:"name"`                                     // 昵称
	Job                  string `json:"job" yaml:"job"`                                       // 职业
	Address              string `json:"address" yaml:"address"`                               // 地址
	Email                string `json:"email" yaml:"email"`                                   // 邮箱
	QQImage              string `json:"qq_image" yaml:"qq_image"`                             // QQ 图片链接
	WechatImage          string `json:"wechat_image" yaml:"wechat_image"`                     // 微信图片链接
}

// Zap 日志配置，详情可参考七米的博客 https://liwenzhou.com/posts/Go/zap/
type Zap struct {
	Level          string `json:"level" yaml:"level"`                       // 日志等级，无特殊需求，用 info 即可
	Filename       string `json:"filename" yaml:"filename"`                 // 日志文件的位置
	MaxSize        int    `json:"max_size" yaml:"max_size"`                 // 在进行切割之前，日志文件的最大大小（以MB为单位）
	MaxBackups     int    `json:"max_backups" yaml:"max_backups"`           // 保留旧文件的最大个数
	MaxAge         int    `json:"max_age" yaml:"max_age"`                   // 保留旧文件的最大天数
	IsConsolePrint bool   `json:"is_console_print" yaml:"is_console_print"` // 是否在控制台打印日志，true 表示打印，false 表示不打印
}

// Redis 缓存数据库配置
type Redis struct {
	Address  string `json:"address" yaml:"address"`   // Redis 服务器的地址，通常为 "localhost:6379" 或其他主机和端口
	Password string `json:"password" yaml:"password"` // 连接 Redis 时的密码，如果没有设置密码则留空
	DB       int    `json:"db" yaml:"db"`             // 指定使用的数据库索引，单实例模式下可选择的数据库，默认为 0
}

// read and parse the configuration file
func LoadConfig(c string) (cfg *Config, err error) {
	var file *os.File
	if c != "" {
		file, err = os.Open(c)
	} else {
		file, err = os.Open("../../conf/config.yaml")
	}

	if err != nil {
		return
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		return
	}
	err = yaml.Unmarshal(content, &cfg)
	if err != nil {
		return
	}
	return
}
