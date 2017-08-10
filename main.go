package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"log"
	"time"
	"flag"
	"bytes"
	"github.com/animotto/thg"
	"github.com/chzyer/readline"
	"github.com/robertkrimen/otto"
	"github.com/spf13/viper"
)

const (
	prompt string = "» "
	histFile string = "sandbox.history"
	configsDir = "configs"
	scriptsDir = "scripts"
)

var (
	context string = "/"
	rl *readline.Instance
	vmScript otto.Otto
)

func rlDynamicCompleter() func(string) []string {
	return func(line string) []string {
		var commands []string
		switch context {
		case "/":
			commands = []string{
				"query", "script", "chat", "auth", "checkconn",
				"mynet", "node", "program", "world", "help",	"exit",
			}
		case "/query":
			commands = []string{
				"rq", "sq", "pm", "help", "exit",
			}
		case "/script":
			commands = []string{
				"run", "go", "eval", "help", "exit",
			}
		case "/chat":
			commands = []string{
				"open", "close", "list", "say", "help", "exit",
			}
		}
		return commands
	}
}

func main() {
	configFile := flag.String("c", configsDir + "/sandbox.conf", "Configuration file")
	scriptFile := flag.String("s", "", "Script file")
	scriptArgs := flag.String("a", "", "Script arguments")
	flag.Parse()

	rlCompleter := readline.NewPrefixCompleter(
		readline.PcItemDynamic(rlDynamicCompleter()),
	)

	var err error
	rl, err = readline.NewEx(&readline.Config{
		HistoryFile: histFile,
		HistoryLimit: 100,
		AutoComplete: rlCompleter,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer rl.Close()

	viper.SetConfigFile(*configFile)
	viper.SetConfigType("json")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

	thgClient := thg.New(thg.Config{
		Address: viper.GetString("Address"),
		ReqUrl: viper.GetString("ReqUrl"),
		HashSalt: viper.GetString("HashSalt"),
		AppVersion: uint16(viper.GetInt("AppVersion")),
		IdPlayer: uint32(viper.GetInt("IdPlayer")),
		Password: viper.GetString("Password"),
	})

	vmScript := otto.New()
	vmScript.Set("sleep", func(timer int) {
		time.Sleep(time.Duration(timer) * time.Second)
	})
	vmScript.Set("print", func(text string) {
		rl.Clean()
		fmt.Print(text)
		//fmt.Fprint(rl, text)
	})
	vmScript.Set("println", func(text string) {
		fmt.Println(text)
		//fmt.Fprintln(rl, text)
	})
	vmScript.Set("config", thgClient.GetConfig())
	vmScript.Set("sessionid", func() otto.Value {
		ret, _ := vmScript.ToValue(thgClient.SessionId)
		return ret
	})
	vmScript.Set("rquery", func(url string) otto.Value {
		var ret otto.Value
		data, err := thgClient.HttpGet(thgClient.GetFullUrl(url))
		if err != nil {
			// log.Println(err)
			ret, _ = vmScript.ToValue(nil)
		} else {
			ret, _ = vmScript.ToValue(data)
		}
		return ret
	})
	vmScript.Set("squery", func(url string) otto.Value {
		var ret otto.Value
		url += "&session_id=" + fmt.Sprint(thgClient.SessionId)
		data, err := thgClient.HttpGet(thgClient.GetFullUrl(url))
		if err != nil {
			// log.Println(err)
			ret, _ = vmScript.ToValue(nil)
		} else {
			ret, _ = vmScript.ToValue(data)
		}
		return ret
	})
	vmScript.Set("auth", func() otto.Value {
		var ret otto.Value
		data, err := thgClient.Auth()
		if err != nil {
			// log.Println(err)
			ret, _ = vmScript.ToValue(nil)
		} else {
			ret, _ = vmScript.ToValue(data)
		}

		return ret
	})
	vmScript.Set("checkconn", func() otto.Value {
		var ret otto.Value
		err := thgClient.CheckConn()
		if err != nil {
			// log.Println(err)
			ret, _ = vmScript.ToValue(false)
		} else {
			ret, _ = vmScript.ToValue(true)
		}

		return ret
	})
	vmScript.Set("mynet", func() otto.Value {
		var ret otto.Value
		data, err := thgClient.NetMaint()
		if err != nil {
			// log.Println(err)
			ret, _ = vmScript.ToValue(nil)
		} else {
			ret, _ = vmScript.ToValue(data)
		}

		return ret
	})
	vmScript.Set("world", func() otto.Value {
		var ret otto.Value
		data, err := thgClient.GetWorld()
		if err != nil {
			ret, _ = vmScript.ToValue(nil)
		} else {
			ret, _ = vmScript.ToValue(data)
		}

		return ret
	})
	vmScript.Set("boncoll", func(id int) otto.Value {
		var ret otto.Value
		err := thgClient.BonusCollect(uint32(id))
		if err != nil {
			ret, _ = vmScript.ToValue(nil)
		} else {
			ret, _ = vmScript.ToValue(true)
		}

		return ret
	})
	vmScript.Set("goaltypes", func() otto.Value {
		var ret otto.Value
		goalTypes, err := thgClient.GetGoalTypes()
		if err != nil {
			ret, _ = vmScript.ToValue(nil)
		} else {
			ret, _ = vmScript.ToValue(goalTypes)
		}

		return ret
	})
	vmScript.Set("goalupd", func(id int, record int) otto.Value {
		var ret otto.Value
		err := thgClient.GoalUpdate(uint32(id), byte(record))
		if err != nil {
			ret, _ = vmScript.ToValue(nil)
		} else {
			ret, _ = vmScript.ToValue(true)
		}

		return ret
	})
	vmScript.Set("getchat", func(room int, lm string) otto.Value {
		var ret otto.Value
		data, err := thgClient.ChatDisplay(uint16(room), lm)
		if err != nil {
			ret, _ = vmScript.ToValue(nil)
		} else {
			ret, _ = vmScript.ToValue(data)
		}

		return ret
	})
	vmScript.Set("netupd", func(net string) otto.Value {
		var ret otto.Value
		err := thgClient.NetUpdate(net)
		if err != nil {
			ret, _ = vmScript.ToValue(false)
		} else {
			ret, _ = vmScript.ToValue(true)
		}

		return ret
	})

	if len(*scriptFile) != 0 {
		f, err := os.Open(*scriptFile)
		if err != nil {
			log.Fatal(err)
		}

		defer f.Close()
		buff := bytes.NewBuffer(nil)
		if _, err := buff.ReadFrom(f); err != nil {
			log.Fatal(err)
		}

		vmScript.Set("args", strings.Fields(*scriptArgs))

		if _, err := vmScript.Run(buff.String()); err != nil {
			log.Println(err)
		}

		os.Exit(0)
	}

	chanChat := make(map[uint16]chan string)

	for {
		rl.SetPrompt("\x1b[36m" + context + "\x1b[31m" + prompt + "\x1b[0m")
		cmdline, err := rl.Readline()
		if err == io.EOF {
			os.Exit(0)
		} else if err == readline.ErrInterrupt {
			if len(cmdline) != 0 {
				continue
			} else {
				os.Exit(0)
			}
		} else if err != nil {
			log.Fatal(err)
		}

		cmd := strings.Fields(cmdline)

		if len(cmd) == 0 {
			continue
		}

		switch cmd[0] {
		case "exit":
			os.Exit(0)
		case "/":
			context = "/"
			continue
		case "..":
			path := strings.Split(context, "/")
			context = "/"
			for i:=0; i<len(path)-2; i++ {
				context += "/"
			}
			continue
		}

		switch context {
		case "/":
			contextRoot(cmd, &thgClient)
		case "/query":
			contextQuery(cmd, &thgClient)
		case "/script":
			contextScript(cmd, vmScript)
		case "/chat":
			contextChat(cmd, &thgClient, chanChat)
		}
	}
}
