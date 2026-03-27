# OCR Image Resize — Plan

---

## 1. Context

- Camera điện thoại hiện tại: 12-50MP, ảnh gốc 5-15MB
- API contract: max 10MB
- Google Vision: max 20MB raw, 10MB base64 JSON
- Baidu OCR: max 4MB (base64), 10MB (URL)
- Latency budget: p50 < 1.5s, p99 < 3s
- Use case: scan vở ghi chép / sách giáo khoa → nhận diện Hán tự
- OCR accuracy cho text recognition: 1-2MP là đủ, >4MP không tăng accuracy

---

## 2. Flow chi tiết — KHÔNG resize

```
User chụp ảnh (12MP, ~8MB JPEG)
    │
    ▼
[Mobile App]
    │  Upload multipart/form-data (8MB, ảnh gốc)
    │  ⏱ Upload time: 3-8s (3G/4G), 1-2s (WiFi)
    │
    ▼
[Nginx / Load Balancer]
    │  Body limit check (max 10MB) ── ✅ Pass (8MB)
    │
    ▼
[Go Server — OCR Handler]
    │  Validate: file size ≤ 10MB, format, type, language ── ✅
    │  Forward 8MB ảnh gốc (12MP) cho OCR engine
    │
    ▼
[OCR Engine]
    │  Xử lý ảnh 12MP — chậm hơn, tốn resource hơn, accuracy KHÔNG tăng
    │
    ▼
[Response] ✅ Hoạt động, nhưng lãng phí
```

### Vấn đề khi không resize

| Vấn đề | Chi tiết |
|---|---|
| **Reject ảnh quá lớn** | Camera 50MP+ có thể > 10MB limit |
| **Tốn bandwidth** | Upload 8MB trên 4G: 3-8s |
| **Tốn server RAM** | Decode 12MP JPEG = ~48MB RAM per request trên server |
| **Lãng phí** | OCR accuracy cho text chỉ cần 1-2MP, gửi 12MP không tăng kết quả |

**Kết luận: Nên có resize trên server để tối ưu resource trước khi gửi cho OCR engine.**

---

## 3. Flow chi tiết — Có resize

### Phương án A: Client-side resize (Mobile)

> ⚠️ Resize trên thiết bị người dùng — máy yếu/cũ có thể bị lag, tốn pin, ảnh hưởng UX.

