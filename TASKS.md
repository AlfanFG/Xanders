# 🎬 Xanders AI Video Platform — Master Developer Task List

> **For AI Developer Agents**: This document is the single source of truth for all remaining work on the Xanders B2B AI Cinematic Video Generation SaaS. Read the entire document before starting. Execute tasks **in order within each phase**. Tasks marked ✅ are already done — do not re-implement them. Tasks marked 🔲 need to be built.

---

## 🗂️ Project Context

| Item | Value |
|---|---|
| **Module** | `xanders-gen-video` |
| **Backend** | Go Fiber v2 · `e:/Work/Project/xanders/project/backend/` |
| **Frontend** | Next.js 16 + React 19 · `e:/Work/Project/xanders/project/web/` |
| **Database** | Neon PostgreSQL 17 · Project ID: `lingering-darkness-93698397` · Region: `aws-ap-southeast-1` |
| **AI Engine** | Google Gen AI SDK (`google.golang.org/genai v1.62.0`) — Veo 3 model |
| **Auth** | Stateless JWT via `golang-jwt/jwt/v5` |
| **DB Driver** | `pgx/v5` + `pgxpool` |

---

## ✅ Already Completed (Do Not Re-implement)

### Backend
- ✅ `go.mod` scaffolded as `xanders-gen-video` with all dependencies installed
- ✅ `internal/config/config.go` — env config loader
- ✅ `internal/database/db.go` — `pgxpool` connection to Neon
- ✅ `internal/models/models.go` — `User`, `Job`, `CreditTransaction` structs
- ✅ `internal/middleware/auth.go` — JWT Bearer token validation middleware
- ✅ `internal/services/auth_service.go` — bcrypt + JWT register/login
- ✅ `internal/services/credit_service.go` — `SELECT FOR UPDATE` transactional credit deduction
- ✅ `internal/services/prompt_service.go` — template → assembled cinematic prompt
- ✅ `internal/services/queue_service.go` — goroutine placeholder (not real Cloud Tasks yet)
- ✅ `internal/services/video_service.go` — **STUBBED** (Veo calls commented out)
- ✅ `internal/handlers/auth_handler.go` — `Register`, `Login`
- ✅ `internal/handlers/user_handler.go` — `GetProfile`
- ✅ `internal/handlers/job_handler.go` — `CreateJob`, `GetJob`
- ✅ `internal/router/router.go` — all routes registered
- ✅ `cmd/api/main.go` — Fiber app bootstrap, graceful shutdown
- ✅ `migrations/001_initial_schema.sql` — full schema
- ✅ Neon DB live with 4 tables: `users`, `credit_transactions`, `video_generation_jobs`, `invoices`
- ✅ `.env` with live Neon `DATABASE_URL`

### Frontend
- ✅ Next.js 16 app with 6 pages: Home (wizard), Login, Register, History, Profile, Settings
- ✅ 4-step video wizard UI (Template → Visuals → Audio → Review)
- ✅ `AppContext.tsx` — theme, lang, user state (currently **mocked — needs real auth**)
- ✅ `Navbar.tsx` — scrollable navbar with credit badge and avatar
- ✅ `globals.css` — design system (dark mode, glass panels, gradient text)
- ✅ All page CSS modules

---

## 🔲 PHASE 1 — Backend: Complete Missing Endpoints & Fix Config

> **Priority: CRITICAL** — Frontend cannot integrate without these.

### Task 1.1 — Fix CORS for Local Development
**File:** `backend/cmd/api/main.go`
**What:** Replace the generic `cors.New()` with explicit origin config.
```go
app.Use(cors.New(cors.Config{
    AllowOrigins:     "http://localhost:3000",
    AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
    AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
    AllowCredentials: true,
}))
```
**Why:** Without this, all browser fetch calls from Next.js will be blocked by CORS policy.

---

