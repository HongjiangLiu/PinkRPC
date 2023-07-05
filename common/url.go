package common

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"pinkrpc/common/constant"
	perrors "pinkrpc/common/errors"
	"strings"
)

const (
	CONSUMER = iota
	CONFIGURATOR
	ROUTER
	PROVIDER
)

var PinkNodes = [...]string{"consumers", "configurators", "routers", "providers"}
var PinkRole = [...]string{"consumer", "", "", "provider"}

type RoleType int

// 获取PinkNodes和PinkRole的string类型值
func (t RoleType) String() string {
	return PinkNodes[t]
}

func (t RoleType) Role() string {
	return PinkRole[t]
}

type baseUrl struct {
	Protocol     string
	Location     string
	Ip           string
	Port         string
	Params       url.Values //用于处理URL查询参数的类型，是一个map类型，用于存储URL中查询参数的键值对
	PrimitiveURL string     //原始的Url
	ctx          context.Context
}

type URL struct {
	baseUrl
	Path     string //例如/com.ikurento.dubbo.UserProvider3
	Username string
	Password string
	Methods  []string
	//为注册中心特殊指定
	SubURL *URL
}

type option func(*URL) //option是一个函数类型，接受一个URL指针作为参数 没有返回值

// WithUsername 选项函数，用于设置URL结构体中Username字段
func WithUsername(username string) option {
	//它返回一个匿名函数，这个匿名函数即是option类型的实例。
	//这个匿名函数接收一个指向URL结构体的指针作为参数，并将传入的username值设置给URL结构体的Username字段。
	return func(url *URL) {
		url.Username = username
	}
	//这种设计模式允许我们在调用函数时，通过调用多个选项函数来设置不同的参数，从而配置函数的行为。
}

func WithPassword(password string) option {
	return func(url *URL) {
		url.Password = password
	}
}

func WithMethods(methods []string) option {
	return func(url *URL) {
		url.Methods = methods
	}
}

func WithParams(params url.Values) option {
	return func(url *URL) {
		url.Params = params
	}
}

func WithParamsValue(key, val string) option {
	return func(url *URL) {
		url.Params.Set(key, val)
	}
}

func WithProtocol(protocol string) option {
	return func(url *URL) {
		url.Protocol = protocol
	}
}
func WithIp(ip string) option {
	return func(url *URL) {
		url.Ip = ip
	}
}

func WithPort(port string) option {
	return func(url *URL) {
		url.Port = port
	}
}

func NewURLWithOptions(service string, opts ...option) *URL {
	//service用来指定URL中的服务名，option类型作为一个函数类型，接受一个*URL类型的指针作为参数，用于设置URL对象的属性
	url := &URL{
		Path: "/" + service,
	}
	//在每次迭代中，将当前选项函数opt应用于url对象。
	//选项函数opt会修改url对象的属性，根据选项函数的具体实现来设置URL的其他属性。
	for _, opt := range opts {
		opt(url)
	}
	url.Location = url.Ip + ":" + url.Port
	return url
}

