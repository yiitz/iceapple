package media

import "C"
import (
	"github.com/yiitz/iceapple/log"
	"sync/atomic"
	"time"
	"unsafe"
)

/*
//https://gstreamer.freedesktop.org/documentation/tutorials/playback/progressive-streaming.html
#cgo pkg-config: gstreamer-1.0
#include <stdio.h>
#include <stdlib.h>
#include <gst/gst.h>

extern void gst_log_error(char *);
extern void gst_log_debug(char *);
extern void player_refresh_ui(gint p0, gint64 p1, gint64 p2);

static char msgbuf[256];

static GMainLoop *loop;

typedef struct _CustomData {
	GstElement *pipeline;
	gint id;
	gint buffering_level;
	gboolean is_live;
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
	GstStreamStatusType stream_status;
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
		case GST_MESSAGE_STREAM_STATUS:
			gst_message_parse_stream_status(msg, &stream_status, NULL);
			log_call(gst_log_debug,"stream status changed: %d",stream_status);
			break;
	}
}

static void got_location(GstObject *gstobject, GstObject *prop_object, GParamSpec *prop, gpointer data) {
	gchar *location;
	g_object_get (G_OBJECT (prop_object), "temp-location", &location, NULL);
	snprintf(msgbuf, sizeof(msgbuf), "temporary file: %s", location);
	gst_log_debug (msgbuf);
	g_free (location);
	//g_object_set (G_OBJECT (prop_object), "temp-remove", FALSE, NULL);
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

static CustomData* player_new(gint id){
	GstBus *bus;
	CustomData *data;
	guint flags;

	data = malloc(sizeof(CustomData));
	data->id = id;
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

static void player_destroy(CustomData *data){
	gst_element_set_state(data->pipeline, GST_STATE_NULL);
	gst_object_unref(data->pipeline);
	free(data);
}

static void player_play(CustomData *data, gchar *path){
	GstStateChangeReturn ret;
	gst_element_set_state (data->pipeline, GST_STATE_NULL);
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

static void player_pause(CustomData *data){
	gst_element_set_state(data->pipeline, GST_STATE_PAUSED);
}

static void player_resume(CustomData *data){
	gst_element_set_state(data->pipeline, GST_STATE_PLAYING);
}

static void player_stop(CustomData *data){
	gst_element_set_state(data->pipeline, GST_STATE_NULL);
}

static void player_set_volume(CustomData *data, float vol)
{
	g_object_set(G_OBJECT(data->pipeline), "volume", vol, NULL);
}

static void player_query_progress(CustomData *data, gint64 *position, gint64 *duration){
	gst_element_query_position (data->pipeline, GST_FORMAT_TIME, position);
	gst_element_query_duration (data->pipeline, GST_FORMAT_TIME, duration);
}

static int player_get_state(CustomData *data){
	GstState state;
	gst_element_get_state(data->pipeline, &state, NULL, GST_SECOND*3);
	return state;
}

static void player_seek(CustomData *data, gint64 pos)
{
	gst_element_seek (data->pipeline, 1.0, GST_FORMAT_TIME, GST_SEEK_FLAG_FLUSH,
                         GST_SEEK_TYPE_SET, pos,
                         GST_SEEK_TYPE_NONE, GST_CLOCK_TIME_NONE);
}

*/
import "C"

var logger = log.NewLogger("media")

//export gst_log_error
func gst_log_error(message *C.char) {
	logger.Error(C.GoString(message))
}

//export gst_log_debug
func gst_log_debug(message *C.char) {
	logger.Debug(C.GoString(message))
}

var gid int32 = 0

type Player struct {
	player      *C.CustomData
	id          int32
	volume      int
	OnPlayStart func()
}

const GstStatePaused = 3
const GstStatePlaying = 4

var players = map[int32]*Player{}

func NewPlayer() *Player {
	p := &Player{}
	p.id = atomic.AddInt32(&gid, 1)
	p.player = C.player_new(C.gint(p.id))
	p.volume = 80
	players[p.id] = p

	C.player_set_volume(p.player, C.float(float32(p.volume)/100))
	return p
}

func ClosePlayer(p *Player) {
	C.player_destroy(p.player)
	delete(players, p.id)
}

func gString(s string) *C.gchar {
	return (*C.gchar)(C.CString(s))
}

func gFree(s unsafe.Pointer) {
	C.g_free(C.gpointer(s))
}

func (p *Player) Play(uri string) {
	logger.Debugf("play media:%s", uri)
	s := gString(uri)
	C.player_play(p.player, s)
	gFree(unsafe.Pointer(s))
	if p.OnPlayStart != nil {
		p.OnPlayStart()
	}
}

func (p *Player) Pause() {
	C.player_pause(p.player)
}

func (p *Player) Resume() {
	C.player_resume(p.player)
	if p.OnPlayStart != nil {
		p.OnPlayStart()
	}
}

func (p *Player) Seek(position int64) {
	C.player_seek(p.player, C.gint64(position))
}

func (p *Player) TriggerPlay() {
	logger.Debug("trigger play")
	state := p.GetState()
	switch state {
	case GstStatePaused:
		p.Resume()
	case GstStatePlaying:
		p.Pause()
	}
}

func (p *Player) VolumeUp() {
	p.volume += 10
	if p.volume > 100 {
		p.volume = 100
	}
	C.player_set_volume(p.player, C.float(float32(p.volume)/100))
}

func (p *Player) VolumeDown() {
	p.volume -= 10
	if p.volume < 0 {
		p.volume = 0
	}
	C.player_set_volume(p.player, C.float(float32(p.volume)/100))
}

func (p *Player) GetVolume() int {
	return p.volume
}

func (p *Player) GetProgress() (time.Duration, time.Duration) {
	var position C.gint64 = -1
	var duration C.gint64 = -1
	C.player_query_progress(p.player, &position, &duration)
	return time.Duration(position), time.Duration(duration)
}

func (p *Player) GetState() int {
	return int(C.player_get_state(p.player))
}

func (p *Player) Stop() {
	C.player_stop(p.player)
}

func Init() {
	C.loop_init()
	go func() {
		C.loop_run()
	}()
}

func Destroy() {
	C.loop_destroy_and_quit()
}
