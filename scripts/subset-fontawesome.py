#!/usr/bin/env python3
"""Font Awesome 子集化脚本。

用法:在仓库根目录执行 `python3 scripts/subset-fontawesome.py`。
前置:源包放在 static/plugins/@fortawesome/fontawesome-free(全量包,生成后可删除,
本仓库保留此脚本与新图标引入时的操作说明)。

新增图标的流程:
1. 在模板/JS 中正常使用 fa- 类名;
2. 把新图标名加进下方 solid/regular/brands 列表;
3. 临时放回全量包目录,执行本脚本重新生成 static/plugins/fontawesome/。

依赖:pip install fonttools brotli
"""
import re
import pathlib
from fontTools.subset import Subsetter, Options, load_font, save_font

SRC = pathlib.Path("static/plugins/@fortawesome/fontawesome-free")  # 全量包(临时)
OUT = pathlib.Path("static/plugins/fontawesome")

# 全站实际使用的图标，按风格分组。手工维护，但 main() 会自动扫描
# internal/web、static/js、pkg、api 中的 fa- 类名做交叉校验——分页条等
# 服务端拼 HTML 的代码在 Go 文件里，只扫模板会漏（chevron 就这样漏过）。
solid = """adjust angle-double-down angle-down archive arrow-up arrows-alt-h bars
bullhorn chart-line check chevron-left chevron-right cog comments envelope
folder-open heart history home inbox link paste search sign-out-alt spinner
stream tag tags thumbs-up thumbtack times""".split()
regular = "calendar-alt clock eye file-word".split()
brands = "git-alt github".split()


def scan_used_icons() -> set[str]:
    """扫描全站源码里的 fa-xxx 图标类名（模板/JS/Go 拼串）。"""
    import re as _re
    roots = ["internal/web", "static/js", "pkg", "api"]
    used: set[str] = set()
    for root in roots:
        for p in pathlib.Path(root).rglob("*"):
            if p.suffix not in (".html", ".js", ".go") or not p.is_file():
                continue
            for m in _re.finditer(r"\bfa[bsr]? fa-([a-z0-9-]+)", p.read_text(errors="ignore")):
                used.add(m.group(1))
    return used


def name_to_codepoint(css: str) -> dict[str, str]:
    mapping: dict[str, str] = {}
    for m in re.finditer(r'([^{}]+)\{content:"\\([0-9a-f]{4})"\}', css):
        for sel in m.group(1).split(","):
            n = re.match(r"\.fa-([a-z0-9-]+):before$", sel.strip())
            if n:
                mapping.setdefault(n.group(1), m.group(2))
    return mapping


def subset_font(ttf: str, names: list[str], name2cp: dict[str, str]) -> None:
    cps = sorted({int(name2cp[n], 16) for n in names})
    opts = Options()
    opts.flavor = "woff2"
    font = load_font(str(SRC / "webfonts" / ttf), opts)
    cmap = font.getBestCmap()
    missing = [hex(c) for c in cps if c not in cmap]
    if missing:
        raise SystemExit(f"{ttf} 缺少码点 {missing}:CSS 与字体版本不匹配")
    ss = Subsetter(opts)
    ss.populate(unicodes=cps)
    ss.subset(font)
    dst = OUT / "webfonts" / ttf.replace(".ttf", ".woff2")
    save_font(font, str(dst), opts)
    print(f"  {dst.name}: {len(cps)} 字形, {dst.stat().st_size / 1024:.1f} KB")


def extract_balanced(text: str, start: int) -> str:
    i = text.index("{", start)
    depth, j = 0, i
    while j < len(text):
        if text[j] == "{":
            depth += 1
        elif text[j] == "}":
            depth -= 1
            if depth == 0:
                break
        j += 1
    return text[start:j + 1]


