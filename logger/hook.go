package logger

import (
	"go.uber.org/zap/zapcore"
)

func ZapLogHandler(entry zapcore.Entry) error {
	// 参数 entry 介绍
	// entry  参数就是单条日志结构体，主要包括字段如下：
	//Level      日志等级
	//Time       当前时间
	//LoggerName  日志名称
	//Message    日志内容
	//Caller     各个文件调用路径
	//Stack      代码调用栈

	//这里启动一个协程，hook丝毫不会影响程序性能，
	go func(paramEntry zapcore.Entry) {
		//fmt.Println(" Zap log hook ....,you can dispose log here....")
		//fmt.Printf("%#+v\n", paramEntry)
	}(entry)
	return nil
}
