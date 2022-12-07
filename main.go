package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"sync"
	"time"

	"gopkg.in/ini.v1"
)

type schedule struct {
	task_id     int
	start_time  time.Time
	twitch_link string
	duration    int // in minutes
}

// =================================

func main() {
	//var twitch_link string
	//var duration_minute int
	task_list := make([]schedule, 0)

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
	// version_number := config.Section("").Key("last_updated").MustInt()

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

		// Parse the config to object
		task_list = append(task_list, schedule{task_id: i + 1, twitch_link: twitch_link, duration: duration_minute, start_time: until})
	}

	// Sort the object due to time order
	sort.Slice(task_list, func(i, j int) bool {

		if task_list[i].start_time.Before(task_list[j].start_time) {
			return true
		}

		if task_list[i].start_time.Equal(task_list[j].start_time) && task_list[i].task_id < task_list[j].task_id {
			return true
		}

		return false
	})

	fmt.Println(task_list)

	// Prepare worker pool
	concurrentGoroutines := make(chan struct{}, 5)

	var wg sync.WaitGroup

	// Start Main loop
	for _, item := range task_list {
		// Wait Group Increment
		wg.Add(1)

		go func(item schedule) {
			// Defer Done() for WaitGroup
			defer wg.Done()
			concurrentGoroutines <- struct{}{}

			// Notify the outdated task
			if item.start_time.Before(time.Now()) {
				fmt.Println("Task", item.task_id, "Outdated.")
				<-concurrentGoroutines
				return
			}

			// Sleep until time come
			// fmt.Println("Sleeping...")
			time.Sleep(time.Until(item.start_time))

			// Time matched, start the batch script with parameters
			cmd1 := exec.Command("cmd.exe", "/C", "start", "https://www.twitch.tv/drops/inventory")
			cmd1.Start()
			time.Sleep(time.Duration(10) * time.Second)

			cmd2 := exec.Command("cmd.exe", "/C", "start", item.twitch_link)
			cmd2.Start()

			// Wait the duration
			// fmt.Println("Sleeping...")
			time.Sleep(time.Duration(item.duration) * time.Minute)

			<-concurrentGoroutines
		}(item)
	}

	wg.Wait()

	// Notify No Task is pending
	fmt.Println("All Task is completed.")

	// // Fetch The Ini text setting from web
	// source_ini_link := "https://raw.githubusercontent.com/dark-person/Twitch-drop-custom/master/schedule.ini"
	// fileURL, err := url.Parse(source_ini_link)
	// if err != nil {
	// 	fmt.Println("Download INI Failed:", err)
	// }
	// path := fileURL.Path
	// segments := strings.Split(path, "/")
	// fileName := "TEMP_" + segments[len(segments)-1]

	// // Create blank file
	// file, err := os.Create(fileName)
	// if err != nil {
	// 	fmt.Println("Create Ini Failed:", err)
	// }
	// client := http.Client{
	// 	CheckRedirect: func(r *http.Request, via []*http.Request) error {
	// 		r.URL.Opaque = r.URL.Path
	// 		return nil
	// 	},
	// }
	// // Put content on file
	// resp, err := client.Get(source_ini_link)
	// if err != nil {
	// 	fmt.Println("Download INI Failed:", err)
	// }
	// defer resp.Body.Close()

	// size, _ := io.Copy(file, resp.Body)
	// defer file.Close()

	// fmt.Printf("Downloaded ini %s with size %d\n", fileName, size)

	// // Verify the downloaded file is newer than current	version
	// // Read the config from external files
	// new_config, err := ini.Load("TEMP_schedule.ini")
	// if err != nil {
	// 	fmt.Printf("Fail to read file: %v", err)
	// 	os.Exit(1)
	// }
	// new_config.SaveTo("TEMP_schedule.ini")

	// // Get Task Number
	// new_version := new_config.Section("").Key("last_updated").MustInt()
	// fmt.Println("Downloaded Version:", new_version)
	// if new_version > version_number {
	// 	os.Rename("TEMP_schedule.ini", "schedule.ini")
	// 	fmt.Println("Newer Version replaced.")
	// } else {
	// 	err := os.Remove("TEMP_schedule.ini")
	// 	if err != nil {
	// 		fmt.Printf("Remove File failed: %v\n", err)
	// 	}

	// 	fmt.Println("Downloaded Version Removed.")
	// }
}
