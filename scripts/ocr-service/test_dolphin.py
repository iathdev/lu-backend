"""
Test Dolphin-v2 (ByteDance) — local document parsing VLM.

Setup:
  1. Install deps: pip install torch torchvision transformers accelerate Pillow
  2. Run:          python test_dolphin.py [image_path]

First run downloads the model from HuggingFace (~6GB).
Requires ~16GB VRAM (bfloat16) or use --4bit for quantized mode (~6GB).

Default test image: test_image.png
"""

import os
import sys
import time

MODEL_ID = "ByteDance/Dolphin-v2"

PROMPTS = {
    "extract": "识别图片中的所有中文文本，逐行列出。只输出中文文字，不要解释。",
    "structured": (
        "识别图片中的所有中文文本。以JSON数组格式输出，每个元素包含："
        '{"text": "识别的文字", "type": "printed或handwritten"}。'
        "只输出JSON，不要其他内容。"
    ),
}


def load_model(use_4bit: bool = False):
    import torch
    from transformers import AutoProcessor, Qwen2_5_VLForConditionalGeneration

    print(f"Loading model: {MODEL_ID} ({'4-bit quantized' if use_4bit else 'bfloat16'})...")
    start = time.time()

    processor = AutoProcessor.from_pretrained(MODEL_ID, trust_remote_code=True)

    load_kwargs = {
        "trust_remote_code": True,
        "device_map": "auto",
    }
    if use_4bit:
        load_kwargs["load_in_4bit"] = True
    else:
        load_kwargs["dtype"] = torch.bfloat16

    model = Qwen2_5_VLForConditionalGeneration.from_pretrained(MODEL_ID, **load_kwargs)

    elapsed = time.time() - start
    print(f"Model loaded in {elapsed:.1f}s")
    return processor, model


def run_inference(processor, model, image_path: str, prompt_key: str = "extract") -> dict:
    import torch
    from PIL import Image
    from qwen_vl_utils import process_vision_info

    prompt = PROMPTS.get(prompt_key, PROMPTS["extract"])

    messages = [
        {
            "role": "user",
            "content": [
                {"type": "image", "image": f"file://{os.path.abspath(image_path)}"},
                {"type": "text", "text": prompt},
            ],
        }
    ]

    text = processor.apply_chat_template(messages, tokenize=False, add_generation_prompt=True)
    image_inputs, video_inputs = process_vision_info(messages)
    inputs = processor(
        text=[text],
        images=image_inputs,
        videos=video_inputs,
        return_tensors="pt",
    ).to(model.device)

    print(f"Running inference ({prompt_key} mode)...")
    start = time.time()

    with torch.no_grad():
        generated_ids = model.generate(
            **inputs,
            max_new_tokens=4096,
            do_sample=False,
        )

    output_text = processor.batch_decode(
        generated_ids[:, inputs["input_ids"].shape[1]:],
        skip_special_tokens=True,
    )[0]

    elapsed = time.time() - start
    return {"output": output_text, "elapsed": elapsed}


def main():
    use_4bit = "--4bit" in sys.argv
    args = [a for a in sys.argv[1:] if a != "--4bit"]
    image_path = args[0] if args else "test_image.png"

    if not os.path.exists(image_path):
        print(f"ERROR: Image not found: {image_path}")
        sys.exit(1)

    print(f"Image:  {image_path}")
    print(f"Size:   {os.path.getsize(image_path) / 1024:.1f} KB")
    print(f"Model:  {MODEL_ID}")
    print("=" * 60)

    processor, model = load_model(use_4bit=use_4bit)

    # Test 1: Simple extraction
    resp = run_inference(processor, model, image_path, "extract")

    print(f"\n--- Extract mode ({resp['elapsed']:.2f}s) ---")
    print(f"Output:\n{resp['output']}")

    # Test 2: Structured JSON
    print("\n" + "=" * 60)
    resp2 = run_inference(processor, model, image_path, "structured")

    print(f"\n--- Structured mode ({resp2['elapsed']:.2f}s) ---")
    print(f"Output:\n{resp2['output']}")

    # Summary
    print("\n" + "=" * 60)
    print("SUMMARY:")
    print(f"  Model:                {MODEL_ID}")
    print(f"  Quantized:            {'4-bit' if use_4bit else 'bfloat16'}")
    print(f"  Latency (extract):    {resp['elapsed']:.2f}s")
    print(f"  Latency (structured): {resp2['elapsed']:.2f}s")
    print(f"  Per-char confidence:  NOT AVAILABLE (VLM generative output)")
    print(f"  Bounding boxes:       NOT AVAILABLE")


if __name__ == "__main__":
    main()
