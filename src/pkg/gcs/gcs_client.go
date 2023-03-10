// Copyright 2021 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gcs

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"strings"

	"cloud.google.com/go/storage"
)

const schemeGCS = "gs"

func readGCSObject(ctx context.Context, gcsClient *storage.Client, inputURL string) (*storage.Reader, error) {
	gcsBucket, name, err := getGCSVariables(inputURL)
	if err != nil {
		return nil, err
	}
	return gcsClient.Bucket(gcsBucket).Object(name).NewReader(ctx)
}

// DownloadGCSObject downloads the object at inputURL and saves it at destinationPath
func DownloadGCSObject(ctx context.Context,
	gcsClient *storage.Client, inputURL, destinationPath string) error {
	r, err := readGCSObject(ctx, gcsClient, inputURL)
	if err != nil {
		return err
	}
	defer r.Close()

	f, err := os.Create(destinationPath)
	if err != nil {
		return err
	}
	if _, err := io.Copy(f, r); err != nil {
		return fmt.Errorf("error copying file from gcs bucket: %v", err)
	}
	if err = f.Close(); err != nil {
		return fmt.Errorf("error closing file: %v", err)
	}
	return nil
}

// DownloadGCSObjectString downloads the object at inputURL and saves the contents of the object file to a string
func DownloadGCSObjectString(ctx context.Context,
	gcsClient *storage.Client, inputURL string) (string, error) {
	r, err := readGCSObject(ctx, gcsClient, inputURL)
	if err != nil {
		return "", err
	}
	defer r.Close()

	ret, err := ioutil.ReadAll(r)
	if err != nil {
		return "", fmt.Errorf("unable to read file from gcs bucket: %v", err)
	}
	return string(ret), nil
}

// GCSObjectExists checks if an object exists at inputURL
func GCSObjectExists(ctx context.Context, gcsClient *storage.Client, inputURL string) (bool, error) {
	gcsBucket, name, err := getGCSVariables(inputURL)
	_, err = gcsClient.Bucket(gcsBucket).Object(name).Attrs(ctx)
	if err == storage.ErrObjectNotExist {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// UploadGCSObject uploads an object at inputPath to destination URL
func UploadGCSObject(ctx context.Context, gcsClient *storage.Client, inputPath, destinationURL string) error {
	fileReader, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	return uploadGCSObject(ctx, gcsClient, fileReader, destinationURL)
}

// UploadGCSObjectString uploads an input string as a file to destination URL
func UploadGCSObjectString(ctx context.Context, gcsClient *storage.Client, inputStr, destinationURL string) error {
	reader := strings.NewReader(inputStr)
	return uploadGCSObject(ctx, gcsClient, reader, destinationURL)
}

func uploadGCSObject(ctx context.Context,
	gcsClient *storage.Client, reader io.Reader, destinationURL string) error {
	gcsBucket, name, err := getGCSVariables(destinationURL)
	if err != nil {
		return fmt.Errorf("error parsing destination URL: %v", err)
	}
	w := gcsClient.Bucket(gcsBucket).Object(name).NewWriter(ctx)
	if _, err := io.Copy(w, reader); err != nil {
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}
	return nil
}

// DeleteGCSObject deletes an object at the input URL
func DeleteGCSObject(ctx context.Context,
	gcsClient *storage.Client, inputURL string) error {
	gcsBucket, name, err := getGCSVariables(inputURL)
	if err != nil {
		return fmt.Errorf("error parsing input URL: %v", err)
	}
	return gcsClient.Bucket(gcsBucket).Object(name).Delete(ctx)
}

// Returns the getGCSVariables(GCSBucket, GCSPath, fileName) based on the input.
func getGCSVariables(gcsPath string) (string, string, error) {
	url, err := url.Parse(gcsPath)
	if err != nil || url.Scheme != schemeGCS {
		return "", "", fmt.Errorf("error parsing the input GCS path: %s", gcsPath)
	}
	// url.EscapedPath returns with the leading /.
	return url.Hostname(), strings.TrimLeft(url.EscapedPath(), "/"), nil
}
