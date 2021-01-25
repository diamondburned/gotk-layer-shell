package layershell

import (
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"unsafe"
)

// #cgo pkg-config: gtk-layer-shell-0 gtk+-3.0 glib-2.0 gio-2.0 glib-2.0 gobject-2.0
// #include <gtk-layer-shell/gtk-layer-shell.h>
// #include <gtk/gtk.h>
// #include <gio/gio.h>
// #include <glib.h>
// #include <glib-object.h>
//
import "C"

// objector is used internally for other interfaces.
type objector interface {
	glib.IObject
	Connect(string, interface{}) glib.SignalHandle
	ConnectAfter(string, interface{}) glib.SignalHandle
	GetProperty(name string) (interface{}, error)
	SetProperty(name string, value interface{}) error
	Native() uintptr
}

// asserting objector interface
var _ objector = (*glib.Object)(nil)

// Caster is the interface that allows casting objects to widgets.
type Caster interface {
	objector
	Cast() (gtk.IWidget, error)
}

func init() {
	glib.RegisterGValueMarshalers([]glib.TypeMarshaler{
		// Enums

		// Objects/Classes

		// Boxed
	})
}

type Edge int

func marshalEdge(p uintptr) (interface{}, error) {
	return Edge(C.g_value_get_enum((*C.GValue)(unsafe.Pointer(p)))), nil
}

const (
	// EdgeLeft the left edge of the screen.
	EdgeLeft Edge = 0
	// EdgeRight the right edge of the screen.
	EdgeRight Edge = 1
	// EdgeTop the top edge of the screen.
	EdgeTop Edge = 2
	// EdgeBottom the bottom edge of the screen.
	EdgeBottom Edge = 3
	// EdgeEntryNumber should not be used except to get the number of entries
	EdgeEntryNumber Edge = 4
)

type Layer int

func marshalLayer(p uintptr) (interface{}, error) {
	return Layer(C.g_value_get_enum((*C.GValue)(unsafe.Pointer(p)))), nil
}

const (
	// LayerBackground the background layer.
	LayerBackground Layer = 0
	// LayerBottom the bottom layer.
	LayerBottom Layer = 1
	// LayerTop the top layer.
	LayerTop Layer = 2
	// LayerOverlay the overlay layer.
	LayerOverlay Layer = 3
	// LayerEntryNumber should not be used except to get the number of entries
	LayerEntryNumber Layer = 4
)

// AutoExclusiveZoneEnable when auto exclusive zone is enabled, exclusive zone
// is automatically set to the size of the window + relevant margin. To disable
// auto exclusive zone, just set the exclusive zone to 0 or any other fixed
// value.
//
// NOTE: you can control the auto exclusive zone by changing the margin on the
// non-anchored edge. This behavior is specific to gtk-layer-shell and not part
// of the underlying protocol
func AutoExclusiveZoneEnable(window *gtk.Window) {
	v1 := (*C.GtkWindow)(unsafe.Pointer(window.Widget.Native()))
	C.gtk_layer_auto_exclusive_zone_enable(v1)
}
func AutoExclusiveZoneIsEnabled(window *gtk.Window) bool {
	v1 := (*C.GtkWindow)(unsafe.Pointer(window.Widget.Native()))
	r := gobool(C.gtk_layer_auto_exclusive_zone_is_enabled(v1))
	return r
}
func GetAnchor(window *gtk.Window, edge Edge) bool {
	v1 := (*C.GtkWindow)(unsafe.Pointer(window.Widget.Native()))
	v2 := C.GtkLayerShellEdge(edge)

	r := gobool(C.gtk_layer_get_anchor(v1, v2))
	return r
}
func GetExclusiveZone(window *gtk.Window) int {
	v1 := (*C.GtkWindow)(unsafe.Pointer(window.Widget.Native()))
	r := int(C.gtk_layer_get_exclusive_zone(v1))
	return r
}
func GetKeyboardInteractivity(window *gtk.Window) bool {
	v1 := (*C.GtkWindow)(unsafe.Pointer(window.Widget.Native()))
	r := gobool(C.gtk_layer_get_keyboard_interactivity(v1))
	return r
}
func GetLayer(window *gtk.Window) Layer {
	v1 := (*C.GtkWindow)(unsafe.Pointer(window.Widget.Native()))
	r := Layer(C.gtk_layer_get_layer(v1))
	return r
}
func GetMajorVersion() uint {
	r := uint(C.gtk_layer_get_major_version())
	return r
}
func GetMargin(window *gtk.Window, edge Edge) int {
	v1 := (*C.GtkWindow)(unsafe.Pointer(window.Widget.Native()))
	v2 := C.GtkLayerShellEdge(edge)

	r := int(C.gtk_layer_get_margin(v1, v2))
	return r
}
func GetMicroVersion() uint {
	r := uint(C.gtk_layer_get_micro_version())
	return r
}
func GetMinorVersion() uint {
	r := uint(C.gtk_layer_get_minor_version())
	return r
}

// GetMonitor nOTE: To get which monitor the surface is actually on, use
// C.gdk_display_get_monitor_at_window().
func GetMonitor(window *gtk.Window) *gdk.Monitor {
	v1 := (*C.GtkWindow)(unsafe.Pointer(window.Widget.Native()))
	obj := glib.Take(unsafe.Pointer(C.gtk_layer_get_monitor(v1)))
	r := &gdk.Monitor{
		Object: obj,
	}
	return r
}

