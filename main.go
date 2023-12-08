package main

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/schollz/progressbar/v3"
)

var (
	client     = flag.Bool("client", false, "Is this the client executable.")
	bulkSize   = flag.Int64("bulkSize", 0, "The amount of bulk data to send.")
	totalReads = flag.Int64("totalReads", 10, "The number of reads to do.")
)

func main() {
	flag.Parse()

	if *client {
		bulkData := make([]byte, *bulkSize)

		_, err := rand.Read(bulkData)
		if err != nil {
			log.Fatal(err)
		}

		bulkDataString := hex.EncodeToString(bulkData)

		for {
			now := time.Now()

			fmt.Fprintf(os.Stdout, "%d,%d,%s\n", now.Unix(), now.Nanosecond(), bulkDataString)
		}
	} else {
		cmd, err := os.Executable()
		if err != nil {
			log.Fatal(err)
		}

		proc := exec.Command(cmd, "-client", "-bulkSize", fmt.Sprintf("%d", *bulkSize))

		out, err := proc.StdoutPipe()
		if err != nil {
			log.Fatal(err)
		}

		err = proc.Start()
		if err != nil {
			log.Fatal(err)
		}
		defer proc.Process.Kill()

		scan := bufio.NewScanner(out)

		scan.Buffer(nil, int(*bulkSize*3))

		currentReads := 0

		pb := progressbar.DefaultBytes(*totalReads * *bulkSize)
		defer pb.Close()

		for scan.Scan() {
			line := scan.Text()

			var unix int64
			var nano int64

			_, err := fmt.Sscanf(line, "%d,%d", &unix, &nano)
			if err != nil {
				log.Fatal(err)
			}

			// sent := time.Unix(unix, nano)

			// log.Printf("time = %s", time.Since(sent))

			pb.Add64(*bulkSize)

			currentReads += 1

			if currentReads >= int(*totalReads) {
				break
			}
		}

		err = scan.Err()
		if err != nil {
			log.Fatal(err)
		}
	}
}
