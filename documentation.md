# Документация
## Объекты
### Product
* **Описание**: Содержит информацию о товаре
* **Поля**:
    * **Title (string)**: Название товара
    * **ID (uint64)**: Код товара. Является уникальным идентификатором товара.
    * **Category (string)**: Категория товара
    
## Методы
### createProduct
* **Описание**: Добавляет товар в список 
* **HttpMethod**: POST
* **URL**: /products
### deleteProduct
* **Описание**: Удаляет товар с указанным ID из списка 
* **HttpMethod**: DELETE
* **URL**: /products/{id}
### getProducts
* **Описание**: Возвращает список всех товаров
* **HttpMethod**: GET
* **URL**: /products
### getProduct
* **Описание**: Возвращает товар с указанным ID
* **HttpMethod**: GET
* **URL**: /products/{id}
### updateProduct
* **Описание**: Обновляет информацию о товаре
* **HttpMethod**: PUT
* **URL**: /products