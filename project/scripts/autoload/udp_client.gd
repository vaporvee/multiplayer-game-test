extends Node

signal player_move(session_id: String, position_mod: Vector2)
signal player_join(player_client: Dictionary, clients: Array)
signal player_disconnect(session_id: String)

var udp_peer: PacketPeerUDP = PacketPeerUDP.new()
var server_ip: String = "127.0.0.1" # Replace with your server's IP
var server_port: int = 4477 # Replace with your server's port

func _ready():
	udp_peer.bind(server_port, "*")
	print("Listening on port: ", server_port)
	# Send a message to the server to initiate session ID assignment
	var payload: Dictionary = {
		"type": "init",
	}
	send_payload(payload)

func _process(_delta):
	if udp_peer.get_available_packet_count() > 0:
		var packet = udp_peer.get_packet()
		var packet_json: Dictionary = JSON.parse_string(packet.get_string_from_utf8())
		match packet_json.type:
			"init_success":
				player_join.emit(packet_json.player_client, Array(packet_json.clients))
			"message":
				print(packet_json.msg)
			"move":
				player_move.emit(packet_json.session_id, Vector2(packet_json.direction.x,packet_json.direction.y))
			"disconnect":
				player_disconnect.emit(packet_json.session_id)

func send_message(message: String) -> void:
	var payload: Dictionary = {
		"type": "message",
		"msg": message
	}
	send_payload(payload)

func send_payload(payload: Dictionary):
	var packet = JSON.stringify(payload).to_utf8_buffer()
	udp_peer.set_dest_address(server_ip, server_port)
	udp_peer.put_packet(packet)

func move(position_mod: Vector2):
	var payload: Dictionary = {
		"type": "move",
		"direction": {
			"x": position_mod.x,
			"y": position_mod.y
		}
	}
	send_payload(payload)

func _exit_tree():
	var payload: Dictionary = {
		"type": "disconnect"
	}
	send_payload(payload)
