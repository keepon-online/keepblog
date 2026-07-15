// API 通用响应类型定义

/** 基础响应结构 */
export interface ApiResponse<T = unknown> {
  code: number;
  message: string;
  payload: T;
}

/** 分页响应数据 */
export interface PaginatedData<T> {
  list: T[];
  total: number;
  pageNum: number;
  pageSize: number;
}

/** 分页请求参数 */
export interface PaginationParams {
  pageNum?: number;
  pageSize?: number;
}

// ============ 文章模块 ============

/** 文章信息 */
export interface Post {
  postId: number;
  title: string;
  postSlug: string;
  summary?: string;
  postContent?: string;
  postContentHtml?: string;
  coverImage?: string;
  categoryId: number;
  categoryName?: string;
  tags?: string[];
  status: number; // 0: 草稿, 1: 已发布
  type: number; // 0: 转载, 1: 原创
  readCount?: number;
  wordCount?: number;
  createAt?: string;
  updateAt?: string;
}

/** 文章列表查询参数 */
export interface PostQueryParams extends PaginationParams {
  title?: string;
  categoryId?: number;
  published?: number;
}

// ============ 分类模块 ============

/** 分类信息 */
export interface Category {
  categoryId: number;
  categoryName: string;
  categorySlug?: string;
  description?: string;
  postCount?: number;
  createAt?: string;
}

// ============ 标签模块 ============

/** 标签信息 */
export interface Tag {
  tagId: number;
  tagName: string;
  tagSlug?: string;
  postCount?: number;
  createAt?: string;
}

// ============ 仪表盘模块 ============

/** 仪表盘面板数据 */
export interface DashboardPanel {
  postTotal: number;
  categoryTotal: number;
  tagTotal: number;
  visit: number;
  totalWords?: number;
  totalReadCount?: number;
  todayVisit?: number;
  totalMusic?: number;
}

/** 图表数据项 */
export interface ChartDataItem {
  name: string;
  value: number;
}

/** 折线图数据项 */
export interface LineChartItem {
  name: string;
  pv: number;
  uv: number;
}

/** 仪表盘完整数据 */
export interface DashboardData {
  panel: DashboardPanel;
  line: LineChartItem[];
  pie: ChartDataItem[];
  bar: ChartDataItem[];
  map?: ChartDataItem[];
}

// ============ 用户模块 ============

/** 登录请求 */
export interface LoginRequest {
  username: string;
  password: string;
}

/** 登录响应 */
export interface LoginResponse {
  username: string;
  roles: string[];
  accessToken: string;
  refreshToken: string;
  expires: string;
}

/** 刷新Token响应 */
export interface RefreshTokenResponse {
  accessToken: string;
  refreshToken: string;
  expires: string;
}

// ============ 系统监控模块 ============

/** 服务器信息 */
export interface ServerInfo {
  general?: {
    hostname?: string;
    os?: string;
    arch?: string;
    kernel?: string;
    uptime?: number;
    uptimeFormat?: string;
    cpuCores?: number;
    goVersion?: string;
  };
  cpu?: {
    usagePercent?: number;
    cores?: number;
    modelName?: string;
    frequency?: number;
    coreDetails?: number[];
  };
  memory?: {
    total?: number;
    used?: number;
    free?: number;
    available?: number;
    usedPercent?: number;
    totalFormat?: string;
    usedFormat?: string;
    freeFormat?: string;
    swapTotal?: number;
    swapUsed?: number;
    swapFree?: number;
    swapUsedPercent?: number;
    swapTotalFormat?: string;
    swapUsedFormat?: string;
  };
  disk?: DiskInfo[];
  network?: NetworkInfo[];
  load?: {
    load1?: number;
    load5?: number;
    load15?: number;
  };
  timestamp?: number;
}

/** 磁盘信息 */
export interface DiskInfo {
  device: string;
  mountpoint: string;
  fstype: string;
  total: number;
  used: number;
  free: number;
  usedPercent: number;
  totalFormat: string;
  usedFormat: string;
  freeFormat: string;
}

/** 网络接口信息 */
export interface NetworkInfo {
  name: string;
  bytesRecv: number;
  bytesSent: number;
  packetsRecv: number;
  packetsSent: number;
  recvFormat: string;
  sentFormat: string;
  isUp: boolean;
}

// ============ 网站设置模块 ============

/** 网站设置 */
export interface WebsiteSetting {
  id?: number;
  title: string;
  url: string;
  notice?: string;
  description?: string;
  stat?: string;
  site?: string;
  icp?: string;
  copyright?: string;
  keywords?: string;
  github?: string;
  gitee?: string;
  siteStartDate?: string;
  email?: string;
}