### Task 1.2 — Add `GET /api/v1/jobs` (List User Jobs)
**File:** `backend/internal/handlers/job_handler.go`
**What:** Add a new `ListJobs` handler that returns paginated list of jobs for the authenticated user.
```go
func (h *JobHandler) ListJobs(c *fiber.Ctx) error {
    userIDStr := c.Locals("user_id").(string)
    page, _ := strconv.Atoi(c.Query("page", "1"))
    limit, _ := strconv.Atoi(c.Query("limit", "20"))
    offset := (page - 1) * limit

    rows, err := database.DB.Query(context.Background(),
        `SELECT id, user_id, status, prompt_input, assembled_prompt, model_name,
         duration_seconds, credits_charged, gcs_video_url, error_message, created_at, updated_at
         FROM video_generation_jobs
         WHERE user_id = $1
         ORDER BY created_at DESC
         LIMIT $2 OFFSET $3`,
        userIDStr, limit, offset,
    )
    // scan rows into []models.Job → return c.JSON(jobs)
}
```
**Also register in router:** `jobs.Get("/", jobHandler.ListJobs)`

---

### Task 1.3 — Add `GET /api/v1/me/credits` (Credit Transaction History)
**File:** `backend/internal/handlers/user_handler.go`
**What:** Return paginated `credit_transactions` for the current user.
**Route:** `me.Get("/credits", userHandler.GetCreditHistory)`
```go
func (h *UserHandler) GetCreditHistory(c *fiber.Ctx) error {
    userIDStr := c.Locals("user_id").(string)
    // SELECT id, amount, reason, job_id, created_at
    // FROM credit_transactions WHERE user_id=$1 ORDER BY created_at DESC LIMIT 20
}
```

---

### Task 1.4 — Create Webhook Handler for Job Completion
**File:** `backend/internal/handlers/webhook_handler.go` — **[CREATE NEW]**
**What:** Internal endpoint to receive callbacks when a video finishes rendering (called by the goroutine worker or Cloud Tasks).
```go
// POST /api/v1/webhooks/job-complete
// Body: { "job_id": "uuid", "status": "completed|failed", "gcs_url": "https://..." }
// Auth: check X-Internal-Secret header against a config env var
func WebhookJobComplete(c *fiber.Ctx) error {
    // 1. Verify X-Internal-Secret header == cfg.InternalSecret
    // 2. Parse body
    // 3. UPDATE video_generation_jobs SET status=$1, gcs_video_url=$2, updated_at=NOW() WHERE id=$3
    // 4. Return 200 OK
}
```
**Route:** `api.Post("/webhooks/job-complete", webhookHandler.JobComplete)` — no JWT required.

---

### Task 1.5 — Run & Verify Backend Compiles
```bash
cd backend
go build ./...
go vet ./...
go run ./cmd/api/main.go
# Test: curl http://localhost:8080/health
# Expected: {"status":"ok","message":"API is running"}
```

---

## 🔲 PHASE 2 — Frontend: API Client & Auth Layer

> **Priority: CRITICAL** — All page wiring depends on this foundation.

### Task 2.1 — Create `.env.local`
**File:** `web/.env.local` — **[CREATE NEW]**
```env
NEXT_PUBLIC_API_URL=http://localhost:8080
```

---

### Task 2.2 — Create API Client
**File:** `web/src/lib/api.ts` — **[CREATE NEW]**
**What:** A single centralized fetch wrapper. All API calls must go through this file — never use raw `fetch` in components.

