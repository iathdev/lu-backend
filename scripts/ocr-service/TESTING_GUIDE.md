# OCR Service — Testing Guide

Co 2 cach test: **Python script** (goi truc tiep API, khong can server) va **HTTP** (qua endpoint `POST /test`).

## Setup

```bash
cd backend/scripts/ocr-service
python3 -m venv .venv
source .venv/bin/activate
pip install -r requirements.txt
```

---

## Cach 1: Python script (truc tiep, khong can server)

Goi thang vao engine Python API. Khong can chay uvicorn.

```bash
source .venv/bin/activate

# PaddleOCR (default, khong can setup them)
python test_paddleocr.py [image_path] [zh|en|vi]

# Tesseract (can: brew install tesseract tesseract-lang && pip install pytesseract Pillow)
python test_tesseract.py [image_path] [zh|en|vi]

# Google Cloud Vision (can: pip install google-auth google-cloud-vision)
python test_google.py [image_path] [zh|en|vi]

# GLM — Ollama local (can: ollama serve && ollama pull glm-ocr)
python test_glm.py [image_path]

# GLM — Zhipu AI cloud
export ZAI_API_KEY=your-key-here
python test_glm.py [image_path] --cloud

# DeepSeek (can: DEEPSEEK_API_KEY)
export DEEPSEEK_API_KEY=sk-...
python test_deepseek.py [image_path]

# Dolphin-v2 — ByteDance local VLM (can: pip install torch torchvision transformers accelerate)
python test_dolphin.py [image_path]

# Dolphin-v2 — 4-bit quantized (it VRAM hon, ~6GB thay vi ~16GB)
python test_dolphin.py [image_path] --4bit
```

Default: `test_image.png`, language `zh`.

---

## Cach 2: HTTP server (qua endpoint)

Can chay server truoc:

```bash
uvicorn main:app --host 0.0.0.0 --port 8000 --reload
```

Swagger UI: http://localhost:8000/docs

### 1. PaddleOCR (default)

**Khong can setup them gi. Chay ngay.**

```bash
curl -F image=@test_image.png http://localhost:8000/test
```

- Type: OCR
- Output: characters + confidence (per-char) + bounding box
- Model: PP-OCRv5

### 2. Tesseract

**Can cai them tesseract.**

```bash
# Setup (1 lan)
brew install tesseract tesseract-lang

# Test
curl -F image=@test_image.png -F engine=tesseract http://localhost:8000/test
```

- Type: OCR
- Output: characters + confidence + bounding box

### 3. Google Cloud Vision

**Khong can set env var. Dung service account `.gcp/ocr.json` (da co san).**

```bash
# Setup (1 lan, trong venv)
pip install google-auth google-cloud-vision

# Test
curl -F image=@test_image.png -F engine=google http://localhost:8000/test
```

- Type: OCR
- Output: characters + confidence + bounding box
- Auth: service account tu `.gcp/ocr.json` (auto detect)

### 4. GLM — Local via Ollama

**Free. Khong can key.**

```bash
# Setup (1 lan)
brew install ollama
ollama serve          # de terminal nay mo
ollama pull glm-ocr   # ~600MB, chi can 1 lan

# Test
curl -F image=@test_image.png -F engine=glm http://localhost:8000/test
```

- Type: VLM
- Output: raw text (khong co confidence, khong co bounding box)
- Model: `glm-ocr` (0.9B) — nhe, chay duoc tren Mac

### 5. GLM — Cloud via Zhipu AI

**Can key. Free tier.**

```bash
# 1. Dang ky tai https://open.bigmodel.cn (can email)
# 2. Tao API key

# 3. Khoi dong server voi key
export ZAI_API_KEY=your-key-here
uvicorn main:app --host 0.0.0.0 --port 8000 --reload

# 4. Test
curl -F image=@test_image.png -F engine=glm http://localhost:8000/test
```

- Type: VLM
- Output: raw text
- Model: `glm-4v-flash` (mien phi)
- Luu y: neu co `ZAI_API_KEY` -> tu dong dung cloud, khong thi fallback ve Ollama local

### 6. DeepSeek

**Can key. ~$2 free credit.**

```bash
# 1. Dang ky tai https://platform.deepseek.com
# 2. Vao API Keys -> tao key

# 3. Khoi dong server voi key
export DEEPSEEK_API_KEY=sk-...
uvicorn main:app --host 0.0.0.0 --port 8000 --reload

# 4. Test
curl -F image=@test_image.png -F engine=deepseek http://localhost:8000/test
```

- Type: VLM
- Output: raw text
- Model: `deepseek-chat` (DeepSeek-V3)
- **Luu y:** DeepSeek chua public vision API -> co the fail khi gui image

---

## Test tat ca engines cung luc (HTTP)

```bash
# Chay lan luot, so sanh ket qua
for engine in paddleocr tesseract google glm deepseek; do
  echo "=== $engine ==="
  curl -s -F image=@test_image.png -F engine=$engine http://localhost:8000/test | python3 -m json.tool
  echo ""
done
```

---

## So sanh

| Engine      | Type | Free | Can key          | Confidence | Bounding box | Toc do  | Test script         |
|-------------|------|------|------------------|------------|--------------|---------|---------------------|
| `paddleocr` | OCR  | Yes  | Khong            | Per-char   | Co           | Nhanh   | `test_paddleocr.py` |
| `tesseract` | OCR  | Yes  | Khong            | Per-word   | Co           | Nhanh   | `test_tesseract.py` |
| `google`    | OCR  | Free tier | Khong (SA auto) | Co      | Co           | Nhanh   | `test_google.py`    |
| `glm`       | VLM  | Yes  | Khong (Ollama)   | Khong      | Khong        | Tuy may | `test_glm.py`       |
| `deepseek`  | VLM  | ~$2  | Co               | Khong      | Khong        | Nhanh   | `test_deepseek.py`  |
| `dolphin`   | VLM  | Yes  | Khong (local)    | Khong      | Khong        | Tuy GPU | `test_dolphin.py`   |

## Thu tu test khuyen nghi

1. **PaddleOCR** — da co san, chay ngay (`python test_paddleocr.py`)
2. **Google Cloud Vision** — SA da co, chi can pip install (`python test_google.py`)
3. **GLM Ollama** — free, khong can key, can cai Ollama (`python test_glm.py`)
4. **Tesseract** — can brew install (`python test_tesseract.py`)
5. **GLM Cloud / DeepSeek** — can dang ky lay key (`python test_glm.py --cloud` / `python test_deepseek.py`)
6. **Dolphin-v2** — can GPU 16GB VRAM hoac dung `--4bit` (~6GB) (`python test_dolphin.py`)
