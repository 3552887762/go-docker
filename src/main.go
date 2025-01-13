package main

import (
	"fmt"
	"golang.org/x/sys/unix" // 使用 unix 包
	"os"
	"os/exec"
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
	cmd.SysProcAttr = &unix.SysProcAttr{
		Cloneflags: unix.CLONE_NEWUTS | unix.CLONE_NEWNS | unix.CLONE_NEWPID, // 设置新的 UTS 名称空间和 PID 名称空间
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

func Init() {
	// 设置容器的主机名为 "container"
	fmt.Println("Setting hostname to 'container'...")
	if err := unix.Sethostname([]byte("container")); err != nil {
		fmt.Printf("Error setting hostname: %v\n", err)
		panic(err)
	}

	// 切换到新的根目录
	fmt.Println("Changing root to 'root'...")
	if err := unix.Chroot("root"); err != nil {
		fmt.Printf("Error changing root: %v\n", err)
		panic(err)
	}

	// 切换到根目录
	fmt.Println("Changing current directory to '/'...")
	if err := unix.Chdir("/"); err != nil {
		fmt.Printf("Error changing directory: %v\n", err)
		panic(err)
	}

	// 挂载 /proc 文件系统
	fmt.Println("Mounting /proc...")
	if err := unix.Mount("proc", "/proc", "proc", 0, ""); err != nil {
		fmt.Printf("Error mounting /proc: %v\n", err)
		panic(err)
	}

	// 执行指定的程序
	fmt.Println("Executing program...")
	fmt.Printf("Program to execute: %v\n", os.Args[2])
	if err := unix.Exec(os.Args[2], os.Args[2:], os.Environ()); err != nil {
		fmt.Printf("Error executing program: %v\n", err)
		panic(err)
	}

	// 卸载 /proc 文件系统
	fmt.Println("Unmounting /proc...")
	if err := unix.Unmount("/proc", 0); err != nil {
		fmt.Printf("Error unmounting /proc: %v\n", err)
	}
}
