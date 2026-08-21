-- +goose NO TRANSACTION

-- +goose Up
CREATE TABLE IF NOT EXISTS post (
    post_id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT DEFAULT '',
    post_slug TEXT DEFAULT '',
    author TEXT DEFAULT '',
    cover_image TEXT DEFAULT '',
    post_content TEXT DEFAULT '',
    post_content_html TEXT DEFAULT '',
    summary TEXT DEFAULT '',
    type INTEGER DEFAULT 1,
    top INTEGER DEFAULT 0,
    read_count INTEGER DEFAULT 0,
    word_count INTEGER DEFAULT 0,
    is_published INTEGER DEFAULT 0,
    is_deleted INTEGER DEFAULT 0,
    status INTEGER DEFAULT 0,
    create_time INTEGER,
    pub_time INTEGER,
    last_modified_time INTEGER,
    category_id INTEGER DEFAULT 0
);

CREATE TABLE IF NOT EXISTS category (
    category_id INTEGER PRIMARY KEY AUTOINCREMENT,
    create_time INTEGER,
    category_name TEXT DEFAULT '',
    note TEXT DEFAULT '',
    state INTEGER DEFAULT 1
);

CREATE TABLE IF NOT EXISTS friend_link (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT DEFAULT '',
    link_url TEXT DEFAULT '',
    link_icon TEXT DEFAULT '',
    state INTEGER DEFAULT 0,
    type INTEGER,
    link_desc TEXT DEFAULT '',
    create_time INTEGER,
    last_modified_time INTEGER
);

CREATE TABLE IF NOT EXISTS tag (
    tag_id INTEGER PRIMARY KEY AUTOINCREMENT,
    create_time INTEGER,
    tag_name TEXT
);

CREATE TABLE IF NOT EXISTS post_tag (
    post_id INTEGER,
    tag_id INTEGER
);

CREATE TABLE IF NOT EXISTS about (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT DEFAULT '',
    note TEXT DEFAULT '',
    create_time INTEGER,
    last_modified_time INTEGER
);

CREATE TABLE IF NOT EXISTS user (
    user_id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT DEFAULT '',
    password TEXT DEFAULT '',
    nick_name TEXT DEFAULT '',
    email TEXT DEFAULT '',
    phonenumber TEXT DEFAULT '',
    sex INTEGER DEFAULT 0,
    avatar TEXT DEFAULT ''
);

CREATE TABLE IF NOT EXISTS music (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT DEFAULT '',
    artist TEXT DEFAULT '',
    url TEXT DEFAULT '',
    cover TEXT DEFAULT '',
    lrc TEXT DEFAULT '',
    sort INTEGER DEFAULT 0,
    state INTEGER DEFAULT 1,
    create_time INTEGER
);

CREATE TABLE IF NOT EXISTS system_access_log (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    ip INTEGER DEFAULT 0,
    url TEXT DEFAULT '',
    pv INTEGER DEFAULT 0,
    uv INTEGER DEFAULT 0,
    ua TEXT DEFAULT '',
    status INTEGER DEFAULT 0,
    referer TEXT DEFAULT '',
    area TEXT DEFAULT '',
    create_at INTEGER,
    updated_at INTEGER
);

CREATE TABLE IF NOT EXISTS web_site (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    icp TEXT,
    notice TEXT,
    title TEXT,
    description TEXT,
    url TEXT,
    keywords TEXT,
    copyright TEXT,
    baidu_stat TEXT,
    baidu_site TEXT,
    github TEXT,
    gitee TEXT,
    email TEXT,
    site_start_date TEXT,
    create_at INTEGER,
    updated_at INTEGER
);

CREATE TABLE IF NOT EXISTS system_login_log (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    ip INTEGER DEFAULT 0,
    ua TEXT DEFAULT '',
    note TEXT,
    referer TEXT DEFAULT '',
    area TEXT DEFAULT '',
    success INTEGER DEFAULT 1,
    create_at INTEGER,
    updated_at INTEGER
);

CREATE TABLE IF NOT EXISTS system_notice (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER,
    type INTEGER DEFAULT 1,
    title TEXT,
    description TEXT,
    avatar TEXT,
    extra TEXT,
    extra_status TEXT,
    status INTEGER DEFAULT 0,
    link TEXT,
    created_at DATETIME,
    updated_at DATETIME
);

CREATE INDEX IF NOT EXISTS idx_post_published ON post (is_published);
CREATE INDEX IF NOT EXISTS idx_post_deleted ON post (is_deleted);
CREATE INDEX IF NOT EXISTS idx_post_category ON post (category_id);
CREATE INDEX IF NOT EXISTS idx_post_status ON post (status);
CREATE INDEX IF NOT EXISTS idx_post_create_time ON post (create_time);
CREATE INDEX IF NOT EXISTS idx_post_pub_time ON post (pub_time);
CREATE INDEX IF NOT EXISTS idx_post_slug ON post (post_slug);
CREATE INDEX IF NOT EXISTS idx_post_list ON post (is_published, is_deleted, create_time DESC);
CREATE INDEX IF NOT EXISTS idx_category_state ON category (state);
CREATE INDEX IF NOT EXISTS idx_tag_name ON tag (tag_name);
CREATE INDEX IF NOT EXISTS idx_post_tag_post ON post_tag (post_id);
CREATE INDEX IF NOT EXISTS idx_post_tag_tag ON post_tag (tag_id);
CREATE INDEX IF NOT EXISTS idx_access_log_ip ON system_access_log (ip);
CREATE INDEX IF NOT EXISTS idx_access_log_url ON system_access_log (url);
CREATE INDEX IF NOT EXISTS idx_access_log_create_time ON system_access_log (create_at);
CREATE INDEX IF NOT EXISTS idx_access_log_status ON system_access_log (status);
CREATE INDEX IF NOT EXISTS idx_login_log_create_time ON system_login_log (create_at);
CREATE INDEX IF NOT EXISTS idx_login_log_success ON system_login_log (success);
CREATE INDEX IF NOT EXISTS idx_friend_link_state ON friend_link (state);
CREATE INDEX IF NOT EXISTS idx_system_notice_user_id ON system_notice (user_id);

-- +goose Down
-- Baseline rollback is intentionally a no-op: dropping production tables would be destructive.
SELECT 1;