// GetNamespace nOTE: this function does not return ownership of the string. Do
// not free the returned string. Future calls into the library may invalidate
// the returned string.
func GetNamespace(window *gtk.Window) string {
	v1 := (*C.GtkWindow)(unsafe.Pointer(window.Widget.Native()))
	r := C.GoString(C.gtk_layer_get_namespace(v1))
	return r
}
func GetZwlrLayerSurfaceV1(window *gtk.Window) unsafe.Pointer {
	v1 := (*C.GtkWindow)(unsafe.Pointer(window.Widget.Native()))
	r := (unsafe.Pointer)(C.gtk_layer_get_zwlr_layer_surface_v1(v1))
	return r
}

// InitForWindow set the window up to be a layer surface once it is mapped. this
// must be called before the window is realized.
func InitForWindow(window *gtk.Window) {
	v1 := (*C.GtkWindow)(unsafe.Pointer(window.Widget.Native()))
	C.gtk_layer_init_for_window(v1)
}
func IsLayerWindow(window *gtk.Window) bool {
	v1 := (*C.GtkWindow)(unsafe.Pointer(window.Widget.Native()))
	r := gobool(C.gtk_layer_is_layer_window(v1))
	return r
}

// IsSupported may block for a Wayland roundtrip the first time it's called.
func IsSupported() bool {
	r := gobool(C.gtk_layer_is_supported())
	return r
}

// SetAnchor set whether window should be anchored to edge. - If two
// perpendicular edges are anchored, the surface with be anchored to that corner
// - If two opposite edges are anchored, the window will be stretched across the
// screen in that direction
//
// Default is false for each LayerShellEdge
func SetAnchor(window *gtk.Window, edge Edge, anchorToEdge bool) {
	v1 := (*C.GtkWindow)(unsafe.Pointer(window.Widget.Native()))
	v2 := C.GtkLayerShellEdge(edge)
	v3 := cbool(anchorToEdge)

	C.gtk_layer_set_anchor(v1, v2, v3)
}

// SetExclusiveZone has no effect unless the surface is anchored to an edge.
// Requests that the compositor does not place other surfaces within the given
// exclusive zone of the anchored edge. For example, a panel can request to not
// be covered by maximized windows. See wlr-layer-shell-unstable-v1.xml for
// details.
//
// Default is 0
func SetExclusiveZone(window *gtk.Window, exclusiveZone int) {
	v1 := (*C.GtkWindow)(unsafe.Pointer(window.Widget.Native()))
	v2 := C.int(exclusiveZone)

	C.gtk_layer_set_exclusive_zone(v1, v2)
}

// SetKeyboardInteractivity whether the window should receive keyboard events
// from the compositor.
//
// Default is false
func SetKeyboardInteractivity(window *gtk.Window, interacitvity bool) {
	v1 := (*C.GtkWindow)(unsafe.Pointer(window.Widget.Native()))
	v2 := cbool(interacitvity)

	C.gtk_layer_set_keyboard_interactivity(v1, v2)
}

// SetLayer set the "layer" on which the surface appears (controls if it is over
// top of or below other surfaces). The layer may be changed on-the-fly in the
// current version of the layer shell protocol, but on compositors that only
// support an older version the window is remapped so the change can take
// effect.
//
// Default is K_LAYER_SHELL_LAYER_TOP
func SetLayer(window *gtk.Window, layer Layer) {
	v1 := (*C.GtkWindow)(unsafe.Pointer(window.Widget.Native()))
	v2 := C.GtkLayerShellLayer(layer)

	C.gtk_layer_set_layer(v1, v2)
}

// SetMargin set the margin for a specific edge of a window. Effects both
// surface's distance from the edge and its exclusive zone size (if auto
// exclusive zone enabled).
//
// Default is 0 for each LayerShellEdge
func SetMargin(window *gtk.Window, edge Edge, marginSize int) {
	v1 := (*C.GtkWindow)(unsafe.Pointer(window.Widget.Native()))
	v2 := C.GtkLayerShellEdge(edge)
	v3 := C.int(marginSize)

	C.gtk_layer_set_margin(v1, v2, v3)
}

// SetMonitor set the output for the window to be placed on, or nil to let the
// compositor choose. If the window is currently mapped, it will get remapped so
// the change can take effect.
//
// Default is nil
func SetMonitor(window *gtk.Window, monitor *gdk.Monitor) {
	v1 := (*C.GtkWindow)(unsafe.Pointer(window.Widget.Native()))
	v2 := (*C.GdkMonitor)(unsafe.Pointer(monitor.Native()))

	C.gtk_layer_set_monitor(v1, v2)
}

// SetNamespace set the "namespace" of the surface.
//
// No one is quite sure what this is for, but it probably should be something
// generic ("panel", "osk", etc). The name_space string is copied, and caller
// maintians ownership of original. If the window is currently mapped, it will
// get remapped so the change can take effect.
//
// Default is "gtk-layer-shell" (which will be used if set to nil)
func SetNamespace(window *gtk.Window, nameSpace string) {
	v1 := (*C.GtkWindow)(unsafe.Pointer(window.Widget.Native()))
	v2 := C.CString(nameSpace)
	defer C.free(unsafe.Pointer(v2))

	C.gtk_layer_set_namespace(v1, v2)
}
