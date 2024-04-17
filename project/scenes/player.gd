extends CharacterBody2D
class_name Player

var session_id: String#
var is_current_player: bool

func _enter_tree() -> void:
	$NameTag.text = session_id

func make_current() -> void:
	$Camera2D.enabled = true
	is_current_player = true
