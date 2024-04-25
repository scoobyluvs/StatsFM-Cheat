package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

func addMilliseconds(ts string, msPlayed int) string {
	parsedTime, err := time.Parse("2006-01-02T15:04:05Z", ts)
	if err != nil {
		panic(err)
	}
	updatedTime := parsedTime.Add(time.Millisecond * time.Duration(msPlayed))
	return updatedTime.Format("2006-01-02T15:04:05Z")
}

func main() {
	customStartTime := "2023-01-21T00:00:00Z"

	fmt.Print("Enter the track ID: ")
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		fmt.Println("Error reading input:", scanner.Err())
		return
	}
	trackID := scanner.Text()

	initialData := []map[string]interface{}{
		{
			"ts":                                "2019-05-28T23:19:27Z",
			"ms_played":                         166606,
			"master_metadata_track_name":        "a",
			"master_metadata_album_artist_name": "a",
			"master_metadata_album_album_name":  "a",
			"spotify_track_uri":                 "spotify:track:" + trackID,
		},
	}

	dataList := make([]map[string]interface{}, 0)
	currentIndex := 0
	currentTS := customStartTime

	for {
		currentData := initialData[currentIndex]

		updatedTS := addMilliseconds(currentTS, int(currentData["ms_played"].(int)))

		updatedData := make(map[string]interface{})
		for key, value := range currentData {
			updatedData[key] = value
		}
		updatedData["ts"] = updatedTS

		dataList = append(dataList, updatedData)

		currentIndex = (currentIndex + 1) % len(initialData)
		currentTS = updatedTS

		parsedTime, err := time.Parse("2006-01-02T15:04:05Z", updatedTS)
		if err != nil {
			panic(err)
		}
		if parsedTime.UTC().After(time.Now().UTC()) {
			break
		}
	}

	outputFile, err := json.MarshalIndent(dataList, "", "    ")
	if err != nil {
		panic(err)
	}

	err = writeToFile("output.json", outputFile)
	if err != nil {
		panic(err)
	}

	fmt.Println("Data written to output.json")
}

func writeToFile(filename string, data []byte) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	err = file.Sync()
	if err != nil {
		return err
	}

	return nil
}
