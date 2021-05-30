from bottle import run, route

def start_server(serverside_function):
    @route("/")
    def process():
        print("CONN")
        return serverside_function()


    run(host='localhost', port=8080, debug=True)
