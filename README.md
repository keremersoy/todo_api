# TODO API

Temel CRUD işlemleri ve JWT ile kimlik doğrulama işlemi gerçekleştiren bir API. `bun` ORM ile PostgreSQL veritabanını kullanmaktadır ve Fiber üzerinde çalışır.

## Gereksinimler

- Go programlama dili
- PostgreSQL veritabanı
- [Fiber](https://gofiber.io/)

## Kurulum

1. Bu projeyi klonlayın.
2. `.env` dosyasını oluşturun ve PostgreSQL veritabanı bağlantı ayarlarınızı ekleyin.
    ```
    DB_HOST=localhost
    DB_PORT=5432
    DB_USER=your_username
    DB_PASSWORD=your_password
    DB_SSLMODE=disable
    DB_DBNAME=your_database_name
    ```
3. Gerekli bağımlılıkları yükleyin:
    ```
    go get -d
    ```
4. Uygulamayı başlatın:
    ```
    go get -d
    ```
Uygulama varsayılan olarak localhost:9090 üzerinde çalışacaktır.

## API Endpoints

#### Kullanıcı İşlemleri
- **POST /user/login:** Kullanıcı girişi
- **POST /user/register:** Kullanıcı kaydı
- **GET /user:** Tüm kullanıcıları listeler.
- **GET /user/info:** Aktif kullanıcının bilgilerini getirir.

#### TODO İşlemleri
- **GET /todos:** Tüm TODO öğelerini listeler.
- **POST /todos:** Yeni bir TODO öğesi ekler.
- **GET /todos/{id}:** Belirli bir TODO öğesini getirir.
- **PATCH /todos/{id}:** Belirli bir TODO öğesinin durumunu değiştirir.
- **DELETE /todos/{id}:** Belirli bir TODO öğesini siler.

