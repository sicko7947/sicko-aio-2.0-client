package communicator

import (
	"fmt"
	"os"
	"time"

	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/os/gmutex"
	"github.com/huandu/go-clone"
	"github.com/sicko7947/sickocommon"
	"sicko-aio-2.0-client/models"
)

var (
	// DEV_ENV : check if its under dev enviroment
	DEV_ENV = true
	// LOGGING : enable logging or not
	LOGGING = true
	// PATH : path for storing config.json file
	PATH = os.Getenv("APPDATA") + "/sicko-aio"
	// Config : global config files that contains all user
	Config *models.Config
	// TaskLogsGMap : global variable that contains all the current task status and message
	TaskLogsGMap *gmap.AnyAnyMap
	// TaskWorkerObjectGMap : global variable that contains all the current task information
	TaskWorkerObjectGMap *gmap.AnyAnyMap
	// TaskScraperObjectGMap : global variable that contains all the current task information
	TaskScraperObjectGMap *gmap.AnyAnyMap
	// TaskMessageGMap : global variable that contains all the current task logs
	TaskMessageGMap *gmap.AnyAnyMap
	// TempMessageGmap : temporary map that contains all the current task logs
	TempMessageGmap *gmap.AnyAnyMap
	// SuccessCountGmap : record success count by each taskgroup
	SuccessCountGmap *gmap.StrIntMap
)

func init() {
	TaskLogsGMap = gmap.New(true)
	TaskWorkerObjectGMap = gmap.New(true)
	TaskScraperObjectGMap = gmap.New(true)
	TaskMessageGMap = gmap.New(true)
	TempMessageGmap = gmap.New(true)
	SuccessCountGmap = gmap.NewStrIntMap(true)

	Config = &models.Config{ // Initiate config object at begining
		TaskGroups: make(map[models.TaskGroupID]*models.TaskGroup),
		Profiles:   make(map[models.ProfileGroupName]map[models.ProfileName]*models.Profile),
		Proxies:    make(map[models.ProxyGroupName][]string),
		Accounts:   make(map[models.AccountGroupName][]*models.Account),
		GiftCards:  make(map[models.GiftCardGroupName][]*models.GiftCard),
		Settings:   new(models.Settings),
	}

	sickocommon.CopyFile(PATH+"/config.json", PATH+"/config-copy.json", true) // backup confifg on start up
	err := sickocommon.ReadJson(PATH+"/config.json", &Config)
	if err != nil {
		fmt.Println("Error loading config file - will overide it with default value")
		sickocommon.WriteJson(PATH+"/config.json", Config, os.ModePerm)
	}

	for _, v := range Config.TaskGroups {
		t := clone.Clone(v).(*models.TaskGroup) // deep copy value and put in TaskObjMap
		for _, worker := range t.TaskWorkers {
			worker.Mutex = gmutex.New()
			TaskWorkerObjectGMap.Set(worker.TaskID, worker)
		}

		for _, scraper := range t.TaskScrapers {
			scraper.Mutex = gmutex.New()
			TaskScraperObjectGMap.Set(scraper.TaskID, scraper)
		}

		SuccessCountGmap.Set(string(v.GroupID), 0)
	}

	go func() {
		ticker := time.NewTicker(3 * time.Second)
		for {
			sickocommon.WriteJson(PATH+"/config.json", Config, os.ModePerm)
			<-ticker.C
		}
	}()
}
