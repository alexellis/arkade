// Copyright (c) arkade author(s) 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.
package archive

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"io/ioutil"
	"bytes"
)



// Added specifically to handle unpacking the argo-workflows
// binary, likely needs some work to be generally applicable
func Ungzip(r io.Reader, dir string, filename string) error {
	return ungzip(r, dir, filename)
}

// Write gunzipped data to a Writer
func gunzipWrite(w io.Writer, data []byte) error {
	// Write gzipped data to the client
	gr, err := gzip.NewReader(bytes.NewBuffer(data))
	defer gr.Close()
	data, err = ioutil.ReadAll(gr)
	if err != nil {
		return err
	}
	w.Write(data)
	return nil
}

func ungzip(r io.Reader, dir string, filename string) (err error) {
	abs := path.Join(dir, filename)
	defer func() {
		if err == nil {
			log.Printf("extracted gzip into %s", abs)
		} else {
			log.Printf("error extracting gzip into %s", abs)
		}
	}()

	zr, err := gzip.NewReader(r)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	defer zr.Close()

	data, err := ioutil.ReadAll(zr)
	if err != nil {
		return err
	}


	
	writer, err := os.Create(abs)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	defer writer.Close()

	if _, err = io.Copy(writer, zr); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	writer.Write(data)

	return nil
}
