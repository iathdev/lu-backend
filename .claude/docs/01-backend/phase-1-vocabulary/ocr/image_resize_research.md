# Research: Image Resize cho OCR Pipeline

---

## 1. Giới hạn từ OCR engines

| Constraint | Google Cloud Vision | Baidu OCR API |
|---|---|---|
| **Recommend resolution** | ≥ 1024×768 | ≤ 2K (1080p-2K) |
| **Max file size (base64)** | ~7.3MB raw (10MB JSON limit) | **3MB raw (4MB encoded)** ← tight nhất |
| **Max file size (URL)** | 20MB | 10MB |
| **Max dimensions** | 75MP (tự resize) | **4096px longest edge** |
| **Min dimensions** | 640×480 (recommend) | 15px shortest edge |

> "Image sizes larger than 640×480 pixels may not gain much in accuracy, while greatly diminishing throughput." — Google Cloud Vision docs

**Tóm lại:** Target resize nằm giữa **1024-2048px longest side**, output JPEG **< 3MB** (để vừa Baidu base64 limit).

---

## 2. Best practices

### Resolution

| Rule | Lý do |
|---|---|
| **Giữ aspect ratio** — không ép fixed dimensions | Méo ảnh → méo chữ → hại OCR accuracy |
| **Không upscale** — ảnh nhỏ hơn target thì pass through | Upscale tạo artifacts, chữ bị blur |
| **Target 2048px longest side** | Trên recommend Google (1024), dưới limit Baidu (4096). Hán tự ~30-80px ở mức này — dư sức |
| **Hán tự phải ≥ 15-20px cao** | Dưới 15px OCR engine không nhận diện được nét |

### Compression

| Rule | Lý do |
|---|---|
| **JPEG quality ≥ 80** (sweet spot: 85) | Dưới 75: "mosquito noise" quanh nét chữ → merge/tạo phantom strokes. Hại nhất cho Hán tự giống nhau (大/太, 已/己/巳) |
| **Dùng JPEG, không PNG** | PNG lossless nhưng 10x lớn hơn (3-6MB vs 300-600KB). Accuracy gain không đáng kể cho cloud OCR |

### Post-processing

| Rule | Lý do |
|---|---|
| **KHÔNG binarize, deskew, denoise** | Cloud OCR engines có preprocessing riêng. Làm thêm có thể mất thông tin hoặc tạo artifacts |
| **Sharpen: optional, nhẹ** (sigma 0.5 max) | Phục hồi nét mất khi downscale. Over-sharpen tạo halo artifacts → hại OCR |

### Safety

| Rule | Lý do |
|---|---|
| **`DecodeConfig()` trước `Decode()`** | Check dimensions trước khi decode → reject ảnh quá lớn sớm, tránh OOM (decode 75MP = 300MB RAM) |
| **Skip resize nếu ảnh đã đủ nhỏ** | File ≤ 2MB AND longest side ≤ 2048 → forward thẳng, tránh tốn CPU |

---

## 3. Resize strategy

### Bối cảnh

Khi resize ảnh, cần chọn cách xác định kích thước đích. Có 3 cách phổ biến:

- **Max longest side**: Scale sao cho cạnh dài nhất = target, giữ tỉ lệ. VD: 4000×3000 → 2048×1536
- **Fixed dimensions**: Ép về kích thước cố định. VD: 4000×3000 → 1600×1200
- **Max area**: Scale sao cho tổng pixel ≤ target. VD: max 2MP → `scale = sqrt(2MP / (w×h))`

### So sánh

| Tiêu chí | Max longest side | Fixed dimensions | Max area |
|---|---|---|---|
| **OCR accuracy** | ✅ Tốt nhất — giữ tỉ lệ | ❌ Méo ảnh → méo chữ | ✅ Tốt — giữ tỉ lệ |
| **Portrait + landscape** | ✅ Xử lý đúng cả 2 | ❌ Portrait bị ép landscape | ✅ Xử lý đúng |
| **Đơn giản** | ✅ `scale = 2048 / longestSide` | ✅ Đơn giản | ⚠️ `scale = sqrt(maxArea/(w*h))` |
| **Industry standard** | ✅ Google Vision, hầu hết OCR pipelines | ❌ Không ai dùng cho OCR | ⚠️ Ít phổ biến |
| **Output size** | ⚠️ Thay đổi theo aspect ratio | ✅ Cố định | ✅ Gần cố định |

