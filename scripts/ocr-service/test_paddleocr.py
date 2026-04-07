"""
Test PaddleOCR PP-OCRv5 directly via Python API.

Setup:
  1. Activate venv: source .venv/bin/activate
  2. Run: python test_paddleocr.py [image_path]

Default test image: test_image.png
"""

import os
import re
import sys
import time

import numpy as np
from pypinyin import pinyin, Style

CJK_PATTERN = re.compile(r"[\u4e00-\u9fff]+")
LANG_MAP = {"zh": "ch", "en": "en", "vi": "vi"}


def get_engine(language: str):
    from paddleocr import PaddleOCR

    paddle_lang = LANG_MAP.get(language, language)
    return PaddleOCR(lang=paddle_lang, use_textline_orientation=True)


def to_pinyin(text: str) -> str:
    return " ".join(syllable[0] for syllable in pinyin(text, style=Style.TONE))


def recognize(image_path: str, language: str = "zh") -> dict:
    engine = get_engine(language)

    print(f"Running PaddleOCR (lang={language})...")
    start = time.time()
    results = engine.predict(image_path)
    elapsed = time.time() - start

    lines = []
    for result in results:
        texts = result.get("rec_texts", [])
        scores = result.get("rec_scores", [])
        char_scores_list = result.get("rec_char_scores", [[] for _ in texts])
        rec_polys = result.get("rec_polys", [None] * len(texts))

        for text, score, char_scores, poly in zip(texts, scores, char_scores_list, rec_polys):
            location = _poly_to_location(poly)
            lines.append({
                "text": text,
                "confidence": round(float(score), 4),
                "char_scores": [round(float(s), 4) for s in char_scores] if char_scores else [],
                "location": location,
            })

    return {"lines": lines, "elapsed": elapsed}


def _poly_to_location(poly) -> list[list[float]]:
    if poly is None:
        return []
    try:
        pts = np.array(poly)
        if pts.ndim == 2 and pts.shape[0] >= 4 and pts.shape[1] >= 2:
            return [[round(float(pts[i][0])), round(float(pts[i][1]))] for i in range(4)]
    except (ValueError, TypeError):
        pass
    return []


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
        if line["char_scores"]:
            chars = list(line["text"])
            scores = line["char_scores"]
            pairs = [f"{c}={s:.2f}" for c, s in zip(chars, scores)]
            print(f"      Per-char:   {', '.join(pairs)}")
        if line["location"]:
            print(f"      Location:   {line['location']}")

    # Summary
    print("\n" + "=" * 60)
    print("SUMMARY:")
    print(f"  Engine:              PaddleOCR PP-OCRv5")
    print(f"  Language:            {language}")
    print(f"  Latency:             {resp['elapsed']:.2f}s")
    print(f"  Lines detected:      {len(resp['lines'])}")
    print(f"  Per-char confidence: YES")
    print(f"  Bounding boxes:      YES")

    if resp["lines"]:
        avg_conf = sum(l["confidence"] for l in resp["lines"]) / len(resp["lines"])
        print(f"  Avg confidence:      {avg_conf:.4f}")


if __name__ == "__main__":
    main()
