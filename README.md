# Quoter

# 1. Добавить цитату (POST)
curl -X POST http://localhost:8080/quotes \
  -H "Content-Type: application/json" \
  -d '{"author":"Confucius", "quote":"Life is simple, but we insist on making it complicated."}'

# 2. Получить все цитаты (GET)
curl http://localhost:8080/quotes

# 3. Получить случайную цитату (GET)
curl http://localhost:8080/quotes/random

# 4. Фильтрация по автору (GET)
curl http://localhost:8080/quotes?author=Confucius

# 5. Удалить цитату (DELETE)
curl -X DELETE http://localhost:8080/quotes/1