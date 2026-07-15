export const baseUrlApi = (url: string) => `http://localhost:8000${url}`;

export const Int2Ip = num => {
  let str;
  const tt = [];
  tt[0] = (num >>> 24) >>> 0;
  tt[1] = ((num << 8) >>> 24) >>> 0;
  tt[2] = (num << 16) >>> 24;
  tt[3] = (num << 24) >>> 24;
  str =
    String(tt[0]) +
    "." +
    String(tt[1]) +
    "." +
    String(tt[2]) +
    "." +
    String(tt[3]);
  return str;
};

export const Ip2Int = ip => {
  let num = 0;
  ip = ip.split(".");
  num =
    Number(ip[0]) * 256 * 256 * 256 +
    Number(ip[1]) * 256 * 256 +
    Number(ip[2]) * 256 +
    Number(ip[3]);
  num = num >>> 0;
  return num;
};

export function computeSize(size: number): string {
  const num = 1024.0;
  if (size < num) return size + " B";
  if (size < Math.pow(num, 2))
    return formattedNumber((size / num).toFixed(2)) + " KB";
  if (size < Math.pow(num, 3))
    return formattedNumber((size / Math.pow(num, 2)).toFixed(2)) + " MB";
  if (size < Math.pow(num, 4))
    return formattedNumber((size / Math.pow(num, 3)).toFixed(2)) + " GB";
  return formattedNumber((size / Math.pow(num, 4)).toFixed(2)) + " TB";
}

export function formattedNumber(num: string) {
  return num.endsWith(".00") ? Number(num.slice(0, -3)) : Number(num);
}
