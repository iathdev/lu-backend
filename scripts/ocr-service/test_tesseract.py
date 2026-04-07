"""
Test Tesseract OCR directly via Python API.

Setup:
  1. Install Tesseract: brew install tesseract tesseract-lang
  2. Install pytesseract: pip install pytesseract Pillow
  3. Run: python test_tesseract.py [image_path]

Default test image: test_image.png
"""

import os
import re
import sys
import time

from pypinyin import pinyin, Style

CJK_PATTERN = re.compile(r"[\u4e00-\u9fff]+")
TESS_LANG_MAP = {"zh": "chi_sim", "en": "eng", "vi": "vie"}


def to_pinyin(text: str) -> str:
    return " ".join(syllable[0] for syllable in pinyin(text, style=Style.TONE))


def recognize(image_path: str, language: str = "zh") -> dict:
    try:
        import pytesseract
    except ImportError:
        print("ERROR: pytesseract not installed.")
        print("  pip install pytesseract Pillow")
        print("  brew install tesseract tesseract-lang")
        sys.exit(1)

    from PIL import Image

    tess_lang = TESS_LANG_MAP.get(language, "eng")
    img = Image.open(image_path)

    print(f"Running Tesseract (lang={tess_lang})...")
    start = time.time()
    data = pytesseract.image_to_data(img, lang=tess_lang, output_type=pytesseract.Output.DICT)
    elapsed = time.time() - start

    # Group words by line
    lines: dict[tuple[int, int, int], list[dict]] = {}
    for i, text in enumerate(data["text"]):
        if data["level"][i] != 5:
            continue
        text = text.strip()
        if not text:
            continue
        key = (data["block_num"][i], data["par_num"][i], data["line_num"][i])
        if key not in lines:
            lines[key] = []
        lines[key].append({
            "text": text,
            "left": data["left"][i],
            "top": data["top"][i],
            "width": data["width"][i],
            "height": data["height"][i],
            "conf": data["conf"][i],
        })

    results = []
    for words in lines.values():
        line_text = " ".join(w["text"] for w in words)
        left = min(w["left"] for w in words)
        top = min(w["top"] for w in words)
        right = max(w["left"] + w["width"] for w in words)
        bottom = max(w["top"] + w["height"] for w in words)
        avg_conf = sum(max(0.0, w["conf"]) for w in words) / len(words) / 100.0
        location = [[left, top], [right, top], [right, bottom], [left, bottom]]
        results.append({
            "text": line_text,
            "confidence": round(avg_conf, 4),
            "location": location,
            "words": words,
        })

    return {"lines": results, "elapsed": elapsed}


def main():
    image_path = sys.argv[1] if len(sys.argv) > 1 else "test_image.png"
    language = sys.argv[2] if len(sys.argv) > 2 else "zh"

    if not os.path.exists(image_path):
        print(f"ERROR: Image not found: {image_path}")
        sys.exit(1)

    print(f"Image:    {image_path}")
    print(f"Size:     {os.path.getsize(image_path) / 1024:.1f} KB")
    print(f"Language: {language}")
    print("=" * 60)

    resp = recognize(image_path, language)

    print(f"\n--- Results ({resp['elapsed']:.2f}s) ---")
    print(f"Lines detected: {len(resp['lines'])}")

    for i, line in enumerate(resp["lines"]):
        print(f"\n  [{i + 1}] {line['text']}")
        if language == "zh" and CJK_PATTERN.search(line["text"]):
            print(f"      Pinyin:     {to_pinyin(line['text'])}")
        print(f"      Confidence: {line['confidence']}")
        print(f"      Location:   {line['location']}")
        if len(line["words"]) > 1:
            word_details = [f"{w['text']}={w['conf']}%" for w in line["words"]]
            print(f"      Words:      {', '.join(word_details)}")

    # Summary
    print("\n" + "=" * 60)
    print("SUMMARY:")
    print(f"  Engine:              Tesseract")
    print(f"  Language:            {language} ({TESS_LANG_MAP.get(language, 'eng')})")
    print(f"  Latency:             {resp['elapsed']:.2f}s")
    print(f"  Lines detected:      {len(resp['lines'])}")
    print(f"  Per-char confidence: NO (per-word only)")
    print(f"  Bounding boxes:      YES")

    if resp["lines"]:
        avg_conf = sum(l["confidence"] for l in resp["lines"]) / len(resp["lines"])
        print(f"  Avg confidence:      {avg_conf:.4f}")


if __name__ == "__main__":
    main()
