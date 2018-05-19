package player

import "C"
import (
	"github.com/yiitz/iceapple/log"
	"unsafe"
)

/*
#cgo pkg-config: gstreamer-1.0
#include <stdio.h>
#include <stdlib.h>
#include <gst/gst.h>

extern void gst_log_error(char *);
extern void gst_log_debug(char *);

static char msgbuf[256];

static GMainLoop *loop;

typedef struct _CustomData {
  gboolean is_live;
  GstElement *pipeline;
  gint buffering_level;
} CustomData;

//非线程安全
static void log_call(void (*func)(char *) ,const char* fmt,...){
	va_list argptr;
	va_start(argptr,fmt);
	snprintf(msgbuf, sizeof(msgbuf), fmt, argptr);
	va_end(argptr);
	func(msgbuf);
}

static gboolean bus_callback(GstBus *bus, GstMessage *msg, CustomData *data){
	GError *err;
	gchar *debug;
  	switch (GST_MESSAGE_TYPE (msg)) {
		case GST_MESSAGE_ERROR:
			gst_message_parse_error (msg, &err, &debug);
			snprintf (msgbuf, sizeof(msgbuf), "bus callback error: %s", err->message);
			gst_log_debug (msgbuf);
			g_error_free (err);
			g_free (debug);
			gst_element_set_state (data->pipeline, GST_STATE_READY);
			break;
		case GST_MESSAGE_EOS:
			gst_element_set_state (data->pipeline, GST_STATE_READY);
			break;
		case GST_MESSAGE_BUFFERING:
			if (data->is_live) break;

			gst_message_parse_buffering (msg, &data->buffering_level);

			if (data->buffering_level < 100)
				gst_element_set_state (data->pipeline, GST_STATE_PAUSED);
			else
				gst_element_set_state (data->pipeline, GST_STATE_PLAYING);
			break;
		case GST_MESSAGE_CLOCK_LOST:
			gst_element_set_state (data->pipeline, GST_STATE_PAUSED);
			gst_element_set_state (data->pipeline, GST_STATE_PLAYING);
			break;
	}
}

static void got_location(GstObject *gstobject, GstObject *prop_object, GParamSpec *prop, gpointer data) {
	gchar *location;
	g_object_get (G_OBJECT (prop_object), "temp-location", &location, NULL);
	snprintf(msgbuf, sizeof(msgbuf), "temporary file: %s", location);
	gst_log_debug (msgbuf);
	g_free (location);
	g_object_set (G_OBJECT (prop_object), "temp-remove", FALSE, NULL);
}

static void loop_init(){
	gst_init(NULL,NULL);
	loop = g_main_loop_new(NULL,FALSE);
}

static void loop_run(){
	g_main_loop_run(loop);
}

static void loop_destroy_and_quit(){
	g_main_loop_quit(loop);
	g_main_loop_unref(loop);
}

static void* player_new(){
	GstBus *bus;
	CustomData *data;
	guint flags;

	data = malloc(sizeof(CustomData));
	data->pipeline = gst_element_factory_make("playbin", "playbin");
	g_object_get (data->pipeline, "flags", &flags, NULL);
	flags |= 0x80;
	g_object_set (data->pipeline, "flags", flags, NULL);

	bus = gst_element_get_bus(data->pipeline);
	if (NULL == bus) {
		gst_log_error("get pipeline bus error");
		return NULL;
	}
	gst_bus_add_signal_watch (bus);
  	g_signal_connect(bus, "message", G_CALLBACK (bus_callback), data);
  	gst_object_unref(bus);

  	g_signal_connect (data->pipeline, "deep-notify::temp-location", G_CALLBACK (got_location), NULL);

	return data;
}

static void player_destroy(void *player){
	CustomData *data = (CustomData*)player;
	gst_element_set_state(data->pipeline, GST_STATE_NULL);
	gst_object_unref(data->pipeline);
}

static void player_play(void *player, gchar *path){
	GstStateChangeReturn ret;
	CustomData *data = (CustomData*)player;
	g_object_set(G_OBJECT(data->pipeline), "uri", path, NULL);
	gst_element_set_state (data->pipeline, GST_STATE_READY);
	ret = gst_element_set_state(data->pipeline, GST_STATE_PLAYING);
	data->is_live = FALSE;
	if (ret == GST_STATE_CHANGE_FAILURE) {
		log_call(gst_log_error, "unable to set the pipeline to the playing state.");
		gst_object_unref (data->pipeline);
	} else if (ret == GST_STATE_CHANGE_NO_PREROLL) {
		data->is_live = TRUE;
	}
}

static void player_pause(void *player){
	CustomData *data = (CustomData*)player;
	gst_element_set_state(data->pipeline, GST_STATE_PAUSED);
}

static void player_resume(void *player){
	CustomData *data = (CustomData*)player;
	gst_element_set_state(data->pipeline, GST_STATE_PLAYING);
}

static void player_stop(void *player){
	CustomData *data = (CustomData*)player;
	gst_element_set_state(data->pipeline, GST_STATE_NULL);
}
*/
import "C"

var logger = log.NewLogger("[player]")

//export gst_log_error
func gst_log_error(message *C.char) {
	logger.Error(C.GoString(message))
}

//export gst_log_debug
func gst_log_debug(message *C.char)  {
	logger.Debug(C.GoString(message))
}

type Player struct {
	player  unsafe.Pointer
	playing bool
}

func NewPlayer() *Player {
	p := &Player{}
	p.player = C.player_new()
	p.playing = false
	return p
}


func gString(s string) *C.gchar {
	return (*C.gchar)(C.CString(s))
}

func gFree(s unsafe.Pointer) {
	C.g_free(C.gpointer(s))
}

func (p *Player) Play(uri string)  {
	s := gString(uri)
	C.player_play(p.player, s)
	gFree(unsafe.Pointer(s))
	p.playing = true
}

func (p *Player) Pause()  {
	C.player_pause(p.player)
}

func (p *Player) Resume()  {
	C.player_resume(p.player)
}

func (p *Player) TriggerPlay()  {
	p.playing = !p.playing
	if p.playing {
		p.Resume()
	}else {
		p.Pause()
	}
}

func (p *Player) Stop()  {
	C.player_stop(p.player)
}

func Init()  {
	C.loop_init()
}

func Start() {
	go func() {
		C.loop_run()
	}()
}

func Stop() {
	C.loop_destroy_and_quit()
}