```typescript
const API_BASE = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080"

async function request<T>(path: string, options?: RequestInit): Promise<T> {
  const res = await fetch(`${API_BASE}${path}`, {
    headers: { "Content-Type": "application/json", ...options?.headers },
    ...options,
  })
  if (!res.ok) {
    const err = await res.json().catch(() => ({}))
    throw new Error(err.error ?? `HTTP ${res.status}`)
  }
  return res.json()
}

export const api = {
  // Auth
  register: (email: string, password: string) =>
    request<{ user: User; token: string }>("/api/v1/auth/register", {
      method: "POST", body: JSON.stringify({ email, password })
    }),

  login: (email: string, password: string) =>
    request<{ user: User; token: string }>("/api/v1/auth/login", {
      method: "POST", body: JSON.stringify({ email, password })
    }),

  // User
  getMe: (token: string) =>
    request<User>("/api/v1/me", { headers: { Authorization: `Bearer ${token}` } }),

  getCreditHistory: (token: string) =>
    request<CreditTransaction[]>("/api/v1/me/credits", { headers: { Authorization: `Bearer ${token}` } }),

  // Jobs
  createJob: (token: string, promptInput: Record<string, string>) =>
    request<{ job_id: string; status: string }>("/api/v1/jobs", {
      method: "POST",
      headers: { Authorization: `Bearer ${token}` },
      body: JSON.stringify({ prompt_input: promptInput }),
    }),

  getJob: (token: string, jobId: string) =>
    request<Job>(`/api/v1/jobs/${jobId}`, { headers: { Authorization: `Bearer ${token}` } }),

  listJobs: (token: string, page = 1) =>
    request<Job[]>(`/api/v1/jobs?page=${page}&limit=20`, { headers: { Authorization: `Bearer ${token}` } }),
}

// Shared types
export type User = { id: string; email: string; plan: string; credit_balance: number; created_at: string }
export type Job = {
  id: string; user_id: string; status: string; prompt_input: Record<string, string>;
  assembled_prompt: string; model_name: string; duration_seconds: number;
  credits_charged: number; gcs_video_url: string; error_message: string;
  created_at: string; updated_at: string;
}
export type CreditTransaction = { id: string; amount: number; reason: string; job_id: string; created_at: string }
```

---

### Task 2.3 — Create Auth Token Helpers
**File:** `web/src/lib/auth.ts` — **[CREATE NEW]**
```typescript
const TOKEN_KEY = "xanders_jwt"

export const saveToken = (token: string): void => localStorage.setItem(TOKEN_KEY, token)
export const getToken = (): string | null => localStorage.getItem(TOKEN_KEY)
export const removeToken = (): void => localStorage.removeItem(TOKEN_KEY)
export const isLoggedIn = (): boolean => !!getToken()
```

---

### Task 2.4 — Migrate `AppContext.tsx` to Real Auth
**File:** `web/src/app/context/AppContext.tsx` — **[MODIFY]**
**What:** Replace the hardcoded mock user with real JWT-based auth state.

**Key changes:**
1. Import `api`, `User` from `lib/api.ts`; import `getToken`, `saveToken`, `removeToken` from `lib/auth.ts`
2. Add `token: string | null` state, `isLoading: boolean` state
3. Add `logout()` function: calls `removeToken()`, sets `user = null`, sets `token = null`
4. On mount (`useEffect`): read token from localStorage → call `api.getMe(token)` → set user state
5. Remove the hardcoded default user `{ name: "John Doe", email: "john@acme.com", credits: 250 }`
6. Expose `token`, `setToken`, `logout`, `isLoading` from context

**New context type:**
```typescript
type AppContextType = {
  theme: "dark" | "light"; setTheme: (t: "dark" | "light") => void
  lang: string; setLang: (l: string) => void
  user: User | null; setUser: (u: User | null) => void
  token: string | null; setToken: (t: string | null) => void
  logout: () => void
  isLoading: boolean
}
```

---

### Task 2.5 — Wire Login Page
**File:** `web/src/app/login/page.tsx` — **[MODIFY]**
**Current state:** Static form — the button does nothing (type="button").

**Implementation:**
1. Add `"use client"` directive
2. Add `useState` for `email`, `password`, `error: string`, `isLoading: boolean`
3. Use `useAppContext()` to get `setUser`, `setToken`
4. Use `useRouter()` from `next/navigation` for redirect
5. Convert `<form>` to controlled inputs
6. Handle submit:
```typescript
const handleSubmit = async (e: React.FormEvent) => {
  e.preventDefault()
  setIsLoading(true); setError("")
  try {
    const data = await api.login(email, password)
    saveToken(data.token)
    setToken(data.token)
    setUser(data.user)
    router.push("/")
  } catch (err: any) {
    setError(err.message ?? "Login failed. Check your credentials.")
  } finally {
    setIsLoading(false)
  }
}
```
7. Show red error message below submit button when `error !== ""`
8. Disable button + show "Signing in..." when `isLoading === true`

---

### Task 2.6 — Wire Register Page
**File:** `web/src/app/register/page.tsx` — **[MODIFY]**
**Current state:** Static form with Full Name, Company, Email, Password.

