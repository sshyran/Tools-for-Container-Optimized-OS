// Copyright 2018 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// cos_customizer is a Cloud Build builder for building custom COS images.
package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"

	"cos.googlesource.com/cos/tools.git/src/pkg/fs"

	"golang.org/x/oauth2/google"

	"cloud.google.com/go/storage"
	"github.com/google/subcommands"
	compute "google.golang.org/api/compute/v1"
	"google.golang.org/api/option"
)

var persistentDir = flag.String("local-state-workdir", ".cos-customizer-workdir",
	"Name of the directory in $HOME to use for storing local state.")

var computeEndpoint = flag.String("compute-endpoint", "", "If set, used as the endpoint for the GCE API.")

func clients(ctx context.Context, anonymousCreds bool) (*compute.Service, *storage.Client, error) {
	var httpClient *http.Client
	var err error
	if anonymousCreds {
		httpClient = &http.Client{}
	} else {
		httpClient, err = google.DefaultClient(ctx)
		if err != nil {
			return nil, nil, err
		}
	}
	computeOpts := []option.ClientOption{option.WithHTTPClient(httpClient)}
	if *computeEndpoint != "" {
		computeOpts = append(computeOpts, option.WithEndpoint(*computeEndpoint))
	}
	svc, err := compute.NewService(ctx, computeOpts...)
	if err != nil {
		return nil, nil, err
	}
	gcsClient, err := storage.NewClient(ctx, option.WithHTTPClient(httpClient))
	if err != nil {
		return nil, nil, err
	}
	return svc, gcsClient, nil
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(new(StartImageBuild), "")
	subcommands.Register(new(RunScript), "")
	subcommands.Register(new(InstallGPU), "")
	subcommands.Register(new(SealOEM), "")
	subcommands.Register(new(DisableAutoUpdate), "")
	subcommands.Register(new(FinishImageBuild), "")
	subcommands.Register(new(InstallPackage), "")
	flag.Parse()
	ctx := context.Background()
	files := fs.DefaultFiles(*persistentDir)
	ret := int(subcommands.Execute(ctx, files, ServiceClients(clients)))
	os.Exit(ret)
}
