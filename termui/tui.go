package termui

// func TermuiInit() {
// 	// time.Sleep(2 * time.Second)
// 	err := termbox.Init()
// 	if err != nil {
// 		log.Logger.Errorw("termui初始化失败", "err", err)
// 		return
// 	}
// 	homeUi()
// }

// func homeUi() {
// 	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
// 	termbox.Flush()
// 	list := process.ProcessCtlService.GetProcessList()
// 	fmt.Println()
// 	for i, v := range list {
// 		if v.User != "" {
// 			fmt.Printf("  [%v] %v  %v  <%v>\n", i, v.Name, v.StartTime, v.User)
// 		} else {
// 			fmt.Printf("  [%v] %v  %v\n", i, v.Name, v.StartTime)
// 		}
// 	}
// 	input := ""
// 	fmt.Scan(&input)
// 	for i, v := range list {
// 		if input == strconv.Itoa(i) {
// 			// prcessUi(v.Uuid)
// 		}
// 	}
// }

// func prcessUi(uuid int) {
// 	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
// 	termbox.Flush()
// 	proc, err := process.ProcessCtlService.GetProcess(uuid)
// 	if err != nil {
// 		log.Logger.Errorw("进程获取失败", "err", err)
// 		return
// 	}
// 	proc.SetControl("")
// 	go func() {
// 		for {
// 			if output := proc.Read(); output != "" {
// 				fmt.Println(output)
// 			}
// 		}
// 	}()
// 	go func() {
// 		input := ""
// 		for {
// 			fmt.Scan(&input)
// 			proc.Write(input + "\n")
// 		}
// 	}()

// }
