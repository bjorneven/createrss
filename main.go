package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
)

func main() {

	folderInPtr := flag.String("inDir", ".", "The folder with .torrents. Default current dir")
	folderOutPtr := flag.String("outDir", ".", "The folder for output. Default current dir")
	feedNamePtr := flag.String("filename", "feed.xml", "Filename without directory, for torrent feed. Default feed.xml")
	baseURLPtr := flag.String("baseUrl", "", "The baseURL which the torrents can be downloaded from. Default null.")

	flag.Parse()

	// read folder In
	folderIn := filepath.Dir(*folderInPtr)
	folderOut := filepath.Dir(*folderOutPtr)
	feedName := *feedNamePtr
	baseURL := *baseURLPtr

	if *baseURLPtr == "" {
		log.Fatal("Missing baseURL")
		os.Exit(-1)
	}

	// output buffer
	var outBuffer bytes.Buffer // A Buffer needs no initialization.

	// head of xml feed
	outBuffer.Write([]byte(`		
<rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom">
<channel>
<title>Zibo Mod Files</title>
<link>https://forums.x-plane.org/index.php?/forums/topic/185685-zibo-install-guide-training-checklist/</link>
<description>Zibo Mod Install Guide, Training Checklist and Updates</description>
<atom:link href="https://1drv.ms/u/s!AjcwFonqlaRWgYdsjr2JbyxDPdcCvA?e=OENfhN" rel="self" type="application/rss+xml" />
`))

	fmt.Printf("Running with options\n inDir= %s\n outDir= %s\n filename= %s\n baseUrl= %s\n", folderIn, folderOut, feedName, baseURL)

	fmt.Println("\nScanning directory " + folderIn + "\n")

	err := filepath.Walk(folderIn, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && (filepath.Ext(path) == ".torrent" || filepath.Ext(path) == ".pdf") {
			fmt.Printf("\t%s/%s\n", folderIn, path)

			title := fmt.Sprint("\t<title>", info.Name(), "</title>\n")
			pubDate := fmt.Sprint("\t<pubDate>", info.ModTime(), "</pubDate>\n")
			url, err := url.Parse(baseURL + "/" + info.Name())

			if err != nil {
				log.Fatal("BaseURL in wrong format", err)

			}

			urlStr := fmt.Sprint("\t<link>", url.String(), "</link>\n")
			description := "\t<description />\n"

			outBuffer.Write([]byte("<item>\n"))
			outBuffer.Write([]byte(title))
			outBuffer.Write([]byte(description))
			outBuffer.Write([]byte(urlStr))
			outBuffer.Write([]byte(pubDate))

			outBuffer.Write([]byte("</item>\n"))
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	outBuffer.Write([]byte(`
	</channel>
	</rss>`))

	file := folderOut + "/" + feedName

	ioutil.WriteFile(file, outBuffer.Bytes(), 0644)

	fmt.Printf("Writing %s  ", file)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Done")

}
