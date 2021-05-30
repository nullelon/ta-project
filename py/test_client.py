import json, urllib.request

ul = "http://localhost:8080/"
with urllib.request.urlopen(ul) as url:
    data = json.loads(url.read().decode())
print(data)