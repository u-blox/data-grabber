/* Entry point for HTTP logging server.
 *
 * Copyright (C) u-blox Melbourn Ltd
 * u-blox Melbourn Ltd, Melbourn, UK
 *
 * All rights reserved.
 *
 * This source file is the sole property of u-blox Melbourn Ltd.
 * Reproduction or utilization of this source in whole or part is
 * forbidden without the written consent of u-blox Melbourn Ltd.
 */
 
package main

import (
    "fmt"
    "log"
    "flag"
    "os"
    "net/http"
    "io/ioutil"
)

//--------------------------------------------------------------------
// Types
//--------------------------------------------------------------------

//--------------------------------------------------------------------
// Variables
//--------------------------------------------------------------------

// File handle
var pFile *os.File = nil

// Port to listen to HTTP requests on
var httpPort uint64 = 0

// Command-line flags
var pHttpPort = flag.String ("p", "5055", "the port number to listen on for HTTP requests.");
var pFileName = flag.String ("t", "grabbed.txt", "the file name to append the body of received HTTP requests to.");
var Usage = func() {
    fmt.Fprintf(os.Stderr, "\n%s: run the HTTP data logging server.  Usage:\n", os.Args[0])
        flag.PrintDefaults()
    }

//--------------------------------------------------------------------
// Functions
//--------------------------------------------------------------------

// HTTP home page handler
func homeHttpHandler (writer http.ResponseWriter, request *http.Request) {
    if (request.Method != "") {
        log.Printf("  > received %s request from \"%s\":\n", request.Method, request.URL)
        body, err := ioutil.ReadAll(request.Body)
        if err == nil {
            log.Printf("  > \"%s\".\n", body)
            writer.WriteHeader(http.StatusOK)
        } else {
            log.Printf("!!> couldn't read body from %s request (%s).\n", request.Method, err.Error())
            writer.WriteHeader(http.StatusBadRequest)
        }
    } else {
        log.Printf("!!> received unsupported REST request from \"%s\".\n", request.URL)
    }
}

// Entry point
func main() {
    var err error

    // Set up logging
    log.SetFlags(log.LstdFlags)
    
    // Open the output file for append
    pFile, err = os.OpenFile(*pFileName, os.O_APPEND | os.O_CREATE, 0666)
    
    if (err == nil) {        
        // Set up the HTTP server
        http.HandleFunc("/", homeHttpHandler)    
        log.Printf("Listening for HTTP requests on port %s.\n", *pHttpPort)
        // Listen forever
        http.ListenAndServe(":" + *pHttpPort, nil)
    } else {
        log.Printf("Couldn't open file %s (%s).\n", *pFileName, err.Error())
    }
}

// End Of File