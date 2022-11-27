package main

import (
  twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
  "github.com/pelletier/go-toml/v2"
  "github.com/twilio/twilio-go"
  "math/rand"
  "encoding/json"
  "time"
  "log"
  "fmt"
  "os"
)

var config string

type Person struct {
  Name  string
  Phone string
}

type Details struct {
  Message  string
  Greeting string
  TwilioSid string
  TwilioToken string
  TwilioNumber string
}

type Config struct {
  People []Person
  Deets  Details
}

type Match struct {
  Person Person
  Match  Person
}

type SMSResult struct {
  Match Match
  Error error
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

func sendSMS(client *twilio.RestClient, config Config, match Match, resChan chan SMSResult) {
  params := &twilioApi.CreateMessageParams{}
  params.SetTo(match.Person.Phone)
  params.SetFrom(config.Deets.TwilioNumber)
  msg := fmt.Sprintf(
    "%s %s. %s %s.",
    config.Deets.Greeting,
    match.Person.Name,
    config.Deets.Message,
    match.Match.Name)
  params.SetBody(msg)

  result := SMSResult{match, nil}
  resp, err := client.Api.CreateMessage(params)
  if err != nil {
    result.Error = err
    resChan <- result
  } else {
    response, _ := json.Marshal(*resp)
    log.Println("Response:", response)
    resChan <- result
  }
}

func tryMatch(people []Person) []Match {
  n := len(people)
  list := rand.Perm(n)
  matches := make([]Match, n)
  for i := 0; i < n; i++ {
    currentPerson := people[i]
    matchedPerson := people[list[i]]
    if i == list[i] {
      log.Println(
        "Unfortunately,",
        currentPerson.Name,
        "got", currentPerson.Name)
      return nil
    }
    if list[i] < i && matches[list[i]].Match.Name == currentPerson.Name {
      log.Println("Unfortunately,",
        currentPerson.Name,
        "got", matchedPerson.Name,
        "but", matchedPerson.Name,
        "matched with", currentPerson.Name,
        "which is a cycle and is not allowed")
      return nil
    }
    matches[i] = Match{Person: currentPerson, Match: matchedPerson}
  }
  return matches
}

func getTwilioClient(config Config) *twilio.RestClient {
  return twilio.NewRestClientWithParams(twilio.ClientParams{
    Username: config.Deets.TwilioSid,
    Password: config.Deets.TwilioToken,
  })
}

func saveFailures(failures []SMSResult) {
  if len(failures) > 0 {
    data, _ := json.Marshal(failures)
    err := os.WriteFile("failures.txt", []byte(data), 0664)
    if err != nil {
      log.Fatal("oh my word")
    }
    log.Fatal("Awful.")
  } else {
    log.Println("Succeeded to an extreme degree.")
  }
}

func main() {
  rand.Seed(time.Now().UnixNano())
  config := parse_config()
  var mixed []Match

  if len(config.People) <= 2 {
    log.Fatal("Big Epic failure. Need more people.")
  }

  for {
    mixed = tryMatch(config.People)
    if mixed != nil {
      break
    }
  }

  resultChan := make(chan SMSResult)

  client := getTwilioClient(config)
  for _, v := range mixed {
    go sendSMS(client, config, v, resultChan)
  }

  length := len(mixed)
  failures := make([]SMSResult, 0)
  completeFailure := true
  for i := 0; i < length; i++ {
    if result := <-resultChan; result.Error != nil {
      log.Println("Failed Sending SMS to", result.Match.Person.Name)
      failures = append(failures, result)
    } else {
      log.Println("Succeeded Sending SMS to", result.Match.Person.Name)
      completeFailure = false
    }
  }

  if !completeFailure {
    saveFailures(failures)
  } else {
    log.Println("CompleteFailure")
    log.Println(failures)
  }
}