```
User chụp ảnh (12MP, ~8MB JPEG)
    │
    ▼
[Mobile App — Image Processing]
    │  1. Resize: 12MP → 2MP (1600×1200)
    │  2. Compress: JPEG quality 85%
    │  3. Kết quả: ~300-600KB
    │  ⏱ Processing: 50-200ms (flagship), 500-2000ms (máy cũ/yếu)
    │
    ▼
[Mobile App — Upload]
    │  Upload multipart/form-data (~500KB)
    │  Fields: image, type, language
    │  Headers: X-Idempotency-Key, Authorization
    │  ⏱ Upload: 0.3-1s (4G), <0.5s (WiFi)
    │
    ▼
[Nginx / Load Balancer]
    │  Body limit check (max 10MB) ── ✅ Pass (~500KB)
    │
    ▼
[Go Server — Middleware]
    │  Auth → RateLimit → RequestID → OTEL
    │
    ▼
[OCR Handler]
    │  1. Parse multipart form
    │  2. Validate: size ≤ 10MB, format (JPEG/PNG), type, language
    │  3. Read image bytes
    │
    ▼
[OCR Use Case — OCRCommand]
    │  1. Check idempotency (Redis) ── cache hit → return cached result
    │  2. Check per-user daily limit (Redis) ── exceeded → 429
    │  3. Resolve engine (type + language → engine chain)
    │  4. Call OCR engine
    │
    ├──────────────────────────────────────────────┐
    ▼                                              ▼
[Google Vision Adapter]                    [Baidu OCR Adapter]
    │  Send image bytes (gRPC)                │  Base64 encode → REST POST
    │  Feature: DOCUMENT_TEXT_DETECTION        │  recognize_granularity=small
    │  ⏱ ~500-1500ms                          │  ⏱ ~500-1500ms
    │                                          │
    ▼                                          ▼
[OCR Raw Result]◄──────────────────────────────┘
    │
    ▼
[Post-processing Pipeline]
    │  1. Normalize (Unicode NFC)
    │  2. Filter CJK (unicode.Han)
    │  3. Deduplicate (giữ confidence cao nhất)
    │  4. Word segmentation + Enrich (CC-CEDICT → pinyin, meaning, HSK)
    │  5. Classify confidence (confirmed / suggest / low_confidence)
    │  6. Generate candidates (confusable_map cho medium + low)
    │  ⏱ < 25ms
    │
    ▼
[Vocabulary Classify] (Phase 6 — pending)
    │  FindByHanziList() → split new_items / existing_items
    │  ⏱ ~10ms (DB query)
    │
    ▼
[Response]
    │  { new_items[], existing_items[], low_confidence_items[], metadata }
    │
    ▼
[Mobile App — Preview Screen]
    │  Hiển thị detected characters
    │  User: Edit / Delete / Add missing / Confirm all
    │
    ▼
[Mobile App — Confirm]
    │  POST /api/vocabularies (bulk create)
    │  Assign to folder/topic
    │
    ▼
[Vocabulary Created] ✅

Tổng latency (client-side resize):
  Resize:   50-200ms
  Upload:   300-1000ms
  Server:   500-1500ms (OCR) + 25ms (post-process) + 10ms (DB)
  ────────────────────
  Total:    ~900-2700ms ✅ Trong budget p50 < 1.5s, p99 < 3s
```

---

### Phương án B: Server-side resize

```
User chụp ảnh (12MP, ~8MB JPEG)
    │
    ▼
[Mobile App — Upload]
    │  Upload multipart/form-data (8MB, ảnh gốc)
    │  ⏱ Upload: 3-8s (4G), 1-2s (WiFi)
    │
    ▼
[Nginx / Load Balancer]
    │  Body limit: tăng lên 15-20MB để chấp nhận ảnh gốc
    │
    ▼
[Go Server — Middleware]
    │  Auth → RateLimit → RequestID → OTEL
    │
    ▼
[OCR Handler]
    │  1. Parse multipart form
    │  2. Validate: size ≤ 20MB, format
    │  3. Read image bytes
    │
    ▼
[Image Resize Service] ◄── NEW COMPONENT
    │  1. Decode image (JPEG/PNG)
    │  2. Resize: → 2MP (1600×1200)
    │  3. Re-encode JPEG quality 85%
    │  4. Kết quả: ~300-600KB
    │  ⏱ 100-500ms (Go image processing, CPU-bound)
    │  ⚠️ Memory spike: decode 12MP = ~48MB RAM per request
    │
    ▼
[OCR Use Case — OCRCommand]
    │  (... giống Phương án A từ đây ...)
    │
    ▼
[Response → Preview → Confirm → Vocabulary Created] ✅

Tổng latency (server-side resize):
  Upload:   3000-8000ms (4G) ❌ ĐÃ VƯỢT BUDGET
  Resize:   100-500ms
  Server:   500-1500ms (OCR) + 25ms + 10ms
  ────────────────────
  Total:    ~3600-10000ms ❌ Vượt xa p99 < 3s
```

---

### Phương án C: Hybrid (Client compress + Server validate/re-resize)

> ⚠️ Cùng vấn đề với A: resize trên thiết bị người dùng, máy yếu/cũ bị lag. Thêm complexity vì phải maintain logic resize ở cả 2 nơi.

