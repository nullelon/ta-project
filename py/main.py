import json

from server import start_server
from chains_pairs import get_weighed_chains

# from, to
base, quote = "ETC", "UAH"

serverside_function = lambda: json.dumps(get_weighed_chains(base, quote))

start_server(serverside_function)