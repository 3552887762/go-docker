package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

// sudo docker run /bin/bash
func main() {
	// 打印当前进程的命令行参数和进程ID
	fmt.Printf("Process => %v [%d]\n", os.Args, os.Getpid())

	// 根据命令行参数选择执行的功能
	switch os.Args[1] {
	case "run":
		// 如果命令是 "run"，调用 Run() 函数
		Run()
	case "init":
		// 如果命令是 "init"，调用 Init() 函数
		Init()
	default:
		// 如果没有匹配的命令，程序 panic
		panic("未定义的命令")
	}
	fmt.Println("结束")
}

// Run 函数模拟容器的运行过程
func Run() {
	// 创建一个新的进程，执行 "init" 命令，并传递相关参数
	cmd := exec.Command(os.Args[0], "init", os.Args[2])
	// 设置进程的属性，使用克隆标志来模拟容器的进程隔离
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID, // 设置新的 UTS 名称空间和 PID 名称空间
	}
	// 将当前进程的输入输出重定向到新进程
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// 执行新进程，捕获执行错误
	if err := cmd.Run(); err != nil {
		// 如果执行失败，程序 panic
		panic(err)
	}
}

// Init 函数模拟初始化容器环境
func Init() {
	// 设置容器的主机名为 "container"
	syscall.Sethostname([]byte("container"))

	// 挂载 /proc 文件系统，模拟容器的 proc 文件系统
	// 		自己的pid的proc ， 主机的proc
	syscall.Mount("proc", "/proc", "proc", 0, "")

	// 执行指定的程序，这里假设 os.Args[2] 是程序路径
	// os.Args[2:] 是参数列表，os.Environ() 获取当前环境变量
	syscall.Exec(os.Args[2], os.Args[2:], os.Environ())

	// 如果 Exec 执行成功，下面的代码不会被执行
	// 卸载 /proc 文件系统，如果 Exec 失败则执行卸载操作
	syscall.Unmount("/proc", 0)
}
