package dlog

//func main() {
//	// 初始化 日志服务
//	topic := "go_passport"
//	dlog.SetTopic(topic) // dlog.SetJsonLogTopic(topic) 这种方式是使用json 日志

//	dlog.DebugLog(true)
//
//	host := ""
//	dlog.Info("server start", topic, "host", host, "addr", constants.HttpListenPort)
//
//	...
//	c := make(chan os.Signal, 1)
//	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
//	s := <-c
//	dlog.Info("signal.Notify", s)
//	dlog.Info("wait work ", "done")
//	dlog.Close()
//	fmt.Println("dlog flush")
//	fmt.Println("server exit")
//}
