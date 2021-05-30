import json

from server import start_server
from chains_pairs import get_weighed_chains, get_chains

# from, to
base, quote = "ETC", "UAH"


get_weighed_chains_wrap = lambda: json.dumps(get_weighed_chains(base, quote))
get_chains_wrap = lambda: json.dumps(get_chains(base, quote))

# @route("/weighed/")
# @route("/basic/")
endpoints = {
    "/weighed/": get_weighed_chains_wrap,
    "/basic/": get_chains_wrap
}

start_server(get_chains_wrap, get_weighed_chains_wrap)