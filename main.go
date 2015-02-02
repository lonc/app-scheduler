package main

import (
	"fmt"
	"github.com/ninjasphere/app-scheduler/model"
	"github.com/ninjasphere/app-scheduler/service"
	"github.com/ninjasphere/go-ninja/api"
	"github.com/ninjasphere/go-ninja/support"
)

var (
	info = ninja.LoadModuleInfo("./package.json")
)

//SchedulerApp describes the scheduler application.
type SchedulerApp struct {
	support.AppSupport
	service *service.SchedulerService
}

// Start is called after the ExportApp call is complete.
func (a *SchedulerApp) Start(m *model.Schedule) error {
	if a.service != nil {
		return fmt.Errorf("illegal state: scheduler is already running")
	}
	a.service = &service.SchedulerService{
		Log:   a.Log,
		Conn:  a.Conn,
		Model: m,
		ConfigStore: func(m *model.Schedule) {
			a.SendEvent("config", m)
		},
	}
	err := a.service.Init(a.Info.ID)
	if err != nil {
		return err
	}
	return nil
}

// Stop the scheduler module.
func (a *SchedulerApp) Stop() error {
	var err error
	if a.service != nil {
		tmp := a.service
		a.service = nil
		err = tmp.Scheduler.Stop()
	}
	return err
}

func main() {
	app := &SchedulerApp{}
	err := app.Init(info)
	if err != nil {
		app.Log.Fatalf("failed to initialize app: %v", err)
	}

	err = app.Export(app)
	if err != nil {
		app.Log.Fatalf("failed to export app: %v", err)
	}

	support.WaitUntilSignal()
}
