package main
import (
"fmt"
)


type request struct {  //创建一个request对象
        a, b    int
        replyc  chan int
}
type binOp func(a, b int) int    //定义一个binOp类型 表示函数 xxx 类型
func run(op binOp, req *request) {   //将函数op 和结构体 req传入run函数
     reply := op(req.a, req.b)
     req.replyc <- reply
}

func server(op binOp, service chan *request, quit chan bool) {
     for {
              select {
              case req := <-service:
                  go run(op, req)  // don't wait for it
              case <-quit:
                  return
              }
          }
}

func startServer(op binOp) (service chan *request, quit chan bool) {
        service = make(chan *request)
        quit = make(chan bool)
        go server(op, service, quit)
        return service, quit
}

func main() {
          adder, quit := startServer(func(a, b int) int { return a + b })
          const N = 100
          var reqs [N]request
          for i := 0; i < N; i++ {
              req := &reqs[i]
              req.a = i
              req.b = i + N
              req.replyc = make(chan int)
              adder <- req
          }
          for i := N-1; i >= 0; i-- {   // doesn't matter what order
              if <-reqs[i].replyc != N + 2*i {
                  fmt.Println("fail at", i)
              }
          }
          quit <- true
          fmt.Println("done")
}

/*
代码简要说明：
首先执行startServer方法，传入两数相加的函数
此时执行 server 函数，server函数里执行一个for循环，监听req chan 是否有消息发送过来，如果一旦收到消息，则执行 run 函数
startServer 返回一个 传递 request类型的管道
设置常量100
设置reqs数组，长度为100
循环reqs数组，对erqs数组的中的request对象进行初始化工作
利用make创建reqs数组中的每项request对象的replyc管道，然后讲此request对象传递到adder管道中
adder管道一旦收到request对象，就开始执行 run 函数
run函数先执行op(req.a, req.b),讲计算结果赋值给局部变量replay，然后讲replay通过req的管道reqplc发到main里
main循环接受run线程发过来的值，如果结果不是N+2*i不正确，就打印 fail at
执行结束打印 done
*/