**Important note:** The backend `POST /api/v1/auth/register` only accepts `email` + `password`. The `Full Name` and `Company` fields do not have DB columns yet. For now: only send `email` + `password` to the API. Add a `// TODO: add name/company to backend` comment.

**Implementation:**
1. Add `"use client"` + state for all fields + `error` + `isLoading`
2. Handle submit: call `api.register(email, password)` → save token → `router.push("/")`
3. Show error if email already exists (backend returns `"email already exists"`)
4. Disable button during loading

---

### Task 2.7 — Create Auth Guard Component
**File:** `web/src/components/AuthGuard.tsx` — **[CREATE NEW]**
**What:** HOC that redirects unauthenticated users to `/login`.

```typescript
"use client"
import { useEffect } from "react"
import { useRouter } from "next/navigation"
import { useAppContext } from "@/app/context/AppContext"

export function AuthGuard({ children }: { children: React.ReactNode }) {
  const { user, isLoading } = useAppContext()
  const router = useRouter()

  useEffect(() => {
    if (!isLoading && !user) router.push("/login")
  }, [user, isLoading, router])

  if (isLoading) return (
    <div style={{ display:'flex', justifyContent:'center', alignItems:'center', height:'100vh' }}>
      <div style={{ color:'var(--text-secondary)' }}>Loading...</div>
    </div>
  )
  if (!user) return null
  return <>{children}</>
}
```

---

## 🔲 PHASE 3 — Frontend: Wire All Pages to Real API

### Task 3.1 — Wire Main Generator (`page.tsx`)
**File:** `web/src/app/page.tsx` — **[MODIFY]**
**What:** Replace the fake `setTimeout` timers with real API calls + polling.

**Step 1 — Delete the fake generate function:**
```typescript
// DELETE THIS entire block from handleGenerate:
setTimeout(() => setStatus("queue"), 1500)
setTimeout(() => setStatus("processing"), 4000)
setTimeout(() => { setStatus("completed"); setVideoUrl("https://...sample...") }, 9000)
```

**Step 2 — Add real implementation:**
```typescript
const { token } = useAppContext()
const [currentJobId, setCurrentJobId] = useState<string | null>(null)
const pollRef = useRef<ReturnType<typeof setInterval> | null>(null)

const handleGenerate = async () => {
  if (!token) { router.push("/login"); return }
  setStatus("pending")
  setVideoUrl("")
  try {
    const data = await api.createJob(token, form)
    setCurrentJobId(data.job_id)
    setStatus("in_queue")
    startPolling(data.job_id)
  } catch (err: any) {
    setStatus("idle")
    alert(err.message ?? "Failed to create job")
  }
}

const startPolling = (jobId: string) => {
  if (pollRef.current) clearInterval(pollRef.current)
  pollRef.current = setInterval(async () => {
    try {
      const job = await api.getJob(token!, jobId)
      setStatus(job.status as any)
      if (job.status === "completed") {
        setVideoUrl(job.gcs_video_url)
        clearInterval(pollRef.current!)
      }
      if (job.status === "failed") {
        clearInterval(pollRef.current!)
        alert("Video generation failed: " + job.error_message)
      }
    } catch (_) { /* keep polling on transient errors */ }
  }, 3000)
}

useEffect(() => () => { if (pollRef.current) clearInterval(pollRef.current) }, [])
```

**Step 3 — Update status text labels** in the preview panel:
```
"pending"    → "Initializing request..."
"in_queue"   → "In Queue — waiting for AI engine"
"processing" → "Rendering AI Video..."
"completed"  → show video player
"failed"     → show error + "Try Again" button that resets to idle
```

**Step 4 — Wrap body with `<AuthGuard>`** (import from `@/components/AuthGuard`)

---

### Task 3.2 — Wire History Page
**File:** `web/src/app/history/page.tsx` — **[MODIFY]**
**Current state:** Renders from `MOCK_HISTORY` hardcoded array.

**Steps:**
1. Add `"use client"` directive
2. Import `api`, `Job` from `@/lib/api`; `useAppContext`; `useEffect`, `useState`
3. Delete the `MOCK_HISTORY` constant
4. Add state: `const [jobs, setJobs] = useState<Job[]>([])`
5. On mount: `api.listJobs(token!).then(setJobs).catch(console.error)`
6. Map `jobs` to the card UI — field mapping:
   - `job.status.toLowerCase()` → badge class (`completed`, `processing`, `failed`)
   - `job.prompt_input?.template` → template label
   - `new Date(job.created_at).toLocaleDateString()` → date
   - `job.gcs_video_url` → thumbnail/play button (only when `status === "completed"`)
