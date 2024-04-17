extends Control

var websocket_url = "wss://acecore.lol:4477/ws"
var websocket_peer = WebSocketPeer.new()

var payload: Dictionary = {
	"type": "broadcast",
	"msg": "TESTBUTTON_PRESSED"
}

func _ready():
	var error = websocket_peer.connect_to_url(websocket_url)
	if error == OK:
		print("Connected to WebSocket server")
		send_message("Connected client")
	else:
		print("Failed to connect to WebSocket server")

func _process(_delta):
	websocket_peer.poll()
	var state = websocket_peer.get_ready_state()
	if state == WebSocketPeer.STATE_OPEN:
		while websocket_peer.get_available_packet_count():
			print("Packet: ", websocket_peer.get_packet().get_string_from_utf8())
	elif state == WebSocketPeer.STATE_CLOSING:
		# Keep polling to achieve proper close.
		pass
	elif state == WebSocketPeer.STATE_CLOSED:
		var code = websocket_peer.get_close_code()
		var reason = websocket_peer.get_close_reason()
		print("WebSocket closed with code: %d, reason %s. Clean: %s" % [code, reason, code != -1])
		set_process(false)

func send_message(message: String) -> void:
	var state = websocket_peer.get_ready_state()
	if state == WebSocketPeer.STATE_OPEN:
		websocket_peer.send(str(payload).to_utf8_buffer())
	else:
		print("WebSocket connection is not open. Current state: ", state)

func _on_button_pressed():
	send_message("Test button pressed")
