package errcode

import (
	"fmt"
	"strings"
	"zinx/lib/logger"
)

// Error 错误
type Error struct {
	Code     int      `json:"code"`    // 错误码
	Msg      string   `json:"msg"`     // 消息
	Details  []string `json:"Details"` // 错误详情
	isFormat bool     // 是否是格式化的错误
}

// codes 存储错误码与错误信息的映射关系
var codes = map[int]string{}

// NewError 新建错误
func NewError(code int, msg string) *Error {
	if _, ok := codes[code]; ok {
		panic(fmt.Sprintf("错误信息 ID: %d 已经存在，请更换一个", code))
	}
	codes[code] = msg
	return &Error{
		Code: code,
		Msg:  msg,
	}
}

// NewFormatError 新建格式化错误
func NewFormatError(code int, format string) (err *Error) {
	err = NewError(code, format)
	err.isFormat = true
	return
}

// NewErrorWithDetails 新建错误, 并添加错误详情
func NewErrorWithDetails(code int, msg string, details ...string) (err *Error) {
	err = NewError(code, msg)
	err.AddDetails(details...)
	return
}

func (e *Error) Error() string {
	if strings.ContainsRune(e.Msg, '%') {
		logger.Warnf("错误信息 ID: %d 中包含格式化字符，建议使用 Msgf 方法格式化错误信息", e.Code)
	}
	return fmt.Sprintf("ErrorCode：%d, Message：%s", e.Code, e.Msg)
}

// Msgf 格式化错误信息，类似于 'fmt.Sprintf
//
// 如果错误信息不是格式化的，则直接返回错误信息
func (e *Error) Msgf(args ...any) string {
	if !e.isFormat {
		logger.Warn("错误信息 ID: %d 不是格式化的，建议使用 Error 方法获取错误信息", e.Code)
		return e.Msg
	}
	return fmt.Sprintf(e.Msg, args)
}

// Format 在原错误上使用 Msgf 替换原错误信息
//
// 返回原错误的副本
func (e *Error) Format(args ...any) *Error {
	err := *e
	err.Msg = e.Msgf(args...)
	return &err
}

// AddDetails 添加错误详情
func (e *Error) AddDetails(details ...string) {
	for _, d := range details {
		e.Details = append(e.Details, d)
	}
}

// WithDetails 创建一个当前错误的副本，并重置错误详情
func (e *Error) WithDetails(details ...string) *Error {
	err := *e
	err.Details = []string{}
	for _, d := range details {
		err.Details = append(err.Details, d)
	}
	return &err
}
