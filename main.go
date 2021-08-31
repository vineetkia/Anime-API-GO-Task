package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

type animeStruct struct {
	ID           string `json:"id"`
	AnimeID      string `json:"animeID"`
	Title        string `json:"title"`
	Rank         string `json:"rank"`
	Score        string `json:"score"`
	Members      string `json:"members"`
	ImageUrl     string `json:"imageUrl"`
	Popularity   string `json:"popularity"`
	Status       string `json:"status"`
	Rating       string `json:"rating"`
	JapaneseWord string `json:"japaneseWord"`
	Synopsis     string `json:"synopsis"`
}

var anim []animeStruct

func getTitle(body string) string { // fetch title of anime from response passed in parameter

	startPoint := "<title>"
	endPoint := strings.Split(string(body), startPoint)
	return (strings.Split(string(endPoint[1]), "</title>"))[0]
}
func getPopularity(body string) string { // fetch popularity of anime from response passed in parameter

	startPoint := "><span class=\"numbers popularity\">Popularity <strong>#"
	endPoint := strings.Split(string(body), startPoint)
	return (strings.Split(string(endPoint[1]), "</strong>"))[0]
}
func getSynopsis(body string) string { // fetch synopsis of anime from response passed in parameter

	startPoint4 := "<p itemprop=\"description\">"
	endPoint4 := strings.Split(string(body), startPoint4)
	return (strings.Split(string(endPoint4[1]), "</p>"))[0]
}
func getMembers(body string) string { // fetch members of anime from response passed in parameter

	startPoint := "<span class=\"numbers members\">Members <strong>"
	endPoint := strings.Split(string(body), startPoint)
	return (strings.Split(string(endPoint[1]), "</strong>"))[0]
}

func getRank(body string) string { // fetch rank of anime from response passed in parameter

	startPoint := "<span class=\"dark_text\">Ranked:</span>\n  #"
	endPoint := strings.Split(string(body), startPoint)
	return (strings.Split(string(endPoint[1]), "<sup>"))[0]
}

func getImageUrl(body string) string { // fetch image url of anime from response passed in parameter

	startPoint := "<img class=\"lazyload\" data-src=\""
	endPoint := strings.Split(string(body), startPoint)
	return (strings.Split(string(endPoint[1]), "\" alt="))[0]
}

func getScore(body string) string { // fetch score of anime from response passed in parameter

	startPoint := "<span itemprop=\"ratingValue\" class=\"score-label score-"
	epOne := strings.Split(string(body), startPoint)
	epTwo := (strings.Split(string(epOne[1]), "ratingCount"))[0]
	result := (strings.Split(string(epTwo), "\">"))[1]
	return (strings.Split(string(result), "</span>"))[0]
}

func getStatus(body string) string { // fetch status of anime from response passed in parameter

	startPoint := "<span class=\"dark_text\">Status:</span>\n  "
	endPoint := strings.Split(string(body), startPoint)
	return (strings.Split(string(endPoint[1]), "\n  </div>"))[0]
}

func getRating(body string) string { // fetch rating of anime from response passed in parameter

	startPoint := "<span class=\"dark_text\">Rating:</span>\n  "
	endPoint := strings.Split(string(body), startPoint)
	return (strings.Split(string(endPoint[1]), "\n  </div>"))[0]
}

func getJapaneseWord(body string) string { // fetch japanese word of anime from response passed in parameter

	startPoint := "<span class=\"dark_text\">Japanese:</span> "
	endPoint := strings.Split(string(body), startPoint)
	return (strings.Split(string(endPoint[1]), "\n  </div><br"))[0]
}

func getData(w http.ResponseWriter, r *http.Request) { // retrieve all the data using GET METHOD from the database structure about anime from animeStruct Object
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(anim)
}

func getDataId(w http.ResponseWriter, r *http.Request) { // retrieve data for specific anime using GET METHOD by parsing the id parameter into the endpoint /anim/{id}
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range anim {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&animeStruct{})
}

func postData(w http.ResponseWriter, r *http.Request) { // POST Method to post animeID into parameter and then store all the anime info in db object and retreive response of details in json format
	w.Header().Set("Content-Type", "application/json")

	var animDat animeStruct
	_ = json.NewDecoder(r.Body).Decode(&animDat)
	animDat.ID = strconv.Itoa(rand.Intn(999999))

	resp, err := http.Get("https://myanimelist.net/anime/" + string(animDat.AnimeID)) // Getting Page data response from myanimelist
	if err != nil {
		log.Fatal(err) // throw error if operation failed
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	// below are the struct variables which stores all the info for the anime data
	animDat.Title = getTitle(string(body))
	animDat.Popularity = getPopularity(string(body))
	animDat.Synopsis = getSynopsis(string(body))
	animDat.Members = getMembers(string(body))
	animDat.Rank = getRank(string(body))
	animDat.ImageUrl = getImageUrl(string(body))
	animDat.Score = getScore(string(body))
	animDat.Status = getStatus(string(body))
	animDat.Rating = getRating(string(body))
	animDat.JapaneseWord = getJapaneseWord(string(body))
	// this will append the final structure to the anim object and store it there in the form of a database
	anim = append(anim, animDat)
	json.NewEncoder(w).Encode(animDat)
}

func main() {
	r := mux.NewRouter()

	anim = append(anim, animeStruct{ID: "1",
		AnimeID:      "1",
		Title:        "-",
		Rank:         "-",
		Score:        "-",
		Members:      "-",
		ImageUrl:     "#",
		Popularity:   "-",
		Status:       "-",
		Rating:       "-",
		JapaneseWord: "-",
		Synopsis:     "-"})

	r.HandleFunc("/anim", getData).Methods("GET")        // GET request method to retrieve all the data available in the structure database
	r.HandleFunc("/anim/{id}", getDataId).Methods("GET") // GET request method with id paramater in the url to get a specific anime details pre-stored in the database.
	r.HandleFunc("/anim", postData).Methods("POST")      // POST request method to retreive anime data from animeID which is passed in JSON form and then resolved in backend from postData function

	log.Fatal(http.ListenAndServe(":80", r)) // To log all the server fatal errors if encountered
}
