extends Node

func _on_button_pressed():
	UdpClient.send_message("TEST BUTTON PRESSED")