7. Add "No jobs yet" empty state when `jobs.length === 0`
8. Keep the search input + status filter select but implement client-side filtering with `useMemo`
9. Wrap with `<AuthGuard>`

---

### Task 3.3 — Wire Profile Page
**File:** `web/src/app/profile/page.tsx` — **[MODIFY]**
**Current state:** Reads from mocked context user. Will partially work after Task 2.4 but needs these fixes:

1. Add null guard: `if (!user) return null`
2. Replace `user.credits` → `user.credit_balance`
3. Replace `user.name` → `user.email` (no name field in backend yet)
4. Replace `user.company` badge → `user.plan` (shows plan tier: free/starter/pro)
5. Replace hardcoded `24` (videos generated) with: `api.listJobs(token!).then(j => setJobCount(j.length))`
6. Replace mock "Credit Usage History" chart with real data:
   - Call `api.getCreditHistory(token!)` on mount
   - Map last 7 transactions to bar heights (scale by max absolute value)
   - Bars for negative amounts (deductions) = red/accent color
   - Bars for positive amounts (top-ups) = green color
7. Replace the 2 hardcoded activity list items with real `credit_transactions` mapped to the `<li>` template
8. Wire "Sign Out" button: `onClick={() => { logout(); router.push("/login") }}`
9. "Buy Credits" button: show a coming-soon `alert()` for now
10. Wrap with `<AuthGuard>`

---

### Task 3.4 — Wire Settings Page
**File:** `web/src/app/settings/page.tsx` — **[MODIFY]**
**Current state:** Theme/Language tabs work (localStorage). Profile tab "Save Changes" does nothing.

**Steps:**
1. Profile tab — "Save Changes" button: show `alert("Profile update coming soon!")` for now and add a `// TODO` comment
2. Notifications tab — wire toggles to `localStorage` for persistence (key: `notif_video_complete`, `notif_low_credit`)
3. Load toggle initial values from `localStorage` in `useEffect`
4. Wrap entire page with `<AuthGuard>`

---

### Task 3.5 — Update Navbar
**File:** `web/src/app/Navbar.tsx` — **[MODIFY]**
**Current state:** Shows hardcoded `user?.credits || 250`, no logout, no guest state.

**Steps:**
1. Get `user`, `logout`, `token` from `useAppContext()`
2. Credit badge: change `user?.credits || 250` → `user?.credit_balance ?? 0`
3. If `user === null` (guest): render "Login" + "Register" links instead of Credits badge + avatar
4. If `user !== null` (authenticated): keep existing layout, add a logout handler to the avatar/profile link area
5. Avatar area: wrap in a dropdown `<div>` with items: "Profile", "Settings", "Sign Out"
   - "Sign Out" onClick: `logout(); router.push("/login")`

---

## 🔲 PHASE 4 — Backend: Real AI Integration (Veo 3)

> **Priority: HIGH** — The core product feature. Complete Phase 1–3 first.

### Task 4.1 — Implement Real Video Generation
**File:** `backend/internal/services/video_service.go` — **[MODIFY — uncomment & complete]**
**Current state:** All Veo code is commented out, function just prints a log line.

**Steps:**
1. Add config injection to `VideoService`:
```go
type VideoService struct{ cfg *config.Config }
func NewVideoService(cfg *config.Config) *VideoService { return &VideoService{cfg: cfg} }
```
2. Uncomment and complete the `GenerateVideo` function:
```go
client, err := genai.NewClient(ctx, &genai.ClientConfig{
    APIKey:  s.cfg.GoogleAPIKey,
    Backend: genai.BackendGeminiAPI,
    // If s.cfg.GoogleGenaiUseVertexAI == "true": use genai.BackendVertexAI
})
op, err := client.Models.GenerateVideos(ctx, "veo-3.0-fast-generate-001", prompt, nil, nil, nil)
// Poll: for !op.Done { time.Sleep(20*time.Second); op, _ = client.Operations.Get(ctx, op.Name, nil) }
// Download: client.Files.Download(ctx, op.Response.GeneratedVideos[0].Video.Name, nil)
```
3. Return `([]byte, error)` — the raw video bytes for GCS upload