**Recommend: Max longest side = 2048px.**

---

## 4. Interpolation method

### Interpolation là gì?

Khi resize 4000×3000 → 2048×1536, ảnh mới có ít pixel hơn. Mỗi pixel mới nằm "giữa" nhiều pixel gốc. **Interpolation** = thuật toán tính màu cho pixel mới từ các pixel gốc xung quanh.

```
Ảnh gốc (4000px)                    Ảnh mới (2048px)
┌─┬─┬─┬─┬─┬─┬─┬─┐                  ┌──┬──┬──┬──┐
│A│B│C│D│E│F│G│H│   ──resize──►    │ ?│ ?│ ?│ ?│
└─┴─┴─┴─┴─┴─┴─┴─┘                  └──┴──┴──┴──┘

Pixel "?" = tính từ các pixel gốc xung quanh.
Càng nhiều pixels tham chiếu → kết quả càng sắc nét → nhưng càng chậm.
```

**Quan trọng cho OCR** vì chữ Hán = nét mảnh high-contrast. Interpolation kém → nét bị blur/blocky → OCR nhận sai.

### So sánh

| Tiêu chí | CatmullRom (Bicubic) | Lanczos | BiLinear | NearestNeighbor |
|---|---|---|---|---|
| **Cách hoạt động** | Đường cong bậc 3, 4×4 = 16 pixels | Hàm sinc, 6×6 = 36+ pixels | Trung bình 2×2 = 4 pixels | Lấy 1 pixel gần nhất |
| **OCR quality** | ✅ Rất tốt — giữ nét sắc | ✅ Tốt nhất — sắc nhất | ⚠️ Hơi mờ edges | ❌ Blocky, phá nét chữ |
| **Tốc độ** | ✅ Nhanh | ❌ Chậm nhất | ✅ Nhanh | ✅ Nhanh nhất |
| **Go stdlib** | ✅ `x/image/draw` | ❌ Cần lib ngoài | ✅ `x/image/draw` | ✅ `x/image/draw` |
| **Ưu điểm chính** | Best tradeoff quality/speed | Chuẩn vàng downscaling | Đơn giản, nhanh | Nhanh nhất |
| **Nhược điểm chính** | Hơi kém Lanczos (nhưng cloud OCR có preprocessing riêng → difference gần mất) | Chậm, cần lib ngoài. Marginal gain vs CatmullRom | Nét mảnh Hán tự handwritten có thể blur | **Không dùng cho OCR** |

**Recommend: CatmullRom.** Tốt nhất trong stdlib, quality gần Lanczos. Cloud OCR engines có preprocessing riêng → difference giữa CatmullRom và Lanczos gần như biến mất.

---

## 5. Output format

### Bối cảnh

Sau resize, ảnh phải encode lại để gửi cho OCR engine. Chọn format (JPEG/PNG) và quality level.

JPEG nén lossy — bỏ bớt detail, file nhỏ. Ở quality thấp tạo **"mosquito noise"** (đốm mờ quanh nét chữ) → merge thin strokes hoặc tạo phantom strokes. Đặc biệt hại cho Hán tự chỉ khác nhau 1 nét (大/太, 已/己/巳).

### So sánh

