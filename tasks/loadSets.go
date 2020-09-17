package tasks

import (
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/arkenproject/arken/config"
	"github.com/arkenproject/arken/keysets"

	"github.com/arkenproject/arken/database"
)

func loadSets(keySets []config.KeySet, new chan database.FileKey, output chan database.FileKey, settings chan string) {
	// Run LoadSets every hour.
	for {
		fmt.Println("\n[Indexing & Updating Keysets]")
		err := keysets.LoadSets(keySets)
		if err != nil {
			log.Fatal(err)
		}

		for keySet := range keySets {
			location := filepath.Join(config.Global.Sources.Repositories, filepath.Base(keySets[keySet].URL))
			lighthouse, err := keysets.ConfigLighthouse(keySets[keySet].LightHouseFileID, keySets[keySet].URL)
			if err != nil {
				log.Fatal(err)
			}
			output <- lighthouse

			fmt.Printf("Indexing: %s\n", filepath.Base(keySets[keySet].URL))
			err = keysets.Index(location, new, output, settings)
			if err != nil {
				log.Fatal(err)
			}
		}

		fmt.Println("[Finished Indexing & Updating Keysets]")
		time.Sleep(1 * time.Hour)
	}
}