---

### Task 4.2 — Create Google Cloud Storage Service
**File:** `backend/internal/services/storage_service.go` — **[CREATE NEW]**
**What:** Upload `.mp4` bytes to GCS bucket → return public URL.

```go
package services
import (
    "bytes"; "context"; "fmt"; "io"
    "cloud.google.com/go/storage"
    "xanders-gen-video/internal/config"
)

type StorageService struct{ cfg *config.Config }
func NewStorageService(cfg *config.Config) *StorageService { return &StorageService{cfg: cfg} }

func (s *StorageService) UploadVideo(ctx context.Context, jobID string, data []byte) (string, error) {
    client, err := storage.NewClient(ctx)
    if err != nil { return "", err }
    defer client.Close()
    objName := fmt.Sprintf("videos/%s.mp4", jobID)
    wc := client.Bucket(s.cfg.GCSBucket).Object(objName).NewWriter(ctx)
    wc.ContentType = "video/mp4"
    if _, err = io.Copy(wc, bytes.NewReader(data)); err != nil { return "", err }
    if err = wc.Close(); err != nil { return "", err }
    return fmt.Sprintf("https://storage.googleapis.com/%s/%s", s.cfg.GCSBucket, objName), nil
}
```

**Also:** Add `GCSBucket string` to `config.go` + add `GCS_BUCKET=xanders-video-output` to `.env`

---

### Task 4.3 — Wire Queue → Video → Storage → DB Update
**File:** `backend/internal/services/queue_service.go` — **[MODIFY]**
**Current state:** Goroutine just prints. Inject `VideoService`, `StorageService`, `DB`.

```go
type QueueService struct {
    videoService   *VideoService
    storageService *StorageService
    db             *pgxpool.Pool
}

func (s *QueueService) EnqueueVideoJob(ctx context.Context, jobID string) error {
    go func() {
        bgCtx := context.Background()

        // 1. Fetch assembled_prompt from DB
        var prompt string
        s.db.QueryRow(bgCtx, "SELECT assembled_prompt FROM video_generation_jobs WHERE id=$1", jobID).Scan(&prompt)

        // 2. Update status to "processing"
        s.db.Exec(bgCtx, "UPDATE video_generation_jobs SET status='processing', updated_at=NOW() WHERE id=$1", jobID)

        // 3. Call Veo API
        videoBytes, err := s.videoService.GenerateVideo(bgCtx, prompt)
        if err != nil {
            s.db.Exec(bgCtx, "UPDATE video_generation_jobs SET status='failed', error_message=$1, updated_at=NOW() WHERE id=$2", err.Error(), jobID)
            return
        }

        // 4. Upload to GCS
        gcsURL, err := s.storageService.UploadVideo(bgCtx, jobID, videoBytes)
        if err != nil {
            s.db.Exec(bgCtx, "UPDATE video_generation_jobs SET status='failed', error_message=$1, updated_at=NOW() WHERE id=$2", err.Error(), jobID)
            return
        }

        // 5. Mark complete
        s.db.Exec(bgCtx, "UPDATE video_generation_jobs SET status='completed', gcs_video_url=$1, updated_at=NOW() WHERE id=$2", gcsURL, jobID)
    }()
    return nil
}
```

---

### Task 4.4 — Expand Prompt Assembler
**File:** `backend/internal/services/prompt_service.go` — **[MODIFY]**
**Current state:** Returns one basic generic string.

