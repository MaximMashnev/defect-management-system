import os
import re
import subprocess
import markdown

def clean_math(text):
    # Ensure blank line before and after tables
    text = re.sub(r'([^\n|])\n(\|)', r'\1\n\n\2', text)
    text = re.sub(r'(\|[^\n]*)\n([^\|\n\s])', r'\1\n\n\2', text)

    # Replace block formulas (1) and (2) with elegant HTML blocks
    formula_1 = (
        '<div class="formula-container">'
        '<span class="formula-text"><i>Priority</i> = max(1, min(10, <i>S</i><sub>общ</sub>))</span>'
        '<span class="formula-number">(1)</span>'
        '</div>'
    )
    formula_2 = (
        '<div class="formula-container">'
        '<span class="formula-text"><i>S</i><sub>общ</sub> = <i>S</i><sub>баз</sub> + <i>S</i><sub>кат</sub> + <i>S</i><sub>опис</sub> + <i>S</i><sub>медиа</sub> + <i>S</i><sub>время</sub> + <i>S</i><sub>повтор</sub> + <i>S</i><sub>объект</sub></span>'
        '<span class="formula-number">(2)</span>'
        '</div>'
    )

    # Match block formula 1
    text = re.sub(r'\$\$[\s\S]*?Priority[\s\S]*?\(1\)\$\$', '\n\n' + formula_1 + '\n\n', text)
    # Match block formula 2
    text = re.sub(r'\$\$[\s\S]*?S_[\s\S]*?\(2\)\$\$', '\n\n' + formula_2 + '\n\n', text)

    # Inline variable notations
    text = re.sub(r'\(\$S_\{?\\text\{баз\}?\}?\s*\\in\s*\[1\.\.4\]\$\)', '(<i>S</i><sub>баз</sub> ∈ [1..4])', text)
    text = re.sub(r'\$S_\{?\\text\{([^\}]+)\}?\}\$', r'<i>S</i><sub>\1</sub>', text)

    # Clean any remaining inline symbols or LaTeX leftovers
    text = text.replace(r"$\ge 50$", "≥ 50")
    text = text.replace(r"$\ge 200$", "≥ 200")
    text = text.replace(r"$\ge 3$", "≥ 3")
    text = text.replace(r"$[1; 10]$", "[1; 10]")
    text = text.replace(r"$1 : N$", "1 : N")
    text = text.replace(r"\ge", "≥")
    text = text.replace(r"\le", "≤")
    text = text.replace(r"\to", "→")
    text = text.replace(r"\in", "∈")
    text = text.replace(r"\forall", "∀")
    text = text.replace(r"\sum", "∑")
    text = text.replace(r"\neq", "≠")
    text = text.replace(r"\implies", "⇒")
    text = text.replace(r"\land", "и")
    text = text.replace(r"\qquad", " ")
    text = text.replace(r"\;", " ")
    text = text.replace(r"\,", " ")
    text = text.replace(r"$$", "")
    text = text.replace(r"$", "")
    return text

