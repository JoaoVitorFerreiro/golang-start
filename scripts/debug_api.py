import requests
import socket

def test_api_connection():
    urls_to_test = [
        "http://localhost:8080/health",
        "http://127.0.0.1:8080/health",
        "http://[::1]:8080/health"
    ]
    
    for url in urls_to_test:
        try:
            print(f"Testando: {url}")
            response = requests.get(url, timeout=5)
            print(f"✅ Sucesso! Status: {response.status_code}")
            print(f"Response: {response.text}")
            return url.replace("/health", "")  # Retorna base URL que funciona
        except Exception as e:
            print(f"❌ Falhou: {e}")
    
    return None

# Teste
working_url = test_api_connection()
if working_url:
    print(f"\n🎯 Use esta URL: {working_url}")