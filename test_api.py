from requests import Session

net = Session()
net.headers['maraki'] = "online.helio3.maraki"

url1 = 'http://localhost/'

r = net.get(url1)
print(r.json())

