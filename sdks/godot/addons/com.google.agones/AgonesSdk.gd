extends HTTPRequest

# This code is generated by go generate.
# DO NOT EDIT BY HAND!

signal on_request(path, params, method)

export (String) var api_endpoint = "http://localhost"


func _init():
	var agones_port = OS.get_environment("AGONES_SDK_HTTP_PORT")
	if ! agones_port:
		agones_port = 9358
	api_endpoint = "http://127.0.0.1:%s" % agones_port


# Retrieve the current GameServer data
func GetGameServer() -> Dictionary:
	return yield(_api_request("/gameserver", {}, HTTPClient.METHOD_GET), "completed")


# Apply a Label to the backing GameServer metadata
func SetLabel(body) -> Dictionary:
	return yield(
		_api_request(
			"/metadata/label",
			{
				"body": body,
			},
			HTTPClient.METHOD_PUT
		),
		"completed"
	)


# Call when the GameServer is ready
func Ready() -> Dictionary:
	return yield(_api_request("/ready", {}, HTTPClient.METHOD_POST), "completed")


# Send GameServer details whenever the GameServer is updated
func WatchGameServer() -> Dictionary:
	return yield(_api_request("/watch/gameserver", {}, HTTPClient.METHOD_GET), "completed")


# Call to self Allocation the GameServer
func Allocate() -> Dictionary:
	return yield(_api_request("/allocate", {}, HTTPClient.METHOD_POST), "completed")


# Send a Empty every d Duration to declare that this GameSever is healthy
func Health() -> Dictionary:
	return yield(_api_request("/health", {}, HTTPClient.METHOD_POST), "completed")


# Apply a Annotation to the backing GameServer metadata
func SetAnnotation(body) -> Dictionary:
	return yield(
		_api_request(
			"/metadata/annotation",
			{
				"body": body,
			},
			HTTPClient.METHOD_PUT
		),
		"completed"
	)


# Marks the GameServer as the Reserved state for Duration
func Reserve(body) -> Dictionary:
	return yield(
		_api_request(
			"/reserve",
			{
				"body": body,
			},
			HTTPClient.METHOD_POST
		),
		"completed"
	)


# Call when the GameServer is shutting down
func Shutdown() -> Dictionary:
	return yield(_api_request("/shutdown", {}, HTTPClient.METHOD_POST), "completed")


func _api_request(path: String, params: Dictionary, method = HTTPClient.METHOD_GET) -> Dictionary:
	emit_signal("on_request", path, params, method)

	# Build required request objects
	var request_string = "%s%s%s" % [api_endpoint, path, _params_to_string(params)]
	var headers: PoolStringArray = ["Content-Type: application/json"]

	# Make HTTP Requesst
	var error = self.request(request_string, headers, false, method)
	if error != OK:
		yield(get_tree().create_timer(0.001), "timeout")
		return _build_error_message("Agones Client encounted an error Godot error code: %s" % error)

	# Get and parse result
	var result = yield(self, "request_completed")
	if len(result) > 3:
		if result[1] == 200:
			var json: JSONParseResult = JSON.parse(result[3].get_string_from_utf8())
			if json.error:
				return _build_error_message("Failed to parse response: %s" % json.error_string)
			return json.result
		else:  # Return response code in error message if possible
			return _build_error_message(
				"Request failed! Response code: %s\n%s" % [str(result[1]), str(result[3])]
			)

	return _build_error_message("Request failed!")


# Helper function for converting a dictionary into HTTP parameters
func _params_to_string(params: Dictionary) -> String:
	var param_strings = []
	for param in params:
		param_strings.append("%s=%s" % [param, str(params[param])])

	var params_string = ""
	for i in range(param_strings.size()):
		if i == 0:
			params_string += "?"

		params_string += param_strings[i]

		if i != params.size():
			params_string += "&"
	return params_string


# Helper function for generating client errors
func _build_error_message(message):
	return {"message": message, "success": false}
