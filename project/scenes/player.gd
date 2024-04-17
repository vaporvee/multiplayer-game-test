extends CharacterBody2D
class_name Player

var session_id: String

func _enter_tree() -> void:
	$NameTag.text = session_id
