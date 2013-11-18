sp
==

Spotify API client

Development status: incomplete

## Usage:
```Go
package main

import (
    "fmt"
    "github.com/shawnps/sp"
)

func main() {
    sp := sp.Spotify{}
    r, err := sp.SearchAlbums("Danny Brown")
    if err != nil {
        fmt.Println("ERROR: ", err.Error())
    }   
    for _, a := range r.Albums {
        fmt.Println(a)
    }   
}
```
