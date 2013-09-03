package main
import (
"fmt"
)


func generate() chan int {
    ch := make(chan int)
    go func(){
        for i := 2; i<100; i++ {
              ch <- i
          }
      }()
      return ch
 }


func filter(in chan int, prime int) chan int {
    out := make(chan int)
    go func() {
        for {
            if i := <-in; i % prime != 0 {
                out <- i
           }
        }
    }()
    return out
}


func sieve() chan int {
       out := make(chan int)
        go func() {
            ch := generate()
            for {
                prime := <-ch
                out <- prime
                ch = filter(ch, prime)
            }
        }()
        return out
 }


func main() {
 primes := sieve()
 for {
     fmt.Println(<-primes)
  }
}

/*
代码简要说明：
main函数入口，执行sieve函数，返回out ch
这样代码很清晰，main函数只是执行循环，等待primes的生产，然后消费不停的循环
sieve开始执行，定义 out ch，在sieve里执行一个线程闭包，获得generate的ch
执行循环，prime等待generate的生产
下面同prime前一个例子相同，只是filter的out步是filter里生成的，ch的赋值不再通过直接赋值ch1来实现，而是通过函数的返回值来实现

*/