```
User chụp ảnh (12MP, ~8MB JPEG)
    │
    ▼
[Mobile App — Image Processing]
    │  1. Resize: → max 2MP (1600×1200)
    │  2. Compress: JPEG quality 85%
    │  3. Kết quả: ~300-600KB
    │  ⏱ 50-200ms (flagship), 500-2000ms (máy cũ/yếu)
    │
    ▼
[Mobile App — Upload]
    │  Upload (~500KB)
    │  ⏱ 0.3-1s
    │
    ▼
[Go Server — OCR Handler]
    │  1. Validate size ≤ 10MB ── ✅
    │  2. Decode image header → check dimensions
    │  3. Nếu dimensions > 2MP (client cũ/bug):
    │     → Server resize xuống 2MP ── fallback safety net
    │  4. Nếu dimensions ≤ 2MP:
    │     → Forward trực tiếp (không xử lý thêm)
    │  ⏱ Header check: <1ms, resize (nếu cần): 100-500ms
    │
    ▼
[OCR Use Case]
    │  (... giống Phương án A ...)
    │
    ▼
[Response → Preview → Confirm → Vocabulary Created] ✅

Tổng latency (hybrid, happy path — client đã resize đúng):
  Resize:   50-200ms (client)
  Upload:   300-1000ms
  Server:   <1ms (validate) + 500-1500ms (OCR) + 25ms + 10ms
  ────────────────────
  Total:    ~900-2700ms ✅

Tổng latency (hybrid, fallback — client gửi ảnh lớn):
  Upload:   300-1000ms (vẫn ≤ 10MB)
  Resize:   100-500ms (server)
  Server:   500-1500ms (OCR) + 25ms + 10ms
  ────────────────────
  Total:    ~900-3000ms ⚠️ Sát budget nhưng chấp nhận được
```

---

### Phương án D: Dedicated OCR Service (tách service riêng)

Tách toàn bộ xử lý ảnh (resize + gọi OCR engine) ra một service riêng biệt, Go API chỉ proxy request.

```
User chụp ảnh (12MP, ~8MB JPEG)
    │
    ▼
[Mobile App — Upload]
    │  Upload multipart/form-data (≤ 10MB)
    │
    ▼
[Go API Server — Middleware]
    │  Auth → RateLimit → RequestID → OTEL
    │
    ▼
[OCR Handler]
    │  1. Validate: size ≤ 10MB, format, type, language
    │  2. Check idempotency (Redis)
    │  3. Check per-user daily limit (Redis)
    │  4. Forward image bytes + metadata sang OCR Service
    │     (gRPC hoặc HTTP multipart)
    │
    ▼ ─── network hop ───
    │
[OCR Service] ◄── SEPARATE DEPLOYMENT (container riêng, scale riêng)
    │
    │  1. Nhận image bytes
    │  2. Resize logic:
    │     ├── size ≤ 2MB → skip resize
    │     └── size > 2MB → decode → resize 1600×1200 → JPEG 85%
    │  3. Resolve engine (type + language → engine chain)
    │  4. Gọi OCR engine (Google Vision / Baidu / PaddleOCR / Tesseract)
    │  5. Post-processing pipeline (normalize, filter CJK, dedup, enrich, classify confidence, candidates)
    │  6. Return OCR result
    │
    ▼ ─── network hop ───
    │
[Go API Server — OCR Use Case]
    │  1. Nhận OCR result từ OCR Service
    │  2. Vocabulary classify: FindByHanziList() → split new/existing
    │  3. Build response
    │
    ▼
[Response → Preview → Confirm → Vocabulary Created] ✅
```

**Ưu điểm:**

| Ưu điểm | Chi tiết |
|---|---|
| **Isolate resource** | Image decode/resize rất tốn CPU + RAM (~48MB/request cho 12MP). Tách ra → Go API server giữ lightweight, không bị ảnh hưởng bởi OCR spike |
| **Scale độc lập** | OCR Service scale theo OCR traffic (spike 18h-22h, trước kỳ thi 10x). Go API scale theo API traffic chung. Không phải scale cả monolith vì OCR |
| **Fault isolation** | OCR Service crash/OOM → Go API vẫn chạy, các endpoint khác (auth, vocabulary CRUD) không bị ảnh hưởng |
| **Tech stack linh hoạt** | OCR Service có thể viết bằng Go, Python, hoặc mix. Python sidecar (PaddleOCR/Tesseract) có thể merge vào luôn thay vì chạy riêng |
| **Resource tuning riêng** | Set CPU/memory limit riêng cho OCR pods (high CPU, high memory) vs API pods (low CPU, low memory) |
| **Deploy độc lập** | Update OCR logic (thêm engine, đổi threshold, update model) không cần redeploy Go API |

