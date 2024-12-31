from requests import Session

net = Session()
net.headers['maraki'] = "online.helio3.maraki"

url1 = 'https://maraki-api-production.up.railway.app/'
url2 = 'https://maraki-api-production.up.railway.app/serie/arcane/temporada-1/1789'

r = net.get(url2)
print(r.json())

