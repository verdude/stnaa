package main

import (
  "github.com/pelletier/go-toml/v2"
  "math/rand"
  "time"
  "log"
  "os"
)

var config string

type Person struct {
  Name string
  Phone string
}

type Details struct {
  Message string
  Greeting string
}

type Config struct {
  People []Person
  Deets Details
}

type Match struct {
  Person Person
  Match Person
}

func parse_config() Config {
  bytes, err := os.ReadFile(config)
  if err != nil {
    log.Fatal("Failed to read config: ", err)
  }
  var conf Config
  e := toml.Unmarshal(bytes, &conf)
  if e != nil {
    log.Fatal("Failed to parse config into toml: ", e)
  }
  return conf
}

func tryMatch(people []Person) []Match {
  n := len(people)
  list := rand.Perm(n)
  matches := make([]Match, n)
  for i := 0; i < n; i++ {
    if i == list[i] {
      log.Println("Unfortunately,", people[i].Name, "got", people[i].Name)
      return nil
    }
    if list[i] < i && matches[list[i]].Match.Name == people[i].Name {
      log.Println("Unfortunately,", people[i].Name, "got", people[list[i]].Name, "but", people[list[i]].Name, "matched with", people[i].Name, "which is a cycle and is not allowed")
      return nil
    }
    matches[i] = Match{ Person: people[i], Match: people[list[i]] }
  }
  return matches
}

func main() {
  rand.Seed(time.Now().UnixNano())
  config := parse_config()
  var mixed []Match

  for {
    mixed = tryMatch(config.People)
    if mixed != nil {
      break
    }
    log.Println("trying again...")
  }

  log.Println(mixed)
}
