#!/usr/bin/env python3
# -*- coding: utf-8 -*-

"""
Генератор интерактивной презентации (HTML 16:9 + PDF) из Markdown-файла PRESENTATION.md.
Работает полностью на стандартной библиотеке Python (без сторонних зависимостей).
Автоматически кодирует скриншоты в base64 для создания автономного Single-File HTML.
"""

import os
import re
import sys
import base64
import subprocess
from pathlib import Path

BASE_DIR = Path(__file__).resolve().parent
MD_PATH = BASE_DIR / "PRESENTATION.md"
HTML_OUT = BASE_DIR / "presentation.html"
PDF_OUT = BASE_DIR / "presentation.pdf"

def img_to_base64(img_path_str: str) -> str:
    path = Path(img_path_str)
    if not path.is_absolute():
        path = BASE_DIR / path
    if path.exists():
        ext = path.suffix.lower().replace(".", "")
        if ext == "jpg":
            ext = "jpeg"
        mime = f"image/{ext}"
        try:
            with open(path, "rb") as f:
                b64 = base64.b64encode(f.read()).decode("utf-8")
                return f"data:{mime};base64,{b64}"
        except Exception as e:
            print(f"Warning: could not encode image {path}: {e}")
    return img_path_str

def parse_inline_markdown(text: str) -> str:
    # Arrow and math replacements
    text = text.replace(r'$\rightarrow$', '→').replace(r'\rightarrow', '→')
    text = text.replace(r'$\to$', '→').replace(r'\to', '→')
    text = text.replace(r'$\leftarrow$', '←').replace(r'\leftarrow', '←')
    text = text.replace(r'\ge', '≥').replace(r'\le', '≤').replace(r'\in', '∈').replace(r'\pm', '±')
    text = text.replace(r'$\ge$', '≥').replace(r'$\le$', '≤')
    # Bold
    text = re.sub(r'\*\*(.+?)\*\*', r'<strong>\1</strong>', text)
    # Italic
    text = re.sub(r'\*(.+?)\*', r'<em>\1</em>', text)
    # Inline code
    text = re.sub(r'`([^`]+)`', r'<code>\1</code>', text)
    text = re.sub(r'\$\$([\s\S]+?)\$\$', r'<div class="formula-block">\1</div>', text)
    text = re.sub(r'\$([^\$]+)\$', r'<span class="formula-inline">\1</span>', text)
    # Clean dangling dollar signs
    text = text.replace('$', '')
    return text

