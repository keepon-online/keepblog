/**
 * 纯前端现代极客风技术博文卡片封面生成器 (Canvas 2D)
 * 16:9 比例 (1280 × 720 高清输出)
 */

export interface TechCoverTheme {
  id: string;
  name: string;
  desc: string;
  bgStart: string;
  bgEnd: string;
  accent: string;
  secondary: string;
  textColor: string;
  tagBg: string;
  tagBorder: string;
  tagText: string;
}

export const TECH_COVER_THEMES: TechCoverTheme[] = [
  {
    id: "galaxy",
    name: "深邃银河 (Galaxy)",
    desc: "深空蓝紫微光与科技网格",
    bgStart: "#070B19",
    bgEnd: "#141A36",
    accent: "#6366F1",
    secondary: "#38BDF8",
    textColor: "#FFFFFF",
    tagBg: "rgba(99, 102, 241, 0.2)",
    tagBorder: "rgba(129, 140, 248, 0.4)",
    tagText: "#C7D2FE"
  },
  {
    id: "cyber",
    name: "赛博霓虹 (Cyber)",
    desc: "深黑炭灰与青橙霓虹光影",
    bgStart: "#0B0D14",
    bgEnd: "#161926",
    accent: "#00F5FF",
    secondary: "#FF6B6B",
    textColor: "#FFFFFF",
    tagBg: "rgba(0, 245, 255, 0.15)",
    tagBorder: "rgba(0, 245, 255, 0.35)",
    tagText: "#7DD3FC"
  },
  {
    id: "emerald",
    name: "极客薄荷 (Emerald)",
    desc: "极简墨绿与清新终端质感",
    bgStart: "#051A14",
    bgEnd: "#0D2E25",
    accent: "#10B981",
    secondary: "#34D399",
    textColor: "#FFFFFF",
    tagBg: "rgba(16, 185, 129, 0.2)",
    tagBorder: "rgba(52, 211, 153, 0.4)",
    tagText: "#A7F3D0"
  },
  {
    id: "sunset",
    name: "暮光极光 (Sunset)",
    desc: "暮紫与暖橙柔和弥散",
    bgStart: "#1A0E2E",
    bgEnd: "#3B185F",
    accent: "#F97316",
    secondary: "#E879F9",
    textColor: "#FFFFFF",
    tagBg: "rgba(249, 115, 22, 0.2)",
    tagBorder: "rgba(251, 146, 60, 0.4)",
    tagText: "#FED7AA"
  },
  {
    id: "slate",
    name: "冷淡黑金 (Obsidian)",
    desc: "工业极简灰度与科技准星",
    bgStart: "#11141A",
    bgEnd: "#1E232E",
    accent: "#F59E0B",
    secondary: "#94A3B8",
    textColor: "#FFFFFF",
    tagBg: "rgba(245, 158, 11, 0.15)",
    tagBorder: "rgba(245, 158, 11, 0.35)",
    tagText: "#FDE68A"
  }
];

export interface TechCoverOptions {
  title: string;
  themeId?: string;
  tags?: string[];
  author?: string;
  siteName?: string;
  width?: number;
  height?: number;
}

/**
 * 动态渲染现代极客风技术卡片封面
 */
