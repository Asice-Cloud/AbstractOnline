package log

/*
 	variable.ZapLog.Info("基本的运行提示类信息")
    variable.ZapLog.Warn("UserCreate接口参数非法警告，相关参数：",zap.String("userName","demo_name"),zap.Int("userAge",18))
    variable.ZapLog.Panic("UserDestory接口参数异常，相关参数：",zap.String("userName","demo_name"),zap.String("password","pass123456")
    variable.ZapLog.Error("UserDestory接口参数错误，相关参数：",zap.Error(error))
    variable.ZapLog.Fatal("Mysql初始化参数错误，退出运行。相关参数：",zap.String("name","root"), zap.Int("端口",3306))
*/

// explict application should use log hook under logger folder
