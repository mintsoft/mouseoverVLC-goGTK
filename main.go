package main

import (
	"fmt"
	"os"

	"internal/utils"

	vlc "github.com/adrg/libvlc-go/v3"
	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

const appID = "com.github.mintsoft.vlc"


func builderGetObject(builder *gtk.Builder, name string) glib.IObject {
	obj, _ := builder.GetObject(name)
	return obj
}

func playerReleaseMedia(player *vlc.Player) {
	player.Stop()
	if media, _ := player.Media(); media != nil {
		media.Release()
	}
}

func main() {

	err := vlc.Init("--quiet", "--no-xlib")
	utils.AssertErr(err)

	player, err := vlc.NewPlayer()
	utils.AssertErr(err)

	app, err := gtk.ApplicationNew(appID, glib.APPLICATION_FLAGS_NONE)
	utils.AssertErr(err)

	app.Connect("activate", func() {
		builder, err := gtk.BuilderNewFromFile("layout.glade")
		utils.AssertErr(err)

		appWin, ok := builderGetObject(builder, "appWindow").(*gtk.ApplicationWindow)
		utils.AssertConv(ok)

		// Add builder signal handlers.
		signals := map[string]interface{}{
			"onRealizePlayerArea": func(playerArea *gtk.DrawingArea) {
				playerWindow, err := playerArea.GetWindow()
				utils.AssertErr(err)
				err = setPlayerWindow(player, playerWindow)
				utils.AssertErr(err)
				player.SetMouseInput(false)
			},
			"onDrawPlayerArea": func(playerArea *gtk.DrawingArea, cr *cairo.Context) {
				cr.SetSourceRGB(0, 0, 0)
				cr.Paint()
			},
			"onActivateQuit": func() {
				app.Quit()
			},
			"appWindow_motion_notify_event_cb": func(entry *gtk.ApplicationWindow, event *gdk.Event) bool {
				eventMotion := gdk.EventMotionNewFromEvent(event)
				x, y := eventMotion.MotionVal()
				fmt.Printf("AWMotionNotify x: %f y: %f\n", x, y)
				return false
			},
			/*
			"playerArea_motion_notify_event_cb": func(entry *gtk.ApplicationWindow, event *gdk.Event) bool {
				eventMotion := gdk.EventMotionNewFromEvent(event)
				x, y := eventMotion.MotionVal()
				fmt.Printf("PAMotionNotify x: %f y: %f\n", x, y)
				return false
			},
			*/
			"open_activate_cb": func() {
				fmt.Println("open activate called")
				player.LoadMediaFromURL("https://gstreamer.freedesktop.org/data/media/sintel_trailer-480p.webm")
				player.Play()
			},
		}
		builder.ConnectSignals(signals)
		appWin.ShowAll()
		app.AddWindow(appWin)
		
		appWin.AddEvents((int)(gdk.POINTER_MOTION_MASK))
	})

	// Cleanup on exit.
	app.Connect("shutdown", func() {
		playerReleaseMedia(player)
		player.Release()
		vlc.Release()
	})
/*
	go func() {
		for {
			pumpVlcMouseLocation(player)
			time.Sleep(1000/60 * time.Millisecond)
		}
	}()
-*/
	os.Exit(app.Run(os.Args))
}

var prevMouseX = 0 
var prevMouseY = 0

func pumpVlcMouseLocation(player *vlc.Player) {
	if(player != nil) {
		cursorLocationRelativeToVideoX,cursorLocationRelativeToVideoY,cursorError := player.CursorPosition()
		dimsX, dimsY, videoDimensionError := player.VideoDimensions()

		if cursorError != nil  || videoDimensionError != nil{
			return
		}

		if(cursorLocationRelativeToVideoX >= 0 && cursorLocationRelativeToVideoY >= 0 && cursorLocationRelativeToVideoX <= int(dimsX) && cursorLocationRelativeToVideoY <= int(dimsY)) {
			if(prevMouseX != cursorLocationRelativeToVideoX || prevMouseY != cursorLocationRelativeToVideoY) {
				
				fmt.Printf("pumpedMouseX: %d, pumped MouseY: %d\n", cursorLocationRelativeToVideoX, cursorLocationRelativeToVideoY)

				prevMouseX = cursorLocationRelativeToVideoX
				prevMouseY = cursorLocationRelativeToVideoY
			}
		}
	}
}