def parse_markdown_slide(raw_slide: str, slide_index: int) -> dict:
    lines = [line.strip() for line in raw_slide.strip().split("\n") if line.strip() != ""]
    if not lines:
        return None

    title = ""
    content_html = []
    images = []
    is_title_slide = (slide_index == 1)

    i = 0
    while i < len(lines):
        line = lines[i]

        # Slide header ## Слайд X. Title
        header_match = re.match(r'^##\s*(?:Слайд\s*\d+\.\s*)?(.*)', line)
        if header_match:
            title = header_match.group(1).strip()
            i += 1
            continue

        # Subtitle or h3
        if line.startswith("### "):
            sub = line[4:].strip()
            content_html.append(f'<h3 class="slide-subtitle">{parse_inline_markdown(sub)}</h3>')
            i += 1
            continue

        # Images: ![alt](url)
        img_match = re.match(r'^!\[(.*?)\]\((.*?)\)', line)
        if img_match:
            alt, src = img_match.groups()
            b64_src = img_to_base64(src)
            caption = ""
            if i + 1 < len(lines) and lines[i+1].startswith("*") and lines[i+1].endswith("*"):
                caption = lines[i+1].strip("*_ ")
                i += 1
            images.append({"alt": alt, "src": b64_src, "caption": caption})
            i += 1
            continue

        # Blockquote: > text
        if line.startswith("> "):
            bq_text = line[2:].strip()
            content_html.append(f'<div class="info-badge">{parse_inline_markdown(bq_text)}</div>')
            i += 1
            continue

        # Code block / diagram: ``` ... ```
        if line.startswith("```"):
            code_lines = []
            i += 1
            while i < len(lines) and not lines[i].startswith("```"):
                code_lines.append(lines[i])
                i += 1
            code_str = "\n".join(code_lines)
            content_html.append(f'<pre class="code-block"><code>{code_str}</code></pre>')
            i += 1
            continue

        # Table: | a | b |
        if line.startswith("|") and line.endswith("|"):
            table_lines = [line]
            i += 1
            while i < len(lines) and lines[i].startswith("|") and lines[i].endswith("|"):
                table_lines.append(lines[i])
                i += 1
            
            # Convert table lines to HTML
            if len(table_lines) >= 2:
                header_cols = [c.strip() for c in table_lines[0].strip("|").split("|")]
                rows = []
                for row_line in table_lines[2:]:
                    cols = [c.strip() for c in row_line.strip("|").split("|")]
                    rows.append(cols)
                
                table_html = ['<table class="slide-table"><thead><tr>']
                for c in header_cols:
                    table_html.append(f'<th>{parse_inline_markdown(c)}</th>')
                table_html.append('</tr></thead><tbody>')
                for r in rows:
                    table_html.append('<tr>')
                    for c in r:
                        table_html.append(f'<td>{parse_inline_markdown(c)}</td>')
                    table_html.append('</tr>')
                table_html.append('</tbody></table>')
                content_html.append("".join(table_html))
            continue

        # Unordered list: - item
        if line.startswith("- ") or line.startswith("* "):
            list_items = [line[2:].strip()]
            i += 1
            while i < len(lines) and (lines[i].startswith("- ") or lines[i].startswith("* ")):
                list_items.append(lines[i][2:].strip())
                i += 1
            ul_html = ['<ul class="slide-list">']
            for item in list_items:
                ul_html.append(f'<li>{parse_inline_markdown(item)}</li>')
            ul_html.append('</ul>')
            content_html.append("".join(ul_html))
            continue

        # Ordered list: 1. item
        ol_match = re.match(r'^\d+\.\s+(.*)', line)
        if ol_match:
            list_items = [ol_match.group(1)]
            i += 1
            while i < len(lines):
                m = re.match(r'^\d+\.\s+(.*)', lines[i])
                if m:
                    list_items.append(m.group(1))
                    i += 1
                else:
                    break
            ol_html = ['<ol class="slide-list ol-list">']
            for item in list_items:
                ol_html.append(f'<li>{parse_inline_markdown(item)}</li>')
            ol_html.append('</ol>')
            content_html.append("".join(ol_html))
            continue

        # Regular paragraph
        content_html.append(f'<p class="slide-p">{parse_inline_markdown(line)}</p>')
        i += 1

    return {
        "index": slide_index,
        "title": title,
        "is_title_slide": is_title_slide,
        "content_html": "".join(content_html),
        "images": images
    }

def get_architecture_diagram_html() -> str:
    return """
    <div class="arch-flow-diagram">
      <div class="arch-card client-card">
        <div class="arch-card-header">
          <span class="arch-badge client-badge">Клиентский уровень</span>
          <span class="arch-tech-pill">SPA • Port :3000</span>
        </div>
        <div class="arch-card-title">React 19 Frontend (TypeScript)</div>
        <div class="arch-card-tags">
          <span class="tag">Feature-Sliced Design</span>
          <span class="tag">MobX State</span>
          <span class="tag">MUI & DataGrid</span>
          <span class="tag">Axios</span>
        </div>
      </div>

      <div class="arch-arrow-bridge">
        <span class="bridge-line"></span>
        <span class="bridge-label">REST API (JSON / HTTP)</span>
        <span class="bridge-pointer">▼</span>
      </div>

      <div class="arch-card server-card">
        <div class="arch-card-header">
          <span class="arch-badge server-badge">Серверный уровень</span>
          <span class="arch-tech-pill">Go 1.27+ • Port :8080</span>
        </div>
        <div class="arch-card-title">Gin Web API (Clean Architecture)</div>
        <div class="arch-subgrid">
          <div class="submodule">🔐 JWT & bcrypt Auth</div>
          <div class="submodule">⚖ Scoring Engine (1–10)</div>
          <div class="submodule">📁 Multipart Uploads</div>
          <div class="submodule">📊 Analytics Aggregator</div>
        </div>
      </div>

      <div class="arch-arrow-bridge">
        <span class="bridge-line"></span>
        <span class="bridge-label">SQL Transactions / File Storage</span>
        <span class="bridge-pointer">▼</span>
      </div>

      <div class="arch-data-row">
        <div class="arch-card data-card">
          <div class="arch-card-header">
            <span class="arch-badge data-badge">СУБД</span>
          </div>
          <div class="arch-card-title">SQLite / PostgreSQL</div>
          <div class="arch-desc">Авто-миграции, внешние ключи, сиды</div>
        </div>
        <div class="arch-card storage-card">
          <div class="arch-card-header">
            <span class="arch-badge storage-badge">Файлы</span>
          </div>
          <div class="arch-card-title">Uploads Storage</div>
          <div class="arch-desc">Хранилище фото- и видеоматериалов</div>
        </div>
      </div>
    </div>
    """

