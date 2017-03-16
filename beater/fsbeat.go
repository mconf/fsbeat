package beater

import (
  "fmt"
  "time"
  "strings"

  "github.com/elastic/beats/libbeat/beat"
  "github.com/elastic/beats/libbeat/common"
  "github.com/elastic/beats/libbeat/logp"
  "github.com/elastic/beats/libbeat/publisher"

  "github.com/mconftec/fsbeat/config"
)

type Fsbeat struct {
  done   chan struct{}
  config config.Config
  client publisher.Client
}

// Creates beater
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
  config := config.DefaultConfig
  if err := cfg.Unpack(&config); err != nil {
    return nil, fmt.Errorf("Error reading config file: %v", err)
  }

  bt := &Fsbeat{
    done: make(chan struct{}),
    config: config,
  }
  return bt, nil
}

func (bt *Fsbeat) Run(b *beat.Beat) error {
  logp.Info("fsbeat is running! Hit CTRL-C to stop it.")

  bt.client = b.Publisher.Connect()
  //ticker := time.NewTicker(bt.config.Period)

  c, err := getConnection(bt.config.FSServer,
                          bt.config.FSPort,
                          bt.config.FSAuth)

  if err != nil {
    logp.Err("Error getting connection to FreeSWITCH server.")
    return nil
  }

  c.configure()

  events := make(chan *Event, bt.config.MaxBuffer)

  go c.getEvents(events)

  var ev *Event
  for {
    select {
    case <-bt.done:
      // TODO: Change to defer?
      c.Close()
      return nil
    case ev = <-events:
    // TODO: Do we need a case for ticker? I don't think so.
    //case <-ticker.C:
    }

    // TODO: Change type to 'esl' or something like that?
    event := common.MapStr{
      "@timestamp": common.Time(time.Now()),
      "type":       b.Name,
    }

    // Copy all fields from ev to event.
    for k, v := range ev.Header {
      event[k] = v
    }

    // TODO: Remove it.
    fmt.Println(event)
    bt.client.PublishEvent(event)
    logp.Info("Event sent")
  }
}

func (bt *Fsbeat) Stop() {
  bt.client.Close()
  close(bt.done)
}

func getConnection(FSServer string,
                   FSPort string,
                   FSAuth string) (*Connection, error) {
  s := []string{FSServer, FSPort}
  address := strings.Join(s, ":")

  c, err := Dial(address, FSAuth)

  return c, err
}

func (c *Connection) configure() {
  const dest = "sofia/internal/1000%127.0.0.1"
  const dialplan = "&socket(localhost:9090 async)"

  c.Send("events json ALL")
  c.Send(fmt.Sprintf("bgapi originate %s %s", dest, dialplan))
}

func (c *Connection) getEvents(events chan<- *Event) {
  for {
    ev, err := c.ReadEvent()
    if err != nil {
      logp.Err("Error reading event.")
    }

    // TODO: Is this really necessary?
    if ev.Get("Answer-State") == "hangup" {
      break
    }

    events <- ev
  }
}
