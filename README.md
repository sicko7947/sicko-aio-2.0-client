# Sicko AIO Client

Hi, I’m **Sicko** — many of you might know me from **SickoAIO**, a sneaker/AIO bot I launched back in **2021**. After a successful release of version 1.0, the scene started shifting fast. With constant CDN security updates and personal life catching up, I burned out. Over time, the passion faded, and I couldn’t keep up with the arms race anymore.

But after 4 years, I’ve decided to **open source the entire SickoAIO 2.0 framework**, fully rewritten in **Golang**. This is a **high-concurrency**, **scalable** system that may be useful when building performance-critical frameworks — whether for botting, automation, or other use cases.

---

## 🧱 Architecture Overview

The framework is modular and split into a few key components:

- `auth-grpc` – authentication and license verification layer  
- `server-grpc` – core orchestration and for getting cookies or bypassing security measures
- `client-backend` – backend logic for task handling and data flow  
- `client-frontend` – UI panel for task management and control  

Each part is designed to be highly concurrent and distributed-friendly.

---

## 🧠 Dev Notes & Insights

Over the years I’ve explored a lot of deep technical areas — from:

- CDN edge caching tricks
- Reverse engineering flows
- Custom TLS fingerprinting
- Many more...

I plan to **slowly release more repos, writeups, or code samples** on these areas over time.

> ❗ **Important Note:** I will **not** be open sourcing the specific logic I used to tackle Akamai's cookie protection — even if it's outdated — due to **legal and ethical reasons**.

---

## ⚠️ Disclaimer

- This is the **infrastructure behind a bot**, not a ready-to-go bot itself.
- You **will need to test, debug, and adapt** the code to your use case.
- This repo is under Apache2.0 License, you can basically do whatever the fuck you want.
---

# Prerequisite
- redis
- asynq

# Running
`go run main.go ` 

