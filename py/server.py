from bottle import route, run


#duh idk how to iterate over this
def start_server(basic_f, weighed_f):

    #symbol as ETH_USDT

    @route("/basic/<symbol>")
    def process(symbol=""):
        return basic_f()

    @route("/weighed/<symbol>")
    def process(symbol=""):
        return weighed_f()

    run(host='localhost', port=8080, debug=True)


if (__name__ == "__main__"):
    start_server()
