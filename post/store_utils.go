package post

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func dumpContent(filePath string, lineLength uint) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	data := make([]byte, 4096)
	reader := bufio.NewReader(f)
	charIdx := uint(0)
	for {
		data = data[:cap(data)]
		n, err := reader.Read(data)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		data = data[:n]
		for _, b := range data {
			for i := 7; i >= 0; i-- { // read bits from msb to lsb...
				bit := GetNthBit(b, uint64(i))
				if bit {
					fmt.Print("1")
				} else {
					fmt.Print("0")
				}

				if charIdx == lineLength-1 {
					fmt.Println()
					charIdx = 0
				} else {
					charIdx += 1
				}
			}
		}
	}

	fmt.Printf("\n")
	return nil
}
