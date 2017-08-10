package main

import (
	"fmt"
	"strconv"
	"os"
	"log"
	"bytes"
	"strings"
	"time"
	"thg"
	"github.com/robertkrimen/otto"
)

var (
	queryPm string = "none"
)

// Root context
func contextRoot(cmd []string, thg *thg.Thg) {
	switch cmd[0] {
	case "query", "script", "chat":
		context = context + cmd[0]
	case "config":
		fmt.Printf("%-15s .. %s\n", "Address", thg.Address)
		fmt.Printf("%-15s .. %s\n", "ReqUrl", thg.ReqUrl)
		fmt.Printf("%-15s .. %s\n", "HashSalt", thg.HashSalt)
		fmt.Printf("%-15s .. %d\n", "AppVersion", thg.AppVersion)
		fmt.Printf("%-15s .. %d\n", "IdPlayer", thg.IdPlayer)
		fmt.Printf("%-15s .. %s\n", "Password", thg.Password)
	case "auth":
		fmt.Print("Getting session ID: ")
		_, err := thg.Auth()
		if err != nil {
			fmt.Println("Error")
			return
		}
		fmt.Println(thg.SessionId)
	case "checkconn":
		if len(thg.SessionId) == 0 {
			fmt.Println("No session ID")
			return
		}

		fmt.Print("Checking connection: ")
		err := thg.CheckConn()
		if err != nil {
			fmt.Println("Error")
			return
		}
		fmt.Println("OK")
	case "mynet":
		var section string
		if len(cmd) == 2 {
			section = cmd[1]
			if section != "profile" && section != "nodes" &&
				section != "programs" && section != "queue" &&
				section != "log" {
				fmt.Println("Unknown section")
				return
				}
		}

		if len(thg.SessionId) == 0 {
			fmt.Println("No session ID")
			return
		}

		fmt.Print("Getting my network: ")
		net, err := thg.NetMaint()
		if err != nil {
			fmt.Println("Error")
			return
		}
		fmt.Println("OK")

		if section == "profile" || len(section) == 0 {
			fmt.Println()
			fmt.Println("\x1b[36mProfile:")
			fmt.Printf(" \x1b[35m %-15s\x1b[0m .. %d\n", "ID", net.Profile.Id)
			fmt.Printf(" \x1b[35m %-15s\x1b[0m .. %s\n", "Name", net.Profile.Name)
			fmt.Printf(" \x1b[35m %-15s\x1b[0m .. %d\n", "Gold", net.Profile.Gold)
			fmt.Printf(" \x1b[35m %-15s\x1b[0m .. %d\n", "Bitcoins", net.Profile.Bitcoins)
			fmt.Printf(" \x1b[35m %-15s\x1b[0m .. %d\n", "Credits", net.Profile.Credits)
			fmt.Printf(" \x1b[35m %-15s\x1b[0m .. %d\n", "Builders", net.Profile.Builders)
			fmt.Printf(" \x1b[35m %-15s\x1b[0m .. %d\n", "Rank", net.Profile.Rank)
			fmt.Printf(" \x1b[35m %-15s\x1b[0m .. %d\n", "Experience", net.Profile.Experience)
			fmt.Printf(" \x1b[35m %-15s\x1b[0m .. %d\n", "Country", net.Profile.Country)
		}
		if section == "nodes" || len(section) == 0 {
			fmt.Println()
			fmt.Println("\x1b[36mNodes:")
			for i := range net.Nodes {
				fmt.Printf(" \x1b[35m%-15d\x1b[0m .. Type %-3d | Level %-3d | Timer %-10d\n",
					net.Nodes[i].Id, net.Nodes[i].Ntype, net.Nodes[i].Level, net.Nodes[i].UpgradeTimer)
			}
		}
		if section == "programs" || len(section) == 0 {
			fmt.Println()
			fmt.Println("\x1b[36mPrograms:")
			for i := range net.Programs {
				fmt.Printf(" \x1b[35m%-15d\x1b[0m .. Type %-3d | Level %-3d | Amount %-4d | Timer %-10d\n",
					net.Programs[i].Id, net.Programs[i].Ptype, net.Programs[i].Level, net.Programs[i].Amount, net.Programs[i].UpgradeTimer)
			}
		}
		if section == "queue" || len(section) == 0 {
			fmt.Println()
			fmt.Println("\x1b[36mQueue:")
			for i := range net.Queue {
				fmt.Printf(" \x1b[35m%-3d\x1b[0m .. Amount %-4d\n",
					net.Queue[i].Ptype, net.Queue[i].Amount)
			}
		}
		if section == "log" || len(section) == 0 {
			fmt.Println()
			fmt.Println("\x1b[36mLog attacks:")
			for i := range net.LogAttacks {
				fmt.Printf(" \x1b[35m%-15s\x1b[0m .. %-15s | Rank %-3d\n",
					net.LogAttacks[i].DateTime, net.LogAttacks[i].AttackerName, net.LogAttacks[i].Rank)
			}
		}
	case "node":
		if len(thg.SessionId) == 0 {
			fmt.Println("No session ID")
			return
		}
		if len(cmd) < 2 {
			return
		}
		switch cmd[1] {
		case "create":
			if len(cmd) != 3 {
				return
			}

			if v, e := strconv.Atoi(cmd[2]); e == nil {
				fmt.Printf("Creating node %d: ", v)
				if err := thg.NodeCreate(byte(v), ""); err != nil {
					fmt.Println("Error")
					return
				}
				fmt.Println("OK")
			}
		case "delete":
			if len(cmd) != 3 {
				return
			}

			if v, e := strconv.Atoi(cmd[2]); e == nil {
				fmt.Printf("Deleting node %d: ", v)
				if err := thg.NodeDelete(uint32(v), ""); err != nil {
					fmt.Println("Error")
					return
				}
				fmt.Println("OK")
			}
		case "upgrade":
			if len(cmd) != 3 {
				return
			}

			if v, e := strconv.Atoi(cmd[2]); e == nil {
				fmt.Printf("Upgrading node %d: ", v)
				if err := thg.NodeUpgrade(uint32(v)); err != nil {
					fmt.Println("Error")
					return
				}
				fmt.Println("OK")
			}
		case "builders":
			if len(cmd) != 4 {
				return
			}

			if v1, e := strconv.Atoi(cmd[2]); e == nil {
				if v2, e := strconv.Atoi(cmd[3]); e == nil {
					fmt.Printf("Setting %d builders for node %d: ", v2, v1)
					if err := thg.NodeSetBuilders(uint32(v1), byte(v2)); err != nil {
						fmt.Println("Error")
						return
					}
					fmt.Println("OK")
				}
			}
		case "collect":
			if len(cmd) != 3 {
				return
			}

			if v, e := strconv.Atoi(cmd[2]); e == nil {
				fmt.Printf("Collecting resources from node %d: ", v)
				if err := thg.NodeCollect(uint32(v)); err != nil {
					fmt.Println("Error")
					return
				}
				fmt.Println("OK")
			}
		}
	case "program":
		if len(thg.SessionId) == 0 {
			fmt.Println("No session ID")
			return
		}
		if len(cmd) < 2 {
			return
		}
		switch cmd[1] {
		case "upgrade":
			if len(cmd) != 3 {
				return
			}

			if v, e := strconv.Atoi(cmd[2]); e == nil {
				fmt.Printf("Upgrading program %d: ", v)
				if err := thg.ProgramUpgrade(uint32(v)); err != nil {
					fmt.Println("Error")
					return
				}
				fmt.Println("OK")
			}
		}
	case "world":
		if len(thg.SessionId) == 0 {
			fmt.Println("No session ID")
			return
		}
		if thg.Profile.Id == 0 {
			fmt.Println("No my network information")
			return
		}

		fmt.Print("Getting world: ")
		world, err := thg.GetWorld()
		if err != nil {
			fmt.Println("Error")
			return
		}
		fmt.Println("OK")

		fmt.Println()
		fmt.Println("Players:")
		for _, v := range world.Players {
			fmt.Printf(" %-15d .. %s\n", v.Id, v.Name)
		}

		fmt.Println()
		fmt.Println("Bonuses:")
		for _, v := range world.Bonuses {
			fmt.Printf(" %-15d .. %d credits\n", v.Id, v.Amount)
		}

		fmt.Println()
		fmt.Println("Goals:")
		for _, v := range world.Goals {
			fmt.Printf(" %-15d\n", v.Id)
		}
	case "help":
		fmt.Println("[query] - Raw queries")
		fmt.Println("[script] - Script engine")
		fmt.Println("[chat] - Chat")
		fmt.Println("auth - Authenticate")
		fmt.Println("checkconn - Check connection")
		fmt.Println("mynet [profile|nodes|programs|queue|log] - Get my network information")
		fmt.Println("node <create|delete|upgrade|builders|collect> <id> - Create/Delete/Upgrade/Builders/Collect nodes")
		fmt.Println("program <upgrade> <id> - Upgrade programs")
		fmt.Println("world - Get world informaion")
	default:
		fmt.Println("Unknown command")
	}
}