def generate_pdf():
    md_path = r"e:\projects\defect-management-system\TZ.md"
    html_path = r"e:\projects\defect-management-system\tz_export.html"
    pdf_path = r"e:\projects\defect-management-system\TZ.pdf"

    with open(md_path, "r", encoding="utf-8") as f:
        md_content = f.read()

    md_content = clean_math(md_content)

    # Convert markdown to html with tables and fenced_code
    html_body = markdown.markdown(
        md_content,
        extensions=["tables", "fenced_code", "sane_lists"]
    )

    full_html = f"""<!DOCTYPE html>
<html lang="ru">
<head>
<meta charset="UTF-8">
<title>Техническое задание — DefectTrack</title>
<style>
  @page {{
    size: A4;
    margin: 20mm 15mm 20mm 25mm;
  }}

  * {{
    box-sizing: border-box;
  }}

  body {{
    font-family: 'Times New Roman', 'Liberation Serif', serif;
    font-size: 12pt;
    line-height: 1.45;
    color: #111827;
    background-color: #ffffff;
    margin: 0;
    padding: 0;
    text-rendering: optimizeLegibility;
  }}

  /* Headings */
  h1, h2, h3, h4, h5, h6 {{
    font-family: 'Times New Roman', 'Liberation Serif', serif;
    font-weight: bold;
    color: #0f2442;
    page-break-after: avoid;
    break-after: avoid;
  }}

  h1 {{
    font-size: 16pt;
    text-align: center;
    margin-top: 18pt;
    margin-bottom: 12pt;
    text-transform: uppercase;
  }}

  h2 {{
    font-size: 13.5pt;
    color: #1a365d;
    border-bottom: 1.5px solid #2b4c7e;
    padding-bottom: 3pt;
    margin-top: 18pt;
    margin-bottom: 8pt;
    page-break-after: avoid;
    break-after: avoid;
  }}

  h3 {{
    font-size: 12.5pt;
    color: #1e3a8a;
    margin-top: 13pt;
    margin-bottom: 6pt;
    page-break-after: avoid;
    break-after: avoid;
  }}

  p {{
    margin-top: 0;
    margin-bottom: 7pt;
    text-align: justify;
    text-indent: 1.25cm;
  }}

  p:empty {{
    display: none;
  }}

  /* Lists */
  ul, ol {{
    margin-top: 3pt;
    margin-bottom: 7pt;
    padding-left: 1.8cm;
  }}

  li {{
    margin-bottom: 3pt;
    text-align: justify;
  }}

  /* Table Caption */
  p:has(> em:first-child) {{
    text-indent: 0 !important;
    margin-top: 12pt;
    margin-bottom: 4pt;
    font-weight: bold;
    color: #1e293b;
    page-break-after: avoid;
    break-after: avoid;
  }}

  p > em {{
    font-style: italic;
  }}

  /* Real Tables Styling (OpticalWatch standard) */
  table {{
    width: 100%;
    border-collapse: collapse;
    margin-top: 4pt;
    margin-bottom: 14pt;
    font-size: 10pt;
    line-height: 1.35;
    page-break-inside: auto;
  }}

  thead {{
    display: table-header-group;
  }}

  tr {{
    page-break-inside: avoid;
    break-inside: avoid;
  }}

  th, td {{
    border: 1px solid #718096;
    padding: 5pt 7pt;
    vertical-align: middle;
    text-align: left;
  }}

  th {{
    background-color: #1f4e78;
    color: #ffffff;
    font-weight: bold;
    text-align: center;
    border: 1px solid #1f4e78;
  }}

  tbody tr:nth-child(even) {{
    background-color: #f7fafc;
  }}

  tbody tr:nth-child(odd) {{
    background-color: #ffffff;
  }}

  /* Specific adjustment for wide tables (e.g. Gantt chart) */
  table:has(th:nth-child(10)) {{
    font-size: 8pt;
  }}

  table:has(th:nth-child(10)) th,
  table:has(th:nth-child(10)) td {{
    padding: 3pt 2pt;
    text-align: center;
  }}

  table:has(th:nth-child(10)) td:first-child {{
    text-align: left;
    padding-left: 4pt;
    font-size: 8.5pt;
  }}

  /* Code and ASCII diagrams */
  pre {{
    font-family: 'Consolas', 'Courier New', monospace;
    font-size: 8.5pt;
    line-height: 1.25;
    background-color: #f8fafc;
    border: 1px solid #cbd5e1;
    border-radius: 4px;
    padding: 7pt 9pt;
    overflow-x: auto;
    page-break-inside: avoid;
    break-inside: avoid;
    margin-top: 6pt;
    margin-bottom: 10pt;
    white-space: pre;
  }}

  code {{
    font-family: 'Consolas', 'Courier New', monospace;
    font-size: 9.5pt;
    background-color: #f1f5f9;
    padding: 1pt 3pt;
    border-radius: 3px;
  }}

  pre code {{
    background-color: transparent;
    padding: 0;
    font-size: 8.5pt;
  }}

  /* Formula Container - GOST Standard Display */
  .formula-container {{
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin: 14pt 0 14pt 1.25cm;
    padding: 10pt 20pt;
    background-color: #f8fafc;
    border: 1px solid #cbd5e1;
    border-left: 5px solid #1f4e78;
    border-radius: 4px;
    page-break-inside: avoid;
    break-inside: avoid;
  }}

  .formula-text {{
    font-family: 'Times New Roman', 'Liberation Serif', serif;
    font-size: 13.5pt;
    font-weight: 500;
    color: #0f172a;
    text-align: center;
    flex-grow: 1;
  }}

  .formula-number {{
    font-family: 'Times New Roman', 'Liberation Serif', serif;
    font-size: 12.5pt;
    font-weight: normal;
    color: #475569;
    padding-left: 18pt;
    white-space: nowrap;
  }}

  /* Page breaks */
  div[style*="page-break-after: always"] {{
    page-break-after: always;
    break-after: page;
    height: 0;
  }}

  /* Title block table override (borderless) */
  table[style*="border: none"] {{
    margin: 0 !important;
  }}

  table[style*="border: none"] th,
  table[style*="border: none"] td {{
    border: none !important;
    background: transparent !important;
    padding: 0 !important;
    color: inherit !important;
  }}

  hr {{
    border: 0;
    height: 1px;
    background: #cbd5e0;
    margin: 12pt 0;
  }}

  strong {{
    font-weight: bold;
  }}
</style>
</head>
<body>
{html_body}
</body>
</html>
"""

    with open(html_path, "w", encoding="utf-8") as f:
        f.write(full_html)

    print("HTML updated successfully.")

    edge_exe = r"C:\Program Files (x86)\Microsoft\Edge\Application\msedge.exe"
    chrome_exe = r"C:\Program Files\Google\Chrome\Application\chrome.exe"
    browser_exe = edge_exe if os.path.exists(edge_exe) else chrome_exe

    cmd = [
        browser_exe,
        "--headless",
        "--disable-gpu",
        "--run-all-compositor-stages-before-draw",
        f'--print-to-pdf={pdf_path}',
        "--no-pdf-header-footer",
        html_path
    ]

    res = subprocess.run(cmd, capture_output=True, text=True)
    if os.path.exists(pdf_path):
        size = os.path.getsize(pdf_path)
        print(f"PDF successfully generated: {pdf_path} (size: {size} bytes)")
    else:
        print("PDF was not created. Error:", res.stderr)

if __name__ == "__main__":
    generate_pdf()
