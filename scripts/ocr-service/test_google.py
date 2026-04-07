"""
Test Google Cloud Vision OCR directly via Python API.

Setup:
  1. Service account: .gcp/ocr.json (da co san)
  2. Install: pip install google-auth google-cloud-vision
  3. Run: python test_google.py [image_path]

Default test image: test_image.png
"""

import os
import re
import sys
import time

from pypinyin import pinyin, Style

CJK_PATTERN = re.compile(r"[\u4e00-\u9fff]+")


def to_pinyin(text: str) -> str:
    return " ".join(syllable[0] for syllable in pinyin(text, style=Style.TONE))


def get_client():
    try:
        from google.cloud import vision
        from google.oauth2 import service_account
    except ImportError:
        print("ERROR: google-cloud-vision not installed.")
        print("  pip install google-auth google-cloud-vision")
        sys.exit(1)

    sa_path = os.path.join(os.path.dirname(__file__), "..", "..", ".gcp", "ocr.json")
    sa_path = os.path.abspath(sa_path)
    if not os.path.exists(sa_path):
        print(f"ERROR: Service account not found: {sa_path}")
        sys.exit(1)

    print(f"Service account: {sa_path}")
    credentials = service_account.Credentials.from_service_account_file(
        sa_path, scopes=["https://www.googleapis.com/auth/cloud-platform"],
    )
    return vision.ImageAnnotatorClient(credentials=credentials)


def recognize(image_path: str, language: str = "zh") -> dict:
    from google.cloud import vision

    client = get_client()

    with open(image_path, "rb") as f:
        content = f.read()

    image = vision.Image(content=content)
    lang_hints = {"zh": ["zh"], "en": ["en"], "vi": ["vi"]}
    hints = lang_hints.get(language, [])
    image_context = vision.ImageContext(language_hints=hints) if hints else None

    print(f"Calling Google Cloud Vision (lang={language})...")
    start = time.time()
    response = client.text_detection(image=image, image_context=image_context)
    elapsed = time.time() - start

    if response.error.message:
        print(f"ERROR: Google Vision API error: {response.error.message}")
        sys.exit(1)

    annotations = []
    for ann in response.text_annotations:
        verts = ann.bounding_poly.vertices
        annotations.append({
            "text": ann.description,
            "score": ann.score,
            "locale": ann.locale,
            "bounding_box": [[v.x, v.y] for v in verts] if verts else [],
        })

    return {"annotations": annotations, "elapsed": elapsed}


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
    annotations = resp["annotations"]

    print(f"\n--- Results ({resp['elapsed']:.2f}s) ---")
    print(f"Annotations: {len(annotations)}")

    if annotations:
        # First annotation is full text
        full_text = annotations[0]["text"]
        print(f"\n  [Full text]")
        for line in full_text.strip().split("\n"):
            print(f"    {line}")
            if language == "zh" and CJK_PATTERN.search(line):
                print(f"    -> {to_pinyin(line)}")

        # Individual words/characters
        print(f"\n  [Individual annotations: {len(annotations) - 1}]")
        for i, ann in enumerate(annotations[1:], 1):
            text = ann["text"]
            bbox = ann["bounding_box"]
            locale = ann["locale"]
            score = ann["score"]
            extra = []
            if locale:
                extra.append(f"locale={locale}")
            if score:
                extra.append(f"score={score}")
            extra_str = f" ({', '.join(extra)})" if extra else ""
            print(f"    [{i}] {text}{extra_str}  bbox={bbox}")

    # Summary
    print("\n" + "=" * 60)
    print("SUMMARY:")
    print(f"  Engine:              Google Cloud Vision")
    print(f"  Language:            {language}")
    print(f"  Latency:             {resp['elapsed']:.2f}s")
    print(f"  Annotations:         {len(annotations)}")
    print(f"  Per-char confidence: YES (per-annotation)")
    print(f"  Bounding boxes:      YES")


if __name__ == "__main__":
    main()
