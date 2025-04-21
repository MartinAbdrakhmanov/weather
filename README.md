# Weather

## How to Run

1. Fill in your credentials in `example.env` and rename it to `.env`.  
2. Then run the following commands:

```bash
go mod tidy
go run main.go
```

By default, the server will start at `localhost:8080`.

---

## What It Does

For each city, LLaMA 3 generates a clothing suggestion based on the current weather.  
The first generation may take a few seconds, but the result is cached in Redis:

- If the same city is requested again, the suggestion is returned instantly.
- Each suggestion is stored for 1 hour, after which it will be regenerated upon the next request.

![Weather UI Preview](https://github.com/user-attachments/assets/5d0aa6c9-97de-45de-b62b-6cc6c1f39ca7)

---

## 📄 Notes

The **Report** and **Form** tabs are not currently in use.  
However, if you upload any `.md` file to the `data` folder, it will be accessible at `/report`.