def generate_presentation():
    if not MD_PATH.exists():
        print(f"Error: {MD_PATH} not found.")
        sys.exit(1)

    with open(MD_PATH, "r", encoding="utf-8") as f:
        md_text = f.read()

    # Split into raw slides by '---'
    raw_slides = re.split(r'\n---\n', md_text)
    slides_data = []
    
    slide_idx = 1
    for raw in raw_slides:
        raw = raw.strip()
        if not raw:
            continue
        slide = parse_markdown_slide(raw, slide_idx)
        if slide:
            slides_data.append(slide)
            slide_idx += 1

    total_slides = len(slides_data)
    print(f"Parsed {total_slides} slides from {MD_PATH.name}")

    slides_html_parts = []
    for s in slides_data:
        idx = s["index"]
        title = s["title"]
        has_images = len(s["images"]) > 0
        is_arch_slide = ("Архитектура" in title)
        
        layout_class = "normal-layout"
        if s["is_title_slide"]:
            layout_class = "title-layout"
        elif is_arch_slide or has_images:
            layout_class = "split-layout"

        slide_dom = [f'<section class="slide {layout_class}" id="slide-{idx}" data-slide="{idx}">']
        
        # Header (except title slide)
        if not s["is_title_slide"]:
            slide_dom.append(f'''
            <header class="slide-header">
                <span class="slide-topic">DefectTrack — Defect Management System</span>
                <h2 class="slide-title">{title}</h2>
            </header>
            ''')

        slide_dom.append('<div class="slide-body">')

        if s["is_title_slide"]:
            slide_dom.append(f'''
            <div class="title-slide-content">
                <div class="project-tag">Строительный контроль & Инженерия 2026</div>
                <h1 class="main-title">DefectTrack</h1>
                <h2 class="sub-title">Defect Management System</h2>
                <p class="title-desc">Полнофункциональная система приёма, алгоритмического скоринга и диспетчеризации строительных рекламаций</p>
                <div class="title-meta">
                    <div class="meta-card">
                        <span class="meta-label">Разработчик:</span>
                        <span class="meta-val">студент гр. бИВТ-231 <b>Машнев М.</b></span>
                    </div>
                    <div class="meta-card">
                        <span class="meta-label">Научный руководитель:</span>
                        <span class="meta-val"><b>Сафронов В. В.</b></span>
                    </div>
                    <div class="meta-card">
                        <span class="meta-label">Организация:</span>
                        <span class="meta-val">ВГТУ, кафедра АВС</span>
                    </div>
                </div>
            </div>
            ''')
        elif is_arch_slide:
            slide_dom.append('<div class="content-col">')
            slide_dom.append(s["content_html"])
            slide_dom.append('</div>')
            slide_dom.append('<div class="image-col">')
            slide_dom.append(get_architecture_diagram_html())
            slide_dom.append('</div>')
        elif has_images:
            slide_dom.append('<div class="content-col">')
            slide_dom.append(s["content_html"])
            slide_dom.append('</div>')
            slide_dom.append('<div class="image-col">')
            for img in s["images"]:
                caption_dom = f'<div class="img-caption">{img["caption"]}</div>' if img["caption"] else ""
                slide_dom.append(f'''
                <div class="slide-image-wrapper">
                    <img src="{img['src']}" alt="{img['alt']}" class="slide-img" onclick="openLightbox(this.src)" />
                    {caption_dom}
                </div>
                ''')
            slide_dom.append('</div>')
        else:
            slide_dom.append('<div class="full-content-col">')
            slide_dom.append(s["content_html"])
            slide_dom.append('</div>')

        slide_dom.append('</div>') # end slide-body

        # Slide footer
        slide_dom.append(f'''
        <footer class="slide-footer">
            <span class="footer-left">ВГТУ • Кафедра АВС • Машнев М.</span>
            <span class="footer-center">Техническое задание на разработку DefectTrack</span>
            <span class="footer-right">{idx} / {total_slides}</span>
        </footer>
        ''')

        slide_dom.append('</section>')
        slides_html_parts.append("".join(slide_dom))

    all_slides_html = "\n".join(slides_html_parts)

    html_template = f"""<!DOCTYPE html>
<html lang="ru">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>Презентация DefectTrack — Defect Management System</title>
<style>
  :root {{
    --bg-dark: #0a1128;
    --primary: #1976d2;
    --primary-light: #4791db;
    --primary-dark: #115293;
    --accent: #00e5ff;
    --text-main: #1f2937;
    --text-muted: #6b7280;
    --card-bg: #ffffff;
    --slide-ratio: 16 / 9;
  }}

  * {{
    box-sizing: border-box;
    margin: 0;
    padding: 0;
  }}

  body {{
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif;
    background-color: var(--bg-dark);
    color: var(--text-main);
    overflow: hidden;
    height: 100vh;
    width: 100vw;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
  }}

  /* Presentation Stage (16:9 Aspect Ratio) */
  #presentation-container {{
    position: relative;
    width: 95vw;
    height: calc(95vw * 9 / 16);
    max-height: 92vh;
    max-width: calc(92vh * 16 / 9);
    background: #ffffff;
    box-shadow: 0 25px 60px rgba(0,0,0,0.6);
    border-radius: 12px;
    overflow: hidden;
  }}

  /* Slide layout */
  .slide {{
    display: none;
    width: 100%;
    height: 100%;
    position: absolute;
    top: 0;
    left: 0;
    padding: 34px 44px 24px 44px;
    flex-direction: column;
    justify-content: space-between;
    background: #ffffff;
    background-image: radial-gradient(#f0f4f8 1px, transparent 1px);
    background-size: 24px 24px;
    animation: fadeIn 0.25s ease-in-out;
  }}

  .slide.active {{
    display: flex;
  }}

  @keyframes fadeIn {{
    from {{ opacity: 0.6; transform: scale(0.995); }}
    to {{ opacity: 1; transform: scale(1); }}
  }}

  /* Header */
  .slide-header {{
    border-bottom: 2px solid #e2e8f0;
    padding-bottom: 8px;
    margin-bottom: 12px;
  }}

  .slide-topic {{
    font-size: 11px;
    text-transform: uppercase;
    letter-spacing: 1.5px;
    color: var(--primary);
    font-weight: 700;
  }}

  .slide-title {{
    font-size: 24px;
    font-weight: 800;
    color: #0f172a;
    margin-top: 2px;
  }}

  /* Body */
  .slide-body {{
    flex: 1;
    overflow: hidden;
    display: flex;
    gap: 26px;
  }}

  .full-content-col {{
    width: 100%;
    display: flex;
    flex-direction: column;
    justify-content: center;
    gap: 12px;
  }}

  .split-layout .slide-body {{
    display: grid;
    grid-template-columns: 1.05fr 1fr;
    gap: 24px;
    align-items: center;
  }}

  .content-col {{
    display: flex;
    flex-direction: column;
    justify-content: center;
    gap: 12px;
    overflow-y: auto;
  }}

  .image-col {{
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    height: 100%;
    width: 100%;
  }}

  .slide-image-wrapper {{
    width: 100%;
    height: 100%;
    max-height: 52vh;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    background: #f8fafc;
    border: 1px solid #e2e8f0;
    border-radius: 8px;
    padding: 8px;
    box-shadow: 0 4px 14px rgba(0,0,0,0.06);
  }}

  .slide-img {{
    max-width: 100%;
    max-height: calc(52vh - 30px);
    width: auto;
    height: auto;
    border-radius: 6px;
    object-fit: contain;
    cursor: zoom-in;
    transition: transform 0.2s;
  }}

  .slide-img:hover {{
    transform: scale(1.015);
  }}

  .img-caption {{
    font-size: 12px;
    color: #475569;
    margin-top: 6px;
    font-style: italic;
    font-weight: 500;
    text-align: center;
  }}

  /* Title Slide Specifics */
  .title-layout {{
    background: linear-gradient(135deg, #091a38 0%, #102a5c 50%, #1e4976 100%);
    color: #ffffff;
    padding: 50px 60px;
    justify-content: center;
    align-items: center;
  }}

  .title-slide-content {{
    text-align: center;
    max-width: 900px;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 14px;
  }}

  .project-tag {{
    background: rgba(0, 229, 255, 0.15);
    color: var(--accent);
    padding: 6px 16px;
    border-radius: 20px;
    font-size: 12px;
    font-weight: 700;
    text-transform: uppercase;
    letter-spacing: 2px;
    border: 1px solid rgba(0, 229, 255, 0.3);
  }}

  .main-title {{
    font-size: 52px;
    font-weight: 900;
    letter-spacing: -1px;
    background: linear-gradient(90deg, #ffffff, #90caf9, #00e5ff);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
  }}

  .sub-title {{
    font-size: 22px;
    font-weight: 500;
    color: #e2e8f0;
  }}

  .title-desc {{
    font-size: 15px;
    color: #cbd5e1;
    max-width: 700px;
    line-height: 1.5;
  }}

  .title-meta {{
    display: flex;
    gap: 16px;
    margin-top: 20px;
  }}

  .meta-card {{
    background: rgba(255, 255, 255, 0.08);
    backdrop-filter: blur(8px);
    border: 1px solid rgba(255, 255, 255, 0.15);
    padding: 10px 18px;
    border-radius: 8px;
    font-size: 13px;
    display: flex;
    flex-direction: column;
    gap: 3px;
  }}

  .meta-label {{
    color: #94a3b8;
    font-size: 11px;
    text-transform: uppercase;
  }}

  .meta-val {{
    color: #ffffff;
    font-size: 13px;
  }}

  /* Content typography */
  .slide-subtitle {{
    font-size: 17px;
    color: var(--primary);
    font-weight: 700;
    margin-top: 4px;
    margin-bottom: 4px;
  }}

  .slide-p {{
    font-size: 14.5px;
    line-height: 1.55;
    color: #334155;
  }}

  .slide-list {{
    list-style: none;
    padding-left: 0;
    display: flex;
    flex-direction: column;
    gap: 9px;
  }}

  .slide-list li {{
    font-size: 14px;
    line-height: 1.5;
    color: #1e293b;
    position: relative;
    padding-left: 20px;
  }}

  .slide-list:not(.ol-list) li::before {{
    content: "•";
    color: var(--primary);
    font-weight: bold;
    font-size: 18px;
    position: absolute;
    left: 4px;
    top: -2px;
  }}

  .ol-list {{
    counter-reset: ol-counter;
  }}

  .ol-list li {{
    counter-increment: ol-counter;
  }}

  .ol-list li::before {{
    content: counter(ol-counter) ".";
    color: var(--primary);
    font-weight: 700;
    position: absolute;
    left: 0;
  }}

  .info-badge {{
    background: #e0f2fe;
    border-left: 4px solid var(--primary);
    padding: 8px 12px;
    border-radius: 0 6px 6px 0;
    font-size: 13.5px;
    color: #0369a1;
    font-weight: 500;
  }}

  .formula-block {{
    background: #f1f5f9;
    border: 1px dashed #cbd5e1;
    padding: 12px 18px;
    border-radius: 8px;
    font-family: 'Cambria Math', Georgia, serif;
    font-size: 16px;
    color: #0f172a;
    font-style: italic;
    text-align: center;
    box-shadow: inset 0 2px 4px rgba(0,0,0,0.02);
  }}

  code {{
    background: #f1f5f9;
    color: #0f172a;
    padding: 2px 5px;
    border-radius: 4px;
    font-size: 13px;
    font-family: Consolas, monospace;
  }}

  /* Table */
  .slide-table {{
    width: 100%;
    border-collapse: collapse;
    font-size: 13px;
    margin-top: 6px;
  }}

  .slide-table th, .slide-table td {{
    padding: 7px 12px;
    border: 1px solid #e2e8f0;
    text-align: left;
  }}

  .slide-table th {{
    background: #f8fafc;
    color: #0f172a;
    font-weight: 700;
  }}

  .slide-table tr:nth-child(even) td {{
    background: #fafafa;
  }}

  /* Architecture Diagram Styles */
  .arch-flow-diagram {{
    width: 100%;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 8px;
    background: #f8fafc;
    border: 1px solid #e2e8f0;
    border-radius: 10px;
    padding: 14px 18px;
    box-shadow: 0 4px 14px rgba(0,0,0,0.04);
  }}

  .arch-card {{
    width: 100%;
    border-radius: 8px;
    padding: 10px 14px;
    border: 1px solid #cbd5e1;
    background: #ffffff;
    box-shadow: 0 2px 6px rgba(0,0,0,0.03);
  }}

  .client-card {{
    border-left: 5px solid #3b82f6;
  }}

  .server-card {{
    border-left: 5px solid #10b981;
  }}

  .data-card {{
    border-left: 4px solid #f59e0b;
    flex: 1;
  }}

  .storage-card {{
    border-left: 4px solid #8b5cf6;
    flex: 1;
  }}

  .arch-card-header {{
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 4px;
  }}

  .arch-badge {{
    font-size: 10px;
    font-weight: 800;
    text-transform: uppercase;
    letter-spacing: 0.5px;
    padding: 2px 7px;
    border-radius: 4px;
  }}

  .client-badge {{ background: #dbeafe; color: #1e40af; }}
  .server-badge {{ background: #d1fae5; color: #065f46; }}
  .data-badge {{ background: #fef3c7; color: #92400e; }}
  .storage-badge {{ background: #ede9fe; color: #5b21b6; }}

  .arch-tech-pill {{
    font-size: 10.5px;
    color: #64748b;
    font-weight: 600;
  }}

  .arch-card-title {{
    font-size: 13.5px;
    font-weight: 700;
    color: #0f172a;
    margin-bottom: 6px;
  }}

  .arch-card-tags {{
    display: flex;
    flex-wrap: wrap;
    gap: 5px;
  }}

  .arch-card-tags .tag {{
    background: #f1f5f9;
    color: #334155;
    font-size: 11px;
    padding: 2px 8px;
    border-radius: 10px;
    font-weight: 500;
  }}

  .arch-subgrid {{
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 5px;
  }}

  .submodule {{
    background: #f8fafc;
    border: 1px solid #e2e8f0;
    border-radius: 4px;
    font-size: 11px;
    padding: 3px 6px;
    color: #1e293b;
    font-weight: 600;
  }}

  .arch-arrow-bridge {{
    display: flex;
    align-items: center;
    gap: 8px;
    color: #94a3b8;
    font-size: 11px;
    font-weight: 600;
    padding: 2px 0;
  }}

  .bridge-line {{
    width: 24px;
    height: 1px;
    background: #cbd5e1;
  }}

  .bridge-pointer {{
    color: #3b82f6;
    font-size: 10px;
  }}

  .arch-data-row {{
    display: flex;
    width: 100%;
    gap: 10px;
  }}

  .arch-desc {{
    font-size: 11px;
    color: #64748b;
    line-height: 1.3;
  }}

  /* Footer */
  .slide-footer {{
    border-top: 1px solid #f1f5f9;
    padding-top: 8px;
    display: flex;
    justify-content: space-between;
    align-items: center;
    font-size: 11px;
    color: #94a3b8;
  }}

  .title-layout .slide-footer {{
    border-top: 1px solid rgba(255,255,255,0.1);
    color: rgba(255,255,255,0.4);
  }}

  /* Controls */
  .controls-bar {{
    position: fixed;
    bottom: 12px;
    display: flex;
    align-items: center;
    gap: 12px;
    background: rgba(15, 23, 42, 0.85);
    backdrop-filter: blur(10px);
    border: 1px solid rgba(255, 255, 255, 0.15);
    padding: 8px 16px;
    border-radius: 30px;
    box-shadow: 0 10px 25px rgba(0,0,0,0.4);
    z-index: 1000;
  }}

  .ctrl-btn {{
    background: rgba(255,255,255,0.1);
    color: #ffffff;
    border: none;
    padding: 6px 14px;
    border-radius: 16px;
    font-size: 13px;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s;
    display: flex;
    align-items: center;
    gap: 5px;
  }}

  .ctrl-btn:hover {{
    background: var(--primary);
    transform: translateY(-1px);
  }}

  .slide-indicator {{
    color: #cbd5e1;
    font-size: 13px;
    font-weight: 600;
    padding: 0 6px;
  }}

  .progress-container {{
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 4px;
    background: rgba(0,0,0,0.06);
    z-index: 10;
  }}

  .progress-bar {{
    height: 100%;
    width: 0%;
    background: linear-gradient(90deg, var(--primary), var(--accent));
    transition: width 0.3s ease;
  }}

  #lightbox {{
    display: none;
    position: fixed;
    top: 0;
    left: 0;
    width: 100vw;
    height: 100vh;
    background: rgba(0,0,0,0.88);
    z-index: 2000;
    align-items: center;
    justify-content: center;
    cursor: zoom-out;
  }}

  #lightbox img {{
    max-width: 92vw;
    max-height: 92vh;
    border-radius: 8px;
    box-shadow: 0 20px 50px rgba(0,0,0,0.8);
  }}

  /* Print / PDF Styles for true 16:9 Widescreen (1920x1080) */
  @page {{
    size: 1920px 1080px;
    margin: 0;
  }}

  @media print {{
    @page {{
      size: 1920px 1080px;
      margin: 0;
    }}
    html, body {{
      width: 1920px !important;
      height: 1080px !important;
      margin: 0 !important;
      padding: 0 !important;
      background: #ffffff !important;
      overflow: visible !important;
      -webkit-print-color-adjust: exact !important;
      print-color-adjust: exact !important;
    }}
    .controls-bar, .progress-container, #lightbox {{
      display: none !important;
    }}
    #presentation-container {{
      width: 1920px !important;
      height: 1080px !important;
      max-width: none !important;
      max-height: none !important;
      box-shadow: none !important;
      border-radius: 0 !important;
      margin: 0 !important;
      padding: 0 !important;
      overflow: visible !important;
    }}
    .slide {{
      display: flex !important;
      position: relative !important;
      page-break-after: always !important;
      page-break-inside: avoid !important;
      break-after: page !important;
      break-inside: avoid !important;
      width: 1920px !important;
      height: 1080px !important;
      min-height: 1080px !important;
      max-height: 1080px !important;
      box-sizing: border-box !important;
      margin: 0 !important;
      padding: 55px 75px 35px 75px !important;
      -webkit-print-color-adjust: exact !important;
      print-color-adjust: exact !important;
    }}
    .slide-title {{
      font-size: 38px !important;
    }}
    .slide-topic {{
      font-size: 16px !important;
    }}
    .slide-subtitle {{
      font-size: 26px !important;
    }}
    .slide-p {{
      font-size: 21px !important;
      line-height: 1.55 !important;
    }}
    .slide-list li {{
      font-size: 21px !important;
      line-height: 1.5 !important;
      padding-left: 28px !important;
    }}
    .slide-list:not(.ol-list) li::before {{
      font-size: 26px !important;
    }}
    .formula-block {{
      font-size: 26px !important;
      padding: 20px 28px !important;
    }}
    .slide-table {{
      font-size: 20px !important;
    }}
    .slide-table th, .slide-table td {{
      padding: 13px 18px !important;
    }}
    .info-badge {{
      font-size: 20px !important;
      padding: 14px 20px !important;
    }}
    .slide-footer {{
      font-size: 16px !important;
      padding-top: 14px !important;
    }}
    .main-title {{
      font-size: 76px !important;
    }}
    .sub-title {{
      font-size: 34px !important;
    }}
    .title-desc {{
      font-size: 22px !important;
    }}
    .project-tag {{
      font-size: 18px !important;
      padding: 8px 24px !important;
    }}
    .meta-card {{
      font-size: 18px !important;
      padding: 14px 24px !important;
    }}
    .meta-label {{
      font-size: 14px !important;
    }}
    .meta-val {{
      font-size: 18px !important;
    }}
    .slide-image-wrapper {{
      max-height: 600px !important;
      padding: 12px !important;
    }}
    .slide-img {{
      max-height: 520px !important;
    }}
    .img-caption {{
      font-size: 18px !important;
      margin-top: 10px !important;
    }}
    /* Print diagram scaling */
    .arch-flow-diagram {{
      padding: 20px 24px !important;
      gap: 12px !important;
    }}
    .arch-card {{
      padding: 14px 18px !important;
    }}
    .arch-card-title {{
      font-size: 18px !important;
    }}
    .arch-badge {{
      font-size: 13px !important;
    }}
    .arch-tech-pill {{
      font-size: 14px !important;
    }}
    .arch-card-tags .tag {{
      font-size: 14px !important;
      padding: 4px 12px !important;
    }}
    .submodule {{
      font-size: 14px !important;
      padding: 6px 10px !important;
    }}
    .arch-arrow-bridge {{
      font-size: 14px !important;
    }}
    .arch-desc {{
      font-size: 14px !important;
    }}
  }}
</style>
</head>
<body>

<div id="presentation-container">
  <div class="progress-container">
    <div class="progress-bar" id="progressBar"></div>
  </div>

  {all_slides_html}
</div>

<!-- Controls Bar -->
<div class="controls-bar">
  <button class="ctrl-btn" onclick="prevSlide()" title="Предыдущий слайд (Стрелка влево)">◀ Назад</button>
  <span class="slide-indicator" id="slideIndicator">1 / {total_slides}</span>
  <button class="ctrl-btn" onclick="nextSlide()" title="Следующий слайд (Стрелка вправо / Пробел)">Вперёд ▶</button>
  <button class="ctrl-btn" onclick="toggleFullScreen()" title="Полноэкранный режим (Клавиша F)">⛶ Полный экран</button>
  <button class="ctrl-btn" onclick="window.print()" title="Печать или экспорт в PDF">🖨 Печать в PDF</button>
</div>

<!-- Lightbox for enlarged images -->
<div id="lightbox" onclick="closeLightbox()">
  <img id="lightboxImg" src="" alt="Enlarged" />
</div>

<script>
  let currentSlide = 1;
  const totalSlides = {total_slides};

  function showSlide(index) {{
    if (index < 1) index = 1;
    if (index > totalSlides) index = totalSlides;
    currentSlide = index;

    document.querySelectorAll('.slide').forEach(s => s.classList.remove('active'));
    const target = document.getElementById(`slide-${{currentSlide}}`);
    if (target) {{
      target.classList.add('active');
    }}

    document.getElementById('slideIndicator').innerText = `${{currentSlide}} / ${{totalSlides}}`;
    const progressPercent = ((currentSlide - 1) / (totalSlides - 1 || 1)) * 100;
    document.getElementById('progressBar').style.width = `${{progressPercent}}%`;
  }}

  function nextSlide() {{
    if (currentSlide < totalSlides) {{
      showSlide(currentSlide + 1);
    }}
  }}

  function prevSlide() {{
    if (currentSlide > 1) {{
      showSlide(currentSlide - 1);
    }}
  }}

  function toggleFullScreen() {{
    if (!document.fullscreenElement) {{
      document.documentElement.requestFullscreen().catch(err => console.log(err));
    }} else {{
      document.exitFullscreen().catch(err => console.log(err));
    }}
  }}

  function openLightbox(src) {{
    const lb = document.getElementById('lightbox');
    const lbImg = document.getElementById('lightboxImg');
    lbImg.src = src;
    lb.style.display = 'flex';
  }}

  function closeLightbox() {{
    document.getElementById('lightbox').style.display = 'none';
  }}

  // Keyboard navigation
  window.addEventListener('keydown', (e) => {{
    if (e.key === 'ArrowRight' || e.key === 'PageDown' || e.key === ' ') {{
      e.preventDefault();
      nextSlide();
    }} else if (e.key === 'ArrowLeft' || e.key === 'PageUp' || e.key === 'Backspace') {{
      e.preventDefault();
      prevSlide();
    }} else if (e.key === 'f' || e.key === 'F' || e.key === 'а' || e.key === 'А') {{
      toggleFullScreen();
    }} else if (e.key === 'Home') {{
      showSlide(1);
    }} else if (e.key === 'End') {{
      showSlide(totalSlides);
    }} else if (e.key === 'Escape') {{
      closeLightbox();
    }}
  }});

  // Touch swipe support
  let touchStartX = 0;
  window.addEventListener('touchstart', (e) => {{
    touchStartX = e.changedTouches[0].screenX;
  }}, false);

  window.addEventListener('touchend', (e) => {{
    let touchEndX = e.changedTouches[0].screenX;
    if (touchEndX < touchStartX - 50) {{
      nextSlide();
    }} else if (touchEndX > touchStartX + 50) {{
      prevSlide();
    }}
  }}, false);

  // Initialize
  showSlide(1);
</script>

</body>
</html>
"""

    with open(HTML_OUT, "w", encoding="utf-8") as f:
        f.write(html_template)

    print(f"Interactive HTML presentation generated: {HTML_OUT}")

    # Generate PDF via Headless Edge/Chrome
    edge_exe = r"C:\Program Files (x86)\Microsoft\Edge\Application\msedge.exe"
    chrome_exe = r"C:\Program Files\Google\Chrome\Application\chrome.exe"
    browser_exe = edge_exe if os.path.exists(edge_exe) else chrome_exe

    if os.path.exists(browser_exe):
        print(f"Generating PDF via {browser_exe}...")
        cmd = [
            browser_exe,
            "--headless",
            "--disable-gpu",
            "--window-size=1920,1080",
            "--run-all-compositor-stages-before-draw",
            f"--print-to-pdf={PDF_OUT}",
            "--no-pdf-header-footer",
            str(HTML_OUT)
        ]
        res = subprocess.run(cmd, capture_output=True, text=True)
        if PDF_OUT.exists():
            size = PDF_OUT.stat().st_size
            print(f"PDF presentation generated: {PDF_OUT} ({size} bytes)")
        else:
            print(f"Notice: PDF generation finished. Code: {res.returncode}")
    else:
        print("Note: Chromium/Edge not found for automated PDF generation.")

if __name__ == "__main__":
    generate_presentation()
