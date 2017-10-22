package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "time"
    "flag"
    "os/exec"
    "io/ioutil"
    "os"
)


//<task>-------------------------------------------------
type task struct {
    id string
    cmd string
    args []string
    interval int
    running bool
    stop chan int
}

//创建任务数据
func NewTask(id string, cmd string, args []string, interval int) task {
    var t task
    t.id = id
    t.cmd= cmd
    t.args=args
    t.interval = interval
    t.running = false
    t.stop = make(chan int)
    return t
}

//  开始任务
func (t *task) Start() {
    if t.running {
        return
    }
    t.running = true
    go t.run()
}

// 结束任务
func (t *task) Stop() {
    if !t.running {
        return
    }
    fmt.Println("t.stop <- 1")
    t.stop <- 1
    select {
        case <-t.stop:
            break
        default :
            break
    }
    t.running = false
    fmt.Println("t.stop")
}

// 执行定时任务
func (t *task) run() {
    fmt.Println("t.running-----------------------------------------------")
    for {
        select {
            case <-t.stop :
                fmt.Println("select { case <-t.stop }")
                return
            default :
                fmt.Println("t.run  break")
                break
        }
        //dateCmd := exec.Command(fmt.Sprintf("%s %s", t.cmd ,t.args[0]))
        dateCmd := exec.Command(t.cmd ,t.args[0])
        dateOut, err := dateCmd.Output()
        if err != nil {
            panic(err)
        }
        fmt.Println(string(dateOut))
        time.Sleep(time.Duration(t.interval) * time.Millisecond)
    }
    t.stop <- 0
}


//<线性表>---存储任务数据-----------------------------------
const MAXSIZE = 20 //定义数组长度
//定义线性表结构
type List struct {
    Element [MAXSIZE]task //存储线性表元素的数组
    length  int          //线性表长度，最大值=MAXSIZE
}
//初始化线性表,d:初始化的元素, p位置
func (l *List) InitList(d task, p int) {
    l.Element[p] = d
    l.length++
}
//存在id为id的task?index:-1
func (l *List) hasExisted(id string) int {
    for k := 0; k < l.length; k++ {
        if l.Element[k].id == id {
            fmt.Println("exist")
            return k
        }
        fmt.Print("go!")
    }
    fmt.Println("l.hasExisted:-1")
    return -1
}
//插入元素
//d:插入的数据taskdata
//p:插入位置
func (l *List) Insert(d task, p int) bool {
    if p < 0 || p >= MAXSIZE || l.length >= MAXSIZE {
        return false
    }
    if p < l.length {
        for k := l.length - 1; k >= p; k-- {
            l.Element[k+1] = l.Element[k]
        }
        l.Element[p] = d
        l.length++
        return true
    } else {
        l.Element[l.length] = d
        l.length++
        return true
    }
}
//追加元素
func (l *List) Append(d task) bool {
    if l.length == 0 {
        l.Insert(d, 0)
        l.Element[0].Start()
        return true
    }
    for k := 0; k < l.length; k++ {
        if l.Element[k].id == "" {
            l.Insert(d, k)
            l.Element[k].Start()
            return true
        }
    }
    return false
}
//删除元素
//p:删除元素的位置
func (l *List) Delete(p int) bool {
    if p < 0 || p > l.length || p >= MAXSIZE {
        return false
    }
    l.Element[p].Stop()
    //<-l.Element[p].stop
    for ; p < l.length-1; p++ {
        l.Element[p] = l.Element[p+1]
    }
    l.Element[l.length-1].id = ""
    l.length--
    return true
}
//根据id删除
func (l *List) DeleteById (id string) bool {
    index := l.hasExisted(id)
    if index < 0 {
        return false
    } else {
        l.Delete(index)
        return true
    }
}


var tasklist List


//<具体任务处理>-----------------------------------------------------
//插入并开始任务
func newAndStartTask(id string, cmd string, args []string, interval int) bool {
    newtask := NewTask(id, cmd, args, interval)
    tasklist.Append(newtask)
    //开始任务在Append()执行了
    if newtask.running == true {
        return true
    }
    //默认返回
    return false
}
//停止并删除任务
func stopAndDeleteTask(id string) bool{
    return tasklist.DeleteById(id)
}