// Query context
func contextQuery(cmd []string, client *thg.Thg) {
	switch cmd[0] {
	case "rq", "sq":
		url := strings.Join(cmd[1:], "&")
		if cmd[0] == "sq" {
			if len(client.SessionId) == 0 {
				fmt.Println("No session ID")
				return
			}
			url += "&session_id=" + client.SessionId
		}
		url = client.GetFullUrl(url)
		fmt.Println("Raw query: " + url)
		res, err := client.HttpGet(url)
		fmt.Println()
		if err != nil {
			// log.Println(err)
			fmt.Println("Error")
			return
		}

		switch queryPm {
		case "none":
			fmt.Println(res)
		case "sc":
			data := thg.ParseDataSc(res)
			for i1, _ := range data {
				fmt.Printf("== [%d] =============\n", i1)
				for _, v2 := range data[i1] {
					fmt.Printf("%s ", v2)
				}
				fmt.Printf("\n\n")
			}
		case "dog":
			data := thg.ParseDataDog(res)
			for i1, _ := range data {
				fmt.Printf("== [%d] =============\n", i1)
				for i2, _ := range data[i1] {
					for _, v3 := range data[i1][i2] {
						fmt.Printf("[%s] ", v3)
						}
					fmt.Println()
				}
			fmt.Println()
			}
		}
	case "pm":
		if len(cmd) == 2 {
			if cmd[1] != "none" && cmd[1] != "sc" &&
				cmd[1] != "dog" && cmd[1] != "bin" {
				fmt.Println("Unknown parse mode")
			} else {
				queryPm = cmd[1]
			}
			return
		}

		fmt.Println("Parse mode:", queryPm)
	case "help":
		fmt.Println("rq <arg1> .. <argN> - Raw query")
		fmt.Println("sq <arg1> .. <argN> - Raw query with session ID")
		fmt.Println("pm <none|sc|dog|bin> - Parse mode")
	default:
		fmt.Println("Unknown command")
	}
}