Add complete mapping tables for all 7 templates × 4 lighting × 4 camera modes:
```go
var templateBase = map[string]string{
    "aviation":    "a sleek private jet in flight",
    "automotive":  "a luxury sports car on a winding mountain road",
    "fashion":     "a high-fashion editorial model in designer clothing",
    "accessories": "premium luxury accessories on a reflective surface",
    "pc-gaming":   "a custom high-end gaming PC build with RGB lighting",
    "electronics": "a cutting-edge smartphone floating in mid-air",
    "real-estate": "a stunning luxury penthouse with floor-to-ceiling windows",
}
var lightingModifier = map[string]string{
    "golden-hour": "bathed in warm golden hour sunlight with lens flares",
    "cyberpunk":   "illuminated by neon city lights in the rain",
    "studio":      "under clean soft studio lighting with neutral background",
    "dramatic":    "in dramatic high-contrast chiaroscuro lighting",
}
var cameraModifier = map[string]string{
    "slow-track": "smooth cinematic slow tracking shot, 35mm anamorphic lens, f/1.4",
    "drone":      "aerial drone flyover, wide establishing shot, 24mm lens",
    "static":     "locked-off tripod shot, shallow depth of field, 85mm portrait lens",
    "orbit":      "360-degree circular orbit, medium focal length, continuous motion",
}
// Final assembled: "Cinematic 4K {cameraModifier} of {templateBase}, {lightingModifier}, photorealistic, hyperdetailed"
```

---

## 🔲 PHASE 5 — Backend: Supplementary Endpoints

### Task 5.1 — Add `PATCH /api/v1/me` (Update Profile)
**File:** `backend/internal/handlers/user_handler.go`
Allow users to update their email. Must check for uniqueness constraint.

### Task 5.2 — Improve Credit Pricing Logic
**File:** `backend/internal/handlers/job_handler.go`
Replace the hardcoded `creditsRequired := 10` with a proper pricing function:
```go
func calculateCredits(durationSeconds int, resolution string) int {
    base := durationSeconds * 2  // 2 credits/second
    if resolution == "4k" { base += 20 }
    return base
}
```

### Task 5.3 — Add `GET /api/v1/invoices`
Return invoices list for the current user. Required for future billing page.

---

## 🔲 PHASE 6 — End-to-End Testing

### Task 6.1 — Manual Integration Test Checklist
Run both servers:
```bash
# Terminal 1 — Backend
cd backend && go run ./cmd/api/main.go   # http://localhost:8080

# Terminal 2 — Frontend
cd web && npm run dev                    # http://localhost:3000
```

Test sequence (in order):
- [ ] `curl http://localhost:8080/health` → `{"status":"ok"}`
- [ ] `/register` → submit form → JWT in localStorage → redirect to `/`
- [ ] Navbar shows `credit_balance: 100` (new users start with 100 credits)
- [ ] Complete 4-step wizard → submit → status = `in_queue`
- [ ] Neon DB: `SELECT * FROM video_generation_jobs;` → row with `status='pending'`
- [ ] Neon DB: `SELECT * FROM users;` → `credit_balance` decreased by 10
- [ ] Neon DB: `SELECT * FROM credit_transactions;` → deduction logged
- [ ] `/history` → real job visible with correct status
- [ ] `/profile` → shows real email, real credit balance, real transactions
- [ ] Sign Out → JWT removed → redirect to `/login`
- [ ] Access `/` without JWT → redirect to `/login`

### Task 6.2 — Test AI Video Generation (After Phase 4)
- [ ] Add real `GOOGLE_API_KEY` to `backend/.env`
- [ ] Submit generate job → backend logs show Veo API call
- [ ] Poll: job transitions `pending → processing → completed`
- [ ] Job in DB has `gcs_video_url` set
- [ ] Video player shows real video in frontend

---

## 🔲 PHASE 7 — Production Readiness

### Task 7.1 — Finalize Environment Variables
```env
# backend (Cloud Run)
DATABASE_URL=postgresql://...neon.tech/neondb?sslmode=require
PORT=8080
JWT_SECRET=<256-bit-random-secret>
GOOGLE_API_KEY=<enterprise-api-key>
GCS_BUCKET=xanders-video-output
GOOGLE_CLOUD_PROJECT=<gcp-project-id>
GOOGLE_CLOUD_LOCATION=us-central1
INTERNAL_SECRET=<random-webhook-secret>

# web (Vercel)
NEXT_PUBLIC_API_URL=https://xanders-backend-xxxx.run.app
```

### Task 7.2 — Dockerize Backend
**File:** `backend/Dockerfile` — **[CREATE NEW]**
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/api/main.go