export function renderTechCover(
  canvas: HTMLCanvasElement,
  options: TechCoverOptions
) {
  const width = options.width || 1280;
  const height = options.height || 720;
  canvas.width = width;
  canvas.height = height;

  const ctx = canvas.getContext("2d");
  if (!ctx) return;

  const theme =
    TECH_COVER_THEMES.find(t => t.id === options.themeId) ||
    TECH_COVER_THEMES[0];

  // 1. 绘制背景线性渐变
  const bgGrad = ctx.createLinearGradient(0, 0, width, height);
  bgGrad.addColorStop(0, theme.bgStart);
  bgGrad.addColorStop(1, theme.bgEnd);
  ctx.fillStyle = bgGrad;
  ctx.fillRect(0, 0, width, height);

  // 2. 绘制科技风光晕与弥散粒子 (Glow spheres)
  ctx.save();
  // 右上角弥散光球
  const glowTopRight = ctx.createRadialGradient(
    width * 0.85,
    height * 0.15,
    20,
    width * 0.85,
    height * 0.15,
    width * 0.45
  );
  glowTopRight.addColorStop(0, hexToRgba(theme.accent, 0.45));
  glowTopRight.addColorStop(0.5, hexToRgba(theme.secondary, 0.18));
  glowTopRight.addColorStop(1, "transparent");
  ctx.fillStyle = glowTopRight;
  ctx.fillRect(0, 0, width, height);

  // 左下角微光球
  const glowBottomLeft = ctx.createRadialGradient(
    width * 0.15,
    height * 0.85,
    10,
    width * 0.15,
    height * 0.85,
    width * 0.35
  );
  glowBottomLeft.addColorStop(0, hexToRgba(theme.secondary, 0.3));
  glowBottomLeft.addColorStop(1, "transparent");
  ctx.fillStyle = glowBottomLeft;
  ctx.fillRect(0, 0, width, height);
  ctx.restore();

  // 3. 绘制细致的技术几何背景网格
  ctx.save();
  ctx.strokeStyle = "rgba(255, 255, 255, 0.04)";
  ctx.lineWidth = 1;
  const gridSize = 40;
  for (let x = 0; x < width; x += gridSize) {
    ctx.beginPath();
    ctx.moveTo(x, 0);
    ctx.lineTo(x, height);
    ctx.stroke();
  }
  for (let y = 0; y < height; y += gridSize) {
    ctx.beginPath();
    ctx.moveTo(0, y);
    ctx.lineTo(width, y);
    ctx.stroke();
  }
  ctx.restore();

  // 4. 绘制四周极客科技装饰框与交叉十字准星
  ctx.save();
  ctx.strokeStyle = hexToRgba(theme.accent, 0.25);
  ctx.lineWidth = 1.5;
  const margin = 48;
  const cornerLen = 24;

  // 四角 L 形装饰
  // 左上
  ctx.beginPath();
  ctx.moveTo(margin, margin + cornerLen);
  ctx.lineTo(margin, margin);
  ctx.lineTo(margin + cornerLen, margin);
  ctx.stroke();
  // 右上
  ctx.beginPath();
  ctx.moveTo(width - margin - cornerLen, margin);
  ctx.lineTo(width - margin, margin);
  ctx.lineTo(width - margin, margin + cornerLen);
  ctx.stroke();
  // 左下
  ctx.beginPath();
  ctx.moveTo(margin, height - margin - cornerLen);
  ctx.lineTo(margin, height - margin);
  ctx.lineTo(margin + cornerLen, height - margin);
  ctx.stroke();
  // 右下
  ctx.beginPath();
  ctx.moveTo(width - margin - cornerLen, height - margin);
  ctx.lineTo(width - margin, height - margin);
  ctx.lineTo(width - margin, height - margin - cornerLen);
  ctx.stroke();
  ctx.restore();

  // 5. 绘制背景代码水印 / 极客标志符号
  ctx.save();
  ctx.font = "900 120px monospace";
  ctx.fillStyle = "rgba(255, 255, 255, 0.025)";
  ctx.textAlign = "right";
  ctx.fillText("</>", width - 80, height * 0.7);
  ctx.restore();

  // 6. 绘制顶部分类/站点胶囊徽标
  ctx.save();
  const siteBadgeText = (options.siteName || "TECH BLOG").toUpperCase();
  ctx.font = "600 14px -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif";
  const badgeWidth = ctx.measureText(siteBadgeText).width + 36;
  const badgeHeight = 28;
  const badgeX = margin + 20;
  const badgeY = margin + 20;

  // 胶囊背景
  drawRoundedRect(ctx, badgeX, badgeY, badgeWidth, badgeHeight, 14);
  ctx.fillStyle = hexToRgba(theme.accent, 0.15);
  ctx.fill();
  ctx.strokeStyle = hexToRgba(theme.accent, 0.4);
  ctx.lineWidth = 1;
  ctx.stroke();

  // 脉冲小圆点
  ctx.beginPath();
  ctx.arc(badgeX + 14, badgeY + badgeHeight / 2, 4, 0, Math.PI * 2);
  ctx.fillStyle = theme.accent;
  ctx.fill();

  // 徽标文本
  ctx.fillStyle = theme.textColor;
  ctx.fillText(siteBadgeText, badgeX + 26, badgeY + 19);
  ctx.restore();

  // 7. 绘制主标题（自适应折行与字号）
  ctx.save();
  const rawTitle = options.title.trim() || "探索全栈技术与架构之道";
  const maxTitleWidth = width - margin * 2 - 80;

  let fontSize = 54;
  if (rawTitle.length > 36) fontSize = 42;
  else if (rawTitle.length > 24) fontSize = 48;

  ctx.font = `800 ${fontSize}px -apple-system, BlinkMacSystemFont, 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei', sans-serif`;
  ctx.fillStyle = theme.textColor;
  ctx.shadowColor = "rgba(0, 0, 0, 0.4)";
  ctx.shadowBlur = 12;
  ctx.shadowOffsetY = 4;

  const lines = wrapText(ctx, rawTitle, maxTitleWidth, 3);
  const lineHeight = fontSize * 1.35;
  const titleStartY = height * 0.38 - ((lines.length - 1) * lineHeight) / 2;

  lines.forEach((line, index) => {
    ctx.fillText(line, margin + 20, titleStartY + index * lineHeight);
  });
  ctx.restore();

  // 8. 绘制标签药丸 (Tags Chips)
  const tags = (options.tags && options.tags.length ? options.tags : ["Engineering", "Architecture"]).slice(0, 5);
  ctx.save();
  let tagCurrentX = margin + 20;
  const tagY = height * 0.68;
  const tagHeight = 32;

  ctx.font = "600 15px -apple-system, BlinkMacSystemFont, 'PingFang SC', sans-serif";

  tags.forEach(t => {
    const textWidth = ctx.measureText(t).width;
    const tagWidth = textWidth + 24;

    if (tagCurrentX + tagWidth < width - margin - 40) {
      drawRoundedRect(ctx, tagCurrentX, tagY, tagWidth, tagHeight, 6);
      ctx.fillStyle = theme.tagBg;
      ctx.fill();
      ctx.strokeStyle = theme.tagBorder;
      ctx.lineWidth = 1;
      ctx.stroke();

      ctx.fillStyle = theme.tagText;
      ctx.fillText(t, tagCurrentX + 12, tagY + 21);

      tagCurrentX += tagWidth + 12;
    }
  });
  ctx.restore();

  // 9. 绘制底部作者与版权署名 (Footer)
  ctx.save();
  const footerY = height - margin - 16;
  const authorName = options.author || "KeepBlog Author";

  ctx.strokeStyle = "rgba(255, 255, 255, 0.08)";
  ctx.lineWidth = 1;
  ctx.beginPath();
  ctx.moveTo(margin + 20, footerY - 24);
  ctx.lineTo(width - margin - 20, footerY - 24);
  ctx.stroke();

  // 作者名字
  ctx.font = "500 16px -apple-system, BlinkMacSystemFont, 'PingFang SC', sans-serif";
  ctx.fillStyle = "rgba(255, 255, 255, 0.7)";
  ctx.fillText(`By ${authorName}`, margin + 20, footerY);

  // 右侧日期或微印记
  ctx.font = "500 14px monospace";
  ctx.fillStyle = "rgba(255, 255, 255, 0.4)";
  ctx.textAlign = "right";
  ctx.fillText(new Date().toISOString().slice(0, 10), width - margin - 20, footerY);
  ctx.restore();
}

