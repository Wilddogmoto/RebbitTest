package logmq

var (
	l = InitializationLogger()
)
func ResponseError(err error,msg string){
	if err!=nil{
		l.Errorf("%v:%v",msg,err)
	}
}