// Script context
func contextScript(cmd []string, vm *otto.Otto) {
	var script string

	switch cmd[0] {
	case "run", "go":
		if len(cmd) >= 2 {
			f, err := os.Open(scriptsDir + "/" + cmd[1])
			if err != nil {
				log.Println(err)
				return
			}
			defer f.Close()
			buff := bytes.NewBuffer(nil)
			if _, err := buff.ReadFrom(f); err != nil {
				log.Println(err)
				return
			}
			script = buff.String()
			vm.Set("args", cmd[2:])
		}
	case "eval":
		script = strings.Join(cmd[1:], " ")
	case "help":
		fmt.Println("run <file> - Run script from file")
		fmt.Println("go <file> - Run script from file in background")
		fmt.Println("eval <script> - Execute script code")
	default:
		fmt.Println("Unknown command")
	}

	if len(script) == 0 {
		return
	}

	switch cmd[0] {
	case "run", "eval":
		if _, err := vm.Run(script); err != nil {
			log.Println(err)
		}
	case "go":
		go func(vm *otto.Otto, script string) {
			if _, err := vm.Run(script); err != nil {
				log.Println(err)
				return
			}
		}(vm, script)
	}
}

// Chat context
func contextChat(cmd []string, client *thg.Thg, channels map[uint16]chan string) {
	switch cmd[0] {
	case "open":
		if len(client.SessionId) == 0 {
			fmt.Println("No session ID")
			return
		}
		if len(cmd) != 2 {
			return
		}
		var room uint16
		if v, e := strconv.Atoi(cmd[1]); e == nil {
			room = uint16(v)
		}

		channels[room] = make(chan string)

		go func(room uint16, channels map[uint16]chan string) {
			lm := ""
			for {
				data, err := client.ChatDisplay(room, lm)
				if err != nil {
					log.Println(err)
					close(channels[room])
					delete(channels, room)
					return
				}
				for i := range data.Messages {
					fmt.Fprintf(rl, "\x1b[33m[%s] \x1b[35m%s (%d): \x1b[0m%s\n",
						data.Messages[i].DateTime, data.Messages[i].Name, data.Messages[i].Id, data.Messages[i].Message)
					lm = data.Messages[i].DateTime
				}
				select {
				case message := <- channels[room]:
					if message == "close" {
						return
					}
				case <- time.After(time.Second * 5):
				}
			}
		}(room, channels)
	case "close":
		if len(cmd) != 2 {
			return
		}
		var room uint16
		if v, e := strconv.Atoi(cmd[1]); e == nil {
			room = uint16(v)
		}
		if channel, ok := channels[room]; ok {
			channel <- "close"
			close(channel)
			delete(channels, room)
		}
	case "list":
		for i := range channels {
			fmt.Fprintln(rl, i)
		}
	case "say":
		if len(client.SessionId) == 0 {
			fmt.Println("No session ID")
			return
		}
		if len(cmd) < 3 {
			return
		}
		var room uint16
		if v, e := strconv.Atoi(cmd[1]); e == nil {
			room = uint16(v)
		}
		_, err := client.ChatSend(room, strings.Join(cmd[2:], " "), "")
		if err != nil {
			log.Println(err)
			return
		}
	case "help":
		fmt.Println("open <room> - Open chat")
		fmt.Println("close <room> - Close chat")
		fmt.Println("list - List opened chats")
		fmt.Println("say <room> <message> - Send message to chat")
	default:
		fmt.Println("Unknown command")
	}
}
