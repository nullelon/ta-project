from dijkstar import find_path
import json
import urllib.request

import graph_constructor


def cost_func(u, v, edge, prev_edge, visited_edges):
    # should be less or equals than 1 / N of currencies
    visited_const = 0.001
    return 1 + visited_edges.count(edge) * visited_const


def cost_func_wrap(visited_edges):
    return lambda u, v, pu, pv: cost_func(u, v, pu, pv, visited_edges)


def get_pair_price(pair):
    reversed = "\\" in pair
    symbol = pair.replace("/", "") if not reversed else "".join(pair.split("\\")[::-1])
    ul = "https://api.binance.com/api/v3/ticker/bookTicker?symbol=" + symbol
    with urllib.request.urlopen(ul) as url:
        data = json.loads(url.read().decode())

        # lower, higher prices
        bidPrice, askPrice = float(data["bidPrice"]), float(data["askPrice"])

        if not reversed:
            return bidPrice
        else:
            return 1 / askPrice


# returns [['ETC/ETH', 'ETH/UAH'], ['ETC/BTC', 'BTC/UAH'], ['ETC/USDT', 'USDT/UAH'], ['ETC/BNB', 'BNB/UAH']]
def get_chains(base, quote):
    baseAsset, quoteAsset = base, quote
    graph = graph_constructor.load()
    visited_edges = []
    conversion_chains = []
    path = find_path(graph, baseAsset, quoteAsset, cost_func=cost_func_wrap(visited_edges))
    while path.edges not in conversion_chains:
        conversion_chains.append(path.edges)
        visited_edges.extend(path.edges)
        path = find_path(graph, baseAsset, quoteAsset, cost_func=cost_func_wrap(visited_edges))

    return conversion_chains

# returns array of dicts of format {'chain': ['ETC/ETH', 'ETH/UAH'], 'rate': 1811.1134006050681}
def get_weighed_chains(base, quote):
    # from, to
    baseAsset, quoteAsset = base, quote
    trade_fee = 1 - 0.001

    weighed_chains = []
    for chain in get_chains(base, quote):
        rate = 1
        for pair in chain:
            rate *= get_pair_price(pair) * trade_fee

        weighed_chain = {"chain": chain, "rate": rate}
        print(weighed_chain)
        weighed_chains.append(weighed_chain)

    return sorted(weighed_chains, key=lambda chain: chain["rate"], reverse=True)