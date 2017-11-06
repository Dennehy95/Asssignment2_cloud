package resource

// Imported resources
import (
	"bytes"
	"encoding/json"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"net/http"
	"strings"
)

// Project struct for the project
type Project struct {
	Id              bson.ObjectId `json:"-" bson:"_id"`
	WebHookURL      string        `json:"webhookURL"`
	BaseCurrency    string        `json:"baseCurrency"`
	TargetCurrency  string        `json:"targetCurrency"`
	MinTriggerValue float64       `json:"minTriggerValue"`
	MaxTriggerValue float64       `json:"maxTriggerValue"`
}

// Rates struct for the project
type Rates struct {
	Base  string                 `json:"base"`
	Date  string                 `json:"date"`
	Rates map[string]interface{} `json:"Rates"`
}

// Constant variables
const (
	GET    = "GET"
	POST   = "POST"
	DELETE = "DELETE"

	/* DB: MongoDB */
	Url         = "mongodb://user:123@ds133465.mlab.com:33465/assignment2"
	Database    = "assignment2"
	Collection  = "tuttut"
	Collection2 = "tottot"
)

// Handles request related to inserting documents to DB
func HandlerPost(w http.ResponseWriter, r *http.Request) {
	if r.Method == POST {
		http.Header.Add(w.Header(), "content-type", "application/json")

		// Fetch data from the request body
		res, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		post := &Project{}
		err = json.Unmarshal(res, &post)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		// Starts a session with the DB
		c, session, err := StartSession(Url, Database, Collection)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		defer session.Close()

		// Creates document id
		post.Id = bson.NewObjectId()

		// Uploads the content to the DB
		err = c.Insert(&post)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		w.Write([]byte(post.Id.Hex()))

	} else { // Invalid method
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
}

// Handles request related to getting or deleting documents in the DB
func HandlerGetDel(w http.ResponseWriter, r *http.Request) {
	if r.Method == GET || r.Method == DELETE {
		http.Header.Add(w.Header(), "content-type", "application/json")

		// Decode URL to fetch document Id
		parts := strings.Split(r.URL.Path, "/")
		Id := parts[len(parts)-1]

		// Starts a session with the DB
		c, session, err := StartSession(Url, Database, Collection)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		defer session.Close()

		if r.Method == GET {
			// Load data from the DB
			result := &Project{}
			err = c.Find(bson.M{"_id": bson.ObjectIdHex(Id)}).One(result)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}

			json.NewEncoder(w).Encode(result) // Return id to client
		} else { // Must be DELETE action
			err = c.Remove(bson.M{"_id": bson.ObjectIdHex(Id)})
			if err != nil {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
		}
	} else { // Invalid method used
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
}

// Handles request related to the latest value in DB
func HandlerLatest(w http.ResponseWriter, r *http.Request) {
	http.Header.Add(w.Header(), "content-type", "application/json")
	if r.Method == POST {

		// Starts a session with the DB
		c, session, err := StartSession(Url, Database, Collection2)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		defer session.Close()

		// Load data from the DB
		result := &Rates{}
		err = c.Find(nil).Sort("-$natural").One(result)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		// Fetch data from the request body
		res, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		post := &Project{}
		err = json.Unmarshal(res, &post)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		// Calculate conversion rate
		result.Rates[result.Base] = float64(1) //Inserts the DB base into DB rates
		if (result.Rates[post.BaseCurrency] != nil) && (result.Rates[post.TargetCurrency] != nil) {
			fmt.Fprint(w, result.Rates[post.TargetCurrency].(float64)/result.Rates[post.BaseCurrency].(float64))
		} else { // Invalid currency
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

	} else {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
}

// Handles request related to the latest values in DB
func HandlerAverage(w http.ResponseWriter, r *http.Request) {
	http.Header.Add(w.Header(), "content-type", "application/json")
	if r.Method == POST {

		c, session, err := StartSession(Url, Database, Collection2)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		defer session.Close()

		currency := []Rates{}
		err = c.Find(nil).Sort("-$natural").Limit(3).All(&currency)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		// Fetch data from the request body
		res, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		post := &Project{}
		err = json.Unmarshal(res, &post)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		// Calculate conversion rate average
		currency[0].Rates[currency[0].Base] = float64(1)
		if (currency[0].Rates[post.BaseCurrency] != nil) && (currency[0].Rates[post.TargetCurrency] != nil) {
			var numerator, denominator float64 //= 3

			for i := range currency {
				currency[i].Rates[currency[i].Base] = float64(1)
				numerator += currency[i].Rates[post.TargetCurrency].(float64)
				denominator += currency[i].Rates[post.BaseCurrency].(float64)
			}

			fmt.Fprint(w, numerator/denominator)
		} else {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
	} else {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
}

// Starts a session with a MongoDB database | Returns collection and session or error
func StartSession(link string, db string, col string) (*mgo.Collection, *mgo.Session, error) {
	session, err := mgo.Dial(link)
	if err != nil {
		return nil, nil, err
	}

	c := session.DB(db).C(col)
	return c, session, nil
}

// Fetches and decodes json. Takes an url and an interface as input.
func GetJSON(url string, target interface{}) error {
	r, err := http.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

//Automatic trigger used by the ticker
func AutoTriggerCheck() error {
	// Starts a session with the DBs
	c, session, err := StartSession(Url, Database, Collection)
	if err != nil {
		return err
	}

	c2, session2, err := StartSession(Url, Database, Collection2)
	if err != nil {
		return err
	}

	// Load data from the DB
	limits := []Project{}
	err = c.Find(nil).Sort("-$natural").All(&limits)
	if err != nil {
		return err
	}
	session.Close()

	rates := Rates{}
	err = c2.Find(nil).Sort("-$natural").One(&rates)
	if err != nil {
		return err
	}
	session2.Close()

	// Check if DB rates are within set limits
	rates.Rates[rates.Base] = float64(1)
	for i := range limits {
		if (rates.Rates[limits[i].BaseCurrency] != nil) && (rates.Rates[limits[i].TargetCurrency] != nil) {
			if limits[i].MaxTriggerValue <= rates.Rates[limits[i].TargetCurrency].(float64) || limits[i].MinTriggerValue >= rates.Rates[limits[i].TargetCurrency].(float64) {
				Invoker(rates.Rates[limits[i].TargetCurrency].(float64), limits[i])
			}
		}
	}

	return nil
}

func FullTriggerCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method == GET {
		http.Header.Add(w.Header(), "content-type", "application/json")
		// Starts a session with the DBs
		c, session, err := StartSession(Url, Database, Collection)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		c2, session2, err := StartSession(Url, Database, Collection2)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		// Load data from the DB
		limits := []Project{}
		err = c.Find(nil).Sort("-$natural").All(&limits)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		session.Close()

		rates := Rates{}
		err = c2.Find(nil).Sort("-$natural").One(&rates)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		session2.Close()

		// Check if DB rates are within set limits
		rates.Rates[rates.Base] = float64(1)
		for i := range limits {
			if (rates.Rates[limits[i].BaseCurrency] != nil) && (rates.Rates[limits[i].TargetCurrency] != nil) {
				err = Invoker(rates.Rates[limits[i].TargetCurrency].(float64), limits[i])
				if err != nil {
					http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
					return
				}
			}
		}

	} else {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
}

func Invoker(newRate float64, limits Project) error {
	output := make(map[string]interface{})

	output["baseCurrency"] = limits.BaseCurrency
	output["targetCurrency"] = limits.TargetCurrency
	output["currentRate"] = newRate
	output["minTriggerValue"] = limits.MinTriggerValue
	output["maxTriggerValue"] = limits.MaxTriggerValue

	response, err := json.Marshal(output)
	if err != nil {
		return err
	}

	reply, err := http.Post(limits.WebHookURL, "application/json", bytes.NewBuffer(response))
	if err != nil {
		fmt.Print("Error occured during http.Post: ", err, "\n")
		return err
	}

	fmt.Print("Status: ", reply.StatusCode)
	return nil
}
