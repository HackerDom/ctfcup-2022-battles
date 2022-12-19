import requests

def main():
    with open('rom', 'rb') as file:
        data = file.read()
    
    r = requests.post('http://localhost:8080/prog', data=data)
    print(r.text)

if __name__ == '__main__':
    main()