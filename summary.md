Title : Google Gen AI Enterprise Commercial Cinematic Video Product
Role  : Expert Golang + Google Cloud + Next.js Developer

---

# 📋 Software Architecture & PRD: B2B AI Video Generation SaaS

## 1. Visi Produk

Sebuah platform SaaS (Software as a Service) B2B yang menyederhanakan pembuatan video komersial sinematik menggunakan AI. Platform ini menghilangkan kompleksitas prompt engineering bagi pengguna bisnis, menyediakan template berkualitas tinggi (khususnya untuk industri bernilai tinggi seperti layanan aviasi dan jet pribadi), dan menghasilkan output video langsung ke format siap pakai.

## 2. Pengalaman Pengguna (UI/UX)

- **Guided Prompt Builder (4-Step Wizard):** Tidak ada kotak teks kosong. Pengguna menggunakan form wizard 4 langkah (Template → Visuals → Audio → Review) untuk memilih elemen (suasana, pencahayaan, pergerakan kamera).
- **Industry-Specific Templates:** Backend menggabungkan input sederhana dari pengguna menjadi prompt teknis yang sangat mendetail. Template yang tersedia: Fashion, Accessories, PC & Gaming, Electronics, Automotive, Luxury Real Estate, Aviation.
- **Sistem Status UI Reaktif:** `Pending` ➡️ `In Queue` ➡️ `Processing` ➡️ `Completed`. Frontend melakukan polling ke backend setiap 3 detik untuk mendapatkan status terbaru.
- **Halaman yang Tersedia:** Home (wizard), Login, Register, History (riwayat job), Profile, Settings.

## 3. Sistem Bisnis & Manajemen "Kredit"

- **Sistem Kredit Internal:** Pengguna berlangganan untuk mendapatkan saldo "Kredit" bulanan.
- **Deduction Logic:** Tindakan berbeda memakan biaya kredit berbeda (contoh: render 5 detik = 10 kredit, Upscale ke 4K = 20 kredit).
- **Validasi Transaksional:** Backend selalu memvalidasi dan memotong saldo kredit di database menggunakan `SELECT ... FOR UPDATE` **sebelum** mengirimkan request ke AI API — mencegah double-spend.

## 4. Tumpukan Teknologi (Tech Stack) — Status Terkini

| Layer | Teknologi | Status |
|---|---|---|
| Frontend | **Next.js 16** / React 19 | ✅ Running |
| Backend | **Go Fiber v2** di Google Cloud Run | ✅ Scaffolded |
| Database | **Neon (PostgreSQL 17)** — serverless | ✅ Live, 4 tables created |
| Task Queue | **Google Cloud Tasks** (sementara: goroutine) | 🔄 Placeholder |
| AI Engine | **Google Vertex AI (Veo 3)** via `google.golang.org/genai` | 🔄 Stubbed |
| Storage | **Google Cloud Storage** (.mp4) | 🔄 Pending |
| Auth | **JWT** (`golang-jwt/jwt/v5`) | ✅ Implemented |

### Database (Neon)
- **Project ID:** `lingering-darkness-93698397`
- **Region:** `aws-ap-southeast-1` (Singapore)
- **Tables:** `users`, `credit_transactions`, `video_generation_jobs`, `invoices`

### Backend API Endpoints (Go Fiber)
| Method | Path | Auth | Status |
|---|---|---|---|
| POST | `/api/v1/auth/register` | ❌ | ✅ Done |
| POST | `/api/v1/auth/login` | ❌ | ✅ Done |
| GET | `/api/v1/me` | JWT | ✅ Done |
| POST | `/api/v1/jobs` | JWT | ✅ Done |
| GET | `/api/v1/jobs/:id` | JWT | ✅ Done |
| GET | `/api/v1/jobs` | JWT | 🔄 Planned |
| GET | `/api/v1/health` | ❌ | ✅ Done |

## 5. Alur Data & Konkurensi

1. **Request Masuk:** Pengguna menekan tombol "Generate Video" di Frontend Next.js.
2. **Validasi & Potong Kredit (Transaksional):** Backend Go menerima request dan mengeksekusi query ke PostgreSQL menggunakan `SELECT ... FOR UPDATE` (Serializable isolation). Ini memastikan tidak ada kredit yang terpotong ganda.
3. **Masuk Antrean:** Backend langsung melempar tugas ke Google Cloud Tasks (atau goroutine lokal) dan mengembalikan `202 Accepted { job_id, status: "in_queue" }` ke Frontend.
4. **Polling Status:** Frontend melakukan polling ke `GET /api/v1/jobs/:id` setiap 3 detik untuk mendapatkan update status.
5. **Drip-Feeding ke AI:** Cloud Tasks mengatur ritme pengerjaan ke Vertex AI sesuai batas kuota API.
6. **Penyimpanan:** Hasil render diupload ke Google Cloud Storage.
7. **Selesai:** Backend update status job menjadi `Completed` + set `gcs_video_url`, Frontend tampilkan video.

## 6. Struktur Proyek

```
project/
├── web/                        # Next.js 16 Frontend
│   └── src/app/
│       ├── page.tsx            # Main 4-step wizard generator
│       ├── login/              # Login page
│       ├── register/           # Register page
│       ├── history/            # Job history list
│       ├── profile/            # User profile
│       ├── settings/           # App settings
│       ├── context/AppContext.tsx
│       ├── Navbar.tsx
│       └── lib/                # [PLANNED] api.ts, auth.ts
│
└── backend/                    # Go Fiber Backend
    ├── cmd/api/main.go
    ├── internal/
    │   ├── config/config.go
    │   ├── database/db.go
    │   ├── middleware/auth.go
    │   ├── handlers/           # auth, user, job handlers
    │   ├── services/           # auth, credit, prompt, queue, video
    │   ├── models/models.go
    │   └── router/router.go
    └── migrations/001_initial_schema.sql
```

---

**Instruksi Developer — Langkah Selanjutnya:**
1. Integrasi Web ↔ Backend: buat `lib/api.ts`, `lib/auth.ts`, dan migrasikan `AppContext` dari mock ke real JWT auth.
2. Implementasi Veo API: aktifkan `video_service.go` dengan `google.golang.org/genai`.
3. Google Cloud Tasks: ganti goroutine placeholder dengan Cloud Tasks yang sebenarnya.
4. GCS Upload: implementasi upload hasil render ke Google Cloud Storage.
