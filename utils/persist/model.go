/*
 * Copyright (c) 2021.  -present, Broos Action, Inc. All rights reserved.
 *
 *  This source code is licensed under the MIT license
 *  found in the LICENSE file in the root directory of this source tree.
 */

package persist

import (
	"archive/zip"
	"bytes"
	"fmt"
	"github.com/broosaction/gotext/tokenizers"
	"io/ioutil"
)


const (
	SerializationFormatVersion = "Joi 10"
)


// in a classification file format.
type Modeldata struct {
	// FormatVersion should always be 1 for this structure
	FormatVersion string
	// Uses the classifier name (provided by the classifier)
	Name string
	// ClassifierVersion is also provided by the classifier
	// and checks whether this version of GoLearn can read what's
	// be written.
	Version string

	Data []byte
}

// Save the persitable model.
func Save(filename string, meta Modeldata) {

	// Create a buffer to write our archive to.
	buf := new(bytes.Buffer)

	//in case the packaging changes, one must first do a compatibility test before actually loading
	meta.FormatVersion = SerializationFormatVersion
	// Create a new zip archive.
	zipWriter := zip.NewWriter(buf)



	// Add some files to the archive.
	var files = []struct {
		Name string
		Body []byte
	}{
		//the model files, not to be changed
		{"model.joi", meta.Data},
		{"meta.conf", []byte(meta.FormatVersion+"\n"+meta.Name+"\n"+meta.Version)},

	}
	for _, file := range files {
		zipFile, err := zipWriter.Create(file.Name)
		if err != nil {
			fmt.Println(err)
		}
		_, err = zipFile.Write([]byte(file.Body))
		if err != nil {
			fmt.Println(err)
		}
	}

	// Make sure to check the error on Close.
	err := zipWriter.Close()
	if err != nil {
		fmt.Println(err)
	}

	//write the zipped file to the disk
	ioutil.WriteFile(filename, buf.Bytes(), 0777)

}


type myCloser interface {
	Close() error
}

// closeFile is a helper function which streamlines closing
// with error checking on different file types.
func closeFile(f myCloser) {
	err := f.Close()
	check(err)
}

// readAll is a wrapper function for ioutil.ReadAll. It accepts a zip.File as
// its parameter, opens it, reads its content and returns it as a byte slice.
func readAll(file *zip.File) []byte {
	fc, err := file.Open()
	check(err)
	defer closeFile(fc)

	content, err := ioutil.ReadAll(fc)
	check(err)

	return content
}

// check is a helper function which streamlines error checking
func check(e error) {
	if e != nil {
		panic(e)
	}
}

//Load the last model that was saved.
func Load(filename string) Modeldata {

	zf, err := zip.OpenReader(filename)
	check(err)
	defer closeFile(zf)

	meta := Modeldata{}

	for _, file := range zf.File {

          if file.Name == "model.joi"{
          	meta.Data = readAll(file)
		  }else if file.Name == "meta.conf"{
			  tokenizer := tokenizers.DefaultTokenizer{
			  	RemoveStopWords: false,
			  }
			  str := tokenizer.Tokenize(string(readAll(file)))
			  meta.FormatVersion = str[0]
			  meta.Name          = str[1]
			  meta.Version       = str[2]

		  }
	}

	return meta
}