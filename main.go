package main

import (
	"fmt"
	"os/exec"
	"sort"
	"sync"
	"time"

	"github.com/spf13/viper"
)

type ScheduleConfig struct {
	Link           string
	Start_time     string
	Duration       int
	start_time_obj time.Time
}

func (config *ScheduleConfig) ToString() string {
	return fmt.Sprintf("Link: %s, Start at: %s, Duration: %d min (%s)", config.Link, config.Start_time, config.Duration, config.start_time_obj)
}

func (config *ScheduleConfig) ToHumanString() string {
	return fmt.Sprintf("Start at: %s, in %s for %d minutes", config.Start_time, config.Link, config.Duration)
}

type Config struct {
	Version   int              `mapstructure:"last_updated"`
	Task_list []ScheduleConfig `mapstructure:"task"`
	//Task []map[string]string
}

type PendingList struct {
	// Threa Safe desgin of pending tasks
	mux     sync.Mutex
	pending []ScheduleConfig
}

// =================================

func print_pending_task(verified_tasks []ScheduleConfig) {
	fmt.Println("No. of Task are waiting:", len(verified_tasks))
	for i, task := range verified_tasks {
		fmt.Printf("%2d : %s\n", i+1, task.ToHumanString())
	}
}

func main() {
	viper.SetConfigFile("schedule.yaml")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {
		// Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		fmt.Println("Error when unmarshalling config file: ", err)
	}
	fmt.Println(config.Version)

	// Parse the start time to Time Object
	var verified_tasks []ScheduleConfig

	for _, task := range config.Task_list {
		task.start_time_obj, err = time.ParseInLocation("2006-01-02 15:04:05", task.Start_time, time.Local)
		if err != nil {
			fmt.Println("[WARN] Ini Date format not correct", err)
			continue
		}

		// Notify the outdated task
		if task.start_time_obj.Before(time.Now()) {
			fmt.Println("Task (Start at", task.Start_time, "Outdated.")
			continue
		}

		verified_tasks = append(verified_tasks, task)
		//fmt.Println(task.ToString())
	}

	// Sort the object due to time order
	sort.Slice(verified_tasks, func(i, j int) bool {
		if verified_tasks[i].start_time_obj.Before(verified_tasks[j].start_time_obj) {
			return true
		}
		if verified_tasks[i].start_time_obj.Equal(verified_tasks[j].start_time_obj) && verified_tasks[i].Duration < verified_tasks[j].Duration {
			return true
		}
		return false
	})

	// Prepare UI for start program.
	print_pending_task(verified_tasks)
	fmt.Printf("\n\n") // Seperation

	// Prepare worker pool
	concurrentGoroutines := make(chan struct{}, 5)
	var wg sync.WaitGroup

	// Create pending tasks object
	var pending_list PendingList
	pending_list.pending = verified_tasks[0:]

	// Start Main loop
	fmt.Println("Main Function Starting...")
	for i, item := range verified_tasks {
		// Wait Group Increment
		wg.Add(1)

		go func(item ScheduleConfig, i int) {
			// Defer Done() for WaitGroup
			defer wg.Done()
			concurrentGoroutines <- struct{}{}

			// Sleep until time come
			time.Sleep(time.Until(item.start_time_obj))

			// Time matched, start the batch script with parameters
			fmt.Println("Task Started:", item.Link, "for", item.Duration, "minutes. End at", item.start_time_obj.Add(time.Duration(10)*time.Minute).Format("2006-01-02 15:04:05"))
			fmt.Println()

			// cmd1 := exec.Command("cmd.exe", "/C", "start", "https://www.twitch.tv/drops/inventory")
			// cmd1.Start()
			// time.Sleep(time.Duration(10) * time.Second)

			cmd2 := exec.Command("cmd.exe", "/C", "start", item.Link)
			cmd2.Start()

			// Wait the duration
			time.Sleep(time.Duration(item.Duration) * time.Minute)

			// Update pending list
			pending_list.mux.Lock()
			fmt.Println("Task for", item.Link, "Completed.")
			pending_list.pending = pending_list.pending[1:]
			print_pending_task(pending_list.pending)
			fmt.Printf("\n\n")
			pending_list.mux.Unlock()

			<-concurrentGoroutines
		}(item, i)
	}

	wg.Wait()
	fmt.Println("All Task is completed.")

	// // Notify No Task is pending
	// fmt.Println("All Task is completed.")

	// // // Fetch The Ini text setting from web
	// // source_ini_link := "https://raw.githubusercontent.com/dark-person/Twitch-drop-custom/master/schedule.ini"
	// // fileURL, err := url.Parse(source_ini_link)
	// // if err != nil {
	// // 	fmt.Println("Download INI Failed:", err)
	// // }
	// // path := fileURL.Path
	// // segments := strings.Split(path, "/")
	// // fileName := "TEMP_" + segments[len(segments)-1]

	// // // Create blank file
	// // file, err := os.Create(fileName)
	// // if err != nil {
	// // 	fmt.Println("Create Ini Failed:", err)
	// // }
	// // client := http.Client{
	// // 	CheckRedirect: func(r *http.Request, via []*http.Request) error {
	// // 		r.URL.Opaque = r.URL.Path
	// // 		return nil
	// // 	},
	// // }
	// // // Put content on file
	// // resp, err := client.Get(source_ini_link)
	// // if err != nil {
	// // 	fmt.Println("Download INI Failed:", err)
	// // }
	// // defer resp.Body.Close()

	// // size, _ := io.Copy(file, resp.Body)
	// // defer file.Close()

	// // fmt.Printf("Downloaded ini %s with size %d\n", fileName, size)

	// // // Verify the downloaded file is newer than current	version
	// // // Read the config from external files
	// // new_config, err := ini.Load("TEMP_schedule.ini")
	// // if err != nil {
	// // 	fmt.Printf("Fail to read file: %v", err)
	// // 	os.Exit(1)
	// // }
	// // new_config.SaveTo("TEMP_schedule.ini")

	// // // Get Task Number
	// // new_version := new_config.Section("").Key("last_updated").MustInt()
	// // fmt.Println("Downloaded Version:", new_version)
	// // if new_version > version_number {
	// // 	os.Rename("TEMP_schedule.ini", "schedule.ini")
	// // 	fmt.Println("Newer Version replaced.")
	// // } else {
	// // 	err := os.Remove("TEMP_schedule.ini")
	// // 	if err != nil {
	// // 		fmt.Printf("Remove File failed: %v\n", err)
	// // 	}

	// // 	fmt.Println("Downloaded Version Removed.")
	// // }
}
