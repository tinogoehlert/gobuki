package main

import (
	"flag"
	"log"
	"net"
	"os"
	"time"
)

type frame struct {
	buffer []byte
}

func prepareBuffer(filedump string) []frame {
	f, err := os.Open(filedump)
	if err != nil {
		log.Fatalf("could not open file: %s", err.Error())
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		log.Fatalf("could not stat file: %s", err.Error())
	}

	data := make([]byte, stat.Size())
	_, err = f.Read(data)
	if err != nil {
		log.Fatalf("could not read file: %s", err.Error())
	}

	frames := make([]frame, 0, 128)
	for i := 0; i+1 < len(data); {
		if data[i] == 0xAA && data[i+1] == 0x55 {
			frameLen := int(data[i+2])
			feedbackFrame := frame{
				buffer: data[i : i+frameLen+4],
			}
			frames = append(frames, feedbackFrame)
			i += (frameLen + 2)
		} else {
			i++
		}
	}

	return frames
}

func init() {
	flag.Parse()
}

func main() {
	listener, err := net.Listen("tcp", ":3333")
	if err != nil {
		log.Fatalf("could not create Listener: %s", err.Error())
	}

	log.Println("use dump: ", flag.Arg(0))
	defer listener.Close()
	dump := prepareBuffer(flag.Arg(0))
	log.Println("start serving")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("could not create Listener: %s", err.Error())
		}

		go func() {
			for {
				for {
					for _, f := range dump {
						conn.Write(f.buffer)
						time.Sleep(5 * time.Millisecond)
					}
				}
			}
		}()
	}
}
