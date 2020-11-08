package main

import (
	"../Calculator/stack"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"unicode"
)
func sayHello1(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()//获取请求参数
	fmt.Printf(r.FormValue("data"))
	stat:=r.FormValue("data")
	stat = strings.TrimSpace(stat)//去掉字符串首尾空白字符
	postfix := infix2ToPostfix(stat)
	fmt.Printf("后缀表达式：%s\n", postfix)
	fmt.Printf("计算结果: %f\n", calculate(postfix))
	fmt.Fprintf(w,"%f", calculate(postfix))//返回前台页面

}
func  index1(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("22.tmpl")
	if err != nil {
		fmt.Fprintf(w, "Parse: %v", err)
		return
	}
	err = tmpl.Execute(w, "11.tmpl")
	if err != nil {
		fmt.Fprintf(w, "Execute: %v", err)
		return
	}
}


func calculate(postfix string) float32 {
	stack := stack.ItemStack{}
	fixLen := len(postfix)
	for i := 0; i < fixLen; i++ {
		nextChar := string(postfix[i])
		// 数字：直接压栈
		if unicode.IsDigit(rune(postfix[i])) {
			stack.Push(nextChar)
		} else {
			// 操作符：取出两个数字计算值，再将结果压栈
			num1, _ := strconv.Atoi(stack.Pop())
			num2, _ := strconv.Atoi(stack.Pop())
			switch nextChar {
			case "+":
				stack.Push(strconv.Itoa(num1 + num2))
			case "-":
				stack.Push(strconv.Itoa(num2 - num1))
			case "*":
				stack.Push(strconv.Itoa(num1 * num2))
			case "/":
				stack.Push(strconv.Itoa(num2 / num1))
             case "√":
                 stack.Push(strconv.Ttoa(num1**0.5))
             case "2":
                 stack.Push(strconv.Ttoa(num1**2))
			}
		}
	}
	result, _ := strconv.Atoi(stack.Top())
	return float32(result)
}

func infix2ToPostfix(exp string) string {
	stack := stack.ItemStack{}
	postfix := ""
	expLen := len(exp)

	// 遍历整个表达式
	for i := 0; i < expLen; i++ {

		char := string(exp[i])

		switch char {
		case " ":
			continue
		case "(":
			// 左括号直接入栈
			stack.Push("(")
		case ")":
			// 右括号则弹出元素直到遇到左括号
			for !stack.IsEmpty() {
				preChar := stack.Top()
				if preChar == "(" {
					stack.Pop() // 弹出 "("
					break
				}
				postfix += preChar
				stack.Pop()
			}

			// 数字则直接输出
		case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
			j := i
			digit := ""
			for ; j < expLen && unicode.IsDigit(rune(exp[j])); j++ {
				digit += string(exp[j])
			}
			postfix += digit
			i = j - 1 // i 向前跨越一个整数，由于执行了一步多余的 j++，需要减 1

		default:
			// 操作符：遇到高优先级的运算符，不断弹出，直到遇见更低优先级运算符
			for !stack.IsEmpty() {
				top := stack.Top()
				if top == "(" || isLower(top, char) {
					break
				}
				postfix += top
				stack.Pop()
			}
			// 低优先级的运算符入栈
			stack.Push(char)
		}
	}

	// 栈不空则全部输出
	for !stack.IsEmpty() {
		postfix += stack.Pop()
	}

	return postfix
}

func isLower(top string, newTop string) bool {
	// 注意 a + b + c 的后缀表达式是 ab + c +，不是 abc + +
	switch top {
	case "+", "-":
		if newTop == "*" || newTop == "/" {
			return true
		}
	case "(":
		return true
	}
	return false
}
func main() {
	http.HandleFunc("/", index1)
	http.HandleFunc("/hello",sayHello1)
	log.Print("Starting server...")
	log.Fatal(http.ListenAndServe(":4000", nil))
}
