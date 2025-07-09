# Criar versão corrigida do script de performance
import requests
import threading
import time
import json
from concurrent.futures import ThreadPoolExecutor, as_completed
import statistics
import uuid

class GoAPIPerformanceTest:
    def __init__(self, base_url="http://localhost:8080"):
        self.base_url = base_url
        self.results = []
        self.created_users = []
        
    def test_connection(self):
        """Testa conexão com a API"""
        print("🔍 Testando conexão com a API...")
        
        # Testa diferentes URLs
        urls_to_test = [
            f"{self.base_url}/health",
            f"http://127.0.0.1:8080/health",
            f"http://[::1]:8080/health"
        ]
        
        for url in urls_to_test:
            try:
                print(f"   Tentando: {url}")
                response = requests.get(url, timeout=5)
                if response.status_code == 200:
                    print(f"   ✅ Sucesso! Status: {response.status_code}")
                    print(f"   Response: {response.text}")
                    # Atualiza base_url para a URL que funciona
                    self.base_url = url.replace("/health", "")
                    return True
                else:
                    print(f"   ⚠️  Status: {response.status_code}")
            except Exception as e:
                print(f"   ❌ Erro: {e}")
        
        return False
        
    def create_user_data(self, user_id):
        """Gera dados de usuário para teste"""
        return {
            "name": f"Usuario{user_id}",
            "email": f"user{user_id}_{int(time.time())}_{user_id}@test.com"
        }
    
    def make_create_request(self, user_data):
        """Faz uma requisição HTTP POST para criar usuário"""
        try:
            start_time = time.time()
            response = requests.post(
                f"{self.base_url}/users",
                json=user_data,
                headers={"Content-Type": "application/json"},
                timeout=10
            )
            end_time = time.time()
            
            # Debug: print first few responses
            if len(self.created_users) < 3:
                print(f"   Debug - Status: {response.status_code}, Response: {response.text[:100]}")
            
            # Armazena ID do usuário criado para cleanup
            success = False
            if response.status_code in [200, 201]:
                try:
                    user_response = response.json()
                    if "id" in user_response:
                        self.created_users.append(user_response["id"])
                        success = True
                except:
                    pass
            
            return {
                "success": success,
                "status_code": response.status_code,
                "response_time": end_time - start_time,
                "endpoint": "POST /users"
            }
        except Exception as e:
            return {
                "success": False,
                "status_code": 0,
                "response_time": 0,
                "endpoint": "POST /users",
                "error": str(e)
            }
    
    def make_get_request(self, user_id=None):
        """Faz uma requisição HTTP GET"""
        try:
            start_time = time.time()
            
            if user_id:
                url = f"{self.base_url}/users/{user_id}"
                endpoint = "GET /users/:id"
            else:
                url = f"{self.base_url}/users"
                endpoint = "GET /users"
            
            response = requests.get(url, timeout=10)
            end_time = time.time()
            
            return {
                "success": response.status_code == 200,
                "status_code": response.status_code,
                "response_time": end_time - start_time,
                "endpoint": endpoint
            }
        except Exception as e:
            return {
                "success": False,
                "status_code": 0,
                "response_time": 0,
                "endpoint": endpoint,
                "error": str(e)
            }
    
    def run_simple_test(self, num_requests=100, num_threads=5):
        """Executa teste simples de criação"""
        print(f"\\n{'='*50}")
        print(f"🚀 TESTE SIMPLES DE CRIAÇÃO")
        print(f"Requisições: {num_requests}")
        print(f"Threads: {num_threads}")
        print(f"URL: {self.base_url}")
        print(f"{'='*50}")
        
        # Teste manual primeiro
        print("\\n🧪 Teste manual primeiro...")
        test_user = self.create_user_data(999)
        result = self.make_create_request(test_user)
        print(f"Resultado teste manual: {result}")
        
        if not result["success"]:
            print("❌ Teste manual falhou! Verifique a API.")
            return None
        
        print("✅ Teste manual OK! Iniciando teste de performance...")
        
        results = []
        start_time = time.time()
        
        with ThreadPoolExecutor(max_workers=num_threads) as executor:
            futures = []
            for i in range(num_requests):
                user_data = self.create_user_data(i)
                future = executor.submit(self.make_create_request, user_data)
                futures.append(future)
            
            for future in as_completed(futures):
                result = future.result()
                results.append(result)
                
                if len(results) % 20 == 0:
                    print(f"📊 Processadas: {len(results)}/{num_requests}")
        
        return self._analyze_results(results, start_time, "CREATE USERS")
    
    def _analyze_results(self, results, start_time, test_name):
        """Analisa e exibe os resultados do teste"""
        end_time = time.time()
        total_time = end_time - start_time
        
        successful_requests = [r for r in results if r["success"]]
        failed_requests = [r for r in results if not r["success"]]
        
        response_times = [r["response_time"] for r in successful_requests]
        
        print(f"\\n{'='*20} RESULTADOS {test_name} {'='*20}")
        print(f"⏱️  Tempo total: {total_time:.2f}s")
        print(f"✅ Requisições bem-sucedidas: {len(successful_requests)}")
        print(f"❌ Requisições falharam: {len(failed_requests)}")
        print(f"📈 Taxa de sucesso: {(len(successful_requests)/len(results))*100:.1f}%")
        
        if successful_requests:
            rps = len(successful_requests)/total_time
            print(f"🚀 Requests/segundo: {rps:.2f}")
            print(f"⚡ Tempo médio de resposta: {statistics.mean(response_times)*1000:.2f}ms")
            print(f"🏃 Tempo mínimo: {min(response_times)*1000:.2f}ms")
            print(f"🐌 Tempo máximo: {max(response_times)*1000:.2f}ms")
            if len(response_times) > 1:
                print(f"📊 Mediana: {statistics.median(response_times)*1000:.2f}ms")
        
        if failed_requests:
            print(f"\\n❌ Erros encontrados:")
            error_counts = {}
            for req in failed_requests:
                error_key = f"Status {req['status_code']}"
                if 'error' in req:
                    error_key += f" - {req['error']}"
                error_counts[error_key] = error_counts.get(error_key, 0) + 1
            
            for error, count in error_counts.items():
                print(f"   {error}: {count} vezes")
        
        return {
            "test_name": test_name,
            "total_time": total_time,
            "successful_requests": len(successful_requests),
            "failed_requests": len(failed_requests),
            "requests_per_second": len(successful_requests)/total_time if total_time > 0 else 0,
            "avg_response_time": statistics.mean(response_times) if response_times else 0,
            "response_times": response_times
        }
    
    def cleanup_users(self):
        """Remove usuários criados durante o teste"""
        if not self.created_users:
            print("\\n🧹 Nenhum usuário para limpar")
            return
            
        print(f"\\n🧹 Limpando {len(self.created_users)} usuários criados...")
        
        deleted = 0
        for user_id in self.created_users:
            try:
                response = requests.delete(f"{self.base_url}/users/{user_id}", timeout=5)
                if response.status_code in [200, 204]:
                    deleted += 1
            except:
                pass
        
        print(f"✅ {deleted} usuários removidos")

def main():
    print("🚀 TESTE DE PERFORMANCE - GO USER API (VERSÃO DEBUG)")
    
    tester = GoAPIPerformanceTest()
    
    # Primeiro testa conexão
    if not tester.test_connection():
        print("❌ Não foi possível conectar com a API!")
        return
    
    print(f"\\n✅ Conectado com sucesso! URL: {tester.base_url}")
    
    try:
        # Executa teste simples
        result = tester.run_simple_test(100, 5)
        
        if result:
            print("\\n🎯 TESTE CONCLUÍDO COM SUCESSO!")
            print(f"RPS: {result['requests_per_second']:.2f}")
        
    finally:
        # Cleanup
        tester.cleanup_users()

if __name__ == "__main__":
    main()