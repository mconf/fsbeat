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

type FSConnection struct {
  c *Connection
  address string
  server string
  port string
  auth string
}

type IFSConnection interface {
  getConnection() (*Connection, error)
  configureConnection()
  closeConnection(events chan *Event)
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

  // Creates a new FreeSWITCH connection struct and initializes it.
  fsc, err := createFSConnection(bt.config.FSServer,
                                 bt.config.FSPort,
                                 bt.config.FSAuth)

  if err != nil {
    logp.Err("Error while connecting to FreeSWITCH (%s).", fsc.address)

    return nil
  }

  // Configures the new FreeSWITCH connection.
  fsc.configureConnection(bt.config.FSEvents)

  // Creates the channel used to communicate events received from FreeSWITCH.
  events := make(chan *Event, bt.config.MaxBuffer)

  // Launches goroutine to received events from FreeSWITCH.
  go fsc.getEvents(events)

  var ev *Event
  for {
    select {
    case <-bt.done:
      fsc.closeConnection(events)

      return nil
    case ev = <-events:
    }

    // TODO: Change type to 'esl' or something like that?
    event := common.MapStr{
      "@timestamp": common.Time(time.Now()),
      "type":       b.Name,
    }

    // Copies all fields from ev to event.
    for k, v := range ev.Header {
      k = normalize_field(k)
      event[k] = v
    }

    bt.client.PublishEvent(event)
  }
}

func (bt *Fsbeat) Stop() {
  bt.client.Close()
  close(bt.done)
}

func createFSConnection(FSServer string,
                        FSPort string,
                        FSAuth string) (*FSConnection, error) {
  s := []string{FSServer, FSPort}
  address := strings.Join(s, ":")

  fsc := &FSConnection{address: address,
                       server: FSServer,
                       port: FSPort,
                       auth: FSAuth}

  c, err := fsc.getConnection()
  if err == nil {
    fsc.c = c
  }

  return fsc, err
}

func (fsc *FSConnection) getConnection() (*Connection, error) {
  logp.Info("Trying to connect to FreeSWITCH (%s).", fsc.address)

  c, err := Dial(fsc.address, fsc.auth)

  return c, err
}

func (fsc *FSConnection) closeConnection(events chan *Event) {
  logp.Info("Closing connection to FreeSWITCH (%s) and channels.", fsc.address)

  fsc.c.Close()
  close(events)
}

func (fsc *FSConnection) configureConnection(FSEvents string) {
  logp.Info("Configuring connection to FreeSWITCH (%s).", fsc.address)

  const dest = "sofia/internal/1000%127.0.0.1"
  const dialplan = "&socket(localhost:9090 async)"

  if fsc.c == nil {
    fmt.Println("fsc.c is nil.")
  }

  fsc.c.Send("events json " + FSEvents)
  fsc.c.Send(fmt.Sprintf("bgapi originate %s %s", dest, dialplan))
}

func (fsc *FSConnection) getEvents(events chan<- *Event) {
  c := fsc.c

  for {
    ev, err := c.ReadEvent()
    if err != nil {
      logp.Err("Error while reading event.")
    }

    events <- ev
  }
}

// It follows the conventions defined by Beat:
// https://www.elastic.co/guide/en/beats/libbeat/current/event-conventions.html
func normalize_field(field string) string {
      field = strings.ToLower(field)
      field = strings.Replace(field, "-", "_", -1)

      return field
}
