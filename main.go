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

		playerArea, ok := builderGetObject(builder, "playerArea").(*gtk.DrawingArea)
		utils.AssertConv(ok)
		playerArea.AddEvents(int(gdk.POINTER_MOTION_MASK))

		signals := map[string]interface{}{
			"onRealizePlayerArea": func(playerArea *gtk.DrawingArea) {
				playerWindow, err := playerArea.GetWindow()
				utils.AssertErr(err)
				err = setPlayerWindow(player, playerWindow)
				utils.AssertErr(err)
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
			"playerArea_motion_notify_event_cb": func(entry *gtk.DrawingArea, event *gdk.Event) bool {
				return false
			},
			"open_activate_cb": func() {
				fmt.Println("open activate called")

				player.SetMouseInput(true) //makes no difference
				player.SetKeyInput(true)   //makes no difference

				player.LoadMediaFromURL("https://gstreamer.freedesktop.org/data/media/sintel_trailer-480p.webm")
				player.Play()
			},
		}
		builder.ConnectSignals(signals)
		appWin.ShowAll()
		app.AddWindow(appWin)

		//causes appWindow_motion_notify_event_cb to receive events no matter whereabouts in the window the cursor
		//is located
		appWin.AddEvents((int)(gdk.POINTER_MOTION_MASK))
	})

	// Cleanup on exit.
	app.Connect("shutdown", func() {
		playerReleaseMedia(player)
		player.Release()
		vlc.Release()
	})

	os.Exit(app.Run(os.Args))
}