| Tiêu chí | JPEG 85 | JPEG 95 | PNG (lossless) | JPEG 75 |
|---|---|---|---|---|
| **File size** (2048px) | ✅ 300-600KB | ⚠️ 800KB-1.5MB | ❌ 3-6MB | ✅ 150-300KB |
| **OCR accuracy (printed)** | ✅ Không ảnh hưởng | ✅ Không ảnh hưởng | ✅ Tốt nhất | ✅ Không ảnh hưởng |
| **OCR accuracy (handwritten)** | ✅ Không ảnh hưởng | ✅ Không ảnh hưởng | ✅ Tốt nhất | ⚠️ Artifacts quanh nét mảnh |
| **Trong Baidu 4MB limit** | ✅ Thoải mái | ✅ Vừa | ❌ Có thể vượt | ✅ Thoải mái |
| **Ưu điểm chính** | **Best balance** size vs quality | Gần lossless | Không mất detail | File nhỏ nhất |
| **Nhược điểm chính** | Lossy (invisible ở 85) | 2-3x lớn hơn, gain = 0 vs 85 | 10x lớn hơn JPEG | Rủi ro mosquito noise |

**Recommend: JPEG quality 85.** Floor: 80.

---

## 6. Go library

### Bối cảnh

Resize ảnh trong Go gồm 3 bước: **decode** (đọc JPEG/PNG → pixel data trong RAM) → **resize** (interpolation) → **encode** (nén lại JPEG/PNG). Go stdlib chỉ có decode/encode, cần thêm library cho resize.

Hai loại:
- **Pure Go**: Dễ deploy (`go build` là xong), chậm hơn, tốn RAM hơn
- **CGO** (libvips): 4-8x nhanh, 3-5x ít RAM, nhưng cần thư viện C → Docker phức tạp hơn

### So sánh

| Tiêu chí | `x/image/draw` (stdlib) | `disintegration/imaging` | `h2non/bimg` (libvips) | `go-scaled-jpeg` + `x/image` |
|---|---|---|---|---|
| **Loại** | Pure Go | Pure Go | CGO (C binding) | Pure Go (hybrid) |
| **Tốc độ** (12MP resize) | 100-300ms | 100-300ms | **15-40ms** | 60-150ms |
| **Peak RAM** (per request) | 60-65MB | 60-65MB | **15-20MB** | ~25MB |
| **Dependencies** | ✅ Zero (stdlib ext) | ⚠️ Ít maintain (03/2023) | ❌ CGO + libvips | ⚠️ Niche lib |
| **Deploy** | ✅ `go build` | ✅ `go build` | ❌ Docker multi-stage, +30-50MB | ✅ `go build` |
| **CatmullRom / Sharpen** | ✅ / ❌ | ✅ / ✅ | ✅ / ✅ | ✅ / ❌ |
| **Ưu điểm chính** | Official, zero deps, đơn giản nhất | Feature-rich, có Sharpen | **Nhanh nhất, ít RAM nhất**, battle-tested | Giảm RAM 60→25MB, vẫn pure Go |
| **Nhược điểm chính** | Chậm nhất, tốn RAM nhất (48MB decode 12MP) | Ít maintain, rủi ro abandon | Deploy phức tạp, debug khó (CGO) | Chỉ baseline JPEG, niche |

### Recommend theo giai đoạn

| Giai đoạn | Library | Lý do |
|---|---|---|
| **MVP** | **`golang.org/x/image/draw`** | Zero deps, đơn giản. 60MB RAM × low traffic = OK |
| **Growth** | + optional `go-scaled-jpeg` | Giảm peak RAM 60→25MB nếu cần |
| **Scale** | **`bimg` (libvips)** | 4-8x nhanh, 3-5x ít RAM. CGO OK trong dedicated OCR container |

### Memory impact ở concurrent

| Concurrent requests | Pure Go (`x/image`) | libvips (`bimg`) |
|---|---|---|
| 1 | 60MB | 15MB |
| 5 | 300MB | 75MB |
| 10 | 600MB | 150MB |
| 20 | 1.2GB | 300MB |

Mitigation: Semaphore limit concurrent resize, `DecodeConfig()` guard reject ảnh quá lớn.

---

## 7. Resize logic

