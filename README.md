# go-icloud-shared-album

A small Go library for reading the contents of a **public iCloud Shared Album** —
the kind of album you get from Apple's *"Share" → "Public Website"* feature
(URLs that look like `https://www.icloud.com/sharedalbum/#B0...`).

It talks to the same undocumented `sharedstreams.icloud.com` endpoints the iCloud
web viewer uses, and returns the album's metadata, photos, and direct download
URLs for each image derivative (thumbnail, full-size, etc.).

No Apple account or authentication is required — it only works with albums that
have been made publicly shareable.

## Installation

```sh
go get github.com/grMLEqomlkkU5Eeinz4brIrOVCUCkJuN/go-icloud-shared-album
```

```go
import (
    "github.com/grMLEqomlkkU5Eeinz4brIrOVCUCkJuN/go-icloud-shared-album/src/core"
    "github.com/grMLEqomlkkU5Eeinz4brIrOVCUCkJuN/go-icloud-shared-album/src/utils"
)
```

## The album token

Every shared album has a **token** embedded in its public URL, after the `#`:

```
https://www.icloud.com/sharedalbum/#B0aBcDeFgHiJkL
                                     ^^^^^^^^^^^^^^^
                                     this is the token
```

That token (e.g. `B0aBcDeFgHiJkL`) is the only input you need.

## How it works

Fetching an album is a short pipeline:

1. **`utils.GetBaseUrl(token)`** — derives the initial API base URL from the
   token. The first characters of the token encode which iCloud server
   partition (`p01`, `p23`, …) hosts the album.
2. **`utils.GetRedirectedBaseUrl(baseUrl, token)`** — iCloud may respond with an
   HTTP `330` redirect pointing at the album's real host (via the
   `X-Apple-MMe-Host` header). This resolves that redirect and returns the final
   base URL to use. If there's no redirect, the original URL is returned
   unchanged.
3. **`core.GetApiResponse(baseUrl)`** — calls the `webstream` endpoint and
   returns album metadata plus every photo. Each photo carries one or more
   **derivatives** (different sizes), identified by a `Checksum`.
4. **`core.GetUrls(baseUrl, photoGuids)`** — calls the `webasseturls` endpoint
   and returns a `map[checksum]downloadURL`.
5. **`utils.EnrichImagesWithUrls(album, urls)`** *(optional helper)* — joins the
   two responses, returning images with each derivative's download URL filled
   in. You can also do this join yourself by looking up each derivative's
   `Checksum` in the map from step 4.

## Quick start

```go
package main

import (
    "fmt"
    "log"

    "github.com/grMLEqomlkkU5Eeinz4brIrOVCUCkJuN/go-icloud-shared-album/src/core"
    "github.com/grMLEqomlkkU5Eeinz4brIrOVCUCkJuN/go-icloud-shared-album/src/utils"
)

func main() {
    token := "B0aBcDeFgHiJkL" // the part of the share URL after '#'

    // 1. Derive the initial base URL from the token.
    baseURL := utils.GetBaseUrl(token)

    // 2. Resolve any iCloud host redirect (330).
    baseURL, err := utils.GetRedirectedBaseUrl(baseURL, token)
    if err != nil {
        log.Fatal(err)
    }

    // 3. Fetch the album: metadata + photos.
    album, err := core.GetApiResponse(baseURL)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Album: %q by %s %s — %d items\n",
        album.Metadata.StreamName,
        album.Metadata.UserFirstName,
        album.Metadata.UserLastName,
        album.Metadata.ItemsReturned,
    )

    // 4. Resolve download URLs (keyed by derivative checksum).
    urls, err := core.GetUrls(baseURL, album.PhotoGuids)
    if err != nil {
        log.Fatal(err)
    }

    // 5. Join photos and URLs into ready-to-use images. After this each
    //    derivative's URL field is populated, and derivatives are keyed by
    //    pixel height (e.g. "342", "2048").
    images := utils.EnrichImagesWithUrls(*album, urls)

    for _, photo := range images {
        fmt.Printf("\nPhoto %s (%dx%d) %q\n",
            photo.PhotoGuid, photo.Width, photo.Height, photo.Caption)

        for height, d := range photo.Derivatives {
            fmt.Printf("  derivative h=%-5s %dx%d  %d bytes  ->  %s\n",
                height, d.Width, d.Height, d.FileSize, *d.URL)
        }
    }
}
```

## API reference

### `utils.GetBaseUrl(token string) string`

Builds the initial `https://pXX-sharedstreams.icloud.com/<token>/sharedstreams/`
base URL by decoding the server partition from the token. Pure function, no
network access.

### `utils.GetRedirectedBaseUrl(baseUrl, token string) (string, error)`

Performs one `webstream` request and, if iCloud answers with a `330`, returns a
new base URL built from the `X-Apple-MMe-Host` header. Otherwise returns
`baseUrl` unchanged. Always call this before `GetApiResponse` to be safe.

### `core.GetApiResponse(baseUrl string) (*types.ApiResponse, error)`

Calls `webstream` and returns the parsed album. Numeric and date fields that
iCloud sends as strings are converted to `int` / `time.Time` here.

```go
type ApiResponse struct {
    Photos     map[string]Image // keyed by photoGuid
    PhotoGuids []string         // every photoGuid, for GetUrls
    Metadata   Metadata
}
```

### `core.GetUrls(baseUrl string, photoGuids []string) (map[string]string, error)`

Calls `webasseturls` and returns a map of **derivative checksum → fully-qualified
download URL** (`https://<location><path>`). Look up each
`Image.Derivatives[size].Checksum` in this map to get its link.

### `utils.EnrichImagesWithUrls(album types.ApiResponse, urls map[string]string) []types.Image`

Convenience helper that joins the two responses for you: it returns the album's
images with every derivative's `URL` field populated from the checksum map.
Derivatives whose checksum is missing from `urls` are dropped, and the surviving
derivatives are re-keyed by their pixel **height** (with `-1`, `-2`, … suffixes
if two derivatives share a height). Note `GetApiResponse` returns a pointer, so
pass `*album`.

## Data types

All types live in `src/types`.

```go
type Image struct {
    BatchGUID            string
    Derivatives          map[string]Derivative // keyed by size, e.g. "342", "2048"
    ContributorLastName  string
    BatchDateCreated     time.Time
    ContributorFirstName string
    ContributorFullName  string
    Caption              string
    Height               int
    Width                int
    MediaAssetType       *string
    PhotoGuid            string
    DateCreated          time.Time
}

type Derivative struct {
    Checksum string  // join key for GetUrls
    FileSize int
    Width    int
    Height   int
    URL      *string // nil until you populate it from GetUrls
}

type Metadata struct {
    StreamName    string
    UserFirstName string
    UserLastName  string
    StreamCtag    string
    ItemsReturned int
    Locations     interface{}
}
```

A photo's `Derivatives` map holds the available sizes; each derivative's
`Checksum` is what you pass through `GetUrls` to obtain the downloadable URL.
The largest derivative is generally the original-resolution image.

## Notes & caveats

- These iCloud endpoints are **undocumented and unofficial**; Apple can change
  them at any time.
- Only **public** shared albums work — there is no authentication path.
- Be considerate with request volume to avoid rate limiting.

## License

See repository for license details.