**Nhược điểm:**

| Nhược điểm | Chi tiết |
|---|---|
| **Thêm infrastructure** | 1 service nữa để deploy, monitor, maintain. Cần thêm Dockerfile, K8s manifest, health check, logging |
| **Network latency** | Thêm 1 hop giữa Go API ↔ OCR Service. Nếu gửi image bytes qua network: +5-50ms (cùng cluster). Ảnh 8MB qua gRPC trong cùng cluster ~20-50ms |
| **Operational complexity** | 2 services phải coordinate: versioning API contract, deploy order, distributed tracing, error handling cross-service |
| **MVP overkill** | Traffic MVP: 100-600 req/ngày. Go API thừa sức xử lý resize + OCR trong cùng process |
| **Debugging khó hơn** | Lỗi OCR phải trace qua 2 services. Cần distributed tracing (OTEL đã có, nhưng thêm effort setup cho service mới) |
| **Image transfer overhead** | Gửi ảnh 8MB qua network giữa 2 services trong cùng cluster → tốn internal bandwidth. Giảm bớt nếu dùng shared storage/volume |

**Khi nào nên tách:**

- Traffic OCR > 10K req/ngày và Go API bắt đầu bị ảnh hưởng bởi memory spike từ image processing
- Cần scale OCR pods riêng (ví dụ: 10 OCR pods vs 3 API pods)
- Team lớn hơn, muốn deploy OCR logic độc lập

**MVP: KHÔNG nên tách.** Giữ resize + OCR trong Go API. Tách khi có dấu hiệu cần thiết (growth phase, >10K req/ngày).

---

## 4. So sánh

| Tiêu chí | A: Client-side | B: Server-side | C: Hybrid | D: Dedicated OCR Service |
|---|---|---|---|---|
| **Thiết bị người dùng** | ❌ Máy yếu/cũ lag, tốn pin | ✅ Không ảnh hưởng | ❌ Cùng vấn đề A | ✅ Không ảnh hưởng |
| **Bandwidth** | ✅ ~500KB upload | ⚠️ ~8MB upload | ✅ ~500KB upload | ⚠️ ~8MB upload + internal transfer |
| **API server CPU/RAM** | ✅ Không tốn | ⚠️ ~48MB RAM/req khi resize | ⚠️ Chỉ khi fallback | ✅ Không tốn (OCR Service chịu) |
| **K8s scaling** | ✅ Không cần thêm resource | ⚠️ Cần tăng memory limit per pod | ✅ Không cần | ✅ Scale OCR riêng, API giữ nhẹ |
| **Fault isolation** | ✅ Không ảnh hưởng server | ⚠️ OOM có thể crash API | ⚠️ Fallback resize có thể spike | ✅ OCR crash không ảnh hưởng API |
| **Reliability** | ⚠️ Phụ thuộc client | ✅ Server kiểm soát 100% | ✅ Client + server safety net | ✅ Server kiểm soát 100% |
| **Complexity** | ✅ Đơn giản | ✅ Đơn giản, 1 nơi | ❌ Maintain logic ở cả client + server | ❌ Thêm 1 service, infra, deploy |
| **Multi-platform** | ⚠️ Implement iOS + Android | ✅ 1 lần trên server | ⚠️ Client + server | ✅ 1 lần trong OCR Service |
| **Old client** | ❌ Reject nếu chưa resize | ✅ Mọi client OK | ✅ Server fallback | ✅ Mọi client OK |
| **MVP phù hợp** | ⚠️ | ✅ **Recommend** | ❌ Over-engineer | ❌ Overkill |