func NewURL(ctx context.Context, urlString string, opts ...option) (URL, error) {
	var (
		err          error
		rawUrlString string
		serviceUrl   *url.URL
		//使用ctx构造一个默认的URL
		s = URL{baseUrl: baseUrl{ctx: ctx}}
	)
	//new a null instance
	if urlString == "" {
		return s, nil
	}

	/*
		在URL中，一些字符（例如空格、特殊字符等）需要进行编码，以便在URL中正确传递和解析。
		URL编码使用百分号（%）后跟两个十六进制数字来表示特殊字符的编码值。
		url.QueryUnescape()函数的作用是将URL编码字符串解码为原始的非编码形式。
		这样做是为了还原URL中的特殊字符和空格等，使其成为原始的可读字符串。
	*/
	rawUrlString, err = url.QueryUnescape(urlString)
	if err != nil {
		return s, perrors.Errorf("url.QueryUnescape(%s),  error{%v}", urlString, err)
	}

	/*
		用于解析一个URL字符串并返回一个*url.URL对象。
		它能够将URL字符串解析为更易于处理的结构体形式，方便我们提取URL的各个部分，如协议、主机名、路径、查询参数等
	*/
	serviceUrl, err = url.Parse(rawUrlString)
	if err != nil {
		return s, perrors.Errorf("url.Parse(url string{%s}),  error{%v}", rawUrlString, err)
	}

	/*
		RawQuery是url.URL结构体中的一个字段，用于存储URL中的原始查询字符串（不包括问号）

		url.ParseQuery(serviceUrl.RawQuery)：这是net/url包中的函数ParseQuery()，用于解析URL查询字符串。
		它接收一个字符串作为参数，返回一个url.Values类型的map。url.Values是一个类似map的结构，用于存储URL查询参数的键值对。
	*/
	s.Params, err = url.ParseQuery(serviceUrl.RawQuery)
	if err != nil {
		return s, perrors.Errorf("url.ParseQuery(raw url string{%s}),  error{%v}", serviceUrl.RawQuery, err)
	}

	s.PrimitiveURL = urlString
	s.Protocol = serviceUrl.Scheme
	s.Username = serviceUrl.User.Username()
	s.Password, _ = serviceUrl.User.Password()
	s.Location = serviceUrl.Host
	s.Path = serviceUrl.Path
	//如果Location中存在分号，那么将它分开之后复制给s.IP和s.Port
	if strings.Contains(s.Location, ":") {
		s.Ip, s.Port, err = net.SplitHostPort(s.Location)
		if err != nil {
			return s, perrors.Errorf("net.SplitHostPort(Url.Host{%s}), error{%v}", s.Location, err)
		}
	}

	for _, opt := range opts {
		opt(&s)
	}

	return s, nil
}

func (c URL) String() string {
	buildString := fmt.Sprintf(
		"%s://%s:%s@%s:%s%s?",
		c.Protocol, c.Username, c.Password, c.Ip, c.Port, c.Path)
	//Params是个K-V类型
	//Encode() 方法会将 url.Values 中的键值对按照 URL 编码规则进行编码，并返回一个字符串，该字符串就是 URL 查询字符串的形式
	buildString += c.Params.Encode()
	return buildString
}

/*
URL 是一个值类型，函数的接收者是 URL 的一个副本（pass-by-value）。
当你调用 URLEqual() 函数时，会将原始 URL 对象复制一份，然后在函数内部使用这个副本进行操作。
因此，对副本的任何修改都不会影响到原始 URL 对象。

*URL 是一个指针类型，函数的接收者是 URL 的指针（pass-by-reference）。
当你调用 URLEqual() 函数时，实际上是将原始 URL 对象的地址传递给函数。
在函数内部，操作的是原始 URL 对象本身，而不是副本。因此，对原始 URL 对象的修改会直接反映在函数外部的对象上。

这里需要注意的是，如果在函数内部不需要修改 URL 对象的状态，而只是进行一些读取操作，那么使用 URL 值类型就足够了。
而如果需要修改 URL 对象的状态，那么最好使用指针类型 *URL，这样可以确保在函数内部对原始对象的修改能够影响到函数外部。
*/
func (c URL) URLEqual(url URL) bool {
	c.Ip = ""
	c.Port = ""
	url.Ip = ""
	url.Port = ""
	if c.Key() != url.Key() {
		return false
	}
	return true
}

func (c URL) Key() string {
	buildString := fmt.Sprintf(
		"%s://%s:%s@%s:%s/%s?group=%s&version=%s",
		c.Protocol, c.Username, c.Password, c.Ip, c.Port, c.GetParam(constant.INTERFACE_KEY, strings.TrimPrefix(c.Path, "/")), c.GetParam(constant.GROUP_KEY, ""), c.GetParam(constant.VERSION_KEY, constant.DEFAULT_VERSION))

	return buildString
}

// GetParam 如果获得的参数为空，则用d(efault)来填充该值
func (c URL) GetParam(s string, d string) string {
	var r string
	if r = c.Params.Get(s); r == "" {
		r = d
	}
	return r
}
