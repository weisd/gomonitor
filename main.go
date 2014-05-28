package main

import(
    "os/exec"
    "bytes"
    "fmt"
    "runtime"
    "log"
    "io"
    "strings"
    "os"
)

func main(){

    args := os.Args
    proName := args[1]

    fmt.Println("[INFO]检测程序", proName)

    if !proExists(proName){
        fmt.Println("[WAMM]进程不存在, 开始重启进程...")
        runCmd := exec.Command(proName)

        var out bytes.Buffer

        runCmd.Stdout = &out 
        err := runCmd.Run()
        if err != nil {
            fmt.Println("[ERRO]程序执行失败", err)
            return
        }

        fmt.Println(out.String())

        if proExists(proName){
            fmt.Println("[INFO]程序执行完毕")
            return
        }
    }

    fmt.Println("[INFO]程序正在运行", proName)
}


func proExists(proName string) bool {

    var tasklistCmd *exec.Cmd
    var isExists = false
    var isWin = runtime.GOOS == "windows"

    if isWin {
        proName = strings.Replace(proName, "\\", "/", -1)
        nIdx := strings.LastIndex(proName, "/")
        if nIdx != -1 {
            proName = proName[nIdx:]
        }
        tasklistCmd = exec.Command("tasklist.exe")
    } else {
        tasklistCmd = exec.Command("sh", "-c", fmt.Sprintf("ps -ef | grep %s | grep -v \"%s\"| grep -v grep", proName, os.Args[0]))
//        tasklistCmd = exec.Command("ps", "-ef")
    }

    var out bytes.Buffer

    tasklistCmd.Stdout = &out

    err := tasklistCmd.Run()
    if err != nil {
        log.Fatal("[ERRO]查询失败：", err)
    }

    for{
        line, err := out.ReadString('\n')
        if err == io.EOF {
            break
        }

        idx := strings.Index(line, proName+ " ")
        if idx != -1 {
            fmt.Println("[info]找到程序 "+line)
            isExists = true
            continue
        }

    }

    return isExists
}