---

## 5. Phương án phù hợp theo giai đoạn scale

### Tổng quan timeline

```
MVP (tháng 1-3)          Growth (tháng 4-12)         Scale (năm 2+)
50-200 MAU               1K-10K MAU                  10K-50K MAU
100-600 req/ngày          3K-50K req/ngày              50K-500K req/ngày
~0.03 QPS peak           ~2-10 QPS peak              ~10-100 QPS peak
    │                        │                           │
    ▼                        ▼                           ▼
Phương án B              Phương án B                  Phương án D
(Server-side resize      (Tăng pod resources,         (Tách OCR Service)
 trong Go API)            HPA autoscale)              + Public URL mode (optional)
```

### Giai đoạn 1: MVP (tháng 1-3) → Phương án B

| | Chi tiết |
|---|---|
| **Traffic** | 100-600 req/ngày, QPS gần như = 0 |
| **Phương án** | **B: Server-side resize** — resize + OCR đều chạy trong Go API |
| **Lý do** | Đơn giản nhất. 1 service, 1 deploy. Traffic quá thấp, không cần optimize. Memory spike 48MB/request × 1 concurrent = không đáng lo |
| **Infra** | 1-2 API pods, mỗi pod 256-512MB RAM là đủ |
| **Dấu hiệu chuyển giai đoạn tiếp** | API response time tăng, memory usage > 70% pod limit khi có concurrent OCR requests |

### Giai đoạn 2: Growth (tháng 4-12) → Phương án B (scale vertically + horizontally)

| | Chi tiết |
|---|---|
| **Traffic** | 3K-50K req/ngày, peak 2-10 QPS, spike giờ homework 3-5x |
| **Phương án** | **Giữ B** — vẫn resize trong Go API, tăng resource per pod + thêm pods |
| **Scale cách nào** | Tăng memory limit 512MB-1GB per pod. HPA (Horizontal Pod Autoscaler) scale 2-4 pods theo CPU/memory. Đủ chịu concurrent OCR |
| **Infra** | 2-4 API pods, 512MB-1GB per pod, HPA auto-scale khi spike |
| **Dấu hiệu chuyển giai đoạn tiếp** | OCR memory spike ảnh hưởng non-OCR endpoints (auth, vocabulary CRUD chậm theo). Pod restart/OOM khi concurrent OCR > 10-20. Scale thêm pods không hiệu quả vì mỗi pod đều phải chịu resize overhead |

### Giai đoạn 3: Scale (năm 2+) → Phương án D

| | Chi tiết |
|---|---|
| **Traffic** | 50K-500K req/ngày, peak 10-100 QPS, spike kỳ thi 10x |
| **Phương án** | **D: Tách OCR Service riêng** |
| **Lý do** | OCR chiếm phần lớn CPU/RAM nhưng chỉ là 1 endpoint. Tách ra để: (1) API server nhẹ, phục vụ auth + vocabulary CRUD ổn định; (2) OCR pods scale riêng theo demand; (3) OCR crash/OOM không kéo theo cả hệ thống |
| **Kết hợp** | D + Public URL mode (optional, giảm internal bandwidth khi cần) |
| **Infra** | API: 3-5 pods (256MB). OCR Service: 5-15 pods (1-2GB, high CPU), HPA scale theo QPS + memory |
| **Thêm** | Merge Python sidecar (PaddleOCR/Tesseract) vào OCR Service luôn → 1 service xử lý tất cả engine. Nếu >100K req/ngày: cân nhắc chuyển sang PaddleOCR self-hosted thay vì cloud APIs (rẻ hơn 4-5x) |

### Migration path

