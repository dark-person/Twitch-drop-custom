package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"gopkg.in/ini.v1"
)

func main() {
	//var twitch_link string
	//var duration_minute int

	// flag.StringVar(&twitch_link, "link", "https://www.twitch.tv/warframe", "[Required] The link for twitch drop")
	// flag.IntVar(&duration_minute, "time", 60, "[Required] The duration for the stream last. Unit is minute. Default is 60 minutes.")
	flag.Parse()

	// Read the config from external files
	config, err := ini.Load("schedule.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	// Get Task Number
	task_count := config.Section("").Key("task_count").MustInt()

	for i := 0; i < task_count; i++ {
		section := strconv.Itoa(i + 1)
		twitch_link := config.Section(section).Key("link").String()
		duration_minute := config.Section(section).Key("duration").MustInt()
		start_time := config.Section(section).Key("start_time").String()

		fmt.Println("Twitch Link:", twitch_link)
		fmt.Println("Duration:", duration_minute, "minute")
		fmt.Println("Start at: ", start_time)
		fmt.Println()

		// Parse the start time to Time Object
		//until, err := time.Parse("2006-01-02 15:04:05 +0800 CST", start_time+" +0800 CST")
		until, err := time.ParseInLocation("2006-01-02 15:04:05", start_time, time.Local)
		if err != nil {
			fmt.Println("Ini Date format not correct", err)
			os.Exit(1)
		}

		fmt.Println(until.String())
		fmt.Println(time.Now().String())

		// fmt.Println(until.Equal(time.Now()))
		// fmt.Println(until.Before(time.Now()))
		// fmt.Println(until.After(time.Now()))

		if until.Before(time.Now()) {
			fmt.Println("Task Outdated.")
			os.Exit(1)
		}

		time.Sleep(time.Until(until))
		fmt.Println("Task Completed.")

		//exec.Command("cmd.exe", "/C", "twitch_drop_general.bat", twitch_link, strconv.Itoa(duration_minute))
		break
	}

	// Parse the config to object

	// Sort the object due to time order

	// Notify the outdated task

	// Start Main loop

	// Sleep until time come

	// Time matched, start the batch script with parameters

	// Wait the batch script stop

	// Prepare the next task (loop counter increment)

	// Notify No Task is pending

	// Fetch The Ini text setting from web?

	//fmt.Println("Twitch Link:", twitch_link)
	//fmt.Println("Duration:", duration_minute, "minute")

	// Prepare to start the bat
}
