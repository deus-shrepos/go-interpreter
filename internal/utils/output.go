package utils

import (
	"fmt"
	"io"
	"os"
)

// OuputStream Streams out to any stream that can read/write
type OutputStream struct {
	Writer io.Writer
	Reader io.Reader
}

func NewOutStream(writer io.Writer, reader io.Reader) OutputStream {
	outStream := OutputStream{}
	if writer == nil {
		outStream.Writer = os.Stdout
	}

	if reader == nil {
		outStream.Reader = os.Stdin
	}
	outStream.Writer = writer
	outStream.Reader = reader
	return outStream

}

func (o OutputStream) Print(msg string) {
	_, err := fmt.Fprint(o.Writer, msg)
	if err != nil {
		fmt.Println(err)
	}
}

func (o OutputStream) Error(msg error) {
	_, err := fmt.Fprint(os.Stdin, msg)
	if err != nil {
		fmt.Println(err)
	}
}