```
Giai đoạn 1 (MVP)                     Giai đoạn 2 (Growth)                   Giai đoạn 3 (Scale)
─────────────────                      ──────────────────────                  ─────────────────────

┌─────────────────┐                    ┌─────────────────┐                    ┌──────────────┐
│   Go API Pod    │                    │   Go API Pod    │                    │  Go API Pod  │
│                 │                    │   (more pods,   │                    │  (nhẹ)       │
│  ┌───────────┐  │                    │    more RAM)    │                    │  Auth        │
│  │  Resize   │  │                    │  ┌───────────┐  │                    │  Vocab CRUD  │
│  │  + OCR    │  │                    │  │  Resize   │  │                    │  Proxy OCR   │
│  │  + Post-  │  │                    │  │  + OCR    │  │                    └──────┬───────┘
│  │  process  │  │                    │  │  + Post-  │  │                           │ gRPC
│  └───────────┘  │                    │  │  process  │  │                           ▼
│                 │                    │  └───────────┘  │                    ┌──────────────┐
│  Auth           │                    │                 │                    │ OCR Service  │
│  Vocab CRUD     │                    │  Auth           │                    │ (scale riêng)│
└─────────────────┘                    │  Vocab CRUD     │                    │  Resize      │
                                       └─────────────────┘                    │  OCR engines │
                                                                              │  Post-process│
                                                                              │  PaddleOCR   │
                                                                              └──────────────┘

1-2 pods                               2-4 pods (HPA)                        API: 3-5 pods
256-512MB/pod                          512MB-1GB/pod                         OCR: 5-15 pods
```

### Quyết định chuyển giai đoạn dựa trên metrics

| Metric | Threshold → hành động |
|---|---|
| API pod memory usage | > 70% do OCR spike → tăng pod limit hoặc thêm pods |
| Pod OOM restart | Bất kỳ OOM nào do OCR → ưu tiên cao chuyển D |
| Non-OCR endpoint latency | p50 tăng > 20% khi có OCR traffic → chuyển D |
| OCR req/ngày | > 10K sustained → bắt đầu plan D |
| OCR req/ngày | > 50K → D + cân nhắc PaddleOCR self-hosted thay cloud APIs |
| Cloud OCR cost | > 70% monthly budget → cân nhắc PaddleOCR self-hosted |

---

## 6. Resize spec

> Chi tiết research: [`research_image_resize.md`](image_resize_research.md)

| Param | Value | Lý do |
|---|---|---|
| Skip resize | file size ≤ 2MB AND longest side ≤ 2048px | Ảnh đã đủ nhỏ → forward thẳng, không tốn CPU |
| Target | **Longest side = 2048px, giữ aspect ratio** | Trên ngưỡng Google (1024×768), dưới limit Baidu (4096). Không méo ảnh |
| Interpolation | **CatmullRom** | Giữ nét chữ sắc, nhanh hơn Lanczos, chất lượng gần tương đương |
| Format | JPEG | 10x nhỏ hơn PNG cho ảnh chụp |
| Quality | 85% | Sweet spot — dưới 75% artifacts ảnh hưởng Hán tự handwritten |
| Max file size | 10MB (hard limit) | API contract |
| Max pixel reject | > 75MP → reject 400 | Prevent OOM khi decode |
| Go library (MVP) | `golang.org/x/image/draw` | Official, zero deps, CatmullRom built-in, ~100-300ms, ~60MB peak RAM |
| Go library (Scale) | `h2non/bimg` (libvips) | 4-8x nhanh, 3-5x ít RAM, dùng khi tách OCR Service |

### Resize logic

```
Nhận image bytes
    │
    ├── file size > 10MB → reject 400
    │
    ├── DecodeConfig() → check dimensions
    │   └── > 75MP → reject 400 (prevent OOM)
    │
    ├── file size ≤ 2MB AND longest side ≤ 2048px
    │   → skip resize, forward nguyên bản
    │
    └── cần resize:
        1. Decode JPEG/PNG
        2. scale = 2048.0 / longestSide
        3. Resize bằng CatmullRom (giữ aspect ratio)
        4. Encode JPEG quality 85
        5. Forward cho OCR engine

Ví dụ:
  4000×3000 (landscape) → 2048×1536
  3000×4000 (portrait)  → 1536×2048
  2000×1500             → skip (đã ≤ 2048)
```
