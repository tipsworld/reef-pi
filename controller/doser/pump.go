package doser

import (
	"encoding/json"
	"github.com/reef-pi/reef-pi/controller/utils"
	"gopkg.in/robfig/cron.v2"
	"log"
)

func (c *Controller) Get(id string) (Pump, error) {
	var p Pump
	return p, c.store.Get(Bucket, id, &p)
}

func (c *Controller) Create(p Pump) error {
	fn := func(id string) interface{} {
		p.ID = id
		return &p
	}
	return c.store.Create(Bucket, fn)
}

func (c *Controller) List() ([]Pump, error) {
	pumps := []Pump{}
	fn := func(v []byte) error {
		var p Pump
		if err := json.Unmarshal(v, &p); err != nil {
			return err
		}
		pumps = append(pumps, p)
		return nil
	}
	return pumps, c.store.List(Bucket, fn)
}

func (c *Controller) Calibrate(id string, cal CalibrationDetails) error {
	p, err := c.Get(id)
	if err != nil {
		return err
	}
	r := &Runner{
		pin:      p.Pin,
		duration: cal.Duration,
		speed:    cal.Speed,
		vv:       c.vv,
	}
	go r.Run()
	return nil
}

func (c *Controller) Update(id string, p Pump) error {
	p.ID = id
	if err := c.store.Update(Bucket, id, p); err != nil {
		return err
	}
	// TODO cross check cron assignment
	return nil
}

func (c *Controller) Schedule(id string, r DosingRegiment) error {
	log.Println(r)
	if err := r.Schedule.Validate(); err != nil {
		log.Printf("CronSpec:'%s'\n", r.Schedule.CronSpec())
		return err
	}
	p, err := c.Get(id)
	if err != nil {
		return err
	}
	p.Regiment = r
	if err := c.Update(id, p); err != nil {
		return err
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if cID, ok := c.cronIDs[id]; ok {
		log.Printf("doser sub-system. Removing cron entry %s for pump id: %s.\n", cID, id)
		c.runner.Remove(cID)
	}
	if p.Regiment.Enable {
		cronID, err := c.runner.AddJob(p.Regiment.Schedule.CronSpec(), p.Runner(c.vv))
		if err != nil {
			return err
		}
		log.Println("Successfully added cron entry. ID:", cronID)
		c.cronIDs[p.ID] = cronID
	}
	return nil
}

func (c *Controller) Delete(id string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if cID, ok := c.cronIDs[id]; ok {
		log.Printf("doser sub-system. Removing cron entry %s for pump id: %s.\n", cID, id)
		c.runner.Remove(cID)
	}
	return c.store.Delete(Bucket, id)
}

func (p *Pump) Runner(vv utils.VariableVoltage) cron.Job {
	return &Runner{
		pin:      p.Pin,
		duration: p.Regiment.Duration,
		speed:    p.Regiment.Speed,
		vv:       vv,
	}
}