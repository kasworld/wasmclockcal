// +build ignore

// Copyright 2015,2016,2017,2018,2019 SeukWon Kang (kasworld@gmail.com)
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"fmt"
	"net/http"
)

func main() {
	port := flag.String("port", ":8080", "Serve port")
	folder := flag.String("dir", ".", "Serve Dir")
	fmt.Printf("dir2http dir=%v port=%v , http://localhost%v/ \n",
		*folder, *port, *port)

	webMux := http.NewServeMux()
	webMux.Handle("/",
		http.FileServer(http.Dir(*folder)),
	)
	if err := http.ListenAndServe(*port, webMux); err != nil {
		fmt.Println(err.Error())
	}
}
