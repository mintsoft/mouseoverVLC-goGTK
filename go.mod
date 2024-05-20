module vlc-mouseover-test

go 1.19

require (
	github.com/adrg/libvlc-go/v3 v3.1.5
	github.com/gotk3/gotk3 v0.6.3
)

require (
	github.com/andlabs/ui v0.0.0-20200610043537-70a69d6ae31e // indirect
	github.com/creack/goselect v0.1.2 // indirect
	golang.org/x/sys v0.18.0 // indirect
)

require internal/utils v1.0.0

replace internal/utils => ./internal/utils