import httpx

params = {
    'num1': 1,
    'num2': 2
}

# use docker inspect to find ip of adder docker
url = 'http://172.17.0.1/add'

try:
    r = httpx.get(url, params=params)
    r.raise_for_status()
    print(r.text)
except httpx.RequestError as exc:
    print(f"Request error: {exc}")
except httpx.HTTPStatusError as exc:
    print(f"HTTP error: {exc}")
except Exception as exc:
    print(f"An unexpected error occurred: {exc}")

