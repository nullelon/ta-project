from dijkstar import Graph
import urllib.request, json
from datetime import date

export_name = "data/"+date.today().strftime("%d.%m.%y")
if __name__ == "__main__":
    graph = Graph()
    with urllib.request.urlopen("https://api.binance.com/api/v1/exchangeInfo") as url:
        data = json.loads(url.read().decode())

    for symbol in data["symbols"]:
        baseAsset, quoteAsset = symbol["baseAsset"], symbol["quoteAsset"]
        graph.add_edge(baseAsset, quoteAsset, f"{baseAsset}/{quoteAsset}")

        # here 'B\A' represents reversed pair relation - i.e. we need to sell B and buy A !BUT! using A/B pair
        graph.add_edge(quoteAsset, baseAsset, f"{quoteAsset}\{baseAsset}")

    graph.dump(export_name)

def load():
    return Graph.load(export_name)