package codec

import "io"

/***
	rpc的通信消息的序列化和反序列化
***/

// client 调用时，需要知道调用1.服务名 2.方法名 3.参数
// server 返回时需要有 1.error 2.reply

type Header struct {
	ServiceMethod string // 格式： 服务名.方法名
	Seq           uint64 // 发送序号，可选
	Error         string // 请求错误时返回
}

// 消息体编解码
type Codec interface {
	io.Closer
	ReadHeader(*Header) error
	ReadBody(interface{}) error
	Write(*Header, interface{}) error
}

type NewCodecFunc func(io.ReadWriteCloser) Codec // 构造消息体的构造函数

type Type string

const (
	GobType  Type = "application/gob"
	JsonType Type = "application/json" // 暂时没有实现
)

// 工厂模式，返回消息体的构造函数
var NewCodecFuncMap map[Type]NewCodecFunc

func init() {
	NewCodecFuncMap = make(map[Type]NewCodecFunc)
	NewCodecFuncMap[GobType] = NewGobCodec
}