FROM alpine:latest
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["./server"]
```

### Task 7.3 — Deploy Backend to Google Cloud Run
```bash
gcloud run deploy xanders-backend \
  --source ./backend \
  --region asia-southeast1 \
  --allow-unauthenticated \
  --set-env-vars "DATABASE_URL=...,JWT_SECRET=...,GOOGLE_API_KEY=...,GCS_BUCKET=..."
```

### Task 7.4 — Deploy Frontend to Vercel
```bash
cd web
vercel --prod
# Add NEXT_PUBLIC_API_URL=https://xanders-backend-xxxx.run.app in Vercel dashboard
```

### Task 7.5 — Security Hardening
- [ ] Add rate limiting: `github.com/gofiber/fiber/v2/middleware/limiter` — max 5 req/min on `/auth/*`
- [ ] Add request size limit: `fiber.Config{ BodyLimit: 1 * 1024 * 1024 }` (1MB)
- [ ] Rotate JWT secret from placeholder value before production
- [ ] Add `context.WithTimeout(ctx, 5*time.Second)` on all DB queries
- [ ] Update CORS `AllowOrigins` from `localhost:3000` to the real Vercel domain in production

---

## 📁 Complete File Reference

```
project/
├── summary.md                          ✅ Updated PRD
├── TASKS.md                            ✅ This file
│
├── web/
│   ├── .env.local                      🔲 Task 2.1 (NEW)
│   └── src/
│       ├── components/
│       │   └── AuthGuard.tsx           🔲 Task 2.7 (NEW)
│       ├── lib/
│       │   ├── api.ts                  🔲 Task 2.2 (NEW)
│       │   └── auth.ts                 🔲 Task 2.3 (NEW)
│       └── app/
│           ├── context/AppContext.tsx  🔲 Task 2.4 (real auth)
│           ├── Navbar.tsx              🔲 Task 3.5
│           ├── page.tsx               🔲 Task 3.1 (real API + polling)
│           ├── login/page.tsx          🔲 Task 2.5
│           ├── register/page.tsx       🔲 Task 2.6
│           ├── history/page.tsx        🔲 Task 3.2
│           ├── profile/page.tsx        🔲 Task 3.3
│           └── settings/page.tsx       🔲 Task 3.4
│
└── backend/
    ├── .env                            ✅ Live Neon URL
    ├── Dockerfile                      🔲 Task 7.2 (NEW)
    ├── cmd/api/main.go                 🔲 Task 1.1 (CORS fix)
    └── internal/
        ├── config/config.go            ✅ Done — add GCSBucket, InternalSecret
        ├── database/db.go              ✅ Done
        ├── models/models.go            ✅ Done
        ├── middleware/auth.go          ✅ Done
        ├── handlers/
        │   ├── auth_handler.go         ✅ Done
        │   ├── user_handler.go         🔲 Task 1.3, 5.1
        │   ├── job_handler.go          🔲 Task 1.2, 5.2
        │   └── webhook_handler.go      🔲 Task 1.4 (NEW)
        ├── services/
        │   ├── auth_service.go         ✅ Done
        │   ├── credit_service.go       ✅ Done — minor pricing improvement in 5.2
        │   ├── prompt_service.go       🔲 Task 4.4 (expand templates)
        │   ├── queue_service.go        🔲 Task 4.3 (wire to video+storage)
        │   ├── video_service.go        🔲 Task 4.1 (uncomment Veo API)
        │   └── storage_service.go      🔲 Task 4.2 (NEW)
        ├── router/router.go            🔲 Task 1.2 (add ListJobs, webhook)
        └── migrations/
            └── 001_initial_schema.sql  ✅ Applied to Neon
```

---

## ⚡ Quick Start for Agent

```bash
# Start backend
cd e:/Work/Project/xanders/project/backend
go run ./cmd/api/main.go
# → Listening on :8080

# Start frontend (new terminal)
cd e:/Work/Project/xanders/project/web
npm run dev
# → http://localhost:3000

# Health check
curl http://localhost:8080/health
```

**Execution order**: Phase 1 → Phase 2 → Phase 3 → Phase 4 → Phase 5 → Phase 6 → Phase 7.
Each phase is a prerequisite for the next. Do not skip ahead.
