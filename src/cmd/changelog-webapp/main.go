// Copyright 2020 Google Inc. All Rights Reserved.
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

package main

import (
	"net/http"
	"os"

	"cos.googlesource.com/cos/tools.git/src/cmd/changelog-webapp/controllers"

	log "github.com/sirupsen/logrus"
)

var (
	staticBasePath string
	port           string
)

func init() {
	staticBasePath = os.Getenv("STATIC_BASE_PATH")
	port = os.Getenv("PORT")
}

func main() {
	log.SetLevel(log.DebugLevel)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(staticBasePath))))
	http.HandleFunc("/", controllers.HandleIndex)
	http.HandleFunc("/readme/", controllers.HandleReadme)
	http.HandleFunc("/changelog/", controllers.HandleChangelog)
	http.HandleFunc("/findbuild/", controllers.HandleFindBuild)
	http.HandleFunc("/findreleasedbuildv2/", controllers.HandleFindReleasedBuild)
	http.HandleFunc("/findreleasedbuild", controllers.HandleFindReleasedBuildGerrit)
	http.HandleFunc("/login/", controllers.HandleLogin)
	http.HandleFunc("/oauth2callback/", controllers.HandleCallback)
	http.HandleFunc("/signout/", controllers.HandleSignOut)

	if port == "" {
		port = "8081"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
