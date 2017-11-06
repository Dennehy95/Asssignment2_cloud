package main

// Imported resources
import (
	"fmt"
	"github.com/yJepo/Asssignment2_cloud/ass2"
	"strings"
	"time"
)

func main() {
	TimedTicker(time.Hour * 24)
}

// Automatically checks for data after a set time. Takes time.Duration as input.
func TimedTicker(tickerTime time.Duration) error { //Errors should be: log and break, not return to not kill the process?
	ticker := time.NewTicker(tickerTime)

	for range ticker.C {
		// Starts a session with the DB
		c, session, err := resource.StartSession(resource.Url, resource.Database, resource.Collection2)
		if err != nil {
			return err
		}

		// Load data from the DB
		result := &resource.Rates{}
		err = c.Find(nil).Sort("-$natural").One(&result)

		//		// Only update if there is new data (not implemented due to confusion with the assignment)
		//		goRates := Rates{}
		//		err      = getJSON("https://api.fixer.io/latest", &goRates)
		//		if err != nil {
		//			return err
		//		}
		//
		//		// Updates the DB if the date in the last DB entry is lower than the date in the last entry in api.fixer.io
		//		// Assumes that errors indicates there are no entries yet
		//		if result.Date < goRates.Date || err != nil {

		// Updates the DB if the date in the last DB entry is lower than system clock
		// Assumes that errors indicates there are no entries yet
		if result.Date < strings.Split(time.Now().Local().String(), " ")[0] || err != nil {
			// Fetches the latest data from api.fixer.io
			goRates := resource.Rates{}
			err := resource.GetJSON("https://api.fixer.io/latest", &goRates)
			if err != nil {
				return err
			}

			// Change the date to today and insert to our DB
			goRates.Date = strings.Split(time.Now().Local().String(), " ")[0]
			err = c.Insert(&goRates)
			fmt.Print(&goRates, "\n")
			if err != nil {
				return err
			}

			fmt.Print("Ticker Updated\n")
		} else {
			fmt.Print("Ticker Not updated\n")
		}
		session.Close()

		err = resource.AutoTriggerCheck()
		if err != nil {
			fmt.Print("An error occured somewhere at autoTriggerCheck: ", err, "\n")
		}
	}
	return nil
}