```
Nhận image bytes
    │
    ├── file size > 10MB → reject 400
    │
    ├── DecodeConfig() check dimensions
    │   └── > 75MP → reject 400 (prevent OOM)
    │
    ├── file size ≤ 2MB AND longest side ≤ 2048
    │   → skip resize, forward nguyên bản
    │
    └── cần resize:
        1. Decode JPEG/PNG
        2. scale = 2048.0 / longestSide
        3. Resize bằng CatmullRom (giữ aspect ratio)
        4. Encode JPEG quality 85
        5. Forward cho OCR engine

Ví dụ:
  4000×3000 → 2048×1536  (~400KB)
  3000×4000 → 1536×2048  (~400KB)
  2000×1500 → skip
  800×600   → skip (không upscale)
```

---

## 8. Phương án: Dùng Prep Image Service có sẵn

Prep platform đã có image processing service. Có thể tái sử dụng thay vì build resize riêng.

```
Mobile → [Go API] → [Prep Image Service] → resize → [Go API] → OCR engine
```

| Ưu điểm | Nhược điểm |
|---|---|
| Không cần build resize logic | Thêm 1 network hop (gửi 8MB qua internal network 2 lần) |
| Không tốn CPU/RAM trên Go API | Coupling: Prep service down → OCR down |
| Đã có infra, đã test production | Config mismatch: Prep resize cho thumbnails/avatars, OCR cần config riêng |
| | Overkill: resize cho OCR = ~20 dòng Go code |

**Kết luận: Phương án dự phòng.** MVP resize tại chỗ đơn giản hơn. Chỉ dùng khi có lý do cụ thể (centralize image processing, hoặc Prep service đã có đúng config OCR).

---

## 9. Tổng hợp: Lựa chọn cho MVP

| Quyết định | Chọn | Phương án khác đã xem xét | Lý do chọn |
|---|---|---|---|
| **Resize strategy** | Max longest side = 2048px, giữ aspect ratio | Fixed dimensions, Max area | Giữ tỉ lệ → không méo chữ. Industry standard cho OCR |
| **Interpolation** | CatmullRom | Lanczos, BiLinear, NearestNeighbor | Tốt nhất trong stdlib Go. Quality gần Lanczos, nhanh hơn nhiều |
| **Output format** | JPEG quality 85 | JPEG 95, JPEG 75, PNG | Best balance size (300-600KB) vs quality. Dưới 75 hại handwritten Hán tự |
| **Go library** | `golang.org/x/image/draw` | `imaging`, `bimg` (libvips), `go-scaled-jpeg` | Official, zero deps. Đủ cho MVP traffic. Scale phase chuyển libvips |
| **Skip resize** | File ≤ 2MB AND longest side ≤ 2048px | Luôn resize, chỉ check size, chỉ check dimensions | Tránh tốn CPU vô ích, cover cả 2 điều kiện |
| **Max file size** | 10MB | 5MB, 20MB | Cover hầu hết ảnh camera. Vừa limit Baidu URL mode (10MB) |
| **OOM guard** | Reject > 75MP trước khi decode | Không guard, limit thấp hơn | 75MP = limit Google Vision. Decode 75MP = 300MB RAM |
| **Post-processing** | Không (chỉ resize) | Sharpen, binarize, deskew, denoise | Cloud OCR engines có preprocessing riêng. Làm thêm có thể hại |
| **Resize service** | Tại chỗ trong Go API | Prep Image Service, tách OCR Service | ~20 dòng code. Tách service = overkill cho MVP |

---

## References

| Source | URL |
|---|---|
| Google Cloud Vision — Supported files | https://docs.cloud.google.com/vision/docs/supported-files |
| Baidu OCR API | https://ai.baidu.com/ai-doc/OCR/Ck3h7y2ia |
| Tesseract — Improve Quality | https://tesseract-ocr.github.io/tessdoc/ImproveQuality.html |
| How OCR Works — Compression | https://how-ocr-works.com/images/compression.html |
| golang.org/x/image/draw | https://pkg.go.dev/golang.org/x/image/draw |
| h2non/bimg (libvips) | https://github.com/h2non/bimg |
| fawick/speedtest-resize benchmarks | https://github.com/fawick/speedtest-resize |
| Go issue #10532 — JPEG memory | https://github.com/golang/go/issues/10532 |
| Roboflow — Image Resizing | https://blog.roboflow.com/you-might-be-resizing-your-images-incorrectly/ |
