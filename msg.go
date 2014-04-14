// 读取连接，解析出消息头，返回消息正文和消息类型
package msg

import (
	"bytes"
	"encoding/binary"
	"errors"
)

const (
	MAX_BUFFER   = 1024 // 读取缓存最大值
	SIZE_OF_TYPE = 4    // sizeof int32
	SIZE_OF_SIZE = 4    // sizeof int32
	SIZE_OF_HEAD = SIZE_OF_TYPE + SIZE_OF_SIZE
)

// 消息的结构
type Msg struct {
	Type    int32  // 消息类型
	Size    int32  // 消息大小（含消息类型和消息大小自身
	Content []byte // 消息正文
}

// 解包消息，返回Msg类型
func UnPack(b []byte) (Msg, error) {
	m := Msg{}
	buf := bytes.NewBuffer(b)
	// 消息类型
	mType := buf.Next(SIZE_OF_TYPE)
	bufType := bytes.NewBuffer(mType)
	binary.Read(bufType, binary.LittleEndian, &m.Type)
	// 消息大小
	mSize := buf.Next(SIZE_OF_SIZE)
	bufSize := bytes.NewBuffer(mSize)
	binary.Read(bufSize, binary.LittleEndian, &m.Size)
	// 超限则返回错误
	if m.Size > MAX_BUFFER {
		return m, errors.New("OVER_MAX_BUFFER")
	}
	// 消息正文
	mContent := buf.Bytes()
	rest := int(m.Size - int32(SIZE_OF_HEAD))
	if rest > 0 {
		m.Content = mContent[:rest]
	}
	return m, nil
}

// 打包消息，返回[]byte
func Pack(mType int32, mContent []byte) []byte {
	buf := new(bytes.Buffer)
	// 消息类型
	binary.Write(buf, binary.LittleEndian, mType)
	// 消息大小
	mSize := int32(SIZE_OF_HEAD + len(mContent))
	binary.Write(buf, binary.LittleEndian, mSize)
	// 消息正文
	binary.Write(buf, binary.LittleEndian, mContent)
	b := buf.Bytes()
	return b
}
