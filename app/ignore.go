package app

import (
	"fmt"
	"os"
)

func DownloadSpecialPacks(path string, f func(filePath, addr string)) {
	for _, server := range specialServers {
		fmt.Println("Downloading packs from", server)
		filePath := path + "/packs/" + server + "/19132/" + server + ".zip"

		if err := os.MkdirAll(path+"/packs/"+server+"/19132", 0777); err != nil {
			fmt.Println(err)
			continue
		}

		if err := downloadPacksFromServer(filePath, server+":19132"); err != nil {
			// Log the error (could use a proper logging framework)
			fmt.Println(err)
			continue
		}

		f(filePath, server)
	}
}
