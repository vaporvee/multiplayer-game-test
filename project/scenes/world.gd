extends Node2D


func _ready():
	UdpClient.player_join.connect(_on_player_join)
	UdpClient.player_disconnect.connect(_on_player_disconnect)

func _on_player_join(player_client: Dictionary, clients: Array) -> void:
	add_player(player_client)
	for client in clients:
		add_player(client)


func add_player(client: Dictionary):
	var player: Player = preload("res://scenes/player.tscn").instantiate()
	player.session_id = client.session_id
	player.position = Vector2(client.position.x, client.position.y)
	add_child(player)

func _on_button_pressed():
	var payload: Dictionary = {
		"type": "init",
	}
	UdpClient.send_message("test")

func _on_player_disconnect(session_id: String):
	for child: Node in get_children():
		if child is Player && child.session_id == session_id:
			child.queue_free()
