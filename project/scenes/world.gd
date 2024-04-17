extends Node2D


func _ready():
	UdpClient.player_join.connect(_on_player_join)
	UdpClient.player_disconnect.connect(_on_player_disconnect)

func _on_player_join(player_client: Dictionary, clients: Array) -> void:
	add_player(player_client, true)
	for client in clients:
		add_player(client, false)


func add_player(client: Dictionary, is_current_player: bool):
	var player: Player = preload("res://scenes/player.tscn").instantiate()
	player.session_id = client.session_id
	player.position = Vector2(client.position.x, client.position.y)
	add_child(player)
	if is_current_player:
		player.make_current()

func _on_button_pressed():
	var payload: Dictionary = {
		"type": "init",
	}
	UdpClient.send_message("test")

func _on_player_disconnect(session_id: String):
	for child: Node in get_children():
		if child is Player && child.session_id == session_id:
			child.queue_free()