/**
 * 辅助函数：绘制圆角矩形
 */
function drawRoundedRect(
  ctx: CanvasRenderingContext2D,
  x: number,
  y: number,
  width: number,
  height: number,
  radius: number
) {
  ctx.beginPath();
  ctx.moveTo(x + radius, y);
  ctx.lineTo(x + width - radius, y);
  ctx.quadraticCurveTo(x + width, y, x + width, y + radius);
  ctx.lineTo(x + width, y + height - radius);
  ctx.quadraticCurveTo(x + width, y + height, x + width - radius, y + height);
  ctx.lineTo(x + radius, y + height);
  ctx.quadraticCurveTo(x, y + height, x, y + height - radius);
  ctx.lineTo(x, y + radius);
  ctx.quadraticCurveTo(x, y, x + radius, y);
  ctx.closePath();
}

/**
 * 辅助函数：文字按最大宽度自动折行
 */
function wrapText(
  ctx: CanvasRenderingContext2D,
  text: string,
  maxWidth: number,
  maxLines = 3
): string[] {
  const words = text.split("");
  const lines: string[] = [];
  let currentLine = "";

  for (let i = 0; i < words.length; i++) {
    const testLine = currentLine + words[i];
    const metrics = ctx.measureText(testLine);
    if (metrics.width > maxWidth && i > 0) {
      lines.push(currentLine);
      currentLine = words[i];
      if (lines.length === maxLines - 1) {
        // 最后一行截断并加省略号
        const remaining = words.slice(i).join("");
        let fit = "";
        for (let j = 0; j < remaining.length; j++) {
          if (ctx.measureText(fit + remaining[j] + "…").width > maxWidth) {
            break;
          }
          fit += remaining[j];
        }
        lines.push(fit + "…");
        return lines;
      }
    } else {
      currentLine = testLine;
    }
  }
  if (currentLine) {
    lines.push(currentLine);
  }
  return lines;
}

/**
 * 辅助函数：十六进制颜色转 RGBA
 */
function hexToRgba(hex: string, alpha: number): string {
  let c = hex.replace("#", "");
  if (c.length === 3) {
    c = c
      .split("")
      .map(x => x + x)
      .join("");
  }
  const num = parseInt(c, 16);
  const r = (num >> 16) & 255;
  const g = (num >> 8) & 255;
  const b = num & 255;
  return `rgba(${r}, ${g}, ${b}, ${alpha})`;
}

/**
 * 将 Canvas 输出为 Blob 对象（用于上传至服务器）
 */
export function canvasToBlob(
  canvas: HTMLCanvasElement,
  type = "image/png",
  quality = 0.95
): Promise<Blob> {
  return new Promise((resolve, reject) => {
    canvas.toBlob(
      blob => {
        if (blob) resolve(blob);
        else reject(new Error("Canvas toBlob 导出失败"));
      },
      type,
      quality
    );
  });
}

/**
 * 触发本地下载
 */
export function downloadCanvas(
  canvas: HTMLCanvasElement,
  filename = "cover.png"
) {
  const url = canvas.toDataURL("image/png");
  const a = document.createElement("a");
  a.href = url;
  a.download = filename;
  document.body.appendChild(a);
  a.click();
  document.body.removeChild(a);
}
