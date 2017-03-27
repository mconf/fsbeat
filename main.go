package main

import (
  "os"

  "github.com/elastic/beats/libbeat/beat"

  "github.com/mconftec/fsbeat/beater"
)

func main() {
  err := beat.Run("fsbeat", "", beater.New)
  if err != nil {
    os.Exit(1)
  }
}