def main() -> None:
    css = (SRC / "css/all.min.css").read_text()
    name2cp = name_to_codepoint(css)
    all_icons = set(solid + regular + brands)
    missing = [n for n in all_icons if n not in name2cp]
    if missing:
        raise SystemExit(f"CSS 中找不到码点: {missing}")

    # 交叉校验：源码扫描到的图标若不在子集列表里，说明列表漏维护
    #（风格归属按 all.min.css 的码点存在性归组，此处仅提示名字）。
    used = scan_used_icons()
    unlisted = sorted(used - all_icons)
    if unlisted:
        raise SystemExit(
            f"源码中使用了但子集列表缺少的图标: {unlisted}，"
            f"请加进 subset-fontawesome.py 对应风格分组后重跑"
        )

    (OUT / "webfonts").mkdir(parents=True, exist_ok=True)
    (OUT / "css").mkdir(parents=True, exist_ok=True)
    print("生成子集字体:")
    subset_font("fa-solid-900.ttf", solid, name2cp)
    subset_font("fa-regular-400.ttf", regular, name2cp)
    subset_font("fa-brands-400.ttf", brands, name2cp)

    # 提取基础类、图标规则与用到的动画工具类
    base_first = {".fas", ".far", ".fab", ".fa-solid", ".fa-regular",
                  ".fa-brands", ".fa-fw", ".fa-spin", ".fa-pulse", ".fa-shake"}
    picked, seen = [], set()
    for sel, body in re.findall(r"([^{}]+)\{([^{}]*)\}", css):
        sels = [x.strip() for x in sel.strip().split(",")]
        hit = any(s in base_first for s in sels) or any(
            (m := re.match(r"\.fa-([a-z0-9-]+):before$", s)) and m.group(1) in all_icons
            for s in sels
        )
        if hit and (key := sel.strip() + body[:40]) not in seen:
            seen.add(key)
            picked.append((sel.strip(), body))

    # 兜底确保核心字体族与字重声明完整
    core_rules = [
        (".fa-brands,.fab", 'font-family:"Font Awesome 6 Brands";font-weight:400'),
        (".fa-regular,.far", 'font-family:"Font Awesome 6 Free";font-weight:400'),
        (".fa-solid,.fas", 'font-family:"Font Awesome 6 Free";font-weight:900'),
    ]
    for sel, body in core_rules:
        if not any(p[0] == sel and "font-family" in p[1] for p in picked):
            picked.append((sel, body))

    keyframes = []
    for pat in (r"@-webkit-keyframes fa-(spin|pulse|shake)\b", r"@keyframes fa-(spin|pulse|shake)\b"):
        for m in re.finditer(pat, css):
            keyframes.append(extract_balanced(css, m.start()))

    fontfaces = "\n".join([
        '@font-face{font-family:"Font Awesome 6 Free";font-style:normal;font-weight:900;'
        'font-display:block;src:url(../webfonts/fa-solid-900.woff2) format("woff2")}',
        '@font-face{font-family:"Font Awesome 6 Free";font-style:normal;font-weight:400;'
        'font-display:block;src:url(../webfonts/fa-regular-400.woff2) format("woff2")}',
        '@font-face{font-family:"Font Awesome 6 Brands";font-style:normal;font-weight:400;'
        'font-display:block;src:url(../webfonts/fa-brands-400.woff2) format("woff2")}',
    ])
    header = """/* Font Awesome 6 子集：仅包含 keepblog 模板与脚本实际使用的图标。
 * 模板中的 fa- 类名保持不变；新增图标时用 scripts/subset-fontawesome.py 重新生成。 */
"""
    body = "\n".join(f"{s}{{{b}}}" for s, b in picked) + "\n" + "\n".join(keyframes) + "\n"
    (OUT / "css/fontawesome.min.css").write_text(header + fontfaces + "\n" + body)
    print(f"  fontawesome.min.css: {(OUT / 'css/fontawesome.min.css').stat().st_size / 1024:.1f} KB")


if __name__ == "__main__":
    main()
