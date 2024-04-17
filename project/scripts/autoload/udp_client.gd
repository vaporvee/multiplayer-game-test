extends Node

var udp_peer: PacketPeerUDP = PacketPeerUDP.new()
var server_ip: String = "45.93.249.177" # Replace with your server's IP
var server_port: int = 4477 # Replace with your server's port

func _ready():
	var err = udp_peer.bind(server_port, "*")
	if err != OK:
		print("Failed to bind to port: ", server_port)
		return
	print("Listening on port: ", server_port)
	# Send a message to the server to initiate session ID assignment
	var payload: Dictionary = {
		"type": "init",
	}
	send_payload(payload)

func _process(_delta):
	if udp_peer.get_available_packet_count() > 0:
		var packet = udp_peer.get_packet()
		var message: String = packet.get_string_from_utf8()
		print("Received: ", message)

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

func move(vec: Vector2i):
	var payload: Dictionary = {
		"type": "move",
		"direction": {
			"x": vec.x,
			"y": vec.y
		}
	}
	send_payload(payload)