//<http>---------------------------------------------------
type ReqBody struct {
    Id          string `json:"id"`
    Cmd         string `json:"cmd"`
    Args      []string `json:"args"`
    Interval    int    `json:"interval"`
}

type SuccessRespose struct {
    Ok     bool     `json:"ok"`
    Id     string   `json:"id"`
}

type ErrorRespose struct {
    Ok       bool     `json:"ok"`
    Error    string   `json:"error"`
}



func Task(w http.ResponseWriter, r *http.Request) {
    fmt.Println("\nTask is running...")

    //获取客户端通过POST方式传递的json
    body, _ := ioutil.ReadAll(r.Body)
    body_str := string(body)
    fmt.Println(body_str)
    var reqBody ReqBody
    err := json.Unmarshal(body, &reqBody)
    if  err == nil {
        fmt.Println(err)
        return
    }
    fmt.Println(reqBody)
    
    id := reqBody.Id
    cmd:= reqBody.Cmd
    args:= reqBody.Args
    interval:= reqBody.Interval
    
    if id == "" {
        fmt.Fprint(w, "请指定id")
        return
    }
    idexist := tasklist.hasExisted(id)
    fmt.Println("tasklist.length:",tasklist.length)  
    fmt.Println("existed task id is:",idexist)
    if cmd=="date" && args[0]=="-R" && interval>0 {
        //为POST
        fmt.Println("POST")
        //新建任务时
        //判断任务是否存在，用id
        if idexist > -1 {
            //任务存在
            //返回错误
            var result ErrorRespose
            result.Ok = false
            result.Error = fmt.Sprintf("%s%s%s",
                 "The task ", id, "already exists.")
            //向客户端返回JSON数据
            bytes, _ := json.Marshal(result)
            fmt.Fprint(w, string(bytes))
        } else {
            newAndStartTask(id, cmd, args, interval) 
            //任务不存在，创建新任务
            //返回成功
            var result SuccessRespose
            result.Id = id
            result.Ok = true
            //向客户端返回JSON数据
            bytes, _ := json.Marshal(result)
            fmt.Fprint(w, string(bytes))
        }
    } else {
        fmt.Println("DELETE")
        //为DELETE
        //删除时
        if (idexist<0) {
            //任务不存在，返回错误
            var result ErrorRespose
            result.Ok = false
            result.Error = fmt.Sprintf("%s%s%s",
                 "The Task ", id, " is not found.")
            //向客户端返回JSON数据
            bytes, _ := json.Marshal(result)
            fmt.Fprint(w, string(bytes))
        } else {
            stopAndDeleteTask(id) 
            fmt.Println("delete")
            //任务存在，根据id删除任务
            var result SuccessRespose
            //返回成功
            result.Id = id
            result.Ok = true
            //向客户端返回JSON数据
            bytes, _ := json.Marshal(result)
            fmt.Fprint(w, string(bytes))
        }
    }


}





func main() {
    fmt.Println("This is webserver base!")
    
    host := flag.String("host", "4567", "监听端口")
    flag.Parse()

    count := len(os.Args)
    fmt.Println("参数总个数:",count)
    fmt.Println("参数详情:")  
    for i := 0 ; i < count ;i++{  
        fmt.Println(i,":",os.Args[i])  
    }
    fmt.Println("host--: ", *host)
    
    if os.Args[0] != "" {
        host = &os.Args[1]
        fmt.Println("os.Args[0] != \"\"\nlistening--: ", *host)
    }
    
    //第一个参数为客户端发起http请求时的接口名，第二个参数是一个func，负责处理这个请求。
    http.HandleFunc("/", Task)
    //服务器要监听的主机地址和端口号
    err := http.ListenAndServe(fmt.Sprintf(":%s", *host), nil)

    if err != nil {
        fmt.Println("ListenAndServe error: ", err.Error())
    }
}



//使用
// $   go build gocron.go
// $   ./gocron 8005
// 8005为端口号



































