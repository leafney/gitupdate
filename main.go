package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
)

var (
	baseDirPath             = "."   // 操作的基目录
	onlyOneOrSomeDirName    = ""    // 只操作基目录下的某一个或多个目录
	excludeOneOrSomeDirName = ""    // 排除基目录下的某一个或多个目录
	branchCount             = 3     // 获取最近分支的个数
	isCheck                 = false // 检查模式，验证目录筛选结果
)

func main() {
	/*
	   1. 遍历基目录下的所有目录
	   2. 执行shell命令

	   3. ➜ go run main.go -bd=./test -od=celery
	   4. ➜ go build -o gitupdate main.go
	   5. ➜ go run main.go -bd=./test -od=celery -bc=5
	*/
	flag.StringVar(&baseDirPath, "bd", ".", "defalut current path")
	flag.StringVar(&onlyOneOrSomeDirName, "od", "", "only one or some dir name mutiple between with ,")
	flag.StringVar(&excludeOneOrSomeDirName, "ed", "", "exclude one or some dir name mutiple between with ,")
	flag.IntVar(&branchCount, "bc", 3, "pull branch count default 3")
	flag.BoolVar(&isCheck, "ck", false, "verify directory filtering results")

	flag.Parse()

	// 获取所有目录
	allDirs, err := GetDirs(baseDirPath, onlyOneOrSomeDirName, excludeOneOrSomeDirName)
	if err != nil {
		panic(err)
	}

	fmt.Printf("筛选后的目录：\n%s\n", allDirs)

	if isCheck {
		fmt.Println("Check mode done.")
		return
	}

	fmt.Println("start range")
	for _, dd := range allDirs {
		fmt.Println("Start for path: " + dd)
		fmt.Println("↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓")

		// 获取最新 git pull
		pllNews, err := exec.Command("git", "-C", dd, "pull").CombinedOutput()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
			// panic(err)
		}
		fmt.Println("**获取最新分支信息**")
		fmt.Println(string(pllNews))
		time.Sleep(3)

		// Todo 获取该项目下最近提交的三个分支并依次下载 git for-each-ref --sort=-committerdate --format='%(refname:short)' --count=3 refs/remotes/origin/
		bchs, err := exec.Command("git", "-C", dd, "for-each-ref", "--sort=-committerdate", "--format=%(refname:short)", fmt.Sprintf("--count=%d", branchCount), "refs/remotes/origin/").CombinedOutput()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
			// panic(err)
		}
		fmt.Println("**筛选出符合条件的分支**")
		fmt.Println(string(bchs))
		time.Sleep(3)

		// 从输出中得到最近分支名
		/*
			origin/login-gsc
			origin/login-rb
			origin/login-rb
		*/
		br := bufio.NewReader(bytes.NewReader(bchs))
		for {
			ba, _, bc := br.ReadLine()
			if bc == io.EOF {
				break
			}
			baBranch := string(ba)

			// 遇到 HEAD 分支，跳过
			if baBranch == "origin/HEAD" {
				continue
			}

			fmt.Println("**开始拉取远端分支最新提交**: ", baBranch)
			// 截取分支名
			newBranch := strings.TrimPrefix(baBranch, "origin/")
			fmt.Println("**当前操作分支**: ", newBranch)

			// 切换到新分支 git checkout master
			a, err := exec.Command("git", "-C", dd, "checkout", newBranch).CombinedOutput()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
				// panic(err)
			}
			fmt.Println(string(a))
			time.Sleep(3)

			// 获取新分支 git pull origin master:master
			b, err := exec.Command("git", "-C", dd, "pull", "origin", fmt.Sprintf("%s:%s", newBranch, newBranch)).CombinedOutput()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
				// panic(err)
			}
			fmt.Println(string(b))
			time.Sleep(5)
		}

		// // a, _ := RunCommands
		// a, err := exec.Command("git", "-C", dd, "checkout", "master").CombinedOutput()
		// if err != nil {
		// 	fmt.Println(err)
		// 	os.Exit(1)
		// }
		// fmt.Println(string(a))
		// time.Sleep(3)

		// b, err := exec.Command("git", "-C", dd, "pull", "origin", "master:master").CombinedOutput()
		// if err != nil {
		// 	fmt.Println(err)
		// 	os.Exit(1)
		// }
		// fmt.Println(string(b))
		fmt.Println("↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑")
		// time.Sleep(10)

		// old
		// RunCommands("git checkout master")
		// time.Sleep(2)
		// RunCommands("git pull origin master:master")
		// time.Sleep(10)
		// d, _ := RunCommands("cd ..")
		// fmt.Println(d)
	}
	fmt.Println("All directories have been updated.")
}

// func RunCommands(cmdStr ...string) (outStr string, err error) {
// 	cmd := exec.Command("/bin/bash", cmdStr...)
// 	// out, err := cmd.Output()
// 	out, err := cmd.CombinedOutput()
// 	if err != nil {
// 		fmt.Println(err)
// 		os.Exit(1)
// 	}

// 	return strings.Trim(string(out), "\n"), nil
// }

func GetDirs(dBasePath, oneDirName, excludeDirName string) (dirs []string, err error) {
	fmt.Println("筛选条件: basepath: " + dBasePath + " ; oneDirName: " + oneDirName + " ; excludeDirName: " + excludeDirName)

	// 获取指定路径下的所有目录
	dir, err := ioutil.ReadDir(dBasePath)
	if err != nil {
		return nil, err
	}

	for _, fi := range dir {
		if fi.IsDir() {
			// 如果是目录，则加入集合中
			theDirName := fi.Name()
			theDirPath := path.Join(dBasePath, theDirName)
			// fmt.Println("path.Dir(dBasePath):  " + dBasePath)
			// fmt.Println("path.Join  " + theDirPath)

			// 排除指定的目录
			if CheckIsIn(theDirName, excludeDirName) {
				fmt.Printf("排除目录: %s 不处理\n", theDirPath)
				continue
			}

			if oneDirName != "" {
				// 如果只处理指定的一个或多个目录
				if CheckIsIn(theDirName, oneDirName) {
					fmt.Printf("只处理目录: %s \n", theDirPath)
					dirs = append(dirs, theDirPath)
				}
			} else {
				fmt.Printf("检索到目录: %s \n", theDirPath)
				dirs = append(dirs, theDirPath)
			}
		}
	}
	return dirs, nil
}

func CheckIsIn(oName, aNames string) (result bool) {
	result = false

	allDirNames := strings.Split(aNames, ",")
	for _, ad := range allDirNames {
		if oName == ad {
			result = true
			break
		}
	}
	return result
}
