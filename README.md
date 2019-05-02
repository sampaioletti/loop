Package Loop allows for a method to schedule loop calls on certain threads.

Inspired by https://github.com/faiface/mainthread but wanted something more focused on GUI loop threads

See examples for usage but in general

```golang
func main(){
  //create the first loop, where its created doesn't matter
  l1 := loop.NewLoop()
  l1.AddCall(func() {
    fmt.Println("Hello From l1")
    time.Sleep(time.Second * 1)
  })
  go func() {
    l1.AddCall(func() {
      fmt.Println("Hello Again From l1")
      time.Sleep(time.Second * 5)
        })
        //l1 will run in goroutine thread
    l1.Run(nil)
  }()

  l2 := loop.NewLoop()
  l2.AddCall(func() {
    fmt.Println("Hello From l2")
    time.Sleep(time.Second * 1)
    })
    //l2 will run in the main loop
  fmt.Println(l2.Run(nil))
}
```
Output:
```
Hello From l1
Hello From l2
Hello Again From l1
Hello From l2
Hello From l2
Hello From l2
Hello From l2
Hello From l2
Hello From l1
Hello From l2
Hello Again From l1
```
