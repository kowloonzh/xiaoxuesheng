package main

import (
	"bufio"
	"container/list"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type English struct {
	Units []EnglishUnit `json:"units"`
}

type EnglishUnit struct {
	Id string `json:"id"`
	// 单元名称
	Name  string        `json:"name"`
	Words []EnglishWord `json:"words"`
}

type EnglishWord struct {
	Sentence string `json:"sentence"`
	Word     string `json:"word"`
	Cue      string `json:"cue"`
}

var english English
var queue = list.New()

func main() {
	start1 := "游戏马上开始，请按照提示进行选择或者填空，然后按回车键\n"
	start2 := "如果你不知道填写什么，输入问号 ? 查看答案\n"
	start3 := "如果你想结束游戏，输入句号。\n"
	printChracter(start1)
	printChracter(start2)
	printChracter(start3)

	loadCourse()

	reader := bufio.NewReader(os.Stdin)

	// 请选择单元
	m := make(map[string]EnglishUnit)
UNIT:
	chooseUnit := "请选择英语单元："
	for _, v := range english.Units {
		chooseUnit += fmt.Sprintf("%s ", v.Id)
		m[v.Id] = v
	}
	chooseUnit += "\n"
	printChracter(chooseUnit)

	unit, _ := reader.ReadString('\n')
	unit = strings.TrimSpace(unit)
	englistUnit, ok := m[unit]
	if ok {
		englistUnit = m[unit]
	} else {
		fmt.Printf("输入的单元 %s 不存在, 请重新选择\n", unit)
		goto UNIT
	}

	// 添加到队列
	for k := range englistUnit.Words {
		queue.PushBack(englistUnit.Words[k])
	}

	for queue.Len() > 0 {
		// 清屏幕
		clear()
		item := queue.Front()

		v := item.Value.(EnglishWord)

		printChracter(fmt.Sprintf("%s %s\n", v.Sentence, v.Cue))

		queue.Remove(item)

		// 读取标准输入
		input, _ := reader.ReadString('\n')
		// fmt.Println("你输入的是:", input)
		input = strings.TrimSpace(input)

		if isQuestion(input) {
			fmt.Printf("答案是: %s", v.Word)
			queue.PushBack(v)
			time.Sleep(time.Second * 2)
		} else if isOver(input) {
			fmt.Println("游戏结束")
			time.Sleep(time.Second * 2)
			os.Exit(0)
		} else {
			// 判断是否正确
			if input == v.Word {
				fmt.Println("✅ 回答正确")
				// 回答正确后，清理 item
				time.Sleep(time.Second * 1)
			} else {
				fmt.Printf("❎,正确答案是%s\n\n", v.Word)
				queue.PushBack(v)
				time.Sleep(time.Second * 2)
			}
		}

	}

	printChracter("💐恭喜你，游戏通关，✿✿ヽ(°▽°)ノ✿")
	time.Sleep(time.Second * 2)
}

func loadCourse() {

	path,err := os.Getwd()
	if err != nil {
		panic(err)
	}

	vip := viper.New()
	vip.AddConfigPath(path)
	vip.SetConfigName("conf")
	vip.SetConfigType("json")

	// 尝试进行配置读取
	if err := vip.ReadInConfig();err != nil {
		panic(err)
	}

	if err := vip.Unmarshal(&english);err != nil {
		panic(err)
	}

	// english = English{}

	// englishUnits := make([]EnglishUnit, 0)

	// unitOne := EnglishUnit{
	// 	Id:   "1",
	// 	Name: "单元一",
	// 	Words: []EnglishWord{
	// 		{Sentence: "What is your __?", Word: "name", Cue: "你叫什么名字？"},
	// 		{Sentence: "__, Miss Zheng?", Word: "hi", Cue: "你好，郑老师？"},
	// 	},
	// }

	// unitFifth := EnglishUnit{
	// 	Id:   "15",
	// 	Name: "单元十五",
	// 	Words: []EnglishWord{
	// 		{Sentence: "The __ Art Show is coming.", Word: "school", Cue: "学校艺术节到来了"},
	// 		{Sentence: "The School Art __ is coming.", Word: "show", Cue: "学校艺术节到来了"},
	// 		{Sentence: "What __ you do?", Word: "can", Cue: "你们能干什么？"},
	// 		{Sentence: "What can you __?", Word: "do", Cue: "你们能干什么？"},
	// 		{Sentence: "I can __.", Word: "dance", Cue: "我能跳舞"},
	// 		{Sentence: "I can __ the piano.", Word: "play", Cue: "我能弹钢琴"},
	// 		{Sentence: "I can play the __.", Word: "piano", Cue: "我能弹钢琴"},
	// 		{Sentence: "I can __ Beijing opera.", Word: "sing", Cue: "我能唱京剧"},
	// 		{Sentence: "__ can play together in the show.", Word: "we", Cue: "我们可以一起表演"},
	// 	},
	// }

	// unitSix := EnglishUnit{
	// 	Id:   "16",
	// 	Name: "单元十六",
	// 	Words: []EnglishWord{
	// 		{Sentence: "Can you play __, Baobao?", Word: "ping-pong", Cue: "宝宝，你可以打乒乓球吗？"},
	// 		{Sentence: "__'s go", Word: "let", Cue: "我们走"},
	// 		{Sentence: "Can you __ rope?", Word: "jump", Cue: "你可以跳绳吗？"},
	// 		{Sentence: "Can you jump __?", Word: "rope", Cue: "你可以跳绳吗？"},
	// 		{Sentence: "__, I can't.", Word: "no", Cue: "不，我不会"},
	// 		{Sentence: "It's __.", Word: "easy", Cue: "这很简单"},
	// 	},
	// }

	// englishUnits = append(englishUnits, unitOne, unitFifth, unitSix)
	// english.Units = englishUnits
}

func printChracter(s string) {
	for _, v := range s {
		fmt.Printf("%c", v)
		time.Sleep(time.Millisecond * 100)
	}
}

func isQuestion(q string) bool {
	return q == "?" || q == "？"
}

func isOver(q string) bool {
	return q == "。"
}

func clear() {
	// Check the operating system
	myos := runtime.GOOS

	// Create the appropriate command based on the operating system
	var cmd *exec.Cmd
	switch myos {
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls")
	case "linux", "darwin": // macOS
		cmd = exec.Command("clear")
	default:
		// Unsupported platform
		fmt.Println("Unsupported platform:", myos)
		return
	}

	// Execute the command
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error clearing screen:", err)
		return
	}
}
