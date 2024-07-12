package velk

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"path/filepath"
)

func LoadZones() error {
	dir := filepath.Join("wld", "zone")
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, file := range files {
		if file.IsDir() {
			zonedir := filepath.Join(dir, file.Name(), "metadata.json")
			log.Println("Loading Zone", zonedir)
			zone := &Zone{}

			data, err := ioutil.ReadFile(zonedir)
			if err != nil {
				return err
			}
			err = json.Unmarshal(data, zone)
			if err != nil {
				return err
			}
			Zones[zone.Id] = zone
			err = zone.LoadRooms()
			if err != nil {
				return err
			}
		}
	}
	return nil
}
