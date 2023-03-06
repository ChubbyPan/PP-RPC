package codec

/***
Gob消息结构相关定义
***/

import (
	"bufio"
	"encoding/gob"
	"io"
	"log"
)

type GobCodec struct {
	conn io.ReadWriteCloser // 由构造函数传入的连接符
	buf  *bufio.Writer      // 防止阻塞
	dec  *gob.Decoder       // 官方库的gob解码器
	enc  *gob.Encoder       // 官方库的gob编码器
}

var _ Codec = (*GobCodec)(nil) //TODO 干嘛用的

func NewGobCodec(conn io.ReadWriteCloser) Codec {
	buf := bufio.NewWriter(conn)
	return &GobCodec{
		conn: conn,
		buf:  buf,
		dec:  gob.NewDecoder(conn), //TODO 确定参数是conn 不是buf么
		enc:  gob.NewEncoder(buf),
	}
}

// 实现GobCodec的四种方法， 具体哪四种看Codec接口定义，这里是工厂模式
func (c *GobCodec) ReadHeader(h *Header) error {
	return c.dec.Decode(h)
}

func (c *GobCodec) ReadBody(body interface{}) error {
	return c.dec.Decode(body)
}

func (c *GobCodec) Close() error {
	return c.conn.Close()
}

func (c *GobCodec) Write(h *Header, body interface{}) (err error) {
	defer func() {
		_ = c.buf.Flush()
		if err != nil {
			_ = c.Close()
		}
	}()

	if err := c.enc.Encode(h); err != nil {
		log.Println("rpc codec: gob error encoding header:", err)
		return err
	}
	if err := c.enc.Encode(body); err != nil {
		log.Println("rpc codec: gob error encoding body:", err)
		return err
	}
	return nil
}
