package main
import (
"fmt"
)


func generate(ch chan int) {
     for i := 2; i<100; i++ {
           ch <- i  // Send 'i' to channel 'ch'.
     }
}

// Copy the values from channel 'in' to channel 'out',
// removing those divisible by 'prime'.
func filter(in, out chan int, prime int) {
  for {
      i := <-in  // Receive value of new variable 'i' from 'in'.
      //fmt.Println(prime)
      if i % prime != 0 {
          out <- i  // Send 'i' to channel 'out'.
      }
  }
}


func main() {
     ch := make(chan int)  // Create a new channel.
     go generate(ch)  // Start generate() as a goroutine.
       for {
           prime := <-ch
           fmt.Println(prime)
           ch1 := make(chan int)
           go filter(ch, ch1, prime)
           ch = ch1
        }
    }

/*

程序说明，虽然会报错，但是不影响功能	
整个程序可以理解为：


generate -> filter1(in,out,2) -> filter2(in,out,3) -> filter3(in,out,5) -> ...
generate负责生产素数，第一次生成2阻塞掉循环
filter1收到in：generate的ch，收到out：传入给prime变量的ch，以及素数实参2
filter1开始执行前，main中的变量prime开始接受filter1的输出
开始执行filter1，in ch收到generate函数传递过来的3，发现不能被2整除，于是 out <- 3

这时main函数的prime收到3，打印3，同时新建一个out ch再启动一个go的线程filter2
在执行filter2之前prime等待filter2的out输出

同时，filter1，filter2，main以及generate是同时执行的
filter2开始执行，阻塞循环，传入的in ch 表示filter1的out，out ch这时是main的prime在等待消费，执行等待filter的out过来的in ch ...
因为filter1消费了一次in，所以generate又生产了4，这时filter1因为不满足 4%2 != 0 这个条件，继续循环，消费generate
filter1继续消费到5，因为 5%2 != 0 成立，这是filter 执行out ch，生产数据5到filter2
filter2收到5，执行 5%3, 5%3 !=0 成立，所以 filter2 输出out到main的prime，打印素数5

接下来就继续执行了

*/